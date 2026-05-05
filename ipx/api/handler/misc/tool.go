package misc

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/andantan/evmlab/api/handler"
	"github.com/andantan/evmlab/core"
	"github.com/andantan/evmlab/core/types"
)

type ToolHandler struct{}

func NewToolHandler() *ToolHandler {
	return &ToolHandler{}
}

// ChecksumEIP55 godoc
// @Summary      Convert address to EIP-55 checksum format
// @Description  Returns the EIP-55 mixed-case checksum encoding for the given address
// @Tags         tool
// @Accept       json
// @Produce      json
// @Param        body  body      ChecksumEIP55Request  true  "Address"
// @Success      200   {object}  ChecksumEIP55Response
// @Failure      400   {object}  map[string]string
// @Router       /evm/tool/address/checksum/eip55 [post]
func (h *ToolHandler) ChecksumEIP55(w http.ResponseWriter, r *http.Request) {
	req := new(ChecksumEIP55Request)
	if err := json.NewDecoder(r.Body).Decode(req); err != nil {
		handler.WriteError(w, http.StatusBadRequest, fmt.Sprintf("invalid request body: %s", err))
		return
	}
	if err := req.ValidateRequest(); err != nil {
		handler.WriteError(w, http.StatusBadRequest, err.Error())
		return
	}

	handler.WriteJSON(w, http.StatusOK, NewChecksumEIP55Response(req.ToAddress()))
}

// DeriveKey godoc
// @Summary      Derive key set from private key
// @Description  Returns the public key and address derived from the given private key
// @Tags         tool
// @Accept       json
// @Produce      json
// @Param        body  body      DeriveKeyRequest   true  "Private key"
// @Success      200   {object}  DeriveKeyResponse
// @Failure      400   {object}  map[string]string
// @Router       /evm/tool/crypto/derive [post]
func (h *ToolHandler) DeriveKey(w http.ResponseWriter, r *http.Request) {
	req := new(DeriveKeyRequest)
	if err := json.NewDecoder(r.Body).Decode(req); err != nil {
		handler.WriteError(w, http.StatusBadRequest, fmt.Sprintf("invalid request body: %s", err))
		return
	}
	if err := req.ValidateRequest(); err != nil {
		handler.WriteError(w, http.StatusBadRequest, err.Error())
		return
	}

	key, err := core.DeriveKeyFromPrivHex(req.PrivateKey)
	if err != nil {
		handler.WriteError(w, http.StatusBadRequest, fmt.Sprintf("invalid private_key: %s", err))
		return
	}

	handler.WriteJSON(w, http.StatusOK, NewDeriveKeyResponse(key))
}

// ConvertUnitDecimal godoc
// @Summary      Convert amount between wei, gwei, and ether
// @Description  Converts a decimal amount from one Ethereum unit to another
// @Tags         tool
// @Accept       json
// @Produce      json
// @Param        body  body      UnitConvertDecimalRequest   true  "Unit conversion"
// @Success      200   {object}  UnitConvertDecimalResponse
// @Failure      400   {object}  map[string]string
// @Router       /evm/tool/unit/convert/decimal [post]
func (h *ToolHandler) ConvertUnitDecimal(w http.ResponseWriter, r *http.Request) {
	req := new(UnitConvertDecimalRequest)
	if err := json.NewDecoder(r.Body).Decode(req); err != nil {
		handler.WriteError(w, http.StatusBadRequest, fmt.Sprintf("invalid request body: %s", err))
		return
	}
	if err := req.ValidateRequest(); err != nil {
		handler.WriteError(w, http.StatusBadRequest, err.Error())
		return
	}

	amount, err := types.ConvertUnitDecimal(req.Amount, req.From, req.To)
	if err != nil {
		handler.WriteError(w, http.StatusBadRequest, err.Error())
		return
	}

	handler.WriteJSON(w, http.StatusOK, NewUnitConvertDecimalResponse(amount, req.To))
}
