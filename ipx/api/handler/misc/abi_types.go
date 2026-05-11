package misc

import (
	"encoding/hex"
	"errors"
	"math/big"
	"strings"

	"github.com/andantan/evmlab/core/types"
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
}

func (r *DecodeResultRequest) ValidateRequest() ([]byte, error) {
	if len(r.Types) == 0 {
		return nil, errors.New("types is required")
	}

	r.Data = strings.TrimSpace(r.Data)
	b, err := util.ParseHex(r.Data)
	if err != nil {
		return nil, errors.New("data: " + err.Error())
	}

	return b, nil
}

type DecodeResultResponse struct {
	Values []string `json:"values"`
}

func NewDecodeResultResponse(v []string) *DecodeResultResponse {
	return &DecodeResultResponse{
		Values: v,
	}
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
	b, err := util.ParseHex(r.Data)
	if err != nil {
		return nil, errors.New("data: " + err.Error())
	}

	return b, nil
}

type DecodeCallResponse struct {
	Selector string            `json:"selector"`
	Values   map[string]string `json:"values"`
}

func NewDecodeCallResponse(b []byte, v map[string]string) *DecodeCallResponse {
	return &DecodeCallResponse{
		Selector: "0x" + hex.EncodeToString(b[:4]),
		Values:   v,
	}
}

type ApproveCalldataRequest struct {
	Spender string `json:"spender" example:"0xDa70aA79f1a329719b9cf9d334b0a82b1d5269f3"`
	Amount  string `json:"amount"  example:"1000000000000000000"`

	s *types.Address
	a *big.Int
}

func (r *ApproveCalldataRequest) ValidateRequest() error {
	r.Spender = strings.TrimSpace(r.Spender)
	s, err := types.NewAddressFromHex(r.Spender)
	if err != nil {
		return errors.New("spender: invalid address")
	}
	r.s = s

	r.Amount = strings.TrimSpace(r.Amount)
	var ok bool
	r.a, ok = new(big.Int).SetString(r.Amount, 10)
	if !ok {
		return errors.New("amount: invalid integer")
	}
	return nil
}

func (r *ApproveCalldataRequest) ToSpender() *types.Address { return r.s }
func (r *ApproveCalldataRequest) ToAmount() *big.Int        { return r.a }

type ApproveCalldataResponse struct {
	Data string `json:"data"`
}

func NewApproveCalldataResponse(b []byte) *ApproveCalldataResponse {
	return &ApproveCalldataResponse{
		Data: "0x" + hex.EncodeToString(b),
	}
}

type TransferCalldataRequest struct {
	To     string `json:"to"     example:"0xDa70aA79f1a329719b9cf9d334b0a82b1d5269f3"`
	Amount string `json:"amount" example:"1000000000000000000"`

	t *types.Address
	a *big.Int
}

func (r *TransferCalldataRequest) ValidateRequest() error {
	r.To = strings.TrimSpace(r.To)
	t, err := types.NewAddressFromHex(r.To)
	if err != nil {
		return errors.New("to: invalid address")
	}
	r.t = t

	r.Amount = strings.TrimSpace(r.Amount)
	var ok bool
	r.a, ok = new(big.Int).SetString(r.Amount, 10)
	if !ok {
		return errors.New("amount: invalid integer")
	}
	return nil
}

func (r *TransferCalldataRequest) ToAddress() *types.Address { return r.t }
func (r *TransferCalldataRequest) ToAmount() *big.Int        { return r.a }

type TransferCalldataResponse struct {
	Data string `json:"data"`
}

func NewTransferCalldataResponse(b []byte) *TransferCalldataResponse {
	return &TransferCalldataResponse{
		Data: "0x" + hex.EncodeToString(b),
	}
}

type AllowanceCalldataRequest struct {
	Owner   string `json:"owner"   example:"0xDa70aA79f1a329719b9cf9d334b0a82b1d5269f3"`
	Spender string `json:"spender" example:"0xDa70aA79f1a329719b9cf9d334b0a82b1d5269f3"`

	o *types.Address
	s *types.Address
}

func (r *AllowanceCalldataRequest) ValidateRequest() error {
	r.Owner = strings.TrimSpace(r.Owner)
	o, err := types.NewAddressFromHex(r.Owner)
	if err != nil {
		return errors.New("owner: invalid address")
	}
	r.o = o

	r.Spender = strings.TrimSpace(r.Spender)
	s, err := types.NewAddressFromHex(r.Spender)
	if err != nil {
		return errors.New("spender: invalid address")
	}
	r.s = s

	return nil
}

func (r *AllowanceCalldataRequest) ToOwner() *types.Address   { return r.o }
func (r *AllowanceCalldataRequest) ToSpender() *types.Address { return r.s }

type AllowanceCalldataResponse struct {
	Data string `json:"data"`
}

func NewAllowanceCalldataResponse(b []byte) *AllowanceCalldataResponse {
	return &AllowanceCalldataResponse{
		Data: "0x" + hex.EncodeToString(b),
	}
}

type BalanceOfCalldataRequest struct {
	Account string `json:"account" example:"0xDa70aA79f1a329719b9cf9d334b0a82b1d5269f3"`

	a *types.Address
}

func (r *BalanceOfCalldataRequest) ValidateRequest() error {
	r.Account = strings.TrimSpace(r.Account)
	a, err := types.NewAddressFromHex(r.Account)
	if err != nil {
		return errors.New("account: invalid address")
	}
	r.a = a
	return nil
}

func (r *BalanceOfCalldataRequest) ToAccount() *types.Address { return r.a }

type BalanceOfCalldataResponse struct {
	Data string `json:"data"`
}

func NewBalanceOfCalldataResponse(b []byte) *BalanceOfCalldataResponse {
	return &BalanceOfCalldataResponse{
		Data: "0x" + hex.EncodeToString(b),
	}
}

type TransferFromCalldataRequest struct {
	From   string `json:"from"   example:"0xDa70aA79f1a329719b9cf9d334b0a82b1d5269f3"`
	To     string `json:"to"     example:"0xDa70aA79f1a329719b9cf9d334b0a82b1d5269f3"`
	Amount string `json:"amount" example:"1000000000000000000"`

	f *types.Address
	t *types.Address
	a *big.Int
}

func (r *TransferFromCalldataRequest) ValidateRequest() error {
	r.From = strings.TrimSpace(r.From)
	f, err := types.NewAddressFromHex(r.From)
	if err != nil {
		return errors.New("from: invalid address")
	}
	r.f = f

	r.To = strings.TrimSpace(r.To)
	t, err := types.NewAddressFromHex(r.To)
	if err != nil {
		return errors.New("to: invalid address")
	}
	r.t = t

	r.Amount = strings.TrimSpace(r.Amount)
	var ok bool
	r.a, ok = new(big.Int).SetString(r.Amount, 10)
	if !ok {
		return errors.New("amount: invalid integer")
	}
	return nil
}

func (r *TransferFromCalldataRequest) ToFrom() *types.Address { return r.f }
func (r *TransferFromCalldataRequest) ToTo() *types.Address   { return r.t }
func (r *TransferFromCalldataRequest) ToAmount() *big.Int     { return r.a }

type TransferFromCalldataResponse struct {
	Data string `json:"data"`
}

func NewTransferFromCalldataResponse(b []byte) *TransferFromCalldataResponse {
	return &TransferFromCalldataResponse{
		Data: "0x" + hex.EncodeToString(b),
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

type NameCalldataResponse struct {
	Data string `json:"data"`
}

func NewNameCalldataResponse(b []byte) *NameCalldataResponse {
	return &NameCalldataResponse{
		Data: "0x" + hex.EncodeToString(b),
	}
}

type VersionCalldataResponse struct {
	Data string `json:"data"`
}

func NewVersionCalldataResponse(b []byte) *VersionCalldataResponse {
	return &VersionCalldataResponse{
		Data: "0x" + hex.EncodeToString(b),
	}
}
