package models

import "time"

type UploadedFile struct {
	Id               int       `json:"id"`
	PublicId         string    `json:"publicId"`
	SecureUrl        string    `json:"secureUrl"`
	Format           string    `json:"format"`
	ResourceType     string    `json:"resourceType"`
	Size             float64   `json:"bytes"`
	Width            int       `json:"width,omitempty"`
	Height           int       `json:"height,omitempty"`
	OriginalFilename string    `json:"originalFilename"`
	CreatedAt        time.Time `json:"created_at"`
	UpdatedAt        time.Time `json:"updated_at"`
}
