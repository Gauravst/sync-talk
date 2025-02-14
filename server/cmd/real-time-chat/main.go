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

	//setup router
	router := http.NewServeMux()

	userRepo := repositories.NewUserRepository(database.DB)
	userService := services.NewUserService(userRepo)

	router.HandleFunc("GET /api/user", middleware.Auth(handlers.GetUser(userService)))
	router.HandleFunc("POST /api/user", handlers.CreateUser(userService))
	router.HandleFunc("PUT /api/user", middleware.Auth(handlers.UpdateUser(userService)))
	router.HandleFunc("DELETE /api/user", middleware.Auth(handlers.DeleteUser(userService)))

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
		if err != nil {
			log.Fatal("failed to start server")
		}
	}()

	<-done

	slog.Info("shutting down the server")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	err := server.Shutdown(ctx)
	if err != nil {
		slog.Error("faild to Shutdown server", slog.String("error", err.Error()))
	}

	slog.Info("server Shutdown successfully")
}
