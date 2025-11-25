package chat

import (
	"sync"
	"time"

	"github.com/magefree/mage-server-go/internal/session"
	pb "github.com/magefree/mage-server-go/pkg/proto/mage/v1"
	"go.uber.org/zap"
	"google.golang.org/protobuf/types/known/anypb"
	"google.golang.org/protobuf/types/known/timestamppb"
)

// Message represents a chat message
type Message struct {
	UserName  string
	Text      string
	Timestamp time.Time
	Color     string
	Type      string
}

// ChatRoom represents a chat room
type ChatRoom struct {
	ID          string
	messages    []Message
	users       map[string]bool // username -> present
	mu          sync.RWMutex
	maxMessages int
}

// NewChatRoom creates a new chat room
func NewChatRoom(id string, maxMessages int) *ChatRoom {
	return &ChatRoom{
		ID:          id,
		messages:    make([]Message, 0),
		users:       make(map[string]bool),
		maxMessages: maxMessages,
	}
}

// AddMessage adds a message to the chat room
func (c *ChatRoom) AddMessage(msg Message) {
	c.mu.Lock()
	defer c.mu.Unlock()

	c.messages = append(c.messages, msg)

	// Keep only last N messages
	if len(c.messages) > c.maxMessages {
		c.messages = c.messages[len(c.messages)-c.maxMessages:]
	}
}

// GetMessages returns recent messages
func (c *ChatRoom) GetMessages(limit int) []Message {
	c.mu.RLock()
	defer c.mu.RUnlock()

	if limit == 0 || limit > len(c.messages) {
		limit = len(c.messages)
	}

	start := len(c.messages) - limit
	if start < 0 {
		start = 0
	}

	return append([]Message{}, c.messages[start:]...)
}

// AddUser adds a user to the chat room
func (c *ChatRoom) AddUser(username string) {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.users[username] = true
}

// RemoveUser removes a user from the chat room
func (c *ChatRoom) RemoveUser(username string) {
	c.mu.Lock()
	defer c.mu.Unlock()
	delete(c.users, username)
}

// GetUsers returns list of users in the chat room
func (c *ChatRoom) GetUsers() []string {
	c.mu.RLock()
	defer c.mu.RUnlock()

	users := make([]string, 0, len(c.users))
	for user := range c.users {
		users = append(users, user)
	}
	return users
}

// Manager manages chat rooms
type Manager struct {
	rooms      map[string]*ChatRoom
	mu         sync.RWMutex
	logger     *zap.Logger
	sessionMgr session.Manager
}

// NewManager creates a new chat manager
func NewManager(logger *zap.Logger, sessionMgr session.Manager) *Manager {
	return &Manager{
		rooms:      make(map[string]*ChatRoom),
		logger:     logger,
		sessionMgr: sessionMgr,
	}
}

// CreateRoom creates a new chat room
func (m *Manager) CreateRoom(id string) *ChatRoom {
	m.mu.Lock()
	defer m.mu.Unlock()

	room := NewChatRoom(id, 100) // Keep last 100 messages
	m.rooms[id] = room

	m.logger.Debug("chat room created", zap.String("room_id", id))

	return room
}

// GetRoom retrieves a chat room by ID
func (m *Manager) GetRoom(id string) (*ChatRoom, bool) {
	m.mu.RLock()
	defer m.mu.RUnlock()

	room, ok := m.rooms[id]
	return room, ok
}

// GetOrCreateRoom gets or creates a chat room
func (m *Manager) GetOrCreateRoom(id string) *ChatRoom {
	room, ok := m.GetRoom(id)
	if !ok {
		return m.CreateRoom(id)
	}
	return room
}

// RemoveRoom removes a chat room
func (m *Manager) RemoveRoom(id string) {
	m.mu.Lock()
	defer m.mu.Unlock()

	delete(m.rooms, id)

	m.logger.Debug("chat room removed", zap.String("room_id", id))
}

// SendMessage sends a message to a chat room
func (m *Manager) SendMessage(roomID, username, text string) error {
	room, ok := m.GetRoom(roomID)
	if !ok {
		m.logger.Warn("chat room not found for SendMessage",
			zap.String("room_id", roomID),
			zap.String("username", username),
		)
		return nil // Room doesn't exist, ignore
	}

	msg := Message{
		UserName:  username,
		Text:      text,
		Timestamp: time.Now(),
		Color:     "BLACK",
		Type:      "TALK",
	}

	room.AddMessage(msg)

	m.logger.Info("chat message received",
		zap.String("room_id", roomID),
		zap.String("username", username),
		zap.String("message", text),
		zap.Int("users_in_room", len(room.GetUsers())),
	)

	// Broadcast message to all users in room via WebSocket
	m.broadcastMessage(roomID, msg)

	return nil
}

// broadcastMessage sends a chat message to all users in a room via WebSocket
func (m *Manager) broadcastMessage(roomID string, msg Message) {
	room, ok := m.GetRoom(roomID)
	if !ok {
		return
	}

	// Convert color string to enum
	var colorEnum pb.MessageColor
	switch msg.Color {
	case "BLACK":
		colorEnum = pb.MessageColor_BLACK
	case "GREEN":
		colorEnum = pb.MessageColor_GREEN
	case "BLUE":
		colorEnum = pb.MessageColor_BLUE
	case "RED":
		colorEnum = pb.MessageColor_RED
	case "YELLOW":
		colorEnum = pb.MessageColor_YELLOW
	case "ORANGE":
		colorEnum = pb.MessageColor_ORANGE
	default:
		colorEnum = pb.MessageColor_BLACK
	}

	// Convert type string to enum
	var typeEnum pb.MessageType
	switch msg.Type {
	case "TALK":
		typeEnum = pb.MessageType_TALK
	case "WHISPER":
		typeEnum = pb.MessageType_WHISPER
	case "STATUS":
		typeEnum = pb.MessageType_STATUS
	case "USER_INFO":
		typeEnum = pb.MessageType_USER_INFO
	default:
		typeEnum = pb.MessageType_TALK
	}

	// Create protobuf message
	protoMsg := &pb.ChatMessage{
		UserName:    msg.UserName,
		Message:     msg.Text,
		Time:        timestamppb.New(msg.Timestamp),
		Color:       colorEnum,
		MessageType: typeEnum,
	}

	// Create ChatMessageData
	chatData := &pb.ChatMessageData{
		Message: protoMsg,
		ChatId:  roomID,
	}

	// Pack into Any
	anyData, err := anypb.New(chatData)
	if err != nil {
		m.logger.Error("failed to pack chat message data",
			zap.Error(err),
			zap.String("room_id", roomID),
		)
		return
	}

	// Get all users in the room
	users := room.GetUsers()

	// Send to each user's session (users can have multiple sessions)
	sentCount := 0
	for _, username := range users {
		// Find all sessions for this user
		sessions := m.sessionMgr.GetSessionsByUser(username)

		for _, sess := range sessions {
			if !sess.IsConnected() {
				continue
			}

			// Create ServerEvent
			event := &pb.ServerEvent{
				SessionId: sess.ID,
				ObjectId:  roomID,
				Method:    pb.CallbackMethod_CHATMESSAGE,
				Data:      anyData,
				MessageId: 0, // Not used for now
			}

			// Send callback
			sent := sess.SendCallback(event)
			if sent {
				sentCount++
			} else {
				m.logger.Warn("failed to send chat message to session",
					zap.String("room_id", roomID),
					zap.String("username", username),
					zap.String("session_id", sess.ID),
				)
			}
		}
	}

	m.logger.Debug("broadcasted chat message",
		zap.String("room_id", roomID),
		zap.String("username", msg.UserName),
		zap.Int("users", len(users)),
		zap.Int("sessions_sent", sentCount),
	)
}

// JoinRoom adds a user to a chat room
func (m *Manager) JoinRoom(roomID, username string) {
	room := m.GetOrCreateRoom(roomID)
	room.AddUser(username)

	m.logger.Info("user joined chat room",
		zap.String("room_id", roomID),
		zap.String("username", username),
		zap.Int("total_users", len(room.GetUsers())),
		zap.Strings("users", room.GetUsers()),
	)
}

// LeaveRoom removes a user from a chat room
func (m *Manager) LeaveRoom(roomID, username string) {
	room, ok := m.GetRoom(roomID)
	if !ok {
		return
	}

	room.RemoveUser(username)

	m.logger.Debug("user left chat room",
		zap.String("room_id", roomID),
		zap.String("username", username),
	)
}
