package main

import (
	"context"
	"log"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gauravst/real-time-chat/internal/api/handlers"
	"github.com/gauravst/real-time-chat/internal/api/middleware"
	"github.com/gauravst/real-time-chat/internal/config"
	"github.com/gauravst/real-time-chat/internal/database"
	"github.com/gauravst/real-time-chat/internal/repositories"
	"github.com/gauravst/real-time-chat/internal/services"
)

func main() {
	// load config
	cfg := config.ConfigMustLoad()

	// database setup
	database.InitDB(cfg.DatabaseUri)
	defer database.CloseDB()

	// setup router
	router := http.NewServeMux()

	// Initialize repositories and services
	userRepo := repositories.NewUserRepository(database.DB)
	userService := services.NewUserService(userRepo)

	authRepo := repositories.NewAuthRepository(database.DB)
	authService := services.NewAuthService(authRepo)

	chatRepo := repositories.NewChatRepository(database.DB)
	chatService := services.NewChatService(chatRepo)

	// REST API routes
	router.Handle("GET /api/users",
		middleware.Auth(cfg)(
			http.HandlerFunc(handlers.GetAllUsers(userService)),
		),
	)

	router.Handle("GET /api/user",
		middleware.Auth(cfg)(
			http.HandlerFunc(handlers.GetUser(userService)),
		),
	)

	router.Handle("GET /api/user/{id}",
		middleware.Auth(cfg)(
			http.HandlerFunc(handlers.GetUserById(userService)),
		),
	)

	router.Handle("PUT /api/user/{id}",
		middleware.Auth(cfg)(
			http.HandlerFunc(handlers.UpdateUser(userService)),
		),
	)

	router.Handle("DELETE /api/user/{id}",
		middleware.Auth(cfg)(
			http.HandlerFunc(handlers.DeleteUser(userService)),
		),
	)

	router.HandleFunc("POST /api/auth/login", handlers.LoginUser(authService, *cfg))

	router.Handle("GET /api/room",
		middleware.Auth(cfg)(
			http.HandlerFunc(handlers.GetAllChatRoom(chatService)),
		),
	)

	router.Handle("GET /api/room/{name}",
		middleware.Auth(cfg)(
			http.HandlerFunc(handlers.GetChatRoomByName(chatService)),
		),
	)

	router.Handle("POST /api/room",
		middleware.Auth(cfg)(
			http.HandlerFunc(handlers.CreateNewChatRoom(chatService)),
		),
	)

	router.Handle("PUT /api/room/{name}",
		middleware.Auth(cfg)(
			http.HandlerFunc(handlers.UpdateChatRoom(chatService)),
		),
	)

	router.Handle("DELETE /api/room/{name}",
		middleware.Auth(cfg)(
			http.HandlerFunc(handlers.DeleteChatRoom(chatService)),
		),
	)

	// WebSocket route
	router.HandleFunc("/chat/{roomName}", handlers.LiveChat(chatService, *cfg))

	// setup server
	server := &http.Server{
		Addr:    cfg.Address,
		Handler: router,
	}

	slog.Info("server started", slog.String("address", cfg.Address))

	done := make(chan os.Signal, 1)
	signal.Notify(done, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		err := server.ListenAndServe()
		if err != nil && err != http.ErrServerClosed {
			log.Fatal("failed to start server", slog.String("error", err.Error()))
		}
	}()

	<-done

	slog.Info("shutting down the server")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	err := server.Shutdown(ctx)
	if err != nil {
		slog.Error("failed to shutdown server", slog.String("error", err.Error()))
	}

	slog.Info("server shutdown successfully")
}
