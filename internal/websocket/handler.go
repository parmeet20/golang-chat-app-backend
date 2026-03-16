package websocket

import (
	"log"
	"net/http"
	"strings"

	"github.com/gorilla/websocket"
	"github.com/parmeet20/golang-chatapp/internal/auth"
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

func ServeWs(authService *auth.AuthService, room *Room, w http.ResponseWriter, r *http.Request) {

	if room == nil {
		http.Error(w, "Room not found", http.StatusBadRequest)
		return
	}

	var token string

	// 1️⃣ Try query param first (for browsers)
	token = r.URL.Query().Get("token")

	// 2️⃣ Fallback to Authorization header (for Postman / CLI)
	if token == "" {
		authHeader := r.Header.Get("Authorization")

		if authHeader != "" {
			parts := strings.SplitN(authHeader, " ", 2)

			if len(parts) == 2 && strings.ToLower(parts[0]) == "bearer" {
				token = parts[1]
			}
		}
	}

	if token == "" {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	claims, err := authService.VerifyToken(token)
	if err != nil {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println("websocket upgrade error:", err)
		return
	}

	client := NewClient(claims.UserID, conn, room)

	room.registerClient <- client

	go client.readPump()
	go client.writePump()
}
