package types

import (
	"fmt"
	"math/big"
	"strings"
)

const (
	UnitWei   = "wei"
	UnitGwei  = "gwei"
	UnitEther = "ether"
)

func WeiToGwei(wei *big.Int) string {
	return formatScaledInt(wei, 9)
}

func WeiToEther(wei *big.Int) string {
	return formatScaledInt(wei, 18)
}

func ConvertUnitDecimal(amount, from, to string) (string, error) {
	var scale int
	switch from {
	case UnitWei:
		scale = 0
	case UnitGwei:
		scale = 9
	case UnitEther:
		scale = 18
	default:
		return "", fmt.Errorf("invalid unit: %s", from)
	}

	intPart, fracPart, _ := strings.Cut(amount, ".")
	if len(fracPart) > scale {
		return "", fmt.Errorf("amount: too many decimal places for %s", from)
	}

	digits := strings.TrimLeft(intPart+fracPart+strings.Repeat("0", scale-len(fracPart)), "0")
	if digits == "" {
		digits = "0"
	}

	n, ok := new(big.Int).SetString(digits, 10)
	if !ok {
		return "", fmt.Errorf("amount: invalid decimal")
	}

	switch to {
	case UnitWei:
		return n.String(), nil
	case UnitGwei:
		return formatScaledInt(n, 9), nil
	case UnitEther:
		return formatScaledInt(n, 18), nil
	default:
		return "", fmt.Errorf("invalid unit: %s", to)
	}
}

func formatScaledInt(value *big.Int, scale int) string {
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
