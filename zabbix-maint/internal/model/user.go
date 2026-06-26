package model

// UserCreateReq 创建用户请求
type UserCreateReq struct {
	Alias        string      // 用户名（必填）
	Name         string      // 姓名
	Surname      string      // 姓氏
	Passwd       string      // 初始密码
	UserGroupIDs []string    // 所属用户组
	RoleID       string      // 角色ID（7.x 必填，5.x 忽略）
	Media        []UserMedia // 告警媒介
}

// UserMedia 用户告警媒介
type UserMedia struct {
	MediatypeID string
	Sendto      []string
	Severity    int64
	Period      string
}

// UnifiedUser 统一用户模型
type UnifiedUser struct {
	ID         string
	Alias      string
	Name       string
	Surname    string
	Enabled    bool
	RoleID     string // 7.x 使用，5.x 忽略
	RoleName   string
	GroupIDs   []string
	GroupNames []string
}
