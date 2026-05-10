package types

// Function represents a parsed ABI function signature.
// Types and Names are both empty when the function takes no parameters (e.g. "totalSupply()").
// Names is empty when the signature contained no parameter names (e.g. "approve(address,uint256)").
// Names is populated only when all parameters include a name (e.g. "approve(address spender,uint256 amount)").
type Function struct {
	Signature string
	Hash      *Hash
	Name      string
	Types     []string
	Names     []string
}

func NewFunction(signature string, hash *Hash, name string) *Function {
	return &Function{
		Signature: signature,
		Hash:      hash,
		Name:      name,
		Types:     []string{},
		Names:     []string{},
	}
}

// Selector returns the 4-byte ABI selector: the first 4 bytes of the function's keccak256 hash.
func (f *Function) Selector() []byte {
	return f.Hash.Bytes()[:4]
}
