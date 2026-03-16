package message

import (
	"github.com/go-chi/chi/v5"
	"github.com/parmeet20/golang-chatapp/internal/auth"
)

func RegisterRoutes(r chi.Router, controller *MessageController, authService *auth.AuthService) chi.Router {
	r.Route("/messages", func(r chi.Router) {
		r.Use(authService.JwtMiddleware)
		r.Get("/{id}", controller.GetMessagesByRoomId)
	})
	return r
}
