@echo off
REM Stardew MCP Interactive TUI Installer

echo Starting Interactive Installer...
echo.

cd /d "%~dp0..\mcp-server"

REM Check if Go is available
where go >nul 2>&1
if %errorlevel% neq 0 (
    echo ERROR: Go is not installed
    echo Please install Go 1.23+ first, then run this installer again.
    pause
    exit /b 1
)

REM Download dependencies and run installer
go run installer.go

if %errorlevel% neq 0 (
    echo.
    echo Installer exited with error code: %errorlevel%
    pause
)
