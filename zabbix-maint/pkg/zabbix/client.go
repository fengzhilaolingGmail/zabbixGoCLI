package zabbix

import (
	"context"
	"fmt"

	"zabbix-maint/internal/api"
)

// Client Zabbix 客户端封装
type Client struct {
	operator ZabbixOperator
}

// NewClient 创建 Zabbix 客户端
func NewClient(operator ZabbixOperator) *Client {
	return &Client{operator: operator}
}

// WithAuth 返回带认证的新客户端
func (c *Client) WithAuth(endpoint, username, password string) (*Client, error) {
	// TODO: implement authenticated client creation
	return nil, fmt.Errorf("not implemented")
}

// Operator 获取底层操作接口
func (c *Client) Operator() ZabbixOperator {
	return c.operator
}

// Close 关闭客户端连接
func (c *Client) Close() error {
	return nil
}
