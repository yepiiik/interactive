package websocket

import (
	"encoding/json"
	"log"
	"sync"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	"polling-app/internal/models"
	"polling-app/pkg/database"
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin: func(r *http.Request) bool {
		return true // Allow all origins in development
	},
}

type Client struct {
	ID     uint
	RoomID string
	Conn   *websocket.Conn
	Send   chan []byte
}

type Room struct {
	ID       string
	Clients  map[uint]*Client
	mu       sync.RWMutex
}

type Message struct {
	Type    string          `json:"type"`
	Payload json.RawMessage `json:"payload"`
}

var rooms = make(map[string]*Room)
var roomsMu sync.RWMutex

func HandleWebSocket(c *gin.Context) {
	roomID := c.Param("id")
	user, exists := c.Get("user")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}
	currentUser := user.(models.User)

	// Verify user is in the room
	var room models.Room
	if err := database.DB.First(&room, "id = ?", roomID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Room not found"})
		return
	}

	var isParticipant bool
	database.DB.Model(&room).Where("user_id = ?", currentUser.ID).Association("Participants").Count(&isParticipant)
	if !isParticipant {
		c.JSON(http.StatusForbidden, gin.H{"error": "Not a participant in this room"})
		return
	}

	conn, err := upgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		log.Printf("Failed to upgrade connection: %v", err)
		return
	}

	client := &Client{
		ID:     currentUser.ID,
		RoomID: roomID,
		Conn:   conn,
		Send:   make(chan []byte, 256),
	}

	// Get or create room
	roomMu.Lock()
	if _, exists := rooms[roomID]; !exists {
		rooms[roomID] = &Room{
			ID:      roomID,
			Clients: make(map[uint]*Client),
		}
	}
	room := rooms[roomID]
	room.mu.Lock()
	room.Clients[client.ID] = client
	room.mu.Unlock()
	roomMu.Unlock()

	// Start goroutines for reading and writing
	go client.writePump()
	go client.readPump(room)
}

func (c *Client) readPump(room *Room) {
	defer func() {
		room.mu.Lock()
		delete(room.Clients, c.ID)
		room.mu.Unlock()
		c.Conn.Close()
	}()

	for {
		_, message, err := c.Conn.ReadMessage()
		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
				log.Printf("error: %v", err)
			}
			break
		}

		var msg Message
		if err := json.Unmarshal(message, &msg); err != nil {
			log.Printf("error unmarshaling message: %v", err)
			continue
		}

		// Handle different message types
		switch msg.Type {
		case "vote":
			handleVote(c, room, msg.Payload)
		case "start_poll":
			handleStartPoll(c, room, msg.Payload)
		case "end_poll":
			handleEndPoll(c, room, msg.Payload)
		}
	}
}

func (c *Client) writePump() {
	defer func() {
		c.Conn.Close()
	}()

	for {
		select {
		case message, ok := <-c.Send:
			if !ok {
				c.Conn.WriteMessage(websocket.CloseMessage, []byte{})
				return
			}

			w, err := c.Conn.NextWriter(websocket.TextMessage)
			if err != nil {
				return
			}
			w.Write(message)

			n := len(c.Send)
			for i := 0; i < n; i++ {
				w.Write([]byte{'\n'})
				w.Write(<-c.Send)
			}

			if err := w.Close(); err != nil {
				return
			}
		}
	}
}

func handleVote(client *Client, room *Room, payload json.RawMessage) {
	// Broadcast vote to all clients in the room
	room.mu.RLock()
	defer room.mu.RUnlock()

	message := Message{
		Type:    "vote",
		Payload: payload,
	}

	msgBytes, err := json.Marshal(message)
	if err != nil {
		log.Printf("error marshaling vote message: %v", err)
		return
	}

	for _, c := range room.Clients {
		select {
		case c.Send <- msgBytes:
		default:
			close(c.Send)
			delete(room.Clients, c.ID)
			c.Conn.Close()
		}
	}
}

func handleStartPoll(client *Client, room *Room, payload json.RawMessage) {
	// Broadcast poll start to all clients
	room.mu.RLock()
	defer room.mu.RUnlock()

	message := Message{
		Type:    "start_poll",
		Payload: payload,
	}

	msgBytes, err := json.Marshal(message)
	if err != nil {
		log.Printf("error marshaling start poll message: %v", err)
		return
	}

	for _, c := range room.Clients {
		select {
		case c.Send <- msgBytes:
		default:
			close(c.Send)
			delete(room.Clients, c.ID)
			c.Conn.Close()
		}
	}
}

func handleEndPoll(client *Client, room *Room, payload json.RawMessage) {
	// Broadcast poll end to all clients
	room.mu.RLock()
	defer room.mu.RUnlock()

	message := Message{
		Type:    "end_poll",
		Payload: payload,
	}

	msgBytes, err := json.Marshal(message)
	if err != nil {
		log.Printf("error marshaling end poll message: %v", err)
		return
	}

	for _, c := range room.Clients {
		select {
		case c.Send <- msgBytes:
		default:
			close(c.Send)
			delete(room.Clients, c.ID)
			c.Conn.Close()
		}
	}
} 