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

# ============================================================
# Script Config
# ============================================================
$ProjectName = "zbx-cli"
$ModuleName = "zabbix-maint"
$CmdPath = "./cmd/$ProjectName"
$OutputDir = "./build"
$LdFlagsBase = "-s -w"

# ============================================================
# Help
# ============================================================
if ($Help) {
    Write-Host ""
    Write-Host "Zabbix CLI Build Script"
    Write-Host ""
    Write-Host "Usage: .\build.ps1 [options]"
    Write-Host ""
    Write-Host "Options:"
    Write-Host "  -Version <string>   Set version (e.g. 1.0.0)"
    Write-Host "  -Commit  <string>   Set Git commit hash"
    Write-Host "  -BuildTime <string> Set build time"
    Write-Host "  -Clean              Clean build directory"
    Write-Host "  -Test               Run tests before build"
    Write-Host "  -Windows            Build Windows only"
    Write-Host "  -Linux              Build Linux only (cross-compile)"
    Write-Host "  -All                Build all platforms (default)"
    Write-Host "  -Help               Show this help"
    Write-Host ""
    Write-Host "Examples:"
    Write-Host "  .\build.ps1                  # Build all"
    Write-Host "  .\build.ps1 -Windows         # Build Windows only"
    Write-Host "  .\build.ps1 -Test -Linux     # Test then build Linux"
    Write-Host "  .\build.ps1 -Version 1.0.0   # Build with version"
    Write-Host "  .\build.ps1 -Clean           # Clean build dir"
    Write-Host ""
    exit 0
}

# ============================================================
# Utility Functions
# ============================================================
function Write-Header($text) {
    Write-Host ""
    Write-Host "============================================================"
    Write-Host "  $text"
    Write-Host "============================================================"
}

function Write-Step($text) {
    Write-Host "  >> $text"
}

function Write-Success($text) {
    Write-Host "  [OK] $text" -ForegroundColor Green
}

function Write-Fail($text) {
    Write-Host "  [FAIL] $text" -ForegroundColor Red
}

function Write-Info($text) {
    Write-Host "  [INFO] $text" -ForegroundColor Gray
}

function Test-GoInstalled {
    try {
        $goVersion = go version 2>$null
        if ($goVersion) {
            Write-Info "Go version: $goVersion"
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

    Write-Step "Build [$GOOS/$GOARCH] => $OutputName"

    $env:GOOS = $GOOS
    $env:GOARCH = $GOARCH
    $env:CGO_ENABLED = "0"

    $cmd = "go build -trimpath -ldflags `"$LdFlagsBase $LdFlags`" -o `"$OutputDir/$OutputName`" $CmdPath"

    try {
        Invoke-Expression $cmd
        Write-Success "Build OK: $OutputName"
    } catch {
        Write-Fail "Build failed: $OutputName"
        throw $_
    }
}

# ============================================================
# Main
# ============================================================
Write-Header "Zabbix CLI Build Script"

# 1. Check Go
if (-not (Test-GoInstalled)) {
    Write-Fail "Go not found. Please install Go 1.21+"
    exit 1
}

# 2. Clean mode
if ($Clean) {
    Write-Step "Cleaning build directory..."
    if (Test-Path $OutputDir) {
        Remove-Item -Recurse -Force $OutputDir
        Write-Success "Build directory cleaned"
    } else {
        Write-Info "Build directory does not exist, nothing to clean"
    }
    exit 0
}

# 3. Run tests
if ($Test) {
    Write-Header "Running Tests"
    Write-Step "go test ./..."
    try {
        go test ./... -v
        Write-Success "All tests passed"
    } catch {
        Write-Fail "Tests failed, build aborted"
        exit 1
    }
}

# 4. Get version info
$ldFlags = Get-VersionInfo
Write-Info "Version info: $ldFlags"

# 5. Create output dir
New-BuildDir

# 6. Determine build targets
$buildWindows = $Windows -or $All -or (-not $Windows -and -not $Linux)
$buildLinux = $Linux -or $All -or (-not $Windows -and -not $Linux)

$buildCount = 0

# 7. Build Windows
if ($buildWindows) {
    Write-Header "Building Windows"

    Invoke-Build -GOOS "windows" -GOARCH "amd64" -OutputName "$ProjectName-windows-amd64.exe" -LdFlags $ldFlags
    $buildCount++

    Invoke-Build -GOOS "windows" -GOARCH "arm64" -OutputName "$ProjectName-windows-arm64.exe" -LdFlags $ldFlags
    $buildCount++

    Write-Success "Windows build done"
}

# 8. Cross-compile Linux
if ($buildLinux) {
    Write-Header "Cross-compiling Linux"

    Invoke-Build -GOOS "linux" -GOARCH "amd64" -OutputName "$ProjectName-linux-amd64" -LdFlags $ldFlags
    $buildCount++

    Invoke-Build -GOOS "linux" -GOARCH "arm64" -OutputName "$ProjectName-linux-arm64" -LdFlags $ldFlags
    $buildCount++

    Write-Success "Linux build done"
}

# 9. Summary
Write-Header "Build Summary"
Write-Info "Targets built: $buildCount"
Write-Info "Output dir: $OutputDir"
Write-Info ""

if (Test-Path $OutputDir) {
    Get-ChildItem $OutputDir | ForEach-Object {
        $size = "{0:N0}" -f $_.Length
        Write-Info "  $($_.Name.PadRight(35)) ($size bytes)"
    }
}

Write-Host ""
Write-Success "Build complete!"
Write-Info ""
Write-Info "Usage examples:"
Write-Info "  Windows: $OutputDir\$ProjectName-windows-amd64.exe -h"
Write-Info "  Linux:   ./$ProjectName-linux-amd64 -h"
Write-Info ""
