package misc

import (
	"errors"
	"strings"

	"github.com/andantan/evmlab/core"
	"github.com/andantan/evmlab/core/types"
	"github.com/andantan/evmlab/internal/util"
)

type EIP55Request struct {
	Address string `json:"address" example:"0xEbD69375..."`

	addr *types.Address
}

func (r *EIP55Request) ValidateRequest() error {
	r.Address = strings.TrimSpace(r.Address)
	if r.Address == "" {
		return errors.New("address is required")
	}

	b, err := util.ParseHex(r.Address)
	if err != nil {
		return errors.New("address: " + err.Error())
	}

	r.addr, err = types.NewAddressFromBytes(b)
	if err != nil {
		return errors.New("address: " + err.Error())
	}

	return nil
}

func (r *EIP55Request) ToAddress() *types.Address {
	return r.addr
}

type EIP55Response struct {
	Address string `json:"address"`
}

func NewEIP55Response(addr *types.Address) *EIP55Response {
	return &EIP55Response{
		Address: addr.Checksum(),
	}
}

type DeriveKeyRequest struct {
	PrivateKey string `json:"private_key" example:"0xea66255f..."`
}

func (r *DeriveKeyRequest) ValidateRequest() error {
	r.PrivateKey = strings.TrimSpace(r.PrivateKey)
	if r.PrivateKey == "" {
		return errors.New("private_key is required")
	}
	return nil
}

type DeriveKeyResponse struct {
	Address    string `json:"address"`
	PublicKey  string `json:"public_key"`
	PrivateKey string `json:"private_key"`
}

func NewDeriveKeyResponse(key *core.EVMSecp256k1Key) *DeriveKeyResponse {
	return &DeriveKeyResponse{
		Address:    key.Address.Checksum(),
		PublicKey:  key.PublicKey.Hex(),
		PrivateKey: key.PrivateKey.Hex(),
	}
}

type UnitConvertRequest struct {
	Amount string `json:"amount" example:"1"`
	From   string `json:"from"   example:"ether"`
	To     string `json:"to"     example:"wei"`
}

func (r *UnitConvertRequest) ValidateRequest() error {
	r.Amount = strings.TrimSpace(r.Amount)
	r.From = strings.ToLower(strings.TrimSpace(r.From))
	r.To = strings.ToLower(strings.TrimSpace(r.To))

	if r.Amount == "" {
		return errors.New("amount is required")
	}
	if r.From == "" {
		return errors.New("from is required")
	}
	if r.To == "" {
		return errors.New("to is required")
	}
	if !util.IsSupportedEthereumUnit(r.From) {
		return errors.New("from: unsupported unit")
	}
	if !util.IsSupportedEthereumUnit(r.To) {
		return errors.New("to: unsupported unit")
	}

	return nil
}

type UnitConvertResponse struct {
	Amount string `json:"amount"`
	Unit   string `json:"unit"`
}

func NewUnitConvertResponse(amount string, unit string) *UnitConvertResponse {
	return &UnitConvertResponse{
		Amount: amount,
		Unit:   unit,
	}
}
