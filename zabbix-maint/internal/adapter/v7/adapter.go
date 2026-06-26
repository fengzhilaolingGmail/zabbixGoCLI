package v7

import (
	"context"
	"fmt"

	"zabbix-maint/internal/api"
	"zabbix-maint/internal/model"
	"zabbix-maint/pkg/zabbix"
)

// V7Adapter 实现 zabbix.ZabbixOperator 接口 (V7 版本)
type V7Adapter struct {
	client *api.JSONRPCClient
	auth   *api.AuthManager
}

// NewV7Adapter 创建 V7 适配器
func NewV7Adapter(client *api.JSONRPCClient, auth *api.AuthManager) zabbix.ZabbixOperator {
	return &V7Adapter{client: client, auth: auth}
}

// APIVersion 获取 API 版本
func (a *V7Adapter) APIVersion(ctx context.Context) (string, error) {
	var version string
	err := a.client.Call(ctx, "apiinfo.version", nil, &version)
	return version, err
}

// ServerStatus 获取服务器状态
func (a *V7Adapter) ServerStatus(ctx context.Context) (map[string]interface{}, error) {
	// TODO: implement V7 status.get
	return nil, fmt.Errorf("not implemented")
}

// RoleList 查询角色列表 (V7 专属)
func (a *V7Adapter) RoleList(ctx context.Context, filter string) ([]model.UnifiedRole, error) {
	// TODO: implement V7 role.get
	return nil, fmt.Errorf("not implemented")
}
