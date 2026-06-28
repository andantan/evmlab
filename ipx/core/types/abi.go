package types

import (
	"fmt"
	"math/big"
	"strings"

	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
)

type ABIArgument = abi.Argument
type ABIArguments = abi.Arguments

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

func DecodeAddress(data []byte) (common.Address, error) {
	values, err := ABIArguments{{Type: ABIAddress}}.Unpack(data)
	if err != nil {
		return common.Address{}, err
	}
	return values[0].(common.Address), nil
}

func DecodeBool(data []byte) (bool, error) {
	values, err := ABIArguments{{Type: ABIBool}}.Unpack(data)
	if err != nil {
		return false, err
	}
	return values[0].(bool), nil
}

func DecodeString(data []byte) (string, error) {
	values, err := ABIArguments{{Type: ABIString}}.Unpack(data)
	if err != nil {
		return "", err
	}
	return values[0].(string), nil
}

func DecodeBytes(data []byte) ([]byte, error) {
	values, err := ABIArguments{{Type: ABIBytes}}.Unpack(data)
	if err != nil {
		return nil, err
	}
	return values[0].([]byte), nil
}

func DecodeUint8(data []byte) (uint8, error) {
	values, err := ABIArguments{{Type: ABIUint8}}.Unpack(data)
	if err != nil {
		return 0, err
	}
	return values[0].(uint8), nil
}

func DecodeUint16(data []byte) (uint16, error) {
	values, err := ABIArguments{{Type: ABIUint16}}.Unpack(data)
	if err != nil {
		return 0, err
	}
	return values[0].(uint16), nil
}

func DecodeUint32(data []byte) (uint32, error) {
	values, err := ABIArguments{{Type: ABIUint32}}.Unpack(data)
	if err != nil {
		return 0, err
	}
	return values[0].(uint32), nil
}

func DecodeUint64(data []byte) (uint64, error) {
	values, err := ABIArguments{{Type: ABIUint64}}.Unpack(data)
	if err != nil {
		return 0, err
	}
	return values[0].(uint64), nil
}

func DecodeUint128(data []byte) (*big.Int, error) {
	values, err := ABIArguments{{Type: ABIUint128}}.Unpack(data)
	if err != nil {
		return nil, err
	}
	return values[0].(*big.Int), nil
}

func DecodeUint256(data []byte) (*big.Int, error) {
	values, err := ABIArguments{{Type: ABIUint256}}.Unpack(data)
	if err != nil {
		return nil, err
	}
	return values[0].(*big.Int), nil
}

func DecodeInt8(data []byte) (int8, error) {
	values, err := ABIArguments{{Type: ABIInt8}}.Unpack(data)
	if err != nil {
		return 0, err
	}
	return values[0].(int8), nil
}

func DecodeInt16(data []byte) (int16, error) {
	values, err := ABIArguments{{Type: ABIInt16}}.Unpack(data)
	if err != nil {
		return 0, err
	}
	return values[0].(int16), nil
}

func DecodeInt32(data []byte) (int32, error) {
	values, err := ABIArguments{{Type: ABIInt32}}.Unpack(data)
	if err != nil {
		return 0, err
	}
	return values[0].(int32), nil
}

func DecodeInt64(data []byte) (int64, error) {
	values, err := ABIArguments{{Type: ABIInt64}}.Unpack(data)
	if err != nil {
		return 0, err
	}
	return values[0].(int64), nil
}

func DecodeInt128(data []byte) (*big.Int, error) {
	values, err := ABIArguments{{Type: ABIInt128}}.Unpack(data)
	if err != nil {
		return nil, err
	}
	return values[0].(*big.Int), nil
}

func DecodeInt256(data []byte) (*big.Int, error) {
	values, err := ABIArguments{{Type: ABIInt256}}.Unpack(data)
	if err != nil {
		return nil, err
	}
	return values[0].(*big.Int), nil
}

func DecodeBytes1(data []byte) ([1]byte, error) {
	values, err := ABIArguments{{Type: ABIBytes1}}.Unpack(data)
	if err != nil {
		return [1]byte{}, err
	}
	return values[0].([1]byte), nil
}

func DecodeBytes2(data []byte) ([2]byte, error) {
	values, err := ABIArguments{{Type: ABIBytes2}}.Unpack(data)
	if err != nil {
		return [2]byte{}, err
	}
	return values[0].([2]byte), nil
}

func DecodeBytes3(data []byte) ([3]byte, error) {
	values, err := ABIArguments{{Type: ABIBytes3}}.Unpack(data)
	if err != nil {
		return [3]byte{}, err
	}
	return values[0].([3]byte), nil
}

func DecodeBytes4(data []byte) ([4]byte, error) {
	values, err := ABIArguments{{Type: ABIBytes4}}.Unpack(data)
	if err != nil {
		return [4]byte{}, err
	}
	return values[0].([4]byte), nil
}

func DecodeBytes5(data []byte) ([5]byte, error) {
	values, err := ABIArguments{{Type: ABIBytes5}}.Unpack(data)
	if err != nil {
		return [5]byte{}, err
	}
	return values[0].([5]byte), nil
}

func DecodeBytes6(data []byte) ([6]byte, error) {
	values, err := ABIArguments{{Type: ABIBytes6}}.Unpack(data)
	if err != nil {
		return [6]byte{}, err
	}
	return values[0].([6]byte), nil
}

func DecodeBytes7(data []byte) ([7]byte, error) {
	values, err := ABIArguments{{Type: ABIBytes7}}.Unpack(data)
	if err != nil {
		return [7]byte{}, err
	}
	return values[0].([7]byte), nil
}

func DecodeBytes8(data []byte) ([8]byte, error) {
	values, err := ABIArguments{{Type: ABIBytes8}}.Unpack(data)
	if err != nil {
		return [8]byte{}, err
	}
	return values[0].([8]byte), nil
}

func DecodeBytes9(data []byte) ([9]byte, error) {
	values, err := ABIArguments{{Type: ABIBytes9}}.Unpack(data)
	if err != nil {
		return [9]byte{}, err
	}
	return values[0].([9]byte), nil
}

func DecodeBytes10(data []byte) ([10]byte, error) {
	values, err := ABIArguments{{Type: ABIBytes10}}.Unpack(data)
	if err != nil {
		return [10]byte{}, err
	}
	return values[0].([10]byte), nil
}

func DecodeBytes11(data []byte) ([11]byte, error) {
	values, err := ABIArguments{{Type: ABIBytes11}}.Unpack(data)
	if err != nil {
		return [11]byte{}, err
	}
	return values[0].([11]byte), nil
}

func DecodeBytes12(data []byte) ([12]byte, error) {
	values, err := ABIArguments{{Type: ABIBytes12}}.Unpack(data)
	if err != nil {
		return [12]byte{}, err
	}
	return values[0].([12]byte), nil
}

func DecodeBytes13(data []byte) ([13]byte, error) {
	values, err := ABIArguments{{Type: ABIBytes13}}.Unpack(data)
	if err != nil {
		return [13]byte{}, err
	}
	return values[0].([13]byte), nil
}

func DecodeBytes14(data []byte) ([14]byte, error) {
	values, err := ABIArguments{{Type: ABIBytes14}}.Unpack(data)
	if err != nil {
		return [14]byte{}, err
	}
	return values[0].([14]byte), nil
}

func DecodeBytes15(data []byte) ([15]byte, error) {
	values, err := ABIArguments{{Type: ABIBytes15}}.Unpack(data)
	if err != nil {
		return [15]byte{}, err
	}
	return values[0].([15]byte), nil
}

func DecodeBytes16(data []byte) ([16]byte, error) {
	values, err := ABIArguments{{Type: ABIBytes16}}.Unpack(data)
	if err != nil {
		return [16]byte{}, err
	}
	return values[0].([16]byte), nil
}

func DecodeBytes17(data []byte) ([17]byte, error) {
	values, err := ABIArguments{{Type: ABIBytes17}}.Unpack(data)
	if err != nil {
		return [17]byte{}, err
	}
	return values[0].([17]byte), nil
}

func DecodeBytes18(data []byte) ([18]byte, error) {
	values, err := ABIArguments{{Type: ABIBytes18}}.Unpack(data)
	if err != nil {
		return [18]byte{}, err
	}
	return values[0].([18]byte), nil
}

func DecodeBytes19(data []byte) ([19]byte, error) {
	values, err := ABIArguments{{Type: ABIBytes19}}.Unpack(data)
	if err != nil {
		return [19]byte{}, err
	}
	return values[0].([19]byte), nil
}

func DecodeBytes20(data []byte) ([20]byte, error) {
	values, err := ABIArguments{{Type: ABIBytes20}}.Unpack(data)
	if err != nil {
		return [20]byte{}, err
	}
	return values[0].([20]byte), nil
}

func DecodeBytes21(data []byte) ([21]byte, error) {
	values, err := ABIArguments{{Type: ABIBytes21}}.Unpack(data)
	if err != nil {
		return [21]byte{}, err
	}
	return values[0].([21]byte), nil
}

func DecodeBytes22(data []byte) ([22]byte, error) {
	values, err := ABIArguments{{Type: ABIBytes22}}.Unpack(data)
	if err != nil {
		return [22]byte{}, err
	}
	return values[0].([22]byte), nil
}

func DecodeBytes23(data []byte) ([23]byte, error) {
	values, err := ABIArguments{{Type: ABIBytes23}}.Unpack(data)
	if err != nil {
		return [23]byte{}, err
	}
	return values[0].([23]byte), nil
}

func DecodeBytes24(data []byte) ([24]byte, error) {
	values, err := ABIArguments{{Type: ABIBytes24}}.Unpack(data)
	if err != nil {
		return [24]byte{}, err
	}
	return values[0].([24]byte), nil
}

func DecodeBytes25(data []byte) ([25]byte, error) {
	values, err := ABIArguments{{Type: ABIBytes25}}.Unpack(data)
	if err != nil {
		return [25]byte{}, err
	}
	return values[0].([25]byte), nil
}

func DecodeBytes26(data []byte) ([26]byte, error) {
	values, err := ABIArguments{{Type: ABIBytes26}}.Unpack(data)
	if err != nil {
		return [26]byte{}, err
	}
	return values[0].([26]byte), nil
}

func DecodeBytes27(data []byte) ([27]byte, error) {
	values, err := ABIArguments{{Type: ABIBytes27}}.Unpack(data)
	if err != nil {
		return [27]byte{}, err
	}
	return values[0].([27]byte), nil
}

func DecodeBytes28(data []byte) ([28]byte, error) {
	values, err := ABIArguments{{Type: ABIBytes28}}.Unpack(data)
	if err != nil {
		return [28]byte{}, err
	}
	return values[0].([28]byte), nil
}

func DecodeBytes29(data []byte) ([29]byte, error) {
	values, err := ABIArguments{{Type: ABIBytes29}}.Unpack(data)
	if err != nil {
		return [29]byte{}, err
	}
	return values[0].([29]byte), nil
}

func DecodeBytes30(data []byte) ([30]byte, error) {
	values, err := ABIArguments{{Type: ABIBytes30}}.Unpack(data)
	if err != nil {
		return [30]byte{}, err
	}
	return values[0].([30]byte), nil
}

func DecodeBytes31(data []byte) ([31]byte, error) {
	values, err := ABIArguments{{Type: ABIBytes31}}.Unpack(data)
	if err != nil {
		return [31]byte{}, err
	}
	return values[0].([31]byte), nil
}

func DecodeBytes32(data []byte) ([32]byte, error) {
	values, err := ABIArguments{{Type: ABIBytes32}}.Unpack(data)
	if err != nil {
		return [32]byte{}, err
	}
	return values[0].([32]byte), nil
}

func DecodeAddressSlice(data []byte) ([]common.Address, error) {
	values, err := ABIArguments{{Type: ABIAddressSlice}}.Unpack(data)
	if err != nil {
		return nil, err
	}
	return values[0].([]common.Address), nil
}

func DecodeBoolSlice(data []byte) ([]bool, error) {
	values, err := ABIArguments{{Type: ABIBoolSlice}}.Unpack(data)
	if err != nil {
		return nil, err
	}
	return values[0].([]bool), nil
}

func DecodeStringSlice(data []byte) ([]string, error) {
	values, err := ABIArguments{{Type: ABIStringSlice}}.Unpack(data)
	if err != nil {
		return nil, err
	}
	return values[0].([]string), nil
}

func DecodeBytesSlice(data []byte) ([][]byte, error) {
	values, err := ABIArguments{{Type: ABIBytesSlice}}.Unpack(data)
	if err != nil {
		return nil, err
	}
	return values[0].([][]byte), nil
}

func DecodeUint8Slice(data []byte) ([]uint8, error) {
	values, err := ABIArguments{{Type: ABIUint8Slice}}.Unpack(data)
	if err != nil {
		return nil, err
	}
	return values[0].([]uint8), nil
}

func DecodeUint16Slice(data []byte) ([]uint16, error) {
	values, err := ABIArguments{{Type: ABIUint16Slice}}.Unpack(data)
	if err != nil {
		return nil, err
	}
	return values[0].([]uint16), nil
}

func DecodeUint32Slice(data []byte) ([]uint32, error) {
	values, err := ABIArguments{{Type: ABIUint32Slice}}.Unpack(data)
	if err != nil {
		return nil, err
	}
	return values[0].([]uint32), nil
}

func DecodeUint64Slice(data []byte) ([]uint64, error) {
	values, err := ABIArguments{{Type: ABIUint64Slice}}.Unpack(data)
	if err != nil {
		return nil, err
	}
	return values[0].([]uint64), nil
}

func DecodeUint128Slice(data []byte) ([]*big.Int, error) {
	values, err := ABIArguments{{Type: ABIUint128Slice}}.Unpack(data)
	if err != nil {
		return nil, err
	}
	return values[0].([]*big.Int), nil
}

func DecodeUint256Slice(data []byte) ([]*big.Int, error) {
	values, err := ABIArguments{{Type: ABIUint256Slice}}.Unpack(data)
	if err != nil {
		return nil, err
	}
	return values[0].([]*big.Int), nil
}

func DecodeInt8Slice(data []byte) ([]int8, error) {
	values, err := ABIArguments{{Type: ABIInt8Slice}}.Unpack(data)
	if err != nil {
		return nil, err
	}
	return values[0].([]int8), nil
}

func DecodeInt16Slice(data []byte) ([]int16, error) {
	values, err := ABIArguments{{Type: ABIInt16Slice}}.Unpack(data)
	if err != nil {
		return nil, err
	}
	return values[0].([]int16), nil
}

func DecodeInt32Slice(data []byte) ([]int32, error) {
	values, err := ABIArguments{{Type: ABIInt32Slice}}.Unpack(data)
	if err != nil {
		return nil, err
	}
	return values[0].([]int32), nil
}

func DecodeInt64Slice(data []byte) ([]int64, error) {
	values, err := ABIArguments{{Type: ABIInt64Slice}}.Unpack(data)
	if err != nil {
		return nil, err
	}
	return values[0].([]int64), nil
}

func DecodeInt128Slice(data []byte) ([]*big.Int, error) {
	values, err := ABIArguments{{Type: ABIInt128Slice}}.Unpack(data)
	if err != nil {
		return nil, err
	}
	return values[0].([]*big.Int), nil
}

func DecodeInt256Slice(data []byte) ([]*big.Int, error) {
	values, err := ABIArguments{{Type: ABIInt256Slice}}.Unpack(data)
	if err != nil {
		return nil, err
	}
	return values[0].([]*big.Int), nil
}

func DecodeBytes32Slice(data []byte) ([][32]byte, error) {
	values, err := ABIArguments{{Type: ABIBytes32Slice}}.Unpack(data)
	if err != nil {
		return nil, err
	}
	return values[0].([][32]byte), nil
}

func EncodeAddress(arg string) (common.Address, error) {
	if !common.IsHexAddress(arg) {
		return common.Address{}, fmt.Errorf("invalid address: %s", arg)
	}
	return common.HexToAddress(arg), nil
}

func EncodeBool(arg string) (bool, error) {
	switch strings.ToLower(arg) {
	case "true", "1":
		return true, nil
	case "false", "0":
		return false, nil
	default:
		return false, fmt.Errorf("invalid bool: %s", arg)
	}
}

func EncodeString(arg string) (string, error) {
	return arg, nil
}

func EncodeBytes(arg string) ([]byte, error) {
	b, err := hexutil.Decode(arg)
	if err != nil {
		return nil, fmt.Errorf("invalid bytes: %s", err)
	}
	return b, nil
}

func EncodeUint8(arg string) (uint8, error) {
	n := new(big.Int)
	if _, ok := n.SetString(arg, 0); !ok {
		return 0, fmt.Errorf("invalid uint8: %s", arg)
	}
	return uint8(n.Uint64()), nil
}

func EncodeUint16(arg string) (uint16, error) {
	n := new(big.Int)
	if _, ok := n.SetString(arg, 0); !ok {
		return 0, fmt.Errorf("invalid uint16: %s", arg)
	}
	return uint16(n.Uint64()), nil
}

func EncodeUint32(arg string) (uint32, error) {
	n := new(big.Int)
	if _, ok := n.SetString(arg, 0); !ok {
		return 0, fmt.Errorf("invalid uint32: %s", arg)
	}
	return uint32(n.Uint64()), nil
}

func EncodeUint64(arg string) (uint64, error) {
	n := new(big.Int)
	if _, ok := n.SetString(arg, 0); !ok {
		return 0, fmt.Errorf("invalid uint64: %s", arg)
	}
	return n.Uint64(), nil
}

func EncodeUint256(arg string) (*big.Int, error) {
	n := new(big.Int)
	if _, ok := n.SetString(arg, 0); !ok {
		return nil, fmt.Errorf("invalid uint256: %s", arg)
	}
	return n, nil
}

func EncodeInt8(arg string) (int8, error) {
	n := new(big.Int)
	if _, ok := n.SetString(arg, 0); !ok {
		return 0, fmt.Errorf("invalid int8: %s", arg)
	}
	return int8(n.Int64()), nil
}

func EncodeInt16(arg string) (int16, error) {
	n := new(big.Int)
	if _, ok := n.SetString(arg, 0); !ok {
		return 0, fmt.Errorf("invalid int16: %s", arg)
	}
	return int16(n.Int64()), nil
}

func EncodeInt32(arg string) (int32, error) {
	n := new(big.Int)
	if _, ok := n.SetString(arg, 0); !ok {
		return 0, fmt.Errorf("invalid int32: %s", arg)
	}
	return int32(n.Int64()), nil
}

func EncodeInt64(arg string) (int64, error) {
	n := new(big.Int)
	if _, ok := n.SetString(arg, 0); !ok {
		return 0, fmt.Errorf("invalid int64: %s", arg)
	}
	return n.Int64(), nil
}

func EncodeInt256(arg string) (*big.Int, error) {
	n := new(big.Int)
	if _, ok := n.SetString(arg, 0); !ok {
		return nil, fmt.Errorf("invalid int256: %s", arg)
	}
	return n, nil
}

func EncodeBytes1(arg string) ([1]byte, error) {
	b, err := hexutil.Decode(arg)
	if err != nil {
		return [1]byte{}, fmt.Errorf("invalid bytes1: %s", err)
	}
	var arr [1]byte
	copy(arr[:], b)
	return arr, nil
}

func EncodeBytes2(arg string) ([2]byte, error) {
	b, err := hexutil.Decode(arg)
	if err != nil {
		return [2]byte{}, fmt.Errorf("invalid bytes2: %s", err)
	}
	var arr [2]byte
	copy(arr[:], b)
	return arr, nil
}

func EncodeBytes3(arg string) ([3]byte, error) {
	b, err := hexutil.Decode(arg)
	if err != nil {
		return [3]byte{}, fmt.Errorf("invalid bytes3: %s", err)
	}
	var arr [3]byte
	copy(arr[:], b)
	return arr, nil
}

func EncodeBytes4(arg string) ([4]byte, error) {
	b, err := hexutil.Decode(arg)
	if err != nil {
		return [4]byte{}, fmt.Errorf("invalid bytes4: %s", err)
	}
	var arr [4]byte
	copy(arr[:], b)
	return arr, nil
}

func EncodeBytes5(arg string) ([5]byte, error) {
	b, err := hexutil.Decode(arg)
	if err != nil {
		return [5]byte{}, fmt.Errorf("invalid bytes5: %s", err)
	}
	var arr [5]byte
	copy(arr[:], b)
	return arr, nil
}

func EncodeBytes6(arg string) ([6]byte, error) {
	b, err := hexutil.Decode(arg)
	if err != nil {
		return [6]byte{}, fmt.Errorf("invalid bytes6: %s", err)
	}
	var arr [6]byte
	copy(arr[:], b)
	return arr, nil
}

func EncodeBytes7(arg string) ([7]byte, error) {
	b, err := hexutil.Decode(arg)
	if err != nil {
		return [7]byte{}, fmt.Errorf("invalid bytes7: %s", err)
	}
	var arr [7]byte
	copy(arr[:], b)
	return arr, nil
}

func EncodeBytes8(arg string) ([8]byte, error) {
	b, err := hexutil.Decode(arg)
	if err != nil {
		return [8]byte{}, fmt.Errorf("invalid bytes8: %s", err)
	}
	var arr [8]byte
	copy(arr[:], b)
	return arr, nil
}

func EncodeBytes9(arg string) ([9]byte, error) {
	b, err := hexutil.Decode(arg)
	if err != nil {
		return [9]byte{}, fmt.Errorf("invalid bytes9: %s", err)
	}
	var arr [9]byte
	copy(arr[:], b)
	return arr, nil
}

func EncodeBytes10(arg string) ([10]byte, error) {
	b, err := hexutil.Decode(arg)
	if err != nil {
		return [10]byte{}, fmt.Errorf("invalid bytes10: %s", err)
	}
	var arr [10]byte
	copy(arr[:], b)
	return arr, nil
}

func EncodeBytes11(arg string) ([11]byte, error) {
	b, err := hexutil.Decode(arg)
	if err != nil {
		return [11]byte{}, fmt.Errorf("invalid bytes11: %s", err)
	}
	var arr [11]byte
	copy(arr[:], b)
	return arr, nil
}

func EncodeBytes12(arg string) ([12]byte, error) {
	b, err := hexutil.Decode(arg)
	if err != nil {
		return [12]byte{}, fmt.Errorf("invalid bytes12: %s", err)
	}
	var arr [12]byte
	copy(arr[:], b)
	return arr, nil
}

func EncodeBytes13(arg string) ([13]byte, error) {
	b, err := hexutil.Decode(arg)
	if err != nil {
		return [13]byte{}, fmt.Errorf("invalid bytes13: %s", err)
	}
	var arr [13]byte
	copy(arr[:], b)
	return arr, nil
}

func EncodeBytes14(arg string) ([14]byte, error) {
	b, err := hexutil.Decode(arg)
	if err != nil {
		return [14]byte{}, fmt.Errorf("invalid bytes14: %s", err)
	}
	var arr [14]byte
	copy(arr[:], b)
	return arr, nil
}

func EncodeBytes15(arg string) ([15]byte, error) {
	b, err := hexutil.Decode(arg)
	if err != nil {
		return [15]byte{}, fmt.Errorf("invalid bytes15: %s", err)
	}
	var arr [15]byte
	copy(arr[:], b)
	return arr, nil
}

func EncodeBytes16(arg string) ([16]byte, error) {
	b, err := hexutil.Decode(arg)
	if err != nil {
		return [16]byte{}, fmt.Errorf("invalid bytes16: %s", err)
	}
	var arr [16]byte
	copy(arr[:], b)
	return arr, nil
}

func EncodeBytes17(arg string) ([17]byte, error) {
	b, err := hexutil.Decode(arg)
	if err != nil {
		return [17]byte{}, fmt.Errorf("invalid bytes17: %s", err)
	}
	var arr [17]byte
	copy(arr[:], b)
	return arr, nil
}

func EncodeBytes18(arg string) ([18]byte, error) {
	b, err := hexutil.Decode(arg)
	if err != nil {
		return [18]byte{}, fmt.Errorf("invalid bytes18: %s", err)
	}
	var arr [18]byte
	copy(arr[:], b)
	return arr, nil
}

func EncodeBytes19(arg string) ([19]byte, error) {
	b, err := hexutil.Decode(arg)
	if err != nil {
		return [19]byte{}, fmt.Errorf("invalid bytes19: %s", err)
	}
	var arr [19]byte
	copy(arr[:], b)
	return arr, nil
}

func EncodeBytes20(arg string) ([20]byte, error) {
	b, err := hexutil.Decode(arg)
	if err != nil {
		return [20]byte{}, fmt.Errorf("invalid bytes20: %s", err)
	}
	var arr [20]byte
	copy(arr[:], b)
	return arr, nil
}

func EncodeBytes21(arg string) ([21]byte, error) {
	b, err := hexutil.Decode(arg)
	if err != nil {
		return [21]byte{}, fmt.Errorf("invalid bytes21: %s", err)
	}
	var arr [21]byte
	copy(arr[:], b)
	return arr, nil
}

func EncodeBytes22(arg string) ([22]byte, error) {
	b, err := hexutil.Decode(arg)
	if err != nil {
		return [22]byte{}, fmt.Errorf("invalid bytes22: %s", err)
	}
	var arr [22]byte
	copy(arr[:], b)
	return arr, nil
}

func EncodeBytes23(arg string) ([23]byte, error) {
	b, err := hexutil.Decode(arg)
	if err != nil {
		return [23]byte{}, fmt.Errorf("invalid bytes23: %s", err)
	}
	var arr [23]byte
	copy(arr[:], b)
	return arr, nil
}

func EncodeBytes24(arg string) ([24]byte, error) {
	b, err := hexutil.Decode(arg)
	if err != nil {
		return [24]byte{}, fmt.Errorf("invalid bytes24: %s", err)
	}
	var arr [24]byte
	copy(arr[:], b)
	return arr, nil
}

func EncodeBytes25(arg string) ([25]byte, error) {
	b, err := hexutil.Decode(arg)
	if err != nil {
		return [25]byte{}, fmt.Errorf("invalid bytes25: %s", err)
	}
	var arr [25]byte
	copy(arr[:], b)
	return arr, nil
}

func EncodeBytes26(arg string) ([26]byte, error) {
	b, err := hexutil.Decode(arg)
	if err != nil {
		return [26]byte{}, fmt.Errorf("invalid bytes26: %s", err)
	}
	var arr [26]byte
	copy(arr[:], b)
	return arr, nil
}

func EncodeBytes27(arg string) ([27]byte, error) {
	b, err := hexutil.Decode(arg)
	if err != nil {
		return [27]byte{}, fmt.Errorf("invalid bytes27: %s", err)
	}
	var arr [27]byte
	copy(arr[:], b)
	return arr, nil
}

func EncodeBytes28(arg string) ([28]byte, error) {
	b, err := hexutil.Decode(arg)
	if err != nil {
		return [28]byte{}, fmt.Errorf("invalid bytes28: %s", err)
	}
	var arr [28]byte
	copy(arr[:], b)
	return arr, nil
}

func EncodeBytes29(arg string) ([29]byte, error) {
	b, err := hexutil.Decode(arg)
	if err != nil {
		return [29]byte{}, fmt.Errorf("invalid bytes29: %s", err)
	}
	var arr [29]byte
	copy(arr[:], b)
	return arr, nil
}

func EncodeBytes30(arg string) ([30]byte, error) {
	b, err := hexutil.Decode(arg)
	if err != nil {
		return [30]byte{}, fmt.Errorf("invalid bytes30: %s", err)
	}
	var arr [30]byte
	copy(arr[:], b)
	return arr, nil
}

func EncodeBytes31(arg string) ([31]byte, error) {
	b, err := hexutil.Decode(arg)
	if err != nil {
		return [31]byte{}, fmt.Errorf("invalid bytes31: %s", err)
	}
	var arr [31]byte
	copy(arr[:], b)
	return arr, nil
}

func EncodeBytes32(arg string) ([32]byte, error) {
	b, err := hexutil.Decode(arg)
	if err != nil {
		return [32]byte{}, fmt.Errorf("invalid bytes32: %s", err)
	}
	var arr [32]byte
	copy(arr[:], b)
	return arr, nil
}
