package contract

import (
	"errors"
	"math/big"
	"strings"

	"github.com/andantan/evmlab/core/types"
	"github.com/andantan/evmlab/internal/util"
)

type ERC20MetadataRequest struct {
	Contract string `json:"contract" example:"0x8336c196..."`
	Block    string `json:"block"    example:"latest"`
}

func (r *ERC20MetadataRequest) ValidateRequest() error {
	r.Contract = strings.TrimSpace(r.Contract)
	if !util.IsHexAddress(r.Contract) {
		return errors.New("contract: invalid address")
	}
	if r.Block == "" {
		r.Block = "latest"
	}
	return nil
}

type ERC20MetadataResponse struct {
	Name     string `json:"name"`
	Symbol   string `json:"symbol"`
	Decimals uint8  `json:"decimals"`
}

func NewERC20MetadataResponse(n, s string, d uint8) *ERC20MetadataResponse {
	return &ERC20MetadataResponse{
		Name:     n,
		Symbol:   s,
		Decimals: d,
	}
}

type ERC20BalanceRequest struct {
	Contract string `json:"contract" example:"0x8336c196..."`
	Account  string `json:"account"  example:"0xAbcD1234..."`
	Block    string `json:"block"    example:"latest"`

	account *types.Address
}

func (r *ERC20BalanceRequest) ValidateRequest() error {
	r.Contract = strings.TrimSpace(r.Contract)
	if !util.IsHexAddress(r.Contract) {
		return errors.New("contract: invalid address")
	}
	r.Account = strings.TrimSpace(r.Account)
	var err error
	if r.account, err = types.NewAddressFromHex(r.Account); err != nil {
		return errors.New("account: invalid address")
	}
	if r.Block == "" {
		r.Block = "latest"
	}
	return nil
}

func (r *ERC20BalanceRequest) ToAccount() *types.Address { return r.account }

type ERC20BalanceResponse struct {
	Balance string `json:"balance"`
}

func NewERC20BalanceResponse(s string) *ERC20BalanceResponse {
	return &ERC20BalanceResponse{
		Balance: s,
	}
}

type ERC20AllowanceRequest struct {
	Contract string `json:"contract" example:"0x8336c196..."`
	Owner    string `json:"owner"    example:"0xAbcD1234..."`
	Spender  string `json:"spender"  example:"0xAbcD1234..."`
	Block    string `json:"block"    example:"latest"`

	owner   *types.Address
	spender *types.Address
}

func (r *ERC20AllowanceRequest) ValidateRequest() error {
	r.Contract = strings.TrimSpace(r.Contract)
	if !util.IsHexAddress(r.Contract) {
		return errors.New("contract: invalid address")
	}
	r.Owner = strings.TrimSpace(r.Owner)
	var err error
	if r.owner, err = types.NewAddressFromHex(r.Owner); err != nil {
		return errors.New("owner: invalid address")
	}
	r.Spender = strings.TrimSpace(r.Spender)
	if r.spender, err = types.NewAddressFromHex(r.Spender); err != nil {
		return errors.New("spender: invalid address")
	}
	if r.Block == "" {
		r.Block = "latest"
	}
	return nil
}

func (r *ERC20AllowanceRequest) ToOwner() *types.Address   { return r.owner }
func (r *ERC20AllowanceRequest) ToSpender() *types.Address { return r.spender }

type ERC20AllowanceResponse struct {
	Allowance string `json:"allowance"`
}

func NewERC20AllowanceResponse(s string) *ERC20AllowanceResponse {
	return &ERC20AllowanceResponse{
		Allowance: s,
	}
}

type ERC20ApproveRequest struct {
	Contract string `json:"contract" example:"0x8336c196..."`
	Spender  string `json:"spender"  example:"0xAbcD1234..."`
	Value    string `json:"value"    example:"1000000000000000000"`
	Block    string `json:"block"    example:"latest"`

	spender *types.Address
	value   *big.Int
}

func (r *ERC20ApproveRequest) ValidateRequest() error {
	r.Contract = strings.TrimSpace(r.Contract)
	if !util.IsHexAddress(r.Contract) {
		return errors.New("contract: invalid address")
	}
	r.Spender = strings.TrimSpace(r.Spender)
	var err error
	if r.spender, err = types.NewAddressFromHex(r.Spender); err != nil {
		return errors.New("spender: invalid address")
	}
	r.value = new(big.Int)
	if _, ok := r.value.SetString(strings.TrimSpace(r.Value), 10); !ok {
		return errors.New("value: invalid uint256")
	}
	if r.Block == "" {
		r.Block = "latest"
	}
	return nil
}

func (r *ERC20ApproveRequest) ToSpender() *types.Address { return r.spender }
func (r *ERC20ApproveRequest) ToValue() *big.Int         { return r.value }

type ERC20ApproveResponse struct {
	Approved bool `json:"approved"`
}

func NewERC20ApproveResponse(a bool) *ERC20ApproveResponse {
	return &ERC20ApproveResponse{
		Approved: a,
	}
}
