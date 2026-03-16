package main

import (
	"context"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/parmeet20/golang-chatapp/internal/auth"
	"github.com/parmeet20/golang-chatapp/internal/config"
	"github.com/parmeet20/golang-chatapp/internal/database"
	"github.com/parmeet20/golang-chatapp/internal/message"
	"github.com/parmeet20/golang-chatapp/internal/room"
	"github.com/parmeet20/golang-chatapp/internal/routes"
	"github.com/parmeet20/golang-chatapp/internal/user"
	"github.com/parmeet20/golang-chatapp/internal/websocket"
)

func main() {
	cfg := config.LoadConfig()

	db, err := database.NewDB(cfg.MONGO_URL, "chat-app")
	if err != nil {
		slog.Error("failed to connect to database", "error", err)
		os.Exit(1)
	}

	authService := auth.NewAuthService(cfg.JWT_SECRET, cfg.JWT_TOKEN_EXPIRATION_TIME)

	userRepo := user.NewUserRepo(db.Database)
	userService := user.NewUserService(userRepo, authService)
	userController := user.NewUserController(userService, authService)

	roomRepo := room.NewRoomRepo(db.Database)
	roomService := room.NewRoomService(roomRepo)
	roomController := room.NewRoomController(roomService)

	messageRepo := message.NewMessageRepo(db.Database)
	messageService := message.NewMessageService(messageRepo)
	messageController := message.NewMessageController(messageService, roomService)

	hub := websocket.NewHub()
	go hub.Run()

	router := routes.SetUpRouter(
		cfg,
		authService,
		userController,
		roomController,
		messageController,
		hub,
		roomService,
		messageService,
	)

	server := &http.Server{
		Addr:    ":" + cfg.PORT,
		Handler: router,
	}

	go func() {
		slog.Info("server started", "port", cfg.PORT)
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			slog.Error("server failed", "error", err)
			os.Exit(1)
		}
	}()

	sig := make(chan os.Signal, 1)
	signal.Notify(sig, os.Interrupt, syscall.SIGTERM)

	<-sig

	slog.Info("shutdown signal received")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := server.Shutdown(ctx); err != nil {
		slog.Error("graceful shutdown failed", "error", err)
	}

	hub.Stop()

	slog.Info("server stopped")
}
