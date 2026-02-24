# Stardew MCP Tools Reference

Complete reference for all tools available to AI agents controlling Stardew Valley.

## Core Game Tools

### get_state

Get complete game state including player position, inventory, time, location, and stats.

**Returns:** Full game state object with player info, time, location, inventory, skills, quests, relationships

### get_surroundings

Get detailed information about nearby tiles, objects, NPCs, and terrain.

**Returns:** ASCII map, nearby objects, terrain features, NPCs, monsters, buildings, animals, warp points

---

## Movement & Navigation

### move_to

Navigate player to specified X,Y coordinates using A* pathfinding.

| Parameter | Type | Description |
|-----------|------|-------------|
| x | integer | Target X coordinate |
| y | integer | Target Y coordinate |

**Example:**
```json
{"action": "move_to", "params": {"x": 42, "y": 18}}
```

### face_direction

Change the direction the player is facing.

| Parameter | Type | Description |
|-----------|------|-------------|
| direction | integer | 0=down, 1=left, 2=right, 3=up |

**Example:**
```json
{"action": "face_direction", "params": {"direction": 2}}
```

---

## Interaction Tools

### interact

Interact with the object or NPC directly in front of the player.

**Used for:** Talking to NPCs, picking up items, opening doors, activating machines

### enter_door

Enter a building or mine shaft when standing at a door/warp point.

---

## Tool Usage

### use_tool

Use the currently selected tool (hoe, watering can, pickaxe, axe, etc.).

### use_tool_repeat

Use tool repeatedly with configurable delay.

| Parameter | Type | Description |
|-----------|------|-------------|
| count | integer | Number of times to use tool |
| direction | integer | Direction to face (0=down, 1=left, 2=right, 3=up) |

### switch_tool

Switch to a specific tool by name.

| Parameter | Type | Description |
|-----------|------|-------------|
| tool | string | Tool name: hoe, watering_can, pickaxe, axe, sword, fishing_rod, etc. |

**Example:**
```json
{"action": "switch_tool", "params": {"tool": "hoe"}}
```

---

## Inventory Management

### select_item

Select an item from inventory by slot number.

| Parameter | Type | Description |
|-----------|------|-------------|
| slot | integer | Inventory slot (0-11) |

**Example:**
```json
{"action": "select_item", "params": {"slot": 0}}
```

### eat_item

Eat an item from inventory to restore energy/health.

| Parameter | Type | Description |
|-----------|------|-------------|
| slot | integer | Inventory slot of food item |

---

## Targeting

### find_best_target

Find the nearest actionable object (harvestable crop, ore, tree, etc.).

| Parameter | Type | Description |
|-----------|------|-------------|
| type | string | Target type: crop, ore, tree, forage, monster, npc |

**Example:**
```json
{"action": "find_best_target", "params": {"type": "crop"}}
```

### clear_target

Clear the current target selection.

---

## Cheat Mode Tools

**IMPORTANT:** Must call `cheat_mode_enable` first before any cheat commands work!

### cheat_mode_enable

Enable god-mode cheat commands. Must be called before any other cheat.

### cheat_mode_disable

Disable cheat mode and all persistent effects.

### cheat_time_freeze

Toggle time freeze on/off.

| Parameter | Type | Description |
|-----------|------|-------------|
| freeze | boolean | true to freeze, false to unfreeze |

### cheat_infinite_energy

Toggle infinite stamina on/off.

| Parameter | Type | Description |
|-----------|------|-------------|
| enable | boolean | true to enable, false to disable |

---

### Teleportation

### cheat_warp

Teleport to any location.

| Parameter | Type | Description |
|-----------|------|-------------|
| location | string | Location name |

**Valid Locations:**
- Farm, FarmHouse
- Town, BusStop
- Beach, Marina
- Mountain, Mine, Mines
- Forest, BusStop
- Desert, Casino
- CommunityCenter, JojaMart
- WizardTower
- WitchHut
- Sewer
- VolcanoDungeon
- Custom: Mine_1 to Mine_120, SkullCavern

**Example:**
```json
{"action": "cheat_warp", "params": {"location": "Farm"}}
```

### cheat_mine_warp

Warp to a specific mine level.

| Parameter | Type | Description |
|-----------|------|-------------|
| level | integer | Mine level (1-120 = Mines, 121+ = Skull Cavern) |

---

### Farming Automation

### cheat_clear_debris

Remove all weeds, stones, twigs, grass in current area.

### cheat_cut_trees

Chop down all trees and collect wood/hardwood.

### cheat_mine_rocks

Break all rocks/boulders and collect ores.

### cheat_hoe_all

Till all diggable tiles in current area.

### cheat_water_all

Water all crops in current area.

### cheat_plant_seeds

Plant seeds on all empty tilled soil.

| Parameter | Type | Description |
|-----------|------|-------------|
| season | string | Season: spring, summer, fall, winter |
| seedId | string | Optional specific seed ID |

### cheat_fertilize_all

Apply fertilizer to all tilled tiles.

| Parameter | Type | Description |
|-----------|------|-------------|
| type | string | Fertilizer: speed_gro, deluxe_speed_gro, hyper_speed_gro, quality, deluxe_quality |

### cheat_grow_crops

Instantly grow all crops to harvest-ready state.

### cheat_harvest_all

Harvest all ready crops in current area.

### cheat_dig_artifacts

Dig up all artifact spots.

---

### Pattern Drawing

### cheat_hoe_tiles

Create a pattern of tilled tiles.

| Parameter | Type | Description |
|-----------|------|-------------|
| pattern | string | ASCII pattern (x = till, . = don't till) |

### cheat_clear_tiles

Clear objects from tiles in a pattern.

| Parameter | Type | Description |
|-----------|------|-------------|
| pattern | string | ASCII pattern (x = clear, . = don't clear) |

### cheat_hoe_custom_pattern

Draw shapes using ASCII grid input.

| Parameter | Type | Description |
|-----------|------|-------------|
| grid | string | Multi-line ASCII grid |

---

### Resources & Items

### cheat_set_money

Set player's money to a specific amount.

| Parameter | Type | Description |
|-----------|------|-------------|
| amount | integer | Amount of gold |

**Example:**
```json
{"action": "cheat_set_money", "params": {"amount": 10000}}
```

### cheat_add_item

Add an item to inventory by ID.

| Parameter | Type | Description |
|-----------|------|-------------|
| id | string | Item ID (e.g., "771" for pumpkin seeds) |
| count | integer | Number of items |

**Common Item IDs:**
- Seeds: 472 (parsnip), 473 (carrot), 474 (corn), 475 (tomato), 476 (potato), 477 (red cabbage), 478 (radish), 479 (eggplant), 480 (pumpkin), 481 (sunflower), 482 (pepper)
- Tools: 0 (hoe), 1 (watering can), 2 (pickaxe), 3 (axe), 4 (sword), 5 (galoshes), 6 (lamp), 7 (boots), 8 (fishing rod), 9 (trident), 10 (wands)

### cheat_spawn_ores

Spawn ore nodes in the mine.

| Parameter | Type | Description |
|-----------|------|-------------|
| type | string | Ore type: copper, iron, gold, iridium, coal |

### cheat_set_energy

Restore stamina to specific level.

| Parameter | Type | Description |
|-----------|------|-------------|
| amount | integer | Energy amount |

### cheat_set_health

Restore health to specific level.

| Parameter | Type | Description |
|-----------|------|-------------|
| amount | integer | Health amount |

---

### Social

### cheat_set_friendship

Set friendship level with an NPC.

| Parameter | Type | Description |
|-----------|------|-------------|
| npc | string | NPC name |
| points | integer | Friendship points (0-2500, ~250 per heart) |

### cheat_max_all_friendships

Max out friendship with all NPCs (8 hearts).

### cheat_give_gift

Give a gift to an NPC instantly (counts as liked/loved).

| Parameter | Type | Description |
|-----------|------|-------------|
| npc | string | NPC name |
| item_id | string | Item ID to give |

---

### Upgrades

### cheat_upgrade_backpack

Upgrade backpack size.

| Parameter | Type | Description |
|-----------|------|-------------|
| level | integer | Upgrade level (1=12→24 slots, 2=24→36 slots) |

### cheat_upgrade_tool

Upgrade a specific tool.

| Parameter | Type | Description |
|-----------|------|-------------|
| tool | string | Tool name: hoe, watering_can, pickaxe, axe |
| level | integer | Upgrade level (0=basic, 1=copper, 2=steel, 3=gold, 4=iridium) |

### cheat_upgrade_all_tools

Upgrade all tools to specified level.

| Parameter | Type | Description |
|-----------|------|-------------|
| level | integer | Upgrade level (0-4) |

### cheat_unlock_all

Max everything: backpack, tools, recipes, skills, and unlocks casino.

---

## Quick Reference

### Common Tasks

**Setup new farm:**
```
cheat_mode_enable
cheat_warp Farm
cheat_clear_debris
cheat_cut_trees
cheat_mine_rocks
cheat_hoe_all
cheat_plant_seeds spring
cheat_grow_crops
cheat_harvest_all
```

**Get money:**
```
cheat_mode_enable
cheat_set_money 50000
```

**Max skills:**
```
cheat_mode_enable
cheat_upgrade_all_tools 4
```
