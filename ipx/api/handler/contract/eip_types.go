package contract

import (
	"errors"
	"math/big"
	"strings"

	"github.com/andantan/evmlab/internal/util"
)

type EIP712DomainRequest struct {
	Contract string `json:"contract" example:"0x8336c196..."`
	Block    string `json:"block"    example:"latest"`
}

func (r *EIP712DomainRequest) ValidateRequest() error {
	r.Contract = strings.TrimSpace(r.Contract)
	if !util.IsHexAddress(r.Contract) {
		return errors.New("contract: invalid address")
	}
	if r.Block == "" {
		r.Block = "latest"
	}
	return nil
}

type EIP712DomainResponse struct {
	Name              string `json:"name"`
	Version           string `json:"version"`
	ChainID           string `json:"chain_id"`
	VerifyingContract string `json:"verifying_contract"`
}

func NewEIP712DomainResponse(n, v string, id *big.Int, ca string) *EIP712DomainResponse {
	return &EIP712DomainResponse{
		Name:              n,
		Version:           v,
		ChainID:           id.String(),
		VerifyingContract: ca,
	}
}
