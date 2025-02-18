package models

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
