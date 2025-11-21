package main

import (
	"encoding/json"
	"log"
	"net/http"
	"sync"
	"time"

	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true // Allow all origins for demo
	},
}

type Card struct {
	ID         string    `json:"id"`
	Name       string    `json:"name"`
	Type       string    `json:"type"`
	Power      string    `json:"power,omitempty"`
	Toughness  string    `json:"toughness,omitempty"`
	Zone       string    `json:"zone"`
	Tapped     bool      `json:"tapped"`
	Attacking  bool      `json:"attacking"`
	Blocking   bool      `json:"blocking"`
	Damage     int       `json:"damage"`
	Controller string    `json:"controller"`
	Owner      string    `json:"owner"`
	Abilities  []Ability `json:"abilities"`
}

type Ability struct {
	ID   string `json:"id"`
	Text string `json:"text"`
}

type Player struct {
	ID             string `json:"id"`
	Name           string `json:"name"`
	Life           int    `json:"life"`
	LibraryCount   int    `json:"library_count"`
	HandCount      int    `json:"hand_count"`
	GraveyardCount int    `json:"graveyard_count"`
}

type GameState struct {
	GameID         string   `json:"game_id"`
	CurrentPlayer  string   `json:"current_player"`
	ActivePlayer   string   `json:"active_player"`
	PriorityPlayer string   `json:"priority_player"`
	Phase          string   `json:"phase"`
	Step           string   `json:"step"`
	Turn           int      `json:"turn"`
	Players        []Player `json:"players"`
	Battlefield    []Card   `json:"battlefield"`
	Hand           []Card   `json:"hand"`
	Graveyard      []Card   `json:"graveyard"`
	Exile          []Card   `json:"exile"`
	Stack          []any    `json:"stack"`
}

type WSMessage struct {
	Type     string `json:"type"`
	GameID   string `json:"game_id,omitempty"`
	PlayerID string `json:"player_id,omitempty"`
	Data     any    `json:"data,omitempty"`
}

type Client struct {
	conn     *websocket.Conn
	send     chan []byte
	playerID string
	gameID   string
}

type Hub struct {
	clients    map[*Client]bool
	broadcast  chan []byte
	register   chan *Client
	unregister chan *Client
	mu         sync.RWMutex
	games      map[string]*GameState
}

func newHub() *Hub {
	return &Hub{
		clients:    make(map[*Client]bool),
		broadcast:  make(chan []byte),
		register:   make(chan *Client),
		unregister: make(chan *Client),
		games:      make(map[string]*GameState),
	}
}

func (h *Hub) run() {
	for {
		select {
		case client := <-h.register:
			h.clients[client] = true
			log.Printf("Client registered: %s", client.playerID)

		case client := <-h.unregister:
			if _, ok := h.clients[client]; ok {
				delete(h.clients, client)
				close(client.send)
				log.Printf("Client unregistered: %s", client.playerID)
			}

		case message := <-h.broadcast:
			for client := range h.clients {
				select {
				case client.send <- message:
				default:
					close(client.send)
					delete(h.clients, client)
				}
			}
		}
	}
}

func (h *Hub) createDemoGame(gameID string) *GameState {
	h.mu.Lock()
	defer h.mu.Unlock()

	game := &GameState{
		GameID:         gameID,
		CurrentPlayer:  "player1",
		ActivePlayer:   "player1",
		PriorityPlayer: "player1",
		Phase:          "Main",
		Step:           "Main1",
		Turn:           1,
		Players: []Player{
			{ID: "player1", Name: "Alice", Life: 20, LibraryCount: 53, HandCount: 7, GraveyardCount: 0},
			{ID: "player2", Name: "Bob", Life: 20, LibraryCount: 53, HandCount: 7, GraveyardCount: 0},
		},
		Battlefield: []Card{
			{
				ID: "card-1", Name: "Grizzly Bears", Type: "Creature - Bear",
				Power: "2", Toughness: "2", Zone: "battlefield",
				Controller: "player1", Owner: "player1",
				Abilities: []Ability{{ID: "none", Text: ""}},
			},
			{
				ID: "card-2", Name: "Serra Angel", Type: "Creature - Angel",
				Power: "4", Toughness: "4", Zone: "battlefield",
				Controller: "player1", Owner: "player1",
				Abilities: []Ability{
					{ID: "FlyingAbility", Text: "Flying"},
					{ID: "VigilanceAbility", Text: "Vigilance"},
				},
			},
			{
				ID: "card-3", Name: "Shivan Dragon", Type: "Creature - Dragon",
				Power: "5", Toughness: "5", Zone: "battlefield",
				Controller: "player2", Owner: "player2",
				Abilities: []Ability{
					{ID: "FlyingAbility", Text: "Flying"},
				},
			},
			{
				ID: "card-4", Name: "Llanowar Elves", Type: "Creature - Elf Druid",
				Power: "1", Toughness: "1", Zone: "battlefield",
				Controller: "player2", Owner: "player2", Tapped: true,
				Abilities: []Ability{{ID: "mana", Text: "{T}: Add {G}"}},
			},
		},
		Hand:      []Card{},
		Graveyard: []Card{},
		Exile:     []Card{},
		Stack:     []any{},
	}

	h.games[gameID] = game
	return game
}

func (h *Hub) handleMessage(client *Client, msg WSMessage) {
	log.Printf("Received message: %s from %s", msg.Type, client.playerID)

	switch msg.Type {
	case "create_game":
		gameID := "game-" + time.Now().Format("20060102-150405")
		game := h.createDemoGame(gameID)
		client.gameID = gameID

		response, _ := json.Marshal(WSMessage{
			Type: "game_state",
			Data: game,
		})
		client.send <- response

	case "join_game":
		h.mu.RLock()
		game, exists := h.games[msg.GameID]
		h.mu.RUnlock()

		if !exists {
			game = h.createDemoGame(msg.GameID)
		}

		client.gameID = msg.GameID
		client.playerID = msg.PlayerID

		response, _ := json.Marshal(WSMessage{
			Type: "game_state",
			Data: game,
		})
		client.send <- response

	case "declare_attacker":
		h.mu.Lock()
		game := h.games[client.gameID]
		if game != nil {
			data := msg.Data.(map[string]any)
			cardID := data["card_id"].(string)

			// Find and update card
			for i := range game.Battlefield {
				if game.Battlefield[i].ID == cardID {
					game.Battlefield[i].Attacking = true
					game.Battlefield[i].Tapped = true
					break
				}
			}
		}
		h.mu.Unlock()

		// Broadcast updated state
		h.broadcastGameState(client.gameID)

	case "pass_priority":
		h.mu.Lock()
		game := h.games[client.gameID]
		if game != nil {
			// Simple turn passing
			if game.CurrentPlayer == "player1" {
				game.CurrentPlayer = "player2"
			} else {
				game.CurrentPlayer = "player1"
				game.Turn++
			}
		}
		h.mu.Unlock()

		h.broadcastGameState(client.gameID)
	}
}

func (h *Hub) broadcastGameState(gameID string) {
	h.mu.RLock()
	game := h.games[gameID]
	h.mu.RUnlock()

	if game == nil {
		return
	}

	response, _ := json.Marshal(WSMessage{
		Type: "game_state",
		Data: game,
	})

	// Send to all clients in this game
	for client := range h.clients {
		if client.gameID == gameID {
			select {
			case client.send <- response:
			default:
			}
		}
	}
}

func (c *Client) readPump(hub *Hub) {
	defer func() {
		hub.unregister <- c
		c.conn.Close()
	}()

	for {
		_, message, err := c.conn.ReadMessage()
		if err != nil {
			break
		}

		var msg WSMessage
		if err := json.Unmarshal(message, &msg); err != nil {
			log.Printf("Error unmarshaling message: %v", err)
			continue
		}

		hub.handleMessage(c, msg)
	}
}

func (c *Client) writePump() {
	defer c.conn.Close()

	for message := range c.send {
		if err := c.conn.WriteMessage(websocket.TextMessage, message); err != nil {
			break
		}
	}
}

func serveWS(hub *Hub, w http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println(err)
		return
	}

	client := &Client{
		conn: conn,
		send: make(chan []byte, 256),
	}

	hub.register <- client

	go client.writePump()
	go client.readPump(hub)
}

func main() {
	hub := newHub()
	go hub.run()

	http.HandleFunc("/ws", func(w http.ResponseWriter, r *http.Request) {
		serveWS(hub, w, r)
	})

	log.Println("🚀 WebSocket server starting on :8080")
	log.Println("📡 WebSocket endpoint: ws://localhost:8080/ws")
	log.Println("🎮 Ready for Svelte client connections!")

	if err := http.ListenAndServe(":8080", nil); err != nil {
		log.Fatal("ListenAndServe: ", err)
	}
}
