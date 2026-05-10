package v1

import (
	"errors"
	"math/big"
	"strings"

	"github.com/andantan/evmlab/core"
	"github.com/andantan/evmlab/core/types"
	"github.com/ethereum/go-ethereum/common"
)

type Keccak256LegacyRequest struct {
	Message string `json:"message" example:"hello world!"`
}

func (r *Keccak256LegacyRequest) ValidateRequest() error {
	r.Message = strings.TrimSpace(r.Message)
	if r.Message == "" {
		return errors.New("message is required")
	}

	return nil
}

type Keccak256LegacyResponse struct {
	Digest string `json:"digest"`
}

func NewKeccak256LegacyResponse(h *types.Hash) *Keccak256LegacyResponse {
	return &Keccak256LegacyResponse{
		Digest: h.String(),
	}
}

type Keccak256EIP191Request struct {
	Message string `json:"message" example:"hello world!"`
}

func (r *Keccak256EIP191Request) ValidateRequest() error {
	r.Message = strings.TrimSpace(r.Message)
	if r.Message == "" {
		return errors.New("message is required")
	}

	return nil
}

type Keccak256EIP191Response struct {
	Digest string `json:"digest"`
}

func NewKeccak256EIP191Response(h *types.Hash) *Keccak256EIP191Response {
	return &Keccak256EIP191Response{
		Digest: h.String(),
	}
}

type Keccak256EIP712Request struct {
	Name      string   `json:"name"`
	Version   string   `json:"version"`
	ChainID   string   `json:"chain_id"`
	Contract  string   `json:"contract"`
	Signature string   `json:"signature"`
	Args      []string `json:"args"`

	fn     *types.Function
	domain *types.EIP712Domain
}

func (r *Keccak256EIP712Request) ValidateRequest() error {
	r.Name = strings.TrimSpace(r.Name)
	if r.Name == "" {
		return errors.New("name is required")
	}
	r.Version = strings.TrimSpace(r.Version)
	if r.Version == "" {
		return errors.New("version is required")
	}
	r.Contract = strings.TrimSpace(r.Contract)
	if !common.IsHexAddress(r.Contract) {
		return errors.New("contract: invalid address")
	}
	r.Signature = strings.TrimSpace(r.Signature)
	if r.Signature == "" {
		return errors.New("signature is required")
	}

	var err error
	r.fn, err = core.ABI.ParseFunctionSignature(r.Signature)
	if err != nil {
		return errors.New("signature: " + err.Error())
	}
	if len(r.fn.Names) == 0 {
		return errors.New("signature must include parameter names for EIP-712")
	}
	if len(r.Args) != len(r.fn.Types) {
		return errors.New("args count does not match signature")
	}

	chainID, ok := new(big.Int).SetString(r.ChainID, 10)
	if !ok {
		return errors.New("chain_id: invalid integer")
	}

	r.domain = &types.EIP712Domain{
		Name:     r.Name,
		Version:  r.Version,
		ChainID:  chainID,
		Contract: types.NewAddress(common.HexToAddress(r.Contract)),
	}

	return nil
}

func (r *Keccak256EIP712Request) ToFn() *types.Function {
	return r.fn
}

func (r *Keccak256EIP712Request) ToEIP712Domain() *types.EIP712Domain {
	return r.domain
}

type Keccak256EIP712Response struct {
	Digest          string `json:"digest"`
	DomainSeparator string `json:"domain_separator"`
	MessageHash     string `json:"message_hash"`
}

func NewKeccak256EIP712Response(r *types.EIP712Result) *Keccak256EIP712Response {
	return &Keccak256EIP712Response{
		Digest:          r.Digest.String(),
		DomainSeparator: r.DomainSeparator.String(),
		MessageHash:     r.MessageHash.String(),
	}
}
