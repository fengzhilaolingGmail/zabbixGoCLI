package cli

import (
	"context"
	"fmt"
	"strings"

	"zabbix-maint/internal/version"
	"zabbix-maint/pkg/zabbix"
)

// MenuNode 菜单节点
type MenuNode struct {
	ID          string
	Title       string
	Description string
	Action      func(ctx context.Context, client zabbix.ZabbixOperator, ver version.Version) error
	SubMenus    []*MenuNode
	Parent      *MenuNode
	VersionHint string
}

// MenuTree 完整菜单树
type MenuTree struct {
	Root    *MenuNode
	Current *MenuNode
}

// BuildMenuTree 构建菜单树
func BuildMenuTree() *MenuTree {
	root := &MenuNode{
		ID:    "root",
		Title: "主菜单",
		SubMenus: []*MenuNode{
			{
				ID:    "user",
				Title: "用户管理",
				SubMenus: []*MenuNode{
					{ID: "user-create", Title: "创建用户", Action: handleUserCreate},
					{ID: "user-disable", Title: "禁用用户", Action: handleUserDisable},
					{ID: "user-enable", Title: "启用用户", Action: handleUserEnable},
					{ID: "user-passwd", Title: "修改密码", Action: handleUserPasswd},
					{ID: "user-delete", Title: "删除用户", Action: handleUserDelete},
					{ID: "user-list", Title: "查看用户列表", Action: handleUserList},
					{ID: "user-detail", Title: "查看用户详情", Action: handleUserDetail},
				},
			},
			{
				ID:    "usergroup",
				Title: "用户组管理",
				SubMenus: []*MenuNode{
					{ID: "ug-create", Title: "创建用户组", Action: handleUserGroupCreate},
					{ID: "ug-disable", Title: "禁用用户组", Action: handleUserGroupDisable},
					{ID: "ug-clear", Title: "清空用户组", Action: handleUserGroupClear},
					{ID: "ug-delete", Title: "删除用户组", Action: handleUserGroupDelete},
					{ID: "ug-list", Title: "查看用户组列表", Action: handleUserGroupList},
				},
			},
			{
				ID:    "host",
				Title: "主机管理",
				SubMenus: []*MenuNode{
					{ID: "host-clone", Title: "主机全克隆", Action: handleHostClone},
					{ID: "host-list", Title: "查看主机列表", Action: handleHostList},
					{ID: "host-detail", Title: "查看主机详情", Action: handleHostDetail},
				},
			},
			{
				ID:    "hostgroup",
				Title: "主机组管理",
				SubMenus: []*MenuNode{
					{ID: "hg-create", Title: "创建主机组", Action: handleHostGroupCreate},
					{ID: "hg-clear", Title: "清空主机组", Action: handleHostGroupClear},
					{ID: "hg-delete", Title: "删除主机组", Action: handleHostGroupDelete},
					{ID: "hg-list", Title: "查看主机组列表", Action: handleHostGroupList},
				},
			},
			{
				ID:    "dashboard",
				Title: "聚合图形 / 仪表盘管理",
				SubMenus: []*MenuNode{
					{ID: "dash-create", Title: "创建聚合图形/仪表盘", Action: handleDashboardCreate},
					{ID: "dash-add-widget", Title: "添加 Widget / 子项", Action: handleDashboardAddWidget, VersionHint: "5.x=Screen, 7.x=Dashboard"},
					{ID: "dash-list", Title: "查看列表", Action: handleDashboardList},
				},
			},
			{
				ID:    "batch",
				Title: "批量操作",
				SubMenus: []*MenuNode{
					{ID: "batch-user-create", Title: "批量创建用户", Action: handleBatchUserCreate},
					{ID: "batch-host-clone", Title: "批量克隆主机", Action: handleBatchHostClone},
				},
			},
			{
				ID:    "system",
				Title: "系统信息",
				SubMenus: []*MenuNode{
					{ID: "sys-version", Title: "API 版本信息", Action: handleSysVersion},
					{ID: "sys-status", Title: "服务器状态", Action: handleSysStatus},
				},
			},
		},
	}
	return &MenuTree{Root: root, Current: root}
}

// ============================================================
// 辅助函数
// ============================================================

func promptString(prompt string, required bool) (string, error) {
	fmt.Printf("  >> %s: ", prompt)
	var input string
	_, err := fmt.Scanln(&input)
	if err != nil {
		return "", err
	}
	if required && strings.TrimSpace(input) == "" {
		return "", fmt.Errorf("输入不能为空")
	}
	return input, nil
}

func promptInt(prompt string) (int, error) {
	fmt.Printf("  >> %s: ", prompt)
	var val int
	_, err := fmt.Scanln(&val)
	if err != nil {
		return 0, err
	}
	return val, nil
}

func promptConfirm(prompt string) bool {
	fmt.Printf("  >> %s [Y/n]: ", prompt)
	var input string
	fmt.Scanln(&input)
	input = strings.ToLower(strings.TrimSpace(input))
	return input == "" || input == "y" || input == "yes"
}

func printUserTable(users []zabbix.UnifiedUser) {
	if len(users) == 0 {
		fmt.Println("  (未找到用户)")
		return
	}
	fmt.Printf("  %-10s %-20s %-20s %-10s\n", "ID", "用户名", "姓名", "状态")
	fmt.Println("  " + strings.Repeat("-", 65))
	for _, u := range users {
		status := "已启用"
		if !u.Enabled {
			status = "已禁用"
		}
		fmt.Printf("  %-10s %-20s %-20s %-10s\n", u.ID, u.Alias, u.Name, status)
	}
}

func printHostTable(hosts []zabbix.UnifiedHost) {
	if len(hosts) == 0 {
		fmt.Println("  (未找到主机)")
		return
	}
	fmt.Printf("  %-10s %-25s %-30s\n", "ID", "主机名", "显示名称")
	fmt.Println("  " + strings.Repeat("-", 70))
	for _, h := range hosts {
		fmt.Printf("  %-10s %-25s %-30s\n", h.ID, h.Host, h.Name)
	}
}

func printUserGroupTable(groups []zabbix.UnifiedUserGroup) {
	if len(groups) == 0 {
		fmt.Println("  (未找到用户组)")
		return
	}
	fmt.Printf("  %-10s %-30s %-10s\n", "ID", "名称", "用户数")
	fmt.Println("  " + strings.Repeat("-", 55))
	for _, g := range groups {
		fmt.Printf("  %-10s %-30s %-10d\n", g.ID, g.Name, g.UserCount)
	}
}

func printHostGroupTable(groups []zabbix.UnifiedHostGroup) {
	if len(groups) == 0 {
		fmt.Println("  (未找到主机组)")
		return
	}
	fmt.Printf("  %-10s %-30s %-10s\n", "ID", "名称", "主机数")
	fmt.Println("  " + strings.Repeat("-", 55))
	for _, g := range groups {
		fmt.Printf("  %-10s %-30s %-10d\n", g.ID, g.Name, g.HostCount)
	}
}

func printDashboardTable(dashes []zabbix.UnifiedDashboard) {
	if len(dashes) == 0 {
		fmt.Println("  (未找到仪表盘/聚合图形)")
		return
	}
	fmt.Printf("  %-10s %-30s %-10s\n", "ID", "名称", "类型")
	fmt.Println("  " + strings.Repeat("-", 55))
	for _, d := range dashes {
		fmt.Printf("  %-10s %-30s %-10s\n", d.ID, d.Name, d.Type)
	}
}

// ============================================================
// 用户管理
// ============================================================

func handleUserCreate(ctx context.Context, client zabbix.ZabbixOperator, ver version.Version) error {
	alias, err := promptString("用户名 (alias)", true)
	if err != nil {
		return err
	}
	name, _ := promptString("姓名", false)
	passwd, err := promptString("密码", true)
	if err != nil {
		return err
	}
	groupID, err := promptString("用户组 ID", true)
	if err != nil {
		return err
	}

	req := zabbix.UserCreateReq{
		Alias:        alias,
		Name:         name,
		Passwd:       passwd,
		UserGroupIDs: []string{groupID},
	}

	if ver == version.Version7 {
		roleID, _ := promptString("角色 ID (7.x 必填)", false)
		req.RoleID = roleID
	}

	if !promptConfirm(fmt.Sprintf("确认创建用户 '%s'", alias)) {
		fmt.Println("  已取消。")
		return nil
	}

	id, err := client.UserCreate(ctx, req)
	if err != nil {
		return err
	}
	fmt.Printf("\n  [OK] 用户创建成功! ID: %s\n", id)
	return nil
}

func handleUserDisable(ctx context.Context, client zabbix.ZabbixOperator, ver version.Version) error {
	id, err := promptString("要禁用的用户 ID", true)
	if err != nil {
		return err
	}
	if !promptConfirm(fmt.Sprintf("确认禁用用户 '%s'", id)) {
		fmt.Println("  已取消。")
		return nil
	}
	if err := client.UserDisable(ctx, id); err != nil {
		return err
	}
	fmt.Println("  [OK] 用户已禁用。")
	return nil
}

func handleUserEnable(ctx context.Context, client zabbix.ZabbixOperator, ver version.Version) error {
	id, err := promptString("要启用的用户 ID", true)
	if err != nil {
		return err
	}
	if !promptConfirm(fmt.Sprintf("确认启用用户 '%s'", id)) {
		fmt.Println("  已取消。")
		return nil
	}
	if err := client.UserEnable(ctx, id); err != nil {
		return err
	}
	fmt.Println("  [OK] 用户已启用。")
	return nil
}

func handleUserPasswd(ctx context.Context, client zabbix.ZabbixOperator, ver version.Version) error {
	id, err := promptString("用户 ID", true)
	if err != nil {
		return err
	}
	passwd, err := promptString("新密码", true)
	if err != nil {
		return err
	}
	if !promptConfirm(fmt.Sprintf("确认修改用户 '%s' 的密码", id)) {
		fmt.Println("  已取消。")
		return nil
	}
	if err := client.UserUpdatePassword(ctx, id, passwd); err != nil {
		return err
	}
	fmt.Println("  [OK] 密码已修改。")
	return nil
}

func handleUserDelete(ctx context.Context, client zabbix.ZabbixOperator, ver version.Version) error {
	id, err := promptString("要删除的用户 ID", true)
	if err != nil {
		return err
	}
	if !promptConfirm(fmt.Sprintf("确认删除用户 '%s' (不可恢复)", id)) {
		fmt.Println("  已取消。")
		return nil
	}
	if err := client.UserDelete(ctx, id); err != nil {
		return err
	}
	fmt.Println("  [OK] 用户已删除。")
	return nil
}

func handleUserList(ctx context.Context, client zabbix.ZabbixOperator, ver version.Version) error {
	users, err := client.UserList(ctx, "")
	if err != nil {
		return err
	}
	fmt.Printf("\n  共 %d 个用户:\n\n", len(users))
	printUserTable(users)
	return nil
}

func handleUserDetail(ctx context.Context, client zabbix.ZabbixOperator, ver version.Version) error {
	id, err := promptString("用户 ID", true)
	if err != nil {
		return err
	}
	u, err := client.UserDetail(ctx, id)
	if err != nil {
		return err
	}
	fmt.Printf("\n  ID:       %s\n", u.ID)
	fmt.Printf("  用户名:   %s\n", u.Alias)
	fmt.Printf("  姓名:     %s\n", u.Name)
	fmt.Printf("  姓氏:     %s\n", u.Surname)
	fmt.Printf("  状态:     %v\n", u.Enabled)
	fmt.Printf("  角色ID:   %s\n", u.RoleID)
	fmt.Printf("  用户组:   %v\n", u.GroupNames)
	return nil
}

// ============================================================
// 用户组管理
// ============================================================

func handleUserGroupCreate(ctx context.Context, client zabbix.ZabbixOperator, ver version.Version) error {
	name, err := promptString("用户组名称", true)
	if err != nil {
		return err
	}
	if !promptConfirm(fmt.Sprintf("确认创建用户组 '%s'", name)) {
		fmt.Println("  已取消。")
		return nil
	}
	req := zabbix.UserGroupCreateReq{Name: name}
	id, err := client.UserGroupCreate(ctx, req)
	if err != nil {
		return err
	}
	fmt.Printf("  [OK] 用户组创建成功! ID: %s\n", id)
	return nil
}

func handleUserGroupDisable(ctx context.Context, client zabbix.ZabbixOperator, ver version.Version) error {
	id, err := promptString("要禁用的用户组 ID", true)
	if err != nil {
		return err
	}
	if !promptConfirm(fmt.Sprintf("确认禁用用户组 '%s'", id)) {
		fmt.Println("  已取消。")
		return nil
	}
	if err := client.UserGroupDisable(ctx, id); err != nil {
		return err
	}
	fmt.Println("  [OK] 用户组已禁用。")
	return nil
}

func handleUserGroupClear(ctx context.Context, client zabbix.ZabbixOperator, ver version.Version) error {
	id, err := promptString("要清空的用户组 ID", true)
	if err != nil {
		return err
	}
	if !promptConfirm(fmt.Sprintf("确认清空用户组 '%s' 中的所有用户", id)) {
		fmt.Println("  已取消。")
		return nil
	}
	if err := client.UserGroupClear(ctx, id); err != nil {
		return err
	}
	fmt.Println("  [OK] 用户组已清空。")
	return nil
}

func handleUserGroupDelete(ctx context.Context, client zabbix.ZabbixOperator, ver version.Version) error {
	id, err := promptString("要删除的用户组 ID", true)
	if err != nil {
		return err
	}
	if !promptConfirm(fmt.Sprintf("确认删除用户组 '%s' (不可恢复)", id)) {
		fmt.Println("  已取消。")
		return nil
	}
	if err := client.UserGroupDelete(ctx, id); err != nil {
		return err
	}
	fmt.Println("  [OK] 用户组已删除。")
	return nil
}

func handleUserGroupList(ctx context.Context, client zabbix.ZabbixOperator, ver version.Version) error {
	groups, err := client.UserGroupList(ctx, "")
	if err != nil {
		return err
	}
	fmt.Printf("\n  共 %d 个用户组:\n\n", len(groups))
	printUserGroupTable(groups)
	return nil
}

// ============================================================
// 主机管理
// ============================================================

func handleHostClone(ctx context.Context, client zabbix.ZabbixOperator, ver version.Version) error {
	srcID, err := promptString("源主机 ID", true)
	if err != nil {
		return err
	}
	newName, err := promptString("新主机名称", true)
	if err != nil {
		return err
	}
	newIP, _ := promptString("新主机 IP (可选)", false)

	req := zabbix.HostCloneReq{
		SourceHostID: srcID,
		NewHostName:  newName,
		NewHostIP:    newIP,
	}

	if !promptConfirm(fmt.Sprintf("确认克隆主机 '%s' 为 '%s'", srcID, newName)) {
		fmt.Println("  已取消。")
		return nil
	}

	id, err := client.HostFullClone(ctx, req)
	if err != nil {
		return err
	}
	fmt.Printf("  [OK] 主机克隆成功! ID: %s\n", id)
	return nil
}

func handleHostList(ctx context.Context, client zabbix.ZabbixOperator, ver version.Version) error {
	hosts, err := client.HostList(ctx, "")
	if err != nil {
		return err
	}
	fmt.Printf("\n  共 %d 个主机:\n\n", len(hosts))
	printHostTable(hosts)
	return nil
}

func handleHostDetail(ctx context.Context, client zabbix.ZabbixOperator, ver version.Version) error {
	id, err := promptString("主机 ID", true)
	if err != nil {
		return err
	}
	h, err := client.HostDetail(ctx, id)
	if err != nil {
		return err
	}
	fmt.Printf("\n  ID:       %s\n", h.ID)
	fmt.Printf("  主机名:   %s\n", h.Host)
	fmt.Printf("  显示名:   %s\n", h.Name)
	fmt.Printf("  状态:     %d\n", h.Status)
	fmt.Printf("  主机组:   %v\n", h.GroupNames)
	fmt.Printf("  模板:     %v\n", h.Templates)
	return nil
}

// ============================================================
// 主机组管理
// ============================================================

func handleHostGroupCreate(ctx context.Context, client zabbix.ZabbixOperator, ver version.Version) error {
	name, err := promptString("主机组名称", true)
	if err != nil {
		return err
	}
	if !promptConfirm(fmt.Sprintf("确认创建主机组 '%s'", name)) {
		fmt.Println("  已取消。")
		return nil
	}
	id, err := client.HostGroupCreate(ctx, name)
	if err != nil {
		return err
	}
	fmt.Printf("  [OK] 主机组创建成功! ID: %s\n", id)
	return nil
}

func handleHostGroupClear(ctx context.Context, client zabbix.ZabbixOperator, ver version.Version) error {
	id, err := promptString("要清空的主机组 ID", true)
	if err != nil {
		return err
	}
	if !promptConfirm(fmt.Sprintf("确认清空主机组 '%s' 中的所有主机", id)) {
		fmt.Println("  已取消。")
		return nil
	}
	if err := client.HostGroupClear(ctx, id); err != nil {
		return err
	}
	fmt.Println("  [OK] 主机组已清空。")
	return nil
}

func handleHostGroupDelete(ctx context.Context, client zabbix.ZabbixOperator, ver version.Version) error {
	id, err := promptString("要删除的主机组 ID", true)
	if err != nil {
		return err
	}
	if !promptConfirm(fmt.Sprintf("确认删除主机组 '%s' (不可恢复)", id)) {
		fmt.Println("  已取消。")
		return nil
	}
	if err := client.HostGroupDelete(ctx, id); err != nil {
		return err
	}
	fmt.Println("  [OK] 主机组已删除。")
	return nil
}

func handleHostGroupList(ctx context.Context, client zabbix.ZabbixOperator, ver version.Version) error {
	groups, err := client.HostGroupList(ctx, "")
	if err != nil {
		return err
	}
	fmt.Printf("\n  共 %d 个主机组:\n\n", len(groups))
	printHostGroupTable(groups)
	return nil
}

// ============================================================
// 聚合图形/仪表盘
// ============================================================

func handleDashboardCreate(ctx context.Context, client zabbix.ZabbixOperator, ver version.Version) error {
	name, err := promptString("仪表盘/聚合图形名称", true)
	if err != nil {
		return err
	}
	if !promptConfirm(fmt.Sprintf("确认创建 '%s'", name)) {
		fmt.Println("  已取消。")
		return nil
	}
	req := zabbix.DashboardCreateReq{Name: name}
	id, err := client.DashboardCreate(ctx, req)
	if err != nil {
		return err
	}
	fmt.Printf("  [OK] 创建成功! ID: %s\n", id)
	return nil
}

func handleDashboardAddWidget(ctx context.Context, client zabbix.ZabbixOperator, ver version.Version) error {
	fmt.Println("  (交互模式下尚未实现)")
	fmt.Println("  请使用一次性命令: zbx-cli -i <实例> dashboard add ...")
	return nil
}

func handleDashboardList(ctx context.Context, client zabbix.ZabbixOperator, ver version.Version) error {
	dashes, err := client.DashboardList(ctx, "")
	if err != nil {
		return err
	}
	fmt.Printf("\n  共 %d 个仪表盘/聚合图形:\n\n", len(dashes))
	printDashboardTable(dashes)
	return nil
}

// ============================================================
// 批量操作
// ============================================================

func handleBatchUserCreate(ctx context.Context, client zabbix.ZabbixOperator, ver version.Version) error {
	fmt.Println("  (交互模式下尚未实现)")
	fmt.Println("  请使用一次性命令: zbx-cli -i <实例> batch user-create --file ...")
	return nil
}

func handleBatchHostClone(ctx context.Context, client zabbix.ZabbixOperator, ver version.Version) error {
	fmt.Println("  (交互模式下尚未实现)")
	fmt.Println("  请使用一次性命令: zbx-cli -i <实例> batch host-clone --file ...")
	return nil
}

// ============================================================
// 系统信息
// ============================================================

func handleSysVersion(ctx context.Context, client zabbix.ZabbixOperator, ver version.Version) error {
	apiVer, err := client.APIVersion(ctx)
	if err != nil {
		return err
	}
	fmt.Printf("\n  Zabbix API 版本: %s\n", apiVer)
	fmt.Printf("  探测版本:         %s\n", ver)
	return nil
}

func handleSysStatus(ctx context.Context, client zabbix.ZabbixOperator, ver version.Version) error {
	status, err := client.ServerStatus(ctx)
	if err != nil {
		return err
	}
	fmt.Printf("\n  服务器状态:\n")
	for k, v := range status {
		fmt.Printf("    %s: %v\n", k, v)
	}
	return nil
}
