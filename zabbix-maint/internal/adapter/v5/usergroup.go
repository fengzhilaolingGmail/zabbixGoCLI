package v5

import (
	"context"
	"fmt"

	"zabbix-maint/internal/api"
	"zabbix-maint/internal/model"
)

// V5Adapter 用户组相关方法
type UserGroupOps struct {
	client *api.JSONRPCClient
}

// UserGroupCreate 创建用户组
func (a *V5Adapter) UserGroupCreate(ctx context.Context, req model.UserGroupCreateReq) (string, error) {
	// TODO: implement V5 usergroup.create
	return "", fmt.Errorf("not implemented")
}

// UserGroupDisable 禁用用户组 (V5: users_status=1)
func (a *V5Adapter) UserGroupDisable(ctx context.Context, groupID string) error {
	// TODO: implement V5 usergroup.update users_status=1
	return fmt.Errorf("not implemented")
}

// UserGroupEnable 启用用户组 (V5: users_status=0)
func (a *V5Adapter) UserGroupEnable(ctx context.Context, groupID string) error {
	// TODO: implement V5 usergroup.update users_status=0
	return fmt.Errorf("not implemented")
}

// UserGroupClear 清空用户组
func (a *V5Adapter) UserGroupClear(ctx context.Context, groupID string) error {
	// TODO: implement V5 usergroup.update (清空 userids)
	return fmt.Errorf("not implemented")
}

// UserGroupDelete 删除用户组
func (a *V5Adapter) UserGroupDelete(ctx context.Context, groupID string) error {
	// TODO: implement V5 usergroup.delete
	return fmt.Errorf("not implemented")
}

// UserGroupList 查询用户组列表
func (a *V5Adapter) UserGroupList(ctx context.Context, filter string) ([]model.UnifiedUserGroup, error) {
	// TODO: implement V5 usergroup.get
	return nil, fmt.Errorf("not implemented")
}
