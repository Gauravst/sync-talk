package handlers

import (
	"log/slog"
	"net/http"

	"github.com/gauravst/real-time-chat/internal/config"
	"github.com/gorilla/websocket"
)

func LiveChat(cfg config.Config, upgrader websocket.Upgrader) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		// Upgrade HTTP connection to WebSocket
		conn, err := upgrader.Upgrade(w, r, nil)
		if err != nil {
			slog.Error("failed to upgrade to WebSocket", slog.String("error", err.Error()))
			return
		}
		defer conn.Close()

		slog.Info("WebSocket connection established")

		// Handle WebSocket messages
		for {
			messageType, message, err := conn.ReadMessage()
			if err != nil {
				slog.Error("failed to read WebSocket message", slog.String("error", err.Error()))
				break
			}

			slog.Info("received WebSocket message", slog.String("message", string(message)))

			// Echo the message back to the client
			if err := conn.WriteMessage(messageType, message); err != nil {
				slog.Error("failed to write WebSocket message", slog.String("error", err.Error()))
				break
			}
		}
	}
}
