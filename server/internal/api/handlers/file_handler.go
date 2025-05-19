package handlers

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"os"

	"github.com/gauravst/real-time-chat/internal/api/middleware"
	"github.com/gauravst/real-time-chat/internal/config"
	"github.com/gauravst/real-time-chat/internal/models"
	"github.com/gauravst/real-time-chat/internal/services"
	"github.com/gauravst/real-time-chat/internal/utils/response"
)

func UploadFileInRoom(fileService services.FileService, cfg config.Config, wsServer *models.WsServer) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

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

		roomName := r.PathValue("roomName")
		if roomName == "" {
			response.WriteJson(w, http.StatusBadRequest, response.GeneralError(fmt.Errorf("Missing room Name")))
			return
		}

		err := r.ParseMultipartForm(10 << 20) // Limit upload size to 10 MB
		fmt.Print(err)
		if err != nil {
			response.WriteJson(w, http.StatusBadRequest, response.GeneralError(fmt.Errorf("Could not parse multipart form")))
			return
		}

		content := r.FormValue("message")
		fmt.Print("Content --------\n")
		fmt.Print(content)
		file, _, err := r.FormFile("file")
		if err != nil {
			response.WriteJson(w, http.StatusBadRequest, response.GeneralError(fmt.Errorf("File missing or invalid")))
			return
		}
		defer file.Close()

		//save file in temp dir
		uploadDir := "uploads"              // relative to your project
		os.MkdirAll(uploadDir, os.ModePerm) // ensure it exists

		tempFile, err := os.CreateTemp(uploadDir, "upload-*.png")
		if err != nil {
			log.Println("Failed to create temp file:", err)
			http.Error(w, "Failed to save file", http.StatusInternalServerError)
			return
		}
		defer tempFile.Close()

		_, err = io.Copy(tempFile, file)
		if err != nil {
			response.WriteJson(w, http.StatusInternalServerError, response.GeneralError(fmt.Errorf("Cannot save file", err)))
			return
		}

		filePath := tempFile.Name()
		err = fileService.UploadFileInRoom(cfg, filePath, content, roomName, userData, wsServer)
		if err != nil {
			fmt.Print(err)
			response.WriteJson(w, http.StatusInternalServerError, response.GeneralError(fmt.Errorf("Something went worng", err)))
			return
		}

		response.WriteJson(w, http.StatusOK, map[string]string{"success": "ok"})
		return
	}
}
