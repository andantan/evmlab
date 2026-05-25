package sanctum

import (
	"encoding/hex"
	"errors"
	"math/big"
	"strings"

	"github.com/andantan/evmlab/core/types"
	"github.com/andantan/evmlab/internal/util"
)

type SanctumRegisterLegacyRequest struct {
	From    string `json:"from"     example:"0xAbcD1234..."`
	Sanctum string `json:"sanctum"  example:"0x1234Abcd..."`

	f *types.Address
	s *types.Address
}

func (r *SanctumRegisterLegacyRequest) ValidateRequest() error {
	r.From = strings.TrimSpace(r.From)
	if !util.IsHexAddress(r.From) {
		return errors.New("from: invalid address")
	}
	var err error
	if r.f, err = types.NewAddressFromHex(r.From); err != nil {
		return errors.New("from: invalid address")
	}

	r.Sanctum = strings.TrimSpace(r.Sanctum)
	if !util.IsHexAddress(r.Sanctum) {
		return errors.New("contract: invalid address")
	}
	if r.s, err = types.NewAddressFromHex(r.Sanctum); err != nil {
		return errors.New("contract: invalid address")
	}

	return nil
}

func (r *SanctumRegisterLegacyRequest) FromAddr() *types.Address    { return r.f }
func (r *SanctumRegisterLegacyRequest) SanctumAddr() *types.Address { return r.s }

type SanctumRegisterLegacyResponse struct {
	UnsignedRLP string `json:"unsigned_rlp"`
	SignedRLP   string `json:"signed_rlp"`
	TxHash      string `json:"tx_hash"`
	R           string `json:"r"`
	S           string `json:"s"`
	V           string `json:"v"`
}

func NewSanctumRegisterLegacyResponse(unsigned, signed []byte, txHash *types.Hash, sig *types.Signature) *SanctumRegisterLegacyResponse {
	return &SanctumRegisterLegacyResponse{
		UnsignedRLP: "0x" + hex.EncodeToString(unsigned),
		SignedRLP:   "0x" + hex.EncodeToString(signed),
		TxHash:      txHash.String(),
		R:           "0x" + sig.R().Text(16),
		S:           "0x" + sig.S().Text(16),
		V:           "0x" + sig.V().Text(16),
	}
}

type SanctumRegisterEIP1559Request struct {
	From    string `json:"from"     example:"0xAbcD1234..."`
	Sanctum string `json:"sanctum"  example:"0x1234Abcd..."`

	f *types.Address
	s *types.Address
}

func (r *SanctumRegisterEIP1559Request) ValidateRequest() error {
	r.From = strings.TrimSpace(r.From)
	if !util.IsHexAddress(r.From) {
		return errors.New("from: invalid address")
	}
	var err error
	if r.f, err = types.NewAddressFromHex(r.From); err != nil {
		return errors.New("from: invalid address")
	}

	r.Sanctum = strings.TrimSpace(r.Sanctum)
	if !util.IsHexAddress(r.Sanctum) {
		return errors.New("contract: invalid address")
	}
	if r.s, err = types.NewAddressFromHex(r.Sanctum); err != nil {
		return errors.New("contract: invalid address")
	}

	return nil
}

func (r *SanctumRegisterEIP1559Request) FromAddr() *types.Address    { return r.f }
func (r *SanctumRegisterEIP1559Request) SanctumAddr() *types.Address { return r.s }

// SanctumTxResponse is the shared response for all Sanctum transaction endpoints.
type SanctumRegisterEIP1559Response struct {
	UnsignedRLP string `json:"unsigned_rlp"`
	SignedRLP   string `json:"signed_rlp"`
	TxHash      string `json:"tx_hash"`
	R           string `json:"r"`
	S           string `json:"s"`
	V           string `json:"v"`
}

func NewSanctumRegisterEIP1559Response(unsigned, signed []byte, txHash *types.Hash, sig *types.Signature) *SanctumRegisterEIP1559Response {
	return &SanctumRegisterEIP1559Response{
		UnsignedRLP: "0x" + hex.EncodeToString(unsigned),
		SignedRLP:   "0x" + hex.EncodeToString(signed),
		TxHash:      txHash.String(),
		R:           "0x" + sig.R().Text(16),
		S:           "0x" + sig.S().Text(16),
		V:           "0x" + sig.V().Text(16),
	}
}

type SanctumRegisterForLegacyRequest struct {
	From    string `json:"from"    example:"0xAbcD1234..."`
	Sanctum string `json:"sanctum" example:"0x1234Abcd..."`
	Target  string `json:"target"  example:"0xDeadBeef..."`

	f *types.Address
	s *types.Address
	t *types.Address
}

func (r *SanctumRegisterForLegacyRequest) ValidateRequest() error {
	r.From = strings.TrimSpace(r.From)
	if !util.IsHexAddress(r.From) {
		return errors.New("from: invalid address")
	}
	var err error
	if r.f, err = types.NewAddressFromHex(r.From); err != nil {
		return errors.New("from: invalid address")
	}

	r.Sanctum = strings.TrimSpace(r.Sanctum)
	if !util.IsHexAddress(r.Sanctum) {
		return errors.New("sanctum: invalid address")
	}
	if r.s, err = types.NewAddressFromHex(r.Sanctum); err != nil {
		return errors.New("sanctum: invalid address")
	}

	r.Target = strings.TrimSpace(r.Target)
	if !util.IsHexAddress(r.Target) {
		return errors.New("target: invalid address")
	}
	if r.t, err = types.NewAddressFromHex(r.Target); err != nil {
		return errors.New("target: invalid address")
	}

	return nil
}

func (r *SanctumRegisterForLegacyRequest) FromAddr() *types.Address    { return r.f }
func (r *SanctumRegisterForLegacyRequest) SanctumAddr() *types.Address { return r.s }
func (r *SanctumRegisterForLegacyRequest) TargetAddr() *types.Address  { return r.t }

type SanctumRegisterForLegacyResponse struct {
	UnsignedRLP string `json:"unsigned_rlp"`
	SignedRLP   string `json:"signed_rlp"`
	TxHash      string `json:"tx_hash"`
	R           string `json:"r"`
	S           string `json:"s"`
	V           string `json:"v"`
}

func NewSanctumRegisterForLegacyResponse(unsigned, signed []byte, txHash *types.Hash, sig *types.Signature) *SanctumRegisterForLegacyResponse {
	return &SanctumRegisterForLegacyResponse{
		UnsignedRLP: "0x" + hex.EncodeToString(unsigned),
		SignedRLP:   "0x" + hex.EncodeToString(signed),
		TxHash:      txHash.String(),
		R:           "0x" + sig.R().Text(16),
		S:           "0x" + sig.S().Text(16),
		V:           "0x" + sig.V().Text(16),
	}
}

type SanctumRegisterForEIP1559Request struct {
	From    string `json:"from"    example:"0xAbcD1234..."`
	Sanctum string `json:"sanctum" example:"0x1234Abcd..."`
	Target  string `json:"target"  example:"0xDeadBeef..."`

	f *types.Address
	s *types.Address
	t *types.Address
}

func (r *SanctumRegisterForEIP1559Request) ValidateRequest() error {
	r.From = strings.TrimSpace(r.From)
	if !util.IsHexAddress(r.From) {
		return errors.New("from: invalid address")
	}
	var err error
	if r.f, err = types.NewAddressFromHex(r.From); err != nil {
		return errors.New("from: invalid address")
	}

	r.Sanctum = strings.TrimSpace(r.Sanctum)
	if !util.IsHexAddress(r.Sanctum) {
		return errors.New("sanctum: invalid address")
	}
	if r.s, err = types.NewAddressFromHex(r.Sanctum); err != nil {
		return errors.New("sanctum: invalid address")
	}

	r.Target = strings.TrimSpace(r.Target)
	if !util.IsHexAddress(r.Target) {
		return errors.New("target: invalid address")
	}
	if r.t, err = types.NewAddressFromHex(r.Target); err != nil {
		return errors.New("target: invalid address")
	}

	return nil
}

func (r *SanctumRegisterForEIP1559Request) FromAddr() *types.Address    { return r.f }
func (r *SanctumRegisterForEIP1559Request) SanctumAddr() *types.Address { return r.s }
func (r *SanctumRegisterForEIP1559Request) TargetAddr() *types.Address  { return r.t }

type SanctumRegisterForEIP1559Response struct {
	UnsignedRLP string `json:"unsigned_rlp"`
	SignedRLP   string `json:"signed_rlp"`
	TxHash      string `json:"tx_hash"`
	R           string `json:"r"`
	S           string `json:"s"`
	V           string `json:"v"`
}

func NewSanctumRegisterForEIP1559Response(unsigned, signed []byte, txHash *types.Hash, sig *types.Signature) *SanctumRegisterForEIP1559Response {
	return &SanctumRegisterForEIP1559Response{
		UnsignedRLP: "0x" + hex.EncodeToString(unsigned),
		SignedRLP:   "0x" + hex.EncodeToString(signed),
		TxHash:      txHash.String(),
		R:           "0x" + sig.R().Text(16),
		S:           "0x" + sig.S().Text(16),
		V:           "0x" + sig.V().Text(16),
	}
}

type SanctumApproveRegisterLegacyRequest struct {
	From    string `json:"from"    example:"0xAbcD1234..."`
	Sanctum string `json:"sanctum" example:"0x1234Abcd..."`
	Target  string `json:"target"  example:"0xDeadBeef..."`

	f *types.Address
	s *types.Address
	t *types.Address
}

func (r *SanctumApproveRegisterLegacyRequest) ValidateRequest() error {
	r.From = strings.TrimSpace(r.From)
	if !util.IsHexAddress(r.From) {
		return errors.New("from: invalid address")
	}
	var err error
	if r.f, err = types.NewAddressFromHex(r.From); err != nil {
		return errors.New("from: invalid address")
	}

	r.Sanctum = strings.TrimSpace(r.Sanctum)
	if !util.IsHexAddress(r.Sanctum) {
		return errors.New("sanctum: invalid address")
	}
	if r.s, err = types.NewAddressFromHex(r.Sanctum); err != nil {
		return errors.New("sanctum: invalid address")
	}

	r.Target = strings.TrimSpace(r.Target)
	if !util.IsHexAddress(r.Target) {
		return errors.New("target: invalid address")
	}
	if r.t, err = types.NewAddressFromHex(r.Target); err != nil {
		return errors.New("target: invalid address")
	}

	return nil
}

func (r *SanctumApproveRegisterLegacyRequest) FromAddr() *types.Address    { return r.f }
func (r *SanctumApproveRegisterLegacyRequest) SanctumAddr() *types.Address { return r.s }
func (r *SanctumApproveRegisterLegacyRequest) TargetAddr() *types.Address  { return r.t }

type SanctumApproveRegisterLegacyResponse struct {
	UnsignedRLP string `json:"unsigned_rlp"`
	SignedRLP   string `json:"signed_rlp"`
	TxHash      string `json:"tx_hash"`
	R           string `json:"r"`
	S           string `json:"s"`
	V           string `json:"v"`
}

func NewSanctumApproveRegisterLegacyResponse(unsigned, signed []byte, txHash *types.Hash, sig *types.Signature) *SanctumApproveRegisterLegacyResponse {
	return &SanctumApproveRegisterLegacyResponse{
		UnsignedRLP: "0x" + hex.EncodeToString(unsigned),
		SignedRLP:   "0x" + hex.EncodeToString(signed),
		TxHash:      txHash.String(),
		R:           "0x" + sig.R().Text(16),
		S:           "0x" + sig.S().Text(16),
		V:           "0x" + sig.V().Text(16),
	}
}

type SanctumApproveRegisterEIP1559Request struct {
	From    string `json:"from"    example:"0xAbcD1234..."`
	Sanctum string `json:"sanctum" example:"0x1234Abcd..."`
	Target  string `json:"target"  example:"0xDeadBeef..."`

	f *types.Address
	s *types.Address
	t *types.Address
}

func (r *SanctumApproveRegisterEIP1559Request) ValidateRequest() error {
	r.From = strings.TrimSpace(r.From)
	if !util.IsHexAddress(r.From) {
		return errors.New("from: invalid address")
	}
	var err error
	if r.f, err = types.NewAddressFromHex(r.From); err != nil {
		return errors.New("from: invalid address")
	}

	r.Sanctum = strings.TrimSpace(r.Sanctum)
	if !util.IsHexAddress(r.Sanctum) {
		return errors.New("sanctum: invalid address")
	}
	if r.s, err = types.NewAddressFromHex(r.Sanctum); err != nil {
		return errors.New("sanctum: invalid address")
	}

	r.Target = strings.TrimSpace(r.Target)
	if !util.IsHexAddress(r.Target) {
		return errors.New("target: invalid address")
	}
	if r.t, err = types.NewAddressFromHex(r.Target); err != nil {
		return errors.New("target: invalid address")
	}

	return nil
}

func (r *SanctumApproveRegisterEIP1559Request) FromAddr() *types.Address    { return r.f }
func (r *SanctumApproveRegisterEIP1559Request) SanctumAddr() *types.Address { return r.s }
func (r *SanctumApproveRegisterEIP1559Request) TargetAddr() *types.Address  { return r.t }

type SanctumApproveRegisterEIP1559Response struct {
	UnsignedRLP string `json:"unsigned_rlp"`
	SignedRLP   string `json:"signed_rlp"`
	TxHash      string `json:"tx_hash"`
	R           string `json:"r"`
	S           string `json:"s"`
	V           string `json:"v"`
}

func NewSanctumApproveRegisterEIP1559Response(unsigned, signed []byte, txHash *types.Hash, sig *types.Signature) *SanctumApproveRegisterEIP1559Response {
	return &SanctumApproveRegisterEIP1559Response{
		UnsignedRLP: "0x" + hex.EncodeToString(unsigned),
		SignedRLP:   "0x" + hex.EncodeToString(signed),
		TxHash:      txHash.String(),
		R:           "0x" + sig.R().Text(16),
		S:           "0x" + sig.S().Text(16),
		V:           "0x" + sig.V().Text(16),
	}
}

type SanctumDeregisterLegacyRequest struct {
	From    string `json:"from"    example:"0xAbcD1234..."`
	Sanctum string `json:"sanctum" example:"0x1234Abcd..."`

	f *types.Address
	s *types.Address
}

func (r *SanctumDeregisterLegacyRequest) ValidateRequest() error {
	r.From = strings.TrimSpace(r.From)
	if !util.IsHexAddress(r.From) {
		return errors.New("from: invalid address")
	}
	var err error
	if r.f, err = types.NewAddressFromHex(r.From); err != nil {
		return errors.New("from: invalid address")
	}

	r.Sanctum = strings.TrimSpace(r.Sanctum)
	if !util.IsHexAddress(r.Sanctum) {
		return errors.New("sanctum: invalid address")
	}
	if r.s, err = types.NewAddressFromHex(r.Sanctum); err != nil {
		return errors.New("sanctum: invalid address")
	}

	return nil
}

func (r *SanctumDeregisterLegacyRequest) FromAddr() *types.Address    { return r.f }
func (r *SanctumDeregisterLegacyRequest) SanctumAddr() *types.Address { return r.s }

type SanctumDeregisterLegacyResponse struct {
	UnsignedRLP string `json:"unsigned_rlp"`
	SignedRLP   string `json:"signed_rlp"`
	TxHash      string `json:"tx_hash"`
	R           string `json:"r"`
	S           string `json:"s"`
	V           string `json:"v"`
}

func NewSanctumDeregisterLegacyResponse(unsigned, signed []byte, txHash *types.Hash, sig *types.Signature) *SanctumDeregisterLegacyResponse {
	return &SanctumDeregisterLegacyResponse{
		UnsignedRLP: "0x" + hex.EncodeToString(unsigned),
		SignedRLP:   "0x" + hex.EncodeToString(signed),
		TxHash:      txHash.String(),
		R:           "0x" + sig.R().Text(16),
		S:           "0x" + sig.S().Text(16),
		V:           "0x" + sig.V().Text(16),
	}
}

type SanctumDeregisterEIP1559Request struct {
	From    string `json:"from"    example:"0xAbcD1234..."`
	Sanctum string `json:"sanctum" example:"0x1234Abcd..."`

	f *types.Address
	s *types.Address
}

func (r *SanctumDeregisterEIP1559Request) ValidateRequest() error {
	r.From = strings.TrimSpace(r.From)
	if !util.IsHexAddress(r.From) {
		return errors.New("from: invalid address")
	}
	var err error
	if r.f, err = types.NewAddressFromHex(r.From); err != nil {
		return errors.New("from: invalid address")
	}

	r.Sanctum = strings.TrimSpace(r.Sanctum)
	if !util.IsHexAddress(r.Sanctum) {
		return errors.New("sanctum: invalid address")
	}
	if r.s, err = types.NewAddressFromHex(r.Sanctum); err != nil {
		return errors.New("sanctum: invalid address")
	}

	return nil
}

func (r *SanctumDeregisterEIP1559Request) FromAddr() *types.Address    { return r.f }
func (r *SanctumDeregisterEIP1559Request) SanctumAddr() *types.Address { return r.s }

type SanctumDeregisterEIP1559Response struct {
	UnsignedRLP string `json:"unsigned_rlp"`
	SignedRLP   string `json:"signed_rlp"`
	TxHash      string `json:"tx_hash"`
	R           string `json:"r"`
	S           string `json:"s"`
	V           string `json:"v"`
}

func NewSanctumDeregisterEIP1559Response(unsigned, signed []byte, txHash *types.Hash, sig *types.Signature) *SanctumDeregisterEIP1559Response {
	return &SanctumDeregisterEIP1559Response{
		UnsignedRLP: "0x" + hex.EncodeToString(unsigned),
		SignedRLP:   "0x" + hex.EncodeToString(signed),
		TxHash:      txHash.String(),
		R:           "0x" + sig.R().Text(16),
		S:           "0x" + sig.S().Text(16),
		V:           "0x" + sig.V().Text(16),
	}
}

type SanctumDeregisterForLegacyRequest struct {
	From    string `json:"from"    example:"0xAbcD1234..."`
	Sanctum string `json:"sanctum" example:"0x1234Abcd..."`
	Target  string `json:"target"  example:"0xDeadBeef..."`

	f *types.Address
	s *types.Address
	t *types.Address
}

func (r *SanctumDeregisterForLegacyRequest) ValidateRequest() error {
	r.From = strings.TrimSpace(r.From)
	if !util.IsHexAddress(r.From) {
		return errors.New("from: invalid address")
	}
	var err error
	if r.f, err = types.NewAddressFromHex(r.From); err != nil {
		return errors.New("from: invalid address")
	}

	r.Sanctum = strings.TrimSpace(r.Sanctum)
	if !util.IsHexAddress(r.Sanctum) {
		return errors.New("sanctum: invalid address")
	}
	if r.s, err = types.NewAddressFromHex(r.Sanctum); err != nil {
		return errors.New("sanctum: invalid address")
	}

	r.Target = strings.TrimSpace(r.Target)
	if !util.IsHexAddress(r.Target) {
		return errors.New("target: invalid address")
	}
	if r.t, err = types.NewAddressFromHex(r.Target); err != nil {
		return errors.New("target: invalid address")
	}

	return nil
}

func (r *SanctumDeregisterForLegacyRequest) FromAddr() *types.Address    { return r.f }
func (r *SanctumDeregisterForLegacyRequest) SanctumAddr() *types.Address { return r.s }
func (r *SanctumDeregisterForLegacyRequest) TargetAddr() *types.Address  { return r.t }

type SanctumDeregisterForLegacyResponse struct {
	UnsignedRLP string `json:"unsigned_rlp"`
	SignedRLP   string `json:"signed_rlp"`
	TxHash      string `json:"tx_hash"`
	R           string `json:"r"`
	S           string `json:"s"`
	V           string `json:"v"`
}

func NewSanctumDeregisterForLegacyResponse(unsigned, signed []byte, txHash *types.Hash, sig *types.Signature) *SanctumDeregisterForLegacyResponse {
	return &SanctumDeregisterForLegacyResponse{
		UnsignedRLP: "0x" + hex.EncodeToString(unsigned),
		SignedRLP:   "0x" + hex.EncodeToString(signed),
		TxHash:      txHash.String(),
		R:           "0x" + sig.R().Text(16),
		S:           "0x" + sig.S().Text(16),
		V:           "0x" + sig.V().Text(16),
	}
}

type SanctumDeregisterForEIP1559Request struct {
	From    string `json:"from"    example:"0xAbcD1234..."`
	Sanctum string `json:"sanctum" example:"0x1234Abcd..."`
	Target  string `json:"target"  example:"0xDeadBeef..."`

	f *types.Address
	s *types.Address
	t *types.Address
}

func (r *SanctumDeregisterForEIP1559Request) ValidateRequest() error {
	r.From = strings.TrimSpace(r.From)
	if !util.IsHexAddress(r.From) {
		return errors.New("from: invalid address")
	}
	var err error
	if r.f, err = types.NewAddressFromHex(r.From); err != nil {
		return errors.New("from: invalid address")
	}

	r.Sanctum = strings.TrimSpace(r.Sanctum)
	if !util.IsHexAddress(r.Sanctum) {
		return errors.New("sanctum: invalid address")
	}
	if r.s, err = types.NewAddressFromHex(r.Sanctum); err != nil {
		return errors.New("sanctum: invalid address")
	}

	r.Target = strings.TrimSpace(r.Target)
	if !util.IsHexAddress(r.Target) {
		return errors.New("target: invalid address")
	}
	if r.t, err = types.NewAddressFromHex(r.Target); err != nil {
		return errors.New("target: invalid address")
	}

	return nil
}

func (r *SanctumDeregisterForEIP1559Request) FromAddr() *types.Address    { return r.f }
func (r *SanctumDeregisterForEIP1559Request) SanctumAddr() *types.Address { return r.s }
func (r *SanctumDeregisterForEIP1559Request) TargetAddr() *types.Address  { return r.t }

type SanctumDeregisterForEIP1559Response struct {
	UnsignedRLP string `json:"unsigned_rlp"`
	SignedRLP   string `json:"signed_rlp"`
	TxHash      string `json:"tx_hash"`
	R           string `json:"r"`
	S           string `json:"s"`
	V           string `json:"v"`
}

func NewSanctumDeregisterForEIP1559Response(unsigned, signed []byte, txHash *types.Hash, sig *types.Signature) *SanctumDeregisterForEIP1559Response {
	return &SanctumDeregisterForEIP1559Response{
		UnsignedRLP: "0x" + hex.EncodeToString(unsigned),
		SignedRLP:   "0x" + hex.EncodeToString(signed),
		TxHash:      txHash.String(),
		R:           "0x" + sig.R().Text(16),
		S:           "0x" + sig.S().Text(16),
		V:           "0x" + sig.V().Text(16),
	}
}

type SanctumGetAccountsRequest struct {
	Sanctum string `json:"sanctum" example:"0x1234Abcd..."`
	Block   string `json:"block"   example:"latest"`

	s *types.Address
}

func (r *SanctumGetAccountsRequest) ValidateRequest() error {
	r.Sanctum = strings.TrimSpace(r.Sanctum)
	if !util.IsHexAddress(r.Sanctum) {
		return errors.New("sanctum: invalid address")
	}
	var err error
	if r.s, err = types.NewAddressFromHex(r.Sanctum); err != nil {
		return errors.New("sanctum: invalid address")
	}
	if r.Block == "" {
		r.Block = "latest"
	}
	return nil
}

func (r *SanctumGetAccountsRequest) SanctumAddr() *types.Address { return r.s }

type SanctumGetAccountsResponse struct {
	Accounts []string `json:"accounts"`
}

type SanctumAccountCountRequest struct {
	Sanctum string `json:"sanctum" example:"0x1234Abcd..."`
	Block   string `json:"block"   example:"latest"`

	s *types.Address
}

func (r *SanctumAccountCountRequest) ValidateRequest() error {
	r.Sanctum = strings.TrimSpace(r.Sanctum)
	if !util.IsHexAddress(r.Sanctum) {
		return errors.New("sanctum: invalid address")
	}
	var err error
	if r.s, err = types.NewAddressFromHex(r.Sanctum); err != nil {
		return errors.New("sanctum: invalid address")
	}
	if r.Block == "" {
		r.Block = "latest"
	}
	return nil
}

func (r *SanctumAccountCountRequest) SanctumAddr() *types.Address { return r.s }

type SanctumAccountCountResponse struct {
	Count string `json:"count"`
}

type SanctumGetAccountInfoRequest struct {
	Sanctum string `json:"sanctum"  example:"0x1234Abcd..."`
	Account string `json:"account"  example:"0xAbcD1234..."`
	Block   string `json:"block"    example:"latest"`

	s *types.Address
	a *types.Address
}

func (r *SanctumGetAccountInfoRequest) ValidateRequest() error {
	r.Sanctum = strings.TrimSpace(r.Sanctum)
	if !util.IsHexAddress(r.Sanctum) {
		return errors.New("sanctum: invalid address")
	}
	var err error
	if r.s, err = types.NewAddressFromHex(r.Sanctum); err != nil {
		return errors.New("sanctum: invalid address")
	}

	r.Account = strings.TrimSpace(r.Account)
	if !util.IsHexAddress(r.Account) {
		return errors.New("account: invalid address")
	}
	if r.a, err = types.NewAddressFromHex(r.Account); err != nil {
		return errors.New("account: invalid address")
	}

	if r.Block == "" {
		r.Block = "latest"
	}
	return nil
}

func (r *SanctumGetAccountInfoRequest) SanctumAddr() *types.Address { return r.s }
func (r *SanctumGetAccountInfoRequest) AccountAddr() *types.Address { return r.a }

type SanctumGetAccountInfoResponse struct {
	Address         string `json:"address"`
	Role            string `json:"role"`
	RegisteredBlock string `json:"registered_block"`
}

func NewSanctumGetAccountInfoResponse(i *AccountInfo) *SanctumGetAccountInfoResponse {
	return &SanctumGetAccountInfoResponse{
		Address:         i.Addr.String(),
		Role:            i.Role.String(),
		RegisteredBlock: i.RegisteredBlock.String(),
	}
}

type SanctumDepositNativeLegacyRequest struct {
	From    string `json:"from"    example:"0xAbcD1234..."`
	Sanctum string `json:"sanctum" example:"0x1234Abcd..."`
	Value   string `json:"value"   example:"1000000000000000000"`

	f *types.Address
	s *types.Address
	v *big.Int
}

func (r *SanctumDepositNativeLegacyRequest) ValidateRequest() error {
	r.From = strings.TrimSpace(r.From)
	if !util.IsHexAddress(r.From) {
		return errors.New("from: invalid address")
	}
	var err error
	if r.f, err = types.NewAddressFromHex(r.From); err != nil {
		return errors.New("from: invalid address")
	}

	r.Sanctum = strings.TrimSpace(r.Sanctum)
	if !util.IsHexAddress(r.Sanctum) {
		return errors.New("sanctum: invalid address")
	}
	if r.s, err = types.NewAddressFromHex(r.Sanctum); err != nil {
		return errors.New("sanctum: invalid address")
	}

	r.Value = strings.TrimSpace(r.Value)
	v, ok := new(big.Int).SetString(r.Value, 0)
	if !ok || v.Sign() <= 0 {
		return errors.New("value: invalid positive integer")
	}
	r.v = v

	return nil
}

func (r *SanctumDepositNativeLegacyRequest) FromAddr() *types.Address    { return r.f }
func (r *SanctumDepositNativeLegacyRequest) SanctumAddr() *types.Address { return r.s }
func (r *SanctumDepositNativeLegacyRequest) Val() *big.Int               { return r.v }

type SanctumDepositNativeLegacyResponse struct {
	UnsignedRLP string `json:"unsigned_rlp"`
	SignedRLP   string `json:"signed_rlp"`
	TxHash      string `json:"tx_hash"`
	R           string `json:"r"`
	S           string `json:"s"`
	V           string `json:"v"`
}

func NewSanctumDepositNativeLegacyResponse(unsigned, signed []byte, txHash *types.Hash, sig *types.Signature) *SanctumDepositNativeLegacyResponse {
	return &SanctumDepositNativeLegacyResponse{
		UnsignedRLP: "0x" + hex.EncodeToString(unsigned),
		SignedRLP:   "0x" + hex.EncodeToString(signed),
		TxHash:      txHash.String(),
		R:           "0x" + sig.R().Text(16),
		S:           "0x" + sig.S().Text(16),
		V:           "0x" + sig.V().Text(16),
	}
}

type SanctumDepositNativeEIP1559Request struct {
	From    string `json:"from"    example:"0xAbcD1234..."`
	Sanctum string `json:"sanctum" example:"0x1234Abcd..."`
	Value   string `json:"value"   example:"1000000000000000000"`

	f *types.Address
	s *types.Address
	v *big.Int
}

func (r *SanctumDepositNativeEIP1559Request) ValidateRequest() error {
	r.From = strings.TrimSpace(r.From)
	if !util.IsHexAddress(r.From) {
		return errors.New("from: invalid address")
	}
	var err error
	if r.f, err = types.NewAddressFromHex(r.From); err != nil {
		return errors.New("from: invalid address")
	}

	r.Sanctum = strings.TrimSpace(r.Sanctum)
	if !util.IsHexAddress(r.Sanctum) {
		return errors.New("sanctum: invalid address")
	}
	if r.s, err = types.NewAddressFromHex(r.Sanctum); err != nil {
		return errors.New("sanctum: invalid address")
	}

	r.Value = strings.TrimSpace(r.Value)
	v, ok := new(big.Int).SetString(r.Value, 0)
	if !ok || v.Sign() <= 0 {
		return errors.New("value: invalid positive integer")
	}
	r.v = v

	return nil
}

func (r *SanctumDepositNativeEIP1559Request) FromAddr() *types.Address    { return r.f }
func (r *SanctumDepositNativeEIP1559Request) SanctumAddr() *types.Address { return r.s }
func (r *SanctumDepositNativeEIP1559Request) Val() *big.Int               { return r.v }

type SanctumDepositNativeEIP1559Response struct {
	UnsignedRLP string `json:"unsigned_rlp"`
	SignedRLP   string `json:"signed_rlp"`
	TxHash      string `json:"tx_hash"`
	R           string `json:"r"`
	S           string `json:"s"`
	V           string `json:"v"`
}

func NewSanctumDepositNativeEIP1559Response(unsigned, signed []byte, txHash *types.Hash, sig *types.Signature) *SanctumDepositNativeEIP1559Response {
	return &SanctumDepositNativeEIP1559Response{
		UnsignedRLP: "0x" + hex.EncodeToString(unsigned),
		SignedRLP:   "0x" + hex.EncodeToString(signed),
		TxHash:      txHash.String(),
		R:           "0x" + sig.R().Text(16),
		S:           "0x" + sig.S().Text(16),
		V:           "0x" + sig.V().Text(16),
	}
}

type SanctumRequestNativeLegacyRequest struct {
	From    string `json:"from"    example:"0xAbcD1234..."`
	Sanctum string `json:"sanctum" example:"0x1234Abcd..."`
	Amount  string `json:"amount"  example:"1000000000000000000"`

	f *types.Address
	s *types.Address
	a *big.Int
}

func (r *SanctumRequestNativeLegacyRequest) ValidateRequest() error {
	r.From = strings.TrimSpace(r.From)
	if !util.IsHexAddress(r.From) {
		return errors.New("from: invalid address")
	}
	var err error
	if r.f, err = types.NewAddressFromHex(r.From); err != nil {
		return errors.New("from: invalid address")
	}

	r.Sanctum = strings.TrimSpace(r.Sanctum)
	if !util.IsHexAddress(r.Sanctum) {
		return errors.New("sanctum: invalid address")
	}
	if r.s, err = types.NewAddressFromHex(r.Sanctum); err != nil {
		return errors.New("sanctum: invalid address")
	}

	r.Amount = strings.TrimSpace(r.Amount)
	a, ok := new(big.Int).SetString(r.Amount, 0)
	if !ok || a.Sign() <= 0 {
		return errors.New("amount: invalid positive integer")
	}
	r.a = a

	return nil
}

func (r *SanctumRequestNativeLegacyRequest) FromAddr() *types.Address    { return r.f }
func (r *SanctumRequestNativeLegacyRequest) SanctumAddr() *types.Address { return r.s }
func (r *SanctumRequestNativeLegacyRequest) Amt() *big.Int               { return r.a }

type SanctumRequestNativeLegacyResponse struct {
	UnsignedRLP string `json:"unsigned_rlp"`
	SignedRLP   string `json:"signed_rlp"`
	TxHash      string `json:"tx_hash"`
	R           string `json:"r"`
	S           string `json:"s"`
	V           string `json:"v"`
}

func NewSanctumRequestNativeLegacyResponse(unsigned, signed []byte, txHash *types.Hash, sig *types.Signature) *SanctumRequestNativeLegacyResponse {
	return &SanctumRequestNativeLegacyResponse{
		UnsignedRLP: "0x" + hex.EncodeToString(unsigned),
		SignedRLP:   "0x" + hex.EncodeToString(signed),
		TxHash:      txHash.String(),
		R:           "0x" + sig.R().Text(16),
		S:           "0x" + sig.S().Text(16),
		V:           "0x" + sig.V().Text(16),
	}
}

type SanctumRequestNativeEIP1559Request struct {
	From    string `json:"from"    example:"0xAbcD1234..."`
	Sanctum string `json:"sanctum" example:"0x1234Abcd..."`
	Amount  string `json:"amount"  example:"1000000000000000000"`

	f *types.Address
	s *types.Address
	a *big.Int
}

func (r *SanctumRequestNativeEIP1559Request) ValidateRequest() error {
	r.From = strings.TrimSpace(r.From)
	if !util.IsHexAddress(r.From) {
		return errors.New("from: invalid address")
	}
	var err error
	if r.f, err = types.NewAddressFromHex(r.From); err != nil {
		return errors.New("from: invalid address")
	}

	r.Sanctum = strings.TrimSpace(r.Sanctum)
	if !util.IsHexAddress(r.Sanctum) {
		return errors.New("sanctum: invalid address")
	}
	if r.s, err = types.NewAddressFromHex(r.Sanctum); err != nil {
		return errors.New("sanctum: invalid address")
	}

	r.Amount = strings.TrimSpace(r.Amount)
	a, ok := new(big.Int).SetString(r.Amount, 0)
	if !ok || a.Sign() <= 0 {
		return errors.New("amount: invalid positive integer")
	}
	r.a = a

	return nil
}

func (r *SanctumRequestNativeEIP1559Request) FromAddr() *types.Address    { return r.f }
func (r *SanctumRequestNativeEIP1559Request) SanctumAddr() *types.Address { return r.s }
func (r *SanctumRequestNativeEIP1559Request) Amt() *big.Int               { return r.a }

type SanctumRequestNativeEIP1559Response struct {
	UnsignedRLP string `json:"unsigned_rlp"`
	SignedRLP   string `json:"signed_rlp"`
	TxHash      string `json:"tx_hash"`
	R           string `json:"r"`
	S           string `json:"s"`
	V           string `json:"v"`
}

func NewSanctumRequestNativeEIP1559Response(unsigned, signed []byte, txHash *types.Hash, sig *types.Signature) *SanctumRequestNativeEIP1559Response {
	return &SanctumRequestNativeEIP1559Response{
		UnsignedRLP: "0x" + hex.EncodeToString(unsigned),
		SignedRLP:   "0x" + hex.EncodeToString(signed),
		TxHash:      txHash.String(),
		R:           "0x" + sig.R().Text(16),
		S:           "0x" + sig.S().Text(16),
		V:           "0x" + sig.V().Text(16),
	}
}

type SanctumApproveNativeLegacyRequest struct {
	From    string `json:"from"    example:"0xAbcD1234..."`
	Sanctum string `json:"sanctum" example:"0x1234Abcd..."`
	User    string `json:"user"    example:"0xDeadBeef..."`
	Amount  string `json:"amount"  example:"1000000000000000000"`

	f *types.Address
	s *types.Address
	u *types.Address
	a *big.Int
}

func (r *SanctumApproveNativeLegacyRequest) ValidateRequest() error {
	r.From = strings.TrimSpace(r.From)
	if !util.IsHexAddress(r.From) {
		return errors.New("from: invalid address")
	}
	var err error
	if r.f, err = types.NewAddressFromHex(r.From); err != nil {
		return errors.New("from: invalid address")
	}

	r.Sanctum = strings.TrimSpace(r.Sanctum)
	if !util.IsHexAddress(r.Sanctum) {
		return errors.New("sanctum: invalid address")
	}
	if r.s, err = types.NewAddressFromHex(r.Sanctum); err != nil {
		return errors.New("sanctum: invalid address")
	}

	r.User = strings.TrimSpace(r.User)
	if !util.IsHexAddress(r.User) {
		return errors.New("user: invalid address")
	}
	if r.u, err = types.NewAddressFromHex(r.User); err != nil {
		return errors.New("user: invalid address")
	}

	r.Amount = strings.TrimSpace(r.Amount)
	a, ok := new(big.Int).SetString(r.Amount, 0)
	if !ok || a.Sign() <= 0 {
		return errors.New("amount: invalid positive integer")
	}
	r.a = a

	return nil
}

func (r *SanctumApproveNativeLegacyRequest) FromAddr() *types.Address    { return r.f }
func (r *SanctumApproveNativeLegacyRequest) SanctumAddr() *types.Address { return r.s }
func (r *SanctumApproveNativeLegacyRequest) UserAddr() *types.Address    { return r.u }
func (r *SanctumApproveNativeLegacyRequest) Amt() *big.Int               { return r.a }

type SanctumApproveNativeLegacyResponse struct {
	UnsignedRLP string `json:"unsigned_rlp"`
	SignedRLP   string `json:"signed_rlp"`
	TxHash      string `json:"tx_hash"`
	R           string `json:"r"`
	S           string `json:"s"`
	V           string `json:"v"`
}

func NewSanctumApproveNativeLegacyResponse(unsigned, signed []byte, txHash *types.Hash, sig *types.Signature) *SanctumApproveNativeLegacyResponse {
	return &SanctumApproveNativeLegacyResponse{
		UnsignedRLP: "0x" + hex.EncodeToString(unsigned),
		SignedRLP:   "0x" + hex.EncodeToString(signed),
		TxHash:      txHash.String(),
		R:           "0x" + sig.R().Text(16),
		S:           "0x" + sig.S().Text(16),
		V:           "0x" + sig.V().Text(16),
	}
}

type SanctumApproveNativeEIP1559Request struct {
	From    string `json:"from"    example:"0xAbcD1234..."`
	Sanctum string `json:"sanctum" example:"0x1234Abcd..."`
	User    string `json:"user"    example:"0xDeadBeef..."`
	Amount  string `json:"amount"  example:"1000000000000000000"`

	f *types.Address
	s *types.Address
	u *types.Address
	a *big.Int
}

func (r *SanctumApproveNativeEIP1559Request) ValidateRequest() error {
	r.From = strings.TrimSpace(r.From)
	if !util.IsHexAddress(r.From) {
		return errors.New("from: invalid address")
	}
	var err error
	if r.f, err = types.NewAddressFromHex(r.From); err != nil {
		return errors.New("from: invalid address")
	}

	r.Sanctum = strings.TrimSpace(r.Sanctum)
	if !util.IsHexAddress(r.Sanctum) {
		return errors.New("sanctum: invalid address")
	}
	if r.s, err = types.NewAddressFromHex(r.Sanctum); err != nil {
		return errors.New("sanctum: invalid address")
	}

	r.User = strings.TrimSpace(r.User)
	if !util.IsHexAddress(r.User) {
		return errors.New("user: invalid address")
	}
	if r.u, err = types.NewAddressFromHex(r.User); err != nil {
		return errors.New("user: invalid address")
	}

	r.Amount = strings.TrimSpace(r.Amount)
	a, ok := new(big.Int).SetString(r.Amount, 0)
	if !ok || a.Sign() <= 0 {
		return errors.New("amount: invalid positive integer")
	}
	r.a = a

	return nil
}

func (r *SanctumApproveNativeEIP1559Request) FromAddr() *types.Address    { return r.f }
func (r *SanctumApproveNativeEIP1559Request) SanctumAddr() *types.Address { return r.s }
func (r *SanctumApproveNativeEIP1559Request) UserAddr() *types.Address    { return r.u }
func (r *SanctumApproveNativeEIP1559Request) Amt() *big.Int               { return r.a }

type SanctumApproveNativeEIP1559Response struct {
	UnsignedRLP string `json:"unsigned_rlp"`
	SignedRLP   string `json:"signed_rlp"`
	TxHash      string `json:"tx_hash"`
	R           string `json:"r"`
	S           string `json:"s"`
	V           string `json:"v"`
}

func NewSanctumApproveNativeEIP1559Response(unsigned, signed []byte, txHash *types.Hash, sig *types.Signature) *SanctumApproveNativeEIP1559Response {
	return &SanctumApproveNativeEIP1559Response{
		UnsignedRLP: "0x" + hex.EncodeToString(unsigned),
		SignedRLP:   "0x" + hex.EncodeToString(signed),
		TxHash:      txHash.String(),
		R:           "0x" + sig.R().Text(16),
		S:           "0x" + sig.S().Text(16),
		V:           "0x" + sig.V().Text(16),
	}
}

type SanctumApproveNativeAllLegacyRequest struct {
	From    string `json:"from"    example:"0xAbcD1234..."`
	Sanctum string `json:"sanctum" example:"0x1234Abcd..."`
	User    string `json:"user"    example:"0xDeadBeef..."`

	f *types.Address
	s *types.Address
	u *types.Address
}

func (r *SanctumApproveNativeAllLegacyRequest) ValidateRequest() error {
	r.From = strings.TrimSpace(r.From)
	if !util.IsHexAddress(r.From) {
		return errors.New("from: invalid address")
	}
	var err error
	if r.f, err = types.NewAddressFromHex(r.From); err != nil {
		return errors.New("from: invalid address")
	}

	r.Sanctum = strings.TrimSpace(r.Sanctum)
	if !util.IsHexAddress(r.Sanctum) {
		return errors.New("sanctum: invalid address")
	}
	if r.s, err = types.NewAddressFromHex(r.Sanctum); err != nil {
		return errors.New("sanctum: invalid address")
	}

	r.User = strings.TrimSpace(r.User)
	if !util.IsHexAddress(r.User) {
		return errors.New("user: invalid address")
	}
	if r.u, err = types.NewAddressFromHex(r.User); err != nil {
		return errors.New("user: invalid address")
	}

	return nil
}

func (r *SanctumApproveNativeAllLegacyRequest) FromAddr() *types.Address    { return r.f }
func (r *SanctumApproveNativeAllLegacyRequest) SanctumAddr() *types.Address { return r.s }
func (r *SanctumApproveNativeAllLegacyRequest) UserAddr() *types.Address    { return r.u }

type SanctumApproveNativeAllLegacyResponse struct {
	UnsignedRLP string `json:"unsigned_rlp"`
	SignedRLP   string `json:"signed_rlp"`
	TxHash      string `json:"tx_hash"`
	R           string `json:"r"`
	S           string `json:"s"`
	V           string `json:"v"`
}

func NewSanctumApproveNativeAllLegacyResponse(unsigned, signed []byte, txHash *types.Hash, sig *types.Signature) *SanctumApproveNativeAllLegacyResponse {
	return &SanctumApproveNativeAllLegacyResponse{
		UnsignedRLP: "0x" + hex.EncodeToString(unsigned),
		SignedRLP:   "0x" + hex.EncodeToString(signed),
		TxHash:      txHash.String(),
		R:           "0x" + sig.R().Text(16),
		S:           "0x" + sig.S().Text(16),
		V:           "0x" + sig.V().Text(16),
	}
}

type SanctumApproveNativeAllEIP1559Request struct {
	From    string `json:"from"    example:"0xAbcD1234..."`
	Sanctum string `json:"sanctum" example:"0x1234Abcd..."`
	User    string `json:"user"    example:"0xDeadBeef..."`

	f *types.Address
	s *types.Address
	u *types.Address
}

func (r *SanctumApproveNativeAllEIP1559Request) ValidateRequest() error {
	r.From = strings.TrimSpace(r.From)
	if !util.IsHexAddress(r.From) {
		return errors.New("from: invalid address")
	}
	var err error
	if r.f, err = types.NewAddressFromHex(r.From); err != nil {
		return errors.New("from: invalid address")
	}

	r.Sanctum = strings.TrimSpace(r.Sanctum)
	if !util.IsHexAddress(r.Sanctum) {
		return errors.New("sanctum: invalid address")
	}
	if r.s, err = types.NewAddressFromHex(r.Sanctum); err != nil {
		return errors.New("sanctum: invalid address")
	}

	r.User = strings.TrimSpace(r.User)
	if !util.IsHexAddress(r.User) {
		return errors.New("user: invalid address")
	}
	if r.u, err = types.NewAddressFromHex(r.User); err != nil {
		return errors.New("user: invalid address")
	}

	return nil
}

func (r *SanctumApproveNativeAllEIP1559Request) FromAddr() *types.Address    { return r.f }
func (r *SanctumApproveNativeAllEIP1559Request) SanctumAddr() *types.Address { return r.s }
func (r *SanctumApproveNativeAllEIP1559Request) UserAddr() *types.Address    { return r.u }

type SanctumApproveNativeAllEIP1559Response struct {
	UnsignedRLP string `json:"unsigned_rlp"`
	SignedRLP   string `json:"signed_rlp"`
	TxHash      string `json:"tx_hash"`
	R           string `json:"r"`
	S           string `json:"s"`
	V           string `json:"v"`
}

func NewSanctumApproveNativeAllEIP1559Response(unsigned, signed []byte, txHash *types.Hash, sig *types.Signature) *SanctumApproveNativeAllEIP1559Response {
	return &SanctumApproveNativeAllEIP1559Response{
		UnsignedRLP: "0x" + hex.EncodeToString(unsigned),
		SignedRLP:   "0x" + hex.EncodeToString(signed),
		TxHash:      txHash.String(),
		R:           "0x" + sig.R().Text(16),
		S:           "0x" + sig.S().Text(16),
		V:           "0x" + sig.V().Text(16),
	}
}

type SanctumWithdrawNativeLegacyRequest struct {
	From    string `json:"from"    example:"0xAbcD1234..."`
	Sanctum string `json:"sanctum" example:"0x1234Abcd..."`
	Amount  string `json:"amount"  example:"1000000000000000000"`

	f *types.Address
	s *types.Address
	a *big.Int
}

func (r *SanctumWithdrawNativeLegacyRequest) ValidateRequest() error {
	r.From = strings.TrimSpace(r.From)
	if !util.IsHexAddress(r.From) {
		return errors.New("from: invalid address")
	}
	var err error
	if r.f, err = types.NewAddressFromHex(r.From); err != nil {
		return errors.New("from: invalid address")
	}

	r.Sanctum = strings.TrimSpace(r.Sanctum)
	if !util.IsHexAddress(r.Sanctum) {
		return errors.New("sanctum: invalid address")
	}
	if r.s, err = types.NewAddressFromHex(r.Sanctum); err != nil {
		return errors.New("sanctum: invalid address")
	}

	r.Amount = strings.TrimSpace(r.Amount)
	a, ok := new(big.Int).SetString(r.Amount, 0)
	if !ok || a.Sign() <= 0 {
		return errors.New("amount: invalid positive integer")
	}
	r.a = a

	return nil
}

func (r *SanctumWithdrawNativeLegacyRequest) FromAddr() *types.Address    { return r.f }
func (r *SanctumWithdrawNativeLegacyRequest) SanctumAddr() *types.Address { return r.s }
func (r *SanctumWithdrawNativeLegacyRequest) ToAmount() *big.Int          { return r.a }

type SanctumWithdrawNativeLegacyResponse struct {
	UnsignedRLP string `json:"unsigned_rlp"`
	SignedRLP   string `json:"signed_rlp"`
	TxHash      string `json:"tx_hash"`
	R           string `json:"r"`
	S           string `json:"s"`
	V           string `json:"v"`
}

func NewSanctumWithdrawNativeLegacyResponse(unsigned, signed []byte, txHash *types.Hash, sig *types.Signature) *SanctumWithdrawNativeLegacyResponse {
	return &SanctumWithdrawNativeLegacyResponse{
		UnsignedRLP: "0x" + hex.EncodeToString(unsigned),
		SignedRLP:   "0x" + hex.EncodeToString(signed),
		TxHash:      txHash.String(),
		R:           "0x" + sig.R().Text(16),
		S:           "0x" + sig.S().Text(16),
		V:           "0x" + sig.V().Text(16),
	}
}

type SanctumWithdrawNativeEIP1559Request struct {
	From    string `json:"from"    example:"0xAbcD1234..."`
	Sanctum string `json:"sanctum" example:"0x1234Abcd..."`
	Amount  string `json:"amount"  example:"1000000000000000000"`

	f *types.Address
	s *types.Address
	a *big.Int
}

func (r *SanctumWithdrawNativeEIP1559Request) ValidateRequest() error {
	r.From = strings.TrimSpace(r.From)
	if !util.IsHexAddress(r.From) {
		return errors.New("from: invalid address")
	}
	var err error
	if r.f, err = types.NewAddressFromHex(r.From); err != nil {
		return errors.New("from: invalid address")
	}

	r.Sanctum = strings.TrimSpace(r.Sanctum)
	if !util.IsHexAddress(r.Sanctum) {
		return errors.New("sanctum: invalid address")
	}
	if r.s, err = types.NewAddressFromHex(r.Sanctum); err != nil {
		return errors.New("sanctum: invalid address")
	}

	r.Amount = strings.TrimSpace(r.Amount)
	a, ok := new(big.Int).SetString(r.Amount, 0)
	if !ok || a.Sign() <= 0 {
		return errors.New("amount: invalid positive integer")
	}
	r.a = a

	return nil
}

func (r *SanctumWithdrawNativeEIP1559Request) FromAddr() *types.Address    { return r.f }
func (r *SanctumWithdrawNativeEIP1559Request) SanctumAddr() *types.Address { return r.s }
func (r *SanctumWithdrawNativeEIP1559Request) ToAmount() *big.Int          { return r.a }

type SanctumWithdrawNativeEIP1559Response struct {
	UnsignedRLP string `json:"unsigned_rlp"`
	SignedRLP   string `json:"signed_rlp"`
	TxHash      string `json:"tx_hash"`
	R           string `json:"r"`
	S           string `json:"s"`
	V           string `json:"v"`
}

func NewSanctumWithdrawNativeEIP1559Response(unsigned, signed []byte, txHash *types.Hash, sig *types.Signature) *SanctumWithdrawNativeEIP1559Response {
	return &SanctumWithdrawNativeEIP1559Response{
		UnsignedRLP: "0x" + hex.EncodeToString(unsigned),
		SignedRLP:   "0x" + hex.EncodeToString(signed),
		TxHash:      txHash.String(),
		R:           "0x" + sig.R().Text(16),
		S:           "0x" + sig.S().Text(16),
		V:           "0x" + sig.V().Text(16),
	}
}

type SanctumWithdrawNativeAllLegacyRequest struct {
	From    string `json:"from"    example:"0xAbcD1234..."`
	Sanctum string `json:"sanctum" example:"0x1234Abcd..."`

	f *types.Address
	s *types.Address
}

func (r *SanctumWithdrawNativeAllLegacyRequest) ValidateRequest() error {
	r.From = strings.TrimSpace(r.From)
	if !util.IsHexAddress(r.From) {
		return errors.New("from: invalid address")
	}
	var err error
	if r.f, err = types.NewAddressFromHex(r.From); err != nil {
		return errors.New("from: invalid address")
	}

	r.Sanctum = strings.TrimSpace(r.Sanctum)
	if !util.IsHexAddress(r.Sanctum) {
		return errors.New("sanctum: invalid address")
	}
	if r.s, err = types.NewAddressFromHex(r.Sanctum); err != nil {
		return errors.New("sanctum: invalid address")
	}

	return nil
}

func (r *SanctumWithdrawNativeAllLegacyRequest) FromAddr() *types.Address    { return r.f }
func (r *SanctumWithdrawNativeAllLegacyRequest) SanctumAddr() *types.Address { return r.s }

type SanctumWithdrawNativeAllLegacyResponse struct {
	UnsignedRLP string `json:"unsigned_rlp"`
	SignedRLP   string `json:"signed_rlp"`
	TxHash      string `json:"tx_hash"`
	R           string `json:"r"`
	S           string `json:"s"`
	V           string `json:"v"`
}

func NewSanctumWithdrawNativeAllLegacyResponse(unsigned, signed []byte, txHash *types.Hash, sig *types.Signature) *SanctumWithdrawNativeAllLegacyResponse {
	return &SanctumWithdrawNativeAllLegacyResponse{
		UnsignedRLP: "0x" + hex.EncodeToString(unsigned),
		SignedRLP:   "0x" + hex.EncodeToString(signed),
		TxHash:      txHash.String(),
		R:           "0x" + sig.R().Text(16),
		S:           "0x" + sig.S().Text(16),
		V:           "0x" + sig.V().Text(16),
	}
}

type SanctumWithdrawNativeAllEIP1559Request struct {
	From    string `json:"from"    example:"0xAbcD1234..."`
	Sanctum string `json:"sanctum" example:"0x1234Abcd..."`

	f *types.Address
	s *types.Address
}

func (r *SanctumWithdrawNativeAllEIP1559Request) ValidateRequest() error {
	r.From = strings.TrimSpace(r.From)
	if !util.IsHexAddress(r.From) {
		return errors.New("from: invalid address")
	}
	var err error
	if r.f, err = types.NewAddressFromHex(r.From); err != nil {
		return errors.New("from: invalid address")
	}

	r.Sanctum = strings.TrimSpace(r.Sanctum)
	if !util.IsHexAddress(r.Sanctum) {
		return errors.New("sanctum: invalid address")
	}
	if r.s, err = types.NewAddressFromHex(r.Sanctum); err != nil {
		return errors.New("sanctum: invalid address")
	}

	return nil
}

func (r *SanctumWithdrawNativeAllEIP1559Request) FromAddr() *types.Address    { return r.f }
func (r *SanctumWithdrawNativeAllEIP1559Request) SanctumAddr() *types.Address { return r.s }

type SanctumWithdrawNativeAllEIP1559Response struct {
	UnsignedRLP string `json:"unsigned_rlp"`
	SignedRLP   string `json:"signed_rlp"`
	TxHash      string `json:"tx_hash"`
	R           string `json:"r"`
	S           string `json:"s"`
	V           string `json:"v"`
}

func NewSanctumWithdrawNativeAllEIP1559Response(unsigned, signed []byte, txHash *types.Hash, sig *types.Signature) *SanctumWithdrawNativeAllEIP1559Response {
	return &SanctumWithdrawNativeAllEIP1559Response{
		UnsignedRLP: "0x" + hex.EncodeToString(unsigned),
		SignedRLP:   "0x" + hex.EncodeToString(signed),
		TxHash:      txHash.String(),
		R:           "0x" + sig.R().Text(16),
		S:           "0x" + sig.S().Text(16),
		V:           "0x" + sig.V().Text(16),
	}
}

type SanctumNativeBalanceRequest struct {
	Sanctum string `json:"sanctum" example:"0x1234Abcd..."`
	Block   string `json:"block"   example:"latest"`

	s *types.Address
}

func (r *SanctumNativeBalanceRequest) ValidateRequest() error {
	r.Sanctum = strings.TrimSpace(r.Sanctum)
	if !util.IsHexAddress(r.Sanctum) {
		return errors.New("sanctum: invalid address")
	}
	var err error
	if r.s, err = types.NewAddressFromHex(r.Sanctum); err != nil {
		return errors.New("sanctum: invalid address")
	}
	if r.Block == "" {
		r.Block = "latest"
	}
	return nil
}

func (r *SanctumNativeBalanceRequest) SanctumAddr() *types.Address { return r.s }

type SanctumNativeBalanceResponse struct {
	Balance string `json:"balance"`
}

type SanctumNativeAvailableRequest struct {
	Sanctum string `json:"sanctum" example:"0x1234Abcd..."`
	Block   string `json:"block"   example:"latest"`

	s *types.Address
}

func (r *SanctumNativeAvailableRequest) ValidateRequest() error {
	r.Sanctum = strings.TrimSpace(r.Sanctum)
	if !util.IsHexAddress(r.Sanctum) {
		return errors.New("sanctum: invalid address")
	}
	var err error
	if r.s, err = types.NewAddressFromHex(r.Sanctum); err != nil {
		return errors.New("sanctum: invalid address")
	}
	if r.Block == "" {
		r.Block = "latest"
	}
	return nil
}

func (r *SanctumNativeAvailableRequest) SanctumAddr() *types.Address { return r.s }

type SanctumNativeAvailableResponse struct {
	Available string `json:"available"`
}

type SanctumNativeAllocationRequest struct {
	Sanctum string `json:"sanctum" example:"0x1234Abcd..."`
	User    string `json:"user"    example:"0xAbcD1234..."`
	Block   string `json:"block"   example:"latest"`

	s *types.Address
	u *types.Address
}

func (r *SanctumNativeAllocationRequest) ValidateRequest() error {
	r.Sanctum = strings.TrimSpace(r.Sanctum)
	if !util.IsHexAddress(r.Sanctum) {
		return errors.New("sanctum: invalid address")
	}
	var err error
	if r.s, err = types.NewAddressFromHex(r.Sanctum); err != nil {
		return errors.New("sanctum: invalid address")
	}

	r.User = strings.TrimSpace(r.User)
	if !util.IsHexAddress(r.User) {
		return errors.New("user: invalid address")
	}
	if r.u, err = types.NewAddressFromHex(r.User); err != nil {
		return errors.New("user: invalid address")
	}

	if r.Block == "" {
		r.Block = "latest"
	}
	return nil
}

func (r *SanctumNativeAllocationRequest) SanctumAddr() *types.Address { return r.s }
func (r *SanctumNativeAllocationRequest) UserAddr() *types.Address    { return r.u }

type SanctumNativeAllocationResponse struct {
	Allocation string `json:"allocation"`
}

type SanctumNativePendingRequest struct {
	Sanctum string `json:"sanctum" example:"0x1234Abcd..."`
	User    string `json:"user"    example:"0xAbcD1234..."`
	Block   string `json:"block"   example:"latest"`

	s *types.Address
	u *types.Address
}

func (r *SanctumNativePendingRequest) ValidateRequest() error {
	r.Sanctum = strings.TrimSpace(r.Sanctum)
	if !util.IsHexAddress(r.Sanctum) {
		return errors.New("sanctum: invalid address")
	}
	var err error
	if r.s, err = types.NewAddressFromHex(r.Sanctum); err != nil {
		return errors.New("sanctum: invalid address")
	}

	r.User = strings.TrimSpace(r.User)
	if !util.IsHexAddress(r.User) {
		return errors.New("user: invalid address")
	}
	if r.u, err = types.NewAddressFromHex(r.User); err != nil {
		return errors.New("user: invalid address")
	}

	if r.Block == "" {
		r.Block = "latest"
	}
	return nil
}

func (r *SanctumNativePendingRequest) SanctumAddr() *types.Address { return r.s }
func (r *SanctumNativePendingRequest) UserAddr() *types.Address    { return r.u }

type SanctumNativePendingResponse struct {
	Pending string `json:"pending"`
}
