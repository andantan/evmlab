package v1

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/andantan/evmlab/api/handler"
	"github.com/andantan/evmlab/core"
	"github.com/andantan/evmlab/internal/config"
)

type TransactionHandler struct {
	cfg *config.Config
}

func NewTransactionHandler(cfg *config.Config) *TransactionHandler {
	return &TransactionHandler{cfg: cfg}
}

// BuildLegacyTransaction godoc
// @Summary      Build an unsigned legacy native transfer transaction
// @Description  Constructs an unsigned EIP-155 legacy transaction for native ETH transfer and returns the RLP encoding and signing hash
// @Tags         transaction
// @Accept       json
// @Produce      json
// @Param        body  body      BuildLegacyTransactionRequest  true  "Transaction fields"
// @Success      200   {object}  BuildLegacyTransactionResponse
// @Failure      400   {object}  map[string]string
// @Router       /evm/v1/transaction/legacy/build [post]
func (h *TransactionHandler) BuildLegacyTransaction(w http.ResponseWriter, r *http.Request) {
	req := new(BuildLegacyTransactionRequest)
	if err := json.NewDecoder(r.Body).Decode(req); err != nil {
		handler.WriteError(w, http.StatusBadRequest, fmt.Sprintf("invalid request body: %s", err))
		return
	}
	if err := req.ValidateRequest(); err != nil {
		handler.WriteError(w, http.StatusBadRequest, err.Error())
		return
	}

	raw, err := core.RLP.EncodeLegacyUnsigned(req.ToLegacyTx())
	if err != nil {
		handler.WriteError(w, http.StatusInternalServerError, fmt.Sprintf("failed to encode tx: %s", err))
		return
	}

	signingHash := core.Hasher.Hash(raw)
	handler.WriteJSON(w, http.StatusOK, NewBuildLegacyNativeTransferResponse(raw, signingHash))
}

// BuildEIP1559Transaction godoc
// @Summary      Build an unsigned dynamic fee native transfer transaction
// @Description  Constructs an unsigned EIP-1559 transaction for native ETH transfer and returns the encoded signing payload and signing hash
// @Tags         transaction
// @Accept       json
// @Produce      json
// @Param        body  body      BuildEIP1559TransactionRequest  true  "Transaction fields"
// @Success      200   {object}  BuildEIP1559TransactionResponse
// @Failure      400   {object}  map[string]string
// @Router       /evm/v1/transaction/eip1559/build [post]
func (h *TransactionHandler) BuildEIP1559Transaction(w http.ResponseWriter, r *http.Request) {
	req := new(BuildEIP1559TransactionRequest)
	if err := json.NewDecoder(r.Body).Decode(req); err != nil {
		handler.WriteError(w, http.StatusBadRequest, fmt.Sprintf("invalid request body: %s", err))
		return
	}
	if err := req.ValidateRequest(); err != nil {
		handler.WriteError(w, http.StatusBadRequest, err.Error())
		return
	}

	raw, err := core.RLP.EncodeDynamicFeeUnsigned(req.ToDynamicFeeTx())
	if err != nil {
		handler.WriteError(w, http.StatusInternalServerError, fmt.Sprintf("failed to encode tx: %s", err))
		return
	}

	signingHash := core.Hasher.Hash(raw)
	handler.WriteJSON(w, http.StatusOK, NewBuildEIP1559TransactionResponse(raw, signingHash))
}
