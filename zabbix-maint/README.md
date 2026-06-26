# Zabbix 维护工具

> 版本: v1.1
> 语言: Go
> 支持 Zabbix 版本: 5.x LTS / 7.x LTS

## 简介

基于 Go 语言封装的 Zabbix JSON-RPC API 维护工具，通过**版本兼容层（Version Adapter）**屏蔽 5.x 与 7.x 的接口差异，向上层提供统一、稳定的业务接口。同时提供**交互式 CLI**，支持菜单式操作，降低使用门槛。

## 特性

- **接口统一**: 上层业务只依赖抽象接口，不感知底层版本差异
- **最小侵入**: 兼容层只做参数转换与路由分发，不修改业务语义
- **可扩展**: 新增 6.x/8.x 版本时，只需新增 Adapter 实现
- **交互友好**: 支持交互式菜单 + 命令行参数两种模式
- **防御式编程**: 所有 API 调用均带重试、超时、错误降级

## 安装

```bash
go install github.com/yourname/zabbix-maint/cmd/zbx-cli@latest
```

## 快速开始

### 1. 配置实例

```bash
zbx-cli config add --name prod-zbx5 \
  --url http://zabbix5.example.com/api_jsonrpc.php \
  --user Admin --pass zabbix
```

### 2. 交互式模式

```bash
zbx-cli -i prod-zbx5
```

### 3. 一次性命令模式

```bash
# 创建用户
zbx-cli -i prod-zbx5 user create --alias zhangsan --group 7 --password Temp123!

# 克隆主机
zbx-cli -i prod-zbx7 host clone --src 10084 --name web-clone-01
```

## 项目结构

```
zabbix-maint/
├── cmd/zbx-cli/          # CLI 入口
├── internal/
│   ├── cli/              # CLI 交互层
│   ├── adapter/          # 版本适配器 (v5 / v7)
│   ├── api/              # JSON-RPC 传输层
│   ├── model/            # 统一数据模型
│   ├── service/          # 业务编排层
│   ├── config/           # 配置文件管理
│   └── version/          # 版本检测与路由
├── pkg/zabbix/           # 对外暴露的 SDK 包
└── config/               # 示例配置文件
```

## 支持的功能

| 功能模块 | 5.x | 7.x |
|---------|-----|-----|
| 用户管理 | ✅ | ✅ |
| 用户组管理 | ✅ | ✅ |
| 主机管理 | ✅ | ✅ |
| 主机组管理 | ✅ | ✅ |
| 聚合图形/仪表盘 | ✅ | ✅ |
| 批量操作 | ✅ | ✅ |
| 系统信息 | ✅ | ✅ |

## 版本差异说明

- **Screen → Dashboard**: Zabbix 7.0 废弃 Screen，改用 Dashboard
- **用户状态控制**: 7.0 移除 `status` 字段，改用 Role 控制
- **主机克隆**: 7.0 原生支持 `host.clone`，5.x 需手动实现
- **用户组状态**: 7.0 用户组禁用逻辑与 5.x 不同

## 开发

```bash
# 克隆仓库
git clone https://github.com/yourname/zabbix-maint.git
cd zabbix-maint

# 下载依赖
go mod tidy

# 编译
go build ./cmd/zbx-cli

# 测试
go test ./...
```

## 许可证

MIT
