package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"os"
	"sync"
	"time"

	"github.com/gorilla/websocket"
	"gopkg.in/yaml.v3"
)

// Config represents the configuration file structure
type Config struct {
	Server struct {
		GameURL   string `yaml:"game_url"`
		AutoStart bool   `yaml:"auto_start"`
		LogLevel  string `yaml:"log_level"`
	} `yaml:"server"`
	Agent struct {
		DefaultGoal string `yaml:"default_goal"`
		LLMTimeout  int    `yaml:"llm_timeout"`
		CheatMode   bool   `yaml:"cheat_mode"`
	} `yaml:"agent"`
	OpenClaw struct {
		Enabled    bool   `yaml:"enabled"`
		GatewayURL string `yaml:"gateway_url"`
		AgentName  string `yaml:"agent_name"`
		Workspace  string `yaml:"workspace"`
	} `yaml:"openclaw"`
}

// LoadConfig reads configuration from YAML file
func LoadConfig(path string) *Config {
	cfg := &Config{}

	// Set defaults
	cfg.Server.GameURL = "ws://localhost:8765/game"
	cfg.Server.AutoStart = true
	cfg.Server.LogLevel = "info"
	cfg.Agent.DefaultGoal = "Setup and manage the farm efficiently using available tools"
	cfg.Agent.LLMTimeout = 60
	cfg.Agent.CheatMode = false
	cfg.OpenClaw.Enabled = false
	cfg.OpenClaw.GatewayURL = "ws://127.0.0.1:18789"
	cfg.OpenClaw.AgentName = "stardew-farmer"
	cfg.OpenClaw.Workspace = "~/.openclaw/workspace/stardew"

	// Try to read config file if it exists
	if data, err := os.ReadFile(path); err == nil {
		if err := yaml.Unmarshal(data, cfg); err != nil {
			log.Printf("Warning: Failed to parse config file: %v", err)
		}
	}

	return cfg
}

// GameClient manages the WebSocket connection to the Stardew Valley mod
type GameClient struct {
	conn        *websocket.Conn
	mu          sync.RWMutex
	state       *GameState
	responses   map[string]chan *WebSocketResponse
	responsesMu sync.Mutex
	connected   bool
	url         string
}

// GameState represents the current state of the game
type GameState struct {
	Player        PlayerState        `json:"player"`
	Time          TimeState          `json:"time"`
	World         WorldState         `json:"world"`
	Surroundings  SurroundingsState  `json:"surroundings"`
	Map           MapInfo            `json:"map"`
	Quests        []QuestInfo        `json:"quests,omitempty"`
	Relationships []RelationshipInfo `json:"relationships,omitempty"`
	Skills        *SkillsInfo        `json:"skills,omitempty"`
}

type PlayerState struct {
	Name                string          `json:"name"`
	X                   int             `json:"x"`
	Y                   int             `json:"y"`
	Location            string          `json:"location"`
	Energy              float64         `json:"energy"`
	MaxEnergy           int             `json:"maxEnergy"`
	Health              int             `json:"health"`
	MaxHealth           int             `json:"maxHealth"`
	Money               int             `json:"money"`
	CurrentTool         string          `json:"currentTool"`
	CurrentToolIndex    int             `json:"currentToolIndex"`
	FacingDirection     int             `json:"facingDirection"`
	FacingDirectionName string          `json:"facingDirectionName"`
	IsMoving            bool            `json:"isMoving"`
	CanMove             bool            `json:"canMove"`
	Inventory           []InventoryItem `json:"inventory"`
}

type TimeState struct {
	TimeOfDay           int    `json:"timeOfDay"`
	TimeString          string `json:"timeString"`
	Day                 int    `json:"day"`
	Season              string `json:"season"`
	Year                int    `json:"year"`
	DayOfWeek           string `json:"dayOfWeek"`
	IsNight             bool   `json:"isNight"`
	MinutesUntilMorning int    `json:"minutesUntilMorning"`
}

type WorldState struct {
	Weather             string `json:"weather"`
	IsOutdoors          bool   `json:"isOutdoors"`
	IsFarm              bool   `json:"isFarm"`
	IsGreenhouse        bool   `json:"isGreenhouse"`
	IsBuildableLocation bool   `json:"isBuildableLocation"`
	LocationType        string `json:"locationType"`
}

type SurroundingsState struct {
	AsciiMap              string                `json:"asciiMap"`
	NearbyObjects         []NearbyObject        `json:"nearbyObjects"`
	NearbyTerrainFeatures []NearbyTerrain       `json:"nearbyTerrainFeatures"`
	NearbyNPCs            []NearbyNPC           `json:"nearbyNPCs"`
	NearbyMonsters        []NearbyMonster       `json:"nearbyMonsters"`
	NearbyResourceClumps  []NearbyResourceClump `json:"nearbyResourceClumps"`
	NearbyDebris          []NearbyDebris        `json:"nearbyDebris"`
	NearbyBuildings       []NearbyBuilding      `json:"nearbyBuildings"`
	NearbyAnimals         []NearbyAnimal        `json:"nearbyAnimals"`
	WarpPoints            []WarpPoint           `json:"warpPoints"`
	TileInFront           TileInFront           `json:"tileInFront"`
}

type MapInfo struct {
	Name        string `json:"name"`
	DisplayName string `json:"displayName"`
	Width       int    `json:"width"`
	Height      int    `json:"height"`
	IsMineLevel bool   `json:"isMineLevel"`
	MineLevel   int    `json:"mineLevel"`
	UniqueId    string `json:"uniqueId"`
}

type TileInfo struct {
	X                 int    `json:"x"`
	Y                 int    `json:"y"`
	IsPassable        bool   `json:"isPassable"`
	IsWater           bool   `json:"isWater"`
	IsTillable        bool   `json:"isTillable"`
	HasObject         bool   `json:"hasObject"`
	HasTerrainFeature bool   `json:"hasTerrainFeature"`
	TileType          string `json:"tileType"`
}

type NearbyObject struct {
	X                 int    `json:"x"`
	Y                 int    `json:"y"`
	Name              string `json:"name"`
	DisplayName       string `json:"displayName"`
	Type              string `json:"type"`
	IsPassable        bool   `json:"isPassable"`
	CanBePickedUp     bool   `json:"canBePickedUp"`
	IsReadyForHarvest bool   `json:"isReadyForHarvest"`
	MinutesUntilReady int    `json:"minutesUntilReady"`
	HeldItemName      string `json:"heldItemName,omitempty"`
	RequiredTool      string `json:"requiredTool,omitempty"`
	HitsRequired      int    `json:"hitsRequired"`
}

type NearbyTerrain struct {
	X                 int    `json:"x"`
	Y                 int    `json:"y"`
	Type              string `json:"type"`
	IsPassable        bool   `json:"isPassable"`
	GrowthStage       int    `json:"growthStage"`
	IsFullyGrown      bool   `json:"isFullyGrown"`
	HasSeed           bool   `json:"hasSeed"`
	CanBeChopped      bool   `json:"canBeChopped"`
	FruitCount        int    `json:"fruitCount"`
	IsWatered         bool   `json:"isWatered"`
	HasCrop           bool   `json:"hasCrop"`
	CropName          string `json:"cropName,omitempty"`
	CropPhase         int    `json:"cropPhase"`
	DaysUntilHarvest  int    `json:"daysUntilHarvest"`
	IsReadyForHarvest bool   `json:"isReadyForHarvest"`
	IsDead            bool   `json:"isDead"`
	GrassType         int    `json:"grassType"`
	RequiredTool      string `json:"requiredTool,omitempty"`
	HitsRequired      int    `json:"hitsRequired"`
}

type NearbyNPC struct {
	X               int    `json:"x"`
	Y               int    `json:"y"`
	Name            string `json:"name"`
	DisplayName     string `json:"displayName"`
	IsFacingPlayer  bool   `json:"isFacingPlayer"`
	CanTalk         bool   `json:"canTalk"`
	FriendshipLevel int    `json:"friendshipLevel"`
	IsMoving        bool   `json:"isMoving"`
}

type NearbyMonster struct {
	X              int    `json:"x"`
	Y              int    `json:"y"`
	Name           string `json:"name"`
	Health         int    `json:"health"`
	MaxHealth      int    `json:"maxHealth"`
	DamageToFarmer int    `json:"damageToFarmer"`
	IsGlider       bool   `json:"isGlider"`
	Distance       int    `json:"distance"`
}

type NearbyResourceClump struct {
	X            int     `json:"x"`
	Y            int     `json:"y"`
	Width        int     `json:"width"`
	Height       int     `json:"height"`
	Type         string  `json:"type"`
	Health       float64 `json:"health"`
	RequiredTool string  `json:"requiredTool,omitempty"`
	HitsRequired int     `json:"hitsRequired"`
}

type NearbyDebris struct {
	X             int    `json:"x"`
	Y             int    `json:"y"`
	Name          string `json:"name"`
	Type          string `json:"type"`
	CanBePickedUp bool   `json:"canBePickedUp"`
}

type NearbyBuilding struct {
	X          int    `json:"x"`
	Y          int    `json:"y"`
	Width      int    `json:"width"`
	Height     int    `json:"height"`
	Type       string `json:"type"`
	DoorX      int    `json:"doorX"`
	DoorY      int    `json:"doorY"`
	HasAnimals bool   `json:"hasAnimals"`
}

type NearbyAnimal struct {
	X          int    `json:"x"`
	Y          int    `json:"y"`
	Name       string `json:"name"`
	Type       string `json:"type"`
	Age        int    `json:"age"`
	Happiness  int    `json:"happiness"`
	CanBePet   bool   `json:"canBePet"`
	HasProduce bool   `json:"hasProduce"`
}

type WarpPoint struct {
	X              int    `json:"x"`
	Y              int    `json:"y"`
	TargetLocation string `json:"targetLocation"`
	TargetX        int    `json:"targetX"`
	TargetY        int    `json:"targetY"`
	IsDoor         bool   `json:"isDoor"`
}

type TileInFront struct {
	X            int    `json:"x"`
	Y            int    `json:"y"`
	IsPassable   bool   `json:"isPassable"`
	IsWater      bool   `json:"isWater"`
	IsTillable   bool   `json:"isTillable"`
	ObjectName   string `json:"objectName,omitempty"`
	ObjectType   string `json:"objectType,omitempty"`
	TerrainType  string `json:"terrainType,omitempty"`
	NPCName      string `json:"npcName,omitempty"`
	CanInteract  bool   `json:"canInteract"`
	RequiredTool string `json:"requiredTool,omitempty"`
}

type InventoryItem struct {
	Slot        int    `json:"slot"`
	Name        string `json:"name"`
	DisplayName string `json:"displayName"`
	Stack       int    `json:"stack"`
	Category    string `json:"category"`
	IsTool      bool   `json:"isTool"`
	IsWeapon    bool   `json:"isWeapon"`
}

type QuestInfo struct {
	ID          string `json:"id"`
	Name        string `json:"name"`
	Description string `json:"description"`
	Objective   string `json:"objective"`
	DaysLeft    int    `json:"daysLeft"`
	Reward      int    `json:"reward"`
	IsComplete  bool   `json:"isComplete"`
}

type RelationshipInfo struct {
	NPCName          string `json:"npcName"`
	FriendshipPoints int    `json:"friendshipPoints"`
	Hearts           int    `json:"hearts"`
	GiftsToday       int    `json:"giftsToday"`
	GiftsThisWeek    int    `json:"giftsThisWeek"`
	TalkedToToday    bool   `json:"talkedToToday"`
	Status           string `json:"status"`
}

type SkillsInfo struct {
	Farming    int `json:"farming"`
	Mining     int `json:"mining"`
	Foraging   int `json:"foraging"`
	Fishing    int `json:"fishing"`
	Combat     int `json:"combat"`
	FarmingXp  int `json:"farmingXp"`
	MiningXp   int `json:"miningXp"`
	ForagingXp int `json:"foragingXp"`
	FishingXp  int `json:"fishingXp"`
	CombatXp   int `json:"combatXp"`
}

type WebSocketMessage struct {
	ID     string                 `json:"id,omitempty"`
	Type   string                 `json:"type"`
	Action string                 `json:"action,omitempty"`
	Params map[string]interface{} `json:"params,omitempty"`
}

type WebSocketResponse struct {
	ID      string      `json:"id,omitempty"`
	Type    string      `json:"type"`
	Success bool        `json:"success"`
	Message string      `json:"message,omitempty"`
	Data    interface{} `json:"data,omitempty"`
}

var gameClient *GameClient

func NewGameClient() *GameClient {
	return &GameClient{
		responses: make(map[string]chan *WebSocketResponse),
	}
}

func (c *GameClient) Connect(url string) error {
	c.url = url
	conn, _, err := websocket.DefaultDialer.Dial(url, nil)
	if err != nil {
		return fmt.Errorf("failed to connect to game: %w", err)
	}

	c.mu.Lock()
	c.conn = conn
	c.connected = true
	c.mu.Unlock()

	go c.listen()
	go c.keepAlive()

	return nil
}

func (c *GameClient) keepAlive() {
	ticker := time.NewTicker(15 * time.Second)
	defer ticker.Stop()

	for range ticker.C {
		c.mu.RLock()
		conn := c.conn
		connected := c.connected
		c.mu.RUnlock()

		if !connected || conn == nil {
			return
		}

		ping := WebSocketMessage{
			ID:   fmt.Sprintf("%d", time.Now().UnixNano()),
			Type: "ping",
		}
		data, _ := json.Marshal(ping)

		c.mu.Lock()
		err := conn.WriteMessage(websocket.TextMessage, data)
		c.mu.Unlock()

		if err != nil {
			log.Printf("Ping failed: %v", err)
			return
		}
	}
}

func (c *GameClient) reconnect() {
	c.mu.Lock()
	c.connected = false
	if c.conn != nil {
		c.conn.Close()
	}
	c.mu.Unlock()

	for {
		log.Printf("Attempting to reconnect...")
		time.Sleep(5 * time.Second)

		conn, _, err := websocket.DefaultDialer.Dial(c.url, nil)
		if err != nil {
			log.Printf("Reconnect failed: %v", err)
			continue
		}

		c.mu.Lock()
		c.conn = conn
		c.connected = true
		c.mu.Unlock()

		log.Printf("Reconnected to Stardew Valley at %s", c.url)

		go c.listen()
		go c.keepAlive()
		return
	}
}

func (c *GameClient) listen() {
	for {
		c.mu.RLock()
		conn := c.conn
		c.mu.RUnlock()

		if conn == nil {
			return
		}

		_, message, err := conn.ReadMessage()
		if err != nil {
			log.Printf("WebSocket read error from %s: %v", c.url, err)
			go c.reconnect()
			return
		}

		var response WebSocketResponse
		if err := json.Unmarshal(message, &response); err != nil {
			log.Printf("Failed to parse response: %v", err)
			continue
		}

		switch response.Type {
		case "state":
			c.handleStateUpdate(&response)
		case "response":
			c.handleCommandResponse(&response)
		case "pong":
			// Heartbeat response, ignore
		case "error":
			log.Printf("Error from game: %s", response.Message)
		}
	}
}

func (c *GameClient) handleStateUpdate(response *WebSocketResponse) {
	data, err := json.Marshal(response.Data)
	if err != nil {
		log.Printf("Failed to marshal state data: %v", err)
		return
	}

	var state GameState
	if err := json.Unmarshal(data, &state); err != nil {
		log.Printf("Failed to parse state: %v", err)
		return
	}

	c.mu.Lock()
	c.state = &state
	c.mu.Unlock()
}

func (c *GameClient) handleCommandResponse(response *WebSocketResponse) {
	if response.ID == "" {
		return
	}

	c.responsesMu.Lock()
	ch, ok := c.responses[response.ID]
	if ok {
		delete(c.responses, response.ID)
	}
	c.responsesMu.Unlock()

	if ok {
		ch <- response
	}
}

func (c *GameClient) GetState() *GameState {
	c.mu.RLock()
	defer c.mu.RUnlock()
	return c.state
}

func (c *GameClient) IsConnected() bool {
	c.mu.RLock()
	defer c.mu.RUnlock()
	return c.connected
}

func (c *GameClient) SendCommand(action string, params map[string]interface{}) (*WebSocketResponse, error) {
	if !c.IsConnected() {
		return nil, fmt.Errorf("not connected to game")
	}

	id := fmt.Sprintf("%d", time.Now().UnixNano())

	msg := WebSocketMessage{
		ID:     id,
		Type:   "command",
		Action: action,
		Params: params,
	}

	ch := make(chan *WebSocketResponse, 1)
	c.responsesMu.Lock()
	c.responses[id] = ch
	c.responsesMu.Unlock()

	data, err := json.Marshal(msg)
	if err != nil {
		return nil, err
	}

	c.mu.Lock()
	err = c.conn.WriteMessage(websocket.TextMessage, data)
	c.mu.Unlock()

	if err != nil {
		c.responsesMu.Lock()
		delete(c.responses, id)
		c.responsesMu.Unlock()
		return nil, err
	}

	// Timeout for command responses (15 seconds is sufficient for most operations)
	select {
	case response := <-ch:
		return response, nil
	case <-time.After(15 * time.Second):
		c.responsesMu.Lock()
		delete(c.responses, id)
		c.responsesMu.Unlock()
		return nil, fmt.Errorf("timeout waiting for response")
	}
}

func main() {
	// Load configuration
	configPath := flag.String("config", "config.yaml", "Path to configuration file")
	autoFlag := flag.Bool("auto", false, "Start in autonomous mode (overrides config)")
	goalFlag := flag.String("goal", "", "Goal for autonomous mode (overrides config)")
	urlFlag := flag.String("url", "", "WebSocket URL for the game mod (overrides config)")
	openclawFlag := flag.Bool("openclaw", false, "Enable OpenClaw Gateway mode")
	flag.Parse()

	cfg := LoadConfig(*configPath)

	// Command-line flags override config file
	gameURL := *urlFlag
	if gameURL == "" {
		gameURL = cfg.Server.GameURL
	}

	autoStart := *autoFlag
	if !flag.Parsed() || autoStart == false && !*autoFlag {
		autoStart = cfg.Server.AutoStart
	}

	goal := *goalFlag
	if goal == "" {
		goal = cfg.Agent.DefaultGoal
	}

	// Set log level
	if cfg.Server.LogLevel == "debug" {
		log.SetFlags(log.LstdFlags | log.Lshortfile)
	}

	log.Printf("Stardew MCP Server starting...")
	log.Printf("Game URL: %s", gameURL)
	log.Printf("Auto-start: %v", autoStart)

	gameClient = NewGameClient()

	// Check for OpenClaw Gateway mode
	if *openclawFlag || cfg.OpenClaw.Enabled {
		log.Printf("OpenClaw Gateway mode enabled")
		startOpenClawGateway(cfg, gameURL, autoStart, goal)
		return
	}

	go func() {
		for {
			if err := gameClient.Connect(gameURL); err != nil {
				log.Printf("Failed to connect to game (will retry): %v", err)
				time.Sleep(5 * time.Second)
				continue
			}
			log.Println("Connected to Stardew Valley!")

			if autoStart {
				log.Printf("Starting autonomous agent with goal: %s", goal)

				agent, err := NewStardewAgent()
				if err != nil {
					log.Printf("Failed to start agent: %v", err)
					return
				}
				if err := agent.StartSession(goal); err != nil {
					log.Printf("Failed to start session: %v", err)
					return
				}
			}
			break
		}
	}()

	// Block forever
	select {}
}

// startOpenClawGateway starts the OpenClaw Gateway integration
func startOpenClawGateway(cfg *Config, gameURL string, autoStart bool, goal string) {
	// Connect to game first
	if err := gameClient.Connect(gameURL); err != nil {
		log.Printf("Failed to connect to game: %v", err)
		return
	}
	log.Println("Connected to Stardew Valley!")

	// Connect to OpenClaw Gateway
	gatewayURL := cfg.OpenClaw.GatewayURL
	log.Printf("Connecting to OpenClaw Gateway at %s...", gatewayURL)

	gwConn, _, err := websocket.DefaultDialer.Dial(gatewayURL, nil)
	if err != nil {
		log.Printf("Failed to connect to OpenClaw Gateway: %v", err)
		log.Println("Falling back to standalone mode...")
		if autoStart {
			startAutonomousAgent(goal)
		}
		return
	}
	defer gwConn.Close()
	log.Println("Connected to OpenClaw Gateway!")

	// Register this tool server
	registerMessage := map[string]interface{}{
		"type":        "tool.register",
		"name":        "stardew-mcp",
		"description": "Stardew Valley AI Controller - Farm management and game control",
		"version":     "1.0.0",
		"tools":       getStardewTools(),
	}
	gwConn.WriteJSON(registerMessage)

	// Start autonomous agent if enabled
	if autoStart {
		startAutonomousAgent(goal)
	}

	// Forward messages between game and OpenClaw Gateway
	go func() {
		for {
			_, msg, err := gwConn.ReadMessage()
			if err != nil {
				log.Printf("Gateway read error: %v", err)
				return
			}

			// Handle tool calls from OpenClaw
			var toolMsg map[string]interface{}
			if err := json.Unmarshal(msg, &toolMsg); err != nil {
				continue
			}

			if toolMsg["type"] == "tool.call" {
				toolName := toolMsg["name"].(string)
				params := toolMsg["params"].(map[string]interface{})

				result, err := executeTool(toolName, params)
				response := map[string]interface{}{
					"type":     "tool.result",
					"call_id":  toolMsg["call_id"],
					"success": err == nil,
					"result":   result,
				}
				if err != nil {
					response["error"] = err.Error()
				}
				gwConn.WriteJSON(response)
			}
		}
	}()

	// Block forever
	select {}
}

func startAutonomousAgent(goal string) {
	agent, err := NewStardewAgent()
	if err != nil {
		log.Printf("Failed to start agent: %v", err)
		return
	}
	if err := agent.StartSession(goal); err != nil {
		log.Printf("Failed to start session: %v", err)
		return
	}
}

// getStardewTools returns the tool definitions for OpenClaw
func getStardewTools() []map[string]interface{} {
	return []map[string]interface{}{
		{
			"name":        "get_state",
			"description": "Get current game state including player position, inventory, time, and surroundings",
		},
		{
			"name":        "get_surroundings",
			"description": "Get detailed information about tiles around the player",
		},
		{
			"name":        "move_to",
			"description": "Move player to specified coordinates",
			"params": map[string]interface{}{
				"x": "target X coordinate",
				"y": "target Y coordinate",
			},
		},
		{
			"name":        "interact",
			"description": "Interact with object in front of player",
		},
		{
			"name":        "use_tool",
			"description": "Use currently selected tool",
		},
		{
			"name":        "select_item",
			"description": "Select item from inventory by slot number",
			"params": map[string]interface{}{
				"slot": "inventory slot number (0-11)",
			},
		},
		{
			"name":        "switch_tool",
			"description": "Switch to tool by name",
			"params": map[string]interface{}{
				"tool": "tool name (hoe, watering_can, pickaxe, etc.)",
			},
		},
		{
			"name":        "face_direction",
			"description": "Face a direction",
			"params": map[string]interface{}{
				"direction": "0=down, 1=left, 2=right, 3=up",
			},
		},
		{
			"name":        "cheat_mode_enable",
			"description": "Enable god-mode cheat commands",
		},
		{
			"name":        "cheat_warp",
			"description": "Teleport to location",
			"params": map[string]interface{}{
				"location": "location name (Farm, Town, Mine, etc.)",
			},
		},
		{
			"name":        "cheat_set_money",
			"description": "Set money amount",
			"params": map[string]interface{}{
				"amount": "amount of gold",
			},
		},
	}
}

// executeTool executes a tool call from OpenClaw
func executeTool(name string, params map[string]interface{}) (interface{}, error) {
	switch name {
	case "get_state":
		return gameClient.GetState(), nil

	case "get_surroundings":
		return gameClient.SendCommand("get_surroundings", nil)

	case "move_to":
		x := int(params["x"].(float64))
		y := int(params["y"].(float64))
		return gameClient.SendCommand("move_to", map[string]interface{}{"x": x, "y": y})

	case "interact":
		return gameClient.SendCommand("interact", nil)

	case "use_tool":
		return gameClient.SendCommand("use_tool", nil)

	case "select_item":
		slot := int(params["slot"].(float64))
		return gameClient.SendCommand("select_item", map[string]interface{}{"slot": slot})

	case "switch_tool":
		tool := params["tool"].(string)
		return gameClient.SendCommand("switch_tool", map[string]interface{}{"tool": tool})

	case "face_direction":
		dir := int(params["direction"].(float64))
		return gameClient.SendCommand("face_direction", map[string]interface{}{"direction": dir})

	case "cheat_mode_enable":
		return gameClient.SendCommand("cheat_mode_enable", nil)

	case "cheat_warp":
		location := params["location"].(string)
		return gameClient.SendCommand("cheat_warp", map[string]interface{}{"location": location})

	case "cheat_set_money":
		amount := int(params["amount"].(float64))
		return gameClient.SendCommand("cheat_set_money", map[string]interface{}{"amount": amount})

	default:
		return nil, fmt.Errorf("unknown tool: %s", name)
	}
}
