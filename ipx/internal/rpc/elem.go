package rpc

type Elem struct {
	Method string
	Params any
	Result any
}

type Elems []Elem

func (e *Elems) With(elem Elem) *Elems {
	*e = append(*e, elem)
	return e
}

func (e *Elems) Len() int {
	return len(*e)
}

func (e *Elems) GetMethod(i int) string {
	return (*e)[i].Method
}

func (e *Elems) GetParams(i int) any {
	return (*e)[i].Params
}

func (e *Elems) GetResult(i int) any {
	return (*e)[i].Result
}

func ETHChainID(result *string) Elem {
	return Elem{Method: "eth_chainId", Result: result}
}

func ETHGasPrice(result *string) Elem {
	return Elem{Method: "eth_gasPrice", Result: result}
}

func ETHBlockNumber(result *string) Elem {
	return Elem{Method: "eth_blockNumber", Result: result}
}

func ETHGetBalance(address, block string, result *string) Elem {
	return Elem{Method: "eth_getBalance", Params: []any{address, block}, Result: result}
}

func ETHGetCode(address, block string, result *string) Elem {
	return Elem{Method: "eth_getCode", Params: []any{address, block}, Result: result}
}

func ETHCall(params any, block string, result *string) Elem {
	return Elem{Method: "eth_call", Params: []any{params, block}, Result: result}
}

func ETHEstimateGas(params any, block string, result *string) Elem {
	return Elem{Method: "eth_estimateGas", Params: []any{params, block}, Result: result}
}

func ETHGetTransactionCount(address, block string, result *string) Elem {
	return Elem{Method: "eth_getTransactionCount", Params: []any{address, block}, Result: result}
}

func ETHMaxPriorityFeePerGas(result *string) Elem {
	return Elem{Method: "eth_maxPriorityFeePerGas", Result: result}
}

func ETHSendRawTransaction(rawTx string, result *string) Elem {
	return Elem{Method: "eth_sendRawTransaction", Params: []any{rawTx}, Result: result}
}

func ETHGetBlockByNumber(block string, fullTx bool, result *map[string]any) Elem {
	return Elem{Method: "eth_getBlockByNumber", Params: []any{block, fullTx}, Result: result}
}

func ETHGetTransactionByHash(txHash string, result *map[string]any) Elem {
	return Elem{Method: "eth_getTransactionByHash", Params: []any{txHash}, Result: result}
}

func ETHGetTransactionReceipt(txHash string, result *map[string]any) Elem {
	return Elem{Method: "eth_getTransactionReceipt", Params: []any{txHash}, Result: result}
}
