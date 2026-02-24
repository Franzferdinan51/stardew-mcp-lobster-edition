@echo off
REM Stardew MCP Server - Remote Bot Mode
REM This runs the server that remote AI agents can connect to

echo ========================================
echo Stardew MCP - Remote Bot Server
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

REM Get local IP address for display
for /f "tokens=2 delims=:" %%a in ('ipconfig ^| findstr /i "IPv4"') do (
    set local_ip=%%a
)
set local_ip=%local_ip:~1%

echo Starting Stardew MCP Server in REMOTE MODE...
echo.
echo Remote bots can connect to:
echo   WebSocket: ws://%local_ip%:8765/mcp
echo.
echo To connect from another computer, use:
echo   ws://YOUR_IP_ADDRESS:8765/mcp
echo.
echo IMPORTANT: Make sure port 8765 is open in your firewall!
echo.
echo Waiting for remote connections...
echo (Press Ctrl+C to stop)
echo.

REM Run in server mode - listens for remote agent connections
REM Pass the port and host configuration
stardew-mcp.exe -server -host "0.0.0.0" -port 8765

if %errorlevel% neq 0 (
    echo.
    echo Server exited with error code: %errorlevel%
    pause
)
