package types

import "math/big"

type EIP712Domain struct {
	Name     string
	Version  string
	ChainID  *big.Int
	Contract *Address
}

type EIP712Result struct {
	Digest          *Hash
	DomainSeparator *Hash
	MessageHash     *Hash
}
