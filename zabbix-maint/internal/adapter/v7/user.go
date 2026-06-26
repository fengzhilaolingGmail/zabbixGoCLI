package v7

import (
	"context"
	"fmt"

	"zabbix-maint/pkg/zabbix"
)

// UserCreate 创建用户 (V7)
func (a *V7Adapter) UserCreate(ctx context.Context, req zabbix.UserCreateReq) (string, error) {
	usrgrps := make([]map[string]string, len(req.UserGroupIDs))
	for i, gid := range req.UserGroupIDs {
		usrgrps[i] = map[string]string{"usrgrpid": gid}
	}
	params := map[string]interface{}{
		"username": req.Alias,
		"name":     req.Name,
		"surname":  req.Surname,
		"passwd":   req.Passwd,
		"roleid":   req.RoleID,
		"usrgrps":  usrgrps,
	}
	var result struct {
		UserIDs []string `json:"userids"`
	}
	if err := a.client.Call(ctx, "user.create", params, &result); err != nil {
		return "", err
	}
	if len(result.UserIDs) == 0 {
		return "", fmt.Errorf("user.create returned empty")
	}
	return result.UserIDs[0], nil
}

// UserDisable 禁用用户 (V7: 设为 Guest role=1)
func (a *V7Adapter) UserDisable(ctx context.Context, userID string) error {
	params := map[string]interface{}{
		"userid": userID,
		"roleid": "1",
	}
	return a.client.Call(ctx, "user.update", params, nil)
}

// UserEnable 启用用户 (V7: 恢复为 User role=2)
func (a *V7Adapter) UserEnable(ctx context.Context, userID string) error {
	params := map[string]interface{}{
		"userid": userID,
		"roleid": "2",
	}
	return a.client.Call(ctx, "user.update", params, nil)
}

// UserUpdatePassword 修改密码
func (a *V7Adapter) UserUpdatePassword(ctx context.Context, userID, newPass string) error {
	params := map[string]interface{}{
		"userid": userID,
		"passwd": newPass,
	}
	return a.client.Call(ctx, "user.update", params, nil)
}

// UserDelete 删除用户
func (a *V7Adapter) UserDelete(ctx context.Context, userID string) error {
	return a.client.Call(ctx, "user.delete", []string{userID}, nil)
}

// UserList 查询用户列表
func (a *V7Adapter) UserList(ctx context.Context, filter string) ([]zabbix.UnifiedUser, error) {
	params := map[string]interface{}{
		"output":        "extend",
		"selectUsrgrps": "extend",
		"selectRole":    "extend",
	}
	if filter != "" {
		params["search"] = map[string]string{"username": filter}
	}
	var users []rawUserV7
	if err := a.client.Call(ctx, "user.get", params, &users); err != nil {
		return nil, err
	}
	result := make([]zabbix.UnifiedUser, len(users))
	for i, u := range users {
		result[i] = u.toUnified()
	}
	return result, nil
}

// UserDetail 查询用户详情
func (a *V7Adapter) UserDetail(ctx context.Context, userID string) (*zabbix.UnifiedUser, error) {
	params := map[string]interface{}{
		"userids":       userID,
		"output":        "extend",
		"selectUsrgrps": "extend",
		"selectRole":    "extend",
	}
	var users []rawUserV7
	if err := a.client.Call(ctx, "user.get", params, &users); err != nil {
		return nil, err
	}
	if len(users) == 0 {
		return nil, fmt.Errorf("user not found: %s", userID)
	}
	u := users[0].toUnified()
	return &u, nil
}

// rawUserV7 Zabbix 7.x 原始用户返回
type rawUserV7 struct {
	UserID  string `json:"userid"`
	Username string `json:"username"`
	Name    string `json:"name"`
	Surname string `json:"surname"`
	RoleID  string `json:"roleid"`
	Role    struct {
		RoleID string `json:"roleid"`
		Name   string `json:"name"`
	} `json:"role"`
	Usrgrps []struct {
		Usrgrpid string `json:"usrgrpid"`
		Name     string `json:"name"`
	} `json:"usrgrps"`
}

func (u rawUserV7) toUnified() zabbix.UnifiedUser {
	groupIDs := make([]string, len(u.Usrgrps))
	groupNames := make([]string, len(u.Usrgrps))
	for i, g := range u.Usrgrps {
		groupIDs[i] = g.Usrgrpid
		groupNames[i] = g.Name
	}
	// V7 中 roleid=1 通常对应 Guest（受限），视为禁用
	enabled := u.RoleID != "1"
	return zabbix.UnifiedUser{
		ID:         u.UserID,
		Alias:      u.Username,
		Name:       u.Name,
		Surname:    u.Surname,
		Enabled:    enabled,
		RoleID:     u.RoleID,
		RoleName:   u.Role.Name,
		GroupIDs:   groupIDs,
		GroupNames: groupNames,
	}
}
