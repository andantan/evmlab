package types

import "fmt"

// Selector is a 4-byte ABI function/error selector.
type Selector [4]byte

func (s Selector) Bytes() []byte {
	return s[:]
}

func (s Selector) String() string {
	return fmt.Sprintf("0x%x", s[:])
}

var (
	BalanceOfSelector    = Selector{0x70, 0xa0, 0x82, 0x31} // balanceOf(address)
	ApproveSelector      = Selector{0x09, 0x5e, 0xa7, 0xb3} // approve(address,uint256)
	TransferSelector     = Selector{0xa9, 0x05, 0x9c, 0xbb} // transfer(address,uint256)
	AllowanceSelector    = Selector{0xdd, 0x62, 0xed, 0x3e} // allowance(address,address)
	TransferFromSelector = Selector{0x23, 0xb8, 0x72, 0xdd} // transferFrom(address,address,uint256)
	EIP712DomainSelector = Selector{0x84, 0xb0, 0x19, 0x6e} // eip712Domain()
	NoncesSelector       = Selector{0x7e, 0xce, 0xbe, 0x00} // nonces(address)
	NameSelector         = Selector{0x06, 0xfd, 0xde, 0x03} // name()
	VersionSelector      = Selector{0x54, 0xfd, 0x4d, 0x50} // version()
	SymbolSelector       = Selector{0x95, 0xd8, 0x9b, 0x41} // symbol()
	DecimalsSelector     = Selector{0x31, 0x3c, 0xe5, 0x67} // decimals()
	TotalSupplySelector  = Selector{0x18, 0x16, 0x0d, 0xdd} // totalSupply()
)

// Sanctum selectors
var (
	RegisterSanctumSelector          = Selector{0x1a, 0xa3, 0xa0, 0x08} // register()
	RegisterForSanctumSelector       = Selector{0xe3, 0x0e, 0x38, 0x34} // registerFor(address)
	ApproveRegisterSanctumSelector   = Selector{0xdc, 0xbd, 0x5f, 0x8d} // approveRegister(address)
	DeregisterSanctumSelector        = Selector{0xaf, 0xf5, 0xed, 0xb1} // deregister()
	DeregisterForSanctumSelector     = Selector{0x1f, 0x2a, 0xe7, 0xe8} // deregisterFor(address)
	GetAccountsSanctumSelector       = Selector{0x8a, 0x48, 0xac, 0x03} // getAccounts()
	AccountCountSanctumSelector      = Selector{0xe4, 0xaf, 0x29, 0xfc} // accountCount()
	GetAccountInfoSanctumSelector    = Selector{0x7b, 0x51, 0x0f, 0xe8} // getAccountInfo(address)
	DepositNativeSanctumSelector     = Selector{0xdb, 0x6b, 0x52, 0x46} // depositNative()
	RequestNativeSanctumSelector     = Selector{0x46, 0x11, 0xed, 0x6c} // requestNative(uint256)
	ApproveNativeSanctumSelector     = Selector{0x23, 0xd5, 0x78, 0x86} // approveNative(address,uint256)
	ApproveNativeAllSanctumSelector  = Selector{0x7e, 0x99, 0xcd, 0xc0} // approveNativeAll(address)
	WithdrawNativeSanctumSelector    = Selector{0x84, 0x27, 0x6d, 0x81} // withdrawNative(uint256)
	WithdrawNativeAllSanctumSelector = Selector{0x01, 0x5e, 0xd3, 0xc3} // withdrawNativeAll()
	NativeBalanceSanctumSelector     = Selector{0x86, 0x5f, 0xc5, 0x01} // nativeBalance()
	NativeAvailableSanctumSelector   = Selector{0xcb, 0xb4, 0x19, 0x17} // nativeAvailable()
	NativeAllocationSanctumSelector  = Selector{0xa0, 0x16, 0x42, 0x72} // nativeAllocation(address)
	NativePendingSanctumSelector     = Selector{0x4b, 0x0f, 0xd7, 0x89} // nativePending(address)
)

// Sanctum error selectors
var (
	UnauthorizedSanctumErrorSelector             = Selector{0x8e, 0x4a, 0x23, 0xd6} // Unauthorized(address)
	AlreadyRegisteredAccountSanctumErrorSelector = Selector{0xba, 0x61, 0x4c, 0x22} // AlreadyRegisteredAccount(address)
	AlreadyApprovedAccountSanctumErrorSelector   = Selector{0xe9, 0x7d, 0xad, 0x3f} // AlreadyApprovedAccount(address)
	NotRegisteredAccountSanctumErrorSelector     = Selector{0xbf, 0x58, 0xe0, 0x76} // NotRegisteredAccount(address)
	InvalidAccountSanctumErrorSelector           = Selector{0x4b, 0x57, 0x9b, 0x22} // InvalidAccount(address)
	InvalidRoleSanctumErrorSelector              = Selector{0xf1, 0x78, 0x64, 0xde} // InvalidRole(uint8)
	CannotRemoveMasterSanctumErrorSelector       = Selector{0x45, 0x37, 0x2b, 0xc9} // CannotRemoveMaster()
	InsufficientBalanceSanctumErrorSelector      = Selector{0xdb, 0x42, 0x14, 0x4d} // InsufficientBalance(address,uint256,uint256)
	InsufficientAllocationSanctumErrorSelector   = Selector{0x0c, 0x17, 0xa2, 0xf3} // InsufficientAllocation(address,uint256,uint256)
	InsufficientPendingSanctumErrorSelector      = Selector{0xb8, 0x91, 0x24, 0xd7} // InsufficientPending(address,uint256,uint256)
	ZeroAmountSanctumErrorSelector               = Selector{0x1f, 0x2a, 0x20, 0x05} // ZeroAmount()
	TransferFailedSanctumErrorSelector           = Selector{0xbf, 0x18, 0x2b, 0xe8} // TransferFailed(address,address,uint256)
	NoPendedBalanceSanctumErrorSelector          = Selector{0x3c, 0x86, 0xb5, 0xd6} // NoPendedBalance(address,address)
	NoAllocatedBalanceSanctumErrorSelector       = Selector{0x50, 0x20, 0x4e, 0x15} // NoAllocatedBalance(address,address)
)

// SanctumErrorSignatures maps Sanctum error selectors to their ABI signatures.
// Used to decode revert data into human-readable messages with actual parameter values.
var SanctumErrorSignatures = map[Selector]string{
	UnauthorizedSanctumErrorSelector:             "Unauthorized(address caller)",
	AlreadyRegisteredAccountSanctumErrorSelector: "AlreadyRegisteredAccount(address account)",
	AlreadyApprovedAccountSanctumErrorSelector:   "AlreadyApprovedAccount(address account)",
	NotRegisteredAccountSanctumErrorSelector:     "NotRegisteredAccount(address account)",
	InvalidAccountSanctumErrorSelector:           "InvalidAccount(address account)",
	InvalidRoleSanctumErrorSelector:              "InvalidRole(uint8 role)",
	CannotRemoveMasterSanctumErrorSelector:       "CannotRemoveMaster()",
	InsufficientBalanceSanctumErrorSelector:      "InsufficientBalance(address asset,uint256 requested,uint256 available)",
	InsufficientAllocationSanctumErrorSelector:   "InsufficientAllocation(address asset,uint256 requested,uint256 available)",
	InsufficientPendingSanctumErrorSelector:      "InsufficientPending(address asset,uint256 requested,uint256 available)",
	ZeroAmountSanctumErrorSelector:               "ZeroAmount()",
	TransferFailedSanctumErrorSelector:           "TransferFailed(address asset,address to,uint256 amount)",
	NoPendedBalanceSanctumErrorSelector:          "NoPendedBalance(address asset,address user)",
	NoAllocatedBalanceSanctumErrorSelector:       "NoAllocatedBalance(address asset,address user)",
}
