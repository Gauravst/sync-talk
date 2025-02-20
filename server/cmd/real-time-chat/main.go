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
	router.HandleFunc("GET /api/users", handlers.GetAllUsers(userService))
	router.HandleFunc("GET /api/user", handlers.GetUser(userService))
	router.HandleFunc("GET /api/user/{id}", handlers.GetUserById(userService))
	router.HandleFunc("PUT /api/user/{id}", handlers.UpdateUser(userService))
	router.HandleFunc("DELETE /api/user/{id}", handlers.DeleteUser(userService))
	router.HandleFunc("POST /api/auth/login", handlers.LoginUser(authService, *cfg))
	router.HandleFunc("GET /api/room", handlers.GetAllChatRoom(chatService))
	router.HandleFunc("GET /api/room/{name}", handlers.GetChatRoomByName(chatService))
	router.HandleFunc("POST /api/room", handlers.CreateNewChatRoom(chatService))
	router.HandleFunc("PUT /api/room/{name}", handlers.UpdateChatRoom(chatService))
	router.HandleFunc("DELETE /api/room/{name}", handlers.DeleteChatRoom(chatService))

	// join room
	router.HandleFunc("POST /api/join/{name}", handlers.JoinRoom(chatService))

	// get all joined room by user
	router.HandleFunc("GET /api/join", handlers.GetAllJoinRoom(chatService))

	// WebSocket route
	// working------ > add here auth middleware Ware
	router.HandleFunc("/chat/{roomName}", handlers.LiveChat(chatService, *cfg))

	// Wrap the router with CORS middleware
	authHandler := middleware.Auth(cfg, authService)(router)
	corsHandler := middleware.CORS(cfg)(authHandler)

	// setup server
	server := &http.Server{
		Addr:    cfg.Address,
		Handler: corsHandler,
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
