package services

import (
	"fmt"
	"log/slog"

	"github.com/gauravst/real-time-chat/internal/models"
	"github.com/gauravst/real-time-chat/internal/repositories"
	randomstring "github.com/gauravst/real-time-chat/internal/utils/randomString"
)

type ChatService interface {
	GetAllChatRoom(userData *models.AccessToken) ([]*models.ChatRoom, error)
	GetPrivateChatRoom(code string) (*models.ChatRoom, error)
	GetChatRoomByName(name string) (*models.ChatRoom, error)
	UpdateChatRoom(data *models.ChatRoomRequest) error
	DeleteChatRoom(name string) error
	CreateNewChatRoom(data *models.ChatRoomRequest) error
	CheckChatRoomMember(userId int, roomName string) (bool, error)
	GetOldMessages(roomName string, limit int) ([]*models.MessageRequest, error)
	CreateNewMessage(data *models.MessageRequest, roomName string) (*models.MessageResponse, error)
	JoinRoom(data *models.JoinRoomRequest) error
	JoinPrivateRoom(code string, userData *models.AccessToken) error
	GetAllJoinRoom(userId int) ([]*models.ChatRoom, error)
	LeaveRoom(userId int, roomName string) error
}

type chatService struct {
	chatRepo repositories.ChatRepository
}

func NewChatService(chatRepo repositories.ChatRepository) ChatService {
	return &chatService{
		chatRepo: chatRepo,
	}
}

func (s *chatService) GetAllChatRoom(userData *models.AccessToken) ([]*models.ChatRoom, error) {
	var data []*models.ChatRoom
	data, err := s.chatRepo.GetAllChatRoom(userData)
	if err != nil {
		return data, err
	}

	return data, nil
}

func (s *chatService) GetPrivateChatRoom(code string) (*models.ChatRoom, error) {
	var data *models.ChatRoom
	data, err := s.chatRepo.GetPrivateChatRoom(code)
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
	code := randomstring.GenerateRandomString(5)
	data.Code = code
	err := s.chatRepo.CreateNewChatRoom(data)
	if err != nil {
		return err
	}

	joinRoomData := &models.JoinRoomRequest{
		UserId:   data.UserId,
		RoomName: data.Name,
	}
	err = s.chatRepo.JoinRoom(joinRoomData)
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
		return nil, err
	}
	return data, nil
}

func (s *chatService) CreateNewMessage(data *models.MessageRequest, roomName string) (*models.MessageResponse, error) {
	messageData, err := s.chatRepo.CreateNewMessage(data, roomName)
	if err != nil {
		return nil, err
	}
	return messageData, nil
}

func (s *chatService) JoinRoom(data *models.JoinRoomRequest) error {
	err := s.chatRepo.JoinRoom(data)
	if err != nil {
		return err
	}
	return nil
}

func (s *chatService) JoinPrivateRoom(code string, userData *models.AccessToken) error {
	//check private room using room code
	roomData, err := s.chatRepo.GetPrivateChatRoom(code)
	if err != nil {
		slog.Error(err.Error())
		return err
	}

	member, err := s.chatRepo.CheckChatRoomMember(userData.UserId, roomData.Name)
	if err != nil {
		return err
	}

	if member {
		return fmt.Errorf("You are already member of this room")
	}

	data := &models.JoinRoomRequest{
		UserId:   userData.UserId,
		RoomName: roomData.Name,
	}

	err = s.chatRepo.JoinRoom(data)
	if err != nil {
		return err
	}
	return nil
}

func (s *chatService) GetAllJoinRoom(userId int) ([]*models.ChatRoom, error) {
	var data []*models.ChatRoom
	data, err := s.chatRepo.GetAllJoinRoom(userId)
	if err != nil {
		return data, err
	}
	return data, nil
}

func (s *chatService) LeaveRoom(userId int, roomName string) error {
	roomData, err := s.chatRepo.GetChatRoomByName(roomName)
	if err != nil {
		return err
	}

	if roomData.UserId == userId {
		return fmt.Errorf("you can not leave from you room")
	}

	err = s.chatRepo.LeaveRoom(userId, roomName)
	if err != nil {
		return err
	}
	return nil
}
