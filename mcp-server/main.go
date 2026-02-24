package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"net/http"
	"sync"
	"time"

	"github.com/gorilla/websocket"
)

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
	autoFlag := flag.Bool("auto", true, "Start in autonomous mode")
	goalFlag := flag.String("goal", `USE CHEAT MODE to setup the farm:
1. cheat_mode_enable first
3. cheat_clear_debris, cheat_cut_trees, cheat_mine_rocks
4. cheat_hoe_all to till soil
5. cheat_plant_seeds season appropriate seeds"
6. cheat_grow_crops then cheat_harvest_all`, "Goal for autonomous mode")
	urlFlag := flag.String("url", "ws://localhost:8765/game", "WebSocket URL for the game mod")

	// Server mode flags for remote agent connections
	serverMode := flag.Bool("server", false, "Run as server to accept remote agent connections")
	hostFlag := flag.String("host", "127.0.0.1", "Host to bind to for remote connections")
	portFlag := flag.Int("port", 8765, "Port to listen on for remote connections")

	// OpenClaw Gateway mode
	openclawMode := flag.Bool("openclaw", false, "Connect to OpenClaw Gateway as tool provider")
	openclawURL := flag.String("openclaw-url", "ws://127.0.0.1:18789", "OpenClaw Gateway URL")
	openclawToken := flag.String("openclaw-token", "", "OpenClaw Gateway token (optional)")

	flag.Parse()

	gameClient = NewGameClient()

	// If OpenClaw Gateway mode
	if *openclawMode {
		runOpenClawGatewayMode(*openclawURL, *urlFlag, *openclawToken, *autoFlag, *goalFlag)
		return
	}

	// If server mode, run as remote agent server
	if *serverMode {
		runServerMode(*hostFlag, *portFlag, *urlFlag)
		return
	}

	// Original behavior - connect to game and optionally run agent
	go func() {
		for {
			if err := gameClient.Connect(*urlFlag); err != nil {
				log.Printf("Failed to connect to game (will retry): %v", err)
				time.Sleep(5 * time.Second)
				continue
			}
			log.Println("Connected to Stardew Valley!")

			if *autoFlag {
				log.Printf("Starting autonomous agent with goal: %s", *goalFlag)

				agent, err := NewStardewAgent()
				if err != nil {
					log.Printf("Failed to start agent: %v", err)
					return
				}
				if err := agent.StartSession(*goalFlag); err != nil {
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

// ============================================================================
// OpenClaw Gateway Protocol Implementation
// ============================================================================

// OpenClaw Gateway message types
type OpenClawRequest struct {
	Type   string                 `json:"type"`
	ID     string                 `json:"id"`
	Method string                 `json:"method"`
	Params map[string]interface{} `json:"params,omitempty"`
}

type OpenClawResponse struct {
	Type    string                 `json:"type"`
	ID      string                 `json:"id"`
	OK      bool                   `json:"ok"`
	Payload map[string]interface{} `json:"payload,omitempty"`
	Error   map[string]interface{} `json:"error,omitempty"`
}

type OpenClawEvent struct {
	Type        string                 `json:"type"`
	Event       string                 `json:"event"`
	Payload     map[string]interface{} `json:"payload,omitempty"`
	Seq         int                    `json:"seq,omitempty"`
	StateVersion int                   `json:"stateVersion,omitempty"`
}

// OpenClaw Gateway connection
func connectToOpenClawGateway(gatewayURL string, token string) (*websocket.Conn, error) {
	log.Printf("Connecting to OpenClaw Gateway at %s...", gatewayURL)

	// Set up header for token authentication
	header := http.Header{}
	if token != "" {
		header.Set("Authorization", "Bearer "+token)
	}

	conn, _, err := websocket.DefaultDialer.Dial(gatewayURL, header)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to OpenClaw Gateway: %w", err)
	}

	// Send connect request
	connectReq := OpenClawRequest{
		Type:   "req",
		ID:     "connect",
		Method: "connect",
		Params: map[string]interface{}{
			"caps": []string{"tools.call", "tools.catalog", "operator.read"},
			"name": "stardew-mcp",
		},
	}

	if err := conn.WriteJSON(connectReq); err != nil {
		conn.Close()
		return nil, fmt.Errorf("failed to send connect request: %w", err)
	}

	// Wait for response
	_, msg, err := conn.ReadMessage()
	if err != nil {
		conn.Close()
		return nil, fmt.Errorf("failed to read connect response: %w", err)
	}

	var resp OpenClawResponse
	if err := json.Unmarshal(msg, &resp); err != nil {
		conn.Close()
		return nil, fmt.Errorf("failed to parse connect response: %w", err)
	}

	if !resp.OK {
		conn.Close()
		return nil, fmt.Errorf("connection rejected: %v", resp.Error)
	}

	log.Println("Connected to OpenClaw Gateway!")
	return conn, nil
}

// Register tools with OpenClaw Gateway
func registerToolsWithGateway(conn *websocket.Conn) error {
	tools := getStardewToolsForGateway()

	// Use tools.register method to register tools
	req := OpenClawRequest{
		Type:   "req",
		ID:     "register-tools",
		Method: "tools.register",
		Params: map[string]interface{}{
			"tools": tools,
		},
	}

	if err := conn.WriteJSON(req); err != nil {
		return fmt.Errorf("failed to register tools: %w", err)
	}

	// Wait for response
	_, msg, err := conn.ReadMessage()
	if err != nil {
		return fmt.Errorf("failed to read register response: %w", err)
	}

	var resp OpenClawResponse
	if err := json.Unmarshal(msg, &resp); err != nil {
		return fmt.Errorf("failed to parse register response: %w", err)
	}

	if !resp.OK {
		return fmt.Errorf("tool registration failed: %v", resp.Error)
	}

	log.Printf("Registered %d tools with OpenClaw Gateway", len(tools))
	return nil
}

// Run in OpenClaw Gateway mode - connects to Gateway as a tool provider
func runOpenClawGatewayMode(gatewayURL string, gameURL string, token string, autoStart bool, goal string) {
	// First connect to the game
	log.Printf("Connecting to Stardew Valley at %s...", gameURL)
	for {
		if err := gameClient.Connect(gameURL); err != nil {
			log.Printf("Failed to connect to game (will retry): %v", err)
			time.Sleep(5 * time.Second)
			continue
		}
		log.Println("Connected to Stardew Valley!")
		break
	}

	// Connect to OpenClaw Gateway
	conn, err := connectToOpenClawGateway(gatewayURL, token)
	if err != nil {
		log.Printf("Failed to connect to OpenClaw Gateway: %v", err)
		log.Println("Falling back to standalone mode...")
		if autoStart {
			startAutonomousAgent(goal)
		}
		return
	}
	defer conn.Close()

	// Register tools
	if err := registerToolsWithGateway(conn); err != nil {
		log.Printf("Failed to register tools: %v", err)
	}

	// Start autonomous agent if enabled
	if autoStart {
		startAutonomousAgent(goal)
	}

	// Handle messages from Gateway
	for {
		_, msg, err := conn.ReadMessage()
		if err != nil {
			log.Printf("Gateway read error: %v", err)
			break
		}

		var req OpenClawRequest
		if err := json.Unmarshal(msg, &req); err != nil {
			continue
		}

		// Handle tool calls
		if req.Type == "req" && req.Method == "tools.call" {
			go handleToolCall(conn, req)
		}
	}
}

// Handle tool call from OpenClaw Gateway
func handleToolCall(conn *websocket.Conn, req OpenClawRequest) {
	toolName, ok := req.Params["name"].(string)
	if !ok {
		sendErrorResponse(conn, req.ID, "missing tool name")
		return
	}

	params, _ := req.Params["params"].(map[string]interface{})

	result, err := executeOpenClawTool(toolName, params)

	resp := OpenClawResponse{
		Type: "res",
		ID:   req.ID,
		OK:   err == nil,
	}

	if err != nil {
		resp.Error = map[string]interface{}{
			"code":    "tool_error",
			"message": err.Error(),
		}
	} else {
		resp.Payload = map[string]interface{}{
			"result": result,
		}
	}

	conn.WriteJSON(resp)
}

// Send error response
func sendErrorResponse(conn *websocket.Conn, id string, message string) {
	resp := OpenClawResponse{
		Type: "res",
		ID:   id,
		OK:   false,
		Error: map[string]interface{}{
			"code":    "error",
			"message": message,
		},
	}
	conn.WriteJSON(resp)
}

// Execute tool and return result
func executeOpenClawTool(name string, params map[string]interface{}) (interface{}, error) {
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

// getStardewToolsForGateway returns tool definitions for OpenClaw Gateway
func getStardewToolsForGateway() []map[string]interface{} {
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
			"inputSchema": map[string]interface{}{
				"type": "object",
				"properties": map[string]interface{}{
					"x": map[string]interface{}{"type": "integer"},
					"y": map[string]interface{}{"type": "integer"},
				},
				"required": []string{"x", "y"},
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
			"inputSchema": map[string]interface{}{
				"type": "object",
				"properties": map[string]interface{}{
					"slot": map[string]interface{}{"type": "integer"},
				},
				"required": []string{"slot"},
			},
		},
		{
			"name":        "switch_tool",
			"description": "Switch to tool by name",
			"inputSchema": map[string]interface{}{
				"type": "object",
				"properties": map[string]interface{}{
					"tool": map[string]interface{}{"type": "string"},
				},
				"required": []string{"tool"},
			},
		},
		{
			"name":        "face_direction",
			"description": "Face a direction",
			"inputSchema": map[string]interface{}{
				"type": "object",
				"properties": map[string]interface{}{
					"direction": map[string]interface{}{"type": "integer", "description": "0=down, 1=left, 2=right, 3=up"},
				},
				"required": []string{"direction"},
			},
		},
		{
			"name":        "cheat_mode_enable",
			"description": "Enable god-mode cheat commands",
		},
		{
			"name":        "cheat_warp",
			"description": "Teleport to location",
			"inputSchema": map[string]interface{}{
				"type": "object",
				"properties": map[string]interface{}{
					"location": map[string]interface{}{"type": "string"},
				},
				"required": []string{"location"},
			},
		},
		{
			"name":        "cheat_set_money",
			"description": "Set money amount",
			"inputSchema": map[string]interface{}{
				"type": "object",
				"properties": map[string]interface{}{
					"amount": map[string]interface{}{"type": "integer"},
				},
				"required": []string{"amount"},
			},
		},
	}
}

// runServerMode runs the MCP server that accepts remote agent connections
func runServerMode(host string, port int, gameURL string) {
	addr := fmt.Sprintf("%s:%d", host, port)

	// First connect to the game
	log.Printf("Connecting to Stardew Valley at %s...", gameURL)
	for {
		if err := gameClient.Connect(gameURL); err != nil {
			log.Printf("Failed to connect to game (will retry): %v", err)
			time.Sleep(5 * time.Second)
			continue
		}
		log.Println("Connected to Stardew Valley!")
		break
	}

	// Set up WebSocket upgrader
	upgrader := websocket.Upgrader{
		CheckOrigin: func(r *http.Request) bool { return true },
	}

	// HTTP server for WebSocket connections
	http.HandleFunc("/mcp", func(w http.ResponseWriter, r *http.Request) {
		conn, err := upgrader.Upgrade(w, r, nil)
		if err != nil {
			log.Printf("WebSocket upgrade failed: %v", err)
			return
		}
		defer conn.Close()

		log.Printf("Remote agent connected from %s", r.RemoteAddr)

		// Handle messages from remote agent
		for {
			_, msg, err := conn.ReadMessage()
			if err != nil {
				log.Printf("Remote agent disconnected: %v", err)
				break
			}

			var req WebSocketMessage
			if err := json.Unmarshal(msg, &req); err != nil {
				log.Printf("Failed to parse message: %v", err)
				continue
			}

			// Process command and send to game
			if req.Type == "command" {
				resp, err := gameClient.SendCommand(req.Action, req.Params)

				// Send response back to agent
				response := map[string]interface{}{
					"id":      req.ID,
					"type":    "response",
					"success": err == nil,
				}
				if err != nil {
					response["error"] = err.Error()
				} else if resp != nil {
					response["success"] = resp.Success
					response["message"] = resp.Message
					response["data"] = resp.Data
				}

				conn.WriteJSON(response)
			} else if req.Type == "get_state" {
				// Return current game state
				state := gameClient.GetState()
				response := map[string]interface{}{
					"id":   req.ID,
					"type": "state",
					"data": state,
				}
				conn.WriteJSON(response)
			} else if req.Type == "ping" {
				response := map[string]interface{}{
					"id":   req.ID,
					"type": "pong",
				}
				conn.WriteJSON(response)
			}
		}
	})

	// Also handle root path
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.Write([]byte(`{"status": "ok", "service": "stardew-mcp-remote"}`))
	})

	log.Printf("========================================")
	log.Printf("Stardew MCP Server - Remote Mode")
	log.Printf("========================================")
	log.Printf("Listening for remote agents on: ws://%s/mcp", addr)
	log.Printf("Game connected at: %s", gameURL)
	log.Printf("========================================")
	log.Printf("Waiting for remote connections...")
	log.Printf("(Press Ctrl+C to stop)")
	log.Printf("========================================")

	if err := http.ListenAndServe(addr, nil); err != nil {
		log.Printf("HTTP server error: %v", err)
	}
}
