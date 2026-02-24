// Type definitions for Stardew MCP Server

export interface WebSocketMessage {
  id?: string;
  type: string;
  action?: string;
  params?: Record<string, any>;
}

export interface WebSocketResponse {
  id?: string;
  type: string;
  success: boolean;
  message?: string;
  data?: any;
}

export interface GameState {
  player: PlayerState;
  time: TimeState;
  world: WorldState;
  surroundings: SurroundingsState;
  map: MapInfo;
  quests?: QuestInfo[];
  relationships?: RelationshipInfo[];
  skills?: SkillsInfo;
}

export interface PlayerState {
  name: string;
  x: number;
  y: number;
  location: string;
  energy: number;
  maxEnergy: number;
  health: number;
  maxHealth: number;
  money: number;
  currentTool: string;
  currentToolIndex: number;
  facingDirection: number;
  facingDirectionName: string;
  isMoving: boolean;
  canMove: boolean;
  inventory: InventoryItem[];
}

export interface TimeState {
  timeOfDay: number;
  timeString: string;
  day: number;
  season: string;
  year: number;
  dayOfWeek: string;
  isNight: boolean;
  minutesUntilMorning: number;
}

export interface WorldState {
  weather: string;
  isOutdoors: boolean;
  isFarm: boolean;
  isGreenhouse: boolean;
  isBuildableLocation: boolean;
  locationType: string;
}

export interface SurroundingsState {
  asciiMap: string;
  nearbyObjects: NearbyObject[];
  nearbyTerrainFeatures: NearbyTerrain[];
  nearbyNPCs: NearbyNPC[];
  nearbyMonsters: NearbyMonster[];
  nearbyResourceClumps: NearbyResourceClump[];
  nearbyDebris: NearbyDebris[];
  nearbyBuildings: NearbyBuilding[];
  nearbyAnimals: NearbyAnimal[];
  warpPoints: WarpPoint[];
  tileInFront: TileInFront;
}

export interface NearbyObject {
  x: number;
  y: number;
  name: string;
  displayName: string;
  type: string;
  isPassable: boolean;
  canBePickedUp: boolean;
  isReadyForHarvest: boolean;
  minutesUntilReady: number;
  heldItemName?: string;
  requiredTool?: string;
  hitsRequired: number;
}

export interface NearbyTerrain {
  x: number;
  y: number;
  type: string;
  isPassable: boolean;
  growthStage: number;
  isFullyGrown: boolean;
  hasSeed: boolean;
  canBeChopped: boolean;
  fruitCount: number;
  isWatered: boolean;
  hasCrop: boolean;
  cropName?: string;
  cropPhase: number;
  daysUntilHarvest: number;
  isReadyForHarvest: boolean;
  isDead: boolean;
  grassType: number;
  requiredTool?: string;
  hitsRequired: number;
}

export interface NearbyNPC {
  x: number;
  y: number;
  name: string;
  displayName: string;
  isFacingPlayer: boolean;
  canTalk: boolean;
  friendshipLevel: number;
  isMoving: boolean;
}

export interface NearbyMonster {
  x: number;
  y: number;
  name: string;
  health: number;
  maxHealth: number;
  damageToFarmer: number;
  isGlider: boolean;
  distance: number;
}

export interface NearbyResourceClump {
  x: number;
  y: number;
  width: number;
  height: number;
  type: string;
  health: number;
  requiredTool?: string;
  hitsRequired: number;
}

export interface NearbyDebris {
  x: number;
  y: number;
  name: string;
  type: string;
  canBePickedUp: boolean;
}

export interface NearbyBuilding {
  x: number;
  y: number;
  width: number;
  height: number;
  type: string;
  doorX: number;
  doorY: number;
  hasAnimals: boolean;
}

export interface NearbyAnimal {
  x: number;
  y: number;
  name: string;
  type: string;
  age: number;
  happiness: number;
  canBePet: boolean;
  hasProduce: boolean;
}

export interface WarpPoint {
  x: number;
  y: number;
  targetLocation: string;
  targetX: number;
  targetY: number;
  isDoor: boolean;
}

export interface TileInFront {
  x: number;
  y: number;
  isPassable: boolean;
  isWater: boolean;
  isTillable: boolean;
  objectName?: string;
  objectType?: string;
  terrainType?: string;
  npcName?: string;
  canInteract: boolean;
  requiredTool?: string;
}

export interface InventoryItem {
  slot: number;
  name: string;
  displayName: string;
  stack: number;
  category: string;
  isTool: boolean;
  isWeapon: boolean;
}

export interface QuestInfo {
  id: string;
  name: string;
  description: string;
  objective: string;
  daysLeft: number;
  reward: number;
  isComplete: boolean;
}

export interface RelationshipInfo {
  npcName: string;
  friendshipPoints: number;
  hearts: number;
  giftsToday: number;
  giftsThisWeek: number;
  talkedToToday: boolean;
  status: string;
}

export interface SkillsInfo {
  farming: number;
  mining: number;
  foraging: number;
  fishing: number;
  combat: number;
  farmingXp: number;
  miningXp: number;
  foragingXp: number;
  fishingXp: number;
  combatXp: number;
}

export interface Config {
  server: {
    gameUrl: string;
    autoStart: boolean;
    logLevel: string;
  };
  remote: {
    host: string;
    port: number;
  };
  openclaw: {
    gatewayUrl: string;
    token: string;
    agentName: string;
  };
  agent: {
    defaultGoal: string;
    llmTimeout: number;
    cheatMode: boolean;
  };
}
