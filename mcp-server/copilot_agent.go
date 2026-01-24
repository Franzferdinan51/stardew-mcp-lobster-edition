package main

import (
	"fmt"
	"log"
	"sort"
	"strings"
	"sync"
	"time"

	copilot "github.com/github/copilot-sdk/go"
)

// Stardew Valley Game Knowledge - Shared with AI agent
const gameKnowledge = `# Stardew Valley AI Agent: High-Intelligence Protocol

## 🧠 CORE LOGIC: PLANNING VS EXECUTION

**1. LONG-TERM PLANNING**: When you receive a goal, think about the sequence of areas you need to clear.
**2. SPATIAL AWARENESS**: Check your surroundings (61x61 map) to find the nearest cluster of targets.
**3. EXECUTION**:
   - Move to a tile NEXT to the target.
   - FACE the target.
   - **CONFIRM** the target is in the "Tile in front" data.
   - Use the **Lowest-Energy** tool required.

## 🗺️ ASCII MAP LEGEND (61x61 Vision)
- @ : YOU (The Player)
- . : BLANK GROUND (Walkable)
- # : WALL / BUILDING / IMPASSABLE (Blocked)
- ~ : WATER (Blocked)
- T : TREE / BUSH (Blocked - Chop with AXE to clear)
- O : OBJECT / STONE / TWIG / WEED (Blocked - Break with Pickaxe/Axe/Scythe)
- C : CROP (Blocked - Do not trample if possible)
- H : HOE DIRT (Walkable)
- " : GRASS (Walkable - Cut with Scythe for 0 energy)
- > : WARP / DOOR / ENTRANCE
- ; : ARTIFACT SPOT (Hoe it!)
- ! : NPC
- M : MONSTER

## 🧭 SPATIAL COORDINATION & PRECISION

- **Coordinates**: X is horizontal (0=left), Y is vertical (0=top).
- **Tool Range**: You can ONLY hit the tile directly in front of you.
- **Distance Rule**: You must be exactly 1 tile away from your target.
  - To hit Target (10, 10): Stand at (9, 10) and face "right", OR at (11, 10) and face "left", etc.
  - **DO NOT** stand on the same tile as the target (10, 10).
  - **DO NOT** use move_to to go TO a target coordinate typed 'O' or 'T'. Use move_to to go to a '.' tile NEXT to it.

## 🛠️ TOOL EFFICIENCY

- **SCYTHE**: Use for weeds/grass. It costs **0 ENERGY**. Highly efficient for cleanup.
- **AXE**: Use for wood/twigs/stumps. Costs energy.
- **PICKAXE**: Use for stones/ore. Costs energy.
- **VERIFICATION**: After using a tool, check if the objectName/terrainType in "Tile in front" has changed to "." (walkable ground). If not, your action failed—do not keep moving, FIX it.

## 🚀 INTELLIGENCE & AUTO-CORRECTION

- **No Path Found?**: The tile you clicked is blocked. Try moving to a tile 1-step away from it.
- **IsMoving Error?**: Movement is now BLOCKING. If a move tool finishes, you are at your destination. Do not issue 10 move commands in a row; wait for each.
- **Cleaning Goals**: Don't just swing randomly. Find a target, move to it, clear it, move to the next.

## 🌙 SURVIVAL & NIGHT

- **2:00 AM** is a hard game-over. You MUST be in bed by **1:00 AM**.
- Farmhouse Entrance is usually around (60, 15) on the standard farm layout, but check surroundings for "FarmHouse" warp.
`

// StardewAgent manages the autonomous AI session using GitHub Copilot SDK
type StardewAgent struct {
	client      *copilot.Client
	session     *copilot.Session
	currentPlan string
	toolMutex   sync.Mutex // Prevents concurrent tool execution
}

// NewStardewAgent creates a new Stardew agent using Copilot SDK
func NewStardewAgent() (*StardewAgent, error) {
	log.Printf("[AGENT] Creating GitHub Copilot SDK agent")

	// Create client with default options
	client := copilot.NewClient(nil)
	if err := client.Start(); err != nil {
		return nil, fmt.Errorf("failed to start copilot client: %w", err)
	}

	return &StardewAgent{
		client: client,
	}, nil
}

func (a *StardewAgent) StartSession(initialGoal string) error {
	log.Printf("[AGENT AGENT] Session started with goal: %s", initialGoal)

	// Define tools inline (matches original implementation pattern)
	moveToTool := copilot.DefineTool("move_to", "Move to a WALKABLE tile. This tool BLOCKS until arrival.",
		func(params MoveToParams, inv copilot.ToolInvocation) (string, error) {
			return a.handleMoveTo(params.X, params.Y)
		})

	getSurroundingsTool := copilot.DefineTool("get_surroundings", "Refresh vision to see 61x61 area coordinates.",
		func(params NoParams, inv copilot.ToolInvocation) (string, error) {
			state := gameClient.GetState()
			if state == nil {
				return "Disconnected", nil
			}
			return a.formatGameStateContext(state), nil
		})

	interactTool := copilot.DefineTool("interact", "Interact with tile in front",
		func(params NoParams, inv copilot.ToolInvocation) (string, error) {
			resp, _ := gameClient.SendCommand("interact", nil)
			return resp.Message, nil
		})

	useToolTool := copilot.DefineTool("use_tool", "Use tool once",
		func(params NoParams, inv copilot.ToolInvocation) (string, error) {
			resp, _ := gameClient.SendCommand("use_tool", nil)
			return resp.Message, nil
		})

	useToolRepeatTool := copilot.DefineTool("use_tool_repeat", "Execute tool multiple times",
		func(params CountParams, inv copilot.ToolInvocation) (string, error) {
			resp, _ := gameClient.SendCommand("use_tool_repeat", map[string]interface{}{"count": params.Count})
			return resp.Message, nil
		})

	faceDirectionTool := copilot.DefineTool("face_direction", "Turn character to face direction",
		func(params DirectionParams, inv copilot.ToolInvocation) (string, error) {
			resp, _ := gameClient.SendCommand("face_direction", map[string]interface{}{"direction": params.Direction})
			return resp.Message, nil
		})

	selectItemTool := copilot.DefineTool("select_item", "Find and equip item by name",
		func(params NameParams, inv copilot.ToolInvocation) (string, error) {
			resp, _ := gameClient.SendCommand("select_item", map[string]interface{}{"name": params.Name})
			return resp.Message, nil
		})

	switchToolTool := copilot.DefineTool("switch_tool", "Equip inventory slot",
		func(params SlotParams, inv copilot.ToolInvocation) (string, error) {
			resp, _ := gameClient.SendCommand("switch_tool", map[string]interface{}{"slot": params.Slot})
			return resp.Message, nil
		})

	eatItemTool := copilot.DefineTool("eat_item", "Eat food from inventory",
		func(params SlotParams, inv copilot.ToolInvocation) (string, error) {
			resp, _ := gameClient.SendCommand("eat_item", map[string]interface{}{"slot": params.Slot})
			return resp.Message, nil
		})

	enterDoorTool := copilot.DefineTool("enter_door", "Enter door/warp point in front of player",
		func(params NoParams, inv copilot.ToolInvocation) (string, error) {
			resp, _ := gameClient.SendCommand("enter_door", nil)
			return resp.Message, nil
		})

	findBestTargetTool := copilot.DefineTool("find_best_target", "Find nearest target of specified type with walkable approach tile",
		func(params TargetTypeParams, inv copilot.ToolInvocation) (string, error) {
			state := gameClient.GetState()
			if state == nil {
				return "Game disconnected", nil
			}
			return a.findBestTarget(state, params.TargetType), nil
		})

	clearTargetTool := copilot.DefineTool("clear_target", "Find and clear the nearest target automatically (does select_item + move_to + face + use_tool in one call)",
		func(params TargetTypeParams, inv copilot.ToolInvocation) (string, error) {
			return a.clearTarget(params.TargetType)
		})

	// Create session with tools
	session, err := a.client.CreateSession(&copilot.SessionConfig{
		Model: "claude-4-5-sonnet",
		SystemMessage: &copilot.SystemMessageConfig{
			Content: gameKnowledge,
		},
		Tools: []copilot.Tool{
			moveToTool, getSurroundingsTool, interactTool, useToolTool,
			useToolRepeatTool, faceDirectionTool, selectItemTool, switchToolTool,
			eatItemTool, enterDoorTool, findBestTargetTool, clearTargetTool,
		},
	})
	if err != nil {
		return fmt.Errorf("failed to create session: %w", err)
	}
	a.session = session

	go a.runAutonomousLoop(initialGoal)
	return nil
}

func (a *StardewAgent) runAutonomousLoop(goal string) {
	a.currentPlan = "Initializing..."
	consecutiveErrors := 0
	lastPosition := ""
	iteration := 0

	log.Printf("[AGENT LOOP] Starting autonomous loop...")

	for {
		iteration++
		log.Printf("[AGENT LOOP] Iteration %d - Getting game state...", iteration)

		state := gameClient.GetState()
		if state == nil {
			log.Printf("[AGENT LOOP] Game state is nil, waiting...")
			time.Sleep(2 * time.Second)
			consecutiveErrors++
			if consecutiveErrors > 10 {
				log.Printf("[AGENT AGENT] Too many connection errors, pausing...")
				time.Sleep(10 * time.Second)
			}
			continue
		}
		consecutiveErrors = 0

		log.Printf("[AGENT LOOP] Got state: %s at (%d,%d), Energy: %.0f, CanMove: %v, IsMoving: %v",
			state.Player.Location, int(state.Player.X), int(state.Player.Y),
			state.Player.Energy, state.Player.CanMove, state.Player.IsMoving)

		// Track if stuck
		currentPosition := fmt.Sprintf("%d,%d", int(state.Player.X), int(state.Player.Y))
		if currentPosition == lastPosition {
			a.currentPlan = "Stuck at same position - finding new approach"
		}
		lastPosition = currentPosition

		// Determine active goal
		activeGoal := goal
		urgency := ""

		if state.Time.TimeOfDay >= 2500 {
			activeGoal = "EMERGENCY: Go to bed NOW! Time is " + state.Time.TimeString
			urgency = "CRITICAL"
		} else if state.Time.TimeOfDay >= 2400 {
			activeGoal = "URGENT: Find your bed and sleep. It's " + state.Time.TimeString
			urgency = "URGENT"
		} else if state.Time.TimeOfDay >= 2200 {
			urgency = "Getting late"
		}

		if state.Player.Energy < 10 {
			activeGoal = "LOW ENERGY: Eat food from inventory OR go to bed immediately!"
			urgency = "LOW ENERGY"
		} else if state.Player.Energy < 30 {
			urgency = "Low energy"
		}

		// Skip if player is busy
		if state.Player.IsMoving {
			time.Sleep(100 * time.Millisecond)
			continue
		}
		if !state.Player.CanMove {
			time.Sleep(100 * time.Millisecond)
			continue
		}

		gameContext := a.formatGameStateContext(state)

		prompt := fmt.Sprintf(`Location: %s | Pos: (%d,%d) | Time: %s | Energy: %.0f/%d
%s

GOAL: %s

%s

Use clear_target to efficiently clear debris. Call it with target_type="debris" to find and clear the nearest debris automatically.`,
			state.Player.Location, int(state.Player.X), int(state.Player.Y),
			state.Time.TimeString, state.Player.Energy, state.Player.MaxEnergy,
			urgency,
			activeGoal,
			gameContext)

		// Send message and wait for response
		log.Printf("[AGENT LOOP] Sending prompt to Copilot CLI (timeout: 60s)...")
		response, err := a.session.SendAndWait(copilot.MessageOptions{
			Prompt: prompt,
		}, 0) // 0 means default 60s timeout
		if err != nil {
			log.Printf("[AGENT AGENT] SendAndWait error: %v", err)
			time.Sleep(5 * time.Second)
			continue
		}
		log.Printf("[AGENT LOOP] Got response from Copilot")

		// Log the response
		if response != nil && response.Data.Content != nil {
			thought := strings.TrimSpace(*response.Data.Content)
			if thought != "" {
				log.Printf("[AGENT THOUGHT] %s", thought)
			}
			if strings.Contains(thought, "PLAN:") {
				parts := strings.SplitN(thought, "PLAN:", 2)
				if len(parts) > 1 {
					planEnd := strings.Index(parts[1], "\n")
					if planEnd > 0 {
						a.currentPlan = strings.TrimSpace(parts[1][:planEnd])
					} else {
						a.currentPlan = strings.TrimSpace(parts[1])
					}
				}
			}
		}

		// Brief pause between iterations (LLM call is the main delay)
		if urgency != "" {
			time.Sleep(100 * time.Millisecond)
		} else {
			time.Sleep(250 * time.Millisecond)
		}
	}
}

// Tool parameter structs
type MoveToParams struct {
	X int `json:"x" jsonschema:"Target tile X coordinate"`
	Y int `json:"y" jsonschema:"Target tile Y coordinate"`
}

type NoParams struct{}

type CountParams struct {
	Count int `json:"count" jsonschema:"Number of times to use the tool (1-100)"`
}

type DirectionParams struct {
	Direction string `json:"direction" jsonschema:"Direction to face (up, down, left, right)"`
}

type NameParams struct {
	Name string `json:"name" jsonschema:"Name of the item to find and select (case-insensitive, partial match)"`
}

type TargetTypeParams struct {
	TargetType string `json:"target_type" jsonschema:"Type of target (debris, tree, crop, npc, warp, any)"`
}

type SlotParams struct {
	Slot int `json:"slot" jsonschema:"Inventory slot number"`
}

// TargetInfo contains all info needed to clear a target
type TargetInfo struct {
	X             int
	Y             int
	Name          string
	RequiredTool  string
	HitsRequired  int
	ApproachX     int
	ApproachY     int
	FaceDirection string
}

// Target represents a potential target for the AI
type Target struct {
	X            int
	Y            int
	Name         string
	Type         string
	RequiredTool string
	HitsRequired int
	Distance     int
}

func (a *StardewAgent) handleMoveTo(x, y int) (string, error) {
	a.toolMutex.Lock()
	defer a.toolMutex.Unlock()
	return a.doMoveTo(x, y)
}

// doMoveTo is the internal movement function (caller must hold toolMutex)
func (a *StardewAgent) doMoveTo(x, y int) (string, error) {
	log.Printf("[AGENT TOOL: move_to] Target: (%d, %d)", x, y)

	state := gameClient.GetState()
	if state == nil {
		return "Game disconnected", nil
	}

	if !state.Player.CanMove {
		return "Player is currently busy. Wait for animation to finish.", nil
	}

	if !a.isTileWalkable(state, x, y) {
		return fmt.Sprintf("Target (%d, %d) is blocked by an obstacle. Choose an adjacent '.' tile instead.", x, y), nil
	}

	resp, err := gameClient.SendCommand("move_to", map[string]interface{}{"x": x, "y": y})
	if err != nil {
		return fmt.Sprintf("Move command failed: %v", err), nil
	}
	if !resp.Success {
		return fmt.Sprintf("Move rejected by game: %s", resp.Message), nil
	}

	timeout := time.After(30 * time.Second)
	ticker := time.NewTicker(200 * time.Millisecond)
	defer ticker.Stop()

	for {
		select {
		case <-timeout:
			return "Movement timed out.", nil
		case <-ticker.C:
			state := gameClient.GetState()
			if state != nil && int(state.Player.X) == x && int(state.Player.Y) == y {
				return "Arrived at destination", nil
			}
			if state != nil && !state.Player.IsMoving {
				return fmt.Sprintf("Stopped at (%d, %d). Check surroundings.", int(state.Player.X), int(state.Player.Y)), nil
			}
		}
	}
}

func (a *StardewAgent) clearTarget(targetType string) (string, error) {
	a.toolMutex.Lock()
	defer a.toolMutex.Unlock()

	log.Printf("[AGENT CLEAR_TARGET] Starting for type: %s", targetType)

	state := gameClient.GetState()
	if state == nil {
		return "Game disconnected", nil
	}

	targetInfo := a.findBestTargetInfo(state, targetType)
	if targetInfo == nil {
		return fmt.Sprintf("No %s targets found nearby.", targetType), nil
	}

	log.Printf("[AGENT CLEAR_TARGET] Found: %s at (%d,%d), tool: %s, hits: %d",
		targetInfo.Name, targetInfo.X, targetInfo.Y, targetInfo.RequiredTool, targetInfo.HitsRequired)

	if targetInfo.RequiredTool != "" {
		log.Printf("[AGENT CLEAR_TARGET] Selecting tool: %s", targetInfo.RequiredTool)
		resp, err := gameClient.SendCommand("select_item", map[string]interface{}{"name": targetInfo.RequiredTool})
		if err != nil || resp == nil {
			return fmt.Sprintf("Failed to select %s: connection error", targetInfo.RequiredTool), nil
		}
		if !resp.Success {
			return fmt.Sprintf("Failed to select %s: %s", targetInfo.RequiredTool, resp.Message), nil
		}
		time.Sleep(50 * time.Millisecond)
	}

	log.Printf("[AGENT CLEAR_TARGET] Moving to approach tile (%d, %d)", targetInfo.ApproachX, targetInfo.ApproachY)
	moveResult, _ := a.doMoveTo(targetInfo.ApproachX, targetInfo.ApproachY)
	if !strings.Contains(moveResult, "Arrived") && !strings.Contains(moveResult, "Stopped") {
		return fmt.Sprintf("Failed to reach approach tile: %s", moveResult), nil
	}

	log.Printf("[AGENT CLEAR_TARGET] Facing: %s", targetInfo.FaceDirection)
	resp, err := gameClient.SendCommand("face_direction", map[string]interface{}{"direction": targetInfo.FaceDirection})
	if err != nil || resp == nil {
		return "Failed to face direction: connection error", nil
	}
	if !resp.Success {
		return fmt.Sprintf("Failed to face %s: %s", targetInfo.FaceDirection, resp.Message), nil
	}
	time.Sleep(50 * time.Millisecond)

	var result string
	if targetInfo.HitsRequired > 1 {
		log.Printf("[AGENT CLEAR_TARGET] Using tool %d times", targetInfo.HitsRequired)
		resp, err = gameClient.SendCommand("use_tool_repeat", map[string]interface{}{"count": targetInfo.HitsRequired})
	} else if targetInfo.HitsRequired == 0 {
		log.Printf("[AGENT CLEAR_TARGET] Interacting (no tool needed)")
		resp, err = gameClient.SendCommand("interact", nil)
	} else {
		log.Printf("[AGENT CLEAR_TARGET] Using tool once")
		resp, err = gameClient.SendCommand("use_tool", nil)
	}

	if err != nil || resp == nil {
		result = "command sent (no response)"
	} else {
		result = resp.Message
	}

	log.Printf("[AGENT CLEAR_TARGET] Done! Result: %s", result)
	return fmt.Sprintf("Cleared %s at (%d,%d): %s", targetInfo.Name, targetInfo.X, targetInfo.Y, result), nil
}

func (a *StardewAgent) findBestTargetInfo(state *GameState, targetType string) *TargetInfo {
	px, py := int(state.Player.X), int(state.Player.Y)
	var targets []Target

	targetTypeLower := strings.ToLower(targetType)

	switch targetTypeLower {
	case "weed", "weeds", "grass", "stone", "stones", "rock", "rocks", "twig", "twigs", "stick", "sticks", "object", "objects":
		targetTypeLower = "debris"
	case "trees", "wood", "log", "logs":
		targetTypeLower = "tree"
	case "crops", "harvest", "vegetables", "fruit":
		targetTypeLower = "crop"
	case "all", "everything", "anything":
		targetTypeLower = "any"
	}

	if targetTypeLower == "debris" || targetTypeLower == "any" {
		for _, obj := range state.Surroundings.NearbyObjects {
			if !obj.IsPassable {
				hitsRequired := obj.HitsRequired
				if hitsRequired == 0 {
					hitsRequired = 1
				}
				t := Target{
					X:            obj.X,
					Y:            obj.Y,
					Name:         obj.DisplayName,
					Type:         "object",
					RequiredTool: obj.RequiredTool,
					HitsRequired: hitsRequired,
					Distance:     abs(obj.X-px) + abs(obj.Y-py),
				}
				if t.RequiredTool == "Scythe" {
					t.Distance -= 100
				}
				targets = append(targets, t)
			}
		}
	}

	if targetTypeLower == "tree" || targetTypeLower == "any" {
		for _, tf := range state.Surroundings.NearbyTerrainFeatures {
			if !tf.IsPassable && (tf.Type == "tree" || tf.Type == "fruit_tree") {
				hitsRequired := tf.HitsRequired
				if hitsRequired == 0 {
					hitsRequired = 10
				}
				targets = append(targets, Target{
					X:            tf.X,
					Y:            tf.Y,
					Name:         tf.Type,
					Type:         "terrain",
					RequiredTool: tf.RequiredTool,
					HitsRequired: hitsRequired,
					Distance:     abs(tf.X-px) + abs(tf.Y-py),
				})
			}
		}
	}

	if targetTypeLower == "crop" || targetTypeLower == "any" {
		for _, tf := range state.Surroundings.NearbyTerrainFeatures {
			if tf.HasCrop && tf.IsReadyForHarvest {
				targets = append(targets, Target{
					X:            tf.X,
					Y:            tf.Y,
					Name:         tf.CropName,
					Type:         "crop",
					RequiredTool: "Scythe",
					HitsRequired: 1,
					Distance:     abs(tf.X-px) + abs(tf.Y-py),
				})
			}
		}
	}

	if len(targets) == 0 {
		return nil
	}

	sort.Slice(targets, func(i, j int) bool {
		return targets[i].Distance < targets[j].Distance
	})

	for _, target := range targets {
		adjacents := []struct {
			x, y      int
			direction string
		}{
			{target.X - 1, target.Y, "right"},
			{target.X + 1, target.Y, "left"},
			{target.X, target.Y - 1, "down"},
			{target.X, target.Y + 1, "up"},
		}

		for _, adj := range adjacents {
			if a.isTileWalkable(state, adj.x, adj.y) {
				return &TargetInfo{
					X:             target.X,
					Y:             target.Y,
					Name:          target.Name,
					RequiredTool:  target.RequiredTool,
					HitsRequired:  target.HitsRequired,
					ApproachX:     adj.x,
					ApproachY:     adj.y,
					FaceDirection: adj.direction,
				}
			}
		}
	}

	return nil
}

func (a *StardewAgent) findBestTarget(state *GameState, targetType string) string {
	px, py := int(state.Player.X), int(state.Player.Y)
	var targets []Target

	targetTypeLower := strings.ToLower(targetType)

	switch targetTypeLower {
	case "weed", "weeds", "grass", "stone", "stones", "rock", "rocks", "twig", "twigs", "stick", "sticks", "object", "objects":
		targetTypeLower = "debris"
	case "trees", "wood", "log", "logs":
		targetTypeLower = "tree"
	case "crops", "harvest", "vegetables", "fruit":
		targetTypeLower = "crop"
	case "npcs", "villager", "villagers", "person", "people":
		targetTypeLower = "npc"
	case "warps", "doors", "exit", "entrance", "portal":
		targetTypeLower = "warp"
	case "all", "everything", "anything":
		targetTypeLower = "any"
	}

	if targetTypeLower == "debris" || targetTypeLower == "any" {
		for _, obj := range state.Surroundings.NearbyObjects {
			if !obj.IsPassable {
				hitsRequired := obj.HitsRequired
				if hitsRequired == 0 {
					hitsRequired = 1
				}
				t := Target{
					X:            obj.X,
					Y:            obj.Y,
					Name:         obj.DisplayName,
					Type:         "object",
					RequiredTool: obj.RequiredTool,
					HitsRequired: hitsRequired,
					Distance:     abs(obj.X-px) + abs(obj.Y-py),
				}
				if t.RequiredTool == "Scythe" {
					t.Distance -= 100
				}
				targets = append(targets, t)
			}
		}
	}

	if targetTypeLower == "tree" || targetTypeLower == "any" {
		for _, tf := range state.Surroundings.NearbyTerrainFeatures {
			if !tf.IsPassable && (tf.Type == "tree" || tf.Type == "fruit_tree") {
				hitsRequired := tf.HitsRequired
				if hitsRequired == 0 {
					hitsRequired = 10
				}
				targets = append(targets, Target{
					X:            tf.X,
					Y:            tf.Y,
					Name:         tf.Type,
					Type:         "terrain",
					RequiredTool: tf.RequiredTool,
					HitsRequired: hitsRequired,
					Distance:     abs(tf.X-px) + abs(tf.Y-py),
				})
			}
		}
	}

	if targetTypeLower == "crop" || targetTypeLower == "any" {
		for _, tf := range state.Surroundings.NearbyTerrainFeatures {
			if tf.HasCrop && tf.IsReadyForHarvest {
				targets = append(targets, Target{
					X:            tf.X,
					Y:            tf.Y,
					Name:         tf.CropName,
					Type:         "crop",
					RequiredTool: "Scythe",
					HitsRequired: 1,
					Distance:     abs(tf.X-px) + abs(tf.Y-py),
				})
			}
		}
	}

	if targetTypeLower == "npc" || targetTypeLower == "any" {
		for _, npc := range state.Surroundings.NearbyNPCs {
			targets = append(targets, Target{
				X:            npc.X,
				Y:            npc.Y,
				Name:         npc.DisplayName,
				Type:         "npc",
				RequiredTool: "",
				HitsRequired: 0,
				Distance:     abs(npc.X-px) + abs(npc.Y-py),
			})
		}
	}

	if targetTypeLower == "warp" || targetTypeLower == "door" || targetTypeLower == "any" {
		for _, warp := range state.Surroundings.WarpPoints {
			targets = append(targets, Target{
				X:            warp.X,
				Y:            warp.Y,
				Name:         warp.TargetLocation,
				Type:         "warp",
				RequiredTool: "",
				HitsRequired: 0,
				Distance:     abs(warp.X-px) + abs(warp.Y-py),
			})
		}
		for _, bldg := range state.Surroundings.NearbyBuildings {
			if bldg.Type == "FarmHouse" || bldg.Type == "Cabin" {
				targets = append(targets, Target{
					X:            bldg.DoorX,
					Y:            bldg.DoorY,
					Name:         bldg.Type + " door",
					Type:         "door",
					RequiredTool: "",
					HitsRequired: 0,
					Distance:     abs(bldg.DoorX-px) + abs(bldg.DoorY-py),
				})
			}
		}
	}

	if len(targets) == 0 {
		return fmt.Sprintf("No targets of type '%s' found nearby.", targetType)
	}

	sort.Slice(targets, func(i, j int) bool {
		return targets[i].Distance < targets[j].Distance
	})

	for _, target := range targets {
		adjacents := []struct {
			x, y      int
			direction string
		}{
			{target.X - 1, target.Y, "right"},
			{target.X + 1, target.Y, "left"},
			{target.X, target.Y - 1, "down"},
			{target.X, target.Y + 1, "up"},
		}

		for _, adj := range adjacents {
			if a.isTileWalkable(state, adj.x, adj.y) {
				toolName := strings.ToLower(target.RequiredTool)
				if toolName == "" {
					toolName = "none"
				}

				finalAction := "use_tool"
				if target.HitsRequired > 1 {
					finalAction = fmt.Sprintf("use_tool_repeat with count=%d", target.HitsRequired)
				} else if target.HitsRequired == 0 {
					finalAction = "interact"
				}

				return fmt.Sprintf(`TARGET: %s at (%d,%d) - Tool: %s - Hits: %d

NOW DO THESE IN ORDER (do NOT call find_best_target again):
Step 1: select_item name="%s"
Step 2: move_to x=%d y=%d
Step 3: face_direction direction="%s"
Step 4: %s`,
					target.Name, target.X, target.Y, target.RequiredTool, target.HitsRequired,
					toolName,
					adj.x, adj.y,
					adj.direction,
					finalAction)
			}
		}
	}

	return fmt.Sprintf("Found %d targets but none have accessible approach tiles. Try moving to a different area.", len(targets))
}

func (a *StardewAgent) isTileWalkable(state *GameState, x, y int) bool {
	radius := 30
	px, py := int(state.Player.X), int(state.Player.Y)
	rx := x - px
	ry := y - py

	if ry < -radius || ry > radius || rx < -radius || rx > radius {
		return false
	}

	if state.Surroundings.AsciiMap == "" {
		return true
	}

	lines := strings.Split(state.Surroundings.AsciiMap, "\n")
	gridY := radius + ry
	gridX := radius + rx

	if gridY < 0 || gridY >= len(lines) {
		return false
	}
	line := lines[gridY]
	if gridX < 0 || gridX >= len(line) {
		return false
	}

	char := line[gridX]
	switch char {
	case '.', '>', 'H', '"', ';', '@':
		return true
	default:
		return false
	}
}

func (a *StardewAgent) formatGameStateContext(state *GameState) string {
	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("Equipped Tool: %s (slot %d)\n", state.Player.CurrentTool, state.Player.CurrentToolIndex))

	tif := state.Surroundings.TileInFront
	sb.WriteString(fmt.Sprintf("\n--- TILE IN FRONT (facing %s) ---\n", state.Player.FacingDirectionName))
	sb.WriteString(fmt.Sprintf("Position: (%d, %d)\n", tif.X, tif.Y))
	if tif.ObjectName != "" {
		toolInfo := ""
		if tif.RequiredTool != "" {
			toolInfo = fmt.Sprintf(" [Use: %s]", tif.RequiredTool)
		}
		sb.WriteString(fmt.Sprintf("Object: %s%s\n", tif.ObjectName, toolInfo))
	}
	if tif.TerrainType != "" {
		toolInfo := ""
		if tif.RequiredTool != "" {
			toolInfo = fmt.Sprintf(" [Use: %s]", tif.RequiredTool)
		}
		sb.WriteString(fmt.Sprintf("Terrain: %s%s\n", tif.TerrainType, toolInfo))
	}
	if tif.NPCName != "" {
		sb.WriteString(fmt.Sprintf("NPC: %s\n", tif.NPCName))
	}
	if tif.IsPassable {
		sb.WriteString("Status: WALKABLE\n")
	} else {
		sb.WriteString("Status: BLOCKED\n")
	}

	if len(state.Surroundings.NearbyBuildings) > 0 {
		sb.WriteString("\n--- BUILDINGS ---\n")
		for _, b := range state.Surroundings.NearbyBuildings {
			sb.WriteString(fmt.Sprintf("- %s: Door at (%d, %d)\n", b.Type, b.DoorX, b.DoorY))
		}
	}

	if len(state.Surroundings.WarpPoints) > 0 {
		sb.WriteString("\n--- WARPS/DOORS ---\n")
		for _, w := range state.Surroundings.WarpPoints {
			sb.WriteString(fmt.Sprintf("- (%d, %d) -> %s\n", w.X, w.Y, w.TargetLocation))
		}
	}

	debrisCount := 0
	scytheTargets := 0
	for _, obj := range state.Surroundings.NearbyObjects {
		if !obj.IsPassable {
			debrisCount++
			if obj.RequiredTool == "Scythe" {
				scytheTargets++
			}
		}
	}
	treeCount := 0
	for _, tf := range state.Surroundings.NearbyTerrainFeatures {
		if tf.Type == "tree" || tf.Type == "fruit_tree" {
			treeCount++
		}
	}

	sb.WriteString("\n--- TARGET SUMMARY ---\n")
	sb.WriteString(fmt.Sprintf("Debris (stones/twigs/weeds): %d (%d use Scythe=0 energy)\n", debrisCount, scytheTargets))
	sb.WriteString(fmt.Sprintf("Trees: %d\n", treeCount))
	sb.WriteString(fmt.Sprintf("NPCs: %d\n", len(state.Surroundings.NearbyNPCs)))

	sb.WriteString("\n--- NEAREST TARGETS (use find_best_target for full list) ---\n")
	shown := 0
	for _, obj := range state.Surroundings.NearbyObjects {
		if shown >= 5 {
			break
		}
		if !obj.IsPassable {
			tool := obj.RequiredTool
			if tool == "" {
				tool = "unknown"
			}
			sb.WriteString(fmt.Sprintf("- %s at (%d, %d) [%s]\n", obj.DisplayName, obj.X, obj.Y, tool))
			shown++
		}
	}

	sb.WriteString("\n--- INVENTORY (food items) ---\n")
	hasFood := false
	for _, item := range state.Player.Inventory {
		if item.Category == "Cooking" || strings.Contains(strings.ToLower(item.Name), "salad") ||
			strings.Contains(strings.ToLower(item.Name), "egg") || strings.Contains(strings.ToLower(item.Name), "milk") {
			sb.WriteString(fmt.Sprintf("- Slot %d: %s (x%d)\n", item.Slot, item.DisplayName, item.Stack))
			hasFood = true
		}
	}
	if !hasFood {
		sb.WriteString("No food items found.\n")
	}

	if state.Surroundings.AsciiMap != "" {
		sb.WriteString("\n--- ASCII MAP (center 21x21 of 61x61) ---\n")
		lines := strings.Split(state.Surroundings.AsciiMap, "\n")
		center := 30
		viewRadius := 10
		for y := center - viewRadius; y <= center+viewRadius && y < len(lines); y++ {
			if y >= 0 && y < len(lines) {
				line := lines[y]
				start := center - viewRadius
				end := center + viewRadius + 1
				if start >= 0 && end <= len(line) {
					sb.WriteString(line[start:end] + "\n")
				} else if len(line) > 0 {
					sb.WriteString(line + "\n")
				}
			}
		}
	}

	return sb.String()
}

// abs returns the absolute value of x
func abs(x int) int {
	if x < 0 {
		return -x
	}
	return x
}
