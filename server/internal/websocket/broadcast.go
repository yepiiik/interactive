package websocket

import (
	"encoding/json"
	"log"
)

// BroadcastToRoom sends a message to all clients in a specific room
func BroadcastToRoom(roomID string, messageType string, payload interface{}) {
	roomsMu.RLock()
	room, exists := rooms[roomID]
	roomsMu.RUnlock()

	if !exists {
		log.Printf("Room %s not found for broadcasting", roomID)
		return
	}

	message := Message{
		Type: messageType,
	}

	payloadBytes, err := json.Marshal(payload)
	if err != nil {
		log.Printf("Error marshaling payload: %v", err)
		return
	}
	message.Payload = payloadBytes

	msgBytes, err := json.Marshal(message)
	if err != nil {
		log.Printf("Error marshaling message: %v", err)
		return
	}

	room.mu.RLock()
	defer room.mu.RUnlock()

	for _, client := range room.Clients {
		select {
		case client.Send <- msgBytes:
		default:
			close(client.Send)
			delete(room.Clients, client.ID)
			client.Conn.Close()
		}
	}
} 