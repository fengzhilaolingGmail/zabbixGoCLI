package model

// UserGroupCreateReq 创建用户组请求
type UserGroupCreateReq struct {
	Name    string
	Rights  []GroupRight // 主机组权限
	UserIDs []string     // 初始用户列表
}

// GroupRight 主机组权限
type GroupRight struct {
	Permission int    // 2=只读, 3=读写
	ID         string // hostgroupid
}

// UnifiedUserGroup 统一用户组模型
type UnifiedUserGroup struct {
	ID        string
	Name      string
	Enabled   bool
	UserCount int
	UserIDs   []string
	Rights    []GroupRight
}

// UnifiedHostGroup 统一主机组模型
type UnifiedHostGroup struct {
	ID        string
	Name      string
	HostCount int
	HostIDs   []string
}

// UnifiedRole 统一角色模型
type UnifiedRole struct {
	ID   string
	Name string
	Type int // 1=用户, 2=管理员, 3=超级管理员
}

// WidgetField Widget 字段 (V7)
type WidgetField struct {
	Type  string
	Value string
}

// DashboardWidget 仪表盘 Widget (V7)
type DashboardWidget struct {
	Type   string
	X      int
	Y      int
	Width  int
	Height int
	Fields []WidgetField
}
