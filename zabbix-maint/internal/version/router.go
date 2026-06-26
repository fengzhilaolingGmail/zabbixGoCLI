package version

import (
	"context"
	"fmt"

	"zabbix-maint/internal/api"
	"zabbix-maint/internal/config"
	"zabbix-maint/pkg/zabbix"
)

// AdapterFactory 适配器工厂
type AdapterFactory interface {
	Create(version Version, endpoint string, auth config.AuthConfig) (zabbix.ZabbixOperator, error)
}

// DefaultAdapterFactory 默认适配器工厂
type DefaultAdapterFactory struct{}

// NewAdapterFactory 创建适配器工厂
func NewAdapterFactory() AdapterFactory {
	return &DefaultAdapterFactory{}
}

// Create 创建适配器
func (f *DefaultAdapterFactory) Create(ver Version, endpoint string, auth config.AuthConfig) (zabbix.ZabbixOperator, error) {
	// TODO: implement adapter creation
	return nil, fmt.Errorf("adapter factory not implemented for version %s", ver)
}

// Router 版本路由
type Router struct {
	factory AdapterFactory
}

// NewRouter 创建版本路由器
func NewRouter() *Router {
	return &Router{factory: NewAdapterFactory()}
}

// Connect 连接并返回适配器
func (r *Router) Connect(ctx context.Context, cfg *config.InstanceConfig) (zabbix.ZabbixOperator, Version, error) {
	client := api.NewJSONRPCClient(cfg.Endpoint)

	ver, err := DetectVersion(ctx, client)
	if err != nil {
		return nil, "", fmt.Errorf("version detection failed: %w", err)
	}

	// 如果配置中强制指定了版本，使用配置版本
	if cfg.Version != "" {
		ver = Version(cfg.Version)
	}

	adapter, err := r.factory.Create(ver, cfg.Endpoint, cfg.Auth)
	if err != nil {
		return nil, "", err
	}

	return adapter, ver, nil
}
