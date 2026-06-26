# Zabbix 维护工具 (zbx-cli)

> 版本: v1.1  
> 语言: Go  
> 支持 Zabbix 版本: 5.x LTS / 7.x LTS  
> 设计目标: 版本兼容层 + 统一业务接口 + 交互式 CLI

---

## 目录

- [功能概览](#功能概览)
- [架构设计](#架构设计)
- [安装](#安装)
- [快速开始](#快速开始)
- [配置管理](#配置管理)
- [命令参考](#命令参考)
  - [交互式模式](#交互式模式)
  - [一次性命令模式](#一次性命令模式)
- [版本差异说明](#版本差异说明)
- [项目结构](#项目结构)
- [开发指南](#开发指南)
- [常见问题](#常见问题)

---

## 功能概览

| 功能模块 | 说明 | 5.x | 7.x |
|---------|------|-----|-----|
| **用户管理** | 创建/禁用/启用/改密/删除/查看 |  |  |
| **用户组管理** | 创建/禁用/清空/删除/查看 |  |  |
| **主机管理** | 全克隆/查看列表/查看详情 |  |  |
| **主机组管理** | 创建/清空/删除/查看 |  |  |
| **聚合图形/仪表盘** | 创建/添加组件/查看 |  (Screen) |  (Dashboard) |
| **批量操作** | 批量创建用户/批量克隆主机 |  |  |
| **系统信息** | API 版本/服务器状态 |  |  |

**核心特性：**
- **版本兼容层**：自动探测 Zabbix 版本，屏蔽 5.x 与 7.x API 差异
- **交互式菜单**：支持菜单式操作，降低使用门槛
- **一次性命令**：支持命令行参数直接执行，便于脚本集成
- **防御式编程**：所有 API 调用带重试、超时、错误降级
- **可扩展**：新增 6.x/8.x 版本时，只需新增 Adapter 实现

---

## 架构设计

```
---------------------------------------------------------------
                      CLI / Interactive Layer
   -----------------    -----------------    --------------
   |  Interactive   |   |  One-Shot CMD   |   |  Config CMD  |
   |  Menu Mode     |   |  Mode           |   |  Mode        |
   --------+---------   --------+---------   -------+-------
            ---------------------+----------------
---------------------------------------------------------------
                        Business Layer
   ---------   ----------   ---------   ------------------
   | UserMgr | |UserGrpMgr| | HostMgr | | DashboardMgr     |
   ----+-----   -----+-----   ----+----   --------+---------
        -----------+------------+----------
                      |
           +----------+----------+
           |   Unified Interface   |
           |  (ZabbixOperator)     |
           +----------+----------+
---------------------------------------------------------------
               Version Compatibility Layer
   +----------+   +----------+   +--------+
   | V5Adapter|   | V7Adapter|   |V6Adapter|  ...
   ----+------    ----+------    ---+-----
        ------------+-------------
                      |
   +------------------+------------------+
   |         Transport Layer (HTTP Client)                  |
   |  +-------------+  +-------------+  +-------------+     |
   |  | JSON-RPC    |  | Retry/Backoff|  | Auth Manager|     |
   |  | Client      |  | Middleware   |  | (Token/Key) |     |
   |  +-------------+  +-------------+  +-------------+     |
   +--------------------------------------------------------+
```

---

## 安装

### 从源码编译

```bash
# 克隆仓库
git clone https://github.com/yourname/zabbix-maint.git
cd zabbix-maint

# 下载依赖
go mod tidy

# 编译
go build -o zbx-cli ./cmd/zbx-cli

# 安装到 $GOPATH/bin
go install ./cmd/zbx-cli
```

### 环境要求

- Go 1.21+
- Zabbix Server 5.0 LTS 或 7.0 LTS

---

## 快速开始

### 1. 添加 Zabbix 实例配置

```bash
# 添加 Zabbix 5.x 实例
zbx-cli config add --name prod-zbx5 \
  --url http://zabbix5.example.com/api_jsonrpc.php \
  --user Admin \
  --pass zabbix

# 添加 Zabbix 7.x 实例
zbx-cli config add --name prod-zbx7 \
  --url http://zabbix7.example.com/api_jsonrpc.php \
  --user Admin \
  --pass zabbix
```

### 2. 列出已配置的实例

```bash
zbx-cli config list
```

输出示例：

```
NAME           ENDPOINT
---------------------------------------------------
prod-zbx5      http://zabbix5.example.com/api_jsonrpc.php
prod-zbx7      http://zabbix7.example.com/api_jsonrpc.php
```

### 3. 测试连接

```bash
zbx-cli config test prod-zbx5
```

### 4. 启动交互式模式

```bash
# 使用指定实例
zbx-cli -i prod-zbx5

# 或使用默认实例
zbx-cli
```

交互式界面：

```
===============================================================
  Zabbix CLI Tool v1.0
  Instance: prod-zbx5 [Zabbix 5.0 LTS]
  Status:   Connected
===============================================================

  Main Menu:

  1. User Management
  2. User Group Management
  3. Host Management
  4. Host Group Management
  5. Dashboard / Screen Management
  6. Batch Operations
  7. System Info
  0. Back/Exit

  Select [0-7]: 1

===============================================================
  User Management
===============================================================

  1. Create User
  2. Disable User
  3. Enable User
  4. Change Password
  5. Delete User
  6. List Users
  7. User Detail
  0. Back/Exit

  Select [0-7]: 1
  >> Enter username: zhangsan
  >> Enter full name: Zhang San
  >> Enter password: ********
  >> Select user group (multi-select, comma-separated):
      1. Guests (ID: 7)
      2. Zabbix administrators (ID: 8)
      3. Disabled (ID: 9)
    Select: 1
  >> Select role (only for 7.x, auto-skipped for 5.x)

  [CONFIRM] Create user: zhangsan
            Group: Guests
            Confirm? [Y/n]: Y

  [OK] User created! UserID: 10001

  Press Enter to continue...
```

### 5. 使用一次性命令模式

```bash
# 创建用户
zbx-cli -i prod-zbx5 user create \
  --alias zhangsan \
  --name "Zhang San" \
  --group 7 \
  --password "Temp123!"

# 禁用用户
zbx-cli -i prod-zbx5 user disable --id 10001

# 修改密码
zbx-cli -i prod-zbx5 user passwd --id 10001 --password "NewPass456!"

# 克隆主机 (7.x native support)
zbx-cli -i prod-zbx7 host clone \
  --src 10084 \
  --name web-clone-01 \
  --ip 192.168.1.100
```

---

## 配置管理

### 配置文件位置

配置文件默认存储在 `~/.zbx-cli/config.yaml`：

```yaml
instances:
  - name: "prod-zbx5"
    endpoint: "http://zabbix5.example.com/api_jsonrpc.php"
    version: "5.0"          # auto-detect, or force specify
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

### 配置命令

```bash
# 添加实例
zbx-cli config add --name prod-zbx5 \
  --url http://zabbix5.example.com/api_jsonrpc.php \
  --user Admin --pass zabbix

# 列出所有实例
zbx-cli config list

# 删除实例
zbx-cli config remove prod-zbx5

# 测试连接
zbx-cli config test prod-zbx5
```

---

## 命令参考

### 全局选项

| 选项 | 说明 | 示例 |
|------|------|------|
| `-i` | 指定实例名称 | `-i prod-zbx5` |
| `-h` | 显示帮助 | `-h` |
| `-v` | 显示版本 | `-v` |

### 交互式模式

```bash
zbx-cli -i <instance_name>
# 或
zbx-cli                    # use default instance
```

进入交互式菜单后，使用数字键选择菜单项，`0` 返回上级或退出。

### 一次性命令模式

一次性命令需要配合 `-i <instance>` 使用。

#### 用户管理

```bash
# 创建用户
zbx-cli -i <instance> user create \
  --alias <username> \
  --name <fullname> \
  --group <group_id> \
  --password <password> \
  [--role <role_id>]        # only for 7.x

# 禁用用户
zbx-cli -i <instance> user disable --id <user_id>

# 启用用户
zbx-cli -i <instance> user enable --id <user_id>

# 修改密码
zbx-cli -i <instance> user passwd --id <user_id> --password <new_password>

# 删除用户
zbx-cli -i <instance> user delete --id <user_id>

# 查看用户列表
zbx-cli -i <instance> user list

# 查看用户详情
zbx-cli -i <instance> user detail --id <user_id>
```

#### 用户组管理

```bash
# 创建用户组
zbx-cli -i <instance> usergroup create --name <group_name>

# 禁用用户组
zbx-cli -i <instance> usergroup disable --id <group_id>

# 启用用户组
zbx-cli -i <instance> usergroup enable --id <group_id>

# 清空用户组（移除所有用户）
zbx-cli -i <instance> usergroup clear --id <group_id>

# 删除用户组
zbx-cli -i <instance> usergroup delete --id <group_id>

# 查看用户组列表
zbx-cli -i <instance> usergroup list
```

#### 主机管理

```bash
# 主机全克隆
zbx-cli -i <instance> host clone \
  --src <source_host_id> \
  --name <new_host_name> \
  [--ip <new_ip>]

# 查看主机列表
zbx-cli -i <instance> host list

# 查看主机详情
zbx-cli -i <instance> host detail --id <host_id>
```

#### 主机组管理

```bash
# 创建主机组
zbx-cli -i <instance> hostgroup create --name <group_name>

# 清空主机组（移除所有主机）
zbx-cli -i <instance> hostgroup clear --id <group_id>

# 删除主机组
zbx-cli -i <instance> hostgroup delete --id <group_id>

# 查看主机组列表
zbx-cli -i <instance> hostgroup list
```

#### 聚合图形 / 仪表盘管理

```bash
# 创建聚合图形/仪表盘 (V5: Screen, V7: Dashboard)
zbx-cli -i <instance> dashboard create --name <name>

# 查看列表
zbx-cli -i <instance> dashboard list
```

#### 系统信息

```bash
# 查看 API 版本
zbx-cli -i <instance> system version

# 查看服务器状态
zbx-cli -i <instance> system status
```

---

## 版本差异说明

本工具通过**版本适配层**自动处理以下差异：

| 功能 | Zabbix 5.x | Zabbix 7.x |
|------|-----------|-----------|
| **用户禁用** | `user.update` (status=1) | `user.update` (roleid=disabled role) |
| **用户启用** | `user.update` (status=0) | `user.update` (roleid=default role) |
| **用户创建** | `user.create` (usrgrps) | `user.create` (requires roleid) |
| **用户组禁用** | `usergroup.update` (users_status=1) | `usergroup.update` (userdirectoryid=0) |
| **主机组清空** | `hostgroup.massremove` | `hostgroup.update` (clear hosts) |
| **主机全克隆** | manual implementation (get->create->copy items) | native `host.clone` |
| **聚合图形** | `screen.create` / `screenitem.create` | `dashboard.create` / `dashboard.widget.create` |
| **角色管理** | not supported | `role.get` |

**关键点：**
- 工具启动时会自动调用 `apiinfo.version` 探测 Zabbix 版本
- 用户无需关心底层 API 差异，统一使用 `zbx-cli` 命令
- 7.x 特有的功能（如 roleid）在 5.x 实例下会自动跳过或提示

---

## 项目结构

```
zabbix-maint/
├── cmd/
│   └── zbx-cli/              # CLI entry
│       └── main.go
├── internal/
│   ├── cli/                  # CLI interactive layer
│   │   ├── interactive.go    # interactive menu loop
│   │   ├── menu.go           # menu tree definition
│   │   ├── prompt.go         # interactive input components
│   │   ├── table.go          # table renderer
│   │   ├── oneshot.go        # one-shot command handler
│   │   └── config.go         # config management commands
│   ├── adapter/              # version adapters
│   │   ├── v5/               # Zabbix 5.x adapter
│   │   │   ├── user.go       # user management (status field)
│   │   │   ├── usergroup.go  # user group management (users_status)
│   │   │   ├── host.go       # host management (manual clone)
│   │   │   ├── hostgroup.go  # host group management (massremove)
│   │   │   ├── screen.go     # screen (screen API)
│   │   │   └── adapter.go    # V5 adapter assembly
│   │   └── v7/               # Zabbix 7.x adapter
│   │       ├── user.go       # user management (roleid)
│   │       ├── usergroup.go  # user group management (userdirectoryid)
│   │       ├── host.go       # host management (host.clone)
│   │       ├── hostgroup.go  # host group management (update clear)
│   │       ├── dashboard.go  # dashboard (dashboard API)
│   │       └── adapter.go    # V7 adapter assembly
│   ├── api/                  # JSON-RPC transport layer
│   │   ├── client.go         # HTTP JSON-RPC client
│   │   ├── auth.go           # auth manager (token refresh)
│   │   ├── retry.go          # retry strategy (exponential backoff)
│   │   └── error.go          # error code definitions
│   ├── model/                # unified data models (version-agnostic)
│   │   ├── user.go           # UserCreateReq / UnifiedUser
│   │   ├── host.go           # HostCloneReq / UnifiedHost
│   │   ├── dashboard.go      # DashboardCreateReq / UnifiedDashboard
│   │   └── common.go         # UnifiedUserGroup / UnifiedRole
│   ├── service/              # business orchestration
│   │   └── batch.go          # batch operations (CSV + transaction)
│   ├── config/               # config file management
│   │   ├── manager.go        # config CRUD operations
│   │   └── model.go          # Config / InstanceConfig / AuthConfig
│   └── version/              # version detection & routing
│       ├── detector.go       # apiinfo.version detection
│       └── router.go         # adapter factory / version routing
├── pkg/
│   └── zabbix/               # exposed SDK package
│       ├── client.go         # SDK client wrapper
│       ├── interface.go      # ZabbixOperator unified interface
│       └── types.go          # SDK type definitions
├── config/
│   └── example.yaml          # config file example
├── go.mod
├── build.ps1                 # Windows build script (PowerShell)
├── build.bat                 # Windows build script (Batch)
├── Makefile                  # Make build script
└── README.md
```

---

## 开发指南

### 编译

```bash
# 开发调试
go run ./cmd/zbx-cli -i prod-zbx5 user list

# 编译可执行文件
go build -o zbx-cli.exe ./cmd/zbx-cli

# 使用一键编译脚本（Windows）
.\build.ps1          # all platforms
.\build.ps1 -Windows # Windows only
.\build.ps1 -Linux   # cross-compile Linux only
```

### 测试

```bash
# 运行所有测试
go test ./...

# 运行指定包测试
go test ./internal/adapter/v5/...
go test ./internal/api/...
```

### 添加新版本适配器（如 8.x）

1. 在 `internal/adapter/` 下创建 `v8/` 目录
2. 复制 `v7/` 下的文件结构
3. 修改各文件中 `Call` 的 API 方法名和参数
4. 在 `version/router.go` 的 `AdapterFactory` 中注册新版本
5. 在 `version/detector.go` 中添加版本前缀判断

### 添加新的菜单功能

1. 在 `internal/cli/menu.go` 的 `BuildMenuTree()` 中定义菜单节点
2. 实现对应的 `handleXxx` 函数
3. 在 `pkg/zabbix/interface.go` 的 `ZabbixOperator` 中添加业务方法
4. 在 `v5/` 和 `v7/` 适配器中分别实现

---

## 常见问题

### Q: 为什么无法连接到 Zabbix？

A: 请检查：
- 实例 URL 是否正确（需包含 `/api_jsonrpc.php`）
- 用户名密码是否正确
- Zabbix 服务器的 API 是否启用
- 网络是否可达（可先用 `curl` 测试）

```bash
curl -X POST http://zabbix.example.com/api_jsonrpc.php \
  -H "Content-Type: application/json" \
  -d '{"jsonrpc":"2.0","method":"apiinfo.version","params":{},"id":1}'
```

### Q: 7.x 实例下创建用户提示需要 roleid？

A: Zabbix 7.0 移除了 `status` 字段，用户状态通过 Role 控制。创建用户时必须指定角色 ID。

```bash
zbx-cli -i prod-zbx7 user create --alias zhangsan --group 7 --password Temp123! --role 1
```

### Q: 5.x 主机克隆很慢？

A: Zabbix 5.x 没有原生 `host.clone` API，工具需要逐个复制 items、triggers、graphs。如果主机包含大量监控项，建议：
- 使用模板替代直接克隆
- 或升级到 Zabbix 7.x 使用原生克隆

### Q: 配置文件中的密码安全吗？

A: 当前版本密码以明文存储在 `~/.zbx-cli/config.yaml` 中。后续版本计划支持：
- AES 加密存储
- 系统密钥环集成（Windows Credential / macOS Keychain / Linux Secret Service）

### Q: 如何调试 API 请求？

A: 查看 `internal/api/client.go` 中的 `Call` 方法，可以添加日志输出请求和响应内容。

```bash
# 或使用 verbose 模式（如已实现）
zbx-cli -v -i prod-zbx5 user list
```

---

## 许可证

MIT License

---

> **文档维护**: 每次 Zabbix 发布新版本时，需更新 `internal/adapter/` 下的适配器实现，并同步更新本文档的 API 映射表。
