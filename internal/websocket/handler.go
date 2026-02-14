package websocket

import (
	"net/http"

	"gosaas/internal/realtime"

	"github.com/gorilla/websocket"
	"log/slog"
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin: func(r *http.Request) bool {
		// Allow all origins in development
		// TODO: Tighten this in production to check specific origins
		return true
	},
}

// Handler returns an HTTP handler function for WebSocket upgrades
func Handler(hub *realtime.Hub) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Extract user info from query parameters
		clientID := r.URL.Query().Get("clientId")
		userID := r.URL.Query().Get("userId")

		// Use default values if not provided
		if clientID == "" {
			clientID = "anonymous"
		}
		if userID == "" {
			userID = "anonymous"
		}

		slog.Info("Serving WebSocket", "clientID", clientID, "userID", userID)

		// Upgrade HTTP connection to WebSocket
		conn, err := upgrader.Upgrade(w, r, nil)
		if err != nil {
			slog.Error("WebSocket upgrade error", "error", err)
			return
		}

		// Delegate to the realtime hub
		realtime.ServeWS(hub, conn, clientID, userID)
	}
}
