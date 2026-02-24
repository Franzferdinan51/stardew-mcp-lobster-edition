#!/bin/bash
# Stardew MCP Interactive TUI Installer

SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
cd "$SCRIPT_DIR/../mcp-server"

echo "Starting Interactive Installer..."
echo ""

# Check if Go is available
if ! command -v go &> /dev/null; then
    echo "ERROR: Go is not installed"
    echo "Please install Go 1.23+ first, then run this installer again."
    exit 1
fi

# Download dependencies and run installer
go run installer.go

if [ $? -ne 0 ]; then
    echo ""
    echo "Installer exited with error code: $?"
fi
