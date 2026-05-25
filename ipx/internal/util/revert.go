package util

import (
	"encoding/json"
	"errors"

	"github.com/andantan/evmlab/internal/rpc"
)

// RevertData extracts the raw ABI-encoded revert bytes from an RPC error.
// Returns (nil, false) if the error is not an RPC revert or has no data.
func RevertData(err error) ([]byte, bool) {
	var rpcErr *rpc.RPCError
	if !errors.As(err, &rpcErr) || len(rpcErr.Data) == 0 {
		return nil, false
	}

	var raw string
	if e := json.Unmarshal(rpcErr.Data, &raw); e != nil {
		return nil, false
	}

	b, e := ParseHex(raw)
	if e != nil || len(b) < 4 {
		return nil, false
	}
	return b, true
}
