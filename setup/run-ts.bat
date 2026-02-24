@echo off
REM Stardew MCP Server - TypeScript Version

echo Starting Stardew MCP Server (TypeScript)...
echo.

REM Change to TypeScript server directory
cd /d "%~dp0..\mcp-server-ts"

REM Check if built
if not exist "dist\index.js" (
    echo ERROR: Server not built
    echo Please run setup-ts.bat first
    exit /b 1
)

REM Check if config exists
if not exist "config.yaml" (
    echo WARNING: config.yaml not found, using defaults
)

echo Connecting to Stardew Valley at ws://localhost:8765/game
echo.

node dist/index.js %*

if %errorlevel% neq 0 (
    echo.
    echo Server exited with error code: %errorlevel%
    pause
)
