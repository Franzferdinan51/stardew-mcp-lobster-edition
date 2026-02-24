/**
 * Stardew Valley SDK - Simple Test Bot
 *
 * Usage:
 *   npm run test
 *
 * This is a simple bot that connects to Stardew Valley
 * and performs basic actions to test the SDK.
 */

import { Bot, GameClient, StardewTools } from './index';

const config = {
  gameUrl: process.env.GAME_URL || 'ws://localhost:8765/game',
  goal: 'Test the Stardew SDK',
};

async function main() {
  console.log('Starting Stardew SDK Test Bot...');
  console.log(`Connecting to: ${config.gameUrl}`);

  const bot = new Bot(config);

  bot.on('connected', () => {
    console.log('[Bot] Connected to Stardew Valley!');
  });

  bot.on('disconnected', () => {
    console.log('[Bot] Disconnected from Stardew Valley');
  });

  bot.on('state', (state) => {
    if (state) {
      console.log(`[Bot] State Update - Location: ${state.player.location}, ` +
        `Money: ${state.player.money}, Energy: ${state.player.energy}/${state.player.maxEnergy}`);
    }
  });

  bot.on('error', (error) => {
    console.error('[Bot] Error:', error.message);
  });

  try {
    await bot.start();
    console.log('[Bot] Started successfully!');

    // Wait for initial state
    await new Promise(resolve => setTimeout(resolve, 2000));

    // Test some basic commands
    const state = bot.getState();
    if (state) {
      console.log(`\n[Bot] Player: ${state.player.name}`);
      console.log(`[Bot] Location: ${state.player.location}`);
      console.log(`[Bot] Money: ${state.player.money}g`);
      console.log(`[Bot] Energy: ${state.player.energy}/${state.player.maxEnergy}`);
      console.log(`[Bot] Time: ${state.time.timeString}`);

      // List inventory
      console.log('\n[Bot] Inventory:');
      state.player.inventory.forEach((item, idx) => {
        if (item.name) {
          console.log(`  ${idx}: ${item.displayName} x${item.stack}`);
        }
      });
    }

    // Test tool commands (will fail if not in game, but shows API works)
    console.log('\n[Bot] Testing tool commands...');
    try {
      await bot.getTools().faceDirection(1);
      console.log('[Bot] faceDirection(1) - OK');
    } catch (e: any) {
      console.log(`[Bot] faceDirection(1) - ${e.message}`);
    }

    console.log('\n[Bot] Test complete! Press Ctrl+C to stop.');

  } catch (error: any) {
    console.error('[Bot] Failed to start:', error.message);
    process.exit(1);
  }
}

main();
