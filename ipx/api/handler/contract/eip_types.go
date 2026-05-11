package contract

import (
	"errors"
	"math/big"
	"strings"

	"github.com/andantan/evmlab/core/types"
	"github.com/andantan/evmlab/internal/util"
)

type EIP712DomainRequest struct {
	Contract string `json:"contract" example:"0x8336c196..."`
	Block    string `json:"block"    example:"latest"`

	c *types.Address
}

func (r *EIP712DomainRequest) ValidateRequest() error {
	r.Contract = strings.TrimSpace(r.Contract)
	if !util.IsHexAddress(r.Contract) {
		return errors.New("contract: invalid address")
	}
	var err error
	r.c, err = types.NewAddressFromHex(r.Contract)
	if err != nil {
		return errors.New("contract: invalid address")
	}

	if r.Block == "" {
		r.Block = "latest"
	}
	return nil
}

func (r *EIP712DomainRequest) ToContract() *types.Address { return r.c }

type EIP712DomainResponse struct {
	Name              string `json:"name"`
	Version           string `json:"version"`
	ChainID           string `json:"chain_id"`
	VerifyingContract string `json:"verifying_contract"`
}

func NewEIP712DomainResponse(d *types.EIP712Domain) *EIP712DomainResponse {
	return &EIP712DomainResponse{
		Name:              d.Name,
		Version:           d.Version,
		ChainID:           d.ChainID.String(),
		VerifyingContract: d.Contract.String(),
	}
}

type EIP2612NoncesRequest struct {
	Contract string `json:"contract" example:"0x8336c196..."`
	Owner    string `json:"owner"    example:"0xAbcD1234..."`
	Block    string `json:"block"    example:"latest"`

	o *types.Address
}

func (r *EIP2612NoncesRequest) ValidateRequest() error {
	r.Contract = strings.TrimSpace(r.Contract)
	if !util.IsHexAddress(r.Contract) {
		return errors.New("contract: invalid address")
	}
	r.Owner = strings.TrimSpace(r.Owner)
	var err error
	if r.o, err = types.NewAddressFromHex(r.Owner); err != nil {
		return errors.New("owner: invalid address")
	}
	if r.Block == "" {
		r.Block = "latest"
	}
	return nil
}

func (r *EIP2612NoncesRequest) ToOwner() *types.Address { return r.o }

type EIP2612NoncesResponse struct {
	Nonce string `json:"nonce"`
}

func NewEIP2612NoncesResponse(n *big.Int) *EIP2612NoncesResponse {
	return &EIP2612NoncesResponse{
		Nonce: n.String(),
	}
}
