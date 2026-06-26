package model

// DashboardCreateReq 创建仪表盘/聚合图形请求
type DashboardCreateReq struct {
	Name      string
	OwnerID   string // 所属用户ID
	ShareType int    // 0=私有, 1=公开
}

// DashboardAddReq 添加 Widget/子项请求
type DashboardAddReq struct {
	DashboardID string
	Widgets     []WidgetDef
}

// WidgetDef Widget 定义
type WidgetDef struct {
	Type     string
	Name     string
	Position WidgetPosition
	Size     WidgetSize
	Config   map[string]interface{} // 版本相关配置
}

// WidgetPosition Widget 位置
type WidgetPosition struct {
	X, Y int
}

// WidgetSize Widget 大小
type WidgetSize struct {
	Width, Height int
}

// DashboardPage 仪表盘页面
type DashboardPage struct {
	ID      string
	Name    string
	Widgets []WidgetDef
}

// UnifiedDashboard 统一仪表盘/聚合图形模型
type UnifiedDashboard struct {
	ID        string
	Name      string
	OwnerID   string
	OwnerName string
	Pages     []DashboardPage
	Private   bool
	Type      string // "screen"(5.x) / "dashboard"(7.x)
}
