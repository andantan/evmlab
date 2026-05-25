package sanctum

import (
	"fmt"
	"math/big"

	"github.com/andantan/evmlab/core/types"
)

type SanctumRole uint8

const (
	SanctumRoleNone SanctumRole = iota
	SanctumRoleMaster
	SanctumRoleMember
	SanctumRolePending
)

func (r SanctumRole) String() string {
	switch r {
	case SanctumRoleNone:
		return "None"
	case SanctumRoleMaster:
		return "Master"
	case SanctumRoleMember:
		return "Member"
	case SanctumRolePending:
		return "Pending"
	default:
		return fmt.Sprintf("Unknown(%d)", r)
	}
}

var accountInfoTypes = [3]string{"address", "uint8", "uint256"}

type AccountInfo struct {
	Addr            *types.Address
	Role            SanctumRole
	RegisteredBlock *big.Int
}
