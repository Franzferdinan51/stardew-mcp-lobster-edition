#!/bin/bash
# Stardew MCP Setup Script for Linux/Mac
# This script sets up and builds the Stardew Valley MCP Server

SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
cd "$SCRIPT_DIR/../mcp-server"

echo "========================================"
echo "Stardew MCP - Setup and Startup"
echo "========================================"
echo ""

# Check for Go
if ! command -v go &> /dev/null; then
    echo "ERROR: Go is not installed or not in PATH"
    echo "Please install Go 1.23+ from https://go.dev/dl/"
    exit 1
fi

echo "[1/3] Checking Go version..."
go version
echo ""

echo "[2/3] Installing Go dependencies..."
go mod download
if [ $? -ne 0 ]; then
    echo "ERROR: Failed to download dependencies"
    exit 1
fi
echo ""

echo "[3/3] Building Stardew MCP Server..."
go build -o stardew-mcp
if [ $? -ne 0 ]; then
    echo "ERROR: Build failed"
    exit 1
fi
echo ""

echo "========================================"
echo "Build complete!"
echo "========================================"
echo ""
echo "To run the MCP server:"
echo "  ./run.sh           - Start with autonomous agent"
echo "  ./run.sh -auto=false - Start without autonomous agent"
echo "  ./run.sh -openclaw  - Start with OpenClaw Gateway mode"
echo ""
echo "Make sure Stardew Valley with SMAPI and"
echo "StardewMCP mod is running first!"
echo "========================================"
