package handlers

import (
	"encoding/json"
	"fmt"
	"log"
	"log/slog"
	"net/http"
	"sync"

	"github.com/gauravst/real-time-chat/internal/config"
	"github.com/gauravst/real-time-chat/internal/models"
	"github.com/gauravst/real-time-chat/internal/services"
	"github.com/gauravst/real-time-chat/internal/utils/response"
	"github.com/go-playground/validator/v10"
	"github.com/gorilla/websocket"
)

type contextKey string

const userDataKey contextKey = "userData"

// upgrader to upgrade HTTP connection to Websocket
var (
	upgrader  = websocket.Upgrader{CheckOrigin: func(r *http.Request) bool { return true }}
	rooms     = make(map[string][]*websocket.Conn)
	roomMutex sync.Mutex
)

func LiveChat(chatService services.ChatService, cfg config.Config) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		userData, ok := r.Context().Value(userDataKey).(models.User)
		if !ok {
			http.Error(w, "unauthorized user", http.StatusUnauthorized)
			return
		}

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
		isMember, err := chatService.CheckChatRoomMember(userData.Id, roomName)
		if err != nil {
			slog.Error(err.Error())
			conn.WriteMessage(websocket.TextMessage, []byte("Error: something went worng."))
			return
		}

		if !isMember {
			conn.WriteMessage(websocket.TextMessage, []byte("Error: You are not a member of this group."))
			return
		}

		oldMessages, err := chatService.GetOldMessages(roomName, 20)
		if err != nil {
			log.Println("Failed to fetch old messages:", err)
			return
		}

		// Send old messages to the user
		for _, msg := range oldMessages {
			msgJSON, _ := json.Marshal(msg)
			conn.WriteMessage(websocket.TextMessage, msgJSON)
		}

		// Add connection to the room
		roomMutex.Lock()
		rooms[roomName] = append(rooms[roomName], conn)
		roomMutex.Unlock()

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
				Content: msg.Content,
				UserId:  userData.Id,
			}
			err = chatService.CreateNewMessage(newMessageData, roomName)
			if err != nil {
				slog.Error("failed to save message", slog.String("error", err.Error()))
				continue
			}

			// send message
			broadcastMessage(roomName, message)
		}

		// remove connection
		removeConnection(roomName, conn)
	}
}

func broadcastMessage(roomName string, message []byte) {
	roomMutex.Lock()
	defer roomMutex.Unlock()

	clients := rooms[roomName]
	for _, conn := range clients {
		if err := conn.WriteMessage(websocket.TextMessage, message); err != nil {
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
		// get user data from auth middleware
		var userData models.User
		userData, ok := r.Context().Value(userDataKey).(models.User)
		if !ok {
			response.WriteJson(w, http.StatusInternalServerError, response.GeneralError(fmt.Errorf("user data not found")))
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
		data.UserId = userData.Id

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
		// get user data from auth middleware
		var userData models.User
		userData, ok := r.Context().Value(userDataKey).(models.User)
		if !ok {
			response.WriteJson(w, http.StatusInternalServerError, response.GeneralError(fmt.Errorf("user data not found")))
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

		if userData.Role != "ADMIN" && roomData.UserId != userData.Id {
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
		// get user data from auth middleware
		var userData models.User
		userData, ok := r.Context().Value(userDataKey).(models.User)
		if !ok {
			response.WriteJson(w, http.StatusInternalServerError, response.GeneralError(fmt.Errorf("user data not found")))
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

		if userData.Role != "ADMIN" && roomData.UserId != userData.Id {
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
		// get user data from auth middleware
		var userData models.User
		userData, ok := r.Context().Value(userDataKey).(models.User)
		if !ok {
			response.WriteJson(w, http.StatusInternalServerError, response.GeneralError(fmt.Errorf("user data not found")))
			return
		}

		// get data from parms
		name := r.PathValue("name")
		if name == " " {
			response.WriteJson(w, http.StatusNotFound, response.GeneralError(fmt.Errorf("name parms not found")))
			return
		}

		member, err := chatService.CheckChatRoomMember(userData.Id, name)
		if err != nil {
			response.WriteJson(w, http.StatusInternalServerError, response.GeneralError(err))
			return
		}

		if member {
			response.WriteJson(w, http.StatusConflict, response.GeneralError(fmt.Errorf("Already Exists")))
			return
		}

		data := &models.JoinRoomRequest{
			UserId:   userData.Id,
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
