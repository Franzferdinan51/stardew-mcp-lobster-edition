# Stardew Valley MCP Bridge

A hybrid AI-controlled game mod that bridges Stardew Valley with AI assistants via the Model Context Protocol (MCP). Enables autonomous AI agents to control and play Stardew Valley through real-time game state synchronization.

## Architecture

```
┌─────────────────────────────────────────────────────────┐
│ STARDEW VALLEY (Game)                                   │
│   SMAPI Mod (C# .NET 6)                                 │
│     ModEntry → GameStateSerializer                      │
│              → CommandExecutor (w/ Pathfinder)          │
│              → WebSocketServer                          │
└─────────────────────────────────────────────────────────┘
              ↕ ws://localhost:8765/game
┌─────────────────────────────────────────────────────────┐
│ MCP Server (Go)                                         │
│   GameClient: WebSocket connection, state tracking      │
│   StardewAgent: 12 tools, autonomous loop, LLM calls    │
└─────────────────────────────────────────────────────────┘
              ↕ Copilot SDK
┌─────────────────────────────────────────────────────────┐
│ Claude Sonnet (via GitHub Copilot SDK)                  │
└─────────────────────────────────────────────────────────┘
```

## Prerequisites

### For the SMAPI Mod
- [Stardew Valley](https://www.stardewvalley.net/) (game)
- [SMAPI](https://smapi.io/) (Stardew Modding API) v4.0.0+
- .NET 6.0 SDK

### For the MCP Server
- Go 1.23+
- GitHub Copilot access (for Claude Sonnet via Copilot SDK)

## Building

### 1. Build the C# Mod

```bash
cd mod/StardewMCP
dotnet build
```

This compiles `StardewMCP.dll` and places it in `bin/Debug/net6.0/`.

### 2. Build the Go MCP Server

```bash
cd mcp-server
go build -o stardew-mcp
```

This produces the `stardew-mcp` executable.

## Installation

### Install the SMAPI Mod

1. Build the mod (see above)
2. Copy the entire `mod/StardewMCP/bin/Debug/net6.0/` folder to your Stardew Valley mods directory:
   - **Windows**: `%AppData%\StardewValley\Mods\StardewMCP\`
   - **macOS**: `~/.config/StardewValley/Mods/StardewMCP/`
   - **Linux**: `~/.config/StardewValley/Mods/StardewMCP/`
3. Copy `manifest.json` to the same folder

Or use SMAPI's mod folder structure:
```
Mods/
└── StardewMCP/
    ├── manifest.json
    ├── StardewMCP.dll
    └── (other build outputs)
```

## Usage

### 1. Start Stardew Valley with SMAPI

Launch the game through SMAPI. The mod will automatically start a WebSocket server on `ws://localhost:8765/game`.

### 2. Load a Save File

The mod activates once you load into a game save.

### 3. Run the MCP Server

```bash
cd mcp-server
./stardew-mcp
```

The server connects to the game via WebSocket and begins the autonomous AI agent loop.

## Available AI Tools

The AI agent has access to 12 tools for controlling the game:

| Tool | Description |
|------|-------------|
| `move_to` | Navigate to specific coordinates using A* pathfinding |
| `get_surroundings` | Get current game state and 61x61 ASCII map vision |
| `interact` | Interact with objects/NPCs at a position |
| `use_tool` | Use the currently equipped tool |
| `use_tool_repeat` | Use tool multiple times in succession |
| `face_direction` | Turn player to face a direction |
| `select_item` | Select an item from inventory by name |
| `switch_tool` | Switch to a specific tool |
| `eat_item` | Consume a food item for energy |
| `enter_door` | Enter a building or warp point |
| `find_best_target` | Find optimal target for current tool |
| `clear_target` | Clear the current target |

## WebSocket Protocol

The mod and server communicate via JSON over WebSocket.

**Request format**:
```json
{
  "id": "uuid",
  "type": "command",
  "action": "move_to",
  "params": {"x": 10, "y": 20}
}
```

**Response format**:
```json
{
  "id": "uuid",
  "type": "response",
  "success": true,
  "message": "Moved to position",
  "data": {}
}
```

## Configuration

- **WebSocket Port**: Default `8765` (configured in `WebSocketServer.cs`)
- **Tool Cooldown**: 30 game ticks between tool swings (~0.5s at 60fps)
- **State Broadcast**: Game state sent every 1 second
- **Pathfinding**: A* with 50,000 iteration limit, 30-tile scan radius

## Troubleshooting

**Mod not loading**: Ensure SMAPI 4.0.0+ is installed and the mod files are in the correct Mods folder structure.

**WebSocket connection failed**: Check that the game is running and a save is loaded. The server retries connection every 5 seconds.

**Pathfinding failures**: The A* algorithm attempts up to 5 path recalculations. Some areas may be unreachable due to obstacles.
