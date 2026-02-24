import { GameClient } from './gameClient';
import { GameState } from './types';

export class StardewTools {
  private gameClient: GameClient;

  constructor(gameClient: GameClient) {
    this.gameClient = gameClient;
  }

  // Core Game Tools

  async getState(): Promise<GameState | null> {
    return this.gameClient.getState();
  }

  async getSurroundings(): Promise<any> {
    const response = await this.gameClient.sendCommand('get_surroundings');
    return response.data;
  }

  async moveTo(x: number, y: number): Promise<string> {
    const response = await this.gameClient.sendCommand('move_to', { x, y });
    return response.message || 'Moved';
  }

  async interact(): Promise<string> {
    const response = await this.gameClient.sendCommand('interact');
    return response.message || 'Interacted';
  }

  async useTool(): Promise<string> {
    const response = await this.gameClient.sendCommand('use_tool');
    return response.message || 'Tool used';
  }

  async useToolRepeat(count: number, direction: number = 0): Promise<string> {
    const response = await this.gameClient.sendCommand('use_tool_repeat', { count, direction });
    return response.message || 'Tools used';
  }

  async faceDirection(direction: number): Promise<string> {
    const response = await this.gameClient.sendCommand('face_direction', { direction });
    return response.message || 'Direction changed';
  }

  async selectItem(slot: number): Promise<string> {
    const response = await this.gameClient.sendCommand('select_item', { slot });
    return response.message || 'Item selected';
  }

  async switchTool(tool: string): Promise<string> {
    const response = await this.gameClient.sendCommand('switch_tool', { tool });
    return response.message || 'Tool switched';
  }

  async eatItem(slot: number): Promise<string> {
    const response = await this.gameClient.sendCommand('eat_item', { slot });
    return response.message || 'Item eaten';
  }

  async enterDoor(): Promise<string> {
    const response = await this.gameClient.sendCommand('enter_door');
    return response.message || 'Door entered';
  }

  async findBestTarget(type: string): Promise<string> {
    const response = await this.gameClient.sendCommand('find_best_target', { type });
    return response.message || 'Target found';
  }

  async clearTarget(): Promise<string> {
    const response = await this.gameClient.sendCommand('clear_target');
    return response.message || 'Target cleared';
  }

  // Cheat Mode Tools

  async cheatModeEnable(): Promise<string> {
    const response = await this.gameClient.sendCommand('cheat_mode_enable');
    return response.message || 'Cheat mode enabled';
  }

  async cheatModeDisable(): Promise<string> {
    const response = await this.gameClient.sendCommand('cheat_mode_disable');
    return response.message || 'Cheat mode disabled';
  }

  async cheatTimeFreeze(freeze: boolean): Promise<string> {
    const response = await this.gameClient.sendCommand('cheat_time_freeze', { freeze });
    return response.message || 'Time freeze toggled';
  }

  async cheatInfiniteEnergy(enable: boolean): Promise<string> {
    const response = await this.gameClient.sendCommand('cheat_infinite_energy', { enable });
    return response.message || 'Infinite energy toggled';
  }

  async cheatWarp(location: string): Promise<string> {
    const response = await this.gameClient.sendCommand('cheat_warp', { location });
    return response.message || 'Warped';
  }

  async cheatMineWarp(level: number): Promise<string> {
    const response = await this.gameClient.sendCommand('cheat_mine_warp', { level });
    return response.message || 'Warped to mine';
  }

  async cheatClearDebris(): Promise<string> {
    const response = await this.gameClient.sendCommand('cheat_clear_debris');
    return response.message || 'Debris cleared';
  }

  async cheatCutTrees(): Promise<string> {
    const response = await this.gameClient.sendCommand('cheat_cut_trees');
    return response.message || 'Trees cut';
  }

  async cheatMineRocks(): Promise<string> {
    const response = await this.gameClient.sendCommand('cheat_mine_rocks');
    return response.message || 'Rocks mined';
  }

  async cheatHoeAll(): Promise<string> {
    const response = await this.gameClient.sendCommand('cheat_hoe_all');
    return response.message || 'All tiles hoed';
  }

  async cheatWaterAll(): Promise<string> {
    const response = await this.gameClient.sendCommand('cheat_water_all');
    return response.message || 'All crops watered';
  }

  async cheatPlantSeeds(season: string, seedId?: string): Promise<string> {
    const response = await this.gameClient.sendCommand('cheat_plant_seeds', { season, seedId });
    return response.message || 'Seeds planted';
  }

  async cheatFertilizeAll(type: string): Promise<string> {
    const response = await this.gameClient.sendCommand('cheat_fertilize_all', { type });
    return response.message || 'Fertilizer applied';
  }

  async cheatGrowCrops(): Promise<string> {
    const response = await this.gameClient.sendCommand('cheat_grow_crops');
    return response.message || 'Crops grown';
  }

  async cheatHarvestAll(): Promise<string> {
    const response = await this.gameClient.sendCommand('cheat_harvest_all');
    return response.message || 'Crops harvested';
  }

  async cheatDigArtifacts(): Promise<string> {
    const response = await this.gameClient.sendCommand('cheat_dig_artifacts');
    return response.message || 'Artifacts dug';
  }

  async cheatSetMoney(amount: number): Promise<string> {
    const response = await this.gameClient.sendCommand('cheat_set_money', { amount });
    return response.message || 'Money set';
  }

  async cheatAddItem(id: string, count: number = 1): Promise<string> {
    const response = await this.gameClient.sendCommand('cheat_add_item', { id, count });
    return response.message || 'Item added';
  }

  async cheatSpawnOres(type: string): Promise<string> {
    const response = await this.gameClient.sendCommand('cheat_spawn_ores', { type });
    return response.message || 'Ores spawned';
  }

  async cheatSetFriendship(npc: string, points: number): Promise<string> {
    const response = await this.gameClient.sendCommand('cheat_set_friendship', { npc, points });
    return response.message || 'Friendship set';
  }

  async cheatMaxAllFriendships(): Promise<string> {
    const response = await this.gameClient.sendCommand('cheat_max_all_friendships');
    return response.message || 'All friendships maxed';
  }

  async cheatGiveGift(npc: string, itemId: string): Promise<string> {
    const response = await this.gameClient.sendCommand('cheat_give_gift', { npc, itemId });
    return response.message || 'Gift given';
  }

  async cheatUpgradeBackpack(level: number): Promise<string> {
    const response = await this.gameClient.sendCommand('cheat_upgrade_backpack', { level });
    return response.message || 'Backpack upgraded';
  }

  async cheatUpgradeTool(tool: string, level: number): Promise<string> {
    const response = await this.gameClient.sendCommand('cheat_upgrade_tool', { tool, level });
    return response.message || 'Tool upgraded';
  }

  async cheatUpgradeAllTools(level: number): Promise<string> {
    const response = await this.gameClient.sendCommand('cheat_upgrade_all_tools', { level });
    return response.message || 'All tools upgraded';
  }

  async cheatUnlockAll(): Promise<string> {
    const response = await this.gameClient.sendCommand('cheat_unlock_all');
    return response.message || 'Everything unlocked';
  }

  // Get all tools as object for rs-sdk
  getToolsObject() {
    return {
      get_state: () => this.getState(),
      get_surroundings: () => this.getSurroundings(),
      move_to: (params: { x: number; y: number }) => this.moveTo(params.x, params.y),
      interact: () => this.interact(),
      use_tool: () => this.useTool(),
      use_tool_repeat: (params: { count: number; direction?: number }) =>
        this.useToolRepeat(params.count, params.direction || 0),
      face_direction: (params: { direction: number }) => this.faceDirection(params.direction),
      select_item: (params: { slot: number }) => this.selectItem(params.slot),
      switch_tool: (params: { tool: string }) => this.switchTool(params.tool),
      eat_item: (params: { slot: number }) => this.eatItem(params.slot),
      enter_door: () => this.enterDoor(),
      find_best_target: (params: { type: string }) => this.findBestTarget(params.type),
      clear_target: () => this.clearTarget(),
      cheat_mode_enable: () => this.cheatModeEnable(),
      cheat_mode_disable: () => this.cheatModeDisable(),
      cheat_time_freeze: (params: { freeze: boolean }) => this.cheatTimeFreeze(params.freeze),
      cheat_infinite_energy: (params: { enable: boolean }) => this.cheatInfiniteEnergy(params.enable),
      cheat_warp: (params: { location: string }) => this.cheatWarp(params.location),
      cheat_mine_warp: (params: { level: number }) => this.cheatMineWarp(params.level),
      cheat_clear_debris: () => this.cheatClearDebris(),
      cheat_cut_trees: () => this.cheatCutTrees(),
      cheat_mine_rocks: () => this.cheatMineRocks(),
      cheat_hoe_all: () => this.cheatHoeAll(),
      cheat_water_all: () => this.cheatWaterAll(),
      cheat_plant_seeds: (params: { season: string; seedId?: string }) =>
        this.cheatPlantSeeds(params.season, params.seedId),
      cheat_fertilize_all: (params: { type: string }) => this.cheatFertilizeAll(params.type),
      cheat_grow_crops: () => this.cheatGrowCrops(),
      cheat_harvest_all: () => this.cheatHarvestAll(),
      cheat_dig_artifacts: () => this.cheatDigArtifacts(),
      cheat_set_money: (params: { amount: number }) => this.cheatSetMoney(params.amount),
      cheat_add_item: (params: { id: string; count?: number }) =>
        this.cheatAddItem(params.id, params.count || 1),
      cheat_spawn_ores: (params: { type: string }) => this.cheatSpawnOres(params.type),
      cheat_set_friendship: (params: { npc: string; points: number }) =>
        this.cheatSetFriendship(params.npc, params.points),
      cheat_max_all_friendships: () => this.cheatMaxAllFriendships(),
      cheat_give_gift: (params: { npc: string; item_id: string }) =>
        this.cheatGiveGift(params.npc, params.item_id),
      cheat_upgrade_backpack: (params: { level: number }) => this.cheatUpgradeBackpack(params.level),
      cheat_upgrade_tool: (params: { tool: string; level: number }) =>
        this.cheatUpgradeTool(params.tool, params.level),
      cheat_upgrade_all_tools: (params: { level: number }) => this.cheatUpgradeAllTools(params.level),
      cheat_unlock_all: () => this.cheatUnlockAll(),
    };
  }
}
