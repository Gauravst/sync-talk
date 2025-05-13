package ws

import (
	"encoding/json"
	"log"

	"github.com/gauravst/real-time-chat/internal/models"
	"github.com/gorilla/websocket"
)

func BroadcastMessage(wsServer *models.WsServer, roomName string, sender *websocket.Conn, message *models.MessageResponse) {

	wsServer.RoomMutex.Lock()
	defer wsServer.RoomMutex.Unlock()

	clients := wsServer.Rooms[roomName]

	// Convert the message struct to JSON
	jsonMessage, err := json.Marshal(message)
	if err != nil {
		log.Println("Failed to marshal message:", err)
		return
	}

	// Send the JSON message to all clients **except the sender**
	for _, conn := range clients {
		if sender != nil && sender == conn {
			continue
		}

		if err := conn.WriteMessage(websocket.TextMessage, jsonMessage); err != nil {
			log.Println("Failed to send message:", err)
		}
	}
}
