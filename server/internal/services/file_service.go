package services

import (
	"context"
	"fmt"
	"log"

	"github.com/cloudinary/cloudinary-go/v2"
	"github.com/cloudinary/cloudinary-go/v2/api/uploader"
	"github.com/gauravst/real-time-chat/internal/config"
	"github.com/gauravst/real-time-chat/internal/models"
	"github.com/gauravst/real-time-chat/internal/repositories"
	"github.com/gauravst/real-time-chat/internal/utils/ws"
)

type FileService interface {
	UploadFileInRoom(cfg config.Config, filePath string, content string, roomName string, userData *models.AccessToken, wsServer *models.WsServer) error
}

type fileService struct {
	fileRepo repositories.FileRepository
	chatRepo repositories.ChatRepository
}

func NewFileService(fileRepo repositories.FileRepository, chatRepo repositories.ChatRepository) FileService {
	return &fileService{
		fileRepo: fileRepo,
		chatRepo: chatRepo,
	}
}

func (s *fileService) UploadFileInRoom(cfg config.Config, filePath string, content string, roomName string, userData *models.AccessToken, wsServer *models.WsServer) error {
	// Initialize Cloudinary instance
	cld, err := cloudinary.NewFromParams(cfg.Cloudinary.Name, cfg.Cloudinary.Key, cfg.Cloudinary.SecretKey)
	if err != nil {
		log.Fatalf("Failed to initialize Cloudinary: %v", err)
		return fmt.Errorf("Failed to initialize Cloudinary: %v", err)
	}

	// Upload an image
	uploadResult, err := cld.Upload.Upload(context.Background(), filePath, uploader.UploadParams{})
	if err != nil {
		log.Fatalf("Failed to upload image: %v", err)
		return fmt.Errorf("Failed to upload image: %v", err)
	}

	fileData := &models.UploadedFile{
		PublicId:         uploadResult.PublicID,
		SecureUrl:        uploadResult.SecureURL,
		Format:           uploadResult.Format,
		ResourceType:     uploadResult.ResourceType,
		Size:             float64(uploadResult.Bytes) / 1024,
		Width:            uploadResult.Width,
		Height:           uploadResult.Height,
		OriginalFilename: uploadResult.OriginalFilename,
		CreatedAt:        uploadResult.CreatedAt,
	}

	// add data in db
	err = s.fileRepo.UploadFileInRoom(fileData)
	if err != nil {
		return fmt.Errorf("something went worng")
	}

	data := &models.MessageRequest{
		Type:     "Chat",
		UserId:   userData.UserId,
		Username: userData.Username,
		RoomName: roomName,
		Content:  content,
		File:     fileData.Id,
	}

	messageData, err := s.chatRepo.CreateNewMessage(data, roomName)
	if err != nil {
		return fmt.Errorf("something went worng", err)
	}

	// send data in websoket
	ws.BroadcastMessage(wsServer, roomName, nil, messageData)

	return nil
}
