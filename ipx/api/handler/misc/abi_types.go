package misc

import (
	"encoding/hex"
	"errors"
	"strings"

	"github.com/andantan/evmlab/internal/util"
)

type SelectorRequest struct {
	Signature string `json:"signature" example:"transfer(address,uint256)"`
}

func (r *SelectorRequest) ValidateRequest() error {
	r.Signature = strings.TrimSpace(r.Signature)
	if r.Signature == "" {
		return errors.New("signature is required")
	}
	return nil
}

type SelectorResponse struct {
	Selector string `json:"selector"`
}

func NewSelectorResponse(s []byte) *SelectorResponse {
	return &SelectorResponse{
		Selector: "0x" + hex.EncodeToString(s),
	}
}

type EncodeRequest struct {
	Signature string   `json:"signature" example:"transfer(address,uint256)"`
	Args      []string `json:"args"      example:"[\"0xDa70aA79...\",\"1000000000000000000\"]"`
}

func (r *EncodeRequest) ValidateRequest() error {
	r.Signature = strings.TrimSpace(r.Signature)
	if r.Signature == "" {
		return errors.New("signature is required")
	}

	return nil
}

type EncodeResponse struct {
	Data string `json:"data"`
}

func NewEncodeResponse(b []byte) *EncodeResponse {
	return &EncodeResponse{
		Data: "0x" + hex.EncodeToString(b),
	}
}

type DecodeResultRequest struct {
	Data  string   `json:"data"  example:"0x000000000000000000000000000000000000000000000000000000000001cf1d"`
	Types []string `json:"types" example:"[\"uint256\"]"`

	d []byte
}

func (r *DecodeResultRequest) ValidateRequest() error {
	if len(r.Types) == 0 {
		return errors.New("types is required")
	}

	r.Data = strings.TrimSpace(r.Data)
	b, err := util.ParseHex(r.Data)
	if err != nil {
		return errors.New("data: " + err.Error())
	}

	r.d = b
	return nil
}

func (r *DecodeResultRequest) ToData() []byte { return r.d }

type DecodeResultResponse struct {
	Values []any `json:"values"`
}

func NewDecodeResultResponse(v []any) *DecodeResultResponse {
	return &DecodeResultResponse{
		Values: v,
	}
}

type DecodeCallRequest struct {
	Signature string `json:"signature" example:"transfer(address,uint256)"`
	Data      string `json:"data"      example:"0xa9059cbb000000000000000000000000da70aa79f1a329719b9cf9d334b0a82b1d5269f300000000000000000000000000000000000000000000000000000000000003e8"`

	d []byte
}

func (r *DecodeCallRequest) ValidateRequest() error {
	r.Signature = strings.TrimSpace(r.Signature)
	if r.Signature == "" {
		return errors.New("signature is required")
	}

	r.Data = strings.TrimSpace(r.Data)
	b, err := util.ParseHex(r.Data)
	if err != nil {
		return errors.New("data: " + err.Error())
	}

	r.d = b
	return nil
}

func (r *DecodeCallRequest) ToData() []byte { return r.d }

type DecodeCallResponse struct {
	Selector string         `json:"selector"`
	Values   map[string]any `json:"values"`
}

func NewDecodeCallResponse(b []byte, v map[string]any) *DecodeCallResponse {
	return &DecodeCallResponse{
		Selector: "0x" + hex.EncodeToString(b[:4]),
		Values:   v,
	}
}

type DecodeRevertRequest struct {
	Signature string `json:"signature" example:"InsufficientBalance(address,uint256,uint256)"`
	Data      string `json:"data"      example:"0x1d2a3b4c..."`

	d []byte
}

func (r *DecodeRevertRequest) ValidateRequest() error {
	r.Signature = strings.TrimSpace(r.Signature)
	if r.Signature == "" {
		return errors.New("signature is required")
	}

	r.Data = strings.TrimSpace(r.Data)
	b, err := util.ParseHex(r.Data)
	if err != nil {
		return errors.New("data: " + err.Error())
	}

	r.d = b
	return nil
}

func (r *DecodeRevertRequest) ToData() []byte { return r.d }

type DecodeRevertResponse struct {
	Error  string         `json:"error"`
	Values map[string]any `json:"values"`
}

func NewDecodeRevertResponse(name string, v map[string]any) *DecodeRevertResponse {
	return &DecodeRevertResponse{
		Error:  name,
		Values: v,
	}
}

type EIP712DomainCalldataResponse struct {
	Data string `json:"data"`
}

func NewEIP712DomainCalldataResponse(b []byte) *EIP712DomainCalldataResponse {
	return &EIP712DomainCalldataResponse{
		Data: "0x" + hex.EncodeToString(b),
	}
}
