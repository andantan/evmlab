package contract

import (
	"encoding/hex"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/andantan/evmlab/api/handler"
	"github.com/andantan/evmlab/core"
	"github.com/andantan/evmlab/internal/rpc"
	"github.com/andantan/evmlab/internal/util"
)

type ERC20Handler struct {
	client *rpc.Client
}

func NewERC20Handler(client *rpc.Client) *ERC20Handler {
	return &ERC20Handler{client: client}
}

// Metadata godoc
// @Summary      Fetch ERC-20 token metadata
// @Description  Returns name, symbol, and decimals for an ERC-20 contract
// @Tags         contract
// @Accept       json
// @Produce      json
// @Param        body  body      ERC20MetadataRequest   true  "Contract address"
// @Success      200   {object}  ERC20MetadataResponse
// @Failure      400   {object}  map[string]string
// @Failure      500   {object}  map[string]string
// @Router       /evm/contract/erc20/metadata [post]
func (h *ERC20Handler) Metadata(w http.ResponseWriter, r *http.Request) {
	req := new(ERC20MetadataRequest)
	if err := json.NewDecoder(r.Body).Decode(req); err != nil {
		handler.WriteError(w, http.StatusBadRequest, fmt.Sprintf("invalid request body: %s", err))
		return
	}
	if err := req.ValidateRequest(); err != nil {
		handler.WriteError(w, http.StatusBadRequest, err.Error())
		return
	}

	p := map[string]string{"to": req.Contract}
	var (
		raw  string
		data []byte
		err  error
	)

	ok, err := h.client.IsContract(r.Context(), req.Contract, req.Block)
	if err != nil {
		handler.WriteError(w, http.StatusInternalServerError, fmt.Sprintf("eth_getCode failed: %s", err))
		return
	}
	if !ok {
		handler.WriteError(w, http.StatusUnprocessableEntity, "requested contract is an EOA: metadata not available")
		return
	}

	p["data"] = "0x" + hex.EncodeToString(core.NameCalldata())
	if raw, err = h.client.CallContract(r.Context(), p, req.Block); err != nil {
		handler.WriteError(w, http.StatusInternalServerError, fmt.Sprintf("eth_call name() failed: %s", err))
		return
	}
	if data, err = util.ParseHex(raw); err != nil {
		handler.WriteError(w, http.StatusInternalServerError, fmt.Sprintf("failed to parse hex response: %s", err))
		return
	}
	name, err := core.ABI.DecodeString(data)
	if err != nil {
		handler.WriteError(w, http.StatusInternalServerError, fmt.Sprintf("failed to decode name: %s", err))
		return
	}

	p["data"] = "0x" + hex.EncodeToString(core.SymbolCalldata())
	if raw, err = h.client.CallContract(r.Context(), p, req.Block); err != nil {
		handler.WriteError(w, http.StatusInternalServerError, fmt.Sprintf("eth_call symbol() failed: %s", err))
		return
	}
	if data, err = util.ParseHex(raw); err != nil {
		handler.WriteError(w, http.StatusInternalServerError, fmt.Sprintf("failed to parse hex response: %s", err))
		return
	}
	symbol, err := core.ABI.DecodeString(data)
	if err != nil {
		handler.WriteError(w, http.StatusInternalServerError, fmt.Sprintf("failed to decode symbol: %s", err))
		return
	}

	p["data"] = "0x" + hex.EncodeToString(core.DecimalsCalldata())
	if raw, err = h.client.CallContract(r.Context(), p, req.Block); err != nil {
		handler.WriteError(w, http.StatusInternalServerError, fmt.Sprintf("eth_call decimals() failed: %s", err))
		return
	}
	if data, err = util.ParseHex(raw); err != nil {
		handler.WriteError(w, http.StatusInternalServerError, fmt.Sprintf("failed to parse hex response: %s", err))
		return
	}
	decimals, err := core.ABI.DecodeUint8(data)
	if err != nil {
		handler.WriteError(w, http.StatusInternalServerError, fmt.Sprintf("failed to decode decimals: %s", err))
		return
	}

	handler.WriteJSON(w, http.StatusOK, NewERC20MetadataResponse(name, symbol, decimals))
}

// Balance godoc
// @Summary      Fetch ERC-20 token balance
// @Description  Returns the formatted token balance of an account
// @Tags         contract
// @Accept       json
// @Produce      json
// @Param        body  body      ERC20BalanceRequest   true  "Contract and account address"
// @Success      200   {object}  ERC20BalanceResponse
// @Failure      400   {object}  map[string]string
// @Failure      500   {object}  map[string]string
// @Router       /evm/contract/erc20/balance [post]
func (h *ERC20Handler) Balance(w http.ResponseWriter, r *http.Request) {
	req := new(ERC20BalanceRequest)
	if err := json.NewDecoder(r.Body).Decode(req); err != nil {
		handler.WriteError(w, http.StatusBadRequest, fmt.Sprintf("invalid request body: %s", err))
		return
	}
	if err := req.ValidateRequest(); err != nil {
		handler.WriteError(w, http.StatusBadRequest, err.Error())
		return
	}

	p := map[string]string{"to": req.Contract}
	var (
		raw  string
		data []byte
		err  error
	)

	ok, err := h.client.IsContract(r.Context(), req.Contract, req.Block)
	if err != nil {
		handler.WriteError(w, http.StatusInternalServerError, fmt.Sprintf("eth_getCode failed: %s", err))
		return
	}
	if !ok {
		handler.WriteError(w, http.StatusUnprocessableEntity, "requested contract is an EOA: metadata not available")
		return
	}

	p["data"] = "0x" + hex.EncodeToString(core.DecimalsCalldata())
	if raw, err = h.client.CallContract(r.Context(), p, req.Block); err != nil {
		handler.WriteError(w, http.StatusInternalServerError, fmt.Sprintf("eth_call decimals() failed: %s", err))
		return
	}
	if data, err = util.ParseHex(raw); err != nil {
		handler.WriteError(w, http.StatusInternalServerError, fmt.Sprintf("failed to parse hex response: %s", err))
		return
	}
	decimals, err := core.ABI.DecodeUint8(data)
	if err != nil {
		handler.WriteError(w, http.StatusInternalServerError, fmt.Sprintf("failed to decode decimals: %s", err))
		return
	}

	p["data"] = "0x" + hex.EncodeToString(core.BalanceOfCalldata(req.ToAccount()))
	if raw, err = h.client.CallContract(r.Context(), p, req.Block); err != nil {
		handler.WriteError(w, http.StatusInternalServerError, fmt.Sprintf("eth_call balanceOf() failed: %s", err))
		return
	}
	if data, err = util.ParseHex(raw); err != nil {
		handler.WriteError(w, http.StatusInternalServerError, fmt.Sprintf("failed to parse hex response: %s", err))
		return
	}
	amount, err := core.ABI.DecodeUint256(data)
	if err != nil {
		handler.WriteError(w, http.StatusInternalServerError, fmt.Sprintf("failed to decode balance: %s", err))
		return
	}

	handler.WriteJSON(w, http.StatusOK, NewERC20BalanceResponse(util.FormatTokenAmount(amount, decimals)))
}

// Allowance godoc
// @Summary      Fetch ERC-20 token allowance
// @Description  Returns the formatted token allowance of a spender for an owner
// @Tags         contract
// @Accept       json
// @Produce      json
// @Param        body  body      ERC20AllowanceRequest   true  "Contract, owner, and spender address"
// @Success      200   {object}  ERC20AllowanceResponse
// @Failure      400   {object}  map[string]string
// @Failure      500   {object}  map[string]string
// @Router       /evm/contract/erc20/allowance [post]
func (h *ERC20Handler) Allowance(w http.ResponseWriter, r *http.Request) {
	req := new(ERC20AllowanceRequest)
	if err := json.NewDecoder(r.Body).Decode(req); err != nil {
		handler.WriteError(w, http.StatusBadRequest, fmt.Sprintf("invalid request body: %s", err))
		return
	}
	if err := req.ValidateRequest(); err != nil {
		handler.WriteError(w, http.StatusBadRequest, err.Error())
		return
	}

	p := map[string]string{"to": req.Contract}
	var (
		raw  string
		data []byte
		err  error
	)

	ok, err := h.client.IsContract(r.Context(), req.Contract, req.Block)
	if err != nil {
		handler.WriteError(w, http.StatusInternalServerError, fmt.Sprintf("eth_getCode failed: %s", err))
		return
	}
	if !ok {
		handler.WriteError(w, http.StatusUnprocessableEntity, "requested contract is an EOA: metadata not available")
		return
	}

	p["data"] = "0x" + hex.EncodeToString(core.DecimalsCalldata())
	if raw, err = h.client.CallContract(r.Context(), p, req.Block); err != nil {
		handler.WriteError(w, http.StatusInternalServerError, fmt.Sprintf("eth_call decimals() failed: %s", err))
		return
	}
	if data, err = util.ParseHex(raw); err != nil {
		handler.WriteError(w, http.StatusInternalServerError, fmt.Sprintf("failed to parse hex response: %s", err))
		return
	}
	decimals, err := core.ABI.DecodeUint8(data)
	if err != nil {
		handler.WriteError(w, http.StatusInternalServerError, fmt.Sprintf("failed to decode decimals: %s", err))
		return
	}

	p["data"] = "0x" + hex.EncodeToString(core.AllowanceCalldata(req.ToOwner(), req.ToSpender()))
	if raw, err = h.client.CallContract(r.Context(), p, req.Block); err != nil {
		handler.WriteError(w, http.StatusInternalServerError, fmt.Sprintf("eth_call allowance() failed: %s", err))
		return
	}
	if data, err = util.ParseHex(raw); err != nil {
		handler.WriteError(w, http.StatusInternalServerError, fmt.Sprintf("failed to parse hex response: %s", err))
		return
	}
	amount, err := core.ABI.DecodeUint256(data)
	if err != nil {
		handler.WriteError(w, http.StatusInternalServerError, fmt.Sprintf("failed to decode allowance: %s", err))
		return
	}

	handler.WriteJSON(w, http.StatusOK, NewERC20AllowanceResponse(util.FormatTokenAmount(amount, decimals)))
}

// Approved godoc
// @Summary      Simulate ERC-20 approve
// @Description  Simulates approve(address,uint256) via eth_call and returns the bool result
// @Tags         contract
// @Accept       json
// @Produce      json
// @Param        body  body      ERC20ApproveRequest   true  "Contract, spender, and value"
// @Success      200   {object}  ERC20ApproveResponse
// @Failure      400   {object}  map[string]string
// @Failure      500   {object}  map[string]string
// @Router       /evm/contract/erc20/approved [post]
func (h *ERC20Handler) Approved(w http.ResponseWriter, r *http.Request) {
	req := new(ERC20ApproveRequest)
	if err := json.NewDecoder(r.Body).Decode(req); err != nil {
		handler.WriteError(w, http.StatusBadRequest, fmt.Sprintf("invalid request body: %s", err))
		return
	}
	if err := req.ValidateRequest(); err != nil {
		handler.WriteError(w, http.StatusBadRequest, err.Error())
		return
	}

	ok, err := h.client.IsContract(r.Context(), req.Contract, req.Block)
	if err != nil {
		handler.WriteError(w, http.StatusInternalServerError, fmt.Sprintf("eth_getCode failed: %s", err))
		return
	}
	if !ok {
		handler.WriteError(w, http.StatusUnprocessableEntity, "requested contract is an EOA: metadata not available")
		return
	}

	p := map[string]string{
		"to":   req.Contract,
		"data": "0x" + hex.EncodeToString(core.ApproveCalldata(req.ToSpender(), req.ToValue())),
	}

	raw, err := h.client.CallContract(r.Context(), p, req.Block)
	if err != nil {
		handler.WriteError(w, http.StatusInternalServerError, fmt.Sprintf("eth_call approve() failed: %s", err))
		return
	}

	data, err := util.ParseHex(raw)
	if err != nil {
		handler.WriteError(w, http.StatusInternalServerError, fmt.Sprintf("failed to parse hex response: %s", err))
		return
	}

	approved, err := core.ABI.DecodeBool(data)
	if err != nil {
		handler.WriteError(w, http.StatusInternalServerError, fmt.Sprintf("failed to decode response: %s", err))
		return
	}

	handler.WriteJSON(w, http.StatusOK, NewERC20ApproveResponse(approved))
}
