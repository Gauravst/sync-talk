package models

import "time"

type PrivateRoomUsingCodeResponse struct {
	Id          int    `json:"id"`
	Name        string `json:"name" validate:"required"`
	Members     int    `json:"members"`
	Code        string `json:"code"`
	Description string `json:"description" validate:"required"`
	UserId      int    `json:"userId"`
	IsMember    bool   `json:"isMember"`
}

type MessageResponse struct {
	Id        int           `json:"id"`
	Type      string        `json:"type"`
	UserId    int           `json:"userId" validate:"required"`
	Username  string        `json:"username"`
	RoomName  string        `json:"roomName" validate:"required"`
	Content   string        `json:"content" validate:"required"`
	File      *UploadedFile `json:"file,omitempty"`
	CreatedAt time.Time     `json:"created_at"`
	UpdatedAt time.Time     `json:"updated_at"`
}
