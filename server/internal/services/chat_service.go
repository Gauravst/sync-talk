package services

import (
	"github.com/gauravst/real-time-chat/internal/models"
	"github.com/gauravst/real-time-chat/internal/repositories"
)

type ChatService interface {
	GetAllChatRoom() ([]*models.ChatRoom, error)
	GetChatRoomByName(name string) (*models.ChatRoom, error)
	UpdateChatRoom(data *models.ChatRoomRequest) error
	DeleteChatRoom(name string) error
	CreateNewChatRoom(data *models.ChatRoomRequest) error
	CheckChatRoomMember(userId int, roomName string) (bool, error)
	GetOldMessages(roomName string, limit int) ([]*models.MessageRequest, error)
	CreateNewMessage(data *models.MessageRequest, roomName string) error
}

type chatService struct {
	chatRepo repositories.ChatRepository
}

func NewChatService(chatRepo repositories.ChatRepository) ChatService {
	return &chatService{
		chatRepo: chatRepo,
	}
}

func (s *chatService) GetAllChatRoom() ([]*models.ChatRoom, error) {
	var data []*models.ChatRoom
	data, err := s.chatRepo.GetAllChatRoom()
	if err != nil {
		return data, err
	}

	return data, nil
}

func (s *chatService) GetChatRoomByName(name string) (*models.ChatRoom, error) {
	var data *models.ChatRoom
	data, err := s.chatRepo.GetChatRoomByName(name)
	if err != nil {
		return data, err
	}

	return data, nil
}

func (s *chatService) UpdateChatRoom(data *models.ChatRoomRequest) error {
	err := s.chatRepo.UpdateChatRoom(data)
	if err != nil {
		return err
	}
	return nil
}

func (s *chatService) DeleteChatRoom(name string) error {
	err := s.chatRepo.DeleteChatRoom(name)
	if err != nil {
		return err
	}
	return nil
}

func (s *chatService) CreateNewChatRoom(data *models.ChatRoomRequest) error {
	err := s.chatRepo.CreateNewChatRoom(data)
	if err != nil {
		return err
	}
	return nil
}

func (s *chatService) CheckChatRoomMember(userId int, roomName string) (bool, error) {
	var exists bool
	exists, err := s.chatRepo.CheckChatRoomMember(userId, roomName)
	if err != nil {
		return false, err
	}
	return exists, nil
}

func (s *chatService) GetOldMessages(roomName string, limit int) ([]*models.MessageRequest, error) {
	var data []*models.MessageRequest
	data, err := s.chatRepo.GetOldMessages(roomName, limit)
	if err != nil {
		return data, err
	}
	return data, nil
}

func (s *chatService) CreateNewMessage(data *models.MessageRequest, roomName string) error {
	err := s.chatRepo.CreateNewMessage(data, roomName)
	if err != nil {
		return err
	}
	return nil
}
