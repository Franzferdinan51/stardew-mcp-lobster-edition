/**
 * Stardew Valley SDK - Bot Creator
 *
 * Usage:
 *   npm run create-bot -- <botname>
 *
 * Creates a new Stardew Valley bot with:
 *   bots/<botname>/bot.ts - Main bot file
 *   bots/<botname>/bot.env - Configuration
 */

import * as fs from 'fs';
import * as path from 'path';

const botName = process.argv[2];

if (!botName) {
  console.log('Usage: npm run create-bot -- <botname>');
  console.log('Example: npm run create-bot -- my-farmer');
  process.exit(1);
}

const botsDir = path.join(__dirname, '..', 'bots');
const botDir = path.join(botsDir, botName);

// Create bots directory
if (!fs.existsSync(botsDir)) {
  fs.mkdirSync(botsDir, { recursive: true });
}

// Create bot directory
if (fs.existsSync(botDir)) {
  console.error(`Bot "${botName}" already exists!`);
  process.exit(1);
}

fs.mkdirSync(botDir, { recursive: true });

// Create bot.ts - Stardew Valley specific bot
const botCode = `import { Bot, GameClient, StardewTools } from '../../src';

/**
 * ${botName} - Stardew Valley Bot
 *
 * Customize this bot to automate your farm!
 *
 * Available tools:
 *   bot.getTools().moveTo(x, y)           - Move to coordinates
 *   bot.getTools().interact()              - Interact with object
 *   bot.getTools().useTool()               - Use current tool
 *   bot.getTools().switchTool('hoe')      - Switch to hoe
 *   bot.getTools().faceDirection(2)       - Face: 0=up, 1=right, 2=down, 3=left
 *   bot.getTools().eatItem(slot)          - Eat food from inventory
 *   bot.getTools().enterDoor()            - Enter door/building
 *
 * Cheat mode (use carefully!):
 *   bot.getTools().cheatModeEnable()
 *   bot.getTools().cheatWarp('Farm')
 *   bot.getTools().cheatSetMoney(10000)
 *   bot.getTools().cheatGrowCrops()
 *   bot.getTools().cheatHarvestAll()
 */

// Bot configuration
const config = {
  // Game connection
  gameUrl: process.env.GAME_URL || 'ws://localhost:8765/game',

  // What should the bot do?
  goal: 'Manage the farm efficiently: plant crops, water them, harvest, and sell produce',

  // Auto-start when connected
  autoStart: true,
};

// Initialize bot
const bot = new Bot(config);

// Bot logic - customize this!
async function farmLoop() {
  const state = bot.getState();
  if (!state) return;

  const tools = bot.getTools();
  const location = state.player.location;
  const energy = state.player.energy;
  const money = state.player.money;
  const time = state.time.timeOfDay;

  console.log(\`[\${location}] Energy: \${energy}/\${state.player.maxEnergy} | Money: \${money}g | Time: \${state.time.timeString}\`);

  // Emergency: Low energy - go eat or go home
  if (energy < 20) {
    console.log('[Bot] Low energy! Finding food...');
    // Try to find food in inventory and eat
    const food = state.player.inventory.find(item =>
      item.category === 'Food' || item.name.includes('Cheese') ||
      item.name.includes('Milk') || item.name.includes('Egg')
    );
    if (food) {
      await tools.eatItem(food.slot);
      console.log(\`[Bot] Ate \${food.displayName}\`);
    } else {
      console.log('[Bot] No food, heading home...');
      await tools.cheatWarp?.('Farm');
    }
    return;
  }

  // Getting late - go home
  if (time > 2200) {
    console.log('[Bot] Getting late, heading home...');
    try {
      await tools.cheatWarp?.('Farm');
    } catch {
      // Warp failed, just wait
    }
    return;
  }

  // Location-specific actions
  if (location === 'Farm') {
    // Farm logic - customize this!
    console.log('[Bot] On farm - add your farming logic here');

    // Example: Move to coordinates and use tool
    // await tools.moveTo(10, 5);
    // await tools.useTool();
  }

  // Add more locations as needed...
  // if (location === 'Town') { ... }
  // if (location === 'Mine') { ... }
}

// Main loop - runs every 5 seconds
let tickCount = 0;
async function loop() {
  tickCount++;

  try {
    await farmLoop();
  } catch (error: any) {
    console.error('[Bot] Error in loop:', error.message);
  }
}

// Run the bot
bot.start().then(() => {
  console.log('Bot started! - ' + config.goal);

  // Start the main loop
  setInterval(loop, 5000);

}).catch(err => {
  console.error('Failed to start bot:', err);
  process.exit(1);
});
`;

// Create bot.ts
fs.writeFileSync(path.join(botDir, 'bot.ts'), botCode);

// Create bot.env - Environment configuration
const envContent = `# Stardew Valley Bot Configuration
# =============================================================================

# Game Connection
# The WebSocket URL to the Stardew MCP server
GAME_URL=ws://localhost:8765/game

# =============================================================================
# Bot Settings
# =============================================================================

# What should the bot do? (used for AI decision making)
GOAL=Manage the farm efficiently: plant crops, water them, harvest, and sell produce

# Auto-start when connected to game
AUTO_START=true

# =============================================================================
# Server Settings (for remote connections)
# =============================================================================

# Leave empty for local, or set to remote server URL
# SERVER=ws://192.168.1.100:8765/mcp

# =============================================================================
# AI Settings (optional - for AI-powered bots)
# =============================================================================

# Claude API Key for AI decision making
# CLAUDE_API_KEY=sk-ant-...

# LLM timeout in seconds
# LLM_TIMEOUT=60
`;

// Create bot.env
fs.writeFileSync(path.join(botDir, 'bot.env'), envContent);

// Create README.md for the bot
const readme = `# ${botName}

Stardew Valley bot created with Stardew SDK.

## Quick Start

\`\`\`bash
# Development mode (auto-restart on changes)
npm run dev -- bots/${botName}/bot.ts

# Or build and run
npm run build
node bots/${botName}/bot.js
\`\`\`

## Configuration

Edit \`bot.env\` to configure:
- **GAME_URL** - WebSocket URL to Stardew MCP server (default: ws://localhost:8765/game)
- **GOAL** - What the bot should do
- **AUTO_START** - Whether to auto-start when connected
- **SERVER** - Remote server URL (optional)
- **CLAUDE_API_KEY** - For AI-powered decision making

## Bot Logic

Edit \`bot.ts\` to customize your bot's behavior:

### Available Tools

\`\`\`typescript
// Movement
await bot.getTools().moveTo(x, y);           // Navigate to position
await bot.getTools().faceDirection(0-3);      // 0=up, 1=right, 2=down, 3=left

// Tools
await bot.getTools().useTool();               // Use current tool
await bot.getTools().switchTool('hoe');       // Switch tool (hoe, water, axe, pickaxe, etc.)
await bot.getTools().useToolRepeat(5, 0);     // Use tool 5 times facing up

// Interaction
await bot.getTools().interact();              // Interact with object in front
await bot.getTools().enterDoor();              // Enter building/door

// Inventory
await bot.getTools().eatItem(slot);           // Eat food from inventory slot
await bot.getTools().selectItem(slot);        // Select item from inventory
await bot.getTools().shipItem('parsnip', 10); // Ship items to bin

// Fishing
await bot.getTools().castFishingRod();         // Cast fishing rod
await bot.getTools().reelFish();              // Reel in fish

// Cheat Mode (use carefully!)
await bot.getTools().cheatModeEnable();
await bot.getTools().cheatWarp('Farm');       // Teleport to location
await bot.getTools().cheatSetMoney(10000);    // Set money
await bot.getTools().cheatGrowCrops();         // Instantly grow all crops
await bot.getTools().cheatHarvestAll();        // Harvest all crops
await bot.getTools().cheatWaterAll();         // Water all crops

// Game State
const state = bot.getState();
console.log(state.player.location);            // Current location
console.log(state.player.energy);              // Current energy
console.log(state.player.money);               // Current money
console.log(state.time.timeString);            // Current time
\`\`\`

### Available Locations for Warp

- Farm, FarmHouse
- Town, CommunityCenter
- Beach, Mountain
- Forest, BusStop
- Mine, SkullCavern
- Desert (requires unlock)
- And more...

### Common Farming Tasks

1. **Planting**: Move to hoed tile, select seeds, use tool
2. **Watering**: Switch to watering can, use tool on crops
3. **Harvesting**: Use tool on ready crops
4. **Selling**: Ship items or sell at shops

## Troubleshooting

- **Not connecting**: Make sure Stardew Valley is running with SMAPI and StardewMCP mod
- **Commands failing**: Check that you have enough energy and the right tool equipped
- **Pathfinding issues**: Some areas may be blocked - try warping instead

## Examples

### Simple Farm Auto-Manager

\`\`\`typescript
async function farmLoop() {
  const state = bot.getState();
  const tools = bot.getTools();

  if (state.player.location === 'Farm') {
    // Check energy
    if (state.player.energy < 10) {
      await tools.cheatSetEnergy(200);
    }

    // Your farming logic here!
  }
}
\`\`\`

---

Happy Farming! 🌾
`;

// Create README.md
fs.writeFileSync(path.join(botDir, 'README.md'), readme);

console.log(`
Bot "${botName}" created successfully!

Files created:
  bots/${botName}/bot.ts    - Main bot code
  bots/${botName}/bot.env   - Configuration
  bots/${botName}/README.md - Documentation

Next steps:
  1. cd bots/${botName}
  2. Edit bot.ts to add your farming logic
  3. Run: npm run dev -- bots/${botName}/bot.ts

Have fun farming! 🌾
`);
