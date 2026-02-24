# 🦞 Stardew Valley MCP Bridge

A hybrid AI-controlled game mod that bridges Stardew Valley with AI assistants via the Model Context Protocol (MCP). Enables autonomous AI agents to control and play Stardew Valley through real-time game state synchronization.

## Quick Start

### Interactive TUI Installer (All OS)
```cmd
cd setup
installer.bat
```

This opens an interactive installer with:
- 🦞 TUI interface
- Auto-detect Stardew Valley location
- Option to enable OpenClaw Gateway
- Option to enable Remote Server
- Auto-builds everything

### Windows (One Command!)
```cmd
cd setup
quick-start.bat
```

This builds everything AND installs the mod automatically!

### Linux/Mac
```bash
cd setup
chmod +x *.sh
./quick-start.sh
```

### Individual Steps (Windows)
```cmd
cd setup
setup.bat          # Build Go server only
install-mod.bat    # Build and install C# mod only
run.bat            # Run the MCP server
```

### Individual Steps (Linux/Mac)
```bash
cd setup
chmod +x *.sh
./setup.sh         # Build Go server only
./install-mod.sh   # Build and install C# mod only
./run.sh           # Run the MCP server
```

## Usage Modes

| Mode | Command | Description |
|------|---------|-------------|
| Autonomous | `run.bat` / `./run.sh` | AI agent plays automatically |
| Manual | `run.bat -auto=false` | Connect without AI, send commands manually |
| Remote | `run-remote.bat` / `./run-remote.sh` | Accept remote agent connections |
| OpenClaw | `run-openclaw.bat` / `./run-openclaw.sh` | Connect to OpenClaw Gateway |

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
│   StardewAgent: 12 tools + 30 cheats, autonomous loop   │
└─────────────────────────────────────────────────────────┘
              ↕ Copilot SDK / OpenClaw Gateway
┌─────────────────────────────────────────────────────────┐
│ Claude Sonnet (via GitHub Copilot SDK)                  │
│ OR                                                      │
│ OpenClaw Agent                                          │
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
- OR OpenClaw Gateway (for OpenClaw integration)

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
./stardew-mcp                    # Run with default autonomous mode
./stardew-mcp -auto=false        # Connect without starting AI agent
./stardew-mcp -goal "Clear the farm and plant parsnips"
./stardew-mcp -url ws://localhost:8765/game  # Custom WebSocket URL
```

The server connects to the game via WebSocket and begins the autonomous AI agent loop.

## OpenClaw Gateway Integration

This server can connect to OpenClaw Gateway as a tool provider, making Stardew Valley accessible to any OpenClaw agent.

### Quick Start with OpenClaw:
```cmd
cd setup
setup.bat
run-openclaw.bat
```

### OpenClaw Gateway Options:
```cmd
stardew-mcp.exe -openclaw                                    # Connect to local Gateway
stardew-mcp.exe -openclaw -openclaw-url ws://host:18789     # Custom Gateway URL
stardew-mcp.exe -openclaw -openclaw-token YOUR_TOKEN        # With authentication
```

### Environment Variables:
- `OPENCLAW_GATEWAY_TOKEN` - Token for Gateway authentication

### Protocol:
The server implements the OpenClaw Gateway protocol:
- **Request**: `{ "type": "req", "id": "...", "method": "...", "params": {...} }`
- **Response**: `{ "type": "res", "id": "...", "ok": true, "payload": {...} }`
- **Tools**: Registered via `tools.register` method

### Available Tools for OpenClaw:
| Tool | Description |
|------|-------------|
| `get_state` | Get complete game state |
| `get_surroundings` | Get detailed surroundings |
| `move_to` | Navigate to coordinates |
| `interact` | Interact with objects |
| `use_tool` | Use current tool |
| `switch_tool` | Switch tool by name |
| `face_direction` | Face direction |
| `cheat_mode_enable` | Enable cheats |
| `cheat_warp` | Teleport |
| `cheat_set_money` | Set gold |

Full tool list in [tools.md](./tools.md).

## Remote Bot Support

You can run the MCP server to accept connections from remote AI agents (even from other computers):

### Host Computer (where Stardew runs):
```cmd
cd setup
setup.bat
run-remote.bat
```

The script displays the IP address to connect to:
```
ws://YOUR_IP:8765/mcp
```

### Remote Bot Connection:
Connect to `ws://HOST_IP:8765/mcp` from any machine.

**Important:** Ensure port 8765 is open in your firewall for remote connections!

## Available AI Tools

The AI agent has access to 12 tools for controlling the game:

| Tool | Description |
|------|-------------|
| `move_to` | Navigate to specific coordinates using A* pathfinding |
| `get_surroundings` | Get current game state and 61x61 ASCII map vision |
| `interact` | Interact with objects/NPCs at a position |
| `use_tool` | Use the currently selected tool |
| `use_tool_repeat` | Use tool multiple times in succession |
| `face_direction` | Turn player to face a direction |
| `select_item` | Select an item from inventory by name |
| `switch_tool` | Switch to a specific tool |
| `eat_item` | Consume a food item for energy |
| `enter_door` | Enter a building or warp point |
| `find_best_target` | Find optimal target for current tool |
| `clear_target` | Clear the current target |

## Cheat Mode

Cheat mode provides instant god-mode capabilities for rapid testing or stress-free gameplay. **Must call `cheat_mode_enable` first** before any other cheat commands work.

### Mode Control
| Tool | Description |
|------|-------------|
| `cheat_mode_enable` | Enable cheat mode (required first) |
| `cheat_mode_disable` | Disable cheat mode and all persistent effects |
| `cheat_time_freeze` | Toggle time freeze on/off |
| `cheat_infinite_energy` | Toggle infinite stamina on/off |

### Teleportation
| Tool | Description |
|------|-------------|
| `cheat_warp` | Teleport to location (Farm, Town, Mountain, Beach, Forest, Mine, Desert) |
| `cheat_mine_warp` | Warp to specific mine level (1-120 = Mines, 121+ = Skull Cavern) |

### Farming Automation
| Tool | Description |
|------|-------------|
| `cheat_clear_debris` | Remove all weeds, stones, twigs, grass |
| `cheat_cut_trees` | Chop all trees, collect wood/hardwood |
| `cheat_mine_rocks` | Mine all rocks/boulders, collect ores |
| `cheat_hoe_all` | Till all diggable tiles |
| `cheat_water_all` | Water all tilled soil |
| `cheat_plant_seeds` | Plant seeds on all empty hoed tiles (requires seedId) |
| `cheat_fertilize_all` | Apply fertilizer to all hoed tiles |
| `cheat_grow_crops` | Instantly grow all crops to harvest-ready |
| `cheat_harvest_all` | Harvest all ready crops |
| `cheat_dig_artifacts` | Dig up all artifact spots |

### Resources & Items
| Tool | Description |
|------|-------------|
| `cheat_set_money` | Set gold amount |
| `cheat_add_item` | Add item by ID |
| `cheat_spawn_ores` | Add ores: copper, iron, gold, iridium, coal |
| `cheat_set_energy` | Restore stamina to max |
| `cheat_set_health` | Restore health to max |

### Social
| Tool | Description |
|------|-------------|
| `cheat_set_friendship` | Set friendship with NPC |
| `cheat_max_all_friendships` | Max out all NPC friendships |
| `cheat_give_gift` | Give gift to NPC instantly |

### Upgrades
| Tool | Description |
|------|-------------|
| `cheat_upgrade_backpack` | Upgrade backpack (12, 24, or 36 slots) |
| `cheat_upgrade_tool` | Upgrade tool (0=Basic to 4=Iridium) |
| `cheat_upgrade_all_tools` | Upgrade all tools to specified level |
| `cheat_unlock_all` | Max everything: backpack, tools, recipes, skills |

### Example: Instant Farm Setup
```
1. cheat_mode_enable
2. cheat_warp Farm
3. cheat_clear_debris, cheat_cut_trees, cheat_mine_rocks
4. cheat_hoe_all
5. cheat_plant_seeds (with seedId "472" for parsnips)
6. cheat_grow_crops
7. cheat_harvest_all
```

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

Edit `mcp-server/config.yaml` to customize:
- Game WebSocket URL
- Auto-start behavior
- Log level
- Remote server settings (host/port)
- OpenClaw Gateway settings

**Command-line options:**
```bash
./stardew-mcp -server              # Run as remote server
./stardew-mcp -host "0.0.0.0"     # Host to bind to
./stardew-mcp -port 8765          # Port to listen on
./stardew-mcp -auto=false         # Disable autonomous mode
./stardew-mcp -goal "your goal"   # Set AI goal
./stardew-mcp -config config.yaml # Custom config file
./stardew-mcp -openclaw           # OpenClaw Gateway mode
./stardew-mcp -openclaw-url      # Custom Gateway URL
./stardew-mcp -openclaw-token    # Gateway token
```

- **WebSocket Port**: Default `8765` (configured in `WebSocketServer.cs`)
- **Tool Cooldown**: 30 game ticks between tool swings (~0.5s at 60fps)
- **State Broadcast**: Game state sent every 1 second
- **Pathfinding**: A* with 50,000 iteration limit, 30-tile scan radius

## Troubleshooting

**Mod not loading**: Ensure SMAPI 4.0.0+ is installed and the mod files are in the correct Mods folder structure.

**WebSocket connection failed**: Check that the game is running and a save is loaded. The server retries connection every 5 seconds.

**Pathfinding failures**: The A* algorithm attempts up to 5 path recalculations. Some areas may be unreachable due to obstacles.

---

🦞 Lobster Edition
