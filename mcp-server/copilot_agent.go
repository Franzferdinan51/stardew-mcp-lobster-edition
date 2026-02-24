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

// Embedded game knowledge - no external file dependency
const gameKnowledge = `# Stardew Valley AI Agent: High-Intelligence Protocol

## CORE LOGIC: PLANNING VS EXECUTION

**1. LONG-TERM PLANNING**: When you receive a goal, think about the sequence of areas you need to clear.
**2. SPATIAL AWARENESS**: Check your surroundings (61x61 map) to find the nearest cluster of targets.
**3. EXECUTION**:
   - Move to a tile NEXT to the target.
   - FACE the target.
   - **CONFIRM** the target is in the "Tile in front" data.
   - Use the **Lowest-Energy** tool required.

## ASCII MAP LEGEND (61x61 Vision)
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

## SPATIAL COORDINATION & PRECISION

- **Coordinates**: X is horizontal (0=left), Y is vertical (0=top).
- **Tool Range**: You can ONLY hit the tile directly in front of you.
- **Distance Rule**: You must be exactly 1 tile away from your target.
  - To hit Target (10, 10): Stand at (9, 10) and face "right", OR at (11, 10) and face "left", etc.
  - **DO NOT** stand on the same tile as the target (10, 10).
  - **DO NOT** use move_to to go TO a target coordinate typed 'O' or 'T'. Use move_to to go to a '.' tile NEXT to it.

## TOOL EFFICIENCY

- **SCYTHE**: Use for weeds/grass. It costs **0 ENERGY**. Highly efficient for cleanup.
- **AXE**: Use for wood/twigs/stumps. Costs energy.
- **PICKAXE**: Use for stones/ore. Costs energy.
- **VERIFICATION**: After using a tool, check if the objectName/terrainType in "Tile in front" has changed to "." (walkable ground). If not, your action failed—do not keep moving, FIX it.

## INTELLIGENCE & AUTO-CORRECTION

- **No Path Found?**: The tile you clicked is blocked. Try moving to a tile 1-step away from it.
- **IsMoving Error?**: Movement is now BLOCKING. If a move tool finishes, you are at your destination. Do not issue 10 move commands in a row; wait for each.
- **Cleaning Goals**: Don't just swing randomly. Find a target, move to it, clear it, move to the next.

## SURVIVAL & NIGHT

- **2:00 AM** is a hard game-over. You MUST be in bed by **1:00 AM**.
- Farmhouse Entrance is usually around (60, 15) on the standard farm layout, but check surroundings for "FarmHouse" warp.

## CHEAT MODE (Optional Power Tools)

Cheat mode provides instant, god-mode capabilities. **You must call cheat_mode_enable first** before any cheat commands work.

### Enabling/Disabling
- **cheat_mode_enable**: Activates cheat mode. Required before using any cheats.
- **cheat_mode_disable**: Deactivates cheat mode and all persistent effects (infinite energy, time freeze).

### Instant Resource Cheats
- **cheat_set_money**: Set gold to any amount (e.g., 1000000 for 1 million gold)
- **cheat_add_item**: Add any item by ID (e.g., "(O)465" for starfruit seeds, "(O)74" for prismatic shard)
- **cheat_spawn_ores**: Add ores directly: copper, iron, gold, iridium, coal

### Teleportation
- **cheat_warp**: Teleport to any location (Farm, Town, Mountain, Beach, Forest, Mine, Desert, etc.)
- **cheat_mine_warp**: Warp to specific mine level (1-120 = regular mines, 121+ = Skull Cavern)

### Farming Automation
- **cheat_water_all**: Instantly water all tilled soil
- **cheat_grow_crops**: Instantly grow all crops to harvest-ready
- **cheat_harvest_all**: Instantly harvest all ready crops
- **cheat_clear_debris**: Remove all weeds, stones, twigs, grass
- **cheat_collect_all_forage**: Collect all forage items in current location

### Land Clearing & Farming Automation
- **cheat_hoe_all**: Instantly hoe/till ALL diggable tiles in current location (optional radius parameter)
- **cheat_cut_trees**: Instantly chop ALL trees in current location, collect wood/hardwood/sap/seeds
- **cheat_mine_rocks**: Instantly mine ALL rocks/stones/boulders, collect stone/ores/coal/geodes
- **cheat_dig_artifacts**: Instantly dig ALL artifact spots, collect artifacts/clay/geodes
- **cheat_plant_seeds**: Instantly plant seeds on ALL empty hoed tiles (requires seedId parameter)
- **cheat_fertilize_all**: Apply fertilizer to ALL hoed tiles (optional fertilizerId parameter)

### Mining Automation
- **cheat_instant_mine**: Mine ALL ore nodes in current mine level, drops go to inventory

### Social/Relationship Cheats
- **cheat_set_friendship**: Set friendship with specific NPC (use hearts: 1-10, or points: 0-2500+)
- **cheat_max_all_friendships**: Max out ALL NPC friendships at once
- **cheat_give_gift**: Give gift to NPC instantly (calculates friendship based on NPC's preferences)

### Time Control
- **cheat_time_set**: Set game time (600=6AM, 1200=noon, 1800=6PM, 2400=midnight)
- **cheat_time_freeze**: Toggle time freeze ON/OFF - time stops passing
- **cheat_infinite_energy**: Toggle infinite stamina ON/OFF - never run out of energy

### Other Cheats
- **cheat_set_energy**: Restore stamina to max
- **cheat_set_health**: Restore health to max
- **cheat_unlock_recipes**: Unlock ALL crafting and cooking recipes
- **cheat_pet_all_animals**: Pet all farm animals instantly (daily affection)
- **cheat_complete_quest**: Complete active quests (all or by specific ID)

### Inventory & Upgrade Cheats
- **cheat_upgrade_backpack**: Upgrade backpack to larger size (12, 24, or 36 slots). Default: 36 (max)
- **cheat_upgrade_tool**: Upgrade a specific tool (Hoe, Pickaxe, Axe, WateringCan, FishingRod, Trash Can). Levels: 0=Basic, 1=Copper, 2=Steel, 3=Gold, 4=Iridium
- **cheat_upgrade_all_tools**: Upgrade ALL tools to specified level (default: Iridium)
- **cheat_unlock_all**: ULTIMATE CHEAT - Max backpack, all tools to iridium, all recipes, all skills to level 10, all special items (Rusty Key, Skull Key, Club Card, etc.)

### Targeted/Selective Cheats (For Precise Control & Creative Farming)
These cheats let you control EXACTLY which tiles to affect - perfect for drawing shapes and patterns!

**IMPORTANT: When drawing patterns/shapes, do NOT call cheat_hoe_all first!** The pattern commands automatically clear the surrounding area and hoe only the pattern tiles.

- **cheat_hoe_tiles**: Hoe SPECIFIC tiles by coordinates. Use for precise control.
  - Parameters: tiles (format: "x,y;x,y;x,y") OR x and y for single tile
  - Example: To hoe tiles at (10,20), (11,20), (12,20): tiles="10,20;11,20;12,20"
  
- **cheat_clear_tiles**: Clear SPECIFIC tiles (remove objects, terrain features, hoed dirt).
  - Parameters: tiles or x,y for coordinates
  - Optional: clearObjects (default true), clearFeatures (default true), clearDirt (default true)

- **cheat_hoe_custom_pattern**: Draw ANY shape by designing it yourself as an ASCII grid!
  - **YOU must think about which tiles to hoe** - design the shape in your mind first
  - Use grid parameter: '#' or 'X' = hoe this tile, '.' or ' ' = empty
  - Use '\n' to separate rows
  - Pattern is centered at player position (or x,y if specified)
  - Surrounding area is auto-cleared so your pattern stands out

### Designing Patterns - THINK SPATIALLY!

When asked to draw a shape, YOU must design the ASCII grid. Think about:
1. What does this shape look like from above (top-down view)?
2. Which tiles need to be filled vs empty?
3. Design it row by row

**Example - Heart (9 wide x 8 tall):**
Row 0: .##...##.   <- two bumps at top
Row 1: ####.####   <- bumps widen
Row 2: #########   <- full width  
Row 3: #########   <- full width
Row 4: .#######.   <- starts narrowing
Row 5: ..#####..   <- narrower
Row 6: ...###...   <- narrower
Row 7: ....#....   <- point at bottom
Grid: .##...##.\n####.####\n#########\n#########\n.#######.\n..#####..\n...###...\n....#....

**Example - Star (5 wide x 5 tall):**
Row 0: ..#..   <- top point
Row 1: ..#..   <- stem
Row 2: #####   <- horizontal bar
Row 3: .#.#.   <- lower diagonals
Row 4: #...#   <- bottom points
Grid: ..#..\n..#..\n#####\n.#.#.\n#...#

**Example - Smiley (7 wide x 7 tall):**
Row 0: .#####.   <- top of head
Row 1: #.....#   <- sides
Row 2: #.#.#.#   <- eyes
Row 3: #.....#   <- middle
Row 4: #.###.#   <- mouth
Row 5: #.....#   <- sides  
Row 6: .#####.   <- bottom
Grid: .#####.\n#.....#\n#.#.#.#\n#.....#\n#.###.#\n#.....#\n.#####.

### Drawing Patterns Workflow
1. cheat_mode_enable
2. cheat_warp Farm (or desired location)
3. cheat_clear_debris (remove obstacles)  
4. **DO NOT call cheat_hoe_all** - go directly to pattern!
5. Design your ASCII grid for the shape
6. cheat_hoe_custom_pattern grid="your_grid_here"

### When to Use Cheats
- For rapid testing or speedrunning specific goals
- When exploring content without survival constraints
- When the user explicitly requests "cheat", "instant", or "god mode"
- When time-sensitive goals would otherwise be impossible

### Cheat Strategy
1. Enable cheat mode first: cheat_mode_enable
2. Consider enabling time freeze and infinite energy for stress-free gameplay
3. Use warp for instant travel instead of walking
4. Use instant_mine for quick resource gathering in mines
5. Use max_all_friendships if relationship goals are needed quickly

### Full Farm Setup Workflow (Example)
1. cheat_mode_enable
2. cheat_time_freeze (stop time)
3. cheat_infinite_energy
4. cheat_warp Farm
5. cheat_clear_debris (remove weeds/stones/twigs)
6. cheat_cut_trees (clear all trees)
7. cheat_mine_rocks (clear boulders)
8. cheat_hoe_all (till the ground)
9. cheat_fertilize_all (apply fertilizer)
10. cheat_plant_seeds with seedId "472" (plant parsnips) or "479" (melons)
11. cheat_grow_crops (instant growth)
12. cheat_harvest_all (collect everything)

## Seed IDs by Season

### Spring Seeds
- 472: Parsnip Seeds
- 474: Cauliflower Seeds  
- 476: Potato Seeds
- 427: Tulip Bulb
- 429: Jazz Seeds
- 477: Kale Seeds
- 745: Strawberry Seeds (Festival only)

### Summer Seeds
- 479: Melon Seeds
- 480: Tomato Seeds
- 482: Pepper Seeds
- 483: Wheat Seeds
- 484: Radish Seeds
- 485: Red Cabbage Seeds
- 486: Starfruit Seeds
- 481: Blueberry Seeds
- 302: Hops Starter
- 453: Poppy Seeds
- 455: Spangle Seeds
- 431: Sunflower Seeds

### Fall Seeds
- 487: Corn Seeds (also summer)
- 488: Eggplant Seeds
- 490: Pumpkin Seeds
- 299: Amaranth Seeds
- 301: Grape Starter
- 489: Artichoke Seeds
- 491: Bok Choy Seeds
- 492: Yam Seeds
- 493: Cranberry Seeds
- 494: Beet Seeds
- 425: Fairy Seeds

### Multi-Season
- 433: Coffee Bean (Spring + Summer)
- 745: Ancient Seeds (Spring, Summer, Fall)
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

	// ========== CHEAT MODE TOOLS ==========
	// These tools require cheat_mode_enable to be called first

	cheatEnableTool := copilot.DefineTool("cheat_mode_enable", "Enable cheat mode. Required before using other cheat commands.",
		func(params NoParams, inv copilot.ToolInvocation) (string, error) {
			log.Printf("[TOOL CALL] cheat_mode_enable")
			resp, _ := gameClient.SendCommand("cheat_mode_enable", nil)
			log.Printf("[TOOL RESULT] %s", resp.Message)
			return resp.Message, nil
		})

	cheatDisableTool := copilot.DefineTool("cheat_mode_disable", "Disable cheat mode",
		func(params NoParams, inv copilot.ToolInvocation) (string, error) {
			log.Printf("[TOOL CALL] cheat_mode_disable")
			resp, _ := gameClient.SendCommand("cheat_mode_disable", nil)
			log.Printf("[TOOL RESULT] %s", resp.Message)
			return resp.Message, nil
		})

	cheatWarpTool := copilot.DefineTool("cheat_warp", "Instantly teleport to any location (Farm, Town, Mountain, Beach, Forest, Mine, etc.)",
		func(params CheatWarpParams, inv copilot.ToolInvocation) (string, error) {
			log.Printf("[TOOL CALL] cheat_warp -> %s", params.Location)
			p := map[string]interface{}{"location": params.Location}
			if params.X != 0 {
				p["x"] = params.X
			}
			if params.Y != 0 {
				p["y"] = params.Y
			}
			resp, _ := gameClient.SendCommand("cheat_warp", p)
			log.Printf("[TOOL RESULT] %s", resp.Message)
			return resp.Message, nil
		})

	cheatSetMoneyTool := copilot.DefineTool("cheat_set_money", "Set player's gold amount",
		func(params CheatSetMoneyParams, inv copilot.ToolInvocation) (string, error) {
			resp, _ := gameClient.SendCommand("cheat_set_money", map[string]interface{}{"amount": params.Amount})
			return resp.Message, nil
		})

	cheatAddItemTool := copilot.DefineTool("cheat_add_item", "Add any item to inventory by ID (e.g., '(O)465' for seeds)",
		func(params CheatAddItemParams, inv copilot.ToolInvocation) (string, error) {
			p := map[string]interface{}{"itemId": params.ItemID}
			if params.Count > 0 {
				p["count"] = params.Count
			}
			if params.Quality > 0 {
				p["quality"] = params.Quality
			}
			resp, _ := gameClient.SendCommand("cheat_add_item", p)
			return resp.Message, nil
		})

	cheatSetEnergyTool := copilot.DefineTool("cheat_set_energy", "Restore stamina to max (or specific amount)",
		func(params NoParams, inv copilot.ToolInvocation) (string, error) {
			resp, _ := gameClient.SendCommand("cheat_set_energy", nil)
			return resp.Message, nil
		})

	cheatSetHealthTool := copilot.DefineTool("cheat_set_health", "Restore health to max (or specific amount)",
		func(params NoParams, inv copilot.ToolInvocation) (string, error) {
			resp, _ := gameClient.SendCommand("cheat_set_health", nil)
			return resp.Message, nil
		})

	cheatSetFriendshipTool := copilot.DefineTool("cheat_set_friendship", "Instantly set friendship with any NPC (hearts or points)",
		func(params CheatSetFriendshipParams, inv copilot.ToolInvocation) (string, error) {
			p := map[string]interface{}{"npcName": params.NPCName}
			if params.Hearts > 0 {
				p["hearts"] = params.Hearts
			} else if params.Points > 0 {
				p["points"] = params.Points
			} else {
				p["hearts"] = 10 // default to max
			}
			resp, _ := gameClient.SendCommand("cheat_set_friendship", p)
			return resp.Message, nil
		})

	cheatMaxFriendshipsTool := copilot.DefineTool("cheat_max_all_friendships", "Max out friendship with ALL NPCs at once",
		func(params NoParams, inv copilot.ToolInvocation) (string, error) {
			resp, _ := gameClient.SendCommand("cheat_max_all_friendships", nil)
			return resp.Message, nil
		})

	cheatHarvestAllTool := copilot.DefineTool("cheat_harvest_all", "Instantly harvest all ready crops in current location",
		func(params NoParams, inv copilot.ToolInvocation) (string, error) {
			log.Printf("[TOOL CALL] cheat_harvest_all")
			resp, _ := gameClient.SendCommand("cheat_harvest_all", nil)
			log.Printf("[TOOL RESULT] %s", resp.Message)
			return resp.Message, nil
		})

	cheatWaterAllTool := copilot.DefineTool("cheat_water_all", "Instantly water all soil in current location",
		func(params NoParams, inv copilot.ToolInvocation) (string, error) {
			log.Printf("[TOOL CALL] cheat_water_all")
			resp, _ := gameClient.SendCommand("cheat_water_all", nil)
			log.Printf("[TOOL RESULT] %s", resp.Message)
			return resp.Message, nil
		})

	cheatGrowCropsTool := copilot.DefineTool("cheat_grow_crops", "Instantly grow all crops to harvest-ready",
		func(params NoParams, inv copilot.ToolInvocation) (string, error) {
			log.Printf("[TOOL CALL] cheat_grow_crops")
			resp, _ := gameClient.SendCommand("cheat_grow_crops", nil)
			log.Printf("[TOOL RESULT] %s", resp.Message)
			return resp.Message, nil
		})

	cheatClearDebrisTool := copilot.DefineTool("cheat_clear_debris", "Remove all weeds, stones, twigs, grass in current location",
		func(params NoParams, inv copilot.ToolInvocation) (string, error) {
			log.Printf("[TOOL CALL] cheat_clear_debris")
			resp, _ := gameClient.SendCommand("cheat_clear_debris", nil)
			log.Printf("[TOOL RESULT] %s", resp.Message)
			return resp.Message, nil
		})

	cheatMineWarpTool := copilot.DefineTool("cheat_mine_warp", "Warp directly to specific mine level (1-120 Mines, 121+ Skull Cavern)",
		func(params CheatMineWarpParams, inv copilot.ToolInvocation) (string, error) {
			log.Printf("[TOOL CALL] cheat_mine_warp -> level %d", params.Level)
			resp, _ := gameClient.SendCommand("cheat_mine_warp", map[string]interface{}{"level": params.Level})
			log.Printf("[TOOL RESULT] %s", resp.Message)
			return resp.Message, nil
		})

	cheatSpawnOresTool := copilot.DefineTool("cheat_spawn_ores", "Add ores directly to inventory (copper, iron, gold, iridium, coal)",
		func(params CheatSpawnOresParams, inv copilot.ToolInvocation) (string, error) {
			p := map[string]interface{}{"oreType": params.OreType}
			if params.Count > 0 {
				p["count"] = params.Count
			}
			resp, _ := gameClient.SendCommand("cheat_spawn_ores", p)
			return resp.Message, nil
		})

	cheatCollectForageTool := copilot.DefineTool("cheat_collect_all_forage", "Instantly collect all forage items in current location",
		func(params NoParams, inv copilot.ToolInvocation) (string, error) {
			resp, _ := gameClient.SendCommand("cheat_collect_all_forage", nil)
			return resp.Message, nil
		})

	cheatInstantMineTool := copilot.DefineTool("cheat_instant_mine", "Mine ALL ore nodes in current mine level instantly",
		func(params NoParams, inv copilot.ToolInvocation) (string, error) {
			resp, _ := gameClient.SendCommand("cheat_instant_mine", nil)
			return resp.Message, nil
		})

	cheatTimeSetTool := copilot.DefineTool("cheat_time_set", "Set the game time (600=6AM, 1200=noon, 1800=6PM, 2400=midnight)",
		func(params CheatTimeSetParams, inv copilot.ToolInvocation) (string, error) {
			resp, _ := gameClient.SendCommand("cheat_time_set", map[string]interface{}{"time": params.Time})
			return resp.Message, nil
		})

	cheatTimeFreezeTool := copilot.DefineTool("cheat_time_freeze", "Toggle time freeze on/off",
		func(params NoParams, inv copilot.ToolInvocation) (string, error) {
			resp, _ := gameClient.SendCommand("cheat_time_freeze", nil)
			return resp.Message, nil
		})

	cheatInfiniteEnergyTool := copilot.DefineTool("cheat_infinite_energy", "Toggle infinite stamina on/off",
		func(params NoParams, inv copilot.ToolInvocation) (string, error) {
			log.Printf("[TOOL CALL] cheat_infinite_energy")
			resp, _ := gameClient.SendCommand("cheat_infinite_energy", nil)
			log.Printf("[TOOL RESULT] %s", resp.Message)
			return resp.Message, nil
		})

	cheatUnlockRecipesTool := copilot.DefineTool("cheat_unlock_recipes", "Unlock ALL crafting and cooking recipes",
		func(params NoParams, inv copilot.ToolInvocation) (string, error) {
			log.Printf("[TOOL CALL] cheat_unlock_recipes")
			resp, _ := gameClient.SendCommand("cheat_unlock_recipes", nil)
			log.Printf("[TOOL RESULT] %s", resp.Message)
			return resp.Message, nil
		})

	cheatPetAnimalsTool := copilot.DefineTool("cheat_pet_all_animals", "Pet ALL farm animals instantly",
		func(params NoParams, inv copilot.ToolInvocation) (string, error) {
			log.Printf("[TOOL CALL] cheat_pet_all_animals")
			resp, _ := gameClient.SendCommand("cheat_pet_all_animals", nil)
			log.Printf("[TOOL RESULT] %s", resp.Message)
			return resp.Message, nil
		})

	cheatCompleteQuestTool := copilot.DefineTool("cheat_complete_quest", "Complete active quests instantly",
		func(params CheatCompleteQuestParams, inv copilot.ToolInvocation) (string, error) {
			log.Printf("[TOOL CALL] cheat_complete_quest")
			p := map[string]interface{}{}
			if params.QuestID != "" {
				p["questId"] = params.QuestID
			}
			resp, _ := gameClient.SendCommand("cheat_complete_quest", p)
			log.Printf("[TOOL RESULT] %s", resp.Message)
			return resp.Message, nil
		})

	cheatGiveGiftTool := copilot.DefineTool("cheat_give_gift", "Give a gift to an NPC instantly (for friendship)",
		func(params CheatGiveGiftParams, inv copilot.ToolInvocation) (string, error) {
			log.Printf("[TOOL CALL] cheat_give_gift -> %s", params.NPCName)
			resp, _ := gameClient.SendCommand("cheat_give_gift", map[string]interface{}{
				"npcName": params.NPCName,
				"itemId":  params.ItemID,
			})
			log.Printf("[TOOL RESULT] %s", resp.Message)
			return resp.Message, nil
		})

	// ========== NEW FARMING CHEAT TOOLS ==========

	cheatHoeAllTool := copilot.DefineTool("cheat_hoe_all", "Instantly hoe/till all diggable tiles in current location",
		func(params CheatHoeAllParams, inv copilot.ToolInvocation) (string, error) {
			log.Printf("[TOOL CALL] cheat_hoe_all")
			p := map[string]interface{}{}
			if params.Radius > 0 {
				p["radius"] = params.Radius
			}
			resp, _ := gameClient.SendCommand("cheat_hoe_all", p)
			log.Printf("[TOOL RESULT] %s", resp.Message)
			return resp.Message, nil
		})

	cheatCutTreesTool := copilot.DefineTool("cheat_cut_trees", "Instantly cut/chop ALL trees in current location, collect wood/hardwood",
		func(params CheatCutTreesParams, inv copilot.ToolInvocation) (string, error) {
			log.Printf("[TOOL CALL] cheat_cut_trees")
			p := map[string]interface{}{}
			if !params.IncludeStumps {
				p["includeStumps"] = "false"
			}
			resp, _ := gameClient.SendCommand("cheat_cut_trees", p)
			log.Printf("[TOOL RESULT] %s", resp.Message)
			return resp.Message, nil
		})

	cheatMineRocksTool := copilot.DefineTool("cheat_mine_rocks", "Instantly mine ALL rocks/stones/boulders in current location, collect ores",
		func(params NoParams, inv copilot.ToolInvocation) (string, error) {
			log.Printf("[TOOL CALL] cheat_mine_rocks")
			resp, _ := gameClient.SendCommand("cheat_mine_rocks", nil)
			log.Printf("[TOOL RESULT] %s", resp.Message)
			return resp.Message, nil
		})

	cheatDigArtifactsTool := copilot.DefineTool("cheat_dig_artifacts", "Instantly dig up ALL artifact spots in current location",
		func(params NoParams, inv copilot.ToolInvocation) (string, error) {
			log.Printf("[TOOL CALL] cheat_dig_artifacts")
			resp, _ := gameClient.SendCommand("cheat_dig_artifacts", nil)
			log.Printf("[TOOL RESULT] %s", resp.Message)
			return resp.Message, nil
		})

	cheatPlantSeedsTool := copilot.DefineTool("cheat_plant_seeds", "Instantly plant seeds on ALL empty hoed tiles",
		func(params CheatPlantSeedsParams, inv copilot.ToolInvocation) (string, error) {
			log.Printf("[TOOL CALL] cheat_plant_seeds -> %s", params.SeedID)
			resp, _ := gameClient.SendCommand("cheat_plant_seeds", map[string]interface{}{
				"seedId": params.SeedID,
			})
			log.Printf("[TOOL RESULT] %s", resp.Message)
			return resp.Message, nil
		})

	cheatFertilizeAllTool := copilot.DefineTool("cheat_fertilize_all", "Apply fertilizer to ALL hoed tiles",
		func(params CheatFertilizeAllParams, inv copilot.ToolInvocation) (string, error) {
			log.Printf("[TOOL CALL] cheat_fertilize_all")
			p := map[string]interface{}{}
			if params.FertilizerID != "" {
				p["fertilizerId"] = params.FertilizerID
			}
			resp, _ := gameClient.SendCommand("cheat_fertilize_all", p)
			log.Printf("[TOOL RESULT] %s", resp.Message)
			return resp.Message, nil
		})

	// Inventory & upgrade cheat tools
	cheatUpgradeBackpackTool := copilot.DefineTool("cheat_upgrade_backpack", "Upgrade backpack to larger size (12, 24, or 36 slots). Default: 36 (max)",
		func(params CheatUpgradeBackpackParams, inv copilot.ToolInvocation) (string, error) {
			log.Printf("[TOOL CALL] cheat_upgrade_backpack -> size=%d", params.Size)
			p := map[string]interface{}{}
			if params.Size > 0 {
				p["size"] = params.Size
			}
			resp, _ := gameClient.SendCommand("cheat_upgrade_backpack", p)
			log.Printf("[TOOL RESULT] %s", resp.Message)
			return resp.Message, nil
		})

	cheatUpgradeToolTool := copilot.DefineTool("cheat_upgrade_tool", "Upgrade a specific tool to higher level. Levels: 0=Basic, 1=Copper, 2=Steel, 3=Gold, 4=Iridium",
		func(params CheatUpgradeToolParams, inv copilot.ToolInvocation) (string, error) {
			log.Printf("[TOOL CALL] cheat_upgrade_tool -> %s level=%d", params.Tool, params.Level)
			p := map[string]interface{}{
				"tool": params.Tool,
			}
			if params.Level >= 0 {
				p["level"] = params.Level
			}
			resp, _ := gameClient.SendCommand("cheat_upgrade_tool", p)
			log.Printf("[TOOL RESULT] %s", resp.Message)
			return resp.Message, nil
		})

	cheatUpgradeAllToolsTool := copilot.DefineTool("cheat_upgrade_all_tools", "Upgrade ALL tools to specified level. Default: 4 (Iridium)",
		func(params CheatUpgradeAllToolsParams, inv copilot.ToolInvocation) (string, error) {
			log.Printf("[TOOL CALL] cheat_upgrade_all_tools -> level=%d", params.Level)
			p := map[string]interface{}{}
			if params.Level >= 0 {
				p["level"] = params.Level
			}
			resp, _ := gameClient.SendCommand("cheat_upgrade_all_tools", p)
			log.Printf("[TOOL RESULT] %s", resp.Message)
			return resp.Message, nil
		})

	cheatUnlockAllTool := copilot.DefineTool("cheat_unlock_all", "UNLOCK EVERYTHING: Max backpack, all tools to iridium, all recipes, all skills to level 10, all special items",
		func(params NoParams, inv copilot.ToolInvocation) (string, error) {
			log.Printf("[TOOL CALL] cheat_unlock_all")
			resp, _ := gameClient.SendCommand("cheat_unlock_all", map[string]interface{}{})
			log.Printf("[TOOL RESULT] %s", resp.Message)
			return resp.Message, nil
		})

	// ========== TARGETED/SELECTIVE CHEAT TOOLS (for precise control like drawing shapes) ==========

	cheatHoeTilesTool := copilot.DefineTool("cheat_hoe_tiles", "Hoe SPECIFIC tiles by coordinates. Perfect for drawing shapes/patterns. Use tiles='x,y;x,y' format or single x,y params.",
		func(params CheatHoeTilesParams, inv copilot.ToolInvocation) (string, error) {
			log.Printf("[TOOL CALL] cheat_hoe_tiles tiles=%s x=%d y=%d", params.Tiles, params.X, params.Y)
			p := map[string]interface{}{}
			if params.Tiles != "" {
				p["tiles"] = params.Tiles
			}
			if params.X != 0 || params.Y != 0 {
				p["x"] = params.X
				p["y"] = params.Y
			}
			resp, _ := gameClient.SendCommand("cheat_hoe_tiles", p)
			log.Printf("[TOOL RESULT] %s", resp.Message)
			return resp.Message, nil
		})

	cheatClearTilesTool := copilot.DefineTool("cheat_clear_tiles", "Clear SPECIFIC tiles (objects, terrain, hoed dirt). Use tiles='x,y;x,y' format or single x,y params.",
		func(params CheatClearTilesParams, inv copilot.ToolInvocation) (string, error) {
			log.Printf("[TOOL CALL] cheat_clear_tiles tiles=%s x=%d y=%d", params.Tiles, params.X, params.Y)
			p := map[string]interface{}{}
			if params.Tiles != "" {
				p["tiles"] = params.Tiles
			}
			if params.X != 0 || params.Y != 0 {
				p["x"] = params.X
				p["y"] = params.Y
			}
			// Only include these if explicitly set to false
			if !params.ClearObjects {
				p["clearObjects"] = "false"
			}
			if !params.ClearFeatures {
				p["clearFeatures"] = "false"
			}
			if !params.ClearDirt {
				p["clearDirt"] = "false"
			}
			resp, _ := gameClient.SendCommand("cheat_clear_tiles", p)
			log.Printf("[TOOL RESULT] %s", resp.Message)
			return resp.Message, nil
		})

	// cheatTillPatternTool removed - AI should design its own patterns using cheatHoeCustomPatternTool
	// The preset patterns were too rigid; letting the AI think about tiles produces better results
	_ = copilot.DefineTool("cheat_till_pattern_UNUSED", "UNUSED",
		func(params CheatTillPatternParams, inv copilot.ToolInvocation) (string, error) {
			return "This tool is disabled", nil
		})

	cheatHoeCustomPatternTool := copilot.DefineTool("cheat_hoe_custom_pattern",
		`Draw ANY shape/pattern by hoeing specific tiles. YOU must design the pattern!

HOW TO USE:
1. Think about what shape you want (heart, star, letter, etc.)
2. Design it as an ASCII grid where '#' = hoe this tile, '.' = skip
3. Pass the grid string with \n for newlines

EXAMPLE - Heart shape:
grid=".##.##.\n#######\n#######\n.#####.\n..###..\n...#..."

EXAMPLE - Letter A:
grid="..#..\n.#.#.\n#####\n#...#\n#...#"

EXAMPLE - Star:
grid="..#..\n..#..\n#####\n.#.#.\n#...#"

The pattern will be centered at your position (or x,y if specified).
Surrounding area is auto-cleared so pattern is visible.`,
		func(params CheatHoeCustomPatternParams, inv copilot.ToolInvocation) (string, error) {
			log.Printf("[TOOL CALL] cheat_hoe_custom_pattern x=%d y=%d grid=%q offsets=%q clearArea=%v",
				params.X, params.Y, params.Grid, params.OffsetString, params.ClearArea)
			p := map[string]interface{}{}
			if params.X != 0 {
				p["x"] = params.X
			}
			if params.Y != 0 {
				p["y"] = params.Y
			}
			if params.Grid != "" {
				p["grid"] = params.Grid
			}
			if params.OffsetString != "" {
				p["offsetString"] = params.OffsetString
			}
			if params.ClearRadius > 0 {
				p["clearRadius"] = params.ClearRadius
			}
			// clearArea defaults to true, only send if explicitly false
			if !params.ClearArea {
				p["clearArea"] = "false"
			}
			resp, _ := gameClient.SendCommand("cheat_hoe_custom_pattern", p)
			log.Printf("[TOOL RESULT] %s", resp.Message)
			return resp.Message, nil
		})

	// Create session with tools (using embedded knowledge)
	session, err := a.client.CreateSession(&copilot.SessionConfig{
		Model: "gpt-4.1",
		SystemMessage: &copilot.SystemMessageConfig{
			Content: gameKnowledge,
		},
		Tools: []copilot.Tool{
			// Standard gameplay tools
			moveToTool, getSurroundingsTool, interactTool, useToolTool,
			useToolRepeatTool, faceDirectionTool, selectItemTool, switchToolTool,
			eatItemTool, enterDoorTool, findBestTargetTool, clearTargetTool,
			// Cheat mode tools
			cheatEnableTool, cheatDisableTool, cheatWarpTool, cheatSetMoneyTool,
			cheatAddItemTool, cheatSetEnergyTool, cheatSetHealthTool,
			cheatSetFriendshipTool, cheatMaxFriendshipsTool,
			cheatHarvestAllTool, cheatWaterAllTool, cheatGrowCropsTool, cheatClearDebrisTool,
			cheatMineWarpTool, cheatSpawnOresTool, cheatCollectForageTool, cheatInstantMineTool,
			cheatTimeSetTool, cheatTimeFreezeTool, cheatInfiniteEnergyTool,
			cheatUnlockRecipesTool, cheatPetAnimalsTool, cheatCompleteQuestTool, cheatGiveGiftTool,
			// New farming cheat tools
			cheatHoeAllTool, cheatCutTreesTool, cheatMineRocksTool, cheatDigArtifactsTool,
			cheatPlantSeedsTool, cheatFertilizeAllTool,
			// Inventory & upgrade cheat tools
			cheatUpgradeBackpackTool, cheatUpgradeToolTool, cheatUpgradeAllToolsTool, cheatUnlockAllTool,
			// Targeted/selective cheat tools (for precise control like drawing shapes)
			cheatHoeTilesTool, cheatClearTilesTool, cheatHoeCustomPatternTool,
			// Note: cheatTillPatternTool removed - AI should design its own patterns using cheatHoeCustomPatternTool
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
	goalCompleted := false
	iteration := 0

	log.Printf("[AGENT LOOP] Starting autonomous loop...")

	for {
		iteration++

		// Stop if goal was completed
		if goalCompleted {
			log.Printf("[AGENT LOOP] Goal completed! Stopping autonomous loop.")
			time.Sleep(30 * time.Second) // Wait before potentially starting new goal
			goalCompleted = false        // Reset for next iteration
		}

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

		// Get season-appropriate seed suggestions (use numeric IDs only, not (O) prefix)
		seasonSeeds := map[string]string{
			"spring": "472 (Parsnip), 474 (Cauliflower), 476 (Potato)",
			"summer": "479 (Melon), 480 (Tomato), 482 (Pepper)",
			"fall":   "487 (Corn), 488 (Eggplant), 490 (Pumpkin)",
			"winter": "No outdoor crops - use greenhouse only",
		}
		seedSuggestion := seasonSeeds[state.Time.Season]
		if seedSuggestion == "" {
			seedSuggestion = "472 (Parsnip)"
		}

		// Clear prompt - execute ALL steps in ONE call, then STOP
		prompt := fmt.Sprintf(`Location: %s | Pos: (%d,%d) | Season: %s | Time: %s | Energy: %.0f/%d
%s

GOAL: %s

SEASON INFO: Current season is %s. Valid seeds: %s

CRITICAL EXECUTION ORDER - These tools have dependencies and MUST be called SEQUENTIALLY (one at a time, waiting for each to complete):
1. cheat_mode_enable (FIRST - enables all other cheats)
2. cheat_clear_debris, cheat_cut_trees, cheat_mine_rocks (can be parallel - clearing the land)
3. cheat_hoe_all (MUST complete before planting - creates hoed tiles)
4. cheat_plant_seeds (MUST run AFTER hoe_all completes - needs hoed tiles to exist)
5. cheat_grow_crops (MUST run AFTER plant_seeds - needs crops to exist)
6. cheat_harvest_all (MUST run AFTER grow_crops - needs mature crops)

DO NOT call plant_seeds, grow_crops, or harvest_all in parallel - they depend on each other!
After ALL tools complete successfully, respond with "GOAL COMPLETE".`,
			state.Player.Location, int(state.Player.X), int(state.Player.Y),
			state.Time.Season, state.Time.TimeString, state.Player.Energy, state.Player.MaxEnergy,
			urgency,
			activeGoal,
			state.Time.Season, seedSuggestion)

		// Only include game context if not using cheats
		if !strings.Contains(strings.ToLower(activeGoal), "cheat") {
			prompt += "\n\n" + gameContext
		}

		// Send message and wait for response
		log.Printf("[AGENT LOOP] Sending prompt (%d chars) to Copilot...", len(prompt))
		response, err := a.session.SendAndWait(copilot.MessageOptions{
			Prompt: prompt,
		}, 120*time.Second) // 120 second timeout for complex cheat operations
		if err != nil {
			log.Printf("[AGENT AGENT] SendAndWait error: %v", err)
			time.Sleep(5 * time.Second)
			continue
		}
		log.Printf("[AGENT LOOP] Got response from Copilot")

		// Log the response and check for goal completion
		if response != nil && response.Data.Content != nil {
			thought := strings.TrimSpace(*response.Data.Content)
			if thought != "" {
				log.Printf("[AGENT THOUGHT] %s", thought)
			}

			// Check for goal completion signal
			thoughtUpper := strings.ToUpper(thought)
			if strings.Contains(thoughtUpper, "GOAL COMPLETE") ||
				strings.Contains(thoughtUpper, "GOAL COMPLETED") ||
				strings.Contains(thoughtUpper, "ALL TASKS COMPLETE") ||
				strings.Contains(thoughtUpper, "MISSION ACCOMPLISHED") {
				log.Printf("[AGENT LOOP] Goal completion detected!")
				goalCompleted = true
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

type NoParams struct {
	// Dummy field required for Gemini schema validation
	Unused string `json:"_unused,omitempty" jsonschema:"Ignore this parameter"`
}

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

// Cheat mode parameter structs
type CheatWarpParams struct {
	Location string `json:"location" jsonschema:"Location name (Farm, Town, Mountain, Beach, Forest, Mine, etc.)"`
	X        int    `json:"x,omitempty" jsonschema:"Optional X coordinate"`
	Y        int    `json:"y,omitempty" jsonschema:"Optional Y coordinate"`
}

type CheatSetMoneyParams struct {
	Amount int `json:"amount" jsonschema:"Amount of gold to set"`
}

type CheatAddItemParams struct {
	ItemID  string `json:"itemId" jsonschema:"Item ID (e.g., '(O)465' for Parsnip Seeds, '(T)Pickaxe' for tools)"`
	Count   int    `json:"count,omitempty" jsonschema:"Number of items (default 1)"`
	Quality int    `json:"quality,omitempty" jsonschema:"Quality (0=normal, 1=silver, 2=gold, 4=iridium)"`
}

type CheatSetFriendshipParams struct {
	NPCName string `json:"npcName" jsonschema:"NPC name (e.g., Abigail, Sebastian)"`
	Hearts  int    `json:"hearts,omitempty" jsonschema:"Friendship hearts (0-14, default 10)"`
	Points  int    `json:"points,omitempty" jsonschema:"Friendship points (250 per heart)"`
}

type CheatMineWarpParams struct {
	Level int `json:"level" jsonschema:"Mine level (1-120 for Mines, 121+ for Skull Cavern, 77377 for Quarry)"`
}

type CheatSpawnOresParams struct {
	OreType string `json:"oreType" jsonschema:"Type of ore (copper, iron, gold, iridium, coal)"`
	Count   int    `json:"count,omitempty" jsonschema:"Number of ores (default 10)"`
}

type CheatTimeSetParams struct {
	Time int `json:"time" jsonschema:"Time in 24-hour format (600=6AM, 1800=6PM, 2600=2AM)"`
}

type CheatGiveGiftParams struct {
	NPCName string `json:"npcName" jsonschema:"NPC name to give gift to"`
	ItemID  string `json:"itemId" jsonschema:"Item ID to give as gift"`
}

type CheatCompleteQuestParams struct {
	QuestID string `json:"questId,omitempty" jsonschema:"Quest ID or name to complete (omit to complete all)"`
}

type CheatHoeAllParams struct {
	Radius int `json:"radius,omitempty" jsonschema:"Radius around player to hoe (default 50)"`
}

type CheatCutTreesParams struct {
	IncludeStumps bool `json:"includeStumps,omitempty" jsonschema:"Whether to include tree stumps (default true)"`
}

type CheatPlantSeedsParams struct {
	SeedID string `json:"seedId" jsonschema:"Seed ID to plant (e.g., '(O)472' for Parsnip Seeds)"`
}

type CheatFertilizeAllParams struct {
	FertilizerID string `json:"fertilizerId,omitempty" jsonschema:"Fertilizer ID (default Quality Fertilizer)"`
}

type CheatUpgradeBackpackParams struct {
	Size int `json:"size,omitempty" jsonschema:"Backpack size: 12, 24, or 36 (default 36)"`
}

type CheatUpgradeToolParams struct {
	Tool  string `json:"tool" jsonschema:"Tool name: Hoe, Pickaxe, Axe, WateringCan, FishingRod, or Trash Can"`
	Level int    `json:"level,omitempty" jsonschema:"Upgrade level: 0=Basic, 1=Copper, 2=Steel, 3=Gold, 4=Iridium (default 4)"`
}

type CheatUpgradeAllToolsParams struct {
	Level int `json:"level,omitempty" jsonschema:"Upgrade level: 0=Basic, 1=Copper, 2=Steel, 3=Gold, 4=Iridium (default 4)"`
}

// Targeted/selective cheat params (for precise control like drawing shapes)
type CheatHoeTilesParams struct {
	Tiles string `json:"tiles,omitempty" jsonschema:"Tile coordinates as 'x,y;x,y;x,y' format OR use x/y arrays"`
	X     int    `json:"x,omitempty" jsonschema:"Single tile X coordinate (use with y for one tile)"`
	Y     int    `json:"y,omitempty" jsonschema:"Single tile Y coordinate (use with x for one tile)"`
}

type CheatClearTilesParams struct {
	Tiles         string `json:"tiles,omitempty" jsonschema:"Tile coordinates as 'x,y;x,y;x,y' format OR use x/y arrays"`
	X             int    `json:"x,omitempty" jsonschema:"Single tile X coordinate"`
	Y             int    `json:"y,omitempty" jsonschema:"Single tile Y coordinate"`
	ClearObjects  bool   `json:"clearObjects,omitempty" jsonschema:"Clear objects like debris, stones (default true)"`
	ClearFeatures bool   `json:"clearFeatures,omitempty" jsonschema:"Clear terrain features like grass, trees (default true)"`
	ClearDirt     bool   `json:"clearDirt,omitempty" jsonschema:"Clear hoed dirt (default true)"`
}

type CheatTillPatternParams struct {
	Pattern     string `json:"pattern" jsonschema:"Pattern to draw: heart, circle, square, filled_square, line, cross, star, diamond, smiley, spiral, arrow"`
	X           int    `json:"x,omitempty" jsonschema:"Center X coordinate (default: player position)"`
	Y           int    `json:"y,omitempty" jsonschema:"Center Y coordinate (default: player position)"`
	Size        int    `json:"size,omitempty" jsonschema:"Size/scale of pattern (default 5, max 50)"`
	Direction   string `json:"direction,omitempty" jsonschema:"Direction for line/arrow patterns: up, down, left, right, ne, nw, se, sw (default: up)"`
	ClearArea   bool   `json:"clearArea,omitempty" jsonschema:"Clear surrounding hoed dirt to make pattern visible (default true)"`
	ClearRadius int    `json:"clearRadius,omitempty" jsonschema:"Radius around pattern to clear (default: size*2+5)"`
}

type CheatHoeCustomPatternParams struct {
	X            int    `json:"x,omitempty" jsonschema:"Center X coordinate (default: player position)"`
	Y            int    `json:"y,omitempty" jsonschema:"Center Y coordinate (default: player position)"`
	Grid         string `json:"grid,omitempty" jsonschema:"ASCII art grid where # or X marks tiles to hoe. Use \\n for newlines. Example: '..#..\\n.###.\\n#####\\n.###.\\n..#..' for diamond"`
	OffsetString string `json:"offsetString,omitempty" jsonschema:"Relative offsets as 'dx,dy;dx,dy'. Example: '0,0;1,0;-1,0;0,1;0,-1' for a cross"`
	ClearArea    bool   `json:"clearArea,omitempty" jsonschema:"Clear surrounding hoed dirt to make pattern visible (default true)"`
	ClearRadius  int    `json:"clearRadius,omitempty" jsonschema:"Radius around pattern to clear (default: pattern size + 5)"`
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
