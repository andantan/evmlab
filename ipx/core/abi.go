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
func (c *abiCodec) EncodeCall(fn *types.Function, args []any) ([]byte, error) {
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

func convertArg(t abi.Type, arg any) (any, error) {
	switch t.T {
	case abi.SliceTy:
		items, ok := arg.([]any)
		if !ok {
			return nil, fmt.Errorf("invalid slice: expected array")
		}
		slice := reflect.MakeSlice(reflect.SliceOf(t.Elem.GetType()), len(items), len(items))
		for i, item := range items {
			v, err := convertArg(*t.Elem, item)
			if err != nil {
				return nil, fmt.Errorf("[%d]: %s", i, err)
			}
			slice.Index(i).Set(reflect.ValueOf(v))
		}
		return slice.Interface(), nil

	case abi.ArrayTy:
		items, ok := arg.([]any)
		if !ok {
			return nil, fmt.Errorf("invalid array: expected array")
		}
		if len(items) != t.Size {
			return nil, fmt.Errorf("array[%d] expects %d elements, got %d", t.Size, t.Size, len(items))
		}
		array := reflect.New(reflect.ArrayOf(t.Size, t.Elem.GetType())).Elem()
		for i, item := range items {
			v, err := convertArg(*t.Elem, item)
			if err != nil {
				return nil, fmt.Errorf("[%d]: %s", i, err)
			}
			array.Index(i).Set(reflect.ValueOf(v))
		}
		return array.Interface(), nil

	case abi.TupleTy:
		items, ok := arg.([]any)
		if !ok {
			return nil, fmt.Errorf("invalid tuple: expected array")
		}
		if len(items) != len(t.TupleElems) {
			return nil, fmt.Errorf("tuple expects %d elements, got %d", len(t.TupleElems), len(items))
		}
		tuple := reflect.New(t.TupleType).Elem()
		for i, elem := range t.TupleElems {
			v, err := convertArg(*elem, items[i])
			if err != nil {
				return nil, fmt.Errorf("tuple[%d]: %s", i, err)
			}
			tuple.Field(i).Set(reflect.ValueOf(v))
		}
		return tuple.Interface(), nil

	default:
		s, ok := arg.(string)
		if !ok {
			return nil, fmt.Errorf("expected string value")
		}
		s = strings.TrimSpace(s)
		return convertScalarArg(t, s)
	}
}

func convertScalarArg(t abi.Type, arg string) (any, error) {
	switch t.T {
	case abi.AddressTy:
		return types.EncodeAddress(arg)
	case abi.UintTy:
		switch t.Size {
		case 8:
			return types.EncodeUint8(arg)
		case 16:
			return types.EncodeUint16(arg)
		case 32:
			return types.EncodeUint32(arg)
		case 64:
			return types.EncodeUint64(arg)
		default:
			return types.EncodeUint256(arg)
		}
	case abi.IntTy:
		switch t.Size {
		case 8:
			return types.EncodeInt8(arg)
		case 16:
			return types.EncodeInt16(arg)
		case 32:
			return types.EncodeInt32(arg)
		case 64:
			return types.EncodeInt64(arg)
		default:
			return types.EncodeInt256(arg)
		}
	case abi.BoolTy:
		return types.EncodeBool(arg)
	case abi.StringTy:
		return types.EncodeString(arg)
	case abi.BytesTy:
		return types.EncodeBytes(arg)
	case abi.FixedBytesTy:
		switch t.Size {
		case 1:
			return types.EncodeBytes1(arg)
		case 2:
			return types.EncodeBytes2(arg)
		case 3:
			return types.EncodeBytes3(arg)
		case 4:
			return types.EncodeBytes4(arg)
		case 5:
			return types.EncodeBytes5(arg)
		case 6:
			return types.EncodeBytes6(arg)
		case 7:
			return types.EncodeBytes7(arg)
		case 8:
			return types.EncodeBytes8(arg)
		case 9:
			return types.EncodeBytes9(arg)
		case 10:
			return types.EncodeBytes10(arg)
		case 11:
			return types.EncodeBytes11(arg)
		case 12:
			return types.EncodeBytes12(arg)
		case 13:
			return types.EncodeBytes13(arg)
		case 14:
			return types.EncodeBytes14(arg)
		case 15:
			return types.EncodeBytes15(arg)
		case 16:
			return types.EncodeBytes16(arg)
		case 17:
			return types.EncodeBytes17(arg)
		case 18:
			return types.EncodeBytes18(arg)
		case 19:
			return types.EncodeBytes19(arg)
		case 20:
			return types.EncodeBytes20(arg)
		case 21:
			return types.EncodeBytes21(arg)
		case 22:
			return types.EncodeBytes22(arg)
		case 23:
			return types.EncodeBytes23(arg)
		case 24:
			return types.EncodeBytes24(arg)
		case 25:
			return types.EncodeBytes25(arg)
		case 26:
			return types.EncodeBytes26(arg)
		case 27:
			return types.EncodeBytes27(arg)
		case 28:
			return types.EncodeBytes28(arg)
		case 29:
			return types.EncodeBytes29(arg)
		case 30:
			return types.EncodeBytes30(arg)
		case 31:
			return types.EncodeBytes31(arg)
		case 32:
			return types.EncodeBytes32(arg)
		default:
			return nil, fmt.Errorf("unsupported bytesN size: %d", t.Size)
		}
	default:
		return nil, fmt.Errorf("unsupported type: %v", t)
	}
}
