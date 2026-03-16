package routes

import (
	"github.com/go-chi/chi/middleware"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/cors"
	"github.com/parmeet20/golang-chatapp/internal/auth"
	"github.com/parmeet20/golang-chatapp/internal/config"
	"github.com/parmeet20/golang-chatapp/internal/healthcheck"
	"github.com/parmeet20/golang-chatapp/internal/message"
	"github.com/parmeet20/golang-chatapp/internal/room"
	"github.com/parmeet20/golang-chatapp/internal/user"
	"github.com/parmeet20/golang-chatapp/internal/websocket"
)

func SetUpRouter(
	config *config.Config,
	authService *auth.AuthService,
	userController *user.UserController,
	roomController *room.RoomController,
	messageController *message.MessageController,
	hub *websocket.Hub,
	roomService *room.RoomService,
	messageService *message.MessageService,
) *chi.Mux {
	r := chi.NewRouter()

	r.Use(cors.Handler(cors.Options{
		AllowedOrigins:   config.AllowedOrigins,
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type"},
		ExposedHeaders:   []string{"Link"},
		AllowCredentials: true,
		MaxAge:           300,
	}))

	r.Use(middleware.RequestID)
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)

	r.Route("/api/v1", func(r chi.Router) {
		user.RegisterRoutes(r, userController)
		room.RegisterRoutes(r, roomController, authService)
		message.RegisterRoutes(r, messageController, authService)
		websocket.RegisterRoutes(r, authService, hub, roomService, messageService)
		r.Get("/health", healthcheck.HealthCheck)
	})

	return r
}
