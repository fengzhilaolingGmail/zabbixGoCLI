package v7

import (
	"context"
	"fmt"

	"zabbix-maint/internal/api"
	"zabbix-maint/internal/model"
)

// DashboardCreate 创建仪表盘 (V7: dashboard.create)
func (a *V7Adapter) DashboardCreate(ctx context.Context, req model.DashboardCreateReq) (string, error) {
	// TODO: implement V7 dashboard.create
	return "", fmt.Errorf("not implemented")
}

// DashboardAddWidget 添加 Widget (V7: dashboard.widget.create)
func (a *V7Adapter) DashboardAddWidget(ctx context.Context, req model.DashboardAddReq) (string, error) {
	// TODO: implement V7 dashboard.widget.create
	return "", fmt.Errorf("not implemented")
}

// DashboardList 查询仪表盘列表
func (a *V7Adapter) DashboardList(ctx context.Context, filter string) ([]model.UnifiedDashboard, error) {
	// TODO: implement V7 dashboard.get
	return nil, fmt.Errorf("not implemented")
}
