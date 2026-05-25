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

// Register godoc
// @Summary      Register caller as a Sanctum member
// @Description  Builds, signs, and broadcasts a register() call to the Sanctum contract. Reverts are decoded into human-readable errors.
// @Tags         sanctum
// @Accept       json
// @Produce      json
// @Param        body  body      SanctumRegisterLegacyRequest     true  "Register request"
// @Success      200   {object}  SanctumRegisterLegacyResponse
// @Failure      400   {object}  map[string]string
// @Failure      502   {object}  map[string]string
// @Router       /evm/sanctum/nexus/register/legacy [post]
func (h *SanctumHandler) RegisterLegacy(w http.ResponseWriter, r *http.Request) {
	var req SanctumRegisterLegacyRequest
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

	calldata := core.RegisterSanctumCalldata()
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

	handler.WriteJSON(w, http.StatusOK, NewSanctumRegisterLegacyResponse(unsigned, signed, txHash, sig))
}

// RegisterEIP1559 godoc
// @Summary      Register caller as a Sanctum member (EIP-1559)
// @Description  Builds, signs, and broadcasts a register() call to the Sanctum contract using an EIP-1559 transaction. Reverts are decoded into human-readable errors.
// @Tags         sanctum
// @Accept       json
// @Produce      json
// @Param        body  body      SanctumRegisterEIP1559Request  true  "Register request"
// @Success      200   {object}  SanctumRegisterEIP1559Response
// @Failure      400   {object}  map[string]string
// @Failure      502   {object}  map[string]string
// @Router       /evm/sanctum/nexus/register/eip1559 [post]
func (h *SanctumHandler) RegisterEIP1559(w http.ResponseWriter, r *http.Request) {
	var req SanctumRegisterEIP1559Request
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

	calldata := core.RegisterSanctumCalldata()
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

	handler.WriteJSON(w, http.StatusOK, NewSanctumRegisterEIP1559Response(unsigned, signed, txHash, sig))
}

// RegisterForLegacy godoc
// @Summary      Register an address as a Sanctum member (legacy)
// @Description  Builds, signs, and broadcasts a registerFor(address) call to the Sanctum contract using a legacy transaction. Reverts are decoded into human-readable errors.
// @Tags         sanctum
// @Accept       json
// @Produce      json
// @Param        body  body      SanctumRegisterForLegacyRequest  true  "RegisterFor request"
// @Success      200   {object}  SanctumRegisterForLegacyResponse
// @Failure      400   {object}  map[string]string
// @Failure      502   {object}  map[string]string
// @Router       /evm/sanctum/nexus/register/for/legacy [post]
func (h *SanctumHandler) RegisterForLegacy(w http.ResponseWriter, r *http.Request) {
	var req SanctumRegisterForLegacyRequest
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

	calldata := core.RegisterForSanctumCalldata(req.TargetAddr())
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

	handler.WriteJSON(w, http.StatusOK, NewSanctumRegisterForLegacyResponse(unsigned, signed, txHash, sig))
}

// RegisterForEIP1559 godoc
// @Summary      Register an address as a Sanctum member (EIP-1559)
// @Description  Builds, signs, and broadcasts a registerFor(address) call to the Sanctum contract using an EIP-1559 transaction. Reverts are decoded into human-readable errors.
// @Tags         sanctum
// @Accept       json
// @Produce      json
// @Param        body  body      SanctumRegisterForEIP1559Request  true  "RegisterFor request"
// @Success      200   {object}  SanctumRegisterForEIP1559Response
// @Failure      400   {object}  map[string]string
// @Failure      502   {object}  map[string]string
// @Router       /evm/sanctum/nexus/register/for/eip1559 [post]
func (h *SanctumHandler) RegisterForEIP1559(w http.ResponseWriter, r *http.Request) {
	var req SanctumRegisterForEIP1559Request
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

	calldata := core.RegisterForSanctumCalldata(req.TargetAddr())
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

	handler.WriteJSON(w, http.StatusOK, NewSanctumRegisterForEIP1559Response(unsigned, signed, txHash, sig))
}

// ApproveRegisterLegacy godoc
// @Summary      Approve a pending Sanctum registration (legacy)
// @Description  Builds, signs, and broadcasts an approveRegister(address) call to the Sanctum contract using a legacy transaction. Reverts are decoded into human-readable errors.
// @Tags         sanctum
// @Accept       json
// @Produce      json
// @Param        body  body      SanctumApproveRegisterLegacyRequest  true  "ApproveRegister request"
// @Success      200   {object}  SanctumApproveRegisterLegacyResponse
// @Failure      400   {object}  map[string]string
// @Failure      502   {object}  map[string]string
// @Router       /evm/sanctum/nexus/register/approve/legacy [post]
func (h *SanctumHandler) ApproveRegisterLegacy(w http.ResponseWriter, r *http.Request) {
	var req SanctumApproveRegisterLegacyRequest
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

	calldata := core.ApproveRegisterSanctumCalldata(req.TargetAddr())
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

	handler.WriteJSON(w, http.StatusOK, NewSanctumApproveRegisterLegacyResponse(unsigned, signed, txHash, sig))
}

// ApproveRegisterEIP1559 godoc
// @Summary      Approve a pending Sanctum registration (EIP-1559)
// @Description  Builds, signs, and broadcasts an approveRegister(address) call to the Sanctum contract using an EIP-1559 transaction. Reverts are decoded into human-readable errors.
// @Tags         sanctum
// @Accept       json
// @Produce      json
// @Param        body  body      SanctumApproveRegisterEIP1559Request  true  "ApproveRegister request"
// @Success      200   {object}  SanctumApproveRegisterEIP1559Response
// @Failure      400   {object}  map[string]string
// @Failure      502   {object}  map[string]string
// @Router       /evm/sanctum/nexus/register/approve/eip1559 [post]
func (h *SanctumHandler) ApproveRegisterEIP1559(w http.ResponseWriter, r *http.Request) {
	var req SanctumApproveRegisterEIP1559Request
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

	calldata := core.ApproveRegisterSanctumCalldata(req.TargetAddr())
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

	handler.WriteJSON(w, http.StatusOK, NewSanctumApproveRegisterEIP1559Response(unsigned, signed, txHash, sig))
}

// DeregisterLegacy godoc
// @Summary      Deregister caller from Sanctum (legacy)
// @Description  Builds, signs, and broadcasts a deregister() call to the Sanctum contract using a legacy transaction.
// @Tags         sanctum
// @Accept       json
// @Produce      json
// @Param        body  body      SanctumDeregisterLegacyRequest   true  "Deregister request"
// @Success      200   {object}  SanctumDeregisterLegacyResponse
// @Failure      400   {object}  map[string]string
// @Failure      502   {object}  map[string]string
// @Router       /evm/sanctum/nexus/deregister/legacy [post]
func (h *SanctumHandler) DeregisterLegacy(w http.ResponseWriter, r *http.Request) {
	var req SanctumDeregisterLegacyRequest
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

	calldata := core.DeregisterSanctumCalldata()
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
	handler.WriteJSON(w, http.StatusOK, NewSanctumDeregisterLegacyResponse(unsigned, signed, txHash, sig))
}

// DeregisterEIP1559 godoc
// @Summary      Deregister caller from Sanctum (EIP-1559)
// @Description  Builds, signs, and broadcasts a deregister() call to the Sanctum contract using an EIP-1559 transaction.
// @Tags         sanctum
// @Accept       json
// @Produce      json
// @Param        body  body      SanctumDeregisterEIP1559Request  true  "Deregister request"
// @Success      200   {object}  SanctumDeregisterEIP1559Response
// @Failure      400   {object}  map[string]string
// @Failure      502   {object}  map[string]string
// @Router       /evm/sanctum/nexus/deregister/eip1559 [post]
func (h *SanctumHandler) DeregisterEIP1559(w http.ResponseWriter, r *http.Request) {
	var req SanctumDeregisterEIP1559Request
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

	calldata := core.DeregisterSanctumCalldata()
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
	handler.WriteJSON(w, http.StatusOK, NewSanctumDeregisterEIP1559Response(unsigned, signed, txHash, sig))
}

// DeregisterForLegacy godoc
// @Summary      Deregister an address from Sanctum (legacy)
// @Description  Builds, signs, and broadcasts a deregisterFor(address) call to the Sanctum contract using a legacy transaction.
// @Tags         sanctum
// @Accept       json
// @Produce      json
// @Param        body  body      SanctumDeregisterForLegacyRequest   true  "DeregisterFor request"
// @Success      200   {object}  SanctumDeregisterForLegacyResponse
// @Failure      400   {object}  map[string]string
// @Failure      502   {object}  map[string]string
// @Router       /evm/sanctum/nexus/deregister/for/legacy [post]
func (h *SanctumHandler) DeregisterForLegacy(w http.ResponseWriter, r *http.Request) {
	var req SanctumDeregisterForLegacyRequest
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

	calldata := core.DeregisterForSanctumCalldata(req.TargetAddr())
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
	handler.WriteJSON(w, http.StatusOK, NewSanctumDeregisterForLegacyResponse(unsigned, signed, txHash, sig))
}

// DeregisterForEIP1559 godoc
// @Summary      Deregister an address from Sanctum (EIP-1559)
// @Description  Builds, signs, and broadcasts a deregisterFor(address) call to the Sanctum contract using an EIP-1559 transaction.
// @Tags         sanctum
// @Accept       json
// @Produce      json
// @Param        body  body      SanctumDeregisterForEIP1559Request  true  "DeregisterFor request"
// @Success      200   {object}  SanctumDeregisterForEIP1559Response
// @Failure      400   {object}  map[string]string
// @Failure      502   {object}  map[string]string
// @Router       /evm/sanctum/nexus/deregister/for/eip1559 [post]
func (h *SanctumHandler) DeregisterForEIP1559(w http.ResponseWriter, r *http.Request) {
	var req SanctumDeregisterForEIP1559Request
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

	calldata := core.DeregisterForSanctumCalldata(req.TargetAddr())
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
	handler.WriteJSON(w, http.StatusOK, NewSanctumDeregisterForEIP1559Response(unsigned, signed, txHash, sig))
}

// GetAccounts godoc
// @Summary      Get all registered Sanctum accounts
// @Description  Calls getAccounts() on the Sanctum contract and returns the list of registered addresses.
// @Tags         sanctum
// @Accept       json
// @Produce      json
// @Param        body  body      SanctumGetAccountsRequest   true  "GetAccounts request"
// @Success      200   {object}  SanctumGetAccountsResponse
// @Failure      400   {object}  map[string]string
// @Failure      502   {object}  map[string]string
// @Router       /evm/sanctum/nexus/account/list [post]
func (h *SanctumHandler) GetAccounts(w http.ResponseWriter, r *http.Request) {
	var req SanctumGetAccountsRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		handler.WriteError(w, http.StatusBadRequest, fmt.Sprintf("invalid request body: %s", err))
		return
	}
	if err := req.ValidateRequest(); err != nil {
		handler.WriteError(w, http.StatusBadRequest, err.Error())
		return
	}

	calldata := core.GetAccountsSanctumCalldata()
	p := map[string]any{
		"to":   req.SanctumAddr().String(),
		"data": "0x" + hex.EncodeToString(calldata),
	}
	raw, err := h.client.CallContract(r.Context(), p, req.Block)
	if err != nil {
		h.writeRevertOrGatewayError(w, err, "eth_call getAccounts() failed")
		return
	}
	data, err := util.ParseHex(raw)
	if err != nil {
		handler.WriteError(w, http.StatusInternalServerError, fmt.Sprintf("failed to parse response: %s", err))
		return
	}

	addrs, err := core.ABI.DecodeAddressSlice(data)
	if err != nil {
		handler.WriteError(w, http.StatusInternalServerError, fmt.Sprintf("failed to decode accounts: %s", err))
		return
	}

	accounts := make([]string, len(addrs))
	for i, a := range addrs {
		accounts[i] = a.Hex()
	}
	handler.WriteJSON(w, http.StatusOK, SanctumGetAccountsResponse{Accounts: accounts})
}

// AccountCount godoc
// @Summary      Get the number of registered Sanctum accounts
// @Description  Calls accountCount() on the Sanctum contract.
// @Tags         sanctum
// @Accept       json
// @Produce      json
// @Param        body  body      SanctumAccountCountRequest   true  "AccountCount request"
// @Success      200   {object}  SanctumAccountCountResponse
// @Failure      400   {object}  map[string]string
// @Failure      502   {object}  map[string]string
// @Router       /evm/sanctum/nexus/account/count [post]
func (h *SanctumHandler) AccountCount(w http.ResponseWriter, r *http.Request) {
	var req SanctumAccountCountRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		handler.WriteError(w, http.StatusBadRequest, fmt.Sprintf("invalid request body: %s", err))
		return
	}
	if err := req.ValidateRequest(); err != nil {
		handler.WriteError(w, http.StatusBadRequest, err.Error())
		return
	}

	calldata := core.AccountCountSanctumCalldata()
	p := map[string]any{
		"to":   req.SanctumAddr().String(),
		"data": "0x" + hex.EncodeToString(calldata),
	}
	raw, err := h.client.CallContract(r.Context(), p, req.Block)
	if err != nil {
		h.writeRevertOrGatewayError(w, err, "eth_call accountCount() failed")
		return
	}
	data, err := util.ParseHex(raw)
	if err != nil {
		handler.WriteError(w, http.StatusInternalServerError, fmt.Sprintf("failed to parse response: %s", err))
		return
	}

	count, err := core.ABI.DecodeUint256(data)
	if err != nil {
		handler.WriteError(w, http.StatusInternalServerError, fmt.Sprintf("failed to decode count: %s", err))
		return
	}
	handler.WriteJSON(w, http.StatusOK, SanctumAccountCountResponse{Count: count.String()})
}

// GetAccountInfo godoc
// @Summary      Get account info for a Sanctum member
// @Description  Calls getAccountInfo(address) on the Sanctum contract and returns the account's address, role, and registration block.
// @Tags         sanctum
// @Accept       json
// @Produce      json
// @Param        body  body      SanctumGetAccountInfoRequest   true  "GetAccountInfo request"
// @Success      200   {object}  SanctumGetAccountInfoResponse
// @Failure      400   {object}  map[string]string
// @Failure      502   {object}  map[string]string
// @Router       /evm/sanctum/nexus/account/info [post]
func (h *SanctumHandler) GetAccountInfo(w http.ResponseWriter, r *http.Request) {
	var req SanctumGetAccountInfoRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		handler.WriteError(w, http.StatusBadRequest, fmt.Sprintf("invalid request body: %s", err))
		return
	}
	if err := req.ValidateRequest(); err != nil {
		handler.WriteError(w, http.StatusBadRequest, err.Error())
		return
	}

	calldata := core.GetAccountInfoSanctumCalldata(req.AccountAddr())
	p := map[string]any{
		"to":   req.SanctumAddr().String(),
		"data": "0x" + hex.EncodeToString(calldata),
	}
	raw, err := h.client.CallContract(r.Context(), p, req.Block)
	if err != nil {
		h.writeRevertOrGatewayError(w, err, "eth_call getAccountInfo() failed")
		return
	}
	data, err := util.ParseHex(raw)
	if err != nil {
		handler.WriteError(w, http.StatusInternalServerError, fmt.Sprintf("failed to parse response: %s", err))
		return
	}

	info, err := h.decodeAccountInfo(data)
	if err != nil {
		handler.WriteError(w, http.StatusInternalServerError, fmt.Sprintf("failed to decode account info: %s", err))
		return
	}
	handler.WriteJSON(w, http.StatusOK, NewSanctumGetAccountInfoResponse(info))
}
