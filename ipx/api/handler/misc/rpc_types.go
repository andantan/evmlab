package misc

import (
	"errors"
	"fmt"
	"math/big"
	"strings"

	"github.com/andantan/evmlab/core/types"
	"github.com/andantan/evmlab/internal/util"
)

type RawRPCRequest struct {
	Method string `json:"method"`
	Params []any  `json:"params"`
}

func (r *RawRPCRequest) ValidateRequest() error {
	r.Method = strings.TrimSpace(r.Method)
	if r.Method == "" {
		return errors.New("method is required")
	}
	return nil
}

type RawRPCResponse struct {
	Result any `json:"result"`
}

func NewRawRPCResponse(r any) *RawRPCResponse {
	return &RawRPCResponse{
		Result: r,
	}
}

type ChainIDResponse struct {
	ChainID    uint64 `json:"chain_id"`
	ChainIDHex string `json:"chain_id_hex"`
}

func NewChainIDResponse(c uint64, s string) *ChainIDResponse {
	return &ChainIDResponse{
		ChainID:    c,
		ChainIDHex: s,
	}
}

type GasPriceResponse struct {
	GasPrice    string `json:"gas_price"`
	GasPriceHex string `json:"gas_price_hex"`
	Wei         string `json:"wei"`
	Gwei        string `json:"gwei"`
	Ether       string `json:"ether"`
}

func NewGasPriceResponse(g *big.Int, s string) *GasPriceResponse {
	return &GasPriceResponse{
		GasPrice:    g.String(),
		GasPriceHex: s,
		Wei:         g.String(),
		Gwei:        types.WeiToGwei(g),
		Ether:       types.WeiToEther(g),
	}
}

type BaseFeePerGasResponse struct {
	BaseFeePerGas    string `json:"base_fee_per_gas"`
	BaseFeePerGasHex string `json:"base_fee_per_gas_hex"`
	Wei              string `json:"wei"`
	Gwei             string `json:"gwei"`
	Ether            string `json:"ether"`
}

func NewBaseFeePerGasResponse(f *big.Int, s string) *BaseFeePerGasResponse {
	return &BaseFeePerGasResponse{
		BaseFeePerGas:    f.String(),
		BaseFeePerGasHex: s,
		Wei:              f.String(),
		Gwei:             types.WeiToGwei(f),
		Ether:            types.WeiToEther(f),
	}
}

type MaxPriorityFeePerGasResponse struct {
	MaxPriorityFeePerGas    string `json:"max_priority_fee_per_gas"`
	MaxPriorityFeePerGasHex string `json:"max_priority_fee_per_gas_hex"`
	Wei                     string `json:"wei"`
	Gwei                    string `json:"gwei"`
	Ether                   string `json:"ether"`
}

func NewMaxPriorityFeePerGasResponse(f *big.Int, s string) *MaxPriorityFeePerGasResponse {
	return &MaxPriorityFeePerGasResponse{
		MaxPriorityFeePerGas:    f.String(),
		MaxPriorityFeePerGasHex: s,
		Wei:                     f.String(),
		Gwei:                    types.WeiToGwei(f),
		Ether:                   types.WeiToEther(f),
	}
}

type MaxFeePerGasResponse struct {
	MaxFeePerGas    string `json:"max_fee_per_gas"`
	MaxFeePerGasHex string `json:"max_fee_per_gas_hex"`
	Wei             string `json:"wei"`
	Gwei            string `json:"gwei"`
	Ether           string `json:"ether"`
}

func NewMaxFeePerGasResponse(f *big.Int, s string) *MaxFeePerGasResponse {
	return &MaxFeePerGasResponse{
		MaxFeePerGas:    f.String(),
		MaxFeePerGasHex: s,
		Wei:             f.String(),
		Gwei:            types.WeiToGwei(f),
		Ether:           types.WeiToEther(f),
	}
}

type BlockNumberResponse struct {
	BlockNumber    uint64 `json:"block_number"`
	BlockNumberHex string `json:"block_number_hex"`
}

func NewBlockNumberResponse(n uint64, s string) *BlockNumberResponse {
	return &BlockNumberResponse{
		BlockNumber:    n,
		BlockNumberHex: s,
	}
}

type NonceRequest struct {
	Address string `json:"address" example:"0xEbD69375..."`
	Block   string `json:"block"   example:"pending"`
}

func (r *NonceRequest) ValidateRequest() error {
	r.Address = strings.TrimSpace(r.Address)
	if r.Address == "" {
		return errors.New("address is required")
	}
	if r.Block == "" {
		r.Block = "pending"
	}
	return nil
}

type NonceResponse struct {
	Nonce    uint64 `json:"nonce"`
	NonceHex string `json:"nonce_hex"`
}

func NewNonceResponse(n uint64, s string) *NonceResponse {
	return &NonceResponse{
		Nonce:    n,
		NonceHex: s,
	}
}

type BalanceRequest struct {
	Address string `json:"address" example:"0xEbD69375..."`
	Block   string `json:"block"   example:"latest"`
}

func (r *BalanceRequest) ValidateRequest() error {
	r.Address = strings.TrimSpace(r.Address)
	r.Block = strings.TrimSpace(r.Block)
	if r.Address == "" {
		return errors.New("address is required")
	}
	if _, err := util.ParseHex(r.Address); err != nil {
		return errors.New("address: " + err.Error())
	}
	return nil
}

type BalanceResponse struct {
	Balance    string `json:"balance"`
	BalanceHex string `json:"balance_hex"`
	Wei        string `json:"wei"`
	Gwei       string `json:"gwei"`
	Ether      string `json:"ether"`
}

func NewBalanceResponse(w *big.Int, s string) *BalanceResponse {
	return &BalanceResponse{
		Balance:    w.String(),
		BalanceHex: s,
		Wei:        w.String(),
		Gwei:       types.WeiToGwei(w),
		Ether:      types.WeiToEther(w),
	}
}

type CodeRequest struct {
	Address string `json:"address" example:"0xEbD69375..."`
	Block   string `json:"block"   example:"latest"`
}

func (r *CodeRequest) ValidateRequest() error {
	r.Address = strings.TrimSpace(r.Address)
	if r.Address == "" {
		return errors.New("address is required")
	}
	if _, err := util.ParseHex(r.Address); err != nil {
		return errors.New("address: " + err.Error())
	}

	r.Block = strings.TrimSpace(r.Block)
	if r.Block == "" {
		r.Block = "latest"
	}

	return nil
}

type CodeResponse struct {
	Code       string `json:"code"`
	IsContract bool   `json:"is_contract"`
}

func NewCodeResponse(c string) *CodeResponse {
	return &CodeResponse{
		Code:       c,
		IsContract: c != "" && c != "0x" && c != "0x0",
	}
}

type TransactionRequest struct {
	TxHash string `json:"tx_hash" example:"0xabc123..."`
}

func (r *TransactionRequest) ValidateRequest() error {
	r.TxHash = strings.TrimSpace(r.TxHash)
	if r.TxHash == "" {
		return errors.New("tx_hash is required")
	}
	bare := strings.TrimPrefix(r.TxHash, "0x")
	if len(bare) != types.HashHexLength {
		return fmt.Errorf("tx_hash: must be %d hex chars (got %d)", types.HashHexLength, len(bare))
	}
	return nil
}

type TransactionReceiptRequest struct {
	TxHash string `json:"tx_hash" example:"0xabc123..."`
}

func (r *TransactionReceiptRequest) ValidateRequest() error {
	r.TxHash = strings.TrimSpace(r.TxHash)
	if r.TxHash == "" {
		return errors.New("tx_hash is required")
	}
	bare := strings.TrimPrefix(r.TxHash, "0x")
	if len(bare) != types.HashHexLength {
		return fmt.Errorf("tx_hash: must be %d hex chars (got %d)", types.HashHexLength, len(bare))
	}
	return nil
}

type EstimateGasRequest struct {
	From  string `json:"from"  example:"0xEbD69375..."`
	To    string `json:"to"    example:"0x8336c196..."`
	Value string `json:"value" example:"1000000000000000000"`
	Data  string `json:"data"  example:"0x"`
	Block string `json:"block" example:"latest"`

	p map[string]string
}

func (r *EstimateGasRequest) ValidateRequest() error {
	r.From = strings.TrimSpace(r.From)
	if r.From == "" {
		return errors.New("from is required")
	}
	if r.Block == "" {
		r.Block = "latest"
	}
	r.p = make(map[string]string, 4)
	r.p["from"] = r.From

	if r.To != "" {
		r.p["to"] = r.To
	}
	if r.Value != "" && r.Value != "0x" {
		if strings.HasPrefix(r.Value, "0x") {
			r.p["value"] = r.Value
		} else {
			n, ok := new(big.Int).SetString(r.Value, 10)
			if !ok {
				return errors.New("value: invalid integer")
			}
			r.p["value"] = "0x" + n.Text(16)
		}
	}
	if r.Data != "" {
		r.p["data"] = r.Data
	}
	return nil
}

func (r *EstimateGasRequest) Params() map[string]string {
	return r.p
}

type EstimateGasResponse struct {
	GasLimit    uint64 `json:"gas_limit"`
	GasLimitHex string `json:"gas_limit_hex"`
	Wei         string `json:"wei"`
	Gwei        string `json:"gwei"`
	Ether       string `json:"ether"`
}

func NewEstimateGasResponse(g uint64, s string) *EstimateGasResponse {
	w := new(big.Int).SetUint64(g)
	return &EstimateGasResponse{
		GasLimit:    g,
		GasLimitHex: s,
		Wei:         w.String(),
		Gwei:        types.WeiToGwei(w),
		Ether:       types.WeiToEther(w),
	}
}

type CallRequest struct {
	From  string `json:"from"  example:"0xEbD69375..."`
	To    string `json:"to"    example:"0x8336c196..."`
	Data  string `json:"data"  example:"0x70a08231..."`
	Block string `json:"block" example:"latest"`

	p map[string]string
}

func (r *CallRequest) ValidateRequest() error {
	r.To = strings.TrimSpace(r.To)
	if r.To == "" {
		return errors.New("to is required")
	}
	r.Data = strings.TrimSpace(r.Data)
	if r.Data == "" {
		return errors.New("data is required")
	}
	if r.Block == "" {
		r.Block = "latest"
	}
	r.p = make(map[string]string, 3)
	r.p["to"] = r.To
	r.p["data"] = r.Data
	if r.From != "" {
		r.p["from"] = r.From
	}
	return nil
}

func (r *CallRequest) Params() map[string]string {
	return r.p
}

type SendTransactionRequest struct {
	RawTx string `json:"raw_tx" example:"0x02f8..."`
}

func (r *SendTransactionRequest) ValidateRequest() error {
	r.RawTx = strings.TrimSpace(r.RawTx)
	if r.RawTx == "" {
		return errors.New("raw_tx is required")
	}
	if _, err := util.ParseHex(r.RawTx); err != nil {
		return errors.New("raw_tx: " + err.Error())
	}
	return nil
}

type SendTransactionResponse struct {
	TxHash string `json:"tx_hash"`
}

func NewSendTransactionResponse(s string) *SendTransactionResponse {
	return &SendTransactionResponse{
		TxHash: s,
	}
}
