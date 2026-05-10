package misc

import (
	"encoding/hex"
	"errors"
	"strings"

	"github.com/ethereum/go-ethereum/common/hexutil"
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

func NewSelectorResponse(selector []byte) *SelectorResponse {
	return &SelectorResponse{
		Selector: "0x" + hex.EncodeToString(selector),
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

func NewEncodeResponse(data []byte) *EncodeResponse {
	return &EncodeResponse{
		Data: "0x" + hex.EncodeToString(data),
	}
}

type DecodeResultRequest struct {
	Data  string   `json:"data"  example:"0x000000000000000000000000000000000000000000000000000000000001cf1d"`
	Types []string `json:"types" example:"[\"uint256\"]"`
}

func (r *DecodeResultRequest) ValidateRequest() ([]byte, error) {
	if len(r.Types) == 0 {
		return nil, errors.New("types is required")
	}
	r.Data = strings.TrimSpace(r.Data)
	b, err := hexutil.Decode(r.Data)
	if err != nil {
		return nil, errors.New("data: " + err.Error())
	}
	return b, nil
}

type DecodeResultResponse struct {
	Values []string `json:"values"`
}

func NewDecodeResultResponse(values []string) *DecodeResultResponse {
	return &DecodeResultResponse{Values: values}
}

type DecodeCallRequest struct {
	Signature string `json:"signature" example:"transfer(address,uint256)"`
	Data      string `json:"data"      example:"0xa9059cbb000000000000000000000000da70aa79f1a329719b9cf9d334b0a82b1d5269f300000000000000000000000000000000000000000000000000000000000003e8"`
}

func (r *DecodeCallRequest) ValidateRequest() ([]byte, error) {
	r.Signature = strings.TrimSpace(r.Signature)
	if r.Signature == "" {
		return nil, errors.New("signature is required")
	}
	r.Data = strings.TrimSpace(r.Data)
	b, err := hexutil.Decode(r.Data)
	if err != nil {
		return nil, errors.New("data: " + err.Error())
	}
	return b, nil
}

type DecodeCallResponse struct {
	Selector string            `json:"selector"`
	Values   map[string]string `json:"values"`
}

func NewDecodeCallResponse(data []byte, values map[string]string) *DecodeCallResponse {
	return &DecodeCallResponse{
		Selector: "0x" + hex.EncodeToString(data[:4]),
		Values:   values,
	}
}
