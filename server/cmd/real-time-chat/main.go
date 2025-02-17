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
	"github.com/gauravst/real-time-chat/internal/config"
	"github.com/gauravst/real-time-chat/internal/database"
	"github.com/gauravst/real-time-chat/internal/repositories"
	"github.com/gauravst/real-time-chat/internal/services"
	"github.com/gorilla/websocket"
)

// upgrader to upgrade HTTP connection to Websocket
var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

func main() {
	// load config
	cfg := config.ConfigMustLoad()

	// database setup
	database.InitDB(cfg.DatabaseUri)
	defer database.CloseDB()

	// setup router
	router := http.NewServeMux()

	// Initialize repositories and services
	// userRepo := repositories.NewUserRepository(database.DB)
	// userService := services.NewUserService(userRepo)

	authRepo := repositories.NewAuthRepository(database.DB)
	authService := services.NewAuthService(authRepo)

	// REST API routes
	router.HandleFunc("POST /api/user", handlers.LoginUser(authService, *cfg))

	// WebSocket route
	router.HandleFunc("/chat/{id}", handlers.LiveChat(*cfg, upgrader))

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
