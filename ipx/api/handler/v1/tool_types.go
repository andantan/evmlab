package v1

import (
	"errors"
	"math/big"
	"strings"

	"github.com/andantan/evmlab/core"
	"github.com/andantan/evmlab/core/types"
	"github.com/andantan/evmlab/internal/util"
)

type ChecksumEIP55Request struct {
	Address string `json:"address" example:"0xEbD69375..."`

	addr *types.Address
}

func (r *ChecksumEIP55Request) ValidateRequest() error {
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

func (r *ChecksumEIP55Request) ToAddress() *types.Address {
	return r.addr
}

type ChecksumEIP55Response struct {
	Address string `json:"address"`
}

func NewChecksumEIP55Response(addr *types.Address) *ChecksumEIP55Response {
	return &ChecksumEIP55Response{
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

type UnitConvertDecimalRequest struct {
	Amount string `json:"amount" example:"1"`
	From   string `json:"from"   example:"ether"`
	To     string `json:"to"     example:"wei"`

	amount *big.Int
}

func (r *UnitConvertDecimalRequest) ValidateRequest() error {
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

	n, ok := new(big.Int).SetString(r.Amount, 10)
	if !ok {
		return errors.New("amount: invalid integer")
	}
	r.amount = n

	return nil
}

func (r *UnitConvertDecimalRequest) ToAmount() *big.Int {
	return new(big.Int).Set(r.amount)
}

type UnitConvertHexRequest struct {
	Amount string `json:"amount" example:"0xde0b6b3a7640000"`
	From   string `json:"from"   example:"wei"`
	To     string `json:"to"     example:"ether"`

	amount *big.Int
}

func (r *UnitConvertHexRequest) ValidateRequest() error {
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

	n, err := util.HexToBigInt(r.Amount)
	if err != nil {
		return errors.New("amount: " + err.Error())
	}
	r.amount = n

	return nil
}

func (r *UnitConvertHexRequest) ToAmount() *big.Int {
	return new(big.Int).Set(r.amount)
}

type UnitConvertDecimalResponse struct {
	Amount string `json:"amount"`
	Unit   string `json:"unit"`
}

func NewUnitConvertDecimalResponse(amount string, unit string) *UnitConvertDecimalResponse {
	return &UnitConvertDecimalResponse{
		Amount: amount,
		Unit:   unit,
	}
}

type UnitConvertHexResponse struct {
	Amount string `json:"amount"`
	Unit   string `json:"unit"`
}

func NewUnitConvertHexResponse(amount string, unit string) *UnitConvertHexResponse {
	return &UnitConvertHexResponse{
		Amount: amount,
		Unit:   unit,
	}
}
