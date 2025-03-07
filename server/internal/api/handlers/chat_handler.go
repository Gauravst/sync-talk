package handlers

import (
	"encoding/json"
	"fmt"
	"log"
	"log/slog"
	"net/http"
	"strconv"
	"sync"

	"github.com/gauravst/real-time-chat/internal/api/middleware"
	"github.com/gauravst/real-time-chat/internal/config"
	"github.com/gauravst/real-time-chat/internal/models"
	"github.com/gauravst/real-time-chat/internal/services"
	"github.com/gauravst/real-time-chat/internal/utils/response"
	"github.com/go-playground/validator/v10"
	"github.com/gorilla/websocket"
)

// upgrader to upgrade HTTP connection to Websocket
var (
	upgrader   = websocket.Upgrader{CheckOrigin: func(r *http.Request) bool { return true }}
	rooms      = make(map[string][]*websocket.Conn)
	onlineUser = make(map[string]int)
	roomMutex  sync.Mutex
)

func LiveChat(chatService services.ChatService, cfg config.Config) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// geting middleware data
		userDataRaw := r.Context().Value(middleware.UserDataKey)
		if userDataRaw == nil {
			http.Error(w, "unauthorized user", http.StatusUnauthorized)
			return
		}

		// Correct the type assertion to *models.AccessToken
		userData, ok := userDataRaw.(*models.AccessToken)
		if !ok {
			http.Error(w, "unauthorized user", http.StatusUnauthorized)
			return
		}

		// Store user data in a local variable
		currentUser := *userData

		roomName := r.PathValue("roomName")
		if roomName == "" {
			http.Error(w, "Missing room Name", http.StatusBadRequest)
			return
		}

		// Upgrade HTTP connection to WebSocket
		conn, err := upgrader.Upgrade(w, r, nil)
		if err != nil {
			slog.Error("failed to upgrade to WebSocket", slog.String("error", err.Error()))
			return
		}
		defer conn.Close()

		//check user join or not in room
		isMember, err := chatService.CheckChatRoomMember(currentUser.UserId, roomName)
		if err != nil {
			slog.Error(err.Error())
			conn.WriteMessage(websocket.TextMessage, []byte("Error: something went worng."))
			return
		}

		if !isMember {
			conn.WriteMessage(websocket.TextMessage, []byte("Error: You are not a member of this group."))
			return
		}

		// Add connection to the room
		roomMutex.Lock()
		rooms[roomName] = append(rooms[roomName], conn)
		onlineUser[roomName]++
		roomMutex.Unlock()

		// Broadcast the updated online user count
		go broadcastOnlineUsers(roomName)

		slog.Info("WebSocket connection established")

		// Handle WebSocket messages
		for {
			// geting message from client
			_, message, err := conn.ReadMessage()
			if err != nil {
				slog.Error("failed to read WebSocket message", slog.String("error", err.Error()))
				break
			}

			// decoding message here
			var msg models.MessageRequest
			err = json.Unmarshal(message, &msg)
			if err != nil {
				slog.Error("failed to parse WebSocket message", slog.String("error", err.Error()))
				continue
			}

			// save message in db here
			newMessageData := &models.MessageRequest{
				Type:    "chat",
				Content: msg.Content,
				UserId:  currentUser.UserId,
			}
			createdMessage, err := chatService.CreateNewMessage(newMessageData, roomName)
			if err != nil {
				slog.Error("failed to save message", slog.String("error", err.Error()))
				continue
			}

			// send message
			createdMessage.Type = "chat"
			go broadcastMessage(roomName, conn, createdMessage)
		}

		// remove connection
		removeConnection(roomName, conn)
	}
}

func broadcastMessage(roomName string, sender *websocket.Conn, message *models.MessageRequest) {
	roomMutex.Lock()
	defer roomMutex.Unlock()

	clients := rooms[roomName]

	// Convert the message struct to JSON
	jsonMessage, err := json.Marshal(message)
	if err != nil {
		log.Println("Failed to marshal message:", err)
		return
	}

	// Send the JSON message to all clients **except the sender**
	for _, conn := range clients {
		if conn == sender {
			continue
		}

		if err := conn.WriteMessage(websocket.TextMessage, jsonMessage); err != nil {
			log.Println("Failed to send message:", err)
		}
	}
}

func broadcastOnlineUsers(roomName string) {
	roomMutex.Lock()
	defer roomMutex.Unlock()

	clients := rooms[roomName]
	count := onlineUser[roomName]

	data := &models.OnlineUserCountRequest{
		Type:  "onlineUser",
		Count: count,
	}

	// Convert the message struct to JSON
	jsonMessage, err := json.Marshal(data)
	if err != nil {
		log.Println("Failed to marshal message:", err)
		return
	}

	for _, conn := range clients {
		if err := conn.WriteMessage(websocket.TextMessage, jsonMessage); err != nil {
			log.Println("Failed to send message:", err)
		}
	}
}

func removeConnection(roomName string, conn *websocket.Conn) {
	roomMutex.Lock()
	defer roomMutex.Unlock()

	clients := rooms[roomName]
	for i, c := range clients {
		if c == conn {
			rooms[roomName] = append(clients[:i], clients[i+1:]...)
			break
		}
	}

	// decrease the online user count
	if onlineUser[roomName] > 0 {
		onlineUser[roomName]--
	}

	// Broadcast updated online users count
	go broadcastOnlineUsers(roomName)
}

func GetAllChatRoom(chatService services.ChatService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		data, err := chatService.GetAllChatRoom()
		if err != nil {
			response.WriteJson(w, http.StatusInternalServerError, response.GeneralError(err))
			return
		}

		response.WriteJson(w, http.StatusOK, data)
		return
	}
}

func GetChatRoomByName(chatService services.ChatService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		name := r.PathValue("name")
		if name == " " {
			response.WriteJson(w, http.StatusNotFound, response.GeneralError(fmt.Errorf("name parms not found")))
			return
		}

		data, err := chatService.GetChatRoomByName(name)
		if err != nil {
			response.WriteJson(w, http.StatusInternalServerError, response.GeneralError(err))
			return
		}

		response.WriteJson(w, http.StatusOK, data)
		return
	}
}

func CreateNewChatRoom(chatService services.ChatService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Get value from context
		userDataRaw := r.Context().Value(middleware.UserDataKey)
		if userDataRaw == nil {
			response.WriteJson(w, http.StatusUnauthorized, response.GeneralError(fmt.Errorf("Unauthorized")))
			return
		}

		// Correct the type assertion to *models.AccessToken
		userData, ok := userDataRaw.(*models.AccessToken)
		if !ok {
			response.WriteJson(w, http.StatusUnauthorized, response.GeneralError(fmt.Errorf("Unauthorized")))
			return
		}

		// get chat data from Request
		var data models.ChatRoomRequest
		err := json.NewDecoder(r.Body).Decode(&data)
		if err != nil {
			response.WriteJson(w, http.StatusNotFound, response.GeneralError(err))
			return
		}

		// update data
		data.UserId = userData.UserId

		// vaildate data here
		err = validator.New().Struct(data)
		if err != nil {
			validateErrs := err.(validator.ValidationErrors)
			response.WriteJson(w, http.StatusBadRequest, response.ValidationError(validateErrs))
			return
		}

		// create new chat
		err = chatService.CreateNewChatRoom(&data)
		if err != nil {
			response.WriteJson(w, http.StatusInternalServerError, response.GeneralError(err))
			return
		}

		response.WriteJson(w, http.StatusOK, data)
		return
	}
}

func UpdateChatRoom(chatService services.ChatService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Get value from context
		userDataRaw := r.Context().Value(middleware.UserDataKey)
		if userDataRaw == nil {
			response.WriteJson(w, http.StatusUnauthorized, response.GeneralError(fmt.Errorf("Unauthorized")))
			return
		}

		// Correct the type assertion to *models.AccessToken
		userData, ok := userDataRaw.(*models.AccessToken)
		if !ok {
			response.WriteJson(w, http.StatusUnauthorized, response.GeneralError(fmt.Errorf("Unauthorized")))
			return
		}

		// get chat data from Request
		var data models.ChatRoomRequest
		err := json.NewDecoder(r.Body).Decode(&data)
		if err != nil {
			response.WriteJson(w, http.StatusNotFound, response.GeneralError(err))
			return
		}

		// get data from parms
		name := r.PathValue("name")
		if name == " " {
			response.WriteJson(w, http.StatusNotFound, response.GeneralError(fmt.Errorf("name parms not found")))
			return
		}

		roomData, err := chatService.GetChatRoomByName(name)
		if err != nil {
			response.WriteJson(w, http.StatusInternalServerError, response.GeneralError(err))
			return
		}

		if userData.Role != "ADMIN" && roomData.UserId != userData.UserId {
			response.WriteJson(w, http.StatusUnauthorized, response.GeneralError(fmt.Errorf("unauthorized user")))
			return
		}

		err = chatService.UpdateChatRoom(&data)
		if err != nil {
			response.WriteJson(w, http.StatusInternalServerError, response.GeneralError(err))
			return
		}

		response.WriteJson(w, http.StatusOK, data)
		return
	}
}

func DeleteChatRoom(chatService services.ChatService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Get value from context
		userDataRaw := r.Context().Value(middleware.UserDataKey)
		if userDataRaw == nil {
			response.WriteJson(w, http.StatusUnauthorized, response.GeneralError(fmt.Errorf("Unauthorized")))
			return
		}

		// Correct the type assertion to *models.AccessToken
		userData, ok := userDataRaw.(*models.AccessToken)
		if !ok {
			response.WriteJson(w, http.StatusUnauthorized, response.GeneralError(fmt.Errorf("Unauthorized")))
			return
		}

		// get data from parms
		name := r.PathValue("name")
		if name == " " {
			response.WriteJson(w, http.StatusNotFound, response.GeneralError(fmt.Errorf("name parms not found")))
			return
		}

		roomData, err := chatService.GetChatRoomByName(name)
		if err != nil {
			response.WriteJson(w, http.StatusInternalServerError, response.GeneralError(err))
			return
		}

		if userData.Role != "ADMIN" && roomData.UserId != userData.UserId {
			response.WriteJson(w, http.StatusUnauthorized, response.GeneralError(fmt.Errorf("unauthorized user")))
			return
		}

		err = chatService.DeleteChatRoom(name)
		if err != nil {
			response.WriteJson(w, http.StatusInternalServerError, response.GeneralError(err))
			return
		}

		response.WriteJson(w, http.StatusOK, "Chat Room Deleted")
		return
	}
}

func JoinRoom(chatService services.ChatService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Get value from context
		userDataRaw := r.Context().Value(middleware.UserDataKey)
		if userDataRaw == nil {
			response.WriteJson(w, http.StatusUnauthorized, response.GeneralError(fmt.Errorf("Unauthorized")))
			return
		}

		// Correct the type assertion to *models.AccessToken
		userData, ok := userDataRaw.(*models.AccessToken)
		if !ok {
			response.WriteJson(w, http.StatusUnauthorized, response.GeneralError(fmt.Errorf("Unauthorized")))
			return
		}

		// get data from parms
		name := r.PathValue("name")
		if name == " " {
			response.WriteJson(w, http.StatusNotFound, response.GeneralError(fmt.Errorf("name parms not found")))
			return
		}

		member, err := chatService.CheckChatRoomMember(userData.UserId, name)
		if err != nil {
			response.WriteJson(w, http.StatusInternalServerError, response.GeneralError(err))
			return
		}

		if member {
			response.WriteJson(w, http.StatusConflict, response.GeneralError(fmt.Errorf("Already Exists")))
			return
		}

		data := &models.JoinRoomRequest{
			UserId:   userData.UserId,
			RoomName: name,
		}
		err = chatService.JoinRoom(data)
		if err != nil {
			response.WriteJson(w, http.StatusInternalServerError, response.GeneralError(err))
			return
		}

		response.WriteJson(w, http.StatusOK, "Chat Room Join")
		return
	}
}

func GetAllJoinRoom(chatService services.ChatService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Get value from context
		userDataRaw := r.Context().Value(middleware.UserDataKey)
		if userDataRaw == nil {
			response.WriteJson(w, http.StatusUnauthorized, response.GeneralError(fmt.Errorf("Unauthorized")))
			return
		}

		// Correct the type assertion to *models.AccessToken
		userData, ok := userDataRaw.(*models.AccessToken)
		if !ok {
			response.WriteJson(w, http.StatusUnauthorized, response.GeneralError(fmt.Errorf("Unauthorized")))
			return
		}

		data, err := chatService.GetAllJoinRoom(userData.UserId)
		if err != nil {
			response.WriteJson(w, http.StatusInternalServerError, response.GeneralError(err))
			return
		}

		response.WriteJson(w, http.StatusOK, data)
		return
	}
}

func LeaveRoom(chatService services.ChatService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Get value from context
		userDataRaw := r.Context().Value(middleware.UserDataKey)
		if userDataRaw == nil {
			response.WriteJson(w, http.StatusUnauthorized, response.GeneralError(fmt.Errorf("Unauthorized")))
			return
		}

		// Correct the type assertion to *models.AccessToken
		userData, ok := userDataRaw.(*models.AccessToken)
		if !ok {
			response.WriteJson(w, http.StatusUnauthorized, response.GeneralError(fmt.Errorf("Unauthorized")))
			return
		}

		// get data from parms
		name := r.PathValue("name")
		if name == " " {
			response.WriteJson(w, http.StatusNotFound, response.GeneralError(fmt.Errorf("name parms not found")))
			return
		}

		err := chatService.LeaveRoom(userData.UserId, name)
		if err != nil {
			response.WriteJson(w, http.StatusInternalServerError, response.GeneralError(err))
			return
		}

		response.WriteJson(w, http.StatusOK, "User Leaved Room")
		return
	}
}

func GetOldChats(chatService services.ChatService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Get value from context
		userDataRaw := r.Context().Value(middleware.UserDataKey)
		if userDataRaw == nil {
			response.WriteJson(w, http.StatusUnauthorized, response.GeneralError(fmt.Errorf("Unauthorized")))
			return
		}

		// Correct the type assertion to *models.AccessToken
		userData, ok := userDataRaw.(*models.AccessToken)
		if !ok {
			response.WriteJson(w, http.StatusUnauthorized, response.GeneralError(fmt.Errorf("Unauthorized")))
			return
		}

		// get data from parms
		name := r.PathValue("roomName")
		limit := r.PathValue("limit")
		if name == " " || limit == " " {
			response.WriteJson(w, http.StatusNotFound, response.GeneralError(fmt.Errorf("parms not found")))
			return
		}

		intLimit, err := strconv.Atoi(limit)
		if err != nil {
			response.WriteJson(w, http.StatusInternalServerError, response.GeneralError(err))
			return
		}

		//check user join or not in room
		isMember, err := chatService.CheckChatRoomMember(userData.UserId, name)
		if err != nil {
			slog.Error(err.Error())
			response.WriteJson(w, http.StatusNotFound, response.GeneralError(fmt.Errorf("Error: something went worng.")))
			return
		}

		if !isMember {
			response.WriteJson(w, http.StatusNotFound, response.GeneralError(fmt.Errorf("Error: You are not a member of this group.")))
			return
		}

		oldMessages, err := chatService.GetOldMessages(name, intLimit)
		if err != nil {
			log.Println("Failed to fetch old messages:", err)
			return
		}

		response.WriteJson(w, http.StatusOK, oldMessages)
		return
	}
}
