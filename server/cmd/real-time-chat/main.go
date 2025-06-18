package main

import (
	"context"
	"log"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"strconv"
	"sync"
	"syscall"
	"time"

	"github.com/gauravst/real-time-chat/internal/api/handlers"
	"github.com/gauravst/real-time-chat/internal/api/middleware"
	"github.com/gauravst/real-time-chat/internal/config"
	"github.com/gauravst/real-time-chat/internal/database"
	"github.com/gauravst/real-time-chat/internal/models"
	"github.com/gauravst/real-time-chat/internal/repositories"
	"github.com/gauravst/real-time-chat/internal/services"
	"github.com/gorilla/websocket"
)

func main() {
	// loading queries
	queryManager, err := database.NewQueryManager()
	if err != nil {
		log.Fatalf("Failed to load SQL queries: %v", err)
	}

	// load config
	cfg := config.ConfigMustLoad()

	// database setup
	database.InitDB(cfg.DatabaseUri)
	defer database.CloseDB()

	// Initialize repositories and services
	userRepo := repositories.NewUserRepository(database.DB)
	userService := services.NewUserService(userRepo)

	authRepo := repositories.NewAuthRepository(database.DB)
	authService := services.NewAuthService(authRepo)

	chatRepo := repositories.NewChatRepository(database.DB, queryManager)
	chatService := services.NewChatService(chatRepo)

	fileRepo := repositories.NewFileRepository(database.DB, queryManager)
	fileService := services.NewFileService(fileRepo, chatRepo)

	// Setup routers
	router := http.NewServeMux()
	publicRouter := http.NewServeMux()
	// publicRouter2 := http.NewServeMux()

	wsServer := &models.WsServer{
		RoomMutex:  &sync.Mutex{},
		Rooms:      make(map[string][]*websocket.Conn),
		OnlineUser: make(map[string]map[string]bool),
		Upgrader: websocket.Upgrader{
			CheckOrigin: func(r *http.Request) bool { return true },
		},
	}

	// health api
	publicRouter.HandleFunc("GET /api/auth/health", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("OK"))
	})

	// Public routes (No Auth)
	publicRouter.HandleFunc("POST /api/auth/login", handlers.LoginUser(authService, *cfg))
	publicRouter.HandleFunc("POST /api/auth/loginWithoutAuth", handlers.LoginWithoutAuth(authService, *cfg))

	// Protected routes (Require Auth)
	router.HandleFunc("GET /api/users", handlers.GetAllUsers(userService))
	router.HandleFunc("GET /api/user", handlers.GetUser(userService))
	router.HandleFunc("POST /api/user/logout", handlers.LogoutUser(authService))
	router.HandleFunc("GET /api/user/{id}", handlers.GetUserById(userService))
	router.HandleFunc("PUT /api/user/{id}", handlers.UpdateUser(userService))
	router.HandleFunc("DELETE /api/user/{id}", handlers.DeleteUser(userService))

	router.HandleFunc("GET /api/room", handlers.GetAllChatRoom(chatService))
	router.HandleFunc("GET /api/room/{name}", handlers.GetChatRoomByName(chatService))
	router.HandleFunc("GET /api/room/private/{code}", handlers.GetPrivateChatRoom(chatService))
	router.HandleFunc("POST /api/room", handlers.CreateNewChatRoom(chatService))
	router.HandleFunc("PUT /api/room/{name}", handlers.UpdateChatRoom(chatService))
	router.HandleFunc("DELETE /api/room/{name}", handlers.DeleteChatRoom(chatService))

	// Join room
	router.HandleFunc("GET /api/join", handlers.GetAllJoinRoom(chatService))
	router.HandleFunc("POST /api/join/{name}", handlers.JoinRoom(chatService))
	router.HandleFunc("POST /api/join/private/{code}", handlers.JoinPrivateRoom(chatService))
	router.HandleFunc("DELETE /api/join/{name}", handlers.LeaveRoom(chatService))

	// WebSocket route
	router.HandleFunc("/chat/{roomName}", handlers.LiveChat(chatService, *cfg, wsServer))

	// upload files
	router.HandleFunc("POST /api/chat/upload/{roomName}", handlers.UploadFileInRoom(fileService, *cfg, wsServer))
	// get old chats for a room
	router.HandleFunc("GET /api/chat/{roomName}/{limit}", handlers.GetOldChats(chatService))

	// Merge both routers
	mainRouter := http.NewServeMux()
	mainRouter.Handle("/api/auth/", publicRouter)                     // Public routes (No Auth)
	mainRouter.Handle("/", middleware.Auth(cfg, authService)(router)) // Protected routes
	// mainRouter.Handle("/chat/", publicRouter2)

	// Wrap everything with CORS middleware
	finalHandler := middleware.CORS(cfg)(mainRouter)

	// Setup server
	port := cfg.EnvPort
	addr := cfg.Address

	if port != 0 {
		addr = "0.0.0.0:" + strconv.Itoa(port) // Convert int to string
	}

	server := &http.Server{
		Addr:    addr,
		Handler: finalHandler,
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

	err = server.Shutdown(ctx)
	if err != nil {
		slog.Error("failed to shutdown server", slog.String("error", err.Error()))
	}

	slog.Info("server shutdown successfully")
}
