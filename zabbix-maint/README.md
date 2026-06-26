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
| **用户管理** | 创建/禁用/启用/改密/删除/查看 | ✅ | ✅ |
| **用户组管理** | 创建/禁用/清空/删除/查看 | ✅ | ✅ |
| **主机管理** | 全克隆/查看列表/查看详情 | ✅ | ✅ |
| **主机组管理** | 创建/清空/删除/查看 | ✅ | ✅ |
| **聚合图形/仪表盘** | 创建/添加组件/查看 | ✅(Screen) | ✅(Dashboard) |
| **批量操作** | 批量创建用户/批量克隆主机 | ✅ | ✅ |
| **系统信息** | API 版本/服务器状态 | ✅ | ✅ |

**核心特性：**
- **版本兼容层**：自动探测 Zabbix 版本，屏蔽 5.x 与 7.x API 差异
- **交互式菜单**：支持菜单式操作，降低使用门槛
- **一次性命令**：支持命令行参数直接执行，便于脚本集成
- **防御式编程**：所有 API 调用带重试、超时、错误降级
- **可扩展**：新增 6.x/8.x 版本时，只需新增 Adapter 实现

---

## 架构设计

```
┌─────────────────────────────────────────────────────────────┐
│                      CLI / Interactive Layer                  │
│  ┌─────────────────┐  ┌─────────────────┐  ┌──────────────┐ │
│  │  Interactive    │  │  One-Shot CMD   │  │  Config CMD  │ │
│  │  Menu Mode      │  │  Mode           │  │  Mode        │ │
│  └────────┬────────┘  └────────┬────────┘  └──────┬───────┘ │
│           └────────────────────┴──────────────────┘         │
├─────────────────────────────────────────────────────────────┤
│                        Business Layer                        │
│  ┌─────────┐ ┌──────────┐ ┌─────────┐ ┌──────────────────┐ │
│  │ UserMgr │ │UserGrpMgr│ │ HostMgr │ │ DashboardMgr     │ │
│  └────┬────┘ └────┬─────┘ └────┬────┘ └────────┬─────────┘ │
│       └───────────┴────────────┴───────────────┘           │
│                         │                                  │
│              ┌──────────┴──────────┐                      │
│              │   Unified Interface   │                      │
│              │  (ZabbixOperator)     │                      │
│              └──────────┬──────────┘                      │
├─────────────────────────────────────────────────────────────┤
│              Version Compatibility Layer                     │
│  ┌──────────┐   ┌──────────┐   ┌────────┐               │
│  │ V5Adapter│   │ V7Adapter│   │V6Adapter│  ...           │
│  └────┬─────┘   └────┬─────┘   └───┬────┘               │
│       └──────────────┴──────────────┘                     │
│                      │                                       │
│  ┌───────────────────┼──────────────────────────────────┐  │
│  │         Transport Layer (HTTP Client)                  │  │
│  │  ┌─────────────┐  ┌─────────────┐  ┌─────────────┐   │  │
│  │  │ JSON-RPC    │  │ Retry/Backoff│  │ Auth Manager│   │  │
│  │  │ Client      │  │ Middleware   │  │ (Token/Key) │   │  │
│  │  └─────────────┘  └─────────────┘  └─────────────┘   │  │
│  └──────────────────────────────────────────────────────┘  │
└─────────────────────────────────────────────────────────────┘
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
名称           URL                                  版本    默认
─────────────────────────────────────────────────────────────────
prod-zbx5      http://zabbix5.example.com/api_jsonrpc.php  auto    *
prod-zbx7      http://zabbix7.example.com/api_jsonrpc.php  auto    
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
═══════════════════════════════════════════════════════════════
  Zabbix 维护工具 v1.0
  当前实例: prod-zbx5 [Zabbix 5.0 LTS]
  连接状态: ✅ 已连接
═══════════════════════════════════════════════════════════════

  主菜单:

  1. 用户管理
  2. 用户组管理
  3. 主机管理
  4. 主机组管理
  5. 聚合图形 / 仪表盘管理
  6. 批量操作
  7. 系统信息
  0. 退出

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

### 5. 使用一次性命令模式

```bash
# 创建用户
zbx-cli -i prod-zbx5 user create \
  --alias zhangsan \
  --name "张三" \
  --group 7 \
  --password "Temp123!"

# 禁用用户
zbx-cli -i prod-zbx5 user disable --id 10001

# 修改密码
zbx-cli -i prod-zbx5 user passwd --id 10001 --password "NewPass456!"

# 克隆主机 (7.x 原生支持)
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

# 版本探测（无需预先配置）
zbx-cli detect --url http://zabbix.example.com/api_jsonrpc.php
```

---

## 命令参考

### 全局选项

| 选项 | 说明 | 示例 |
|------|------|------|
| `-i, --instance` | 指定实例名称 | `-i prod-zbx5` |
| `-v, --verbose` | 详细输出模式 | `-v` |
| `-h, --help` | 显示帮助 | `-h` |

### 交互式模式

```bash
zbx-cli -i <instance_name>
# 或
zbx-cli                    # 使用默认实例
```

进入交互式菜单后，使用数字键选择菜单项，`0` 返回上级或退出。

### 一次性命令模式

#### 用户管理

```bash
# 创建用户
zbx-cli -i <instance> user create \
  --alias <username> \
  [--name <fullname>] \
  [--group <group_id>] \
  [--password <password>] \
  [--role <role_id>]        # 仅 7.x 生效

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
zbx-cli -i <instance> usergroup create \
  --name <group_name> \
  [--right <permission:hostgroup_id>]  # 例如: --right 2:15 (2=只读, 3=读写)

# 禁用用户组
zbx-cli -i <instance> usergroup disable --id <group_id>

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
  [--ip <new_ip>] \
  [--group <group_ids>] \
  [--items] [--triggers] [--graphs]

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
zbx-cli -i <instance> dashboard create \
  --name <name> \
  [--owner <user_id>]

# 添加 Widget / 子项
zbx-cli -i <instance> dashboard add \
  --id <dashboard_id> \
  --type <widget_type> \
  [--resource <resource_id>] \
  [--graphid <graph_id>] \
  --x <x> --y <y> --w <width> --h <height>

# 查看列表
zbx-cli -i <instance> dashboard list
```

#### 批量操作

```bash
# 批量创建用户（从 CSV 文件）
zbx-cli -i <instance> batch user-create --file users.csv

# 批量克隆主机（从 CSV 文件）
zbx-cli -i <instance> batch host-clone --file hosts.csv
```

**CSV 格式示例（users.csv）：**

```csv
alias,name,password,group_ids,role_id
zhangsan,张三,Temp123!,7,1
lisi,李四,Temp123!,8,1
```

**CSV 格式示例（hosts.csv）：**

```csv
source_host_id,new_host_name,new_ip,group_ids
10084,web-clone-01,192.168.1.100,2
10084,web-clone-02,192.168.1.101,2
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
| **用户禁用** | `user.update` (status=1) | `user.update` (roleid=禁用角色) |
| **用户启用** | `user.update` (status=0) | `user.update` (roleid=默认角色) |
| **用户创建** | `user.create` (usrgrps) | `user.create` (需 roleid) |
| **用户组禁用** | `usergroup.update` (users_status=1) | `usergroup.update` (userdirectoryid=0) |
| **主机组清空** | `hostgroup.massremove` | `hostgroup.update` (清空 hosts) |
| **主机全克隆** | 手动实现（get→create→copy items） | 原生 `host.clone` |
| **聚合图形** | `screen.create` / `screenitem.create` | `dashboard.create` / `dashboard.widget.create` |
| **角色管理** | 不支持 | `role.get` |

**关键点：**
- 工具启动时会自动调用 `apiinfo.version` 探测 Zabbix 版本
- 用户无需关心底层 API 差异，统一使用 `zbx-cli` 命令
- 7.x 特有的功能（如 roleid）在 5.x 实例下会自动跳过或提示

---

## 项目结构

```
zabbix-maint/
├── cmd/
│   └── zbx-cli/              # CLI 入口
│       └── main.go
├── internal/
│   ├── cli/                  # CLI 交互层
│   │   ├── interactive.go    # 交互式菜单主循环
│   │   ├── menu.go           # 菜单树定义
│   │   ├── prompt.go         # 交互式输入组件（String/Password/Select/Confirm）
│   │   ├── table.go          # 表格渲染器
│   │   ├── oneshot.go        # 一次性命令处理
│   │   └── config.go         # 配置管理命令
│   ├── adapter/              # 版本适配器
│   │   ├── v5/               # Zabbix 5.x 适配器
│   │   │   ├── user.go       # 用户管理（status 字段控制）
│   │   │   ├── usergroup.go  # 用户组管理（users_status）
│   │   │   ├── host.go       # 主机管理（手动全克隆）
│   │   │   ├── hostgroup.go  # 主机组管理（massremove）
│   │   │   ├── screen.go     # 聚合图形（screen API）
│   │   │   └── adapter.go    # V5 适配器组装
│   │   └── v7/               # Zabbix 7.x 适配器
│   │       ├── user.go       # 用户管理（roleid 控制）
│   │       ├── usergroup.go  # 用户组管理（userdirectoryid）
│   │       ├── host.go       # 主机管理（host.clone）
│   │       ├── hostgroup.go  # 主机组管理（update 清空）
│   │       ├── dashboard.go  # 仪表盘（dashboard API）
│   │       └── adapter.go    # V7 适配器组装
│   ├── api/                  # JSON-RPC 传输层
│   │   ├── client.go         # HTTP JSON-RPC 客户端
│   │   ├── auth.go           # 认证管理器（Token 刷新）
│   │   ├── retry.go          # 重试策略（指数退避）
│   │   └── error.go          # 错误码定义
│   ├── model/                # 统一数据模型（Version-Agnostic）
│   │   ├── user.go           # UserCreateReq / UnifiedUser
│   │   ├── host.go           # HostCloneReq / UnifiedHost
│   │   ├── dashboard.go      # DashboardCreateReq / UnifiedDashboard
│   │   └── common.go         # UnifiedUserGroup / UnifiedRole
│   ├── service/              # 业务编排层
│   │   └── batch.go          # 批量操作（CSV 解析 + 事务）
│   ├── config/               # 配置文件管理
│   │   ├── manager.go        # 配置 CRUD 操作
│   │   └── model.go          # Config / InstanceConfig / AuthConfig
│   └── version/              # 版本检测与路由
│       ├── detector.go       # apiinfo.version 探测
│       └── router.go         # 适配器工厂 / 版本路由
├── pkg/
│   └── zabbix/               # 对外暴露的 SDK 包
│       ├── client.go         # SDK 客户端封装
│       ├── interface.go      # ZabbixOperator 统一接口
│       └── types.go          # SDK 类型定义
├── config/
│   └── example.yaml          # 配置文件示例
├── go.mod
└── README.md
```

---

## 开发指南

### 编译

```bash
# 开发调试
go run ./cmd/zbx-cli -- -i prod-zbx5 user list

# 编译可执行文件
go build -o zbx-cli.exe ./cmd/zbx-cli
```

### 测试

```bash
# 运行所有测试
go test ./...

# 运行指定包测试
go test ./internal/adapter/v5/...
go test ./internal/api/...

# 集成测试（需要 Docker）
go test ./... -tags=integration
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

A: Zabbix 7.0 移除了 `status` 字段，用户状态通过 Role 控制。创建用户时必须指定角色 ID。可用以下命令查看可用角色：

```bash
zbx-cli -i prod-zbx7 system role-list
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

A: 使用 `-v` 选项开启详细输出：

```bash
zbx-cli -v -i prod-zbx5 user list
```

---

## 许可证

MIT License

---

> **文档维护**: 每次 Zabbix 发布新版本时，需更新 `internal/adapter/` 下的适配器实现，并同步更新本文档的 API 映射表。
