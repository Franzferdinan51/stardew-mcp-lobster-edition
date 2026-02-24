#!/bin/bash
# Stardew MCP Server - OpenClaw Gateway Mode
# This connects to OpenClaw Gateway as a tool provider

SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
cd "$SCRIPT_DIR/../mcp-server"

echo "========================================"
echo "Stardew MCP - OpenClaw Gateway Mode"
echo "========================================"
echo ""

# Check if executable exists
if [ ! -f "stardew-mcp" ]; then
    echo "ERROR: stardew-mcp not found"
    echo "Please run setup.sh first to build the server"
    exit 1
fi

# Get token from environment if set
OC_TOKEN="${OPENCLAW_GATEWAY_TOKEN:-}"

echo "Starting Stardew MCP with OpenClaw Gateway..."
echo ""

# Run in OpenClaw Gateway mode
./stardew-mcp -openclaw -openclaw-url "ws://127.0.0.1:18789" -openclaw-token "$OC_TOKEN"

if [ $? -ne 0 ]; then
    echo ""
    echo "Server exited with error code: $?"
fi
