package core

import (
	"math/big"

	"github.com/andantan/evmlab/core/types"
)

// BalanceOfCalldata builds calldata for balanceOf(address).
func BalanceOfCalldata(account *types.Address) []byte {
	data := make([]byte, 36)
	copy(data[:4], types.BalanceOfSelector.Bytes())
	copy(data[16:36], account.Bytes())
	return data
}

// ApproveCalldata builds calldata for approve(address,uint256).
func ApproveCalldata(spender *types.Address, amount *big.Int) []byte {
	data := make([]byte, 68)
	copy(data[:4], types.ApproveSelector.Bytes())
	copy(data[16:36], spender.Bytes())
	amount.FillBytes(data[36:68])
	return data
}

// TransferCalldata builds calldata for transfer(address,uint256).
func TransferCalldata(to *types.Address, amount *big.Int) []byte {
	data := make([]byte, 68)
	copy(data[:4], types.TransferSelector.Bytes())
	copy(data[16:36], to.Bytes())
	amount.FillBytes(data[36:68])
	return data
}

// TransferFromCalldata builds calldata for transferFrom(address,address,uint256).
func TransferFromCalldata(from, to *types.Address, amount *big.Int) []byte {
	data := make([]byte, 100)
	copy(data[:4], types.TransferFromSelector.Bytes())
	copy(data[16:36], from.Bytes())
	copy(data[48:68], to.Bytes())
	amount.FillBytes(data[68:100])
	return data
}

// Multicall3Aggregator3CallData builds calldata for aggregate3((address,bool,bytes)[]).
func Multicall3Aggregator3CallData(calls types.Aggregate3s) []byte {
	n := len(calls)

	elemSizes := make([]int, n)
	totalElemBytes := 0
	for i, c := range calls {
		elemSizes[i] = 96 + 32 + (len(c.CallData)+31)&^31
		totalElemBytes += elemSizes[i]
	}

	b := make([]byte, 4+32+32+n*32+totalElemBytes)

	copy(b[0:4], types.MultiCall3Aggregate3Selector.Bytes())

	b[35] = 0x20 // offset to array = 32

	v := uint64(n)
	b[60], b[61], b[62], b[63] = byte(v>>56), byte(v>>48), byte(v>>40), byte(v>>32)
	b[64], b[65], b[66], b[67] = byte(v>>24), byte(v>>16), byte(v>>8), byte(v)

	off := uint64(n * 32)
	for i := range n {
		base := 68 + i*32
		o := off
		b[base+24], b[base+25], b[base+26], b[base+27] = byte(o>>56), byte(o>>48), byte(o>>40), byte(o>>32)
		b[base+28], b[base+29], b[base+30], b[base+31] = byte(o>>24), byte(o>>16), byte(o>>8), byte(o)
		off += uint64(elemSizes[i])
	}

	pos := 68 + n*32
	for _, c := range calls {
		copy(b[pos+12:pos+32], c.Target.Bytes())
		pos += 32

		if c.AllowFail {
			b[pos+31] = 1
		}
		pos += 32

		b[pos+31] = 96 // offset to bytes within tuple
		pos += 32

		l := uint64(len(c.CallData))
		b[pos+24], b[pos+25], b[pos+26], b[pos+27] = byte(l>>56), byte(l>>48), byte(l>>40), byte(l>>32)
		b[pos+28], b[pos+29], b[pos+30], b[pos+31] = byte(l>>24), byte(l>>16), byte(l>>8), byte(l)
		pos += 32

		copy(b[pos:pos+len(c.CallData)], c.CallData)
		pos += (len(c.CallData) + 31) &^ 31
	}

	return b
}

// EIP712DomainCalldata builds calldata for eip712Domain().
func EIP712DomainCalldata() []byte {
	data := make([]byte, 4)
	copy(data, types.EIP712DomainSelector.Bytes())
	return data
}

// NameCalldata builds calldata for name().
func NameCalldata() []byte {
	data := make([]byte, 4)
	copy(data, types.NameSelector.Bytes())
	return data
}

// VersionCalldata builds calldata for version().
func VersionCalldata() []byte {
	data := make([]byte, 4)
	copy(data, types.VersionSelector.Bytes())
	return data
}

// SymbolCalldata builds calldata for symbol().
func SymbolCalldata() []byte {
	data := make([]byte, 4)
	copy(data, types.SymbolSelector.Bytes())
	return data
}

// DecimalsCalldata builds calldata for decimals().
func DecimalsCalldata() []byte {
	data := make([]byte, 4)
	copy(data, types.DecimalsSelector.Bytes())
	return data
}

// TotalSupplyCalldata builds calldata for totalSupply().
func TotalSupplyCalldata() []byte {
	data := make([]byte, 4)
	copy(data, types.TotalSupplySelector.Bytes())
	return data
}

// NoncesCalldata builds calldata for nonces(address).
func NoncesCalldata(owner *types.Address) []byte {
	data := make([]byte, 36)
	copy(data[:4], types.NoncesSelector.Bytes())
	copy(data[16:36], owner.Bytes())
	return data
}

// AllowanceCalldata builds calldata for allowance(address,address).
func AllowanceCalldata(owner, spender *types.Address) []byte {
	data := make([]byte, 68)
	copy(data[:4], types.AllowanceSelector.Bytes())
	copy(data[16:36], owner.Bytes())
	copy(data[48:68], spender.Bytes())
	return data
}

// RegisterSanctumCalldata builds calldata for register().
func RegisterSanctumCalldata() []byte {
	data := make([]byte, 4)
	copy(data, types.RegisterSanctumSelector.Bytes())
	return data
}

// RegisterForSanctumCalldata builds calldata for registerFor(address).
func RegisterForSanctumCalldata(target *types.Address) []byte {
	data := make([]byte, 36)
	copy(data[:4], types.RegisterForSanctumSelector.Bytes())
	copy(data[16:36], target.Bytes())
	return data
}

// ApproveRegisterSanctumCalldata builds calldata for approveRegister(address).
func ApproveRegisterSanctumCalldata(target *types.Address) []byte {
	data := make([]byte, 36)
	copy(data[:4], types.ApproveRegisterSanctumSelector.Bytes())
	copy(data[16:36], target.Bytes())
	return data
}

// DeregisterSanctumCalldata builds calldata for deregister().
func DeregisterSanctumCalldata() []byte {
	data := make([]byte, 4)
	copy(data, types.DeregisterSanctumSelector.Bytes())
	return data
}

// DeregisterForSanctumCalldata builds calldata for deregisterFor(address).
func DeregisterForSanctumCalldata(target *types.Address) []byte {
	data := make([]byte, 36)
	copy(data[:4], types.DeregisterForSanctumSelector.Bytes())
	copy(data[16:36], target.Bytes())
	return data
}

// GetAccountsSanctumCalldata builds calldata for getAccounts().
func GetAccountsSanctumCalldata() []byte {
	data := make([]byte, 4)
	copy(data, types.GetAccountsSanctumSelector.Bytes())
	return data
}

// AccountCountSanctumCalldata builds calldata for accountCount().
func AccountCountSanctumCalldata() []byte {
	data := make([]byte, 4)
	copy(data, types.AccountCountSanctumSelector.Bytes())
	return data
}

// GetAccountInfoSanctumCalldata builds calldata for getAccountInfo(address).
func GetAccountInfoSanctumCalldata(account *types.Address) []byte {
	data := make([]byte, 36)
	copy(data[:4], types.GetAccountInfoSanctumSelector.Bytes())
	copy(data[16:36], account.Bytes())
	return data
}

// DepositNativeSanctumCalldata builds calldata for depositNative().
func DepositNativeSanctumCalldata() []byte {
	data := make([]byte, 4)
	copy(data, types.DepositNativeSanctumSelector.Bytes())
	return data
}

// RequestNativeSanctumCalldata builds calldata for requestNative(uint256).
func RequestNativeSanctumCalldata(amount *big.Int) []byte {
	data := make([]byte, 36)
	copy(data[:4], types.RequestNativeSanctumSelector.Bytes())
	amount.FillBytes(data[4:36])
	return data
}

// ApproveNativeSanctumCalldata builds calldata for approveNative(address,uint256).
func ApproveNativeSanctumCalldata(user *types.Address, amount *big.Int) []byte {
	data := make([]byte, 68)
	copy(data[:4], types.ApproveNativeSanctumSelector.Bytes())
	copy(data[16:36], user.Bytes())
	amount.FillBytes(data[36:68])
	return data
}

// ApproveNativeAllSanctumCalldata builds calldata for approveNativeAll(address).
func ApproveNativeAllSanctumCalldata(user *types.Address) []byte {
	data := make([]byte, 36)
	copy(data[:4], types.ApproveNativeAllSanctumSelector.Bytes())
	copy(data[16:36], user.Bytes())
	return data
}

// WithdrawNativeSanctumCalldata builds calldata for withdrawNative(uint256).
func WithdrawNativeSanctumCalldata(amount *big.Int) []byte {
	data := make([]byte, 36)
	copy(data[:4], types.WithdrawNativeSanctumSelector.Bytes())
	amount.FillBytes(data[4:36])
	return data
}

// WithdrawNativeAllSanctumCalldata builds calldata for withdrawNativeAll().
func WithdrawNativeAllSanctumCalldata() []byte {
	data := make([]byte, 4)
	copy(data, types.WithdrawNativeAllSanctumSelector.Bytes())
	return data
}

// NativeBalanceSanctumCalldata builds calldata for nativeBalance().
func NativeBalanceSanctumCalldata() []byte {
	data := make([]byte, 4)
	copy(data, types.NativeBalanceSanctumSelector.Bytes())
	return data
}

// NativeAvailableSanctumCalldata builds calldata for nativeAvailable().
func NativeAvailableSanctumCalldata() []byte {
	data := make([]byte, 4)
	copy(data, types.NativeAvailableSanctumSelector.Bytes())
	return data
}

// NativeAllocationSanctumCalldata builds calldata for nativeAllocation(address).
func NativeAllocationSanctumCalldata(user *types.Address) []byte {
	data := make([]byte, 36)
	copy(data[:4], types.NativeAllocationSanctumSelector.Bytes())
	copy(data[16:36], user.Bytes())
	return data
}

// NativePendingSanctumCalldata builds calldata for nativePending(address).
func NativePendingSanctumCalldata(user *types.Address) []byte {
	data := make([]byte, 36)
	copy(data[:4], types.NativePendingSanctumSelector.Bytes())
	copy(data[16:36], user.Bytes())
	return data
}
