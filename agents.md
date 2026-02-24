# Agent Instructions for Stardew MCP

You are an expert Stardew Valley farmer AI agent. Your goal is to efficiently manage the farm, grow profitable crops, maximize profits, and build relationships with villagers.

## Project Structure

```
stardew-mcp/
├── setup/                    # Quick setup scripts
│   ├── setup.bat            # Windows build script
│   ├── setup.sh            # Linux/Mac build script
│   ├── run.bat             # Windows run script
│   ├── run.sh              # Linux/Mac run script
│   ├── openclaw-tools.json # OpenClaw tool definitions
│   └── openclaw-config.json# OpenClaw agent config
├── mcp-server/              # Go MCP Server
│   ├── config.yaml          # Configuration file
│   ├── main.go              # Server entry point
│   ├── copilot_agent.go     # AI agent implementation
│   └── go.mod               # Go dependencies
└── mod/
    └── StardewMCP/           # C# SMAPI mod
        ├── ModEntry.cs
        ├── CommandExecutor.cs
        ├── GameStateSerializer.cs
        ├── Pathfinder.cs
        └── WebSocketServer.cs
```

## Setup Instructions

### Step 1: Build the Go MCP Server

**Windows:**
```cmd
cd setup
setup.bat
```

**Linux/Mac:**
```bash
cd setup
chmod +x *.sh
./setup.sh
```

### Step 2: Install the SMAPI Mod

1. Build the C# mod:
   ```bash
   cd mod/StardewMCP
   dotnet build
   ```

2. Copy `mod/StardewMCP/bin/Debug/net6.0/` contents to:
   - Windows: `%AppData%\StardewValley\Mods\StardewMCP\`
   - Linux/Mac: `~/.config/StardewValley/Mods/StardewMCP/`

3. Ensure `manifest.json` is in the mod folder.

### Step 3: Run Everything

1. **Start Stardew Valley** through SMAPI (the mod auto-starts)
2. **Load a save file** (the mod activates when you enter the world)
3. **Run the MCP server:**

   **Windows:**
   ```cmd
   cd setup
   run.bat
   ```

   **Linux/Mac:**
   ```bash
   cd setup
   ./run.sh
   ```

The AI agent will automatically connect and start farming!

## Game Mechanics to Know

### Time System
- In-game day: 6:00 AM to 2:00 AM (next day)
- Time moves faster: 10 minutes per real-world minute
- Each tile takes about 1 second to walk
- Sleep at 2 AM or when energy depletes

### Energy Management
- Start with 270 energy (upgrades available)
- Crops need watering daily unless using fertilizer
- Eating food restores energy
- Use `cheat_mode_enable` then `cheat_infinite_energy` if needed

### Money-Making Strategies
1. **Farming**: Most profitable - plant seasonal crops
2. **Mining**: Get ores for crafting and selling
3. **Fishing**: Easy early-game money
4. **Foraging**: Gather wild items

### Seasonal Crops (by profit)
- **Spring**: Strawberries (from Festival), cauliflower, potatoes
- **Summer**: Blueberries, corn, sunflowers
- ** Fall**: Pumpkins, cranberries, corn

## Available Tools

### Core Movement & Actions
| Tool | Usage |
|------|-------|
| `move_to(x, y)` | Navigate to coordinates |
| `get_surroundings` | See nearby objects/map |
| `interact` | Talk/pickup/activate |
| `use_tool` | Swing current tool |
| `switch_tool(name)` | Pick tool (hoe, watering_can, etc.) |

### Cheat Mode (Use Judiciously)
**ALWAYS call `cheat_mode_enable` first!**

```python
# Setup a new farm quickly
cheat_mode_enable()
cheat_warp("Farm")
cheat_clear_debris()
cheat_cut_trees()
cheat_mine_rocks()
cheat_hoe_all()
cheat_plant_seeds("parsnips")  # or other seasonal seeds
cheat_grow_crops()
cheat_harvest_all()
cheat_set_money(10000)  # Set some starting gold
```

### Full Cheat Commands
- `cheat_time_freeze(bool)` - Stop time
- `cheat_warp(location)` - Teleport (Farm, Town, Mine, Beach, etc.)
- `cheat_set_money(amount)` - Set gold
- `cheat_upgrade_all_tools()` - Max upgrades
- `cheat_max_all_friendships()` - Max NPC hearts

## ASCII Map Legend

When you see the map, these characters represent terrain:
- `.` = Grass (walkable)
- `:` = Tilled soil
- `~` = Water
- `#` = Wall/fence
- `O` = Object (machine, crop, etc.)
- `T` = Tree
- `R` = Rock
- `@` = Player

## Decision Making

1. **Check time** - Plan activities before 9 PM (sell before 10 PM)
2. **Check energy** - Eat food if low, sleep if exhausted
3. **Prioritize**: Water crops > Harvest > Plant > Mining > Foraging
4. **Check weather** - Rain = no watering needed, good for mining
5. **Festival days** - Attend festivals for exclusive items

## Emergency Protocols

If stuck or need help:
1. Use `get_surroundings` to see where you are
2. Use `cheat_warp("Farm")` to return to safety
3. Check inventory with `get_state` to see items

## Configuration

Edit `mcp-server/config.yaml`:
```yaml
server:
  game_url: "ws://localhost:8765/game"
  auto_start: true
  log_level: "info"

agent:
  default_goal: "Efficiently farm and maximize profits"
  llm_timeout: 60
```

## Troubleshooting

- **"Connection refused"**: Start Stardew Valley + load save first
- **"Timeout"**: Game might be paused or at 2 AM
- **Tools not working**: Use `switch_tool("hoe")` etc. to select tool first

## Have Fun!

Stardew Valley is meant to be enjoyed. Don't stress about optimization - focus on:
- Planting crops you enjoy
- Talking to villagers you like
- Exploring the mines at your own pace
- Completing bundles at your leisure

The farm will grow, one season at a time.
