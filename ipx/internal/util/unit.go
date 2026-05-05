package util

import (
	"math/big"
	"strings"
)

func MultiplyUnit(value *big.Int, multiplier *big.Int) string {
	if value == nil {
		return "0"
	}
	return new(big.Int).Mul(value, multiplier).String()
}

func FormatScaledInt(value *big.Int, scale int) string {
	if value == nil {
		return "0"
	}

	sign := ""
	abs := new(big.Int).Set(value)
	if abs.Sign() < 0 {
		sign = "-"
		abs.Neg(abs)
	}

	digits := abs.String()
	if scale == 0 {
		return sign + digits
	}

	if len(digits) <= scale {
		return trimTrailingFractionZeros(sign + "0." + strings.Repeat("0", scale-len(digits)) + digits)
	}

	split := len(digits) - scale
	return trimTrailingFractionZeros(sign + digits[:split] + "." + digits[split:])
}

// IsSupportedEthereumUnit reports whether the given unit is one of the
// supported Ethereum denominations used by the conversion helpers.
func IsSupportedEthereumUnit(unit string) bool {
	switch strings.ToLower(strings.TrimSpace(unit)) {
	case "wei", "gwei", "ether":
		return true
	default:
		return false
	}
}

func trimTrailingFractionZeros(s string) string {
	if !strings.Contains(s, ".") {
		return s
	}
	s = strings.TrimRight(s, "0")
	s = strings.TrimRight(s, ".")
	if s == "" || s == "-" {
		return "0"
	}
	return s
}
