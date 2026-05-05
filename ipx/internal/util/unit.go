package util

import (
	"strings"
)

func IsSupportedEthereumUnit(unit string) bool {
	switch strings.ToLower(strings.TrimSpace(unit)) {
	case "wei", "gwei", "ether":
		return true
	default:
		return false
	}
}
