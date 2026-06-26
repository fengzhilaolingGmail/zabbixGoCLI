package v5

import (
	"context"
	"fmt"

	"zabbix-maint/pkg/zabbix"
)

// UserCreate 创建用户 (V5 版本)
func (a *V5Adapter) UserCreate(ctx context.Context, req zabbix.UserCreateReq) (string, error) {
	// TODO: implement V5 user.create
	return "", fmt.Errorf("not implemented")
}

// UserDisable 禁用用户 (V5: status=1)
func (a *V5Adapter) UserDisable(ctx context.Context, userID string) error {
	// TODO: implement V5 user.update status=1
	return fmt.Errorf("not implemented")
}

// UserEnable 启用用户 (V5: status=0)
func (a *V5Adapter) UserEnable(ctx context.Context, userID string) error {
	// TODO: implement V5 user.update status=0
	return fmt.Errorf("not implemented")
}

// UserUpdatePassword 修改密码
func (a *V5Adapter) UserUpdatePassword(ctx context.Context, userID, newPass string) error {
	// TODO: implement V5 user.update passwd
	return fmt.Errorf("not implemented")
}

// UserDelete 删除用户
func (a *V5Adapter) UserDelete(ctx context.Context, userID string) error {
	// TODO: implement V5 user.delete
	return fmt.Errorf("not implemented")
}

// UserList 查询用户列表
func (a *V5Adapter) UserList(ctx context.Context, filter string) ([]zabbix.UnifiedUser, error) {
	// TODO: implement V5 user.get
	return nil, fmt.Errorf("not implemented")
}

// UserDetail 查询用户详情
func (a *V5Adapter) UserDetail(ctx context.Context, userID string) (*zabbix.UnifiedUser, error) {
	// TODO: implement V5 user.get
	return nil, fmt.Errorf("not implemented")
}
