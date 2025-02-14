package models

import "time"

type User struct {
	ID        int       `json:"id"`
	Name      string    `json:"name" validate:"required"`
	Username  string    `json:"username" validate:"required"`
	Email     string    `json:"email" validate:"required,email"`
	Password  string    `json:"-"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}
