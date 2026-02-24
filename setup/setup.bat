@echo off
REM Stardew MCP Setup Script for Windows
REM This script sets up and builds the Stardew Valley MCP Server

echo ========================================
echo Stardew MCP - Setup and Startup
echo ========================================
echo.

REM Check for Go
where go >nul 2>&1
if %errorlevel% neq 0 (
    echo ERROR: Go is not installed or not in PATH
    echo Please install Go 1.23+ from https://go.dev/dl/
    exit /b 1
)

echo [1/3] Checking Go version...
go version
echo.

REM Change to mcp-server directory
cd /d "%~dp0..\mcp-server"

echo [2/3] Installing Go dependencies...
go mod download
if %errorlevel% neq 0 (
    echo ERROR: Failed to download dependencies
    exit /b 1
)
echo.

echo [3/3] Building Stardew MCP Server...
go build -o stardew-mcp.exe
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
echo   run.bat           - Start with autonomous agent
echo   run.bat -auto=false - Start without autonomous agent
echo   run.bat -openclaw  - Start with OpenClaw Gateway mode
echo.
echo Make sure Stardew Valley with SMAPI and
echo StardewMCP mod is running first!
echo ========================================
