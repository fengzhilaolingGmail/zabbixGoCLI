<!--
 * @Author: fengzhilaoling
 * @Date: 2026-06-26 10:28:15
 * @LastEditors: fengzhilaoling
 * @LastEditTime: 2026-06-26 10:28:31
 * @FilePath: \zabbixGoCLI\zabbix_maint_design.md
 * @Description: 文件详情
-->
# Zabbix API 维护工具设计文档

> 版本: v1.1  
> 语言: Go  
> 支持 Zabbix 版本: 5.x LTS / 7.x LTS  
> 设计目标: 版本兼容层 + 统一业务接口 + 交互式 CLI

---

## 1. 概述

本工具基于 Go 语言封装 Zabbix JSON-RPC API，通过**版本兼容层（Version Adapter）**屏蔽 5.x 与 7.x 的接口差异，向上层提供统一、稳定的业务接口。同时提供**交互式 CLI**，支持菜单式操作，降低使用门槛。

### 1.1 核心设计原则

| 原则 | 说明 |
|------|------|
| **接口统一** | 上层业务只依赖抽象接口，不感知底层版本差异 |
| **最小侵入** | 兼容层只做参数转换与路由分发，不修改业务语义 |
| **可扩展** | 新增 6.x/8.x 版本时，只需新增 Adapter 实现 |
| **交互友好** | 支持交互式菜单 + 命令行参数两种模式 |
| **防御式编程** | 所有 API 调用均带重试、超时、错误降级 |

---

## 2. 版本差异分析（5.x vs 7.x）

### 2.1 API 方法映射表

| 业务功能 | Zabbix 5.x API | Zabbix 7.x API | 差异说明 |
|---------|---------------|---------------|---------|
| 用户创建 | `user.create` | `user.create` | 7.x 新增 `roleid` 必填，5.x 用 `usrgrps` |
| 用户禁用 | `user.update` (status=1) | `user.update` (roleid=禁用角色) | 7.x 移除 `status` 字段，改用角色控制 |
| 密码修改 | `user.update` (passwd) | `user.update` (passwd) | 一致，但 7.x 有密码复杂度策略 |
| 用户组创建 | `usergroup.create` | `usergroup.create` | 一致 |
| 用户组禁用 | `usergroup.update` (users_status=1) | `usergroup.update` (userdirectoryid=0) | 7.x 用户组状态逻辑变化 |
| 用户组清空 | `usergroup.update` (清空 userids) | `usergroup.update` (清空 userids) | 一致 |
| 主机组创建 | `hostgroup.create` | `hostgroup.create` | 一致 |
| 主机组清空 | `hostgroup.massremove` (hostids) | `hostgroup.update` (清空 hostids) | 7.x 推荐用 update 清空 |
| 主机组删除 | `hostgroup.delete` | `hostgroup.delete` | 一致 |
| 主机全克隆 | `host.fullclone` (前端动作) | `host.clone` | 5.x 无原生 API，需手动实现；7.x 提供 `host.clone` |
| 聚合图形 | `screen.create/update` | `dashboard.create/update` | 7.x 废弃 Screen，改用 Dashboard |
| 聚合图形添加 | `screenitem.create` | `dashboard.widget.create` | 数据结构完全不同 |

### 2.2 关键 Breaking Changes

1. **Screen → Dashboard**：Zabbix 7.0 彻底废弃 Screen 概念，所有聚合图形迁移为 Dashboard Widget。
2. **用户状态控制**：7.0 移除用户 `status` 字段，用户启用/禁用通过分配不同 Role（如 Disabled role）实现。
3. **克隆简化**：7.0 将 "Full clone" 重命名为 "Clone"，并移除了旧的简单 Clone 选项。
4. **用户组状态**：7.0 用户组禁用逻辑与 5.x 不同，需通过 `userdirectoryid` 或权限配置控制。

---

## 3. 总体架构

```
┌─────────────────────────────────────────────────────────────┐
│                      CLI / Interactive Layer                  │
│  ┌─────────────────┐  ┌─────────────────┐  ┌──────────────┐ │
│  │  Interactive    │  │  One-Shot CMD   │  │  Config CMD  │ │
│  │  Menu Mode      │  │  Mode           │  │  Mode        │ │
│  │  (-i instance)  │  │  (direct args)  │  │  (setup)     │ │
│  └────────┬────────┘  └────────┬────────┘  └──────┬───────┘ │
│           │                    │                  │         │
│           └────────────────────┴──────────────────┘         │
│                              │                              │
└──────────────────────────────┼──────────────────────────────┘
                               │
┌──────────────────────────────┼──────────────────────────────┐
│                        Business Layer                        │
│  ┌─────────┐ ┌──────────┐ ┌─────────┐ ┌──────────────────┐ │
│  │ UserMgr │ │UserGrpMgr│ │ HostMgr │ │ DashboardMgr     │ │
│  └────┬────┘ └────┬─────┘ └────┬────┘ └────────┬─────────┘ │
│       │           │            │               │           │
│       └───────────┴────────────┴───────────────┘           │
│                         │                                  │
│              ┌──────────┴──────────┐                      │
│              │   Unified Interface   │                      │
│              │  (ZabbixOperator)     │                      │
│              └──────────┬──────────┘                      │
└─────────────────────────┬──────────────────────────────────┘
                          │
┌─────────────────────────┼──────────────────────────────────┐
│              Version Compatibility Layer                     │
│  ┌──────────────────────┼──────────────────────┐          │
│  │                      │                      │          │
│  │   ┌──────────┐   ┌──────────┐   ┌────────┐ │          │
│  │   │ V5Adapter│   │ V7Adapter│   │V6Adapter│ │  ...    │
│  │   └────┬─────┘   └────┬─────┘   └───┬────┘ │          │
│  │        │              │              │      │          │
│  │        └──────────────┴──────────────┘      │          │
│  │                   │                        │          │
│  │          ┌────────┴────────┐              │          │
│  │          │  Version Router   │              │          │
│  │          │  (Factory)        │              │          │
│  │          └────────┬────────┘              │          │
│  └───────────────────┼──────────────────────┘          │
│                      │                                   │
│  ┌───────────────────┼──────────────────────────────────┐│
│  │         Transport Layer (HTTP Client)                  ││
│  │  ┌─────────────┐  ┌─────────────┐  ┌─────────────┐   ││
│  │  │ JSON-RPC    │  │ Retry/Backoff│  │ Auth Manager│   ││
│  │  │ Client      │  │ Middleware   │  │ (Token/Key) │   ││
│  │  └─────────────┘  └─────────────┘  └─────────────┘   ││
│  └──────────────────────────────────────────────────────┘│
└──────────────────────────────────────────────────────────┘
```

---

## 4. CLI 交互层设计

### 4.1 命令行入口

```bash
# 交互式模式（主模式）
zbx-cli -i prod-zbx5

# 一次性命令模式
zbx-cli -i prod-zbx5 user create --alias zhangsan --group 7 --password Temp123!
zbx-cli -i prod-zbx7 host clone --src 10084 --name web-clone-01

# 配置管理
zbx-cli config add --name prod-zbx5 --url http://zabbix5/api_jsonrpc.php --user Admin --pass zabbix
zbx-cli config list
zbx-cli config remove prod-zbx5

# 版本探测
zbx-cli detect --url http://zabbix.example.com/api_jsonrpc.php
```

### 4.2 交互式菜单流程

```
$ zbx-cli -i prod-zbx5

═══════════════════════════════════════════════════════════════
  Zabbix 维护工具 v1.0
  当前实例: prod-zbx5 [Zabbix 5.0 LTS]
  连接状态: ✅ 已连接
═══════════════════════════════════════════════════════════════

  请选择功能模块:

  ┌────────────────────────────────────────────┐
  │  1. 用户管理                                │
  │  2. 用户组管理                              │
  │  3. 主机管理                                │
  │  4. 主机组管理                              │
  │  5. 聚合图形 / 仪表盘管理                    │
  │  6. 批量操作                                │
  │  7. 系统信息                                │
  │  0. 退出                                    │
  └────────────────────────────────────────────┘

  请选择 [0-7]: 1

═══════════════════════════════════════════════════════════════
  用户管理
═══════════════════════════════════════════════════════════════

  1. 创建用户
  2. 禁用用户
  3. 启用用户
  4. 修改密码
  5. 删除用户
  6. 查看用户列表
  7. 查看用户详情
  0. 返回上级

  请选择 [0-7]: 1

  → 请输入用户名: zhangsan
  → 请输入姓名: 张三
  → 请输入密码: ********
  → 请选择用户组 (多选, 逗号分隔):
      1. Guests (ID: 7)
      2. Zabbix administrators (ID: 8)
      3. Disabled (ID: 9)
    请选择: 1
  → 请选择角色 (仅 7.x 生效, 当前实例为 5.x 自动跳过)

  ⚠️  请确认操作:
      创建用户: zhangsan
      用户组: Guests
      确认执行? [Y/n]: Y

  ✅ 用户创建成功! 用户ID: 10001

  按 Enter 继续...
```

### 4.3 交互式菜单结构定义

```go
// MenuNode 菜单节点
type MenuNode struct {
    ID          string
    Title       string
    Description string
    Action      func(ctx context.Context) error  // 叶子节点执行动作
    SubMenus    []*MenuNode                      // 子菜单
    Parent      *MenuNode                        // 父节点（用于返回）
    VersionHint string                           // 版本提示，如 "仅 7.x"
}

// MenuTree 完整菜单树
type MenuTree struct {
    Root *MenuNode
    Current *MenuNode
}

func BuildMenuTree() *MenuTree {
    root := &MenuNode{
        ID: "root",
        Title: "主菜单",
        SubMenus: []*MenuNode{
            {
                ID: "user",
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
                ID: "usergroup",
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
                ID: "host",
                Title: "主机管理",
                SubMenus: []*MenuNode{
                    {ID: "host-clone", Title: "主机全克隆", Action: handleHostClone},
                    {ID: "host-list", Title: "查看主机列表", Action: handleHostList},
                    {ID: "host-detail", Title: "查看主机详情", Action: handleHostDetail},
                },
            },
            {
                ID: "hostgroup",
                Title: "主机组管理",
                SubMenus: []*MenuNode{
                    {ID: "hg-create", Title: "创建主机组", Action: handleHostGroupCreate},
                    {ID: "hg-clear", Title: "清空主机组", Action: handleHostGroupClear},
                    {ID: "hg-delete", Title: "删除主机组", Action: handleHostGroupDelete},
                    {ID: "hg-list", Title: "查看主机组列表", Action: handleHostGroupList},
                },
            },
            {
                ID: "dashboard",
                Title: "聚合图形 / 仪表盘管理",
                SubMenus: []*MenuNode{
                    {ID: "dash-create", Title: "创建聚合图形/仪表盘", Action: handleDashboardCreate},
                    {ID: "dash-add-widget", Title: "添加 Widget / 子项", Action: handleDashboardAddWidget, VersionHint: "5.x=Screen, 7.x=Dashboard"},
                    {ID: "dash-list", Title: "查看列表", Action: handleDashboardList},
                },
            },
            {
                ID: "batch",
                Title: "批量操作",
                SubMenus: []*MenuNode{
                    {ID: "batch-user-create", Title: "批量创建用户", Action: handleBatchUserCreate},
                    {ID: "batch-host-clone", Title: "批量克隆主机", Action: handleBatchHostClone},
                },
            },
            {
                ID: "system",
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
```

### 4.4 交互式输入组件

```go
// PromptReader 交互式输入读取器
type PromptReader interface {
    // 基础输入
    String(prompt string, required bool) (string, error)
    Password(prompt string) (string, error)
    Int(prompt string, min, max int) (int, error)
    Confirm(prompt string, defaultYes bool) (bool, error)

    // 选择器
    Select(prompt string, options []Option) (string, error)           // 单选
    MultiSelect(prompt string, options []Option) ([]string, error)    // 多选
    TableSelect(prompt string, headers []string, rows [][]string) (int, error) // 表格选择

    // 搜索
    SearchSelect(prompt string, searcher func(query string) ([]Option, error)) (string, error)
}

type Option struct {
    ID      string
    Label   string
    Detail  string
    Disabled bool // 不可选
}
```

### 4.5 交互式流程示例（用户创建）

```go
func handleUserCreate(ctx context.Context, client zabbix.ZabbixOperator, version version.Version) error {
    reader := NewPromptReader()

    // 1. 输入用户名
    alias, err := reader.String("请输入用户名", true)
    if err != nil { return err }

    // 2. 输入姓名（可选）
    name, _ := reader.String("请输入姓名(可选)", false)

    // 3. 输入密码
    passwd, err := reader.Password("请输入初始密码")
    if err != nil { return err }

    // 4. 查询并选择用户组（搜索+选择）
    groupID, err := reader.SearchSelect("请选择用户组", func(query string) ([]Option, error) {
        groups, err := client.UserGroupList(ctx, query)
        if err != nil { return nil, err }
        opts := make([]Option, len(groups))
        for i, g := range groups {
            opts[i] = Option{ID: g.ID, Label: g.Name}
        }
        return opts, nil
    })
    if err != nil { return err }

    // 5. 7.x 专属：选择角色
    var roleID string
    if version == version.Version7 {
        roleID, err = reader.SearchSelect("请选择角色", func(query string) ([]Option, error) {
            roles, err := client.RoleList(ctx, query)
            // ...
        })
    }

    // 6. 确认
    ok, err := reader.Confirm(fmt.Sprintf("创建用户 %s, 用户组 %s?", alias, groupID), true)
    if !ok { return fmt.Errorf("用户取消") }

    // 7. 执行
    req := zabbix.UserCreateReq{
        Alias: alias,
        Name:  name,
        Passwd: passwd,
        UserGroupIDs: []string{groupID},
        RoleID: roleID,
    }
    userID, err := client.UserCreate(ctx, req)
    if err != nil { return err }

    // 8. 结果展示
    fmt.Printf("✅ 用户创建成功! ID: %s\n", userID)
    return nil
}
```

### 4.6 表格展示组件

```go
// TableRenderer 表格渲染器
type TableRenderer interface {
    Render(headers []string, rows [][]string)
    RenderWithIndex(headers []string, rows [][]string) // 带序号
    RenderJSON(data interface{}) // JSON 原始输出
}

// 使用示例：用户列表
func handleUserList(ctx context.Context, client zabbix.ZabbixOperator) error {
    users, err := client.UserList(ctx)
    if err != nil { return err }

    headers := []string{"ID", "用户名", "姓名", "状态", "用户组", "角色"}
    rows := make([][]string, len(users))
    for i, u := range users {
        status := "✅ 启用"
        if !u.Enabled { status = "❌ 禁用" }
        rows[i] = []string{u.ID, u.Alias, u.Name, status, strings.Join(u.GroupNames, ","), u.RoleName}
    }

    table := NewTableRenderer()
    table.RenderWithIndex(headers, rows)
    return nil
}
```

---

## 5. 模块详细设计

### 5.1 版本兼容层（Version Adapter）

#### 5.1.1 核心接口定义

```go
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
```

#### 5.1.2 版本工厂

```go
type Version string

const (
    Version5  Version = "5.x"
    Version7  Version = "7.x"
)

type AdapterFactory interface {
    Create(version Version, endpoint string, auth AuthConfig) (ZabbixOperator, error)
}
```

#### 5.1.3 V5Adapter 与 V7Adapter 差异处理

| 方法 | V5Adapter 实现 | V7Adapter 实现 |
|------|---------------|---------------|
| `UserDisable` | 调用 `user.update` 设置 `status: 1` | 调用 `user.update` 设置 `roleid: <disabled_role_id>` |
| `UserEnable` | 调用 `user.update` 设置 `status: 0` | 调用 `user.update` 设置 `roleid: <default_user_role_id>` |
| `UserGroupDisable` | 调用 `usergroup.update` 设置 `users_status: 1` | 调用 `usergroup.update` 设置 `userdirectoryid: 0` 或移除权限 |
| `UserGroupEnable` | 调用 `usergroup.update` 设置 `users_status: 0` | 恢复默认权限配置 |
| `HostFullClone` | 手动实现：get host → create host → copy items/triggers/graphs → copy interfaces | 直接调用 `host.clone` API |
| `DashboardAddWidget` | 调用 `screenitem.create` | 调用 `dashboard.widget.create` |

### 5.2 传输层（Transport Layer）

#### 5.2.1 JSON-RPC 客户端

```go
type JSONRPCClient struct {
    endpoint   string
    authToken  string
    httpClient *http.Client

    // 配置
    retryCount    int
    retryBackoff  time.Duration
    requestTimeout time.Duration
}

func (c *JSONRPCClient) Call(ctx context.Context, method string, params interface{}, result interface{}) error
```

#### 5.2.2 认证管理器

```go
type AuthManager struct {
    username string
    password string
    token    string
    expiry   time.Time
}

func (a *AuthManager) Refresh(ctx context.Context) error  // 调用 user.login
func (a *AuthManager) Token() string
```

### 5.3 业务模块设计

#### 5.3.1 用户管理（User Manager）

```go
type UserCreateReq struct {
    Alias      string   // 用户名（必填）
    Name       string   // 姓名
    Surname    string   // 姓氏
    Passwd     string   // 初始密码
    UserGroupIDs []string // 所属用户组
    RoleID     string   // 角色ID（7.x 必填，5.x 忽略）
    Media      []UserMedia // 告警媒介
}

type UserMedia struct {
    MediatypeID string
    Sendto      []string
    Severity    int64
    Period      string
}
```

**业务逻辑：**
- 创建时，V5 直接传入 `usrgrps`，V7 需额外传入 `roleid`（默认 Guest 或指定角色）
- 禁用时，V5 设置 `status=1`，V7 需先查询系统预置的 "Disabled" 角色 ID，再赋值给 `roleid`
- 启用时，V5 设置 `status=0`，V7 恢复为默认用户角色（如 User role）
- 密码修改，两版本一致，直接调用 `user.update` 传入 `passwd`

#### 5.3.2 用户组管理（User Group Manager）

```go
type UserGroupCreateReq struct {
    Name        string
    Rights      []GroupRight  // 主机组权限
    UserIDs     []string      // 初始用户列表
}

type GroupRight struct {
    Permission int  // 2=只读, 3=读写
    ID         string // hostgroupid
}
```

**业务逻辑：**
- 创建时，V5 可设置 `users_status`（默认启用），V7 该字段废弃
- 禁用时，V5 设置 `users_status=1`，V7 通过移除所有 `rights` 并设置 `userdirectoryid=0` 实现等效禁用
- 启用时，V5 设置 `users_status=0`，V7 恢复权限配置
- 清空时，统一调用 `usergroup.update` 将 `userids` 置为空数组

#### 5.3.3 主机组管理（Host Group Manager）

```go
type HostGroupOps struct {
    // 无额外结构，直接操作 hostgroupid
}
```

**业务逻辑：**
- 创建：`hostgroup.create` 两版本一致
- 清空：V5 调用 `hostgroup.massremove` 移除所有主机；V7 调用 `hostgroup.update` 将 `hosts` 置空
- 删除：V5/V7 均调用 `hostgroup.delete`，但需先确保组内无主机/模板

#### 5.3.4 主机全克隆（Host Full Clone）

```go
type HostCloneReq struct {
    SourceHostID string // 源主机ID
    NewHostName  string // 新主机名
    NewHostIP    string // 新主机IP（可选，覆盖接口IP）
    NewGroupIDs  []string // 新主机组（可选，继承则留空）
    CloneItems   bool     // 是否克隆监控项
    CloneTriggers bool    // 是否克隆触发器
    CloneGraphs  bool     // 是否克隆图形
}
```

**业务逻辑：**
- V5：无原生全克隆 API，需手动实现：
  1. `host.get` 获取源主机完整配置（含 interfaces, items, triggers, graphs, macros, inventory）
  2. `host.create` 创建新主机，复制所有配置
  3. `item.create` 复制监控项（需替换 hostid）
  4. `trigger.create` 复制触发器
  5. `graph.create` 复制图形
  6. `hostinterface.create` 复制接口（可选替换 IP）
- V7：直接调用 `host.clone`（原 Full clone 功能），传入 `hostid` 和 `host` 对象

#### 5.3.5 聚合图形/仪表盘（Dashboard Manager）

```go
type DashboardCreateReq struct {
    Name      string
    OwnerID   string // 所属用户ID
    ShareType int    // 0=私有, 1=公开
}

type DashboardAddReq struct {
    DashboardID string
    Widgets     []WidgetDef
}

type WidgetDef struct {
    Type       string // 如 "graph", "map", "clock"
    Name       string
    Position   WidgetPosition
    Size       WidgetSize
    Config     map[string]interface{} // 版本相关配置
}

type WidgetPosition struct{ X, Y int }
type WidgetSize struct{ Width, Height int }
```

**业务逻辑：**
- V5：
  - 先 `screen.create` 创建聚合图形
  - 再 `screenitem.create` 添加子项，参数包括 `resourcetype` (0=graph, 1=map...), `resourceid`, `x`, `y`, `width`, `height`
- V7：
  - `dashboard.create` 创建仪表盘
  - `dashboard.widget.create` 添加 Widget，参数为 `type`, `name`, `x`, `y`, `width`, `height`, `fields`（数组形式）

---

## 6. 数据模型与版本转换

### 6.1 通用请求模型（Version-Agnostic）

```go
// 内部统一模型，与 Zabbix 版本无关
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

type UnifiedUserGroup struct {
    ID          string
    Name        string
    Enabled     bool
    UserCount   int
    UserIDs     []string
    Rights      []GroupRight
}

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
    Status     int // 0=enabled, 1=disabled
}

type UnifiedHostGroup struct {
    ID        string
    Name      string
    HostCount int
    HostIDs   []string
}

type UnifiedDashboard struct {
    ID       string
    Name     string
    OwnerID  string
    OwnerName string
    Pages    []DashboardPage
    Private  bool
    Type     string // "screen"(5.x) / "dashboard"(7.x)
}

type UnifiedRole struct {
    ID   string
    Name string
    Type int // 1=用户, 2=管理员, 3=超级管理员
}
```

### 6.2 V5 → V7 转换器示例

```go
// ScreenItem → DashboardWidget 转换
type ScreenToDashboardConverter struct{}

func (c *ScreenToDashboardConverter) Convert(screenItems []ScreenItem) []DashboardWidget {
    widgets := make([]DashboardWidget, 0, len(screenItems))
    for _, item := range screenItems {
        w := DashboardWidget{
            Type:   c.mapResourceType(item.ResourceType), // 0→graph, 1→map...
            X:      item.X,
            Y:      item.Y,
            Width:  item.Width,
            Height: item.Height,
        }
        // V7 widget fields 为数组对象，需转换
        w.Fields = []WidgetField{
            {Type: "graphid", Value: item.ResourceID}, // 假设是图形
        }
        widgets = append(widgets, w)
    }
    return widgets
}
```

---

## 7. 错误处理与日志

### 7.1 错误码设计

```go
const (
    ErrCodeAuthFailed       = "AUTH_001"  // 认证失败
    ErrCodeVersionMismatch  = "VER_001"   // 版本不匹配
    ErrCodeAPINotFound      = "API_001"   // API 方法不存在（版本差异）
    ErrCodeParamInvalid     = "PARAM_001" // 参数校验失败
    ErrCodeResourceExist    = "RES_001"   // 资源已存在
    ErrCodeResourceNotFound = "RES_002"   // 资源不存在
    ErrCodeCloneFailed      = "CLONE_001" // 克隆失败
    ErrCodeUserCancel       = "USER_001"  // 用户取消操作
    ErrCodeVersionNotSupport = "VER_002"  // 当前版本不支持该功能
)

type ZabbixError struct {
    Code    string
    Message string
    Version Version
    Raw     error
}
```

### 7.2 重试策略

```go
type RetryPolicy struct {
    MaxRetries     int
    BaseDelay      time.Duration
    MaxDelay       time.Duration
    RetryableCodes []int // HTTP 状态码：502, 503, 504, 429
}

// 默认策略：指数退避，最多 3 次重试
var DefaultRetryPolicy = RetryPolicy{
    MaxRetries:     3,
    BaseDelay:      500 * time.Millisecond,
    MaxDelay:       5 * time.Second,
    RetryableCodes: []int{502, 503, 504, 429},
}
```

---

## 8. 目录结构

```
zabbix-maint/
├── cmd/
│   └── zbx-cli/              # CLI 入口
│       └── main.go
├── internal/
│   ├── cli/                  # CLI 交互层
│   │   ├── interactive.go    # 交互式菜单主循环
│   │   ├── menu.go           # 菜单树定义
│   │   ├── prompt.go         # 交互式输入组件
│   │   ├── table.go          # 表格渲染
│   │   ├── oneshot.go        # 一次性命令处理
│   │   └── config.go         # 配置管理命令
│   ├── adapter/              # 版本适配器
│   │   ├── v5/
│   │   │   ├── user.go
│   │   │   ├── usergroup.go
│   │   │   ├── host.go
│   │   │   ├── hostgroup.go
│   │   │   ├── screen.go
│   │   │   └── adapter.go
│   │   └── v7/
│   │       ├── user.go
│   │       ├── usergroup.go
│   │       ├── host.go
│   │       ├── hostgroup.go
│   │       ├── dashboard.go
│   │       └── adapter.go
│   ├── api/                  # JSON-RPC 传输层
│   │   ├── client.go
│   │   ├── auth.go
│   │   ├── retry.go
│   │   └── error.go
│   ├── model/                # 统一数据模型
│   │   ├── user.go
│   │   ├── host.go
│   │   ├── dashboard.go
│   │   └── common.go
│   ├── service/              # 业务编排层
│   │   └── batch.go
│   ├── config/               # 配置文件管理
│   │   ├── manager.go
│   │   └── model.go
│   └── version/              # 版本检测与路由
│       ├── detector.go
│       └── router.go
├── pkg/
│   └── zabbix/               # 对外暴露的 SDK 包
│       ├── client.go
│       ├── interface.go
│       └── types.go
├── config/
│   └── example.yaml
├── go.mod
└── README.md
```

---

## 9. 配置与使用

### 9.1 配置文件

配置文件默认存储在 `~/.zbx-cli/config.yaml`：

```yaml
instances:
  - name: "prod-zbx5"
    endpoint: "http://zabbix5.example.com/api_jsonrpc.php"
    version: "5.0"          # 自动检测，也可强制指定
    auth:
      username: "Admin"
      password: "zabbix"

  - name: "prod-zbx7"
    endpoint: "http://zabbix7.example.com/api_jsonrpc.php"
    version: "7.0"
    auth:
      username: "Admin"
      password: "zabbix"

global:
  timeout: 30s
  retry:
    max_retries: 3
    base_delay: 500ms
  log_level: "info"
  default_instance: "prod-zbx5"
```

### 9.2 配置管理命令

```bash
# 添加实例
zbx-cli config add --name prod-zbx5   --url http://zabbix5.example.com/api_jsonrpc.php   --user Admin --pass zabbix

# 列出所有实例
zbx-cli config list

# 删除实例
zbx-cli config remove prod-zbx5

# 测试连接
zbx-cli config test prod-zbx5
```

### 9.3 交互式模式使用示例

```bash
# 进入交互式模式（选择实例）
zbx-cli -i prod-zbx5

# 或使用默认实例
zbx-cli
```

交互流程：

```
═══════════════════════════════════════════════════════════════
  Zabbix 维护工具 v1.0
  当前实例: prod-zbx5 [Zabbix 5.0 LTS]
  连接状态: ✅ 已连接
═══════════════════════════════════════════════════════════════

  请选择功能模块:

  ┌────────────────────────────────────────────┐
  │  1. 用户管理                                │
  │  2. 用户组管理                              │
  │  3. 主机管理                                │
  │  4. 主机组管理                              │
  │  5. 聚合图形 / 仪表盘管理                    │
  │  6. 批量操作                                │
  │  7. 系统信息                                │
  │  0. 退出                                    │
  └────────────────────────────────────────────┘

  请选择 [0-7]: 1

═══════════════════════════════════════════════════════════════
  用户管理
═══════════════════════════════════════════════════════════════

  1. 创建用户
  2. 禁用用户
  3. 启用用户
  4. 修改密码
  5. 删除用户
  6. 查看用户列表
  7. 查看用户详情
  0. 返回上级

  请选择 [0-7]: 1

  → 请输入用户名: zhangsan
  → 请输入姓名: 张三
  → 请输入密码: ********
  → 请选择用户组 (搜索+选择):
      [1] Guests (ID: 7)
      [2] Zabbix administrators (ID: 8)
      [3] Disabled (ID: 9)
    请选择: 1
  → 请选择角色 (仅 7.x 生效, 当前实例为 5.x 自动跳过)

  ⚠️  请确认操作:
      创建用户: zhangsan
      用户组: Guests
      确认执行? [Y/n]: Y

  ✅ 用户创建成功! 用户ID: 10001

  按 Enter 继续...
```

### 9.4 一次性命令模式

```bash
# 用户管理
zbx-cli -i prod-zbx5 user create --alias zhangsan --name 张三 --group 7 --password Temp123!
zbx-cli -i prod-zbx5 user disable --id 10001
zbx-cli -i prod-zbx5 user enable --id 10001
zbx-cli -i prod-zbx5 user passwd --id 10001 --password NewPass456!
zbx-cli -i prod-zbx5 user delete --id 10001
zbx-cli -i prod-zbx5 user list
zbx-cli -i prod-zbx5 user detail --id 10001

# 用户组管理
zbx-cli -i prod-zbx5 usergroup create --name "运维组" --right 2:15
zbx-cli -i prod-zbx5 usergroup disable --id 12
zbx-cli -i prod-zbx5 usergroup clear --id 12
zbx-cli -i prod-zbx5 usergroup delete --id 12
zbx-cli -i prod-zbx5 usergroup list

# 主机管理
zbx-cli -i prod-zbx7 host clone --src 10084 --name web-clone-01 --ip 192.168.1.100
zbx-cli -i prod-zbx5 host list
zbx-cli -i prod-zbx5 host detail --id 10084

# 主机组管理
zbx-cli -i prod-zbx5 hostgroup create --name "Linux Servers"
zbx-cli -i prod-zbx5 hostgroup clear --id 16
zbx-cli -i prod-zbx5 hostgroup delete --id 16
zbx-cli -i prod-zbx5 hostgroup list

# 聚合图形/仪表盘
zbx-cli -i prod-zbx5 screen create --name "Server Overview" --owner 10001
zbx-cli -i prod-zbx5 screen add --id 20001 --type graph --resource 1234 --x 0 --y 0 --w 12 --h 5
zbx-cli -i prod-zbx7 dashboard create --name "Server Overview" --owner 10001
zbx-cli -i prod-zbx7 dashboard add --id 20001 --type graph --graphid 1234 --x 0 --y 0 --w 12 --h 5

# 批量操作
zbx-cli -i prod-zbx5 batch user-create --file users.csv
zbx-cli -i prod-zbx5 batch host-clone --file hosts.csv

# 系统信息
zbx-cli -i prod-zbx5 system version
zbx-cli -i prod-zbx5 system status
```

---

## 10. 关键实现细节

### 10.1 版本自动探测

```go
func DetectVersion(ctx context.Context, client *api.JSONRPCClient) (Version, error) {
    // 方法1: 调用 apiinfo.version（无需认证）
    var version string
    err := client.Call(ctx, "apiinfo.version", nil, &version)
    if err != nil {
        return "", err
    }

    // 解析主版本号
    if strings.HasPrefix(version, "5.") {
        return Version5, nil
    } else if strings.HasPrefix(version, "7.") {
        return Version7, nil
    }
    return "", fmt.Errorf("unsupported version: %s", version)
}
```

### 10.2 V5 主机全克隆手动实现

```go
func (a *V5Adapter) HostFullClone(ctx context.Context, req HostCloneReq) (string, error) {
    // 1. 获取源主机完整配置
    var hosts []map[string]interface{}
    err := a.client.Call(ctx, "host.get", map[string]interface{}{
        "hostids": req.SourceHostID,
        "output":  "extend",
        "selectInterfaces":   "extend",
        "selectItems":        "extend",
        "selectTriggers":     "extend",
        "selectGraphs":       "extend",
        "selectMacros":       "extend",
        "selectInventory":    "extend",
        "selectParentTemplates": "extend",
        "selectGroups":       "extend",
    }, &hosts)
    if err != nil || len(hosts) == 0 {
        return "", fmt.Errorf("get source host failed: %w", err)
    }
    src := hosts[0]

    // 2. 构造新主机（移除 ID，修改名称）
    newHost := deepCopy(src)
    delete(newHost, "hostid")
    delete(newHost, "items")
    delete(newHost, "triggers")
    delete(newHost, "graphs")
    newHost["host"] = req.NewHostName
    newHost["name"] = req.NewHostName

    // 3. 创建主机
    var created map[string]interface{}
    err = a.client.Call(ctx, "host.create", newHost, &created)
    if err != nil {
        return "", fmt.Errorf("create host failed: %w", err)
    }
    newHostID := created["hostids"].([]interface{})[0].(string)

    // 4. 复制 Items（替换 hostid）
    if req.CloneItems {
        if items, ok := src["items"].([]interface{}); ok && len(items) > 0 {
            if err := a.cloneItems(ctx, newHostID, items); err != nil {
                // 事务补偿：删除已创建的主机
                a.client.Call(ctx, "host.delete", []string{newHostID}, nil)
                return "", fmt.Errorf("clone items failed: %w", err)
            }
        }
    }

    // 5. 复制 Triggers、Graphs、Interfaces...
    // ... 类似逻辑

    return newHostID, nil
}
```

### 10.3 V7 主机全克隆（原生支持）

```go
func (a *V7Adapter) HostFullClone(ctx context.Context, req HostCloneReq) (string, error) {
    params := map[string]interface{}{
        "hostid": req.SourceHostID,
        "host": map[string]interface{}{
            "host": req.NewHostName,
            "name": req.NewHostName,
        },
    }
    var result map[string]interface{}
    err := a.client.Call(ctx, "host.clone", params, &result)
    if err != nil {
        return "", err
    }
    return result["hostids"].([]interface{})[0].(string), nil
}
```

### 10.4 交互式菜单主循环

```go
func RunInteractiveMode(ctx context.Context, instanceName string) error {
    // 1. 加载配置并连接
    cfg, err := config.Load(instanceName)
    if err != nil { return err }

    client, version, err := connect(ctx, cfg)
    if err != nil { return err }

    // 2. 构建菜单树
    tree := BuildMenuTree()

    // 3. 主循环
    for {
        // 清屏 + 打印标题
        clearScreen()
        printHeader(cfg, version)

        // 打印当前菜单
        printMenu(tree.Current)

        // 读取用户选择
        choice, err := readChoice(len(tree.Current.SubMenus))
        if err != nil { continue }

        if choice == 0 {
            // 返回上级或退出
            if tree.Current.Parent == nil {
                fmt.Println("再见!")
                return nil
            }
            tree.Current = tree.Current.Parent
            continue
        }

        selected := tree.Current.SubMenus[choice-1]

        if len(selected.SubMenus) > 0 {
            // 进入子菜单
            tree.Current = selected
        } else if selected.Action != nil {
            // 执行叶子节点动作
            clearScreen()
            printHeader(cfg, version)
            fmt.Printf("\n【%s】\n\n", selected.Title)

            if err := selected.Action(ctx, client, version); err != nil {
                if zerr, ok := err.(*ZabbixError); ok && zerr.Code == ErrCodeUserCancel {
                    fmt.Println("⚠️  操作已取消")
                } else {
                    fmt.Printf("❌ 操作失败: %v\n", err)
                }
            }

            fmt.Println("\n按 Enter 继续...")
            bufio.NewReader(os.Stdin).ReadBytes('\n')
        }
    }
}
```

---

## 11. 测试策略

| 测试层级 | 内容 | 工具 |
|---------|------|------|
| **单元测试** | Adapter 参数转换、版本路由、错误码映射、菜单树遍历 | `go test` + `testify` |
| **集成测试** | 对接真实 Zabbix 5.0/7.0 容器 | `docker-compose` + `testcontainers-go` |
| **契约测试** | 验证 5.x/7.x API 响应结构 | `zabbix/api` 包固定 schema |
| **E2E 测试** | 完整业务流程：创建→克隆→禁用→删除 | `go test` + 真实环境 |
| **CLI 测试** | 交互式菜单模拟、命令行参数解析 | `go test` + 模拟输入输出 |

### 11.1 Docker Compose 测试环境

```yaml
version: '3.8'
services:
  zabbix-db-5:
    image: mysql:5.7
    environment:
      MYSQL_ROOT_PASSWORD: root_pwd
      MYSQL_DATABASE: zabbix
      MYSQL_USER: zabbix
      MYSQL_PASSWORD: zabbix_pwd

  zabbix-server-5:
    image: zabbix/zabbix-server-mysql:5.0-alpine
    environment:
      DB_SERVER_HOST: zabbix-db-5
      MYSQL_DATABASE: zabbix
      MYSQL_USER: zabbix
      MYSQL_PASSWORD: zabbix_pwd
    depends_on:
      - zabbix-db-5

  zabbix-web-5:
    image: zabbix/zabbix-web-nginx-mysql:5.0-alpine
    ports:
      - "8085:8080"
    depends_on:
      - zabbix-server-5

  zabbix-db-7:
    image: mysql:8.0
    environment:
      MYSQL_ROOT_PASSWORD: root_pwd
      MYSQL_DATABASE: zabbix
      MYSQL_USER: zabbix
      MYSQL_PASSWORD: zabbix_pwd

  zabbix-server-7:
    image: zabbix/zabbix-server-mysql:7.0-alpine
    environment:
      DB_SERVER_HOST: zabbix-db-7
      MYSQL_DATABASE: zabbix
      MYSQL_USER: zabbix
      MYSQL_PASSWORD: zabbix_pwd
    depends_on:
      - zabbix-db-7

  zabbix-web-7:
    image: zabbix/zabbix-web-nginx-mysql:7.0-alpine
    ports:
      - "8087:8080"
    depends_on:
      - zabbix-server-7
```

---

## 12. 风险与应对

| 风险 | 影响 | 应对措施 |
|------|------|---------|
| Zabbix 5.x 与 7.x API 字段差异过大 | 兼容层复杂 | 建立字段映射表 + 单元测试覆盖 |
| V5 主机全克隆非原子操作 | 部分失败导致脏数据 | 实现事务补偿：失败时回滚已创建资源 |
| 7.x Role 机制变化 | 用户禁用逻辑不兼容 | 首次连接时缓存系统 Role 列表（含 Disabled） |
| Screen → Dashboard 数据结构差异 | 聚合图形无法 100% 兼容 | 定义最小公共子集，不支持的功能返回 `ErrNotSupported` |
| 性能问题（V5 全克隆大量 Items） | API 调用次数多，耗时长 | 批量创建 API（`items.create` 支持数组）+ 并发控制 |
| 交互式 CLI 终端兼容性 | 不同终端显示异常 | 使用 `termenv` 检测终端能力，降级为纯文本模式 |
| 密码明文存储 | 安全风险 | 配置文件密码使用 AES 加密，支持密钥环存储 |

---

## 13. 附录：API 速查表

### 13.1 用户相关

| 操作 | 5.x 方法 | 7.x 方法 | 关键参数差异 |
|------|---------|---------|------------|
| 创建 | `user.create` | `user.create` | 7.x 需 `roleid` |
| 禁用 | `user.update` (status=1) | `user.update` (roleid=disabled) | 字段名不同 |
| 启用 | `user.update` (status=0) | `user.update` (roleid=default) | 字段名不同 |
| 改密 | `user.update` (passwd) | `user.update` (passwd) | 一致 |
| 删除 | `user.delete` | `user.delete` | 一致 |

### 13.2 主机相关

| 操作 | 5.x 方法 | 7.x 方法 | 关键参数差异 |
|------|---------|---------|------------|
| 全克隆 | 手动实现 | `host.clone` | 5.x 无原生 API |
| 创建 | `host.create` | `host.create` | 基本一致 |
| 删除 | `host.delete` | `host.delete` | 一致 |

### 13.3 聚合图形/仪表盘

| 操作 | 5.x 方法 | 7.x 方法 | 关键参数差异 |
|------|---------|---------|------------|
| 创建 | `screen.create` | `dashboard.create` | 对象结构完全不同 |
| 添加项 | `screenitem.create` | `dashboard.widget.create` | V5: resourcetype+resourceid；V7: type+fields |

---

> **文档维护**: 每次 Zabbix 发布新版本时，需更新 `internal/adapter/` 下的适配器实现，并同步更新本文档的 API 映射表。
