package v3

import (
	"encoding/hex"
	"errors"
	"math/big"
	"strings"

	"github.com/andantan/evmlab/core/types"
	"github.com/ethereum/go-ethereum/common"
)

type BuildNativeLegacyTransactionRequest struct {
	From  string `json:"from"`
	To    string `json:"to"`
	Value string `json:"value"` // wei as decimal string

	f *types.Address
	t *types.Address
	v *big.Int
}

func (r *BuildNativeLegacyTransactionRequest) ValidateRequest() error {
	r.From = strings.TrimSpace(r.From)
	if r.From == "" {
		return errors.New("from is required")
	}
	if !common.IsHexAddress(r.From) {
		return errors.New("from: invalid address")
	}
	r.f = types.NewAddress(common.HexToAddress(r.From))

	r.To = strings.TrimSpace(r.To)
	if r.To == "" {
		return errors.New("to is required")
	}
	if !common.IsHexAddress(r.To) {
		return errors.New("to: invalid address")
	}
	r.t = types.NewAddress(common.HexToAddress(r.To))

	r.Value = strings.TrimSpace(r.Value)
	if r.Value == "" {
		return errors.New("value is required")
	}
	v, ok := new(big.Int).SetString(r.Value, 10)
	if !ok {
		return errors.New("value: must be a decimal wei amount")
	}
	r.v = v

	return nil
}

func (r *BuildNativeLegacyTransactionRequest) FromAddr() *types.Address { return r.f }
func (r *BuildNativeLegacyTransactionRequest) ToAddr() *types.Address   { return r.t }
func (r *BuildNativeLegacyTransactionRequest) Amount() *big.Int         { return r.v }

type BuildNativeLegacyTransactionResponse struct {
	UnsignedRLP string `json:"unsigned_rlp"`
	SignedRLP   string `json:"signed_rlp"`
	TxHash      string `json:"tx_hash"`
	R           string `json:"r"`
	S           string `json:"s"`
	V           string `json:"v"`
}

func NewBuildNativeLegacyTransactionResponse(unsigned, signed []byte, txHash *types.Hash, sig *types.Signature) *BuildNativeLegacyTransactionResponse {
	return &BuildNativeLegacyTransactionResponse{
		UnsignedRLP: "0x" + hex.EncodeToString(unsigned),
		SignedRLP:   "0x" + hex.EncodeToString(signed),
		TxHash:      txHash.String(),
		R:           "0x" + sig.R().Text(16),
		S:           "0x" + sig.S().Text(16),
		V:           "0x" + sig.V().Text(16),
	}
}

type BuildNativeEIP1559TransactionRequest struct {
	From  string `json:"from"`
	To    string `json:"to"`
	Value string `json:"value"` // wei as decimal string

	f *types.Address
	t *types.Address
	v *big.Int
}

func (r *BuildNativeEIP1559TransactionRequest) ValidateRequest() error {
	r.From = strings.TrimSpace(r.From)
	if r.From == "" {
		return errors.New("from is required")
	}
	if !common.IsHexAddress(r.From) {
		return errors.New("from: invalid address")
	}
	r.f = types.NewAddress(common.HexToAddress(r.From))

	r.To = strings.TrimSpace(r.To)
	if r.To == "" {
		return errors.New("to is required")
	}
	if !common.IsHexAddress(r.To) {
		return errors.New("to: invalid address")
	}
	r.t = types.NewAddress(common.HexToAddress(r.To))

	r.Value = strings.TrimSpace(r.Value)
	if r.Value == "" {
		return errors.New("value is required")
	}
	v, ok := new(big.Int).SetString(r.Value, 10)
	if !ok {
		return errors.New("value: must be a decimal wei amount")
	}
	r.v = v

	return nil
}

func (r *BuildNativeEIP1559TransactionRequest) FromAddr() *types.Address { return r.f }
func (r *BuildNativeEIP1559TransactionRequest) ToAddr() *types.Address   { return r.t }
func (r *BuildNativeEIP1559TransactionRequest) Amount() *big.Int         { return r.v }

type BuildNativeEIP1559TransactionResponse struct {
	UnsignedRLP string `json:"unsigned_rlp"`
	SignedRLP   string `json:"signed_rlp"`
	TxHash      string `json:"tx_hash"`
	R           string `json:"r"`
	S           string `json:"s"`
	V           string `json:"v"`
}

func NewBuildNativeEIP1559TransactionResponse(unsigned, signed []byte, txHash *types.Hash, sig *types.Signature) *BuildNativeEIP1559TransactionResponse {
	return &BuildNativeEIP1559TransactionResponse{
		UnsignedRLP: "0x" + hex.EncodeToString(unsigned),
		SignedRLP:   "0x" + hex.EncodeToString(signed),
		TxHash:      txHash.String(),
		R:           "0x" + sig.R().Text(16),
		S:           "0x" + sig.S().Text(16),
		V:           "0x" + sig.V().Text(16),
	}
}
