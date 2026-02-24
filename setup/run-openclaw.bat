@echo off
REM Stardew MCP Server - OpenClaw Gateway Mode
REM This connects to OpenClaw Gateway as a tool provider

echo ========================================
echo Stardew MCP - OpenClaw Gateway Mode
echo ========================================
echo.

REM Change to mcp-server directory
cd /d "%~dp0..\mcp-server"

REM Check if executable exists
if not exist stardew-mcp.exe (
    echo ERROR: stardew-mcp.exe not found
    echo Please run setup.bat first to build the server
    exit /b 1
)

echo Starting Stardew MCP with OpenClaw Gateway...
echo.

REM Get token from environment if set
set OC_TOKEN=
if defined OPENCLAW_GATEWAY_TOKEN set OC_TOKEN=%OPENCLAW_GATEWAY_TOKEN%

REM Run in OpenClaw Gateway mode
stardew-mcp.exe -openclaw -openclaw-url "ws://127.0.0.1:18789" -openclaw-token "%OC_TOKEN%"

if %errorlevel% neq 0 (
    echo.
    echo Server exited with error code: %errorlevel%
    pause
)
