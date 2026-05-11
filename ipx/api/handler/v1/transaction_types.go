package v1

import (
	"encoding/hex"
	"errors"
	"math/big"
	"strings"

	"github.com/andantan/evmlab/core/types"
	"github.com/andantan/evmlab/internal/util"
)

type BuildLegacyTransactionRequest struct {
	ChainID  string `json:"chain_id"  example:"20001209"`
	Nonce    uint64 `json:"nonce"     example:"0"`
	GasPrice string `json:"gas_price" example:"20000000000"`
	GasLimit uint64 `json:"gas_limit" example:"21000"`
	To       string `json:"to"        example:"0x8336c196ABb9E7092C879C28D352b39d3f2f3D7A"`
	Value    string `json:"value"     example:"1000000000000000000"`
	Data     string `json:"data"      example:"0x"`

	chainID  *big.Int
	gasPrice *big.Int
	to       *types.Address
	value    *big.Int
	data     []byte
}

func (r *BuildLegacyTransactionRequest) ValidateRequest() error {
	var ok bool

	r.chainID, ok = new(big.Int).SetString(strings.TrimSpace(r.ChainID), 10)
	if !ok {
		return errors.New("chain_id: must be a decimal integer")
	}
	if r.chainID.Sign() <= 0 {
		return errors.New("chain_id: must be positive")
	}

	r.gasPrice, ok = new(big.Int).SetString(strings.TrimSpace(r.GasPrice), 10)
	if !ok {
		return errors.New("gas_price: must be a decimal integer")
	}
	if r.gasPrice.Sign() <= 0 {
		return errors.New("gas_price: must be positive")
	}

	if r.GasLimit == 0 {
		return errors.New("gas_limit: must be greater than zero")
	}

	var err error
	r.to, err = types.NewAddressFromHex(r.To)
	if err != nil {
		return errors.New("to: " + err.Error())
	}

	if v := strings.TrimSpace(r.Value); v == "" {
		r.value = new(big.Int)
	} else {
		r.value, ok = new(big.Int).SetString(v, 10)
		if !ok {
			return errors.New("value: must be a decimal integer")
		}
		if r.value.Sign() < 0 {
			return errors.New("value: must be non-negative")
		}
	}

	if d := strings.TrimSpace(r.Data); d == "" || d == "0x" {
		return nil
	}

	if r.data, err = util.ParseHex(r.Data); err != nil {
		return errors.New("data: " + err.Error())
	}

	return nil
}

func (r *BuildLegacyTransactionRequest) ToLegacyTx() *types.LegacyTx {
	return &types.LegacyTx{
		ChainID:  r.chainID,
		Nonce:    r.Nonce,
		GasPrice: r.gasPrice,
		GasLimit: r.GasLimit,
		To:       &r.to.Addr,
		Value:    r.value,
		Data:     r.data,
	}
}

type BuildLegacyTransactionResponse struct {
	UnsignedRLP string `json:"unsigned_rlp"`
	SigningHash string `json:"signing_hash"`
}

func NewBuildLegacyNativeTransferResponse(b []byte, h *types.Hash) *BuildLegacyTransactionResponse {
	return &BuildLegacyTransactionResponse{
		UnsignedRLP: "0x" + hex.EncodeToString(b),
		SigningHash: h.String(),
	}
}

type BuildEIP1559TransactionRequest struct {
	ChainID              string `json:"chain_id"                 example:"20001209"`
	Nonce                uint64 `json:"nonce"                    example:"0"`
	MaxPriorityFeePerGas string `json:"max_priority_fee_per_gas" example:"1500000000"`
	MaxFeePerGas         string `json:"max_fee_per_gas"          example:"3000000000"`
	GasLimit             uint64 `json:"gas_limit"                example:"21000"`
	To                   string `json:"to"                       example:"0x8336c196ABb9E7092C879C28D352b39d3f2f3D7A"`
	Value                string `json:"value"                    example:"1000000000000000000"`
	Data                 string `json:"data"                     example:"0x"`

	chainID   *big.Int
	gasTipCap *big.Int
	gasFeeCap *big.Int
	to        *types.Address
	value     *big.Int
	data      []byte
}

func (r *BuildEIP1559TransactionRequest) ValidateRequest() error {
	var ok bool

	r.chainID, ok = new(big.Int).SetString(strings.TrimSpace(r.ChainID), 10)
	if !ok {
		return errors.New("chain_id: must be a decimal integer")
	}
	if r.chainID.Sign() <= 0 {
		return errors.New("chain_id: must be positive")
	}

	r.gasTipCap, ok = new(big.Int).SetString(strings.TrimSpace(r.MaxPriorityFeePerGas), 10)
	if !ok {
		return errors.New("max_priority_fee_per_gas: must be a decimal integer")
	}
	if r.gasTipCap.Sign() <= 0 {
		return errors.New("max_priority_fee_per_gas: must be positive")
	}

	r.gasFeeCap, ok = new(big.Int).SetString(strings.TrimSpace(r.MaxFeePerGas), 10)
	if !ok {
		return errors.New("max_fee_per_gas: must be a decimal integer")
	}
	if r.gasFeeCap.Sign() <= 0 {
		return errors.New("max_fee_per_gas: must be positive")
	}

	if r.GasLimit == 0 {
		return errors.New("gas_limit: must be greater than zero")
	}

	var err error
	r.to, err = types.NewAddressFromHex(r.To)
	if err != nil {
		return errors.New("to: " + err.Error())
	}

	if v := strings.TrimSpace(r.Value); v == "" {
		r.value = new(big.Int)
	} else {
		r.value, ok = new(big.Int).SetString(v, 10)
		if !ok {
			return errors.New("value: must be a decimal integer")
		}
		if r.value.Sign() < 0 {
			return errors.New("value: must be non-negative")
		}
	}

	if d := strings.TrimSpace(r.Data); d == "" || d == "0x" {
		return nil
	}

	if r.data, err = util.ParseHex(r.Data); err != nil {
		return errors.New("data: " + err.Error())
	}

	return nil
}

func (r *BuildEIP1559TransactionRequest) ToDynamicFeeTx() *types.DynamicFeeTx {
	return &types.DynamicFeeTx{
		ChainID:   r.chainID,
		Nonce:     r.Nonce,
		GasTipCap: r.gasTipCap,
		GasFeeCap: r.gasFeeCap,
		GasLimit:  r.GasLimit,
		To:        &r.to.Addr,
		Value:     r.value,
		Data:      r.data,
	}
}

type BuildEIP1559TransactionResponse struct {
	UnsignedRLP string `json:"unsigned_rlp"`
	SigningHash string `json:"signing_hash"`
}

func NewBuildEIP1559TransactionResponse(b []byte, h *types.Hash) *BuildEIP1559TransactionResponse {
	return &BuildEIP1559TransactionResponse{
		UnsignedRLP: "0x" + hex.EncodeToString(b),
		SigningHash: h.String(),
	}
}
