package websocket

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/parmeet20/golang-chatapp/internal/auth"
	"github.com/parmeet20/golang-chatapp/internal/message"
	"github.com/parmeet20/golang-chatapp/internal/room"
)

func RegisterRoutes(r chi.Router, authService *auth.AuthService, hub *Hub, roomService *room.RoomService, messageService *message.MessageService) {
	r.Get("/ws/{roomId}", func(w http.ResponseWriter, r *http.Request) {

		roomId := chi.URLParam(r, "roomId")
		if roomId == "" {
			http.Error(w, "Room ID is required", http.StatusBadRequest)
			return
		}

		room := hub.GetOrCreateRoom(roomId, roomService, messageService)
		if room == nil {
			http.Error(w, "Failed to get or create room", http.StatusInternalServerError)
			return
		}

		ServeWs(authService, room, w, r)
	})
}
