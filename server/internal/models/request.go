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
	Id          int    `json:"id"`
	Name        string `json:"name" validate:"required"`
	Members     int    `json:"members"`
	Code        string `json:"code"`
	Description string `json:"description" validate:"required"`
	UserId      int    `json:"userId"`
}

type MessageRequest struct {
	Id        int       `json:"id"`
	Type      string    `json:"type"`
	UserId    int       `json:"userId" validate:"required"`
	Username  string    `json:"username"`
	RoomName  string    `json:"roomName" validate:"required"`
	Content   string    `json:"content" validate:"required"`
	File      int       `json:"file"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type JoinRoomRequest struct {
	Id       int    `json:"id"`
	UserId   int    `json:"userId" validate:"required"`
	RoomName string `json:"roomName" validate:"required"`
}

type OnlineUserCountRequest struct {
	Type  string `json:"type"`
	Count int    `json:"count"`
}
