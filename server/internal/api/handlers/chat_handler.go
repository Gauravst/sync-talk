package handlers

import (
	"encoding/json"
	"fmt"
	"log/slog"
	"net/http"

	"github.com/gauravst/real-time-chat/internal/config"
	"github.com/gauravst/real-time-chat/internal/models"
	"github.com/gauravst/real-time-chat/internal/services"
	"github.com/gauravst/real-time-chat/internal/utils/response"
	"github.com/go-playground/validator/v10"
	"github.com/gorilla/websocket"
)

type contextKey string

const userDataKey contextKey = "userData"

func LiveChat(cfg config.Config, upgrader websocket.Upgrader) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		// Upgrade HTTP connection to WebSocket
		conn, err := upgrader.Upgrade(w, r, nil)
		if err != nil {
			slog.Error("failed to upgrade to WebSocket", slog.String("error", err.Error()))
			return
		}
		defer conn.Close()

		slog.Info("WebSocket connection established")

		// Handle WebSocket messages
		for {
			messageType, message, err := conn.ReadMessage()
			if err != nil {
				slog.Error("failed to read WebSocket message", slog.String("error", err.Error()))
				break
			}

			slog.Info("received WebSocket message", slog.String("message", string(message)))

			// Echo the message back to the client
			if err := conn.WriteMessage(messageType, message); err != nil {
				slog.Error("failed to write WebSocket message", slog.String("error", err.Error()))
				break
			}
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
