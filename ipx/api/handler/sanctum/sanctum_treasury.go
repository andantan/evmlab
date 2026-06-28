package sanctum

import (
	"encoding/hex"
	"encoding/json"
	"fmt"
	"math/big"
	"net/http"

	"github.com/andantan/evmlab/api/handler"
	"github.com/andantan/evmlab/core"
	"github.com/andantan/evmlab/core/types"
	"github.com/andantan/evmlab/internal/util"
)

// DepositNativeLegacy godoc
// @Summary      Deposit native tokens into Sanctum treasury (legacy)
// @Description  Builds, signs, and broadcasts a depositNative() payable call to the Sanctum contract using a legacy transaction.
// @Tags         sanctum
// @Accept       json
// @Produce      json
// @Param        body  body      SanctumDepositNativeLegacyRequest   true  "DepositNative request"
// @Success      200   {object}  SanctumDepositNativeLegacyResponse
// @Failure      400   {object}  map[string]string
// @Failure      502   {object}  map[string]string
// @Router       /evm/sanctum/treasury/native/deposit/legacy [post]
func (h *SanctumHandler) DepositNativeLegacy(w http.ResponseWriter, r *http.Request) {
	var req SanctumDepositNativeLegacyRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		handler.WriteError(w, http.StatusBadRequest, fmt.Sprintf("invalid request body: %s", err))
		return
	}
	if err := req.ValidateRequest(); err != nil {
		handler.WriteError(w, http.StatusBadRequest, err.Error())
		return
	}

	key, err := h.cfg.KeyByAddress(req.From)
	if err != nil {
		handler.WriteError(w, http.StatusBadRequest, fmt.Sprintf("from: %s", err))
		return
	}
	evmKey, err := core.DeriveKeyFromHex(key.PrivateKey, key.PublicKey, key.Address)
	if err != nil {
		handler.WriteError(w, http.StatusInternalServerError, fmt.Sprintf("failed to derive key: %s", err))
		return
	}

	chainIDHex, err := h.client.ChainID(r.Context())
	if err != nil {
		handler.WriteError(w, http.StatusBadGateway, fmt.Sprintf("failed to get chain id: %s", err))
		return
	}
	chainID, err := util.HexToBigInt(chainIDHex)
	if err != nil {
		handler.WriteError(w, http.StatusInternalServerError, fmt.Sprintf("failed to parse chain id: %s", err))
		return
	}

	nonceHex, err := h.client.GetTransactionCount(r.Context(), req.From, "pending")
	if err != nil {
		handler.WriteError(w, http.StatusBadGateway, fmt.Sprintf("failed to get nonce: %s", err))
		return
	}
	nonce, err := util.HexToUint64(nonceHex)
	if err != nil {
		handler.WriteError(w, http.StatusInternalServerError, fmt.Sprintf("failed to parse nonce: %s", err))
		return
	}

	gasPriceHex, err := h.client.GasPrice(r.Context())
	if err != nil {
		handler.WriteError(w, http.StatusBadGateway, fmt.Sprintf("failed to get gas price: %s", err))
		return
	}
	gasPrice, err := util.HexToBigInt(gasPriceHex)
	if err != nil {
		handler.WriteError(w, http.StatusInternalServerError, fmt.Sprintf("failed to parse gas price: %s", err))
		return
	}

	calldata := core.DepositNativeSanctumCalldata()
	p := map[string]any{
		"from":  req.From,
		"to":    req.SanctumAddr().String(),
		"value": "0x" + req.Val().Text(16),
		"data":  "0x" + hex.EncodeToString(calldata),
	}
	gasEstHex, err := h.client.EstimateGas(r.Context(), p, "latest")
	if err != nil {
		h.writeRevertOrGatewayError(w, err, "failed to estimate gas")
		return
	}
	gasEst, err := util.HexToUint64(gasEstHex)
	if err != nil {
		handler.WriteError(w, http.StatusInternalServerError, fmt.Sprintf("failed to parse gas estimate: %s", err))
		return
	}

	tx := &types.LegacyTx{
		ChainID:  chainID,
		Nonce:    nonce,
		GasPrice: gasPrice,
		GasLimit: gasEst * 12 / 10,
		To:       &req.SanctumAddr().Addr,
		Value:    req.Val(),
		Data:     calldata,
	}

	unsigned, err := core.RLP.EncodeLegacyUnsigned(tx)
	if err != nil {
		handler.WriteError(w, http.StatusInternalServerError, fmt.Sprintf("failed to encode unsigned tx: %s", err))
		return
	}
	txHash := core.Hasher.Hash(unsigned)
	sig, err := core.Signer.Sign(txHash, *evmKey.PrivateKey)
	if err != nil {
		handler.WriteError(w, http.StatusInternalServerError, fmt.Sprintf("failed to sign tx: %s", err))
		return
	}
	signed, err := core.RLP.EncodeLegacySigned(tx, sig)
	if err != nil {
		handler.WriteError(w, http.StatusInternalServerError, fmt.Sprintf("failed to encode signed tx: %s", err))
		return
	}
	if _, err = h.client.SendRawTransaction(r.Context(), "0x"+hex.EncodeToString(signed)); err != nil {
		handler.WriteError(w, http.StatusBadGateway, fmt.Sprintf("failed to send transaction: %s", err))
		return
	}
	handler.WriteJSON(w, http.StatusOK, NewSanctumDepositNativeLegacyResponse(unsigned, signed, core.Hasher.Hash(signed), sig))
}

// DepositNativeEIP1559 godoc
// @Summary      Deposit native tokens into Sanctum treasury (EIP-1559)
// @Description  Builds, signs, and broadcasts a depositNative() payable call to the Sanctum contract using an EIP-1559 transaction.
// @Tags         sanctum
// @Accept       json
// @Produce      json
// @Param        body  body      SanctumDepositNativeEIP1559Request  true  "DepositNative request"
// @Success      200   {object}  SanctumDepositNativeEIP1559Response
// @Failure      400   {object}  map[string]string
// @Failure      502   {object}  map[string]string
// @Router       /evm/sanctum/treasury/native/deposit/eip1559 [post]
func (h *SanctumHandler) DepositNativeEIP1559(w http.ResponseWriter, r *http.Request) {
	var req SanctumDepositNativeEIP1559Request
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		handler.WriteError(w, http.StatusBadRequest, fmt.Sprintf("invalid request body: %s", err))
		return
	}
	if err := req.ValidateRequest(); err != nil {
		handler.WriteError(w, http.StatusBadRequest, err.Error())
		return
	}

	key, err := h.cfg.KeyByAddress(req.From)
	if err != nil {
		handler.WriteError(w, http.StatusBadRequest, fmt.Sprintf("from: %s", err))
		return
	}
	evmKey, err := core.DeriveKeyFromHex(key.PrivateKey, key.PublicKey, key.Address)
	if err != nil {
		handler.WriteError(w, http.StatusInternalServerError, fmt.Sprintf("failed to derive key: %s", err))
		return
	}

	chainIDHex, err := h.client.ChainID(r.Context())
	if err != nil {
		handler.WriteError(w, http.StatusBadGateway, fmt.Sprintf("failed to get chain id: %s", err))
		return
	}
	chainID, err := util.HexToBigInt(chainIDHex)
	if err != nil {
		handler.WriteError(w, http.StatusInternalServerError, fmt.Sprintf("failed to parse chain id: %s", err))
		return
	}

	nonceHex, err := h.client.GetTransactionCount(r.Context(), req.From, "pending")
	if err != nil {
		handler.WriteError(w, http.StatusBadGateway, fmt.Sprintf("failed to get nonce: %s", err))
		return
	}
	nonce, err := util.HexToUint64(nonceHex)
	if err != nil {
		handler.WriteError(w, http.StatusInternalServerError, fmt.Sprintf("failed to parse nonce: %s", err))
		return
	}

	tipCapHex, err := h.client.MaxPriorityFeePerGas(r.Context())
	if err != nil {
		handler.WriteError(w, http.StatusBadGateway, fmt.Sprintf("failed to get tip cap: %s", err))
		return
	}
	tipCap, err := util.HexToBigInt(tipCapHex)
	if err != nil {
		handler.WriteError(w, http.StatusInternalServerError, fmt.Sprintf("failed to parse tip cap: %s", err))
		return
	}

	baseFeeHex, err := h.client.BaseFeePerGas(r.Context())
	if err != nil {
		handler.WriteError(w, http.StatusBadGateway, fmt.Sprintf("failed to get base fee: %s", err))
		return
	}
	baseFee, err := util.HexToBigInt(baseFeeHex)
	if err != nil {
		handler.WriteError(w, http.StatusInternalServerError, fmt.Sprintf("failed to parse base fee: %s", err))
		return
	}
	feeCap := new(big.Int).Add(new(big.Int).Mul(baseFee, big.NewInt(2)), tipCap)

	calldata := core.DepositNativeSanctumCalldata()
	p := map[string]any{
		"from":  req.From,
		"to":    req.SanctumAddr().String(),
		"value": "0x" + req.Val().Text(16),
		"data":  "0x" + hex.EncodeToString(calldata),
	}
	gasEstHex, err := h.client.EstimateGas(r.Context(), p, "latest")
	if err != nil {
		h.writeRevertOrGatewayError(w, err, "failed to estimate gas")
		return
	}
	gasEst, err := util.HexToUint64(gasEstHex)
	if err != nil {
		handler.WriteError(w, http.StatusInternalServerError, fmt.Sprintf("failed to parse gas estimate: %s", err))
		return
	}

	tx := &types.DynamicFeeTx{
		ChainID:   chainID,
		Nonce:     nonce,
		GasTipCap: tipCap,
		GasFeeCap: feeCap,
		GasLimit:  gasEst * 12 / 10,
		To:        &req.SanctumAddr().Addr,
		Value:     req.Val(),
		Data:      calldata,
	}

	unsigned, err := core.RLP.EncodeDynamicFeeUnsigned(tx)
	if err != nil {
		handler.WriteError(w, http.StatusInternalServerError, fmt.Sprintf("failed to encode unsigned tx: %s", err))
		return
	}
	txHash := core.Hasher.Hash(unsigned)
	sig, err := core.Signer.Sign(txHash, *evmKey.PrivateKey)
	if err != nil {
		handler.WriteError(w, http.StatusInternalServerError, fmt.Sprintf("failed to sign tx: %s", err))
		return
	}
	signed, err := core.RLP.EncodeDynamicFeeSigned(tx, sig)
	if err != nil {
		handler.WriteError(w, http.StatusInternalServerError, fmt.Sprintf("failed to encode signed tx: %s", err))
		return
	}
	if _, err = h.client.SendRawTransaction(r.Context(), "0x"+hex.EncodeToString(signed)); err != nil {
		handler.WriteError(w, http.StatusBadGateway, fmt.Sprintf("failed to send transaction: %s", err))
		return
	}
	handler.WriteJSON(w, http.StatusOK, NewSanctumDepositNativeEIP1559Response(unsigned, signed, core.Hasher.Hash(signed), sig))
}

// RequestNativeLegacy godoc
// @Summary      Request native token withdrawal from Sanctum treasury (legacy)
// @Description  Builds, signs, and broadcasts a requestNative(uint256) call to the Sanctum contract using a legacy transaction.
// @Tags         sanctum
// @Accept       json
// @Produce      json
// @Param        body  body      SanctumRequestNativeLegacyRequest   true  "RequestNative request"
// @Success      200   {object}  SanctumRequestNativeLegacyResponse
// @Failure      400   {object}  map[string]string
// @Failure      502   {object}  map[string]string
// @Router       /evm/sanctum/treasury/native/request/legacy [post]
func (h *SanctumHandler) RequestNativeLegacy(w http.ResponseWriter, r *http.Request) {
	var req SanctumRequestNativeLegacyRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		handler.WriteError(w, http.StatusBadRequest, fmt.Sprintf("invalid request body: %s", err))
		return
	}
	if err := req.ValidateRequest(); err != nil {
		handler.WriteError(w, http.StatusBadRequest, err.Error())
		return
	}

	key, err := h.cfg.KeyByAddress(req.From)
	if err != nil {
		handler.WriteError(w, http.StatusBadRequest, fmt.Sprintf("from: %s", err))
		return
	}
	evmKey, err := core.DeriveKeyFromHex(key.PrivateKey, key.PublicKey, key.Address)
	if err != nil {
		handler.WriteError(w, http.StatusInternalServerError, fmt.Sprintf("failed to derive key: %s", err))
		return
	}

	chainIDHex, err := h.client.ChainID(r.Context())
	if err != nil {
		handler.WriteError(w, http.StatusBadGateway, fmt.Sprintf("failed to get chain id: %s", err))
		return
	}
	chainID, err := util.HexToBigInt(chainIDHex)
	if err != nil {
		handler.WriteError(w, http.StatusInternalServerError, fmt.Sprintf("failed to parse chain id: %s", err))
		return
	}

	nonceHex, err := h.client.GetTransactionCount(r.Context(), req.From, "pending")
	if err != nil {
		handler.WriteError(w, http.StatusBadGateway, fmt.Sprintf("failed to get nonce: %s", err))
		return
	}
	nonce, err := util.HexToUint64(nonceHex)
	if err != nil {
		handler.WriteError(w, http.StatusInternalServerError, fmt.Sprintf("failed to parse nonce: %s", err))
		return
	}

	gasPriceHex, err := h.client.GasPrice(r.Context())
	if err != nil {
		handler.WriteError(w, http.StatusBadGateway, fmt.Sprintf("failed to get gas price: %s", err))
		return
	}
	gasPrice, err := util.HexToBigInt(gasPriceHex)
	if err != nil {
		handler.WriteError(w, http.StatusInternalServerError, fmt.Sprintf("failed to parse gas price: %s", err))
		return
	}

	calldata := core.RequestNativeSanctumCalldata(req.Amt())
	p := map[string]any{
		"from": req.From,
		"to":   req.SanctumAddr().String(),
		"data": "0x" + hex.EncodeToString(calldata),
	}
	gasEstHex, err := h.client.EstimateGas(r.Context(), p, "latest")
	if err != nil {
		h.writeRevertOrGatewayError(w, err, "failed to estimate gas")
		return
	}
	gasEst, err := util.HexToUint64(gasEstHex)
	if err != nil {
		handler.WriteError(w, http.StatusInternalServerError, fmt.Sprintf("failed to parse gas estimate: %s", err))
		return
	}

	tx := &types.LegacyTx{
		ChainID:  chainID,
		Nonce:    nonce,
		GasPrice: gasPrice,
		GasLimit: gasEst * 12 / 10,
		To:       &req.SanctumAddr().Addr,
		Value:    big.NewInt(0),
		Data:     calldata,
	}

	unsigned, err := core.RLP.EncodeLegacyUnsigned(tx)
	if err != nil {
		handler.WriteError(w, http.StatusInternalServerError, fmt.Sprintf("failed to encode unsigned tx: %s", err))
		return
	}
	txHash := core.Hasher.Hash(unsigned)
	sig, err := core.Signer.Sign(txHash, *evmKey.PrivateKey)
	if err != nil {
		handler.WriteError(w, http.StatusInternalServerError, fmt.Sprintf("failed to sign tx: %s", err))
		return
	}
	signed, err := core.RLP.EncodeLegacySigned(tx, sig)
	if err != nil {
		handler.WriteError(w, http.StatusInternalServerError, fmt.Sprintf("failed to encode signed tx: %s", err))
		return
	}
	if _, err = h.client.SendRawTransaction(r.Context(), "0x"+hex.EncodeToString(signed)); err != nil {
		handler.WriteError(w, http.StatusBadGateway, fmt.Sprintf("failed to send transaction: %s", err))
		return
	}
	handler.WriteJSON(w, http.StatusOK, NewSanctumRequestNativeLegacyResponse(unsigned, signed, core.Hasher.Hash(signed), sig))
}

// RequestNativeEIP1559 godoc
// @Summary      Request native token withdrawal from Sanctum treasury (EIP-1559)
// @Description  Builds, signs, and broadcasts a requestNative(uint256) call to the Sanctum contract using an EIP-1559 transaction.
// @Tags         sanctum
// @Accept       json
// @Produce      json
// @Param        body  body      SanctumRequestNativeEIP1559Request  true  "RequestNative request"
// @Success      200   {object}  SanctumRequestNativeEIP1559Response
// @Failure      400   {object}  map[string]string
// @Failure      502   {object}  map[string]string
// @Router       /evm/sanctum/treasury/native/request/eip1559 [post]
func (h *SanctumHandler) RequestNativeEIP1559(w http.ResponseWriter, r *http.Request) {
	var req SanctumRequestNativeEIP1559Request
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		handler.WriteError(w, http.StatusBadRequest, fmt.Sprintf("invalid request body: %s", err))
		return
	}
	if err := req.ValidateRequest(); err != nil {
		handler.WriteError(w, http.StatusBadRequest, err.Error())
		return
	}

	key, err := h.cfg.KeyByAddress(req.From)
	if err != nil {
		handler.WriteError(w, http.StatusBadRequest, fmt.Sprintf("from: %s", err))
		return
	}
	evmKey, err := core.DeriveKeyFromHex(key.PrivateKey, key.PublicKey, key.Address)
	if err != nil {
		handler.WriteError(w, http.StatusInternalServerError, fmt.Sprintf("failed to derive key: %s", err))
		return
	}

	chainIDHex, err := h.client.ChainID(r.Context())
	if err != nil {
		handler.WriteError(w, http.StatusBadGateway, fmt.Sprintf("failed to get chain id: %s", err))
		return
	}
	chainID, err := util.HexToBigInt(chainIDHex)
	if err != nil {
		handler.WriteError(w, http.StatusInternalServerError, fmt.Sprintf("failed to parse chain id: %s", err))
		return
	}

	nonceHex, err := h.client.GetTransactionCount(r.Context(), req.From, "pending")
	if err != nil {
		handler.WriteError(w, http.StatusBadGateway, fmt.Sprintf("failed to get nonce: %s", err))
		return
	}
	nonce, err := util.HexToUint64(nonceHex)
	if err != nil {
		handler.WriteError(w, http.StatusInternalServerError, fmt.Sprintf("failed to parse nonce: %s", err))
		return
	}

	tipCapHex, err := h.client.MaxPriorityFeePerGas(r.Context())
	if err != nil {
		handler.WriteError(w, http.StatusBadGateway, fmt.Sprintf("failed to get tip cap: %s", err))
		return
	}
	tipCap, err := util.HexToBigInt(tipCapHex)
	if err != nil {
		handler.WriteError(w, http.StatusInternalServerError, fmt.Sprintf("failed to parse tip cap: %s", err))
		return
	}

	baseFeeHex, err := h.client.BaseFeePerGas(r.Context())
	if err != nil {
		handler.WriteError(w, http.StatusBadGateway, fmt.Sprintf("failed to get base fee: %s", err))
		return
	}
	baseFee, err := util.HexToBigInt(baseFeeHex)
	if err != nil {
		handler.WriteError(w, http.StatusInternalServerError, fmt.Sprintf("failed to parse base fee: %s", err))
		return
	}
	feeCap := new(big.Int).Add(new(big.Int).Mul(baseFee, big.NewInt(2)), tipCap)

	calldata := core.RequestNativeSanctumCalldata(req.Amt())
	p := map[string]any{
		"from": req.From,
		"to":   req.SanctumAddr().String(),
		"data": "0x" + hex.EncodeToString(calldata),
	}
	gasEstHex, err := h.client.EstimateGas(r.Context(), p, "latest")
	if err != nil {
		h.writeRevertOrGatewayError(w, err, "failed to estimate gas")
		return
	}
	gasEst, err := util.HexToUint64(gasEstHex)
	if err != nil {
		handler.WriteError(w, http.StatusInternalServerError, fmt.Sprintf("failed to parse gas estimate: %s", err))
		return
	}

	tx := &types.DynamicFeeTx{
		ChainID:   chainID,
		Nonce:     nonce,
		GasTipCap: tipCap,
		GasFeeCap: feeCap,
		GasLimit:  gasEst * 12 / 10,
		To:        &req.SanctumAddr().Addr,
		Value:     big.NewInt(0),
		Data:      calldata,
	}

	unsigned, err := core.RLP.EncodeDynamicFeeUnsigned(tx)
	if err != nil {
		handler.WriteError(w, http.StatusInternalServerError, fmt.Sprintf("failed to encode unsigned tx: %s", err))
		return
	}
	txHash := core.Hasher.Hash(unsigned)
	sig, err := core.Signer.Sign(txHash, *evmKey.PrivateKey)
	if err != nil {
		handler.WriteError(w, http.StatusInternalServerError, fmt.Sprintf("failed to sign tx: %s", err))
		return
	}
	signed, err := core.RLP.EncodeDynamicFeeSigned(tx, sig)
	if err != nil {
		handler.WriteError(w, http.StatusInternalServerError, fmt.Sprintf("failed to encode signed tx: %s", err))
		return
	}
	if _, err = h.client.SendRawTransaction(r.Context(), "0x"+hex.EncodeToString(signed)); err != nil {
		handler.WriteError(w, http.StatusBadGateway, fmt.Sprintf("failed to send transaction: %s", err))
		return
	}
	handler.WriteJSON(w, http.StatusOK, NewSanctumRequestNativeEIP1559Response(unsigned, signed, core.Hasher.Hash(signed), sig))
}

// ApproveNativeLegacy godoc
// @Summary      Approve a user's native token withdrawal request (legacy)
// @Description  Builds, signs, and broadcasts an approveNative(address,uint256) call to the Sanctum contract using a legacy transaction.
// @Tags         sanctum
// @Accept       json
// @Produce      json
// @Param        body  body      SanctumApproveNativeLegacyRequest   true  "ApproveNative request"
// @Success      200   {object}  SanctumApproveNativeLegacyResponse
// @Failure      400   {object}  map[string]string
// @Failure      502   {object}  map[string]string
// @Router       /evm/sanctum/treasury/native/approve/legacy [post]
func (h *SanctumHandler) ApproveNativeLegacy(w http.ResponseWriter, r *http.Request) {
	var req SanctumApproveNativeLegacyRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		handler.WriteError(w, http.StatusBadRequest, fmt.Sprintf("invalid request body: %s", err))
		return
	}
	if err := req.ValidateRequest(); err != nil {
		handler.WriteError(w, http.StatusBadRequest, err.Error())
		return
	}

	key, err := h.cfg.KeyByAddress(req.From)
	if err != nil {
		handler.WriteError(w, http.StatusBadRequest, fmt.Sprintf("from: %s", err))
		return
	}
	evmKey, err := core.DeriveKeyFromHex(key.PrivateKey, key.PublicKey, key.Address)
	if err != nil {
		handler.WriteError(w, http.StatusInternalServerError, fmt.Sprintf("failed to derive key: %s", err))
		return
	}

	chainIDHex, err := h.client.ChainID(r.Context())
	if err != nil {
		handler.WriteError(w, http.StatusBadGateway, fmt.Sprintf("failed to get chain id: %s", err))
		return
	}
	chainID, err := util.HexToBigInt(chainIDHex)
	if err != nil {
		handler.WriteError(w, http.StatusInternalServerError, fmt.Sprintf("failed to parse chain id: %s", err))
		return
	}

	nonceHex, err := h.client.GetTransactionCount(r.Context(), req.From, "pending")
	if err != nil {
		handler.WriteError(w, http.StatusBadGateway, fmt.Sprintf("failed to get nonce: %s", err))
		return
	}
	nonce, err := util.HexToUint64(nonceHex)
	if err != nil {
		handler.WriteError(w, http.StatusInternalServerError, fmt.Sprintf("failed to parse nonce: %s", err))
		return
	}

	gasPriceHex, err := h.client.GasPrice(r.Context())
	if err != nil {
		handler.WriteError(w, http.StatusBadGateway, fmt.Sprintf("failed to get gas price: %s", err))
		return
	}
	gasPrice, err := util.HexToBigInt(gasPriceHex)
	if err != nil {
		handler.WriteError(w, http.StatusInternalServerError, fmt.Sprintf("failed to parse gas price: %s", err))
		return
	}

	calldata := core.ApproveNativeSanctumCalldata(req.UserAddr(), req.Amt())
	p := map[string]any{
		"from": req.From,
		"to":   req.SanctumAddr().String(),
		"data": "0x" + hex.EncodeToString(calldata),
	}
	gasEstHex, err := h.client.EstimateGas(r.Context(), p, "latest")
	if err != nil {
		h.writeRevertOrGatewayError(w, err, "failed to estimate gas")
		return
	}
	gasEst, err := util.HexToUint64(gasEstHex)
	if err != nil {
		handler.WriteError(w, http.StatusInternalServerError, fmt.Sprintf("failed to parse gas estimate: %s", err))
		return
	}

	tx := &types.LegacyTx{
		ChainID:  chainID,
		Nonce:    nonce,
		GasPrice: gasPrice,
		GasLimit: gasEst * 12 / 10,
		To:       &req.SanctumAddr().Addr,
		Value:    big.NewInt(0),
		Data:     calldata,
	}

	unsigned, err := core.RLP.EncodeLegacyUnsigned(tx)
	if err != nil {
		handler.WriteError(w, http.StatusInternalServerError, fmt.Sprintf("failed to encode unsigned tx: %s", err))
		return
	}
	txHash := core.Hasher.Hash(unsigned)
	sig, err := core.Signer.Sign(txHash, *evmKey.PrivateKey)
	if err != nil {
		handler.WriteError(w, http.StatusInternalServerError, fmt.Sprintf("failed to sign tx: %s", err))
		return
	}
	signed, err := core.RLP.EncodeLegacySigned(tx, sig)
	if err != nil {
		handler.WriteError(w, http.StatusInternalServerError, fmt.Sprintf("failed to encode signed tx: %s", err))
		return
	}
	if _, err = h.client.SendRawTransaction(r.Context(), "0x"+hex.EncodeToString(signed)); err != nil {
		handler.WriteError(w, http.StatusBadGateway, fmt.Sprintf("failed to send transaction: %s", err))
		return
	}
	handler.WriteJSON(w, http.StatusOK, NewSanctumApproveNativeLegacyResponse(unsigned, signed, core.Hasher.Hash(signed), sig))
}

// ApproveNativeEIP1559 godoc
// @Summary      Approve a user's native token withdrawal request (EIP-1559)
// @Description  Builds, signs, and broadcasts an approveNative(address,uint256) call to the Sanctum contract using an EIP-1559 transaction.
// @Tags         sanctum
// @Accept       json
// @Produce      json
// @Param        body  body      SanctumApproveNativeEIP1559Request  true  "ApproveNative request"
// @Success      200   {object}  SanctumApproveNativeEIP1559Response
// @Failure      400   {object}  map[string]string
// @Failure      502   {object}  map[string]string
// @Router       /evm/sanctum/treasury/native/approve/eip1559 [post]
func (h *SanctumHandler) ApproveNativeEIP1559(w http.ResponseWriter, r *http.Request) {
	var req SanctumApproveNativeEIP1559Request
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		handler.WriteError(w, http.StatusBadRequest, fmt.Sprintf("invalid request body: %s", err))
		return
	}
	if err := req.ValidateRequest(); err != nil {
		handler.WriteError(w, http.StatusBadRequest, err.Error())
		return
	}

	key, err := h.cfg.KeyByAddress(req.From)
	if err != nil {
		handler.WriteError(w, http.StatusBadRequest, fmt.Sprintf("from: %s", err))
		return
	}
	evmKey, err := core.DeriveKeyFromHex(key.PrivateKey, key.PublicKey, key.Address)
	if err != nil {
		handler.WriteError(w, http.StatusInternalServerError, fmt.Sprintf("failed to derive key: %s", err))
		return
	}

	chainIDHex, err := h.client.ChainID(r.Context())
	if err != nil {
		handler.WriteError(w, http.StatusBadGateway, fmt.Sprintf("failed to get chain id: %s", err))
		return
	}
	chainID, err := util.HexToBigInt(chainIDHex)
	if err != nil {
		handler.WriteError(w, http.StatusInternalServerError, fmt.Sprintf("failed to parse chain id: %s", err))
		return
	}

	nonceHex, err := h.client.GetTransactionCount(r.Context(), req.From, "pending")
	if err != nil {
		handler.WriteError(w, http.StatusBadGateway, fmt.Sprintf("failed to get nonce: %s", err))
		return
	}
	nonce, err := util.HexToUint64(nonceHex)
	if err != nil {
		handler.WriteError(w, http.StatusInternalServerError, fmt.Sprintf("failed to parse nonce: %s", err))
		return
	}

	tipCapHex, err := h.client.MaxPriorityFeePerGas(r.Context())
	if err != nil {
		handler.WriteError(w, http.StatusBadGateway, fmt.Sprintf("failed to get tip cap: %s", err))
		return
	}
	tipCap, err := util.HexToBigInt(tipCapHex)
	if err != nil {
		handler.WriteError(w, http.StatusInternalServerError, fmt.Sprintf("failed to parse tip cap: %s", err))
		return
	}

	baseFeeHex, err := h.client.BaseFeePerGas(r.Context())
	if err != nil {
		handler.WriteError(w, http.StatusBadGateway, fmt.Sprintf("failed to get base fee: %s", err))
		return
	}
	baseFee, err := util.HexToBigInt(baseFeeHex)
	if err != nil {
		handler.WriteError(w, http.StatusInternalServerError, fmt.Sprintf("failed to parse base fee: %s", err))
		return
	}
	feeCap := new(big.Int).Add(new(big.Int).Mul(baseFee, big.NewInt(2)), tipCap)

	calldata := core.ApproveNativeSanctumCalldata(req.UserAddr(), req.Amt())
	p := map[string]any{
		"from": req.From,
		"to":   req.SanctumAddr().String(),
		"data": "0x" + hex.EncodeToString(calldata),
	}
	gasEstHex, err := h.client.EstimateGas(r.Context(), p, "latest")
	if err != nil {
		h.writeRevertOrGatewayError(w, err, "failed to estimate gas")
		return
	}
	gasEst, err := util.HexToUint64(gasEstHex)
	if err != nil {
		handler.WriteError(w, http.StatusInternalServerError, fmt.Sprintf("failed to parse gas estimate: %s", err))
		return
	}

	tx := &types.DynamicFeeTx{
		ChainID:   chainID,
		Nonce:     nonce,
		GasTipCap: tipCap,
		GasFeeCap: feeCap,
		GasLimit:  gasEst * 12 / 10,
		To:        &req.SanctumAddr().Addr,
		Value:     big.NewInt(0),
		Data:      calldata,
	}

	unsigned, err := core.RLP.EncodeDynamicFeeUnsigned(tx)
	if err != nil {
		handler.WriteError(w, http.StatusInternalServerError, fmt.Sprintf("failed to encode unsigned tx: %s", err))
		return
	}
	txHash := core.Hasher.Hash(unsigned)
	sig, err := core.Signer.Sign(txHash, *evmKey.PrivateKey)
	if err != nil {
		handler.WriteError(w, http.StatusInternalServerError, fmt.Sprintf("failed to sign tx: %s", err))
		return
	}
	signed, err := core.RLP.EncodeDynamicFeeSigned(tx, sig)
	if err != nil {
		handler.WriteError(w, http.StatusInternalServerError, fmt.Sprintf("failed to encode signed tx: %s", err))
		return
	}
	if _, err = h.client.SendRawTransaction(r.Context(), "0x"+hex.EncodeToString(signed)); err != nil {
		handler.WriteError(w, http.StatusBadGateway, fmt.Sprintf("failed to send transaction: %s", err))
		return
	}
	handler.WriteJSON(w, http.StatusOK, NewSanctumApproveNativeEIP1559Response(unsigned, signed, core.Hasher.Hash(signed), sig))
}

// ApproveNativeAllLegacy godoc
// @Summary      Approve all pending native token withdrawal for a user (legacy)
// @Description  Builds, signs, and broadcasts an approveNativeAll(address) call to the Sanctum contract using a legacy transaction.
// @Tags         sanctum
// @Accept       json
// @Produce      json
// @Param        body  body      SanctumApproveNativeAllLegacyRequest   true  "ApproveNativeAll request"
// @Success      200   {object}  SanctumApproveNativeAllLegacyResponse
// @Failure      400   {object}  map[string]string
// @Failure      502   {object}  map[string]string
// @Router       /evm/sanctum/treasury/native/approve/all/legacy [post]
func (h *SanctumHandler) ApproveNativeAllLegacy(w http.ResponseWriter, r *http.Request) {
	var req SanctumApproveNativeAllLegacyRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		handler.WriteError(w, http.StatusBadRequest, fmt.Sprintf("invalid request body: %s", err))
		return
	}
	if err := req.ValidateRequest(); err != nil {
		handler.WriteError(w, http.StatusBadRequest, err.Error())
		return
	}

	key, err := h.cfg.KeyByAddress(req.From)
	if err != nil {
		handler.WriteError(w, http.StatusBadRequest, fmt.Sprintf("from: %s", err))
		return
	}
	evmKey, err := core.DeriveKeyFromHex(key.PrivateKey, key.PublicKey, key.Address)
	if err != nil {
		handler.WriteError(w, http.StatusInternalServerError, fmt.Sprintf("failed to derive key: %s", err))
		return
	}

	chainIDHex, err := h.client.ChainID(r.Context())
	if err != nil {
		handler.WriteError(w, http.StatusBadGateway, fmt.Sprintf("failed to get chain id: %s", err))
		return
	}
	chainID, err := util.HexToBigInt(chainIDHex)
	if err != nil {
		handler.WriteError(w, http.StatusInternalServerError, fmt.Sprintf("failed to parse chain id: %s", err))
		return
	}

	nonceHex, err := h.client.GetTransactionCount(r.Context(), req.From, "pending")
	if err != nil {
		handler.WriteError(w, http.StatusBadGateway, fmt.Sprintf("failed to get nonce: %s", err))
		return
	}
	nonce, err := util.HexToUint64(nonceHex)
	if err != nil {
		handler.WriteError(w, http.StatusInternalServerError, fmt.Sprintf("failed to parse nonce: %s", err))
		return
	}

	gasPriceHex, err := h.client.GasPrice(r.Context())
	if err != nil {
		handler.WriteError(w, http.StatusBadGateway, fmt.Sprintf("failed to get gas price: %s", err))
		return
	}
	gasPrice, err := util.HexToBigInt(gasPriceHex)
	if err != nil {
		handler.WriteError(w, http.StatusInternalServerError, fmt.Sprintf("failed to parse gas price: %s", err))
		return
	}

	calldata := core.ApproveNativeAllSanctumCalldata(req.UserAddr())
	p := map[string]any{
		"from": req.From,
		"to":   req.SanctumAddr().String(),
		"data": "0x" + hex.EncodeToString(calldata),
	}
	gasEstHex, err := h.client.EstimateGas(r.Context(), p, "latest")
	if err != nil {
		h.writeRevertOrGatewayError(w, err, "failed to estimate gas")
		return
	}
	gasEst, err := util.HexToUint64(gasEstHex)
	if err != nil {
		handler.WriteError(w, http.StatusInternalServerError, fmt.Sprintf("failed to parse gas estimate: %s", err))
		return
	}

	tx := &types.LegacyTx{
		ChainID:  chainID,
		Nonce:    nonce,
		GasPrice: gasPrice,
		GasLimit: gasEst * 12 / 10,
		To:       &req.SanctumAddr().Addr,
		Value:    big.NewInt(0),
		Data:     calldata,
	}

	unsigned, err := core.RLP.EncodeLegacyUnsigned(tx)
	if err != nil {
		handler.WriteError(w, http.StatusInternalServerError, fmt.Sprintf("failed to encode unsigned tx: %s", err))
		return
	}
	txHash := core.Hasher.Hash(unsigned)
	sig, err := core.Signer.Sign(txHash, *evmKey.PrivateKey)
	if err != nil {
		handler.WriteError(w, http.StatusInternalServerError, fmt.Sprintf("failed to sign tx: %s", err))
		return
	}
	signed, err := core.RLP.EncodeLegacySigned(tx, sig)
	if err != nil {
		handler.WriteError(w, http.StatusInternalServerError, fmt.Sprintf("failed to encode signed tx: %s", err))
		return
	}
	if _, err = h.client.SendRawTransaction(r.Context(), "0x"+hex.EncodeToString(signed)); err != nil {
		handler.WriteError(w, http.StatusBadGateway, fmt.Sprintf("failed to send transaction: %s", err))
		return
	}
	handler.WriteJSON(w, http.StatusOK, NewSanctumApproveNativeAllLegacyResponse(unsigned, signed, core.Hasher.Hash(signed), sig))
}

// ApproveNativeAllEIP1559 godoc
// @Summary      Approve all pending native token withdrawal for a user (EIP-1559)
// @Description  Builds, signs, and broadcasts an approveNativeAll(address) call to the Sanctum contract using an EIP-1559 transaction.
// @Tags         sanctum
// @Accept       json
// @Produce      json
// @Param        body  body      SanctumApproveNativeAllEIP1559Request  true  "ApproveNativeAll request"
// @Success      200   {object}  SanctumApproveNativeAllEIP1559Response
// @Failure      400   {object}  map[string]string
// @Failure      502   {object}  map[string]string
// @Router       /evm/sanctum/treasury/native/approve/all/eip1559 [post]
func (h *SanctumHandler) ApproveNativeAllEIP1559(w http.ResponseWriter, r *http.Request) {
	var req SanctumApproveNativeAllEIP1559Request
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		handler.WriteError(w, http.StatusBadRequest, fmt.Sprintf("invalid request body: %s", err))
		return
	}
	if err := req.ValidateRequest(); err != nil {
		handler.WriteError(w, http.StatusBadRequest, err.Error())
		return
	}

	key, err := h.cfg.KeyByAddress(req.From)
	if err != nil {
		handler.WriteError(w, http.StatusBadRequest, fmt.Sprintf("from: %s", err))
		return
	}
	evmKey, err := core.DeriveKeyFromHex(key.PrivateKey, key.PublicKey, key.Address)
	if err != nil {
		handler.WriteError(w, http.StatusInternalServerError, fmt.Sprintf("failed to derive key: %s", err))
		return
	}

	chainIDHex, err := h.client.ChainID(r.Context())
	if err != nil {
		handler.WriteError(w, http.StatusBadGateway, fmt.Sprintf("failed to get chain id: %s", err))
		return
	}
	chainID, err := util.HexToBigInt(chainIDHex)
	if err != nil {
		handler.WriteError(w, http.StatusInternalServerError, fmt.Sprintf("failed to parse chain id: %s", err))
		return
	}

	nonceHex, err := h.client.GetTransactionCount(r.Context(), req.From, "pending")
	if err != nil {
		handler.WriteError(w, http.StatusBadGateway, fmt.Sprintf("failed to get nonce: %s", err))
		return
	}
	nonce, err := util.HexToUint64(nonceHex)
	if err != nil {
		handler.WriteError(w, http.StatusInternalServerError, fmt.Sprintf("failed to parse nonce: %s", err))
		return
	}

	tipCapHex, err := h.client.MaxPriorityFeePerGas(r.Context())
	if err != nil {
		handler.WriteError(w, http.StatusBadGateway, fmt.Sprintf("failed to get tip cap: %s", err))
		return
	}
	tipCap, err := util.HexToBigInt(tipCapHex)
	if err != nil {
		handler.WriteError(w, http.StatusInternalServerError, fmt.Sprintf("failed to parse tip cap: %s", err))
		return
	}

	baseFeeHex, err := h.client.BaseFeePerGas(r.Context())
	if err != nil {
		handler.WriteError(w, http.StatusBadGateway, fmt.Sprintf("failed to get base fee: %s", err))
		return
	}
	baseFee, err := util.HexToBigInt(baseFeeHex)
	if err != nil {
		handler.WriteError(w, http.StatusInternalServerError, fmt.Sprintf("failed to parse base fee: %s", err))
		return
	}
	feeCap := new(big.Int).Add(new(big.Int).Mul(baseFee, big.NewInt(2)), tipCap)

	calldata := core.ApproveNativeAllSanctumCalldata(req.UserAddr())
	p := map[string]any{
		"from": req.From,
		"to":   req.SanctumAddr().String(),
		"data": "0x" + hex.EncodeToString(calldata),
	}
	gasEstHex, err := h.client.EstimateGas(r.Context(), p, "latest")
	if err != nil {
		h.writeRevertOrGatewayError(w, err, "failed to estimate gas")
		return
	}
	gasEst, err := util.HexToUint64(gasEstHex)
	if err != nil {
		handler.WriteError(w, http.StatusInternalServerError, fmt.Sprintf("failed to parse gas estimate: %s", err))
		return
	}

	tx := &types.DynamicFeeTx{
		ChainID:   chainID,
		Nonce:     nonce,
		GasTipCap: tipCap,
		GasFeeCap: feeCap,
		GasLimit:  gasEst * 12 / 10,
		To:        &req.SanctumAddr().Addr,
		Value:     big.NewInt(0),
		Data:      calldata,
	}

	unsigned, err := core.RLP.EncodeDynamicFeeUnsigned(tx)
	if err != nil {
		handler.WriteError(w, http.StatusInternalServerError, fmt.Sprintf("failed to encode unsigned tx: %s", err))
		return
	}
	txHash := core.Hasher.Hash(unsigned)
	sig, err := core.Signer.Sign(txHash, *evmKey.PrivateKey)
	if err != nil {
		handler.WriteError(w, http.StatusInternalServerError, fmt.Sprintf("failed to sign tx: %s", err))
		return
	}
	signed, err := core.RLP.EncodeDynamicFeeSigned(tx, sig)
	if err != nil {
		handler.WriteError(w, http.StatusInternalServerError, fmt.Sprintf("failed to encode signed tx: %s", err))
		return
	}
	if _, err = h.client.SendRawTransaction(r.Context(), "0x"+hex.EncodeToString(signed)); err != nil {
		handler.WriteError(w, http.StatusBadGateway, fmt.Sprintf("failed to send transaction: %s", err))
		return
	}
	handler.WriteJSON(w, http.StatusOK, NewSanctumApproveNativeAllEIP1559Response(unsigned, signed, core.Hasher.Hash(signed), sig))
}

// WithdrawNativeLegacy godoc
// @Summary      Withdraw native tokens from Sanctum treasury (legacy)
// @Description  Builds, signs, and broadcasts a withdrawNative(uint256) call to the Sanctum contract using a legacy transaction.
// @Tags         sanctum
// @Accept       json
// @Produce      json
// @Param        body  body      SanctumWithdrawNativeLegacyRequest   true  "WithdrawNative request"
// @Success      200   {object}  SanctumWithdrawNativeLegacyResponse
// @Failure      400   {object}  map[string]string
// @Failure      502   {object}  map[string]string
// @Router       /evm/sanctum/treasury/native/withdraw/legacy [post]
func (h *SanctumHandler) WithdrawNativeLegacy(w http.ResponseWriter, r *http.Request) {
	var req SanctumWithdrawNativeLegacyRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		handler.WriteError(w, http.StatusBadRequest, fmt.Sprintf("invalid request body: %s", err))
		return
	}
	if err := req.ValidateRequest(); err != nil {
		handler.WriteError(w, http.StatusBadRequest, err.Error())
		return
	}

	key, err := h.cfg.KeyByAddress(req.From)
	if err != nil {
		handler.WriteError(w, http.StatusBadRequest, fmt.Sprintf("from: %s", err))
		return
	}
	evmKey, err := core.DeriveKeyFromHex(key.PrivateKey, key.PublicKey, key.Address)
	if err != nil {
		handler.WriteError(w, http.StatusInternalServerError, fmt.Sprintf("failed to derive key: %s", err))
		return
	}

	chainIDHex, err := h.client.ChainID(r.Context())
	if err != nil {
		handler.WriteError(w, http.StatusBadGateway, fmt.Sprintf("failed to get chain id: %s", err))
		return
	}
	chainID, err := util.HexToBigInt(chainIDHex)
	if err != nil {
		handler.WriteError(w, http.StatusInternalServerError, fmt.Sprintf("failed to parse chain id: %s", err))
		return
	}

	nonceHex, err := h.client.GetTransactionCount(r.Context(), req.From, "pending")
	if err != nil {
		handler.WriteError(w, http.StatusBadGateway, fmt.Sprintf("failed to get nonce: %s", err))
		return
	}
	nonce, err := util.HexToUint64(nonceHex)
	if err != nil {
		handler.WriteError(w, http.StatusInternalServerError, fmt.Sprintf("failed to parse nonce: %s", err))
		return
	}

	gasPriceHex, err := h.client.GasPrice(r.Context())
	if err != nil {
		handler.WriteError(w, http.StatusBadGateway, fmt.Sprintf("failed to get gas price: %s", err))
		return
	}
	gasPrice, err := util.HexToBigInt(gasPriceHex)
	if err != nil {
		handler.WriteError(w, http.StatusInternalServerError, fmt.Sprintf("failed to parse gas price: %s", err))
		return
	}

	calldata := core.WithdrawNativeSanctumCalldata(req.ToAmount())
	p := map[string]any{
		"from": req.From,
		"to":   req.SanctumAddr().String(),
		"data": "0x" + hex.EncodeToString(calldata),
	}
	gasEstHex, err := h.client.EstimateGas(r.Context(), p, "latest")
	if err != nil {
		h.writeRevertOrGatewayError(w, err, "failed to estimate gas")
		return
	}
	gasEst, err := util.HexToUint64(gasEstHex)
	if err != nil {
		handler.WriteError(w, http.StatusInternalServerError, fmt.Sprintf("failed to parse gas estimate: %s", err))
		return
	}

	tx := &types.LegacyTx{
		ChainID:  chainID,
		Nonce:    nonce,
		GasPrice: gasPrice,
		GasLimit: gasEst * 12 / 10,
		To:       &req.SanctumAddr().Addr,
		Value:    big.NewInt(0),
		Data:     calldata,
	}

	unsigned, err := core.RLP.EncodeLegacyUnsigned(tx)
	if err != nil {
		handler.WriteError(w, http.StatusInternalServerError, fmt.Sprintf("failed to encode unsigned tx: %s", err))
		return
	}
	txHash := core.Hasher.Hash(unsigned)
	sig, err := core.Signer.Sign(txHash, *evmKey.PrivateKey)
	if err != nil {
		handler.WriteError(w, http.StatusInternalServerError, fmt.Sprintf("failed to sign tx: %s", err))
		return
	}
	signed, err := core.RLP.EncodeLegacySigned(tx, sig)
	if err != nil {
		handler.WriteError(w, http.StatusInternalServerError, fmt.Sprintf("failed to encode signed tx: %s", err))
		return
	}
	if _, err = h.client.SendRawTransaction(r.Context(), "0x"+hex.EncodeToString(signed)); err != nil {
		handler.WriteError(w, http.StatusBadGateway, fmt.Sprintf("failed to send transaction: %s", err))
		return
	}
	handler.WriteJSON(w, http.StatusOK, NewSanctumWithdrawNativeLegacyResponse(unsigned, signed, core.Hasher.Hash(signed), sig))
}

// WithdrawNativeEIP1559 godoc
// @Summary      Withdraw native tokens from Sanctum treasury (EIP-1559)
// @Description  Builds, signs, and broadcasts a withdrawNative(uint256) call to the Sanctum contract using an EIP-1559 transaction.
// @Tags         sanctum
// @Accept       json
// @Produce      json
// @Param        body  body      SanctumWithdrawNativeEIP1559Request  true  "WithdrawNative request"
// @Success      200   {object}  SanctumWithdrawNativeEIP1559Response
// @Failure      400   {object}  map[string]string
// @Failure      502   {object}  map[string]string
// @Router       /evm/sanctum/treasury/native/withdraw/eip1559 [post]
func (h *SanctumHandler) WithdrawNativeEIP1559(w http.ResponseWriter, r *http.Request) {
	var req SanctumWithdrawNativeEIP1559Request
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		handler.WriteError(w, http.StatusBadRequest, fmt.Sprintf("invalid request body: %s", err))
		return
	}
	if err := req.ValidateRequest(); err != nil {
		handler.WriteError(w, http.StatusBadRequest, err.Error())
		return
	}

	key, err := h.cfg.KeyByAddress(req.From)
	if err != nil {
		handler.WriteError(w, http.StatusBadRequest, fmt.Sprintf("from: %s", err))
		return
	}
	evmKey, err := core.DeriveKeyFromHex(key.PrivateKey, key.PublicKey, key.Address)
	if err != nil {
		handler.WriteError(w, http.StatusInternalServerError, fmt.Sprintf("failed to derive key: %s", err))
		return
	}

	chainIDHex, err := h.client.ChainID(r.Context())
	if err != nil {
		handler.WriteError(w, http.StatusBadGateway, fmt.Sprintf("failed to get chain id: %s", err))
		return
	}
	chainID, err := util.HexToBigInt(chainIDHex)
	if err != nil {
		handler.WriteError(w, http.StatusInternalServerError, fmt.Sprintf("failed to parse chain id: %s", err))
		return
	}

	nonceHex, err := h.client.GetTransactionCount(r.Context(), req.From, "pending")
	if err != nil {
		handler.WriteError(w, http.StatusBadGateway, fmt.Sprintf("failed to get nonce: %s", err))
		return
	}
	nonce, err := util.HexToUint64(nonceHex)
	if err != nil {
		handler.WriteError(w, http.StatusInternalServerError, fmt.Sprintf("failed to parse nonce: %s", err))
		return
	}

	tipCapHex, err := h.client.MaxPriorityFeePerGas(r.Context())
	if err != nil {
		handler.WriteError(w, http.StatusBadGateway, fmt.Sprintf("failed to get tip cap: %s", err))
		return
	}
	tipCap, err := util.HexToBigInt(tipCapHex)
	if err != nil {
		handler.WriteError(w, http.StatusInternalServerError, fmt.Sprintf("failed to parse tip cap: %s", err))
		return
	}

	baseFeeHex, err := h.client.BaseFeePerGas(r.Context())
	if err != nil {
		handler.WriteError(w, http.StatusBadGateway, fmt.Sprintf("failed to get base fee: %s", err))
		return
	}
	baseFee, err := util.HexToBigInt(baseFeeHex)
	if err != nil {
		handler.WriteError(w, http.StatusInternalServerError, fmt.Sprintf("failed to parse base fee: %s", err))
		return
	}
	feeCap := new(big.Int).Add(new(big.Int).Mul(baseFee, big.NewInt(2)), tipCap)

	calldata := core.WithdrawNativeSanctumCalldata(req.ToAmount())
	p := map[string]any{
		"from": req.From,
		"to":   req.SanctumAddr().String(),
		"data": "0x" + hex.EncodeToString(calldata),
	}
	gasEstHex, err := h.client.EstimateGas(r.Context(), p, "latest")
	if err != nil {
		h.writeRevertOrGatewayError(w, err, "failed to estimate gas")
		return
	}
	gasEst, err := util.HexToUint64(gasEstHex)
	if err != nil {
		handler.WriteError(w, http.StatusInternalServerError, fmt.Sprintf("failed to parse gas estimate: %s", err))
		return
	}

	tx := &types.DynamicFeeTx{
		ChainID:   chainID,
		Nonce:     nonce,
		GasTipCap: tipCap,
		GasFeeCap: feeCap,
		GasLimit:  gasEst * 12 / 10,
		To:        &req.SanctumAddr().Addr,
		Value:     big.NewInt(0),
		Data:      calldata,
	}

	unsigned, err := core.RLP.EncodeDynamicFeeUnsigned(tx)
	if err != nil {
		handler.WriteError(w, http.StatusInternalServerError, fmt.Sprintf("failed to encode unsigned tx: %s", err))
		return
	}
	txHash := core.Hasher.Hash(unsigned)
	sig, err := core.Signer.Sign(txHash, *evmKey.PrivateKey)
	if err != nil {
		handler.WriteError(w, http.StatusInternalServerError, fmt.Sprintf("failed to sign tx: %s", err))
		return
	}
	signed, err := core.RLP.EncodeDynamicFeeSigned(tx, sig)
	if err != nil {
		handler.WriteError(w, http.StatusInternalServerError, fmt.Sprintf("failed to encode signed tx: %s", err))
		return
	}
	if _, err = h.client.SendRawTransaction(r.Context(), "0x"+hex.EncodeToString(signed)); err != nil {
		handler.WriteError(w, http.StatusBadGateway, fmt.Sprintf("failed to send transaction: %s", err))
		return
	}
	handler.WriteJSON(w, http.StatusOK, NewSanctumWithdrawNativeEIP1559Response(unsigned, signed, core.Hasher.Hash(signed), sig))
}

// WithdrawNativeAllLegacy godoc
// @Summary      Withdraw all available native tokens from Sanctum treasury (legacy)
// @Description  Builds, signs, and broadcasts a withdrawNativeAll() call to the Sanctum contract using a legacy transaction.
// @Tags         sanctum
// @Accept       json
// @Produce      json
// @Param        body  body      SanctumWithdrawNativeAllLegacyRequest   true  "WithdrawNativeAll request"
// @Success      200   {object}  SanctumWithdrawNativeAllLegacyResponse
// @Failure      400   {object}  map[string]string
// @Failure      502   {object}  map[string]string
// @Router       /evm/sanctum/treasury/native/withdraw/all/legacy [post]
func (h *SanctumHandler) WithdrawNativeAllLegacy(w http.ResponseWriter, r *http.Request) {
	var req SanctumWithdrawNativeAllLegacyRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		handler.WriteError(w, http.StatusBadRequest, fmt.Sprintf("invalid request body: %s", err))
		return
	}
	if err := req.ValidateRequest(); err != nil {
		handler.WriteError(w, http.StatusBadRequest, err.Error())
		return
	}

	key, err := h.cfg.KeyByAddress(req.From)
	if err != nil {
		handler.WriteError(w, http.StatusBadRequest, fmt.Sprintf("from: %s", err))
		return
	}
	evmKey, err := core.DeriveKeyFromHex(key.PrivateKey, key.PublicKey, key.Address)
	if err != nil {
		handler.WriteError(w, http.StatusInternalServerError, fmt.Sprintf("failed to derive key: %s", err))
		return
	}

	chainIDHex, err := h.client.ChainID(r.Context())
	if err != nil {
		handler.WriteError(w, http.StatusBadGateway, fmt.Sprintf("failed to get chain id: %s", err))
		return
	}
	chainID, err := util.HexToBigInt(chainIDHex)
	if err != nil {
		handler.WriteError(w, http.StatusInternalServerError, fmt.Sprintf("failed to parse chain id: %s", err))
		return
	}

	nonceHex, err := h.client.GetTransactionCount(r.Context(), req.From, "pending")
	if err != nil {
		handler.WriteError(w, http.StatusBadGateway, fmt.Sprintf("failed to get nonce: %s", err))
		return
	}
	nonce, err := util.HexToUint64(nonceHex)
	if err != nil {
		handler.WriteError(w, http.StatusInternalServerError, fmt.Sprintf("failed to parse nonce: %s", err))
		return
	}

	gasPriceHex, err := h.client.GasPrice(r.Context())
	if err != nil {
		handler.WriteError(w, http.StatusBadGateway, fmt.Sprintf("failed to get gas price: %s", err))
		return
	}
	gasPrice, err := util.HexToBigInt(gasPriceHex)
	if err != nil {
		handler.WriteError(w, http.StatusInternalServerError, fmt.Sprintf("failed to parse gas price: %s", err))
		return
	}

	calldata := core.WithdrawNativeAllSanctumCalldata()
	p := map[string]any{
		"from": req.From,
		"to":   req.SanctumAddr().String(),
		"data": "0x" + hex.EncodeToString(calldata),
	}
	gasEstHex, err := h.client.EstimateGas(r.Context(), p, "latest")
	if err != nil {
		h.writeRevertOrGatewayError(w, err, "failed to estimate gas")
		return
	}
	gasEst, err := util.HexToUint64(gasEstHex)
	if err != nil {
		handler.WriteError(w, http.StatusInternalServerError, fmt.Sprintf("failed to parse gas estimate: %s", err))
		return
	}

	tx := &types.LegacyTx{
		ChainID:  chainID,
		Nonce:    nonce,
		GasPrice: gasPrice,
		GasLimit: gasEst * 12 / 10,
		To:       &req.SanctumAddr().Addr,
		Value:    big.NewInt(0),
		Data:     calldata,
	}

	unsigned, err := core.RLP.EncodeLegacyUnsigned(tx)
	if err != nil {
		handler.WriteError(w, http.StatusInternalServerError, fmt.Sprintf("failed to encode unsigned tx: %s", err))
		return
	}
	txHash := core.Hasher.Hash(unsigned)
	sig, err := core.Signer.Sign(txHash, *evmKey.PrivateKey)
	if err != nil {
		handler.WriteError(w, http.StatusInternalServerError, fmt.Sprintf("failed to sign tx: %s", err))
		return
	}
	signed, err := core.RLP.EncodeLegacySigned(tx, sig)
	if err != nil {
		handler.WriteError(w, http.StatusInternalServerError, fmt.Sprintf("failed to encode signed tx: %s", err))
		return
	}
	if _, err = h.client.SendRawTransaction(r.Context(), "0x"+hex.EncodeToString(signed)); err != nil {
		handler.WriteError(w, http.StatusBadGateway, fmt.Sprintf("failed to send transaction: %s", err))
		return
	}
	handler.WriteJSON(w, http.StatusOK, NewSanctumWithdrawNativeAllLegacyResponse(unsigned, signed, core.Hasher.Hash(signed), sig))
}

// WithdrawNativeAllEIP1559 godoc
// @Summary      Withdraw all available native tokens from Sanctum treasury (EIP-1559)
// @Description  Builds, signs, and broadcasts a withdrawNativeAll() call to the Sanctum contract using an EIP-1559 transaction.
// @Tags         sanctum
// @Accept       json
// @Produce      json
// @Param        body  body      SanctumWithdrawNativeAllEIP1559Request  true  "WithdrawNativeAll request"
// @Success      200   {object}  SanctumWithdrawNativeAllEIP1559Response
// @Failure      400   {object}  map[string]string
// @Failure      502   {object}  map[string]string
// @Router       /evm/sanctum/treasury/native/withdraw/all/eip1559 [post]
func (h *SanctumHandler) WithdrawNativeAllEIP1559(w http.ResponseWriter, r *http.Request) {
	var req SanctumWithdrawNativeAllEIP1559Request
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		handler.WriteError(w, http.StatusBadRequest, fmt.Sprintf("invalid request body: %s", err))
		return
	}
	if err := req.ValidateRequest(); err != nil {
		handler.WriteError(w, http.StatusBadRequest, err.Error())
		return
	}

	key, err := h.cfg.KeyByAddress(req.From)
	if err != nil {
		handler.WriteError(w, http.StatusBadRequest, fmt.Sprintf("from: %s", err))
		return
	}
	evmKey, err := core.DeriveKeyFromHex(key.PrivateKey, key.PublicKey, key.Address)
	if err != nil {
		handler.WriteError(w, http.StatusInternalServerError, fmt.Sprintf("failed to derive key: %s", err))
		return
	}

	chainIDHex, err := h.client.ChainID(r.Context())
	if err != nil {
		handler.WriteError(w, http.StatusBadGateway, fmt.Sprintf("failed to get chain id: %s", err))
		return
	}
	chainID, err := util.HexToBigInt(chainIDHex)
	if err != nil {
		handler.WriteError(w, http.StatusInternalServerError, fmt.Sprintf("failed to parse chain id: %s", err))
		return
	}

	nonceHex, err := h.client.GetTransactionCount(r.Context(), req.From, "pending")
	if err != nil {
		handler.WriteError(w, http.StatusBadGateway, fmt.Sprintf("failed to get nonce: %s", err))
		return
	}
	nonce, err := util.HexToUint64(nonceHex)
	if err != nil {
		handler.WriteError(w, http.StatusInternalServerError, fmt.Sprintf("failed to parse nonce: %s", err))
		return
	}

	tipCapHex, err := h.client.MaxPriorityFeePerGas(r.Context())
	if err != nil {
		handler.WriteError(w, http.StatusBadGateway, fmt.Sprintf("failed to get tip cap: %s", err))
		return
	}
	tipCap, err := util.HexToBigInt(tipCapHex)
	if err != nil {
		handler.WriteError(w, http.StatusInternalServerError, fmt.Sprintf("failed to parse tip cap: %s", err))
		return
	}

	baseFeeHex, err := h.client.BaseFeePerGas(r.Context())
	if err != nil {
		handler.WriteError(w, http.StatusBadGateway, fmt.Sprintf("failed to get base fee: %s", err))
		return
	}
	baseFee, err := util.HexToBigInt(baseFeeHex)
	if err != nil {
		handler.WriteError(w, http.StatusInternalServerError, fmt.Sprintf("failed to parse base fee: %s", err))
		return
	}
	feeCap := new(big.Int).Add(new(big.Int).Mul(baseFee, big.NewInt(2)), tipCap)

	calldata := core.WithdrawNativeAllSanctumCalldata()
	p := map[string]any{
		"from": req.From,
		"to":   req.SanctumAddr().String(),
		"data": "0x" + hex.EncodeToString(calldata),
	}
	gasEstHex, err := h.client.EstimateGas(r.Context(), p, "latest")
	if err != nil {
		h.writeRevertOrGatewayError(w, err, "failed to estimate gas")
		return
	}
	gasEst, err := util.HexToUint64(gasEstHex)
	if err != nil {
		handler.WriteError(w, http.StatusInternalServerError, fmt.Sprintf("failed to parse gas estimate: %s", err))
		return
	}

	tx := &types.DynamicFeeTx{
		ChainID:   chainID,
		Nonce:     nonce,
		GasTipCap: tipCap,
		GasFeeCap: feeCap,
		GasLimit:  gasEst * 12 / 10,
		To:        &req.SanctumAddr().Addr,
		Value:     big.NewInt(0),
		Data:      calldata,
	}

	unsigned, err := core.RLP.EncodeDynamicFeeUnsigned(tx)
	if err != nil {
		handler.WriteError(w, http.StatusInternalServerError, fmt.Sprintf("failed to encode unsigned tx: %s", err))
		return
	}
	txHash := core.Hasher.Hash(unsigned)
	sig, err := core.Signer.Sign(txHash, *evmKey.PrivateKey)
	if err != nil {
		handler.WriteError(w, http.StatusInternalServerError, fmt.Sprintf("failed to sign tx: %s", err))
		return
	}
	signed, err := core.RLP.EncodeDynamicFeeSigned(tx, sig)
	if err != nil {
		handler.WriteError(w, http.StatusInternalServerError, fmt.Sprintf("failed to encode signed tx: %s", err))
		return
	}
	if _, err = h.client.SendRawTransaction(r.Context(), "0x"+hex.EncodeToString(signed)); err != nil {
		handler.WriteError(w, http.StatusBadGateway, fmt.Sprintf("failed to send transaction: %s", err))
		return
	}
	handler.WriteJSON(w, http.StatusOK, NewSanctumWithdrawNativeAllEIP1559Response(unsigned, signed, core.Hasher.Hash(signed), sig))
}

// NativeBalance godoc
// @Summary      Get the total native token balance of the Sanctum treasury
// @Description  Calls nativeBalance() on the Sanctum contract.
// @Tags         sanctum
// @Accept       json
// @Produce      json
// @Param        body  body      SanctumNativeBalanceRequest   true  "NativeBalance request"
// @Success      200   {object}  SanctumNativeBalanceResponse
// @Failure      400   {object}  map[string]string
// @Failure      502   {object}  map[string]string
// @Router       /evm/sanctum/treasury/native/balance [post]
func (h *SanctumHandler) NativeBalance(w http.ResponseWriter, r *http.Request) {
	var req SanctumNativeBalanceRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		handler.WriteError(w, http.StatusBadRequest, fmt.Sprintf("invalid request body: %s", err))
		return
	}
	if err := req.ValidateRequest(); err != nil {
		handler.WriteError(w, http.StatusBadRequest, err.Error())
		return
	}

	calldata := core.NativeBalanceSanctumCalldata()
	p := map[string]any{
		"to":   req.SanctumAddr().String(),
		"data": "0x" + hex.EncodeToString(calldata),
	}
	raw, err := h.client.CallContract(r.Context(), p, req.Block)
	if err != nil {
		h.writeRevertOrGatewayError(w, err, "eth_call nativeBalance() failed")
		return
	}
	data, err := util.ParseHex(raw)
	if err != nil {
		handler.WriteError(w, http.StatusInternalServerError, fmt.Sprintf("failed to parse response: %s", err))
		return
	}

	balance, err := types.DecodeUint256(data)
	if err != nil {
		handler.WriteError(w, http.StatusInternalServerError, fmt.Sprintf("failed to decode balance: %s", err))
		return
	}
	handler.WriteJSON(w, http.StatusOK, SanctumNativeBalanceResponse{Balance: balance.String()})
}

// NativeAvailable godoc
// @Summary      Get the available (unallocated) native token balance of the Sanctum treasury
// @Description  Calls nativeAvailable() on the Sanctum contract.
// @Tags         sanctum
// @Accept       json
// @Produce      json
// @Param        body  body      SanctumNativeAvailableRequest   true  "NativeAvailable request"
// @Success      200   {object}  SanctumNativeAvailableResponse
// @Failure      400   {object}  map[string]string
// @Failure      502   {object}  map[string]string
// @Router       /evm/sanctum/treasury/native/available [post]
func (h *SanctumHandler) NativeAvailable(w http.ResponseWriter, r *http.Request) {
	var req SanctumNativeAvailableRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		handler.WriteError(w, http.StatusBadRequest, fmt.Sprintf("invalid request body: %s", err))
		return
	}
	if err := req.ValidateRequest(); err != nil {
		handler.WriteError(w, http.StatusBadRequest, err.Error())
		return
	}

	calldata := core.NativeAvailableSanctumCalldata()
	p := map[string]any{
		"to":   req.SanctumAddr().String(),
		"data": "0x" + hex.EncodeToString(calldata),
	}
	raw, err := h.client.CallContract(r.Context(), p, req.Block)
	if err != nil {
		h.writeRevertOrGatewayError(w, err, "eth_call nativeAvailable() failed")
		return
	}
	data, err := util.ParseHex(raw)
	if err != nil {
		handler.WriteError(w, http.StatusInternalServerError, fmt.Sprintf("failed to parse response: %s", err))
		return
	}

	available, err := types.DecodeUint256(data)
	if err != nil {
		handler.WriteError(w, http.StatusInternalServerError, fmt.Sprintf("failed to decode available: %s", err))
		return
	}
	handler.WriteJSON(w, http.StatusOK, SanctumNativeAvailableResponse{Available: available.String()})
}

// NativeAllocation godoc
// @Summary      Get the allocated native token amount for a user
// @Description  Calls nativeAllocation(address) on the Sanctum contract.
// @Tags         sanctum
// @Accept       json
// @Produce      json
// @Param        body  body      SanctumNativeAllocationRequest   true  "NativeAllocation request"
// @Success      200   {object}  SanctumNativeAllocationResponse
// @Failure      400   {object}  map[string]string
// @Failure      502   {object}  map[string]string
// @Router       /evm/sanctum/treasury/native/allocation [post]
func (h *SanctumHandler) NativeAllocation(w http.ResponseWriter, r *http.Request) {
	var req SanctumNativeAllocationRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		handler.WriteError(w, http.StatusBadRequest, fmt.Sprintf("invalid request body: %s", err))
		return
	}
	if err := req.ValidateRequest(); err != nil {
		handler.WriteError(w, http.StatusBadRequest, err.Error())
		return
	}

	calldata := core.NativeAllocationSanctumCalldata(req.UserAddr())
	p := map[string]any{
		"to":   req.SanctumAddr().String(),
		"data": "0x" + hex.EncodeToString(calldata),
	}
	raw, err := h.client.CallContract(r.Context(), p, req.Block)
	if err != nil {
		h.writeRevertOrGatewayError(w, err, "eth_call nativeAllocation() failed")
		return
	}
	data, err := util.ParseHex(raw)
	if err != nil {
		handler.WriteError(w, http.StatusInternalServerError, fmt.Sprintf("failed to parse response: %s", err))
		return
	}

	allocation, err := types.DecodeUint256(data)
	if err != nil {
		handler.WriteError(w, http.StatusInternalServerError, fmt.Sprintf("failed to decode allocation: %s", err))
		return
	}
	handler.WriteJSON(w, http.StatusOK, SanctumNativeAllocationResponse{Allocation: allocation.String()})
}

// NativePending godoc
// @Summary      Get the pending native token withdrawal amount for a user
// @Description  Calls nativePending(address) on the Sanctum contract.
// @Tags         sanctum
// @Accept       json
// @Produce      json
// @Param        body  body      SanctumNativePendingRequest   true  "NativePending request"
// @Success      200   {object}  SanctumNativePendingResponse
// @Failure      400   {object}  map[string]string
// @Failure      502   {object}  map[string]string
// @Router       /evm/sanctum/treasury/native/pending [post]
func (h *SanctumHandler) NativePending(w http.ResponseWriter, r *http.Request) {
	var req SanctumNativePendingRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		handler.WriteError(w, http.StatusBadRequest, fmt.Sprintf("invalid request body: %s", err))
		return
	}
	if err := req.ValidateRequest(); err != nil {
		handler.WriteError(w, http.StatusBadRequest, err.Error())
		return
	}

	calldata := core.NativePendingSanctumCalldata(req.UserAddr())
	p := map[string]any{
		"to":   req.SanctumAddr().String(),
		"data": "0x" + hex.EncodeToString(calldata),
	}
	raw, err := h.client.CallContract(r.Context(), p, req.Block)
	if err != nil {
		h.writeRevertOrGatewayError(w, err, "eth_call nativePending() failed")
		return
	}
	data, err := util.ParseHex(raw)
	if err != nil {
		handler.WriteError(w, http.StatusInternalServerError, fmt.Sprintf("failed to parse response: %s", err))
		return
	}

	pending, err := types.DecodeUint256(data)
	if err != nil {
		handler.WriteError(w, http.StatusInternalServerError, fmt.Sprintf("failed to decode pending: %s", err))
		return
	}
	handler.WriteJSON(w, http.StatusOK, SanctumNativePendingResponse{Pending: pending.String()})
}
