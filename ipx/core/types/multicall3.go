package types

import (
	"encoding/binary"
	"fmt"
)

type Aggregate3 struct {
	Target    *Address
	AllowFail bool
	CallData  []byte
}

type Aggregate3s []Aggregate3

func (a *Aggregate3s) With(target *Address, callData []byte) *Aggregate3s {
	*a = append(*a, Aggregate3{Target: target, AllowFail: true, CallData: callData})
	return a
}

type Aggregator3Result struct {
	Success    bool
	ReturnData []byte
}

// DecodeAggregate3Results decodes the ABI-encoded (bool,bytes)[] returned by aggregate3.
func DecodeAggregate3Results(data []byte) ([]Aggregator3Result, error) {
	if len(data) < 64 {
		return nil, fmt.Errorf("aggregate3 result: data too short (%d bytes)", len(data))
	}

	// [0:32] offset to array (typically 32)
	// [32:64] array length
	n := binary.BigEndian.Uint64(data[56:64])
	base := 64 // start of array content (element offset table)

	if len(data) < base+int(n)*32 {
		return nil, fmt.Errorf("aggregate3 result: element offset table truncated")
	}

	results := make([]Aggregator3Result, n)
	for i := range n {
		elemOff := binary.BigEndian.Uint64(data[base+int(i)*32+24 : base+int(i)*32+32])
		elemStart := base + int(elemOff)

		if len(data) < elemStart+96 {
			return nil, fmt.Errorf("aggregate3 result: element %d: data too short", i)
		}

		results[i].Success = data[elemStart+31] != 0

		// offset to bytes within element (from elemStart), typically 64
		bytesRelOff := binary.BigEndian.Uint64(data[elemStart+56 : elemStart+64])
		bytesStart := elemStart + int(bytesRelOff)

		if len(data) < bytesStart+32 {
			return nil, fmt.Errorf("aggregate3 result: element %d: bytes length out of range", i)
		}

		bytesLen := binary.BigEndian.Uint64(data[bytesStart+24 : bytesStart+32])
		bytesEnd := bytesStart + 32 + int(bytesLen)

		if len(data) < bytesEnd {
			return nil, fmt.Errorf("aggregate3 result: element %d: bytes data out of range", i)
		}

		results[i].ReturnData = make([]byte, bytesLen)
		copy(results[i].ReturnData, data[bytesStart+32:bytesEnd])
	}

	return results, nil
}
