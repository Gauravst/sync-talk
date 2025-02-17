package models

type LoginRequest struct {
	Username    string `json:"username" validate:"required"`
	Password    string `json:"password,omitempty" validate:"required"`
	AccessToken string `json:"accessToken"`
}
