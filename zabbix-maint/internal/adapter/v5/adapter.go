package v5

import (
	"context"
	"fmt"

	"zabbix-maint/internal/api"
	"zabbix-maint/pkg/zabbix"
)

// V5Adapter 实现 zabbix.ZabbixOperator 接口 (V5 版本)
type V5Adapter struct {
	client *api.JSONRPCClient
	auth   *api.AuthManager
}

// NewV5Adapter 创建 V5 适配器
func NewV5Adapter(client *api.JSONRPCClient, auth *api.AuthManager) zabbix.ZabbixOperator {
	return &V5Adapter{client: client, auth: auth}
}

// APIVersion 获取 API 版本
func (a *V5Adapter) APIVersion(ctx context.Context) (string, error) {
	var version string
	err := a.client.Call(ctx, "apiinfo.version", nil, &version)
	return version, err
}

// ServerStatus 获取服务器状态 (V5: status.get)
func (a *V5Adapter) ServerStatus(ctx context.Context) (map[string]interface{}, error) {
	var status map[string]interface{}
	err := a.client.Call(ctx, "status.get", nil, &status)
	return status, err
}

// RoleList V5 不支持角色管理
func (a *V5Adapter) RoleList(ctx context.Context, filter string) ([]zabbix.UnifiedRole, error) {
	return nil, fmt.Errorf("role management not supported in V5")
}

func deepCopyMap(m map[string]interface{}) map[string]interface{} {
	copy := make(map[string]interface{}, len(m))
	for k, v := range m {
		copy[k] = v
	}
	return copy
}
