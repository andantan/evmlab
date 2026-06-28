package sanctum

import (
	"fmt"
	"math/big"
	"net/http"
	"strconv"

	"github.com/andantan/evmlab/api/handler"
	"github.com/andantan/evmlab/core"
	"github.com/andantan/evmlab/core/types"
	"github.com/andantan/evmlab/internal/config"
	"github.com/andantan/evmlab/internal/rpc"
	"github.com/andantan/evmlab/internal/util"
)

type SanctumHandler struct {
	cfg    *config.Config
	client *rpc.Client
}

func NewSanctumHandler(cfg *config.Config, client *rpc.Client) *SanctumHandler {
	return &SanctumHandler{cfg: cfg, client: client}
}

// writeRevertOrGatewayError writes a 400 with a decoded Sanctum revert message if the error
// is a known contract revert, otherwise writes a 502 with the provided context prefix.
func (h *SanctumHandler) writeRevertOrGatewayError(w http.ResponseWriter, err error, context string) {
	d, ok := util.RevertData(err)
	if !ok {
		handler.WriteError(w, http.StatusBadGateway, fmt.Sprintf("%s: %s", context, err))
		return
	}

	fn, params, e := core.ABI.DecodeErrorData(d, types.SanctumErrorSignatures)
	if e != nil {
		handler.WriteError(w, http.StatusBadGateway, fmt.Sprintf("%s: %s", context, err))
		return
	}

	body := map[string]any{"error": fn.Name}
	if len(params) > 0 {
		body["detail"] = params
	}
	handler.WriteJSON(w, http.StatusBadRequest, body)
}

func (h *SanctumHandler) decodeAccountInfo(b []byte) (*AccountInfo, error) {
	values, err := core.ABI.DecodeResult(accountInfoTypes[:], b)
	if err != nil {
		return nil, err
	}

	addr, err := types.NewAddressFromHex(values[0].(string))
	if err != nil {
		return nil, fmt.Errorf("parse addr: %s", err)
	}

	roleVal, err := strconv.ParseUint(values[1].(string), 10, 8)
	if err != nil {
		return nil, fmt.Errorf("parse role: %s", err)
	}

	block := new(big.Int)
	if _, ok := block.SetString(values[2].(string), 10); !ok {
		return nil, fmt.Errorf("parse block: invalid uint256")
	}

	return &AccountInfo{
		Addr:            addr,
		Role:            SanctumRole(roleVal),
		RegisteredBlock: block,
	}, nil
}
