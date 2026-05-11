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

// IsHexAddress reports whether s is a valid 20-byte hex address (with or without 0x prefix).
func IsHexAddress(s string) bool {
	s = strings.TrimPrefix(s, "0x")
	if len(s) != 40 {
		return false
	}
	_, err := hex.DecodeString(s)
	return err == nil
}
