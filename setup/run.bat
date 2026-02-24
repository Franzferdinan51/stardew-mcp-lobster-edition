@echo off
REM Stardew MCP Server Startup Script for Windows
REM Usage: run.bat [options]
REM   -auto=true/false     Enable/disable autonomous mode
REM   -server              Run as remote server
REM   -openclaw            Connect to OpenClaw Gateway
REM   -openclaw-url        OpenClaw Gateway URL
REM   -goal "..."         Set AI goal

echo Starting Stardew MCP Server...
echo.

REM Change to mcp-server directory
cd /d "%~dp0..\mcp-server"

REM Check if executable exists
if not exist stardew-mcp.exe (
    echo ERROR: stardew-mcp.exe not found
    echo Please run setup.bat first to build the server
    exit /b 1
)

REM Check if config exists
if not exist config.yaml (
    echo WARNING: config.yaml not found, using defaults
)

REM Pass all arguments to the server
echo Connecting to Stardew Valley at ws://localhost:8765/game
echo.

stardew-mcp.exe %*

if %errorlevel% neq 0 (
    echo.
    echo Server exited with error code: %errorlevel%
    pause
)
