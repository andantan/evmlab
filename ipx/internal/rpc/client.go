package rpc

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"sync/atomic"
	"time"
)

type Client struct {
	url        string
	httpClient *http.Client
	nextID     uint64
}

type Request struct {
	JSONRPC string `json:"jsonrpc"`
	ID      uint64 `json:"id"`
	Method  string `json:"method"`
	Params  any    `json:"params"`
}

type Response[T any] struct {
	JSONRPC string    `json:"jsonrpc"`
	ID      uint64    `json:"id"`
	Result  T         `json:"result"`
	Error   *RPCError `json:"error,omitempty"`
}

type RPCError struct {
	Code    int             `json:"code"`
	Message string          `json:"message"`
	Data    json.RawMessage `json:"data,omitempty"`
}

func (e *RPCError) Error() string {
	if len(e.Data) == 0 {
		return fmt.Sprintf("rpc error %d: %s", e.Code, e.Message)
	}

	return fmt.Sprintf("rpc error %d: %s: %s", e.Code, e.Message, string(e.Data))
}

func NewClient(url string) *Client {
	return &Client{
		url: url,
		httpClient: &http.Client{
			Timeout: time.Second * 15,
		},
	}
}

func (c *Client) Batch(ctx context.Context, elems *Elems) error {
	n := elems.Len()
	reqs := make([]Request, n)
	ids := make([]uint64, n)
	for i := range n {
		id := atomic.AddUint64(&c.nextID, 1)
		ids[i] = id
		reqs[i] = Request{JSONRPC: "2.0", ID: id, Method: elems.GetMethod(i), Params: elems.GetParams(i)}
	}

	bodyBytes, err := json.Marshal(reqs)
	if err != nil {
		return fmt.Errorf("marshal batch request: %w", err)
	}

	httpReq, err := http.NewRequestWithContext(ctx, http.MethodPost, c.url, bytes.NewReader(bodyBytes))
	if err != nil {
		return fmt.Errorf("create http request: %w", err)
	}

	httpReq.Header.Set("Content-Type", "application/json")
	httpResp, err := c.httpClient.Do(httpReq)
	if err != nil {
		return fmt.Errorf("send batch request: %w", err)
	}
	defer httpResp.Body.Close()

	respBytes, err := io.ReadAll(httpResp.Body)
	if err != nil {
		return fmt.Errorf("read batch response: %w", err)
	}

	if httpResp.StatusCode < 200 || httpResp.StatusCode >= 300 {
		return fmt.Errorf("http status %d: %s", httpResp.StatusCode, string(respBytes))
	}

	var responses []Response[json.RawMessage]
	if err = json.Unmarshal(respBytes, &responses); err != nil {
		return fmt.Errorf("decode batch response: %w: %s", err, string(respBytes))
	}

	byID := make(map[uint64]Response[json.RawMessage], len(responses))
	for _, resp := range responses {
		byID[resp.ID] = resp
	}

	for i, id := range ids {
		resp, ok := byID[id]
		if !ok {
			return fmt.Errorf("missing response for %s", elems.GetMethod(i))
		}
		if resp.Error != nil {
			return fmt.Errorf("%s: %w", elems.GetMethod(i), resp.Error)
		}
		if elems.GetResult(i) == nil {
			continue
		}
		if err = json.Unmarshal(resp.Result, elems.GetResult(i)); err != nil {
			return fmt.Errorf("decode result for %s: %w: %s", elems.GetMethod(i), err, string(resp.Result))
		}
	}

	return nil
}

func (c *Client) Call(ctx context.Context, e Elem) error {
	id := atomic.AddUint64(&c.nextID, 1)
	reqBody := Request{
		JSONRPC: "2.0",
		ID:      id,
		Method:  e.Method,
		Params:  e.Params,
	}

	bodyBytes, err := json.Marshal(reqBody)
	if err != nil {
		return fmt.Errorf("marshal rpc request: %w", err)
	}

	httpReq, err := http.NewRequestWithContext(ctx, http.MethodPost, c.url, bytes.NewReader(bodyBytes))
	if err != nil {
		return fmt.Errorf("create http request: %w", err)
	}

	httpReq.Header.Set("Content-Type", "application/json")
	httpResp, err := c.httpClient.Do(httpReq)
	if err != nil {
		return fmt.Errorf("send rpc request: %w", err)
	}
	defer httpResp.Body.Close()

	respBytes, err := io.ReadAll(httpResp.Body)
	if err != nil {
		return fmt.Errorf("read rpc response: %w", err)
	}

	if httpResp.StatusCode < 200 || httpResp.StatusCode >= 300 {
		return fmt.Errorf("http status %d: %s", httpResp.StatusCode, string(respBytes))
	}

	var rpcResp Response[json.RawMessage]
	if err = json.Unmarshal(respBytes, &rpcResp); err != nil {
		return fmt.Errorf("decode rpc response: %w: %s", err, string(respBytes))
	}

	if rpcResp.Error != nil {
		return rpcResp.Error
	}

	if e.Result == nil {
		return nil
	}

	if err = json.Unmarshal(rpcResp.Result, e.Result); err != nil {
		return fmt.Errorf("decode rpc result for %s: %w: %s", e.Method, err, string(rpcResp.Result))
	}

	return nil
}
