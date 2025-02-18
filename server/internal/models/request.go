package models

import "time"

type UserRequest struct {
	Id       int    `json:"id"`
	Username string `json:"username" validate:"required"`
	Password string `json:"password, omitempty" validate:"required"`
	Role     string `json:"role"`
}

type LoginRequest struct {
	Username    string `json:"username" validate:"required"`
	Password    string `json:"password,omitempty" validate:"required"`
	AccessToken string `json:"accessToken"`
}

type ChatRoomRequest struct {
	Id         int    `json:"id"`
	Name       string `json:"name" validate:"required"`
	ProfilePic string `json:"profilePic"`
	UserId     int    `json:"userId", validate:"required"`
}

type MessageRequest struct {
	Id        int       `json:"id"`
	UserId    int       `json:"userId" validate:"required"`
	RoomName  string    `json:"roomName" validate:"required"`
	Content   string    `json:"content" validate:"required"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type JoinRoomRequest struct {
	Id       int    `json:"id"`
	UserId   int    `json:"userId" validate:"required"`
	RoomName string `json:"roomName" validate:"required"`
}
