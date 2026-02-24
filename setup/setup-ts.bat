@echo off
REM Stardew MCP - TypeScript Setup Script

echo ========================================
echo Stardew MCP - Setup (TypeScript)
echo ========================================
echo.

REM Check for Node.js
where node >nul 2>&1
if %errorlevel% neq 0 (
    echo ERROR: Node.js is not installed
    echo Please install Node.js from https://nodejs.org/
    exit /b 1
)

echo [1/3] Checking Node.js version...
node --version
echo.

REM Change to TypeScript server directory
cd /d "%~dp0..\mcp-server-ts"

echo [2/3] Installing dependencies...
call npm install
if %errorlevel% neq 0 (
    echo ERROR: Failed to install dependencies
    exit /b 1
)
echo.

echo [3/3] Building TypeScript...
call npm run build
if %errorlevel% neq 0 (
    echo ERROR: Build failed
    exit /b 1
)
echo.

echo ========================================
echo Build complete!
echo ========================================
echo.
echo To run the MCP server:
echo   run-ts.bat           - Start with default settings
echo   run-ts.bat -server  - Start in remote server mode
echo.
echo Make sure Stardew Valley with SMAPI and
echo StardewMCP mod is running first!
echo ========================================
