package zabbix

import "context"

// ZabbixOperator 统一业务操作接口
type ZabbixOperator interface {
	// ========== 用户管理 ==========
	UserCreate(ctx context.Context, req UserCreateReq) (string, error)
	UserDisable(ctx context.Context, userID string) error
	UserEnable(ctx context.Context, userID string) error
	UserUpdatePassword(ctx context.Context, userID, newPass string) error
	UserDelete(ctx context.Context, userID string) error
	UserList(ctx context.Context, filter string) ([]UnifiedUser, error)
	UserDetail(ctx context.Context, userID string) (*UnifiedUser, error)

	// ========== 用户组管理 ==========
	UserGroupCreate(ctx context.Context, req UserGroupCreateReq) (string, error)
	UserGroupDisable(ctx context.Context, groupID string) error
	UserGroupEnable(ctx context.Context, groupID string) error
	UserGroupClear(ctx context.Context, groupID string) error
	UserGroupDelete(ctx context.Context, groupID string) error
	UserGroupList(ctx context.Context, filter string) ([]UnifiedUserGroup, error)

	// ========== 主机组管理 ==========
	HostGroupCreate(ctx context.Context, name string) (string, error)
	HostGroupClear(ctx context.Context, groupID string) error
	HostGroupDelete(ctx context.Context, groupID string) error
	HostGroupList(ctx context.Context, filter string) ([]UnifiedHostGroup, error)

	// ========== 主机管理 ==========
	HostFullClone(ctx context.Context, req HostCloneReq) (string, error)
	HostList(ctx context.Context, filter string) ([]UnifiedHost, error)
	HostDetail(ctx context.Context, hostID string) (*UnifiedHost, error)

	// ========== 聚合图形/仪表盘 ==========
	DashboardCreate(ctx context.Context, req DashboardCreateReq) (string, error)
	DashboardAddWidget(ctx context.Context, req DashboardAddReq) (string, error)
	DashboardList(ctx context.Context, filter string) ([]UnifiedDashboard, error)

	// ========== 系统信息 ==========
	APIVersion(ctx context.Context) (string, error)
	ServerStatus(ctx context.Context) (map[string]interface{}, error)

	// ========== 角色管理 (7.x) ==========
	RoleList(ctx context.Context, filter string) ([]UnifiedRole, error)
}
