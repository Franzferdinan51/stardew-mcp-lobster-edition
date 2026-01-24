# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Project Overview

Stardew Valley MCP Bridge - A hybrid AI-controlled game mod that bridges Stardew Valley with AI assistants via the Model Context Protocol (MCP). Enables autonomous AI agents to control and play Stardew Valley through:
- A C# SMAPI mod running inside the game
- A Go MCP server that communicates with GitHub Copilot SDK (Claude Sonnet)
- WebSocket-based real-time game state synchronization

## Build Commands

### C# Mod (SMAPI)
```bash
cd mod/StardewMCP
dotnet build                    # Compile to StardewMCP.dll
```

### Go MCP Server
```bash
cd mcp-server
go build -o stardew-mcp         # Compile binary
```

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

### C# Mod Components (`mod/StardewMCP/`)

- **ModEntry.cs**: SMAPI entry point. Initializes WebSocket server, wires components, registers game loop events (Update, OneSecond, SaveLoaded). Broadcasts game state every 1 second, processes commands each frame.

- **CommandExecutor.cs**: Executes game commands (move_to, use_tool, select_item, etc.). Contains A* pathfinding integration, continuous movement processing, tool cooldown tracking (30 ticks between swings). Queue-based command processing - one command per game tick.

- **GameStateSerializer.cs**: Captures complete game state: Player, Time, World, Surroundings. Generates 61x61 ASCII map vision (30-tile scan radius). Serializes NPCs, items, terrain, quests, relationships, skills.

- **WebSocketServer.cs**: Server on `ws://localhost:8765/game`. Message types: "command", "get_state", "ping". Response types: "state", "response", "error", "pong".

- **Pathfinder.cs**: A* algorithm for navigation. 4-directional movement, 50,000 iteration limit, Manhattan distance heuristic. Checks walkability across tiles, objects, terrain features, buildings, furniture, water.

### Go MCP Server Components (`mcp-server/`)

- **main.go**: GameClient WebSocket manager with reconnection logic (5-second retry). GameState struct definitions for all game entities. Keep-alive pings every 15 seconds, 15-second command timeout.

- **copilot_agent.go**: StardewAgent using GitHub Copilot SDK. 12 tools: move_to, get_surroundings, interact, use_tool, use_tool_repeat, face_direction, select_item, switch_tool, eat_item, enter_door, find_best_target, clear_target. Autonomous loop with emergency handling (time/energy checks), 60-second LLM call timeout.

## WebSocket Protocol

**Request**:
```json
{"id": "uuid", "type": "command", "action": "move_to", "params": {"x": 10, "y": 20}}
```

**Response**:
```json
{"id": "uuid", "type": "response", "success": true, "message": "...", "data": {...}}
```

## Key Design Patterns

- **Queue-based command execution**: One command per game frame prevents desync
- **Async state broadcasting**: Game state sent every 1 second, separate from command processing
- **Path recalculation**: Up to 5 attempts if initial pathfinding fails
- **Tool cooldown**: 30-tick gaps between tool swings (0.5s at 60fps)
- **Mutex-protected tool execution**: Prevents concurrent tool usage in agent
