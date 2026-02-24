#!/bin/bash
# Stardew MCP Server - TypeScript Version

SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
cd "$SCRIPT_DIR/../mcp-server-ts"

echo "Starting Stardew MCP Server (TypeScript)..."
echo ""

# Check if built
if [ ! -f "dist/index.js" ]; then
    echo "ERROR: Server not built"
    echo "Please run setup-ts.sh first"
    exit 1
fi

# Check if config exists
if [ ! -f "config.yaml" ]; then
    echo "WARNING: config.yaml not found, using defaults"
fi

echo "Connecting to Stardew Valley at ws://localhost:8765/game"
echo ""

node dist/index.js "$@"

if [ $? -ne 0 ]; then
    echo ""
    echo "Server exited with error code: $?"
fi
