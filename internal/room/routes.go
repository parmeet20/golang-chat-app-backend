package room

import (
	"github.com/go-chi/chi/v5"
	"github.com/parmeet20/golang-chatapp/internal/auth"
)

func RegisterRoutes(r chi.Router, controller *RoomController, authService *auth.AuthService) chi.Router {
	r.Route("/rooms", func(r chi.Router) {
		r.Use(authService.JwtMiddleware)
		r.Get("/", controller.FindAllRooms)
		r.Post("/", controller.CreateRoom)
		r.Get("/{id}", controller.GetRoomByID)
		r.Put("/{id}/join", controller.JoinRoom)
		r.Put("/{id}/leave", controller.LeaveRoom)
		r.Get("/user/{userId}", controller.GetRoomsByUserID)
	})
	return r
}
