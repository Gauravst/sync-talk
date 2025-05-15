package handlers

import (
	"encoding/json"
	"fmt"
	"log"
	"log/slog"
	"net/http"
	"strconv"

	"github.com/gauravst/real-time-chat/internal/api/middleware"
	"github.com/gauravst/real-time-chat/internal/config"
	"github.com/gauravst/real-time-chat/internal/models"
	"github.com/gauravst/real-time-chat/internal/services"
	"github.com/gauravst/real-time-chat/internal/utils/response"
	"github.com/gauravst/real-time-chat/internal/utils/ws"
	"github.com/go-playground/validator/v10"
	"github.com/gorilla/websocket"
)

func LiveChat(chatService services.ChatService, cfg config.Config, wsServer *models.WsServer) http.HandlerFunc {
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
		conn, err := wsServer.Upgrader.Upgrade(w, r, nil)
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
		wsServer.RoomMutex.Lock()
		wsServer.Rooms[roomName] = append(wsServer.Rooms[roomName], conn)
		if wsServer.OnlineUser[roomName] == nil {
			wsServer.OnlineUser[roomName] = make(map[string]bool)
		}
		wsServer.OnlineUser[roomName][userData.Username] = true

		// check user Already in connection so we not get worng online count
		// if users, ok := wsServer.OnlineUser[roomName]; ok {
		// 	for _, user := range users {
		// 		if user != userData.Username {
		// 			wsServer.OnlineUser[roomName] = append(wsServer.OnlineUser[roomName], userData.Username)
		// 		}
		// 	}
		// }

		wsServer.RoomMutex.Unlock()

		//   count = len()
		// slog.Info(fmt.Sprintf("Number of conn : %d", count))
		// users := []string{}
		// for user := range wsServer.OnlineUser[roomName] {
		// 	users = append(users, user)
		// }
		// slog.Info(fmt.Sprintf("Online users in room %s: %+v", roomName, users))

		// Broadcast the updated online user count
		go broadcastOnlineUsers(roomName, wsServer)

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
			go ws.BroadcastMessage(wsServer, roomName, conn, createdMessage)
		}

		// remove connection
		removeConnection(roomName, conn, wsServer, userData.Username)
	}
}

func broadcastOnlineUsers(roomName string, wsServer *models.WsServer) {
	wsServer.RoomMutex.Lock()
	defer wsServer.RoomMutex.Unlock()

	clients := wsServer.Rooms[roomName]
	count := len(wsServer.OnlineUser[roomName])

	slog.Info(fmt.Sprintf("Number of online users: %d", count))
	users := []string{}
	for user := range wsServer.OnlineUser[roomName] {
		users = append(users, user)
	}
	slog.Info(fmt.Sprintf("Online users in room %s: %+v", roomName, users))

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

func removeConnection(roomName string, conn *websocket.Conn, wsServer *models.WsServer, username string) {
	wsServer.RoomMutex.Lock()
	defer wsServer.RoomMutex.Unlock()

	clients := wsServer.Rooms[roomName]
	for i, c := range clients {
		if c == conn {
			wsServer.Rooms[roomName] = append(clients[:i], clients[i+1:]...)
			break
		}
	}

	count := len(wsServer.OnlineUser[roomName])

	// decrease the online user count
	if count > 0 {
		delete(wsServer.OnlineUser[roomName], username)
	}

	// Broadcast updated online users count
	go broadcastOnlineUsers(roomName, wsServer)
}

func GetAllChatRoom(chatService services.ChatService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
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
		data, err := chatService.GetAllChatRoom(userData)
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

		response.WriteJson(w, http.StatusOK, "Chat Room Joined")
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

func GetPrivateChatRoom(chatService services.ChatService) http.HandlerFunc {
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
		code := r.PathValue("code")
		if code == " " {
			response.WriteJson(w, http.StatusNotFound, response.GeneralError(fmt.Errorf("parms not found")))
			return
		}

		//check private room using room code
		roomData, err := chatService.GetPrivateChatRoom(code)
		if err != nil {
			slog.Error(err.Error())
			response.WriteJson(w, http.StatusNotFound, response.GeneralError(fmt.Errorf("Error: room not found")))
			return
		}

		//check user join or not in room
		isMember, err := chatService.CheckChatRoomMember(userData.UserId, roomData.Name)
		if err != nil {
			slog.Error(err.Error())
			response.WriteJson(w, http.StatusNotFound, response.GeneralError(fmt.Errorf("Error: something went worng.")))
			return
		}

		data := &models.PrivateRoomUsingCodeResponse{
			Id:          roomData.Id,
			Name:        roomData.Name,
			Members:     roomData.Members,
			Code:        roomData.Code,
			Description: roomData.Description,
			UserId:      roomData.UserId,
			IsMember:    isMember,
		}

		response.WriteJson(w, http.StatusOK, data)
		return
	}
}

func JoinPrivateRoom(chatService services.ChatService) http.HandlerFunc {
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
		code := r.PathValue("code")
		if code == " " {
			response.WriteJson(w, http.StatusNotFound, response.GeneralError(fmt.Errorf("parms not found")))
			return
		}

		// JoinPrivateRoom
		err := chatService.JoinPrivateRoom(code, userData)
		if err != nil {
			slog.Error(err.Error())
			response.WriteJson(w, http.StatusInternalServerError, err)
			return
		}

		response.WriteJson(w, http.StatusOK, "Room Joined")
		return
	}
}
