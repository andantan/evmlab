package util

import (
	"encoding/hex"
	"fmt"
	"math/big"
	"strconv"
	"strings"
)

func HexToUint64(s string) (uint64, error) {
	return strconv.ParseUint(strings.TrimPrefix(s, "0x"), 16, 64)
}

func HexToBigInt(s string) (*big.Int, error) {
	n := new(big.Int)
	if _, ok := n.SetString(strings.TrimPrefix(s, "0x"), 16); !ok {
		return nil, fmt.Errorf("invalid hex integer: %s", s)
	}
	return n, nil
}

// ParseHex decodes a hex string (with or without 0x prefix) into bytes.
func ParseHex(s string) ([]byte, error) {
	s = strings.TrimPrefix(strings.TrimSpace(s), "0x")
	b, err := hex.DecodeString(s)
	if err != nil {
		return nil, fmt.Errorf("invalid hex: %w", err)
	}
	return b, nil
}

// FormatTokenAmount formats a raw uint256 token amount with the given decimal places,
// trimming trailing zeros after the decimal point.
func FormatTokenAmount(amount *big.Int, decimals uint8) string {
	if decimals == 0 {
		return amount.String()
	}
	divisor := new(big.Int).Exp(big.NewInt(10), big.NewInt(int64(decimals)), nil)
	intPart := new(big.Int)
	fracPart := new(big.Int)
	intPart.DivMod(amount, divisor, fracPart)
	if fracPart.Sign() == 0 {
		return intPart.String()
	}
	fracStr := fracPart.String()
	for len(fracStr) < int(decimals) {
		fracStr = "0" + fracStr
	}
	return intPart.String() + "." + strings.TrimRight(fracStr, "0")
}

// IsHexAddress reports whether s is a valid 20-byte hex address (with or without 0x prefix).
func IsHexAddress(s string) bool {
	s = strings.TrimPrefix(s, "0x")
	if len(s) != 40 {
		return false
	}
	_, err := hex.DecodeString(s)
	return err == nil
}
