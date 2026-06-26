package v5

import (
	"context"
	"fmt"

	"zabbix-maint/pkg/zabbix"
)

// DashboardCreate 创建聚合图形 (V5: screen.create)
func (a *V5Adapter) DashboardCreate(ctx context.Context, req zabbix.DashboardCreateReq) (string, error) {
	// TODO: implement V5 screen.create
	return "", fmt.Errorf("not implemented")
}

// DashboardAddWidget 添加聚合图形子项 (V5: screenitem.create)
func (a *V5Adapter) DashboardAddWidget(ctx context.Context, req zabbix.DashboardAddReq) (string, error) {
	// TODO: implement V5 screenitem.create
	return "", fmt.Errorf("not implemented")
}

// DashboardList 查询聚合图形列表
func (a *V5Adapter) DashboardList(ctx context.Context, filter string) ([]zabbix.UnifiedDashboard, error) {
	// TODO: implement V5 screen.get
	return nil, fmt.Errorf("not implemented")
}
