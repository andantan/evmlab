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

// EncodeCall ABI-encodes a function call: 4-byte selector + packed arguments.
// TODO: tuple, array, slice types are not yet supported.
func (c *abiCodec) EncodeCall(fn *types.Function, args []string) ([]byte, error) {
	if len(fn.Types) != len(args) {
		return nil, fmt.Errorf("signature has %d param(s) but got %d arg(s)", len(fn.Types), len(args))
	}

	abiArgs := make(abi.Arguments, len(fn.Types))
	for i, ts := range fn.Types {
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

// DecodeResult ABI-decodes raw return data from an eth_call into string values.
func (c *abiCodec) DecodeResult(typeStrs []string, data []byte) ([]string, error) {
	abiArgs := make(abi.Arguments, len(typeStrs))
	for i, ts := range typeStrs {
		t, err := abi.NewType(ts, "", nil)
		if err != nil {
			return nil, fmt.Errorf("invalid type %q: %s", ts, err)
		}
		abiArgs[i] = abi.Argument{Type: t}
	}

	values, err := abiArgs.Unpack(data)
	if err != nil {
		return nil, fmt.Errorf("unpack: %s", err)
	}

	result := make([]string, len(values))
	for i, v := range values {
		s, err := formatValue(v)
		if err != nil {
			return nil, fmt.Errorf("value[%d]: %s", i, err)
		}
		result[i] = s
	}
	return result, nil
}

// DecodeCall ABI-decodes calldata (selector + args) into a name→value map.
// If the signature contains no parameter names, keys are "arg0", "arg1", etc.
func (c *abiCodec) DecodeCall(fn *types.Function, data []byte) (map[string]string, error) {
	if len(data) < 4 {
		return nil, fmt.Errorf("data too short: need at least 4 bytes for selector")
	}

	abiArgs := make(abi.Arguments, len(fn.Types))
	for i, ts := range fn.Types {
		t, err := abi.NewType(ts, "", nil)
		if err != nil {
			return nil, fmt.Errorf("invalid type %q: %s", ts, err)
		}
		abiArgs[i] = abi.Argument{Type: t}
	}

	values, err := abiArgs.Unpack(data[4:])
	if err != nil {
		return nil, fmt.Errorf("unpack: %s", err)
	}

	result := make(map[string]string, len(values))
	for i, v := range values {
		key := fmt.Sprintf("arg%d", i)
		if i < len(fn.Names) && fn.Names[i] != "" {
			key = fn.Names[i]
		}
		s, err := formatValue(v)
		if err != nil {
			return nil, fmt.Errorf("%s: %s", key, err)
		}
		result[key] = s
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
		return "", fmt.Errorf("unsupported type: %T", v)
	}
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
