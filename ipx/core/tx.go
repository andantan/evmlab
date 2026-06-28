package core

import (
	"context"
	"encoding/hex"
	"fmt"
	"math/big"

	"github.com/andantan/evmlab/core/types"
	"github.com/andantan/evmlab/internal/rpc"
	"github.com/andantan/evmlab/internal/util"
)

type CallMsg struct {
	From  *types.Address
	To    *types.Address
	Value *big.Int
	Data  []byte
}

func NewCallMsg(from, to *types.Address, value *big.Int, data []byte) CallMsg {
	return CallMsg{From: from, To: to, Value: value, Data: data}
}

func (m CallMsg) Map() map[string]any {
	p := map[string]any{
		"from":  m.From.String(),
		"to":    m.To.String(),
		"value": "0x" + m.Value.Text(16),
	}
	if len(m.Data) > 0 {
		p["data"] = "0x" + hex.EncodeToString(m.Data)
	}
	return p
}

type LegacyTransactionState struct {
	ChainID  *big.Int
	Nonce    uint64
	GasPrice *big.Int
	GasLimit uint64
	From     *types.Address
	To       *types.Address
	Value    *big.Int
	Data     []byte
}

func (s *LegacyTransactionState) ToTx() *types.LegacyTx {
	return &types.LegacyTx{
		ChainID:  s.ChainID,
		Nonce:    s.Nonce,
		GasPrice: s.GasPrice,
		GasLimit: s.GasLimit,
		To:       &s.To.Addr,
		Value:    s.Value,
		Data:     s.Data,
	}
}

type EIP1559TransactionState struct {
	ChainID   *big.Int
	Nonce     uint64
	GasTipCap *big.Int
	GasFeeCap *big.Int
	GasLimit  uint64
	From      *types.Address
	To        *types.Address
	Value     *big.Int
	Data      []byte
}

func (s *EIP1559TransactionState) ToTx() *types.DynamicFeeTx {
	return &types.DynamicFeeTx{
		ChainID:   s.ChainID,
		Nonce:     s.Nonce,
		GasTipCap: s.GasTipCap,
		GasFeeCap: s.GasFeeCap,
		GasLimit:  s.GasLimit,
		To:        &s.To.Addr,
		Value:     s.Value,
		Data:      s.Data,
	}
}

func GenerateLegacyTransactionState(ctx context.Context, client *rpc.Client, msg CallMsg) (*LegacyTransactionState, error) {
	s := &LegacyTransactionState{
		From:  msg.From,
		To:    msg.To,
		Value: msg.Value,
		Data:  msg.Data,
	}

	var chainIDHex, nonceHex, gasPriceHex, gasEstHex string
	if err := client.Batch(ctx, new(rpc.Elems).
		With(rpc.ETHChainID(&chainIDHex)).
		With(rpc.ETHGetTransactionCount(msg.From.String(), "pending", &nonceHex)).
		With(rpc.ETHGasPrice(&gasPriceHex)).
		With(rpc.ETHEstimateGas(msg.Map(), "latest", &gasEstHex)),
	); err != nil {
		return nil, err
	}

	var err error
	s.ChainID, err = util.HexToBigInt(chainIDHex)
	if err != nil {
		return nil, fmt.Errorf("parse chain id: %w", err)
	}
	s.Nonce, err = util.HexToUint64(nonceHex)
	if err != nil {
		return nil, fmt.Errorf("parse nonce: %w", err)
	}
	s.GasPrice, err = util.HexToBigInt(gasPriceHex)
	if err != nil {
		return nil, fmt.Errorf("parse gas price: %w", err)
	}
	s.GasPrice = new(big.Int).Div(new(big.Int).Mul(s.GasPrice, big.NewInt(11)), big.NewInt(10))
	gasEst, err := util.HexToUint64(gasEstHex)
	if err != nil {
		return nil, fmt.Errorf("parse gas estimate: %w", err)
	}
	s.GasLimit = gasEst * 2

	return s, nil
}

func GenerateEIP1559TransactionState(ctx context.Context, client *rpc.Client, msg CallMsg) (*EIP1559TransactionState, error) {
	s := &EIP1559TransactionState{
		From:  msg.From,
		To:    msg.To,
		Value: msg.Value,
		Data:  msg.Data,
	}

	var chainIDHex, nonceHex, tipCapHex, gasEstHex string
	var block map[string]any
	if err := client.Batch(ctx, new(rpc.Elems).
		With(rpc.ETHChainID(&chainIDHex)).
		With(rpc.ETHGetTransactionCount(msg.From.String(), "pending", &nonceHex)).
		With(rpc.ETHMaxPriorityFeePerGas(&tipCapHex)).
		With(rpc.ETHGetBlockByNumber("latest", false, &block)).
		With(rpc.ETHEstimateGas(msg.Map(), "latest", &gasEstHex)),
	); err != nil {
		return nil, err
	}

	var err error
	s.ChainID, err = util.HexToBigInt(chainIDHex)
	if err != nil {
		return nil, fmt.Errorf("parse chain id: %w", err)
	}
	s.Nonce, err = util.HexToUint64(nonceHex)
	if err != nil {
		return nil, fmt.Errorf("parse nonce: %w", err)
	}
	s.GasTipCap, err = util.HexToBigInt(tipCapHex)
	if err != nil {
		return nil, fmt.Errorf("parse tip cap: %w", err)
	}

	baseFeeHex, ok := block["baseFeePerGas"].(string)
	if !ok || baseFeeHex == "" {
		return nil, fmt.Errorf("baseFeePerGas not found in latest block")
	}
	baseFee, err := util.HexToBigInt(baseFeeHex)
	if err != nil {
		return nil, fmt.Errorf("parse base fee: %w", err)
	}
	s.GasFeeCap = new(big.Int).Add(new(big.Int).Mul(baseFee, big.NewInt(2)), s.GasTipCap)

	gasEst, err := util.HexToUint64(gasEstHex)
	if err != nil {
		return nil, fmt.Errorf("parse gas estimate: %w", err)
	}
	s.GasLimit = gasEst * 2

	return s, nil
}
