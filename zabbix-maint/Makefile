# Zabbix CLI 一键编译脚本
# 支持 Windows 本地编译 + 跨平台编译 Linux 版本

.PHONY: all windows linux clean test help

# ── 变量 ──
PROJECT   := zbx-cli
MODULE    := zabbix-maint
CMD       := ./cmd/$(PROJECT)
OUTDIR    := ./build
VERSION   := $(shell git describe --tags --always 2>/dev/null || echo dev)
COMMIT    := $(shell git rev-parse --short HEAD 2>/dev/null || echo unknown)
BUILDTIME := $(shell powershell -Command "Get-Date -Format 'yyyy-MM-ddTHH:mm:ssK'")
LDFLAGS   := -s -w \
		     -X '$(MODULE)/internal/version.BuildVersion=$(VERSION)' \
		     -X '$(MODULE)/internal/version.BuildCommit=$(COMMIT)' \
		     -X '$(MODULE)/internal/version.BuildTime=$(BUILDTIME)'

# ── 默认目标 ──
all: windows linux

help: ## 显示帮助
	@echo.
	@echo Zabbix CLI 编译脚本
	@echo.
	@echo 用法: make [目标]
	@echo.
	@echo 目标:
	@echo   make all       编译所有平台 (Windows + Linux)
	@echo   make windows   仅编译 Windows 版本 (amd64 + arm64)
	@echo   make linux     仅编译 Linux 版本 (amd64 + arm64)
	@echo   make test      先运行测试，再编译所有平台
	@echo   make clean     清理 build 目录
	@echo.

# ── 编译 Windows ──
windows: ## 编译 Windows 版本
	@echo ═══════════════════════════════════════════
	@echo  编译 Windows 版本
	@echo ═══════════════════════════════════════════
	@powershell -Command "if (-not (Test-Path '$(OUTDIR)')) { New-Item -ItemType Directory -Path '$(OUTDIR)' | Out-Null }"
	@echo → Windows amd64...
	@set GOOS=windows&& set GOARCH=amd64&& set CGO_ENABLED=0&& \
		go build -trimpath -ldflags "$(LDFLAGS)" -o "$(OUTDIR)/$(PROJECT)-windows-amd64.exe" $(CMD)
	@echo → Windows arm64...
	@set GOOS=windows&& set GOARCH=arm64&& set CGO_ENABLED=0&& \
		go build -trimpath -ldflags "$(LDFLAGS)" -o "$(OUTDIR)/$(PROJECT)-windows-arm64.exe" $(CMD)
	@echo ✅ Windows 版本编译完成

# ── 跨平台编译 Linux ──
linux: ## 跨平台编译 Linux 版本
	@echo ═══════════════════════════════════════════
	@echo  跨平台编译 Linux 版本
	@echo ═══════════════════════════════════════════
	@powershell -Command "if (-not (Test-Path '$(OUTDIR)')) { New-Item -ItemType Directory -Path '$(OUTDIR)' | Out-Null }"
	@echo → Linux amd64...
	@set GOOS=linux&& set GOARCH=amd64&& set CGO_ENABLED=0&& \
		go build -trimpath -ldflags "$(LDFLAGS)" -o "$(OUTDIR)/$(PROJECT)-linux-amd64" $(CMD)
	@echo → Linux arm64...
	@set GOOS=linux&& set GOARCH=arm64&& set CGO_ENABLED=0&& \
		go build -trimpath -ldflags "$(LDFLAGS)" -o "$(OUTDIR)/$(PROJECT)-linux-arm64" $(CMD)
	@echo ✅ Linux 版本编译完成

# ── 测试 + 编译 ──
test: ## 运行测试后编译
	@echo ═══════════════════════════════════════════
	@echo  运行测试
	@echo ═══════════════════════════════════════════
	go test ./...
	@echo ✅ 测试通过
	@$(MAKE) all

# ── 清理 ──
clean: ## 清理 build 目录
	@echo 清理 build 目录...
	@powershell -Command "if (Test-Path '$(OUTDIR)') { Remove-Item -Recurse -Force '$(OUTDIR)' }"
	@echo ✅ 已清理

# ── 信息 ──
info: ## 显示版本信息
	@echo Version: $(VERSION)
	@echo Commit:  $(COMMIT)
	@echo Time:    $(BUILDTIME)
