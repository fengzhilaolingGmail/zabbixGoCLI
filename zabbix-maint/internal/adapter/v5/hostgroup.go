package v5

import (
	"context"
	"fmt"

	"zabbix-maint/pkg/zabbix"
)

// HostGroupCreate 创建主机组
func (a *V5Adapter) HostGroupCreate(ctx context.Context, name string) (string, error) {
	// TODO: implement V5 hostgroup.create
	return "", fmt.Errorf("not implemented")
}

// HostGroupClear 清空主机组 (V5: hostgroup.massremove)
func (a *V5Adapter) HostGroupClear(ctx context.Context, groupID string) error {
	// TODO: implement V5 hostgroup.massremove
	return fmt.Errorf("not implemented")
}

// HostGroupDelete 删除主机组
func (a *V5Adapter) HostGroupDelete(ctx context.Context, groupID string) error {
	// TODO: implement V5 hostgroup.delete
	return fmt.Errorf("not implemented")
}

// HostGroupList 查询主机组列表
func (a *V5Adapter) HostGroupList(ctx context.Context, filter string) ([]zabbix.UnifiedHostGroup, error) {
	// TODO: implement V5 hostgroup.get
	return nil, fmt.Errorf("not implemented")
}
