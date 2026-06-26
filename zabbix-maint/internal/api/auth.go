package api

import (
	"context"
	"fmt"
	"time"

	"zabbix-maint/internal/log"
)

// AuthManager 认证管理器
type AuthManager struct {
	username string
	password string
	token    string
	expiry   time.Time
	client   *JSONRPCClient
}

// NewAuthManager 创建认证管理器
func NewAuthManager(username, password string, client *JSONRPCClient) *AuthManager {
	return &AuthManager{
		username: username,
		password: password,
		client:   client,
	}
}

// Refresh 刷新认证令牌
func (a *AuthManager) Refresh(ctx context.Context) error {
	var token string
	if err := a.client.Call(ctx, "user.login", map[string]interface{}{
		"user":     a.username,
		"password": a.password,
	}, &token); err != nil {
		log.Errorf("Login failed: %v", err)
		return fmt.Errorf("login failed: %w", err)
	}

	log.Infof("Login succeeded, token obtained")
	a.token = token
	a.expiry = time.Now().Add(15 * time.Minute) // Zabbix token 默认 15 分钟
	a.client.SetAuthToken(token)
	return nil
}

// Token 获取当前认证令牌
func (a *AuthManager) Token() string {
	return a.token
}

// IsExpired 检查令牌是否过期
func (a *AuthManager) IsExpired() bool {
	return time.Now().After(a.expiry)
}
