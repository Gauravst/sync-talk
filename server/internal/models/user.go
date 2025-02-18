package models

import "time"

type User struct {
	Id        int       `json:"id"`
	Username  string    `json:"username" validate:"required"`
	Password  string    `json:"password"`
	Role      string    `json:"role"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type LoginSession struct {
	Id     int    `json:"id"`
	UserId int    `json:"userId"`
	Token  string `json:"token"`
}
