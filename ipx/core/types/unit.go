package types

import (
	"fmt"
	"math/big"
	"strings"
)

var (
	weiPerGwei  = new(big.Int).Exp(big.NewInt(10), big.NewInt(9), nil)
	weiPerEther = new(big.Int).Exp(big.NewInt(10), big.NewInt(18), nil)
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

func GweiToWei(gwei *big.Int) string {
	return multiplyUnit(gwei, weiPerGwei)
}

func GweiToEther(gwei *big.Int) string {
	return formatScaledInt(gwei, 9)
}

func EtherToWei(ether *big.Int) string {
	return multiplyUnit(ether, weiPerEther)
}

func EtherToGwei(ether *big.Int) string {
	return multiplyUnit(ether, weiPerGwei)
}

func ConvertUnitDecimal(amount *big.Int, from string, to string) (string, error) {
	switch from {
	case UnitWei:
		switch to {
		case UnitWei:
			return amount.String(), nil
		case UnitGwei:
			return WeiToGwei(amount), nil
		case UnitEther:
			return WeiToEther(amount), nil
		}
	case UnitGwei:
		switch to {
		case UnitWei:
			return GweiToWei(amount), nil
		case UnitGwei:
			return amount.String(), nil
		case UnitEther:
			return GweiToEther(amount), nil
		}
	case UnitEther:
		switch to {
		case UnitWei:
			return EtherToWei(amount), nil
		case UnitGwei:
			return EtherToGwei(amount), nil
		case UnitEther:
			return amount.String(), nil
		}
	}

	return "", fmt.Errorf("unsupported unit conversion: %s -> %s", from, to)
}

func ConvertUnitHex(amount *big.Int, from string, to string) (string, error) {
	switch from {
	case UnitWei:
		switch to {
		case UnitWei:
			return bigIntToHex(amount), nil
		case UnitGwei:
			return divideUnitHex(amount, weiPerGwei, from, to)
		case UnitEther:
			return divideUnitHex(amount, weiPerEther, from, to)
		}
	case UnitGwei:
		switch to {
		case UnitWei:
			return bigIntToHex(new(big.Int).Mul(new(big.Int).Set(amount), weiPerGwei)), nil
		case UnitGwei:
			return bigIntToHex(amount), nil
		case UnitEther:
			return divideUnitHex(amount, weiPerGwei, from, to)
		}
	case UnitEther:
		switch to {
		case UnitWei:
			return bigIntToHex(new(big.Int).Mul(new(big.Int).Set(amount), weiPerEther)), nil
		case UnitGwei:
			return bigIntToHex(new(big.Int).Mul(new(big.Int).Set(amount), weiPerGwei)), nil
		case UnitEther:
			return bigIntToHex(amount), nil
		}
	}

	return "", fmt.Errorf("unsupported unit conversion: %s -> %s", from, to)
}

func divideUnitHex(amount *big.Int, divisor *big.Int, from string, to string) (string, error) {
	quotient, remainder := new(big.Int).QuoRem(new(big.Int).Set(amount), divisor, new(big.Int))
	if remainder.Sign() != 0 {
		return "", fmt.Errorf("unit conversion %s -> %s does not produce an integer hex value", from, to)
	}
	return bigIntToHex(quotient), nil
}

func bigIntToHex(n *big.Int) string {
	if n == nil {
		return "0x0"
	}
	return "0x" + n.Text(16)
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

func multiplyUnit(value *big.Int, multiplier *big.Int) string {
	if value == nil {
		return "0"
	}
	return new(big.Int).Mul(value, multiplier).String()
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
