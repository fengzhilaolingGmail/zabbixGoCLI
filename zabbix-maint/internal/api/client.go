package api

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"zabbix-maint/internal/log"
)

// JSONRPCClient JSON-RPC 客户端
type JSONRPCClient struct {
	endpoint       string
	authToken      string
	httpClient     *http.Client
	retryCount     int
	retryBackoff   time.Duration
	requestTimeout time.Duration
}

// NewJSONRPCClient 创建 JSON-RPC 客户端
func NewJSONRPCClient(endpoint string) *JSONRPCClient {
	return &JSONRPCClient{
		endpoint:       endpoint,
		httpClient:     &http.Client{Timeout: 30 * time.Second},
		retryCount:     3,
		retryBackoff:   500 * time.Millisecond,
		requestTimeout: 30 * time.Second,
	}
}

// Call 执行 JSON-RPC 调用
func (c *JSONRPCClient) Call(ctx context.Context, method string, params interface{}, result interface{}) error {
	// 修复：Zabbix JSON-RPC 不接受 params 为 null，转为空对象
	if params == nil {
		params = map[string]interface{}{}
	}

	reqBody := map[string]interface{}{
		"jsonrpc": "2.0",
		"method":  method,
		"params":  params,
		"id":      1,
	}
	if c.authToken != "" {
		reqBody["auth"] = c.authToken
	}

	body, err := json.Marshal(reqBody)
	if err != nil {
		return fmt.Errorf("marshal request failed: %w", err)
	}

	log.Debugf("JSON-RPC Request -> %s %s", method, string(body))

	req, err := http.NewRequestWithContext(ctx, "POST", c.endpoint, bytes.NewReader(body))
	if err != nil {
		return fmt.Errorf("create request failed: %w", err)
	}
	req.Header.Set("Content-Type", "application/json")

	resp, err := c.httpClient.Do(req)
	if err != nil {
		log.Errorf("HTTP request failed: %v", err)
		return fmt.Errorf("http request failed: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		log.Errorf("Unexpected status code: %d", resp.StatusCode)
		return fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}

	var rpcResp struct {
		JSONRPC string          `json:"jsonrpc"`
		Result  json.RawMessage `json:"result"`
		Error   *struct {
			Code    int    `json:"code"`
			Message string `json:"message"`
			Data    string `json:"data"`
		} `json:"error"`
		ID int `json:"id"`
	}

	if err := json.NewDecoder(resp.Body).Decode(&rpcResp); err != nil {
		log.Errorf("Decode response failed: %v", err)
		return fmt.Errorf("decode response failed: %w", err)
	}

	if rpcResp.Error != nil {
		log.Errorf("Zabbix error: %s (code: %d)", rpcResp.Error.Message, rpcResp.Error.Code)
		return fmt.Errorf("zabbix error: %s (code: %d)", rpcResp.Error.Message, rpcResp.Error.Code)
	}

	log.Debugf("JSON-RPC Response <- result: %s", string(rpcResp.Result))

	if result != nil && rpcResp.Result != nil {
		if err := json.Unmarshal(rpcResp.Result, result); err != nil {
			log.Errorf("Unmarshal result failed: %v", err)
			return fmt.Errorf("unmarshal result failed: %w", err)
		}
	}

	return nil
}

// SetAuthToken 设置认证令牌
func (c *JSONRPCClient) SetAuthToken(token string) {
	c.authToken = token
}
