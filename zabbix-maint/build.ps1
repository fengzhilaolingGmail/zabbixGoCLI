param(
    [string]$Version = "",
    [string]$Commit = "",
    [string]$BuildTime = "",
    [switch]$Clean = $false,
    [switch]$Test = $false,
    [switch]$Windows = $false,
    [switch]$Linux = $false,
    [switch]$All = $false,
    [switch]$Help = $false
)

$ErrorActionPreference = "Stop"

# ─────────────────────────────────────────────
# 脚本配置
# ─────────────────────────────────────────────
$ProjectName = "zbx-cli"
$ModuleName = "zabbix-maint"
$CmdPath = "./cmd/$ProjectName"
$OutputDir = "./build"
$LdFlagsBase = "-s -w"

# ─────────────────────────────────────────────
# 帮助信息
# ─────────────────────────────────────────────
if ($Help) {
    @"
Zabbix CLI 一键编译脚本

用法:
    .\build.ps1 [选项]

选项:
    -Version <string>   指定版本号 (如: 1.0.0)
    -Commit <string>    指定 Git Commit Hash
    -BuildTime <string> 指定编译时间
    -Clean              清理 build 目录
    -Test               先运行测试再编译
    -Windows            仅编译 Windows 版本
    -Linux              仅编译 Linux 版本 (amd64 + arm64)
    -All                编译所有平台 (默认行为)
    -Help               显示帮助

示例:
    # 编译所有平台 (默认)
    .\build.ps1

    # 仅编译 Windows 版本
    .\build.ps1 -Windows

    # 先运行测试，再编译 Linux 版本
    .\build.ps1 -Test -Linux

    # 带版本号编译所有平台
    .\build.ps1 -Version "1.0.0" -Commit "abc123"

    # 清理构建产物
    .\build.ps1 -Clean
"@ | Write-Host
    exit 0
}

# ─────────────────────────────────────────────
# 工具函数
# ─────────────────────────────────────────────
function Write-Header($text) {
    Write-Host ""
    Write-Host "╔══════════════════════════════════════════════════════════════╗" -ForegroundColor Cyan
    Write-Host "║  $text" -ForegroundColor Cyan -NoNewline
    Write-Host (" " * (60 - $text.Length)) -NoNewline
    Write-Host "║" -ForegroundColor Cyan
    Write-Host "╚══════════════════════════════════════════════════════════════╝" -ForegroundColor Cyan
}

function Write-Step($text) {
    Write-Host "  → $text" -ForegroundColor Yellow
}

function Write-Success($text) {
    Write-Host "  ✅ $text" -ForegroundColor Green
}

function Write-Error($text) {
    Write-Host "  ❌ $text" -ForegroundColor Red
}

function Write-Info($text) {
    Write-Host "  ℹ️  $text" -ForegroundColor Gray
}

function Test-GoInstalled {
    try {
        $goVersion = go version 2>$null
        if ($goVersion) {
            Write-Info "Go 版本: $goVersion"
            return $true
        }
    } catch {
        return $false
    }
    return $false
}

function Get-VersionInfo {
    $v = $Version
    $c = $Commit
    $t = $BuildTime

    if ([string]::IsNullOrEmpty($v)) {
        $v = git describe --tags --always 2>$null
        if ([string]::IsNullOrEmpty($v)) { $v = "dev" }
    }
    if ([string]::IsNullOrEmpty($c)) {
        $c = git rev-parse --short HEAD 2>$null
        if ([string]::IsNullOrEmpty($c)) { $c = "unknown" }
    }
    if ([string]::IsNullOrEmpty($t)) {
        $t = Get-Date -Format "yyyy-MM-ddTHH:mm:ssK"
    }

    return "-X '$ModuleName/internal/version.BuildVersion=$v' -X '$ModuleName/internal/version.BuildCommit=$c' -X '$ModuleName/internal/version.BuildTime=$t'"
}

function New-BuildDir {
    if (-not (Test-Path $OutputDir)) {
        New-Item -ItemType Directory -Path $OutputDir | Out-Null
    }
}

function Invoke-Build {
    param(
        [string]$GOOS,
        [string]$GOARCH,
        [string]$OutputName,
        [string]$LdFlags
    )

    Write-Step "编译 [$GOOS/$GOARCH] => $OutputName"

    $env:GOOS = $GOOS
    $env:GOARCH = $GOARCH
    $env:CGO_ENABLED = "0"

    $cmd = "go build -trimpath -ldflags `"$LdFlagsBase $LdFlags`" -o `"$OutputDir/$OutputName`" $CmdPath"

    try {
        Invoke-Expression $cmd
        Write-Success "编译成功: $OutputName"
    } catch {
        Write-Error "编译失败: $OutputName"
        throw $_
    }
}

# ─────────────────────────────────────────────
# 主流程
# ─────────────────────────────────────────────
Write-Header "Zabbix CLI 一键编译脚本"

# 1. 检查 Go 环境
if (-not (Test-GoInstalled)) {
    Write-Error "未检测到 Go 环境，请先安装 Go 1.21+"
    exit 1
}

# 2. 清理模式
if ($Clean) {
    Write-Step "清理 build 目录..."
    if (Test-Path $OutputDir) {
        Remove-Item -Recurse -Force $OutputDir
        Write-Success "build 目录已清理"
    } else {
        Write-Info "build 目录不存在，无需清理"
    }
    exit 0
}

# 3. 先运行测试（如果指定了 -Test）
if ($Test) {
    Write-Header "运行测试"
    Write-Step "执行 go test ./..."
    try {
        go test ./... -v
        Write-Success "所有测试通过"
    } catch {
        Write-Error "测试失败，编译终止"
        exit 1
    }
}

# 4. 获取版本信息
$ldFlags = Get-VersionInfo
Write-Info "版本信息: $ldFlags"

# 5. 创建输出目录
New-BuildDir

# 6. 确定编译目标
$buildWindows = $Windows -or $All -or (-not $Windows -and -not $Linux)
$buildLinux = $Linux -or $All -or (-not $Windows -and -not $Linux)

$buildCount = 0

# 7. 编译 Windows 版本
if ($buildWindows) {
    Write-Header "编译 Windows 版本"

    # Windows amd64
    Invoke-Build -GOOS "windows" -GOARCH "amd64" -OutputName "$ProjectName-windows-amd64.exe" -LdFlags $ldFlags
    $buildCount++

    # Windows arm64
    Invoke-Build -GOOS "windows" -GOARCH "arm64" -OutputName "$ProjectName-windows-arm64.exe" -LdFlags $ldFlags
    $buildCount++

    Write-Success "Windows 版本编译完成"
}

# 8. 跨平台编译 Linux 版本
if ($buildLinux) {
    Write-Header "跨平台编译 Linux 版本"

    # Linux amd64
    Invoke-Build -GOOS "linux" -GOARCH "amd64" -OutputName "$ProjectName-linux-amd64" -LdFlags $ldFlags
    $buildCount++

    # Linux arm64
    Invoke-Build -GOOS "linux" -GOARCH "arm64" -OutputName "$ProjectName-linux-arm64" -LdFlags $ldFlags
    $buildCount++

    Write-Success "Linux 版本编译完成"
}

# 9. 输出摘要
Write-Header "编译摘要"
Write-Info "编译目标数: $buildCount"
Write-Info "输出目录: $OutputDir"
Write-Info ""

if (Test-Path $OutputDir) {
    Get-ChildItem $OutputDir | ForEach-Object {
        $size = "{0:N0}" -f $_.Length
        Write-Info "  $($_.Name.PadRight(35)) ($size bytes)"
    }
}

Write-Success "编译完成！"
Write-Info ""
Write-Info "使用示例:"
Write-Info "  Windows: $OutputDir\\$ProjectName-windows-amd64.exe -h"
Write-Info "  Linux:   ./$ProjectName-linux-amd64 -h"
Write-Info ""
