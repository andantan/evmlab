package rpc

import (
	"context"
	"fmt"
	"time"
)

func (c *Client) ChainID(ctx context.Context) (string, error) {
	var result string
	err := c.Call(ctx, ETHChainID(&result))
	return result, err
}

func (c *Client) GasPrice(ctx context.Context) (string, error) {
	var result string
	err := c.Call(ctx, ETHGasPrice(&result))
	return result, err
}

func (c *Client) BlockNumber(ctx context.Context) (string, error) {
	var result string
	err := c.Call(ctx, ETHBlockNumber(&result))
	return result, err
}

func (c *Client) GetBalance(ctx context.Context, address string, block string) (string, error) {
	var result string
	err := c.Call(ctx, ETHGetBalance(address, block, &result))
	return result, err
}

func (c *Client) GetCode(ctx context.Context, address string, block string) (string, error) {
	var result string
	err := c.Call(ctx, ETHGetCode(address, block, &result))
	return result, err
}

func (c *Client) IsContract(ctx context.Context, address string, block string) (bool, error) {
	r, err := c.GetCode(ctx, address, block)
	if err != nil {
		return false, err
	}

	// EIP-7702: 0xef0100 + 20-byte address
	if len(r) == 48 && r[:8] == "0xef0100" {
		return false, nil
	}

	return r != "0x", nil
}

func (c *Client) CallContract(ctx context.Context, params any, block string) (string, error) {
	var result string
	err := c.Call(ctx, ETHCall(params, block, &result))
	return result, err
}

func (c *Client) EstimateGas(ctx context.Context, params any, block string) (string, error) {
	var result string
	err := c.Call(ctx, ETHEstimateGas(params, block, &result))
	return result, err
}

func (c *Client) GetTransactionCount(ctx context.Context, address string, block string) (string, error) {
	var result string
	err := c.Call(ctx, ETHGetTransactionCount(address, block, &result))
	return result, err
}

func (c *Client) MaxPriorityFeePerGas(ctx context.Context) (string, error) {
	var result string
	err := c.Call(ctx, ETHMaxPriorityFeePerGas(&result))
	return result, err
}

func (c *Client) SendRawTransaction(ctx context.Context, rawTx string) (string, error) {
	var result string
	err := c.Call(ctx, ETHSendRawTransaction(rawTx, &result))
	return result, err
}

func (c *Client) BlockByNumber(ctx context.Context, block string) (map[string]any, error) {
	var result map[string]any
	err := c.Call(ctx, ETHGetBlockByNumber(block, false, &result))
	return result, err
}

func (c *Client) BaseFeePerGas(ctx context.Context) (string, error) {
	var result map[string]any
	if err := c.Call(ctx, ETHGetBlockByNumber("latest", false, &result)); err != nil {
		return "", err
	}

	baseFee, ok := result["baseFeePerGas"].(string)
	if !ok || baseFee == "" {
		return "", fmt.Errorf("baseFeePerGas not found in latest block")
	}

	return baseFee, nil
}

func (c *Client) GetTransactionByHash(ctx context.Context, txHash string) (map[string]any, error) {
	var result map[string]any
	err := c.Call(ctx, ETHGetTransactionByHash(txHash, &result))
	return result, err
}

func (c *Client) TransactionReceipt(ctx context.Context, txHash string) (map[string]any, error) {
	var result map[string]any
	err := c.Call(ctx, ETHGetTransactionReceipt(txHash, &result))
	return result, err
}

func (c *Client) WaitForReceipt(ctx context.Context, txHash string, timeout time.Duration) (map[string]any, error) {
	deadline := time.Now().Add(timeout)

	for {
		receipt, err := c.TransactionReceipt(ctx, txHash)
		if err == nil && receipt != nil {
			return receipt, nil
		}

		if time.Now().After(deadline) {
			return nil, fmt.Errorf("timeout waiting for receipt: %s", txHash)
		}

		select {
		case <-ctx.Done():
			return nil, ctx.Err()
		case <-time.After(500 * time.Millisecond):
		}
	}
}
