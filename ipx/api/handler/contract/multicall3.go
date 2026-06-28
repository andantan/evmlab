package contract

import (
	"encoding/hex"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/andantan/evmlab/api/handler"
	"github.com/andantan/evmlab/core"
	"github.com/andantan/evmlab/core/types"
	"github.com/andantan/evmlab/internal/config"
	"github.com/andantan/evmlab/internal/rpc"
	"github.com/andantan/evmlab/internal/util"
)

type Multicall3Handler struct {
	client         *rpc.Client
	multicall3Addr *types.Address
}

func NewMulticall3Handler(cfg *config.Config, client *rpc.Client) *Multicall3Handler {
	if cfg.Multicall3 == "" {
		panic("multicall3 address not configured")
	}

	addr, err := types.NewAddressFromHex(cfg.Multicall3)
	if err != nil {
		panic(fmt.Sprintf("multicall3: invalid address: %s", err))
	}

	return &Multicall3Handler{
		client:         client,
		multicall3Addr: addr,
	}
}

// Aggregate3 godoc
// @Summary      Multicall3 aggregate3
// @Description  Batch multiple eth_call invocations via Multicall3 aggregate3. allow_failure is always true.
// @Tags         contract
// @Accept       json
// @Produce      json
// @Param        body  body      Multicall3Aggregate3Request   true  "Multicall3 calls"
// @Success      200   {object}  Multicall3Aggregate3Response
// @Failure      400   {object}  map[string]string
// @Failure      502   {object}  map[string]string
// @Router       /evm/contract/multicall3/aggregate3 [post]
func (h *Multicall3Handler) Aggregate3(w http.ResponseWriter, r *http.Request) {
	req := new(Multicall3Aggregate3Request)
	if err := json.NewDecoder(r.Body).Decode(req); err != nil {
		handler.WriteError(w, http.StatusBadRequest, fmt.Sprintf("invalid request body: %s", err))
		return
	}
	if err := req.ValidateRequest(); err != nil {
		handler.WriteError(w, http.StatusBadRequest, err.Error())
		return
	}

	p := map[string]string{
		"to":   h.multicall3Addr.String(),
		"data": "0x" + hex.EncodeToString(core.Multicall3Aggregator3CallData(req.ToCalls())),
	}

	raw, err := h.client.CallContract(r.Context(), p, "latest")
	if err != nil {
		handler.WriteError(w, http.StatusBadGateway, fmt.Sprintf("eth_call failed: %s", err))
		return
	}

	resultBytes, err := util.ParseHex(raw)
	if err != nil {
		handler.WriteError(w, http.StatusInternalServerError, fmt.Sprintf("parse response: %s", err))
		return
	}

	decoded, err := types.DecodeAggregate3Results(resultBytes)
	if err != nil {
		handler.WriteError(w, http.StatusInternalServerError, fmt.Sprintf("decode results: %s", err))
		return
	}

	handler.WriteJSON(w, http.StatusOK, NewMulticall3Aggregate3Response(decoded))
}
