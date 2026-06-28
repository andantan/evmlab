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
	var err error
	s := &LegacyTransactionState{
		From:  msg.From,
		To:    msg.To,
		Value: msg.Value,
		Data:  msg.Data,
	}

	chainIDHex, err := client.ChainID(ctx)
	if err != nil {
		return nil, fmt.Errorf("get chain id: %w", err)
	}
	s.ChainID, err = util.HexToBigInt(chainIDHex)
	if err != nil {
		return nil, fmt.Errorf("parse chain id: %w", err)
	}

	nonceHex, err := client.GetTransactionCount(ctx, msg.From.String(), "pending")
	if err != nil {
		return nil, fmt.Errorf("get nonce: %w", err)
	}
	s.Nonce, err = util.HexToUint64(nonceHex)
	if err != nil {
		return nil, fmt.Errorf("parse nonce: %w", err)
	}

	gasPriceHex, err := client.GasPrice(ctx)
	if err != nil {
		return nil, fmt.Errorf("get gas price: %w", err)
	}
	s.GasPrice, err = util.HexToBigInt(gasPriceHex)
	if err != nil {
		return nil, fmt.Errorf("parse gas price: %w", err)
	}

	gasEstHex, err := client.EstimateGas(ctx, msg.Map(), "latest")
	if err != nil {
		return nil, fmt.Errorf("estimate gas: %w", err)
	}
	gasEst, err := util.HexToUint64(gasEstHex)
	if err != nil {
		return nil, fmt.Errorf("parse gas estimate: %w", err)
	}
	s.GasLimit = gasEst * 2

	return s, nil
}

func GenerateEIP1559TransactionState(ctx context.Context, client *rpc.Client, msg CallMsg) (*EIP1559TransactionState, error) {
	var err error
	s := &EIP1559TransactionState{
		From:  msg.From,
		To:    msg.To,
		Value: msg.Value,
		Data:  msg.Data,
	}

	chainIDHex, err := client.ChainID(ctx)
	if err != nil {
		return nil, fmt.Errorf("get chain id: %w", err)
	}
	s.ChainID, err = util.HexToBigInt(chainIDHex)
	if err != nil {
		return nil, fmt.Errorf("parse chain id: %w", err)
	}

	nonceHex, err := client.GetTransactionCount(ctx, msg.From.String(), "pending")
	if err != nil {
		return nil, fmt.Errorf("get nonce: %w", err)
	}
	s.Nonce, err = util.HexToUint64(nonceHex)
	if err != nil {
		return nil, fmt.Errorf("parse nonce: %w", err)
	}

	tipCapHex, err := client.MaxPriorityFeePerGas(ctx)
	if err != nil {
		return nil, fmt.Errorf("get tip cap: %w", err)
	}
	s.GasTipCap, err = util.HexToBigInt(tipCapHex)
	if err != nil {
		return nil, fmt.Errorf("parse tip cap: %w", err)
	}

	baseFeeHex, err := client.BaseFeePerGas(ctx)
	if err != nil {
		return nil, fmt.Errorf("get base fee: %w", err)
	}
	baseFee, err := util.HexToBigInt(baseFeeHex)
	if err != nil {
		return nil, fmt.Errorf("parse base fee: %w", err)
	}
	s.GasFeeCap = new(big.Int).Add(new(big.Int).Mul(baseFee, big.NewInt(2)), s.GasTipCap)

	gasEstHex, err := client.EstimateGas(ctx, msg.Map(), "latest")
	if err != nil {
		return nil, fmt.Errorf("estimate gas: %w", err)
	}
	gasEst, err := util.HexToUint64(gasEstHex)
	if err != nil {
		return nil, fmt.Errorf("parse gas estimate: %w", err)
	}
	s.GasLimit = gasEst * 2

	return s, nil
}
