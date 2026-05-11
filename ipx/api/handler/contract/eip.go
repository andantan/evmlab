package contract

import (
	"encoding/hex"
	"encoding/json"
	"fmt"
	"math/big"
	"net/http"
	"strings"

	"github.com/andantan/evmlab/api/handler"
	"github.com/andantan/evmlab/core"
	"github.com/andantan/evmlab/internal/rpc"
	"github.com/andantan/evmlab/internal/util"
)

type EIPHandler struct {
	client *rpc.Client
}

func NewEIPHandler(client *rpc.Client) *EIPHandler {
	return &EIPHandler{
		client: client,
	}
}

// EIP712Domain godoc
// @Summary      Fetch EIP-712 domain
// @Description  Resolves name, version, chain_id, and verifying_contract for a contract
// @Tags         contract
// @Accept       json
// @Produce      json
// @Param        body  body      EIP712DomainRequest   true  "Contract address"
// @Success      200   {object}  EIP712DomainResponse
// @Failure      400   {object}  map[string]string
// @Failure      500   {object}  map[string]string
// @Router       /evm/contract/eip712/domain [post]
func (h *EIPHandler) EIP712Domain(w http.ResponseWriter, r *http.Request) {
	req := new(EIP712DomainRequest)
	if err := json.NewDecoder(r.Body).Decode(req); err != nil {
		handler.WriteError(w, http.StatusBadRequest, fmt.Sprintf("invalid request body: %s", err))
		return
	}
	if err := req.ValidateRequest(); err != nil {
		handler.WriteError(w, http.StatusBadRequest, err.Error())
		return
	}

	chainIDHex, err := h.client.ChainID(r.Context())
	if err != nil {
		handler.WriteError(w, http.StatusInternalServerError, fmt.Sprintf("eth_chainId: %s", err))
		return
	}
	chainID := new(big.Int)
	chainID.SetString(strings.TrimPrefix(chainIDHex, "0x"), 16)

	p := map[string]string{
		"to": req.Contract,
	}
	var (
		name string
		raw  string
		data []byte
	)

	p["data"] = "0x" + hex.EncodeToString(core.NameCalldata())
	if raw, err = h.client.CallContract(r.Context(), p, req.Block); err != nil {
		handler.WriteError(w, http.StatusInternalServerError, fmt.Sprintf("eth_call name() failed: %s", err))
		return
	}

	if data, err = util.ParseHex(raw); err != nil {
		handler.WriteError(w, http.StatusInternalServerError, fmt.Sprintf("failed to parse hex response: %s", err))
		return
	}

	if name, err = core.ABI.DecodeString(data); err != nil {
		handler.WriteError(w, http.StatusInternalServerError, fmt.Sprintf("failed to decode response: %s", err))
		return
	}

	version := "1"
	p["data"] = "0x" + hex.EncodeToString(core.VersionCalldata())
	if raw, err = h.client.CallContract(r.Context(), p, req.Block); err != nil {
		goto respond
	}
	if data, err = util.ParseHex(raw); err != nil {
		goto respond
	}
	if version, err = core.ABI.DecodeString(data); err != nil {
		goto respond
	}

respond:
	handler.WriteJSON(w, http.StatusOK, NewEIP712DomainResponse(name, version, chainID, req.Contract))
}
