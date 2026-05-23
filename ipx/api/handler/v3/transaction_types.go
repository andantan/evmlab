package v3

import (
	"encoding/hex"
	"errors"
	"math/big"
	"strings"

	"github.com/andantan/evmlab/core/types"
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
	var err error
	r.From = strings.TrimSpace(r.From)
	if r.From == "" {
		return errors.New("from is required")
	}

	r.f, err = types.NewAddressFromHex(r.From)
	if err != nil {
		return errors.New("from: invalid address")
	}

	r.To = strings.TrimSpace(r.To)
	if r.To == "" {
		return errors.New("to is required")
	}

	r.t, err = types.NewAddressFromHex(r.To)
	if err != nil {
		return errors.New("to: invalid address")
	}

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
	var err error
	r.From = strings.TrimSpace(r.From)
	if r.From == "" {
		return errors.New("from is required")
	}

	r.f, err = types.NewAddressFromHex(r.From)
	if err != nil {
		return errors.New("from: invalid address")
	}

	r.To = strings.TrimSpace(r.To)
	if r.To == "" {
		return errors.New("to is required")
	}

	r.t, err = types.NewAddressFromHex(r.To)
	if err != nil {
		return errors.New("to: invalid address")
	}

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

type BuildERC20LegacyTransactionRequest struct {
	From     string `json:"from"     example:"0xDa70aA79f1a329719b9cf9d334b0a82b1d5269f3"`
	To       string `json:"to"       example:"0xDa70aA79f1a329719b9cf9d334b0a82b1d5269f3"`
	Contract string `json:"contract" example:"0x5FbDB2315678afecb367f032d93F642f64180aa3"`
	Amount   string `json:"amount"   example:"1000000000000000000"`

	f *types.Address
	t *types.Address
	c *types.Address
	a *big.Int
}

func (r *BuildERC20LegacyTransactionRequest) ValidateRequest() error {
	var err error
	r.From = strings.TrimSpace(r.From)
	r.f, err = types.NewAddressFromHex(r.From)
	if err != nil {
		return errors.New("from: invalid address")
	}

	r.To = strings.TrimSpace(r.To)
	r.t, err = types.NewAddressFromHex(r.To)
	if err != nil {
		return errors.New("to: invalid address")
	}

	r.Contract = strings.TrimSpace(r.Contract)
	r.c, err = types.NewAddressFromHex(r.Contract)
	if err != nil {
		return errors.New("contract: invalid address")
	}

	r.Amount = strings.TrimSpace(r.Amount)
	a, ok := new(big.Int).SetString(r.Amount, 10)
	if !ok {
		return errors.New("amount: invalid integer")
	}
	r.a = a

	return nil
}

func (r *BuildERC20LegacyTransactionRequest) FromAddr() *types.Address     { return r.f }
func (r *BuildERC20LegacyTransactionRequest) ToAddr() *types.Address       { return r.t }
func (r *BuildERC20LegacyTransactionRequest) ContractAddr() *types.Address { return r.c }
func (r *BuildERC20LegacyTransactionRequest) ToAmount() *big.Int           { return r.a }

type BuildERC20LegacyTransactionResponse struct {
	UnsignedRLP string `json:"unsigned_rlp"`
	SignedRLP   string `json:"signed_rlp"`
	TxHash      string `json:"tx_hash"`
	R           string `json:"r"`
	S           string `json:"s"`
	V           string `json:"v"`
}

func NewBuildERC20LegacyTransactionResponse(unsigned, signed []byte, txHash *types.Hash, sig *types.Signature) *BuildERC20LegacyTransactionResponse {
	return &BuildERC20LegacyTransactionResponse{
		UnsignedRLP: "0x" + hex.EncodeToString(unsigned),
		SignedRLP:   "0x" + hex.EncodeToString(signed),
		TxHash:      txHash.String(),
		R:           "0x" + sig.R().Text(16),
		S:           "0x" + sig.S().Text(16),
		V:           "0x" + sig.V().Text(16),
	}
}

type BuildERC20EIP1559TransactionRequest struct {
	From     string `json:"from"     example:"0xDa70aA79f1a329719b9cf9d334b0a82b1d5269f3"`
	To       string `json:"to"       example:"0xDa70aA79f1a329719b9cf9d334b0a82b1d5269f3"`
	Contract string `json:"contract" example:"0x5FbDB2315678afecb367f032d93F642f64180aa3"`
	Amount   string `json:"amount"   example:"1000000000000000000"`

	f *types.Address
	t *types.Address
	c *types.Address
	a *big.Int
}

func (r *BuildERC20EIP1559TransactionRequest) ValidateRequest() error {
	var err error
	r.From = strings.TrimSpace(r.From)
	r.f, err = types.NewAddressFromHex(r.From)
	if err != nil {
		return errors.New("from: invalid address")
	}

	r.To = strings.TrimSpace(r.To)
	r.t, err = types.NewAddressFromHex(r.To)
	if err != nil {
		return errors.New("to: invalid address")
	}

	r.Contract = strings.TrimSpace(r.Contract)
	r.c, err = types.NewAddressFromHex(r.Contract)
	if err != nil {
		return errors.New("contract: invalid address")
	}

	r.Amount = strings.TrimSpace(r.Amount)
	a, ok := new(big.Int).SetString(r.Amount, 10)
	if !ok {
		return errors.New("amount: invalid integer")
	}
	r.a = a

	return nil
}

func (r *BuildERC20EIP1559TransactionRequest) FromAddr() *types.Address     { return r.f }
func (r *BuildERC20EIP1559TransactionRequest) ToAddr() *types.Address       { return r.t }
func (r *BuildERC20EIP1559TransactionRequest) ContractAddr() *types.Address { return r.c }
func (r *BuildERC20EIP1559TransactionRequest) ToAmount() *big.Int           { return r.a }

type BuildERC20EIP1559TransactionResponse struct {
	UnsignedRLP string `json:"unsigned_rlp"`
	SignedRLP   string `json:"signed_rlp"`
	TxHash      string `json:"tx_hash"`
	R           string `json:"r"`
	S           string `json:"s"`
	V           string `json:"v"`
}

func NewBuildERC20EIP1559TransactionResponse(unsigned, signed []byte, txHash *types.Hash, sig *types.Signature) *BuildERC20EIP1559TransactionResponse {
	return &BuildERC20EIP1559TransactionResponse{
		UnsignedRLP: "0x" + hex.EncodeToString(unsigned),
		SignedRLP:   "0x" + hex.EncodeToString(signed),
		TxHash:      txHash.String(),
		R:           "0x" + sig.R().Text(16),
		S:           "0x" + sig.S().Text(16),
		V:           "0x" + sig.V().Text(16),
	}
}

type BuildContractCallLegacyTransactionRequest struct {
	From  string `json:"from"  example:"0xEbD69375d51a8472DF22A3C18405b5A2586c2Aa2"`
	To    string `json:"to"    example:"0xF0a619CDA27104b969086d15Ad7fcDaa4a251Eb2"`
	Data  string `json:"data"  example:"0xe30e3834000000000000000000000000..."`
	Value string `json:"value" example:"0"`

	f *types.Address
	t *types.Address
	d []byte
	v *big.Int
}

func (r *BuildContractCallLegacyTransactionRequest) ValidateRequest() error {
	var err error

	r.From = strings.TrimSpace(r.From)
	r.f, err = types.NewAddressFromHex(r.From)
	if err != nil {
		return errors.New("from: invalid address")
	}

	r.To = strings.TrimSpace(r.To)
	r.t, err = types.NewAddressFromHex(r.To)
	if err != nil {
		return errors.New("to: invalid address")
	}

	r.Data = strings.TrimSpace(r.Data)
	if r.Data == "" {
		return errors.New("data is required")
	}
	r.d, err = hex.DecodeString(strings.TrimPrefix(r.Data, "0x"))
	if err != nil {
		return errors.New("data: invalid hex")
	}

	r.Value = strings.TrimSpace(r.Value)
	if r.Value == "" {
		r.Value = "0"
	}
	v, ok := new(big.Int).SetString(r.Value, 10)
	if !ok {
		return errors.New("value: must be a decimal wei amount")
	}
	r.v = v

	return nil
}

func (r *BuildContractCallLegacyTransactionRequest) FromAddr() *types.Address { return r.f }
func (r *BuildContractCallLegacyTransactionRequest) ToAddr() *types.Address   { return r.t }
func (r *BuildContractCallLegacyTransactionRequest) Calldata() []byte         { return r.d }
func (r *BuildContractCallLegacyTransactionRequest) Amount() *big.Int         { return r.v }

type BuildContractCallLegacyTransactionResponse struct {
	UnsignedRLP string `json:"unsigned_rlp"`
	SignedRLP   string `json:"signed_rlp"`
	TxHash      string `json:"tx_hash"`
	R           string `json:"r"`
	S           string `json:"s"`
	V           string `json:"v"`
}

func NewBuildContractCallLegacyTransactionResponse(unsigned, signed []byte, txHash *types.Hash, sig *types.Signature) *BuildContractCallLegacyTransactionResponse {
	return &BuildContractCallLegacyTransactionResponse{
		UnsignedRLP: "0x" + hex.EncodeToString(unsigned),
		SignedRLP:   "0x" + hex.EncodeToString(signed),
		TxHash:      txHash.String(),
		R:           "0x" + sig.R().Text(16),
		S:           "0x" + sig.S().Text(16),
		V:           "0x" + sig.V().Text(16),
	}
}

type BuildContractCallEIP1559TransactionRequest struct {
	From  string `json:"from"  example:"0xEbD69375d51a8472DF22A3C18405b5A2586c2Aa2"`
	To    string `json:"to"    example:"0xF0a619CDA27104b969086d15Ad7fcDaa4a251Eb2"`
	Data  string `json:"data"  example:"0xe30e3834000000000000000000000000..."`
	Value string `json:"value" example:"0"`

	f *types.Address
	t *types.Address
	d []byte
	v *big.Int
}

func (r *BuildContractCallEIP1559TransactionRequest) ValidateRequest() error {
	var err error

	r.From = strings.TrimSpace(r.From)
	r.f, err = types.NewAddressFromHex(r.From)
	if err != nil {
		return errors.New("from: invalid address")
	}

	r.To = strings.TrimSpace(r.To)
	r.t, err = types.NewAddressFromHex(r.To)
	if err != nil {
		return errors.New("to: invalid address")
	}

	r.Data = strings.TrimSpace(r.Data)
	if r.Data == "" {
		return errors.New("data is required")
	}
	r.d, err = hex.DecodeString(strings.TrimPrefix(r.Data, "0x"))
	if err != nil {
		return errors.New("data: invalid hex")
	}

	r.Value = strings.TrimSpace(r.Value)
	if r.Value == "" {
		r.Value = "0"
	}
	v, ok := new(big.Int).SetString(r.Value, 10)
	if !ok {
		return errors.New("value: must be a decimal wei amount")
	}
	r.v = v

	return nil
}

func (r *BuildContractCallEIP1559TransactionRequest) FromAddr() *types.Address { return r.f }
func (r *BuildContractCallEIP1559TransactionRequest) ToAddr() *types.Address   { return r.t }
func (r *BuildContractCallEIP1559TransactionRequest) Calldata() []byte         { return r.d }
func (r *BuildContractCallEIP1559TransactionRequest) Amount() *big.Int         { return r.v }

type BuildContractCallEIP1559TransactionResponse struct {
	UnsignedRLP string `json:"unsigned_rlp"`
	SignedRLP   string `json:"signed_rlp"`
	TxHash      string `json:"tx_hash"`
	R           string `json:"r"`
	S           string `json:"s"`
	V           string `json:"v"`
}

func NewBuildContractCallEIP1559TransactionResponse(unsigned, signed []byte, txHash *types.Hash, sig *types.Signature) *BuildContractCallEIP1559TransactionResponse {
	return &BuildContractCallEIP1559TransactionResponse{
		UnsignedRLP: "0x" + hex.EncodeToString(unsigned),
		SignedRLP:   "0x" + hex.EncodeToString(signed),
		TxHash:      txHash.String(),
		R:           "0x" + sig.R().Text(16),
		S:           "0x" + sig.S().Text(16),
		V:           "0x" + sig.V().Text(16),
	}
}
