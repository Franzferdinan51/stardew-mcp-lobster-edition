#!/bin/bash
# Stardew MCP Server - Remote Bot Mode
# This runs the server that remote AI agents can connect to

SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
cd "$SCRIPT_DIR/../mcp-server"

echo "========================================"
echo "Stardew MCP - Remote Bot Server"
echo "========================================"
echo ""

# Check if executable exists
if [ ! -f "stardew-mcp" ]; then
    echo "ERROR: stardew-mcp not found"
    echo "Please run setup.sh first to build the server"
    exit 1
fi

# Get local IP address for display
local_ip=$(hostname -I | awk '{print $1}')

echo "Starting Stardew MCP Server in REMOTE MODE..."
echo ""
echo "Remote bots can connect to:"
echo "  WebSocket: ws://$local_ip:8765/mcp"
echo ""
echo "To connect from another computer, use:"
echo "  ws://YOUR_IP_ADDRESS:8765/mcp"
echo ""
echo "IMPORTANT: Make sure port 8765 is open in your firewall!"
echo ""
echo "Waiting for remote connections..."
echo "(Press Ctrl+C to stop)"
echo ""

# Run in server mode - listens for remote agent connections
./stardew-mcp -server -host "0.0.0.0" -port 8765

if [ $? -ne 0 ]; then
    echo ""
    echo "Server exited with error code: $?"
fi
