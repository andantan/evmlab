package v2

import (
	"encoding/json"
	"fmt"
	"math/big"
	"net/http"

	"github.com/andantan/evmlab/api/handler"
	"github.com/andantan/evmlab/core"
	"github.com/andantan/evmlab/core/types"
	"github.com/andantan/evmlab/internal/rpc"
	"github.com/andantan/evmlab/internal/util"
)

type TransactionHandler struct {
	client *rpc.Client
}

func NewTransactionHandler(client *rpc.Client) *TransactionHandler {
	return &TransactionHandler{client: client}
}

// BuildNativeLegacyTransaction godoc
// @Summary      Build unsigned legacy transfer tx
// @Description  Fetches chain state (chainID, nonce, gas price) and returns the unsigned RLP-encoded EIP-155 legacy transfer transaction
// @Tags         transfer
// @Accept       json
// @Produce      json
// @Param        body  body      BuildNativeLegacyTransactionRequest   true  "Transfer request"
// @Success      200   {object}  BuildNativeLegacyTransactionResponse
// @Failure      400   {object}  map[string]string
// @Failure      502   {object}  map[string]string
// @Router       /evm/v2/transfers/native/legacy [post]
func (h *TransactionHandler) BuildNativeLegacyTransaction(w http.ResponseWriter, r *http.Request) {
	var req BuildNativeLegacyTransactionRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		handler.WriteError(w, http.StatusBadRequest, fmt.Sprintf("invalid request body: %s", err))
		return
	}
	if err := req.ValidateRequest(); err != nil {
		handler.WriteError(w, http.StatusBadRequest, err.Error())
		return
	}

	chainIDHex, err := h.client.ChainID(r.Context())
	if err != nil {
		handler.WriteError(w, http.StatusBadGateway, fmt.Sprintf("failed t get chain id: %s", err))
		return
	}
	chainID, err := util.HexToBigInt(chainIDHex)
	if err != nil {
		handler.WriteError(w, http.StatusInternalServerError, fmt.Sprintf("failed t parse chain id: %s", err))
		return
	}

	nonceHex, err := h.client.GetTransactionCount(r.Context(), req.From, "pending")
	if err != nil {
		handler.WriteError(w, http.StatusBadGateway, fmt.Sprintf("failed t get nonce: %s", err))
		return
	}
	nonce, err := util.HexToUint64(nonceHex)
	if err != nil {
		handler.WriteError(w, http.StatusInternalServerError, fmt.Sprintf("failed t parse nonce: %s", err))
		return
	}

	gasPriceHex, err := h.client.GasPrice(r.Context())
	if err != nil {
		handler.WriteError(w, http.StatusBadGateway, fmt.Sprintf("failed t get gas price: %s", err))
		return
	}
	gasPrice, err := util.HexToBigInt(gasPriceHex)
	if err != nil {
		handler.WriteError(w, http.StatusInternalServerError, fmt.Sprintf("failed t parse gas price: %s", err))
		return
	}

	tx := &types.LegacyTx{
		ChainID:  chainID,
		Nonce:    nonce,
		GasPrice: gasPrice,
		GasLimit: 21000,
		To:       &req.ToAddr().Addr,
		Value:    req.Amount(),
		Data:     nil,
	}

	unsigned, err := core.RLP.EncodeLegacyUnsigned(tx)
	if err != nil {
		handler.WriteError(w, http.StatusInternalServerError, fmt.Sprintf("failed t encode tx: %s", err))
		return
	}

	handler.WriteJSON(w, http.StatusOK, NewBuildNativeLegacyTransaction(unsigned))
}

// BuildNativeEIP1559Transaction godoc
// @Summary      Build unsigned EIP-1559 transfer tx
// @Description  Fetches chain state (chainID, nonce, fees) and returns the unsigned RLP-encoded EIP-1559 transfer transaction
// @Tags         transfer
// @Accept       json
// @Produce      json
// @Param        body  body      BuildNativeEIP1559TransactionRequest   true  "BuildNativeEIP1559Transaction request"
// @Success      200   {object}  BuildNativeEIP1559TransactionResponse
// @Failure      400   {object}  map[string]string
// @Failure      502   {object}  map[string]string
// @Router       /evm/v2/transfers/native/eip1559 [post]
func (h *TransactionHandler) BuildNativeEIP1559Transaction(w http.ResponseWriter, r *http.Request) {
	var req BuildNativeEIP1559TransactionRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		handler.WriteError(w, http.StatusBadRequest, fmt.Sprintf("invalid request body: %s", err))
		return
	}
	if err := req.ValidateRequest(); err != nil {
		handler.WriteError(w, http.StatusBadRequest, err.Error())
		return
	}

	chainIDHex, err := h.client.ChainID(r.Context())
	if err != nil {
		handler.WriteError(w, http.StatusBadGateway, fmt.Sprintf("failed t get chain id: %s", err))
		return
	}
	chainID, err := util.HexToBigInt(chainIDHex)
	if err != nil {
		handler.WriteError(w, http.StatusInternalServerError, fmt.Sprintf("failed t parse chain id: %s", err))
		return
	}

	nonceHex, err := h.client.GetTransactionCount(r.Context(), req.From, "pending")
	if err != nil {
		handler.WriteError(w, http.StatusBadGateway, fmt.Sprintf("failed t get nonce: %s", err))
		return
	}
	nonce, err := util.HexToUint64(nonceHex)
	if err != nil {
		handler.WriteError(w, http.StatusInternalServerError, fmt.Sprintf("failed t parse nonce: %s", err))
		return
	}

	tipCapHex, err := h.client.MaxPriorityFeePerGas(r.Context())
	if err != nil {
		handler.WriteError(w, http.StatusBadGateway, fmt.Sprintf("failed t get tip cap: %s", err))
		return
	}
	tipCap, err := util.HexToBigInt(tipCapHex)
	if err != nil {
		handler.WriteError(w, http.StatusInternalServerError, fmt.Sprintf("failed t parse tip cap: %s", err))
		return
	}

	baseFeeHex, err := h.client.BaseFeePerGas(r.Context())
	if err != nil {
		handler.WriteError(w, http.StatusBadGateway, fmt.Sprintf("failed t get base fee: %s", err))
		return
	}
	baseFee, err := util.HexToBigInt(baseFeeHex)
	if err != nil {
		handler.WriteError(w, http.StatusInternalServerError, fmt.Sprintf("failed t parse base fee: %s", err))
		return
	}
	feeCap := new(big.Int).Add(new(big.Int).Mul(baseFee, big.NewInt(2)), tipCap)

	tx := &types.DynamicFeeTx{
		ChainID:   chainID,
		Nonce:     nonce,
		GasTipCap: tipCap,
		GasFeeCap: feeCap,
		GasLimit:  21000,
		To:        &req.ToAddr().Addr,
		Value:     req.Amount(),
		Data:      nil,
	}

	unsigned, err := core.RLP.EncodeDynamicFeeUnsigned(tx)
	if err != nil {
		handler.WriteError(w, http.StatusInternalServerError, fmt.Sprintf("failed t encode tx: %s", err))
		return
	}

	handler.WriteJSON(w, http.StatusOK, NewBuildNativeEIP1559Transaction(unsigned))
}
