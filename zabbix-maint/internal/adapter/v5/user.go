package v5

import (
	"context"
	"fmt"

	"zabbix-maint/pkg/zabbix"
)

// UserCreate 创建用户 (V5)
func (a *V5Adapter) UserCreate(ctx context.Context, req zabbix.UserCreateReq) (string, error) {
	usrgrps := make([]map[string]string, len(req.UserGroupIDs))
	for i, gid := range req.UserGroupIDs {
		usrgrps[i] = map[string]string{"usrgrpid": gid}
	}
	params := map[string]interface{}{
		"alias":   req.Alias,
		"name":    req.Name,
		"surname": req.Surname,
		"passwd":  req.Passwd,
		"usrgrps": usrgrps,
		"type":    "1",
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

// UserDisable 禁用用户 (V5: status=1)
func (a *V5Adapter) UserDisable(ctx context.Context, userID string) error {
	params := map[string]interface{}{
		"userid": userID,
		"status": "1",
	}
	return a.client.Call(ctx, "user.update", params, nil)
}

// UserEnable 启用用户 (V5: status=0)
func (a *V5Adapter) UserEnable(ctx context.Context, userID string) error {
	params := map[string]interface{}{
		"userid": userID,
		"status": "0",
	}
	return a.client.Call(ctx, "user.update", params, nil)
}

// UserUpdatePassword 修改密码
func (a *V5Adapter) UserUpdatePassword(ctx context.Context, userID, newPass string) error {
	params := map[string]interface{}{
		"userid": userID,
		"passwd": newPass,
	}
	return a.client.Call(ctx, "user.update", params, nil)
}

// UserDelete 删除用户
func (a *V5Adapter) UserDelete(ctx context.Context, userID string) error {
	return a.client.Call(ctx, "user.delete", []string{userID}, nil)
}

// UserList 查询用户列表
func (a *V5Adapter) UserList(ctx context.Context, filter string) ([]zabbix.UnifiedUser, error) {
	params := map[string]interface{}{
		"output":        "extend",
		"selectUsrgrps": "extend",
	}
	if filter != "" {
		params["search"] = map[string]string{"alias": filter}
	}
	var users []rawUserV5
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
func (a *V5Adapter) UserDetail(ctx context.Context, userID string) (*zabbix.UnifiedUser, error) {
	params := map[string]interface{}{
		"userids":       userID,
		"output":        "extend",
		"selectUsrgrps": "extend",
	}
	var users []rawUserV5
	if err := a.client.Call(ctx, "user.get", params, &users); err != nil {
		return nil, err
	}
	if len(users) == 0 {
		return nil, fmt.Errorf("user not found: %s", userID)
	}
	u := users[0].toUnified()
	return &u, nil
}

// rawUserV5 Zabbix 5.x 原始用户返回
type rawUserV5 struct {
	UserID  string `json:"userid"`
	Alias   string `json:"alias"`
	Name    string `json:"name"`
	Surname string `json:"surname"`
	Status  string `json:"status"`
	Usrgrps []struct {
		Usrgrpid string `json:"usrgrpid"`
		Name     string `json:"name"`
	} `json:"usrgrps"`
}

func (u rawUserV5) toUnified() zabbix.UnifiedUser {
	enabled := u.Status == "0"
	groupIDs := make([]string, len(u.Usrgrps))
	groupNames := make([]string, len(u.Usrgrps))
	for i, g := range u.Usrgrps {
		groupIDs[i] = g.Usrgrpid
		groupNames[i] = g.Name
	}
	return zabbix.UnifiedUser{
		ID:         u.UserID,
		Alias:      u.Alias,
		Name:       u.Name,
		Surname:    u.Surname,
		Enabled:    enabled,
		GroupIDs:   groupIDs,
		GroupNames: groupNames,
	}
}
