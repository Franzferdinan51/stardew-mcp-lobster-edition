#!/bin/bash
# Stardew MCP Server Startup Script for Linux/Mac

SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
cd "$SCRIPT_DIR/../mcp-server"

echo "Starting Stardew MCP Server..."
echo ""

# Check if executable exists
if [ ! -f "stardew-mcp" ]; then
    echo "ERROR: stardew-mcp not found"
    echo "Please run setup.sh first to build the server"
    exit 1
fi

# Check if config exists
if [ ! -f "config.yaml" ]; then
    echo "WARNING: config.yaml not found, using defaults"
fi

# Pass all arguments to the server
echo "Connecting to Stardew Valley at ws://localhost:8765/game"
echo ""

./stardew-mcp "$@"

if [ $? -ne 0 ]; then
    echo ""
    echo "Server exited with error code: $?"
fi
