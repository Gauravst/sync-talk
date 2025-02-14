package models_test

import (
	"testing"
	"time"

	"github.com/gauravst/real-time-chat/internal/models"
	"github.com/go-playground/validator/v10"
)

func TestUserValidation(t *testing.T) {
	validate := validator.New()

	tests := []struct {
		name    string
		user    models.User
		wantErr bool
	}{
		{
			name: "valid user",
			user: models.User{
				Name:      "John Doe",
				Username:  "johndoe",
				Email:     "john@example.com",
				Password:  "password123",
				CreatedAt: time.Now(),
				UpdatedAt: time.Now(),
			},
			wantErr: false,
		},
		{
			name: "missing name",
			user: models.User{
				Username:  "johndoe",
				Email:     "john@example.com",
				Password:  "password123",
				CreatedAt: time.Now(),
				UpdatedAt: time.Now(),
			},
			wantErr: true,
		},
		{
			name: "missing username",
			user: models.User{
				Name:      "John Doe",
				Email:     "john@example.com",
				Password:  "password123",
				CreatedAt: time.Now(),
				UpdatedAt: time.Now(),
			},
			wantErr: true,
		},
		{
			name: "invalid email",
			user: models.User{
				Name:      "John Doe",
				Username:  "johndoe",
				Email:     "invalid-email",
				Password:  "password123",
				CreatedAt: time.Now(),
				UpdatedAt: time.Now(),
			},
			wantErr: true,
		},
		{
			name: "missing email",
			user: models.User{
				Name:      "John Doe",
				Username:  "johndoe",
				Password:  "password123",
				CreatedAt: time.Now(),
				UpdatedAt: time.Now(),
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := validate.Struct(tt.user)
			if (err != nil) != tt.wantErr {
				t.Errorf("User validation failed: %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestUserTimestamps(t *testing.T) {
	user := models.User{
		Name:     "John Doe",
		Username: "johndoe",
		Email:    "john@example.com",
		Password: "password123",
	}

	if !user.CreatedAt.IsZero() || !user.UpdatedAt.IsZero() {
		t.Error("expected CreatedAt and UpdatedAt to be zero values initially")
	}

	now := time.Now()
	user.CreatedAt = now
	user.UpdatedAt = now

	if user.CreatedAt.IsZero() || user.UpdatedAt.IsZero() {
		t.Error("expected CreatedAt and UpdatedAt to be set")
	}

	if user.CreatedAt != now || user.UpdatedAt != now {
		t.Error("expected CreatedAt and UpdatedAt to match the current time")
	}
}
