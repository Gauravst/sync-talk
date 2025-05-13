package repositories

import (
	"database/sql"

	"github.com/gauravst/real-time-chat/internal/database"
	"github.com/gauravst/real-time-chat/internal/models"
)

type FileRepository interface {
	UploadFileInRoom(fileData *models.UploadedFile) error
}

type fileRepository struct {
	db      *sql.DB
	queries *database.QueryManager
}

func NewFileRepository(db *sql.DB, qm *database.QueryManager) FileRepository {
	return &fileRepository{
		db:      db,
		queries: qm,
	}
}

func (r *fileRepository) UploadFileInRoom(fileData *models.UploadedFile) error {
	query, err := r.queries.Get("file", "UploadFile")
	if err != nil {
		return err
	}

	err = r.db.QueryRow(query, fileData.PublicId, fileData.SecureUrl, fileData.Format, fileData.ResourceType, fileData.Size, fileData.Width, fileData.Height, fileData.OriginalFilename).Scan(&fileData.Id, &fileData.PublicId, &fileData.SecureUrl, &fileData.Format, &fileData.ResourceType, &fileData.Size, &fileData.Width, &fileData.Height, &fileData.OriginalFilename, &fileData.CreatedAt)
	if err != nil {
		return err
	}

	return nil
}
