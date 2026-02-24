#!/bin/bash
# Stardew MCP - TypeScript Setup Script

SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
cd "$SCRIPT_DIR/../mcp-server-ts"

echo "========================================"
echo "Stardew MCP - Setup (TypeScript)"
echo "========================================"
echo ""

# Check for Node.js
if ! command -v node &> /dev/null; then
    echo "ERROR: Node.js is not installed"
    echo "Please install Node.js from https://nodejs.org/"
    exit 1
fi

echo "[1/3] Checking Node.js version..."
node --version
echo ""

echo "[2/3] Installing dependencies..."
npm install
if [ $? -ne 0 ]; then
    echo "ERROR: Failed to install dependencies"
    exit 1
fi
echo ""

echo "[3/3] Building TypeScript..."
npm run build
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
echo "  ./run-ts.sh           - Start with default settings"
echo "  ./run-ts.sh -server  - Start in remote server mode"
echo ""
echo "Make sure Stardew Valley with SMAPI and"
echo "StardewMCP mod is running first!"
echo "========================================"
