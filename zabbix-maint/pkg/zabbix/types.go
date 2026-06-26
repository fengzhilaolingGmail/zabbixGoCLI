package zabbix

// UserCreateReq 创建用户请求 (SDK 对外暴露)
type UserCreateReq struct {
	Alias        string
	Name         string
	Surname      string
	Passwd       string
	UserGroupIDs []string
	RoleID       string
	Media        []UserMedia
}

// UserMedia 用户告警媒介 (SDK 对外暴露)
type UserMedia struct {
	MediatypeID string
	Sendto      []string
	Severity    int64
	Period      string
}

// UnifiedUser 统一用户模型 (SDK 对外暴露)
type UnifiedUser struct {
	ID         string
	Alias      string
	Name       string
	Surname    string
	Enabled    bool
	RoleID     string
	RoleName   string
	GroupIDs   []string
	GroupNames []string
}

// UserGroupCreateReq 创建用户组请求 (SDK 对外暴露)
type UserGroupCreateReq struct {
	Name    string
	Rights  []GroupRight
	UserIDs []string
}

// GroupRight 主机组权限 (SDK 对外暴露)
type GroupRight struct {
	Permission int
	ID         string
}

// UnifiedUserGroup 统一用户组模型 (SDK 对外暴露)
type UnifiedUserGroup struct {
	ID        string
	Name      string
	Enabled   bool
	UserCount int
	UserIDs   []string
	Rights    []GroupRight
}

// UnifiedHostGroup 统一主机组模型 (SDK 对外暴露)
type UnifiedHostGroup struct {
	ID        string
	Name      string
	HostCount int
	HostIDs   []string
}

// HostCloneReq 主机克隆请求 (SDK 对外暴露)
type HostCloneReq struct {
	SourceHostID  string
	NewHostName   string
	NewHostIP     string
	NewGroupIDs   []string
	CloneItems    bool
	CloneTriggers bool
	CloneGraphs   bool
}

// UnifiedHost 统一主机模型 (SDK 对外暴露)
type UnifiedHost struct {
	ID         string
	Host       string
	Name       string
	Interfaces []HostInterface
	GroupIDs   []string
	GroupNames []string
	Templates  []string
	Macros     []HostMacro
	Inventory  map[string]string
	Status     int
}

// HostInterface 主机接口 (SDK 对外暴露)
type HostInterface struct {
	Type    int
	Main    int
	UseIP   int
	IP      string
	DNS     string
	Port    string
	Details map[string]interface{}
}

// HostMacro 主机宏 (SDK 对外暴露)
type HostMacro struct {
	Macro       string
	Value       string
	Type        int
	Description string
}

// DashboardCreateReq 创建仪表盘请求 (SDK 对外暴露)
type DashboardCreateReq struct {
	Name      string
	OwnerID   string
	ShareType int
}

// DashboardAddReq 添加 Widget 请求 (SDK 对外暴露)
type DashboardAddReq struct {
	DashboardID string
	Widgets     []WidgetDef
}

// WidgetDef Widget 定义 (SDK 对外暴露)
type WidgetDef struct {
	Type     string
	Name     string
	Position WidgetPosition
	Size     WidgetSize
	Config   map[string]interface{}
}

// WidgetPosition Widget 位置 (SDK 对外暴露)
type WidgetPosition struct {
	X, Y int
}

// WidgetSize Widget 大小 (SDK 对外暴露)
type WidgetSize struct {
	Width, Height int
}

// UnifiedDashboard 统一仪表盘模型 (SDK 对外暴露)
type UnifiedDashboard struct {
	ID        string
	Name      string
	OwnerID   string
	OwnerName string
	Pages     []DashboardPage
	Private   bool
	Type      string
}

// DashboardPage 仪表盘页面 (SDK 对外暴露)
type DashboardPage struct {
	ID      string
	Name    string
	Widgets []WidgetDef
}

// UnifiedRole 统一角色模型 (SDK 对外暴露)
type UnifiedRole struct {
	ID   string
	Name string
	Type int
}
