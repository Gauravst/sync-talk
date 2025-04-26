package models

type PrivateRoomUsingCodeResponse struct {
	Id          int    `json:"id"`
	Name        string `json:"name" validate:"required"`
	Members     int    `json:"members"`
	Code        string `json:"code"`
	Description string `json:"description" validate:"required"`
	UserId      int    `json:"userId"`
	IsMember    bool   `json:"isMember"`
}
