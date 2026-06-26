package v7

import (
	"context"
	"fmt"

	"zabbix-maint/internal/api"
	"zabbix-maint/internal/model"
)

// V7Adapter 实现 Zabbix 7.x 版本的适配器
type V7Adapter struct {
	client *api.JSONRPCClient
	auth   *api.AuthManager
}

// NewV7Adapter 创建 V7 适配器
func NewV7Adapter(client *api.JSONRPCClient, auth *api.AuthManager) *V7Adapter {
	return &V7Adapter{client: client, auth: auth}
}

// UserCreate 创建用户 (V7 版本，需要 roleid)
func (a *V7Adapter) UserCreate(ctx context.Context, req model.UserCreateReq) (string, error) {
	// TODO: implement V7 user.create with roleid
	return "", fmt.Errorf("not implemented")
}

// UserDisable 禁用用户 (V7: roleid=<disabled_role_id>)
func (a *V7Adapter) UserDisable(ctx context.Context, userID string) error {
	// TODO: implement V7 user.update with disabled role
	return fmt.Errorf("not implemented")
}

// UserEnable 启用用户 (V7: roleid=<default_user_role_id>)
func (a *V7Adapter) UserEnable(ctx context.Context, userID string) error {
	// TODO: implement V7 user.update with default role
	return fmt.Errorf("not implemented")
}

// UserUpdatePassword 修改密码
func (a *V7Adapter) UserUpdatePassword(ctx context.Context, userID, newPass string) error {
	// TODO: implement V7 user.update passwd
	return fmt.Errorf("not implemented")
}

// UserDelete 删除用户
func (a *V7Adapter) UserDelete(ctx context.Context, userID string) error {
	// TODO: implement V7 user.delete
	return fmt.Errorf("not implemented")
}

// UserList 查询用户列表
func (a *V7Adapter) UserList(ctx context.Context, filter string) ([]model.UnifiedUser, error) {
	// TODO: implement V7 user.get
	return nil, fmt.Errorf("not implemented")
}

// UserDetail 查询用户详情
func (a *V7Adapter) UserDetail(ctx context.Context, userID string) (*model.UnifiedUser, error) {
	// TODO: implement V7 user.get
	return nil, fmt.Errorf("not implemented")
}
