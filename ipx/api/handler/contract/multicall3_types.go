package contract

import (
	"encoding/hex"
	"errors"
	"fmt"
	"strings"

	"github.com/andantan/evmlab/core/types"
	"github.com/andantan/evmlab/internal/config"
)

type Multicall3Call struct {
	To   string `json:"to"`
	Data string `json:"data"`
}

type Multicall3Aggregate3Request struct {
	Target string           `json:"target"`
	Calls  []Multicall3Call `json:"calls"`

	calls []types.Aggregate3
}

func (r *Multicall3Aggregate3Request) ValidateRequest() error {
	r.Target = strings.TrimSpace(r.Target)
	if r.Target == "" {
		r.Target = config.Multicall3CanonicalAddress
	}

	if len(r.Calls) == 0 {
		return errors.New("calls is required")
	}

	r.calls = make([]types.Aggregate3, len(r.Calls))
	for i, c := range r.Calls {
		addr, err := types.NewAddressFromHex(c.To)
		if err != nil {
			return fmt.Errorf("calls[%d]: to: invalid address", i)
		}
		callData, err := hex.DecodeString(strings.TrimPrefix(strings.TrimSpace(c.Data), "0x"))
		if err != nil {
			return fmt.Errorf("calls[%d]: data: %s", i, err)
		}
		r.calls[i] = types.Aggregate3{Target: addr, AllowFail: true, CallData: callData}
	}
	return nil
}

func (r *Multicall3Aggregate3Request) ToCalls() []types.Aggregate3 { return r.calls }

type Multicall3Aggregate3Result struct {
	Success    bool   `json:"success"`
	ReturnData string `json:"return_data"`
}

type Multicall3Aggregate3Response struct {
	Results []Multicall3Aggregate3Result `json:"results"`
}

func NewMulticall3Aggregate3Response(decoded []types.Aggregator3Result) *Multicall3Aggregate3Response {
	results := make([]Multicall3Aggregate3Result, len(decoded))
	for i, d := range decoded {
		results[i] = Multicall3Aggregate3Result{
			Success:    d.Success,
			ReturnData: "0x" + hex.EncodeToString(d.ReturnData),
		}
	}
	return &Multicall3Aggregate3Response{Results: results}
}
