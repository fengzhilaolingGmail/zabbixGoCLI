package v7

import (
	"context"
	"fmt"

	"zabbix-maint/internal/api"
	"zabbix-maint/internal/model"
)

// HostGroupCreate 创建主机组
func (a *V7Adapter) HostGroupCreate(ctx context.Context, name string) (string, error) {
	// TODO: implement V7 hostgroup.create
	return "", fmt.Errorf("not implemented")
}

// HostGroupClear 清空主机组 (V7: hostgroup.update 清空 hostids)
func (a *V7Adapter) HostGroupClear(ctx context.Context, groupID string) error {
	// TODO: implement V7 hostgroup.update with empty hosts
	return fmt.Errorf("not implemented")
}

// HostGroupDelete 删除主机组
func (a *V7Adapter) HostGroupDelete(ctx context.Context, groupID string) error {
	// TODO: implement V7 hostgroup.delete
	return fmt.Errorf("not implemented")
}

// HostGroupList 查询主机组列表
func (a *V7Adapter) HostGroupList(ctx context.Context, filter string) ([]model.UnifiedHostGroup, error) {
	// TODO: implement V7 hostgroup.get
	return nil, fmt.Errorf("not implemented")
}
