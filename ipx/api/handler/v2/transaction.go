package v2

import (
	"encoding/json"
	"fmt"
	"math/big"
	"net/http"

	"github.com/andantan/evmlab/api/handler"
	"github.com/andantan/evmlab/core"
	"github.com/andantan/evmlab/internal/rpc"
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
// @Tags         transaction
// @Accept       json
// @Produce      json
// @Param        body  body      BuildNativeLegacyTransactionRequest   true  "Transfer request"
// @Success      200   {object}  BuildNativeLegacyTransactionResponse
// @Failure      400   {object}  map[string]string
// @Failure      502   {object}  map[string]string
// @Router       /evm/v2/transaction/native/legacy [post]
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

	msg := core.NewCallMsg(req.FromAddr(), req.ToAddr(), req.Amount(), nil)
	state, err := core.GenerateLegacyTransactionState(r.Context(), h.client, msg)
	if err != nil {
		handler.WriteError(w, http.StatusBadGateway, err.Error())
		return
	}

	unsigned, err := core.RLP.EncodeLegacyUnsigned(state.ToTx())
	if err != nil {
		handler.WriteError(w, http.StatusInternalServerError, fmt.Sprintf("failed to encode tx: %s", err))
		return
	}

	handler.WriteJSON(w, http.StatusOK, NewBuildNativeLegacyTransaction(unsigned))
}

// BuildNativeEIP1559Transaction godoc
// @Summary      Build unsigned EIP-1559 transfer tx
// @Description  Fetches chain state (chainID, nonce, fees) and returns the unsigned RLP-encoded EIP-1559 transfer transaction
// @Tags         transaction
// @Accept       json
// @Produce      json
// @Param        body  body      BuildNativeEIP1559TransactionRequest   true  "BuildNativeEIP1559Transaction request"
// @Success      200   {object}  BuildNativeEIP1559TransactionResponse
// @Failure      400   {object}  map[string]string
// @Failure      502   {object}  map[string]string
// @Router       /evm/v2/transaction/native/eip1559 [post]
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

	msg := core.NewCallMsg(req.FromAddr(), req.ToAddr(), req.Amount(), nil)
	state, err := core.GenerateEIP1559TransactionState(r.Context(), h.client, msg)
	if err != nil {
		handler.WriteError(w, http.StatusBadGateway, err.Error())
		return
	}

	unsigned, err := core.RLP.EncodeDynamicFeeUnsigned(state.ToTx())
	if err != nil {
		handler.WriteError(w, http.StatusInternalServerError, fmt.Sprintf("failed to encode tx: %s", err))
		return
	}

	handler.WriteJSON(w, http.StatusOK, NewBuildNativeEIP1559Transaction(unsigned))
}

// BuildERC20LegacyTransaction godoc
// @Summary      Build unsigned ERC-20 legacy transfer tx
// @Description  Builds transfer(address,uint256) calldata internally and returns the unsigned RLP-encoded EIP-155 legacy transaction. Gas limit is estimated on-chain with a 2x buffer.
// @Tags         transaction
// @Accept       json
// @Produce      json
// @Param        body  body      BuildERC20LegacyTransactionRequest   true  "ERC-20 transfer request"
// @Success      200   {object}  BuildERC20LegacyTransactionResponse
// @Failure      400   {object}  map[string]string
// @Failure      502   {object}  map[string]string
// @Router       /evm/v2/transaction/erc20/legacy [post]
func (h *TransactionHandler) BuildERC20LegacyTransaction(w http.ResponseWriter, r *http.Request) {
	var req BuildERC20LegacyTransactionRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		handler.WriteError(w, http.StatusBadRequest, fmt.Sprintf("invalid request body: %s", err))
		return
	}
	if err := req.ValidateRequest(); err != nil {
		handler.WriteError(w, http.StatusBadRequest, err.Error())
		return
	}

	calldata := core.TransferCalldata(req.ToAddr(), req.ToAmount())
	msg := core.NewCallMsg(req.FromAddr(), req.ContractAddr(), big.NewInt(0), calldata)
	state, err := core.GenerateLegacyTransactionState(r.Context(), h.client, msg)
	if err != nil {
		handler.WriteError(w, http.StatusBadGateway, err.Error())
		return
	}

	unsigned, err := core.RLP.EncodeLegacyUnsigned(state.ToTx())
	if err != nil {
		handler.WriteError(w, http.StatusInternalServerError, fmt.Sprintf("failed to encode tx: %s", err))
		return
	}

	handler.WriteJSON(w, http.StatusOK, NewBuildERC20LegacyTransaction(unsigned))
}

// BuildERC20EIP1559Transaction godoc
// @Summary      Build unsigned ERC-20 EIP-1559 transfer tx
// @Description  Builds transfer(address,uint256) calldata internally and returns the unsigned RLP-encoded EIP-1559 transaction. Gas limit is estimated on-chain with a 2x buffer.
// @Tags         transaction
// @Accept       json
// @Produce      json
// @Param        body  body      BuildERC20EIP1559TransactionRequest   true  "ERC-20 transfer request"
// @Success      200   {object}  BuildERC20EIP1559TransactionResponse
// @Failure      400   {object}  map[string]string
// @Failure      502   {object}  map[string]string
// @Router       /evm/v2/transaction/erc20/eip1559 [post]
func (h *TransactionHandler) BuildERC20EIP1559Transaction(w http.ResponseWriter, r *http.Request) {
	var req BuildERC20EIP1559TransactionRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		handler.WriteError(w, http.StatusBadRequest, fmt.Sprintf("invalid request body: %s", err))
		return
	}
	if err := req.ValidateRequest(); err != nil {
		handler.WriteError(w, http.StatusBadRequest, err.Error())
		return
	}

	calldata := core.TransferCalldata(req.ToAddr(), req.ToAmount())
	msg := core.NewCallMsg(req.FromAddr(), req.ContractAddr(), big.NewInt(0), calldata)
	state, err := core.GenerateEIP1559TransactionState(r.Context(), h.client, msg)
	if err != nil {
		handler.WriteError(w, http.StatusBadGateway, err.Error())
		return
	}

	unsigned, err := core.RLP.EncodeDynamicFeeUnsigned(state.ToTx())
	if err != nil {
		handler.WriteError(w, http.StatusInternalServerError, fmt.Sprintf("failed to encode tx: %s", err))
		return
	}

	handler.WriteJSON(w, http.StatusOK, NewBuildERC20EIP1559Transaction(unsigned))
}

// BuildContractCallLegacyTransaction godoc
// @Summary      Build unsigned legacy contract call tx
// @Description  Fetches chain state (chainID, nonce, gas price), estimates gas, and returns the unsigned RLP-encoded EIP-155 legacy contract call transaction
// @Tags         transaction
// @Accept       json
// @Produce      json
// @Param        body  body      BuildContractCallLegacyTransactionRequest   true  "Contract call request"
// @Success      200   {object}  BuildContractCallLegacyTransactionResponse
// @Failure      400   {object}  map[string]string
// @Failure      502   {object}  map[string]string
// @Router       /evm/v2/transaction/contract/legacy [post]
func (h *TransactionHandler) BuildContractCallLegacyTransaction(w http.ResponseWriter, r *http.Request) {
	var req BuildContractCallLegacyTransactionRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		handler.WriteError(w, http.StatusBadRequest, fmt.Sprintf("invalid request body: %s", err))
		return
	}
	if err := req.ValidateRequest(); err != nil {
		handler.WriteError(w, http.StatusBadRequest, err.Error())
		return
	}

	msg := core.NewCallMsg(req.FromAddr(), req.ToAddr(), req.Amount(), req.Calldata())
	state, err := core.GenerateLegacyTransactionState(r.Context(), h.client, msg)
	if err != nil {
		handler.WriteError(w, http.StatusBadGateway, err.Error())
		return
	}

	unsigned, err := core.RLP.EncodeLegacyUnsigned(state.ToTx())
	if err != nil {
		handler.WriteError(w, http.StatusInternalServerError, fmt.Sprintf("failed to encode tx: %s", err))
		return
	}

	handler.WriteJSON(w, http.StatusOK, NewBuildContractCallLegacyTransaction(unsigned))
}

// BuildContractCallEIP1559Transaction godoc
// @Summary      Build unsigned EIP-1559 contract call tx
// @Description  Fetches chain state (chainID, nonce, fees), estimates gas, and returns the unsigned RLP-encoded EIP-1559 contract call transaction
// @Tags         transaction
// @Accept       json
// @Produce      json
// @Param        body  body      BuildContractCallEIP1559TransactionRequest   true  "Contract call request"
// @Success      200   {object}  BuildContractCallEIP1559TransactionResponse
// @Failure      400   {object}  map[string]string
// @Failure      502   {object}  map[string]string
// @Router       /evm/v2/transaction/contract/eip1559 [post]
func (h *TransactionHandler) BuildContractCallEIP1559Transaction(w http.ResponseWriter, r *http.Request) {
	var req BuildContractCallEIP1559TransactionRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		handler.WriteError(w, http.StatusBadRequest, fmt.Sprintf("invalid request body: %s", err))
		return
	}
	if err := req.ValidateRequest(); err != nil {
		handler.WriteError(w, http.StatusBadRequest, err.Error())
		return
	}

	msg := core.NewCallMsg(req.FromAddr(), req.ToAddr(), req.Amount(), req.Calldata())
	state, err := core.GenerateEIP1559TransactionState(r.Context(), h.client, msg)
	if err != nil {
		handler.WriteError(w, http.StatusBadGateway, err.Error())
		return
	}

	unsigned, err := core.RLP.EncodeDynamicFeeUnsigned(state.ToTx())
	if err != nil {
		handler.WriteError(w, http.StatusInternalServerError, fmt.Sprintf("failed to encode tx: %s", err))
		return
	}

	handler.WriteJSON(w, http.StatusOK, NewBuildContractCallEIP1559Transaction(unsigned))
}
