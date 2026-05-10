package misc

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/andantan/evmlab/api/handler"
	"github.com/andantan/evmlab/core"
)

type HashHandler struct{}

func NewHashHandler() *HashHandler {
	return &HashHandler{}
}

// Keccak256Legacy godoc
// @Summary      Compute raw Keccak256 hash
// @Description  Computes the Keccak256 hash of the given message with no prefix applied (no EIP standard)
// @Tags         hash
// @Accept       json
// @Produce      json
// @Param        body  body      Keccak256LegacyRequest   true  "Message to hash"
// @Success      200   {object}  Keccak256LegacyResponse
// @Failure      400   {object}  map[string]string
// @Router       /evm/hash/keccak256/legacy [post]
func (h *HashHandler) Keccak256Legacy(w http.ResponseWriter, r *http.Request) {
	req := new(Keccak256LegacyRequest)
	if err := json.NewDecoder(r.Body).Decode(req); err != nil {
		handler.WriteError(w, http.StatusBadRequest, fmt.Sprintf("invalid request body: %s", err))
		return
	}
	if err := req.ValidateRequest(); err != nil {
		handler.WriteError(w, http.StatusBadRequest, err.Error())
		return
	}

	hash := core.Hasher.HashString(req.Message)
	handler.WriteJSON(w, http.StatusOK, NewKeccak256LegacyResponse(hash))
}

// Keccak256EIP191 godoc
// @Summary      Compute Keccak256 hash with EIP-191 prefix
// @Description  Prepends the EIP-191 personal sign prefix ("\x19Ethereum Signed Message:\n" + length) to the message and returns the Keccak256 hash — matches the digest produced by eth_sign / personal_sign
// @Tags         hash
// @Accept       json
// @Produce      json
// @Param        body  body      Keccak256EIP191Request   true  "Message to hash"
// @Success      200   {object}  Keccak256EIP191Response
// @Failure      400   {object}  map[string]string
// @Router       /evm/hash/keccak256/eip191 [post]
func (h *HashHandler) Keccak256EIP191(w http.ResponseWriter, r *http.Request) {
	req := new(Keccak256EIP191Request)
	if err := json.NewDecoder(r.Body).Decode(req); err != nil {
		handler.WriteError(w, http.StatusBadRequest, fmt.Sprintf("invalid request body: %s", err))
		return
	}
	if err := req.ValidateRequest(); err != nil {
		handler.WriteError(w, http.StatusBadRequest, err.Error())
		return
	}

	hash := core.Hasher.EIP191([]byte(req.Message))
	handler.WriteJSON(w, http.StatusOK, NewKeccak256EIP191Response(hash))
}

// Keccak256EIP712 godoc
// @Summary      Compute EIP-712 typed data hash
// @Description  Parses the function signature to derive the type schema, builds the EIP-712 typed data from domain and args, and returns the digest (\x19\x01 || domainSeparator || hashStruct(message))
// @Tags         hash
// @Accept       json
// @Produce      json
// @Param        body  body      Keccak256EIP712Request   true  "EIP-712 domain, signature and args"
// @Success      200   {object}  Keccak256EIP712Response
// @Failure      400   {object}  map[string]string
// @Router       /evm/hash/keccak256/eip712 [post]
func (h *HashHandler) Keccak256EIP712(w http.ResponseWriter, r *http.Request) {
	req := new(Keccak256EIP712Request)
	if err := json.NewDecoder(r.Body).Decode(req); err != nil {
		handler.WriteError(w, http.StatusBadRequest, fmt.Sprintf("invalid request body: %s", err))
		return
	}
	if err := req.ValidateRequest(); err != nil {
		handler.WriteError(w, http.StatusBadRequest, err.Error())
		return
	}

	result, err := core.Hasher.EIP712(req.ToEIP712Domain(), req.ToFn(), req.Args)
	if err != nil {
		handler.WriteError(w, http.StatusBadRequest, fmt.Sprintf("failed to hash: %s", err))
		return
	}
	handler.WriteJSON(w, http.StatusOK, NewKeccak256EIP712Response(result))
}
