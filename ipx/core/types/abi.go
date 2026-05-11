package types

import (
	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/common"
)

type ABIArgument = abi.Argument
type ABIArguments = abi.Arguments

// FormatABIAddress returns the checksummed hex string for an ABI-decoded address value.
func FormatABIAddress(v any) string {
	return v.(common.Address).Hex()
}

func mustABIType(t string) abi.Type {
	typ, err := abi.NewType(t, "", nil)
	if err != nil {
		panic("abi: invalid type " + t + ": " + err.Error())
	}
	return typ
}

var (
	// Dynamic types
	ABIBool    = mustABIType("bool")
	ABIAddress = mustABIType("address")
	ABIString  = mustABIType("string")
	ABIBytes   = mustABIType("bytes")

	// Unsigned integers
	ABIUint8   = mustABIType("uint8")
	ABIUint16  = mustABIType("uint16")
	ABIUint32  = mustABIType("uint32")
	ABIUint64  = mustABIType("uint64")
	ABIUint128 = mustABIType("uint128")
	ABIUint256 = mustABIType("uint256")

	// Signed integers
	ABIInt8   = mustABIType("int8")
	ABIInt16  = mustABIType("int16")
	ABIInt32  = mustABIType("int32")
	ABIInt64  = mustABIType("int64")
	ABIInt128 = mustABIType("int128")
	ABIInt256 = mustABIType("int256")

	// Fixed-size bytes
	ABIBytes1  = mustABIType("bytes1")
	ABIBytes2  = mustABIType("bytes2")
	ABIBytes3  = mustABIType("bytes3")
	ABIBytes4  = mustABIType("bytes4")
	ABIBytes5  = mustABIType("bytes5")
	ABIBytes6  = mustABIType("bytes6")
	ABIBytes7  = mustABIType("bytes7")
	ABIBytes8  = mustABIType("bytes8")
	ABIBytes9  = mustABIType("bytes9")
	ABIBytes10 = mustABIType("bytes10")
	ABIBytes11 = mustABIType("bytes11")
	ABIBytes12 = mustABIType("bytes12")
	ABIBytes13 = mustABIType("bytes13")
	ABIBytes14 = mustABIType("bytes14")
	ABIBytes15 = mustABIType("bytes15")
	ABIBytes16 = mustABIType("bytes16")
	ABIBytes17 = mustABIType("bytes17")
	ABIBytes18 = mustABIType("bytes18")
	ABIBytes19 = mustABIType("bytes19")
	ABIBytes20 = mustABIType("bytes20")
	ABIBytes21 = mustABIType("bytes21")
	ABIBytes22 = mustABIType("bytes22")
	ABIBytes23 = mustABIType("bytes23")
	ABIBytes24 = mustABIType("bytes24")
	ABIBytes25 = mustABIType("bytes25")
	ABIBytes26 = mustABIType("bytes26")
	ABIBytes27 = mustABIType("bytes27")
	ABIBytes28 = mustABIType("bytes28")
	ABIBytes29 = mustABIType("bytes29")
	ABIBytes30 = mustABIType("bytes30")
	ABIBytes31 = mustABIType("bytes31")
	ABIBytes32 = mustABIType("bytes32")

	// Slices
	ABIBoolSlice    = mustABIType("bool[]")
	ABIAddressSlice = mustABIType("address[]")
	ABIStringSlice  = mustABIType("string[]")
	ABIBytesSlice   = mustABIType("bytes[]")
	ABIUint8Slice   = mustABIType("uint8[]")
	ABIUint16Slice  = mustABIType("uint16[]")
	ABIUint32Slice  = mustABIType("uint32[]")
	ABIUint64Slice  = mustABIType("uint64[]")
	ABIUint128Slice = mustABIType("uint128[]")
	ABIUint256Slice = mustABIType("uint256[]")
	ABIInt8Slice    = mustABIType("int8[]")
	ABIInt16Slice   = mustABIType("int16[]")
	ABIInt32Slice   = mustABIType("int32[]")
	ABIInt64Slice   = mustABIType("int64[]")
	ABIInt128Slice  = mustABIType("int128[]")
	ABIInt256Slice  = mustABIType("int256[]")
	ABIBytes32Slice = mustABIType("bytes32[]")
)
