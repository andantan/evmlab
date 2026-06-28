package core

import (
	"fmt"
	"math/big"
	"reflect"
	"strconv"
	"strings"

	"github.com/andantan/evmlab/core/types"
	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
)

type abiCodec struct{}

var ABI = new(abiCodec)

// ParseFunctionSignature parses a function signature with or without parameter names.
// Accepts both "approve(address,uint256)" and "approve(address spender,uint256 amount)".
// Returns the function name, parameter types, and parameter names.
// names is empty (len 0) when the signature contains no parameter names.
func (c *abiCodec) ParseFunctionSignature(sig string) (*types.Function, error) {
	sig = strings.TrimSpace(sig)
	idx := strings.Index(sig, "(")
	if idx < 0 || !strings.HasSuffix(sig, ")") {
		return nil, fmt.Errorf("invalid signature: expected name(...)")
	}

	name := strings.TrimSpace(sig[:idx])
	if name == "" {
		return nil, fmt.Errorf("invalid signature: missing function name")
	}

	inner := strings.TrimSpace(sig[idx+1 : len(sig)-1])
	if inner == "" {
		signature := name + "()"
		return types.NewFunction(signature, Hasher.HashString(signature), name), nil
	}

	params := strings.Split(inner, ",")
	paramTypes := make([]string, len(params))
	paramNames := make([]string, len(params))
	hasNames := false

	for i, p := range params {
		fields := strings.Fields(p)

		switch len(fields) {
		case 1:
			paramTypes[i] = fields[0]
		case 2:
			paramTypes[i], paramNames[i], hasNames = fields[0], fields[1], true
		default:
			return nil, fmt.Errorf("invalid param %q", strings.TrimSpace(p))
		}
	}

	if !hasNames {
		paramNames = []string{}
	}

	signature := name + "(" + strings.Join(paramTypes, ",") + ")"
	fn := types.NewFunction(signature, Hasher.HashString(signature), name)
	fn.Types = paramTypes
	fn.Names = paramNames
	return fn, nil
}

// EncodeArgs ABI-encodes arguments only (no selector). Used for constructor calldata.
func (c *abiCodec) EncodeArgs(typeStrs []string, args []string) ([]byte, error) {
	if len(typeStrs) != len(args) {
		return nil, fmt.Errorf("type count (%d) does not match arg count (%d)", len(typeStrs), len(args))
	}

	abiArgs := make(abi.Arguments, len(typeStrs))
	for i, ts := range typeStrs {
		t, err := abi.NewType(ts, "", nil)
		if err != nil {
			return nil, fmt.Errorf("invalid type %q: %s", ts, err)
		}
		abiArgs[i] = abi.Argument{Type: t}
	}

	goArgs := make([]any, len(args))
	for i, arg := range args {
		v, err := convertArg(abiArgs[i].Type, arg)
		if err != nil {
			return nil, fmt.Errorf("arg[%d] (%s): %s", i, typeStrs[i], err)
		}
		goArgs[i] = v
	}

	return abiArgs.Pack(goArgs...)
}

// EncodeCall ABI-encodes a function call: 4-byte selector + packed arguments.
// TODO: tuple, array, slice types are not yet supported.
func (c *abiCodec) EncodeCall(fn *types.Function, args []string) ([]byte, error) {
	if len(fn.Types) != len(args) {
		return nil, fmt.Errorf("signature has %d param(s) but got %d arg(s)", len(fn.Types), len(args))
	}

	abiArgs, err := buildABIArgs(fn.Types)
	if err != nil {
		return nil, err
	}

	goArgs := make([]any, len(args))
	for i, arg := range args {
		v, err := convertArg(abiArgs[i].Type, arg)
		if err != nil {
			return nil, fmt.Errorf("arg[%d] (%s): %s", i, fn.Types[i], err)
		}
		goArgs[i] = v
	}

	packed, err := abiArgs.Pack(goArgs...)
	if err != nil {
		return nil, fmt.Errorf("pack: %s", err)
	}

	return append(fn.Selector(), packed...), nil
}

// DecodeResult ABI-decodes raw return data from an eth_call into values.
func (c *abiCodec) DecodeResult(typeStrs []string, data []byte) ([]any, error) {
	abiArgs, err := buildABIArgs(typeStrs)
	if err != nil {
		return nil, err
	}

	values, err := abiArgs.Unpack(data)
	if err != nil {
		return nil, fmt.Errorf("unpack: %s", err)
	}

	return formatDecodedValues(values)
}

// DecodeCall ABI-decodes calldata (selector + args) into a name→value map.
// If the signature contains no parameter names, keys are "arg0", "arg1", etc.
func (c *abiCodec) DecodeCall(fn *types.Function, data []byte) (map[string]any, error) {
	if len(data) < 4 {
		return nil, fmt.Errorf("data too short: need at least 4 bytes for selector")
	}

	abiArgs, err := buildABIArgs(fn.Types)
	if err != nil {
		return nil, err
	}

	values, err := abiArgs.Unpack(data[4:])
	if err != nil {
		return nil, fmt.Errorf("unpack: %s", err)
	}

	return formatDecodedMap(values, fn.Names)
}

// DecodeRevert ABI-decodes revert data (4-byte selector + args) into a name→value map.
// Validates that the selector matches the given error signature. Returns an error if
// the selector does not match or the data cannot be unpacked.
func (c *abiCodec) DecodeRevert(fn *types.Function, data []byte) (map[string]any, error) {
	if len(data) < 4 {
		return nil, fmt.Errorf("data too short: need at least 4 bytes for selector")
	}

	sel := fn.Selector()
	if sel[0] != data[0] || sel[1] != data[1] || sel[2] != data[2] || sel[3] != data[3] {
		return nil, fmt.Errorf("selector mismatch: expected 0x%x, got 0x%x", sel, data[:4])
	}

	if len(fn.Types) == 0 {
		return map[string]any{}, nil
	}

	abiArgs, err := buildABIArgs(fn.Types)
	if err != nil {
		return nil, err
	}

	values, err := abiArgs.Unpack(data[4:])
	if err != nil {
		return nil, fmt.Errorf("unpack: %s", err)
	}

	return formatDecodedMap(values, fn.Names)
}

// DecodeErrorData decodes ABI-encoded revert data using the provided selector→signature map.
// Returns the parsed Function (name + parameter schema) and a map of parameter name→value pairs.
// If the signature has no parameter names, keys fall back to "arg0", "arg1", etc.
func (c *abiCodec) DecodeErrorData(data []byte, signatures map[types.Selector]string) (*types.Function, map[string]any, error) {
	if len(data) < 4 {
		return nil, nil, fmt.Errorf("data too short: need at least 4 bytes for selector")
	}

	sel := types.Selector{data[0], data[1], data[2], data[3]}
	sig, ok := signatures[sel]
	if !ok {
		return nil, nil, fmt.Errorf("unknown error selector: %s", sel)
	}

	fn, err := c.ParseFunctionSignature(sig)
	if err != nil {
		return nil, nil, err
	}

	if len(fn.Types) == 0 {
		return fn, map[string]any{}, nil
	}

	abiArgs, err := buildABIArgs(fn.Types)
	if err != nil {
		return nil, nil, err
	}

	values, err := abiArgs.Unpack(data[4:])
	if err != nil {
		return nil, nil, fmt.Errorf("unpack: %s", err)
	}

	params, err := formatDecodedMap(values, fn.Names)
	if err != nil {
		return nil, nil, err
	}

	return fn, params, nil
}

func isTypedSlice(rv reflect.Value) bool {
	return (rv.Kind() == reflect.Slice && rv.Type() != reflect.TypeOf([]byte(nil))) ||
		(rv.Kind() == reflect.Array && rv.Type().Elem().Kind() != reflect.Uint8)
}

func formatSliceValue(rv reflect.Value) ([]any, error) {
	result := make([]any, rv.Len())
	for i := range result {
		elem := rv.Index(i).Interface()
		erv := reflect.ValueOf(elem)
		if isTypedSlice(erv) {
			sub, err := formatSliceValue(erv)
			if err != nil {
				return nil, fmt.Errorf("[%d]: %s", i, err)
			}
			result[i] = sub
		} else {
			s, err := formatValue(elem)
			if err != nil {
				return nil, fmt.Errorf("[%d]: %s", i, err)
			}
			result[i] = s
		}
	}
	return result, nil
}

func formatValue(v any) (string, error) {
	switch val := v.(type) {
	case common.Address:
		return val.Hex(), nil
	case *big.Int:
		return val.String(), nil
	case uint8:
		return strconv.FormatUint(uint64(val), 10), nil
	case uint16:
		return strconv.FormatUint(uint64(val), 10), nil
	case uint32:
		return strconv.FormatUint(uint64(val), 10), nil
	case uint64:
		return strconv.FormatUint(val, 10), nil
	case int8:
		return strconv.FormatInt(int64(val), 10), nil
	case int16:
		return strconv.FormatInt(int64(val), 10), nil
	case int32:
		return strconv.FormatInt(int64(val), 10), nil
	case int64:
		return strconv.FormatInt(val, 10), nil
	case bool:
		if val {
			return "true", nil
		}
		return "false", nil
	case string:
		return val, nil
	case []byte:
		return hexutil.Encode(val), nil
	default:
		rv := reflect.ValueOf(v)
		if rv.Kind() == reflect.Array && rv.Type().Elem().Kind() == reflect.Uint8 {
			b := make([]byte, rv.Len())
			for i := range b {
				b[i] = rv.Index(i).Interface().(byte)
			}
			return hexutil.Encode(b), nil
		}
		if rv.Kind() == reflect.Array || rv.Kind() == reflect.Slice {
			return formatArrayValue(rv)
		}
		return "", fmt.Errorf("unsupported type: %T", v)
	}
}

func formatArrayValue(rv reflect.Value) (string, error) {
	values := make([]string, rv.Len())
	for i := range values {
		s, err := formatValue(rv.Index(i).Interface())
		if err != nil {
			return "", fmt.Errorf("[%d]: %s", i, err)
		}
		values[i] = s
	}
	return "[" + strings.Join(values, ",") + "]", nil
}

func buildABIArgs(typeStrs []string) (abi.Arguments, error) {
	args := make(abi.Arguments, len(typeStrs))
	for i, ts := range typeStrs {
		t, err := abi.NewType(ts, "", nil)
		if err != nil {
			return nil, fmt.Errorf("invalid type %q: %s", ts, err)
		}
		args[i] = abi.Argument{Type: t}
	}
	return args, nil
}

func formatDecodedValues(values []any) ([]any, error) {
	result := make([]any, len(values))
	for i, v := range values {
		rv := reflect.ValueOf(v)
		if isTypedSlice(rv) {
			arr, err := formatSliceValue(rv)
			if err != nil {
				return nil, fmt.Errorf("value[%d]: %s", i, err)
			}
			result[i] = arr
		} else {
			s, err := formatValue(v)
			if err != nil {
				return nil, fmt.Errorf("value[%d]: %s", i, err)
			}
			result[i] = s
		}
	}
	return result, nil
}

func formatDecodedMap(values []any, names []string) (map[string]any, error) {
	result := make(map[string]any, len(values))
	for i, v := range values {
		key := fmt.Sprintf("arg%d", i)
		if i < len(names) && names[i] != "" {
			key = names[i]
		}
		rv := reflect.ValueOf(v)
		if isTypedSlice(rv) {
			arr, err := formatSliceValue(rv)
			if err != nil {
				return nil, fmt.Errorf("%s: %s", key, err)
			}
			result[key] = arr
		} else {
			s, err := formatValue(v)
			if err != nil {
				return nil, fmt.Errorf("%s: %s", key, err)
			}
			result[key] = s
		}
	}
	return result, nil
}

func convertArg(t abi.Type, arg string) (any, error) {
	arg = strings.TrimSpace(arg)
	switch t.T {
	case abi.AddressTy:
		if !common.IsHexAddress(arg) {
			return nil, fmt.Errorf("invalid address: %s", arg)
		}
		return common.HexToAddress(arg), nil

	case abi.UintTy:
		n := new(big.Int)
		if _, ok := n.SetString(arg, 0); !ok {
			return nil, fmt.Errorf("invalid uint: %s", arg)
		}

		switch t.Size {
		case 8:
			return uint8(n.Uint64()), nil
		case 16:
			return uint16(n.Uint64()), nil
		case 32:
			return uint32(n.Uint64()), nil
		case 64:
			return n.Uint64(), nil
		default:
			return n, nil
		}

	case abi.IntTy:
		n := new(big.Int)
		if _, ok := n.SetString(arg, 0); !ok {
			return nil, fmt.Errorf("invalid int: %s", arg)
		}

		switch t.Size {
		case 8:
			return int8(n.Int64()), nil
		case 16:
			return int16(n.Int64()), nil
		case 32:
			return int32(n.Int64()), nil
		case 64:
			return n.Int64(), nil
		default:
			return n, nil
		}

	case abi.BoolTy:
		switch strings.ToLower(arg) {
		case "true", "1":
			return true, nil
		case "false", "0":
			return false, nil
		default:
			return nil, fmt.Errorf("invalid bool: %s", arg)
		}

	case abi.StringTy:
		return arg, nil

	case abi.BytesTy:
		b, err := hexutil.Decode(arg)
		if err != nil {
			return nil, fmt.Errorf("invalid bytes: %s", err)
		}
		return b, nil

	case abi.FixedBytesTy:
		b, err := hexutil.Decode(arg)
		if err != nil {
			return nil, fmt.Errorf("invalid bytes%d: %s", t.Size, err)
		}
		arrType := reflect.ArrayOf(t.Size, reflect.TypeOf(byte(0)))
		arr := reflect.New(arrType).Elem()
		for i := 0; i < len(b) && i < t.Size; i++ {
			arr.Index(i).Set(reflect.ValueOf(b[i]))
		}
		return arr.Interface(), nil

	default:
		return nil, fmt.Errorf("unsupported type: %v", t)
	}
}

func (c *abiCodec) DecodeAddress(data []byte) (common.Address, error) {
	values, err := types.ABIArguments{{Type: types.ABIAddress}}.Unpack(data)
	if err != nil {
		return common.Address{}, err
	}
	return values[0].(common.Address), nil
}

func (c *abiCodec) DecodeBool(data []byte) (bool, error) {
	values, err := types.ABIArguments{{Type: types.ABIBool}}.Unpack(data)
	if err != nil {
		return false, err
	}
	return values[0].(bool), nil
}

func (c *abiCodec) DecodeString(data []byte) (string, error) {
	values, err := types.ABIArguments{{Type: types.ABIString}}.Unpack(data)
	if err != nil {
		return "", err
	}
	return values[0].(string), nil
}

func (c *abiCodec) DecodeBytes(data []byte) ([]byte, error) {
	values, err := types.ABIArguments{{Type: types.ABIBytes}}.Unpack(data)
	if err != nil {
		return nil, err
	}
	return values[0].([]byte), nil
}

func (c *abiCodec) DecodeUint8(data []byte) (uint8, error) {
	values, err := types.ABIArguments{{Type: types.ABIUint8}}.Unpack(data)
	if err != nil {
		return 0, err
	}
	return values[0].(uint8), nil
}

func (c *abiCodec) DecodeUint16(data []byte) (uint16, error) {
	values, err := types.ABIArguments{{Type: types.ABIUint16}}.Unpack(data)
	if err != nil {
		return 0, err
	}
	return values[0].(uint16), nil
}

func (c *abiCodec) DecodeUint32(data []byte) (uint32, error) {
	values, err := types.ABIArguments{{Type: types.ABIUint32}}.Unpack(data)
	if err != nil {
		return 0, err
	}
	return values[0].(uint32), nil
}

func (c *abiCodec) DecodeUint64(data []byte) (uint64, error) {
	values, err := types.ABIArguments{{Type: types.ABIUint64}}.Unpack(data)
	if err != nil {
		return 0, err
	}
	return values[0].(uint64), nil
}

func (c *abiCodec) DecodeUint128(data []byte) (*big.Int, error) {
	values, err := types.ABIArguments{{Type: types.ABIUint128}}.Unpack(data)
	if err != nil {
		return nil, err
	}
	return values[0].(*big.Int), nil
}

func (c *abiCodec) DecodeUint256(data []byte) (*big.Int, error) {
	values, err := types.ABIArguments{{Type: types.ABIUint256}}.Unpack(data)
	if err != nil {
		return nil, err
	}
	return values[0].(*big.Int), nil
}

func (c *abiCodec) DecodeInt8(data []byte) (int8, error) {
	values, err := types.ABIArguments{{Type: types.ABIInt8}}.Unpack(data)
	if err != nil {
		return 0, err
	}
	return values[0].(int8), nil
}

func (c *abiCodec) DecodeInt16(data []byte) (int16, error) {
	values, err := types.ABIArguments{{Type: types.ABIInt16}}.Unpack(data)
	if err != nil {
		return 0, err
	}
	return values[0].(int16), nil
}

func (c *abiCodec) DecodeInt32(data []byte) (int32, error) {
	values, err := types.ABIArguments{{Type: types.ABIInt32}}.Unpack(data)
	if err != nil {
		return 0, err
	}
	return values[0].(int32), nil
}

func (c *abiCodec) DecodeInt64(data []byte) (int64, error) {
	values, err := types.ABIArguments{{Type: types.ABIInt64}}.Unpack(data)
	if err != nil {
		return 0, err
	}
	return values[0].(int64), nil
}

func (c *abiCodec) DecodeInt128(data []byte) (*big.Int, error) {
	values, err := types.ABIArguments{{Type: types.ABIInt128}}.Unpack(data)
	if err != nil {
		return nil, err
	}
	return values[0].(*big.Int), nil
}

func (c *abiCodec) DecodeInt256(data []byte) (*big.Int, error) {
	values, err := types.ABIArguments{{Type: types.ABIInt256}}.Unpack(data)
	if err != nil {
		return nil, err
	}
	return values[0].(*big.Int), nil
}

func (c *abiCodec) DecodeBytes1(data []byte) ([1]byte, error) {
	values, err := types.ABIArguments{{Type: types.ABIBytes1}}.Unpack(data)
	if err != nil {
		return [1]byte{}, err
	}
	return values[0].([1]byte), nil
}

func (c *abiCodec) DecodeBytes2(data []byte) ([2]byte, error) {
	values, err := types.ABIArguments{{Type: types.ABIBytes2}}.Unpack(data)
	if err != nil {
		return [2]byte{}, err
	}
	return values[0].([2]byte), nil
}

func (c *abiCodec) DecodeBytes3(data []byte) ([3]byte, error) {
	values, err := types.ABIArguments{{Type: types.ABIBytes3}}.Unpack(data)
	if err != nil {
		return [3]byte{}, err
	}
	return values[0].([3]byte), nil
}

func (c *abiCodec) DecodeBytes4(data []byte) ([4]byte, error) {
	values, err := types.ABIArguments{{Type: types.ABIBytes4}}.Unpack(data)
	if err != nil {
		return [4]byte{}, err
	}
	return values[0].([4]byte), nil
}

func (c *abiCodec) DecodeBytes5(data []byte) ([5]byte, error) {
	values, err := types.ABIArguments{{Type: types.ABIBytes5}}.Unpack(data)
	if err != nil {
		return [5]byte{}, err
	}
	return values[0].([5]byte), nil
}

func (c *abiCodec) DecodeBytes6(data []byte) ([6]byte, error) {
	values, err := types.ABIArguments{{Type: types.ABIBytes6}}.Unpack(data)
	if err != nil {
		return [6]byte{}, err
	}
	return values[0].([6]byte), nil
}

func (c *abiCodec) DecodeBytes7(data []byte) ([7]byte, error) {
	values, err := types.ABIArguments{{Type: types.ABIBytes7}}.Unpack(data)
	if err != nil {
		return [7]byte{}, err
	}
	return values[0].([7]byte), nil
}

func (c *abiCodec) DecodeBytes8(data []byte) ([8]byte, error) {
	values, err := types.ABIArguments{{Type: types.ABIBytes8}}.Unpack(data)
	if err != nil {
		return [8]byte{}, err
	}
	return values[0].([8]byte), nil
}

func (c *abiCodec) DecodeBytes9(data []byte) ([9]byte, error) {
	values, err := types.ABIArguments{{Type: types.ABIBytes9}}.Unpack(data)
	if err != nil {
		return [9]byte{}, err
	}
	return values[0].([9]byte), nil
}

func (c *abiCodec) DecodeBytes10(data []byte) ([10]byte, error) {
	values, err := types.ABIArguments{{Type: types.ABIBytes10}}.Unpack(data)
	if err != nil {
		return [10]byte{}, err
	}
	return values[0].([10]byte), nil
}

func (c *abiCodec) DecodeBytes11(data []byte) ([11]byte, error) {
	values, err := types.ABIArguments{{Type: types.ABIBytes11}}.Unpack(data)
	if err != nil {
		return [11]byte{}, err
	}
	return values[0].([11]byte), nil
}

func (c *abiCodec) DecodeBytes12(data []byte) ([12]byte, error) {
	values, err := types.ABIArguments{{Type: types.ABIBytes12}}.Unpack(data)
	if err != nil {
		return [12]byte{}, err
	}
	return values[0].([12]byte), nil
}

func (c *abiCodec) DecodeBytes13(data []byte) ([13]byte, error) {
	values, err := types.ABIArguments{{Type: types.ABIBytes13}}.Unpack(data)
	if err != nil {
		return [13]byte{}, err
	}
	return values[0].([13]byte), nil
}

func (c *abiCodec) DecodeBytes14(data []byte) ([14]byte, error) {
	values, err := types.ABIArguments{{Type: types.ABIBytes14}}.Unpack(data)
	if err != nil {
		return [14]byte{}, err
	}
	return values[0].([14]byte), nil
}

func (c *abiCodec) DecodeBytes15(data []byte) ([15]byte, error) {
	values, err := types.ABIArguments{{Type: types.ABIBytes15}}.Unpack(data)
	if err != nil {
		return [15]byte{}, err
	}
	return values[0].([15]byte), nil
}

func (c *abiCodec) DecodeBytes16(data []byte) ([16]byte, error) {
	values, err := types.ABIArguments{{Type: types.ABIBytes16}}.Unpack(data)
	if err != nil {
		return [16]byte{}, err
	}
	return values[0].([16]byte), nil
}

func (c *abiCodec) DecodeBytes17(data []byte) ([17]byte, error) {
	values, err := types.ABIArguments{{Type: types.ABIBytes17}}.Unpack(data)
	if err != nil {
		return [17]byte{}, err
	}
	return values[0].([17]byte), nil
}

func (c *abiCodec) DecodeBytes18(data []byte) ([18]byte, error) {
	values, err := types.ABIArguments{{Type: types.ABIBytes18}}.Unpack(data)
	if err != nil {
		return [18]byte{}, err
	}
	return values[0].([18]byte), nil
}

func (c *abiCodec) DecodeBytes19(data []byte) ([19]byte, error) {
	values, err := types.ABIArguments{{Type: types.ABIBytes19}}.Unpack(data)
	if err != nil {
		return [19]byte{}, err
	}
	return values[0].([19]byte), nil
}

func (c *abiCodec) DecodeBytes20(data []byte) ([20]byte, error) {
	values, err := types.ABIArguments{{Type: types.ABIBytes20}}.Unpack(data)
	if err != nil {
		return [20]byte{}, err
	}
	return values[0].([20]byte), nil
}

func (c *abiCodec) DecodeBytes21(data []byte) ([21]byte, error) {
	values, err := types.ABIArguments{{Type: types.ABIBytes21}}.Unpack(data)
	if err != nil {
		return [21]byte{}, err
	}
	return values[0].([21]byte), nil
}

func (c *abiCodec) DecodeBytes22(data []byte) ([22]byte, error) {
	values, err := types.ABIArguments{{Type: types.ABIBytes22}}.Unpack(data)
	if err != nil {
		return [22]byte{}, err
	}
	return values[0].([22]byte), nil
}

func (c *abiCodec) DecodeBytes23(data []byte) ([23]byte, error) {
	values, err := types.ABIArguments{{Type: types.ABIBytes23}}.Unpack(data)
	if err != nil {
		return [23]byte{}, err
	}
	return values[0].([23]byte), nil
}

func (c *abiCodec) DecodeBytes24(data []byte) ([24]byte, error) {
	values, err := types.ABIArguments{{Type: types.ABIBytes24}}.Unpack(data)
	if err != nil {
		return [24]byte{}, err
	}
	return values[0].([24]byte), nil
}

func (c *abiCodec) DecodeBytes25(data []byte) ([25]byte, error) {
	values, err := types.ABIArguments{{Type: types.ABIBytes25}}.Unpack(data)
	if err != nil {
		return [25]byte{}, err
	}
	return values[0].([25]byte), nil
}

func (c *abiCodec) DecodeBytes26(data []byte) ([26]byte, error) {
	values, err := types.ABIArguments{{Type: types.ABIBytes26}}.Unpack(data)
	if err != nil {
		return [26]byte{}, err
	}
	return values[0].([26]byte), nil
}

func (c *abiCodec) DecodeBytes27(data []byte) ([27]byte, error) {
	values, err := types.ABIArguments{{Type: types.ABIBytes27}}.Unpack(data)
	if err != nil {
		return [27]byte{}, err
	}
	return values[0].([27]byte), nil
}

func (c *abiCodec) DecodeBytes28(data []byte) ([28]byte, error) {
	values, err := types.ABIArguments{{Type: types.ABIBytes28}}.Unpack(data)
	if err != nil {
		return [28]byte{}, err
	}
	return values[0].([28]byte), nil
}

func (c *abiCodec) DecodeBytes29(data []byte) ([29]byte, error) {
	values, err := types.ABIArguments{{Type: types.ABIBytes29}}.Unpack(data)
	if err != nil {
		return [29]byte{}, err
	}
	return values[0].([29]byte), nil
}

func (c *abiCodec) DecodeBytes30(data []byte) ([30]byte, error) {
	values, err := types.ABIArguments{{Type: types.ABIBytes30}}.Unpack(data)
	if err != nil {
		return [30]byte{}, err
	}
	return values[0].([30]byte), nil
}

func (c *abiCodec) DecodeBytes31(data []byte) ([31]byte, error) {
	values, err := types.ABIArguments{{Type: types.ABIBytes31}}.Unpack(data)
	if err != nil {
		return [31]byte{}, err
	}
	return values[0].([31]byte), nil
}

func (c *abiCodec) DecodeBytes32(data []byte) ([32]byte, error) {
	values, err := types.ABIArguments{{Type: types.ABIBytes32}}.Unpack(data)
	if err != nil {
		return [32]byte{}, err
	}
	return values[0].([32]byte), nil
}

// -- Slice decoders --

func (c *abiCodec) DecodeAddressSlice(data []byte) ([]common.Address, error) {
	values, err := types.ABIArguments{{Type: types.ABIAddressSlice}}.Unpack(data)
	if err != nil {
		return nil, err
	}
	return values[0].([]common.Address), nil
}

func (c *abiCodec) DecodeBoolSlice(data []byte) ([]bool, error) {
	values, err := types.ABIArguments{{Type: types.ABIBoolSlice}}.Unpack(data)
	if err != nil {
		return nil, err
	}
	return values[0].([]bool), nil
}

func (c *abiCodec) DecodeStringSlice(data []byte) ([]string, error) {
	values, err := types.ABIArguments{{Type: types.ABIStringSlice}}.Unpack(data)
	if err != nil {
		return nil, err
	}
	return values[0].([]string), nil
}

func (c *abiCodec) DecodeBytesSlice(data []byte) ([][]byte, error) {
	values, err := types.ABIArguments{{Type: types.ABIBytesSlice}}.Unpack(data)
	if err != nil {
		return nil, err
	}
	return values[0].([][]byte), nil
}

func (c *abiCodec) DecodeUint8Slice(data []byte) ([]uint8, error) {
	values, err := types.ABIArguments{{Type: types.ABIUint8Slice}}.Unpack(data)
	if err != nil {
		return nil, err
	}
	return values[0].([]uint8), nil
}

func (c *abiCodec) DecodeUint16Slice(data []byte) ([]uint16, error) {
	values, err := types.ABIArguments{{Type: types.ABIUint16Slice}}.Unpack(data)
	if err != nil {
		return nil, err
	}
	return values[0].([]uint16), nil
}

func (c *abiCodec) DecodeUint32Slice(data []byte) ([]uint32, error) {
	values, err := types.ABIArguments{{Type: types.ABIUint32Slice}}.Unpack(data)
	if err != nil {
		return nil, err
	}
	return values[0].([]uint32), nil
}

func (c *abiCodec) DecodeUint64Slice(data []byte) ([]uint64, error) {
	values, err := types.ABIArguments{{Type: types.ABIUint64Slice}}.Unpack(data)
	if err != nil {
		return nil, err
	}
	return values[0].([]uint64), nil
}

func (c *abiCodec) DecodeUint128Slice(data []byte) ([]*big.Int, error) {
	values, err := types.ABIArguments{{Type: types.ABIUint128Slice}}.Unpack(data)
	if err != nil {
		return nil, err
	}
	return values[0].([]*big.Int), nil
}

func (c *abiCodec) DecodeUint256Slice(data []byte) ([]*big.Int, error) {
	values, err := types.ABIArguments{{Type: types.ABIUint256Slice}}.Unpack(data)
	if err != nil {
		return nil, err
	}
	return values[0].([]*big.Int), nil
}

func (c *abiCodec) DecodeInt8Slice(data []byte) ([]int8, error) {
	values, err := types.ABIArguments{{Type: types.ABIInt8Slice}}.Unpack(data)
	if err != nil {
		return nil, err
	}
	return values[0].([]int8), nil
}

func (c *abiCodec) DecodeInt16Slice(data []byte) ([]int16, error) {
	values, err := types.ABIArguments{{Type: types.ABIInt16Slice}}.Unpack(data)
	if err != nil {
		return nil, err
	}
	return values[0].([]int16), nil
}

func (c *abiCodec) DecodeInt32Slice(data []byte) ([]int32, error) {
	values, err := types.ABIArguments{{Type: types.ABIInt32Slice}}.Unpack(data)
	if err != nil {
		return nil, err
	}
	return values[0].([]int32), nil
}

func (c *abiCodec) DecodeInt64Slice(data []byte) ([]int64, error) {
	values, err := types.ABIArguments{{Type: types.ABIInt64Slice}}.Unpack(data)
	if err != nil {
		return nil, err
	}
	return values[0].([]int64), nil
}

func (c *abiCodec) DecodeInt128Slice(data []byte) ([]*big.Int, error) {
	values, err := types.ABIArguments{{Type: types.ABIInt128Slice}}.Unpack(data)
	if err != nil {
		return nil, err
	}
	return values[0].([]*big.Int), nil
}

func (c *abiCodec) DecodeInt256Slice(data []byte) ([]*big.Int, error) {
	values, err := types.ABIArguments{{Type: types.ABIInt256Slice}}.Unpack(data)
	if err != nil {
		return nil, err
	}
	return values[0].([]*big.Int), nil
}

func (c *abiCodec) DecodeBytes32Slice(data []byte) ([][32]byte, error) {
	values, err := types.ABIArguments{{Type: types.ABIBytes32Slice}}.Unpack(data)
	if err != nil {
		return nil, err
	}
	return values[0].([][32]byte), nil
}
