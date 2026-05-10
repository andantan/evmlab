package core

import (
	"strconv"

	"github.com/andantan/evmlab/core/types"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/math"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/signer/core/apitypes"
)

const (
	// EIP191Prefix is the personal sign prefix defined by EIP-191.
	EIP191Prefix = "\x19Ethereum Signed Message:\n"

	// EIP712Prefix is the two-byte version prefix defined by EIP-712.
	EIP712Prefix = "\x19\x01"
)

type hasher struct{}

var Hasher = new(hasher)

// Hash computes the raw Keccak256 hash of the input data and returns it as a wrapped Hash.
// No prefix or transformation is applied — the digest is the direct Keccak256 output of m.
func (h *hasher) Hash(m []byte) *types.Hash {
	digest := crypto.Keccak256(m)
	hash := common.BytesToHash(digest)
	return types.NewHash(hash)
}

func (h *hasher) HashString(s string) *types.Hash {
	return h.Hash([]byte(s))
}

// EIP191 applies personal sign prefix and returns the Keccak256 hash.
// The prefix "\x19Ethereum Signed Message:\n" + len(m) is prepended before hashing,
// matching the digest produced by eth_sign / personal_sign in wallets such as MetaMask.
func (h *hasher) EIP191(m []byte) *types.Hash {
	msg := make([]byte, 0, len(EIP191Prefix)+len(strconv.Itoa(len(m)))+len(m))
	msg = append(msg, EIP191Prefix...)
	msg = strconv.AppendInt(msg, int64(len(m)), 10)
	msg = append(msg, m...)

	return h.Hash(msg)
}

// EIP712 computes the EIP-712 digest for the given domain, function, and arguments.
//
//	digest = keccak256("\x19\x01" || domainSeparator || hashStruct(message))
func (h *hasher) EIP712(domain *types.EIP712Domain, fn *types.Function, args []string) (*types.EIP712Result, error) {
	primaryType := make([]apitypes.Type, len(fn.Types))
	for i := range fn.Types {
		primaryType[i] = apitypes.Type{Name: fn.Names[i], Type: fn.Types[i]}
	}

	message := apitypes.TypedDataMessage{}
	for i, name := range fn.Names {
		message[name] = args[i]
	}

	typedData := apitypes.TypedData{
		Types: apitypes.Types{
			"EIP712Domain": []apitypes.Type{
				{Name: "name", Type: "string"},
				{Name: "version", Type: "string"},
				{Name: "chainId", Type: "uint256"},
				{Name: "verifyingContract", Type: "address"},
			},
			fn.Name: primaryType,
		},
		PrimaryType: fn.Name,
		Domain: apitypes.TypedDataDomain{
			Name:              domain.Name,
			Version:           domain.Version,
			ChainId:           (*math.HexOrDecimal256)(domain.ChainID),
			VerifyingContract: domain.Contract.String(),
		},
		Message: message,
	}

	domainSepBytes, err := typedData.HashStruct("EIP712Domain", typedData.Domain.Map())
	if err != nil {
		return nil, err
	}

	messageHashBytes, err := typedData.HashStruct(typedData.PrimaryType, typedData.Message)
	if err != nil {
		return nil, err
	}

	digest := crypto.Keccak256(
		[]byte(EIP712Prefix),
		domainSepBytes,
		messageHashBytes,
	)

	return &types.EIP712Result{
		Digest:          types.NewHash(common.BytesToHash(digest)),
		DomainSeparator: types.NewHash(common.BytesToHash(domainSepBytes)),
		MessageHash:     types.NewHash(common.BytesToHash(messageHashBytes)),
	}, nil
}
