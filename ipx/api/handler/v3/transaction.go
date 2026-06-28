package v3

import (
	"encoding/json"
	"fmt"
	"math/big"
	"net/http"

	"github.com/andantan/evmlab/api/handler"
	"github.com/andantan/evmlab/core"
	"github.com/andantan/evmlab/internal/config"
	"github.com/andantan/evmlab/internal/rpc"
)

type TransactionHandler struct {
	cfg    *config.Config
	client *rpc.Client
}

func NewTransactionHandler(cfg *config.Config, client *rpc.Client) *TransactionHandler {
	return &TransactionHandler{
		cfg:    cfg,
		client: client,
	}
}

// BuildNativeLegacyTransaction godoc
// @Summary      Build and sign legacy transfer tx
// @Description  Fetches chain state, builds and signs an EIP-155 legacy native transfer, returns unsigned_rlp, signed_rlp, and tx_hash
// @Tags         transaction
// @Accept       json
// @Produce      json
// @Param        body  body      BuildNativeLegacyTransactionRequest  true  "Transfer request"
// @Success      200   {object}  BuildNativeLegacyTransactionResponse
// @Failure      400   {object}  map[string]string
// @Failure      502   {object}  map[string]string
// @Router       /evm/v3/transaction/native/legacy [post]
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

	msg := core.NewCallMsg(req.FromAddr(), req.ToAddr(), req.Amount(), nil)
	state, err := core.GenerateLegacyTransactionState(r.Context(), h.client, msg)
	if err != nil {
		handler.WriteError(w, http.StatusBadGateway, err.Error())
		return
	}

	tx := state.ToTx()
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

	handler.WriteJSON(w, http.StatusOK, NewBuildNativeLegacyTransactionResponse(unsigned, signed, core.Hasher.Hash(signed), sig))
}

// BuildNativeEIP1559Transaction godoc
// @Summary      Build and sign EIP-1559 transfer tx
// @Description  Fetches chain state, builds and signs an EIP-1559 native transfer, returns unsigned_rlp, signed_rlp, and tx_hash
// @Tags         transaction
// @Accept       json
// @Produce      json
// @Param        body  body      BuildNativeEIP1559TransactionRequest  true  "Transfer request"
// @Success      200   {object}  BuildNativeEIP1559TransactionResponse
// @Failure      400   {object}  map[string]string
// @Failure      502   {object}  map[string]string
// @Router       /evm/v3/transaction/native/eip1559 [post]
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

	msg := core.NewCallMsg(req.FromAddr(), req.ToAddr(), req.Amount(), nil)
	state, err := core.GenerateEIP1559TransactionState(r.Context(), h.client, msg)
	if err != nil {
		handler.WriteError(w, http.StatusBadGateway, err.Error())
		return
	}

	tx := state.ToTx()
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

	handler.WriteJSON(w, http.StatusOK, NewBuildNativeEIP1559TransactionResponse(unsigned, signed, core.Hasher.Hash(signed), sig))
}

// BuildERC20LegacyTransaction godoc
// @Summary      Build and sign ERC-20 legacy transfer tx
// @Description  Builds transfer(address,uint256) calldata, estimates gas, signs with the configured key, returns unsigned_rlp, signed_rlp, and tx_hash
// @Tags         transaction
// @Accept       json
// @Produce      json
// @Param        body  body      BuildERC20LegacyTransactionRequest  true  "ERC-20 transfer request"
// @Success      200   {object}  BuildERC20LegacyTransactionResponse
// @Failure      400   {object}  map[string]string
// @Failure      502   {object}  map[string]string
// @Router       /evm/v3/transaction/erc20/legacy [post]
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

	calldata := core.TransferCalldata(req.ToAddr(), req.ToAmount())
	msg := core.NewCallMsg(req.FromAddr(), req.ContractAddr(), big.NewInt(0), calldata)
	state, err := core.GenerateLegacyTransactionState(r.Context(), h.client, msg)
	if err != nil {
		handler.WriteError(w, http.StatusBadGateway, err.Error())
		return
	}

	tx := state.ToTx()
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

	handler.WriteJSON(w, http.StatusOK, NewBuildERC20LegacyTransactionResponse(unsigned, signed, core.Hasher.Hash(signed), sig))
}

// BuildERC20EIP1559Transaction godoc
// @Summary      Build and sign ERC-20 EIP-1559 transfer tx
// @Description  Builds transfer(address,uint256) calldata, estimates gas, signs with the configured key, returns unsigned_rlp, signed_rlp, and tx_hash
// @Tags         transaction
// @Accept       json
// @Produce      json
// @Param        body  body      BuildERC20EIP1559TransactionRequest  true  "ERC-20 transfer request"
// @Success      200   {object}  BuildERC20EIP1559TransactionResponse
// @Failure      400   {object}  map[string]string
// @Failure      502   {object}  map[string]string
// @Router       /evm/v3/transaction/erc20/eip1559 [post]
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

	calldata := core.TransferCalldata(req.ToAddr(), req.ToAmount())
	msg := core.NewCallMsg(req.FromAddr(), req.ContractAddr(), big.NewInt(0), calldata)
	state, err := core.GenerateEIP1559TransactionState(r.Context(), h.client, msg)
	if err != nil {
		handler.WriteError(w, http.StatusBadGateway, err.Error())
		return
	}

	tx := state.ToTx()
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

	handler.WriteJSON(w, http.StatusOK, NewBuildERC20EIP1559TransactionResponse(unsigned, signed, core.Hasher.Hash(signed), sig))
}

// BuildContractCallLegacyTransaction godoc
// @Summary      Build and sign legacy contract call tx
// @Description  Fetches chain state, estimates gas, signs with the configured key, returns unsigned_rlp, signed_rlp, and tx_hash
// @Tags         transaction
// @Accept       json
// @Produce      json
// @Param        body  body      BuildContractCallLegacyTransactionRequest  true  "Contract call request"
// @Success      200   {object}  BuildContractCallLegacyTransactionResponse
// @Failure      400   {object}  map[string]string
// @Failure      502   {object}  map[string]string
// @Router       /evm/v3/transaction/contract/legacy [post]
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

	msg := core.NewCallMsg(req.FromAddr(), req.ToAddr(), req.Amount(), req.Calldata())
	state, err := core.GenerateLegacyTransactionState(r.Context(), h.client, msg)
	if err != nil {
		handler.WriteError(w, http.StatusBadGateway, err.Error())
		return
	}

	tx := state.ToTx()
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

	handler.WriteJSON(w, http.StatusOK, NewBuildContractCallLegacyTransactionResponse(unsigned, signed, core.Hasher.Hash(signed), sig))
}

// BuildContractCallEIP1559Transaction godoc
// @Summary      Build and sign EIP-1559 contract call tx
// @Description  Fetches chain state, estimates gas, signs with the configured key, returns unsigned_rlp, signed_rlp, and tx_hash
// @Tags         transaction
// @Accept       json
// @Produce      json
// @Param        body  body      BuildContractCallEIP1559TransactionRequest  true  "Contract call request"
// @Success      200   {object}  BuildContractCallEIP1559TransactionResponse
// @Failure      400   {object}  map[string]string
// @Failure      502   {object}  map[string]string
// @Router       /evm/v3/transaction/contract/eip1559 [post]
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

	msg := core.NewCallMsg(req.FromAddr(), req.ToAddr(), req.Amount(), req.Calldata())
	state, err := core.GenerateEIP1559TransactionState(r.Context(), h.client, msg)
	if err != nil {
		handler.WriteError(w, http.StatusBadGateway, err.Error())
		return
	}

	tx := state.ToTx()
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

	handler.WriteJSON(w, http.StatusOK, NewBuildContractCallEIP1559TransactionResponse(unsigned, signed, core.Hasher.Hash(signed), sig))
}
