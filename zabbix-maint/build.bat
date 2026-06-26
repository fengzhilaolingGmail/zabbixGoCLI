@echo off
chcp 65001 >nul
setlocal enabledelayedexpansion

:: Zabbix CLI 一键编译脚本 (Windows 入口)
:: 调用 PowerShell 脚本执行实际编译

cd /d "%~dp0"

if "%~1"=="" (
    echo 正在编译所有平台...
    powershell -ExecutionPolicy Bypass -File "build.ps1" -All
) else (
    powershell -ExecutionPolicy Bypass -File "build.ps1" %*
)

if %errorlevel% neq 0 (
    echo.
    echo 编译失败，请检查错误信息。
    pause
    exit /b 1
)

echo.
echo 编译完成！
pause
