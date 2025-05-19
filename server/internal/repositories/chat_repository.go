package repositories

import (
	"database/sql"
	_ "embed"
	"errors"
	"log"

	"github.com/gauravst/real-time-chat/internal/database"
	"github.com/gauravst/real-time-chat/internal/models"
)

type ChatRepository interface {
	GetAllChatRoom(userData *models.AccessToken) ([]*models.ChatRoom, error)
	GetPrivateChatRoom(code string) (*models.ChatRoom, error)
	GetChatRoomByName(name string) (*models.ChatRoom, error)
	UpdateChatRoom(data *models.ChatRoomRequest) error
	DeleteChatRoom(name string) error
	CreateNewChatRoom(data *models.ChatRoomRequest) error
	CheckChatRoomMember(userId int, roomName string) (bool, error)
	GetOldMessages(roomName string, limit int) ([]*models.MessageResponse, error)
	CreateNewMessage(data *models.MessageResponse, roomName string) (*models.MessageResponse, error)
	JoinRoom(data *models.JoinRoomRequest) error
	JoinPrivateRoom(data *models.JoinRoomRequest) error
	GetAllJoinRoom(userId int) ([]*models.ChatRoom, error)
	LeaveRoom(userId int, roomName string) error
	GetFile(fileId *int) (*models.UploadedFile, error)
}

// userRepository implements the AuthRepository interface
type chatRepository struct {
	db      *sql.DB
	queries *database.QueryManager
}

func NewChatRepository(db *sql.DB, qm *database.QueryManager) ChatRepository {
	return &chatRepository{
		db:      db,
		queries: qm,
	}
}

func (r *chatRepository) GetAllChatRoom(userData *models.AccessToken) ([]*models.ChatRoom, error) {
	query, err := r.queries.Get("chat", "GetAllChatRoom")
	if err != nil {
		return nil, err
	}

	rows, err := r.db.Query(query, userData.UserId)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var data []*models.ChatRoom
	for rows.Next() {
		room := &models.ChatRoom{}
		err := rows.Scan(&room.Id, &room.Name, &room.Private, &room.Description, &room.UserId, &room.Members)
		if err != nil {
			return nil, err
		}

		data = append(data, room)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return data, nil
}

func (r *chatRepository) GetPrivateChatRoom(code string) (*models.ChatRoom, error) {
	data := &models.ChatRoom{}
	query, err := r.queries.Get("chat", "GetPrivateRoomUsingCode")
	if err != nil {
		return nil, err
	}

	err = r.db.QueryRow(query, code).Scan(&data.Id, &data.Name, &data.Private, &data.Description, &data.UserId, &data.Members)
	if err != nil {
		return data, err
	}
	return data, nil
}

func (r *chatRepository) GetChatRoomByName(name string) (*models.ChatRoom, error) {
	data := &models.ChatRoom{}
	query := `SELECT id, name, private, description, userId FROM chatRoom WHERE name = $1`
	err := r.db.QueryRow(query, name).Scan(&data.Id, &data.Name, &data.Private, &data.Description, &data.UserId)
	if err != nil {
		return data, err
	}
	return data, nil
}

func (r *chatRepository) UpdateChatRoom(data *models.ChatRoomRequest) error {
	query := `UPDATE chatRoom SET name = $1, userId = $2, members = $3, description = $4 WHERE name = $5 RETURNING id, name, members, description, userId`
	row := r.db.QueryRow(query, data.Name, data.UserId, data.Members, data.Description, data.Name)
	err := row.Scan(&data.Id, &data.Name, &data.Members, &data.Description, &data.UserId)
	if err != nil {
		return err
	}
	return nil
}

func (r *chatRepository) DeleteChatRoom(name string) error {
	query := `DELETE chatRoom WHERE id = $1`
	_, err := r.db.Exec(query)
	if err != nil {
		return err
	}
	return nil
}

func (r *chatRepository) CreateNewChatRoom(data *models.ChatRoomRequest) error {
	query := `INSERT INTO chatRoom (name, code, userId, description) VALUES ($1, $2, $3, $4)`
	_, err := r.db.Exec(query, data.Name, data.Code, data.UserId, data.Description)
	if err != nil {
		return err
	}
	return nil
}

func (r *chatRepository) CheckChatRoomMember(userId int, roomName string) (bool, error) {
	query := `SELECT EXISTS(SELECT 1 FROM groupMembers WHERE userId = $1 AND roomname = $2)`
	var exists bool
	err := r.db.QueryRow(query, userId, roomName).Scan(&exists)
	if err != nil {
		log.Printf("Error checking chat room member: %v", err)
		return false, err
	}
	return exists, nil
}

func (r *chatRepository) GetOldMessages(roomName string, limit int) ([]*models.MessageResponse, error) {
	if roomName == "" || limit <= 0 {
		return nil, errors.New("invalid room name or limit")
	}

	query, err := r.queries.Get("chat", "GetOldMessages")
	if err != nil {
		return nil, err
	}

	rows, err := r.db.Query(query, roomName, limit)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var messages []*models.MessageResponse

	for rows.Next() {
		msg := &models.MessageResponse{}
		file := &models.UploadedFile{}
		var fileId sql.NullInt64
		var publicId, secureUrl, format, resourceType, originalFilename sql.NullString
		var size sql.NullFloat64
		var width, height sql.NullInt64
		var fileCreatedAt, fileUpdatedAt sql.NullTime

		err := rows.Scan(
			&msg.Id, &msg.UserId, &msg.Username, &msg.Content, &msg.RoomName,
			&msg.CreatedAt, &msg.UpdatedAt,
			&fileId, &publicId, &secureUrl, &format, &resourceType, &size,
			&width, &height, &originalFilename, &fileCreatedAt, &fileUpdatedAt,
		)
		if err != nil {
			return nil, err
		}

		if fileId.Valid {
			msg.FileId = intPtr(int(fileId.Int64))
			file.Id = int(fileId.Int64)
			file.PublicId = publicId.String
			file.SecureUrl = secureUrl.String
			file.Format = format.String
			file.ResourceType = resourceType.String
			file.Size = size.Float64
			file.Width = int(width.Int64)
			file.Height = int(height.Int64)
			file.OriginalFilename = originalFilename.String
			file.CreatedAt = fileCreatedAt.Time
			file.UpdatedAt = fileUpdatedAt.Time
			msg.File = file
		}

		messages = append(messages, msg)
	}

	if len(messages) == 0 {
		return nil, errors.New("no messages found")
	}

	return messages, nil
}

func intPtr(i int) *int {
	return &i
}

func (r *chatRepository) CreateNewMessage(data *models.MessageResponse, roomName string) (*models.MessageResponse, error) {
	message := &models.MessageResponse{}

	if data.FileId != nil {
		query := `
			INSERT INTO messages (userId, roomName, content, fileId) 
			VALUES ($1, $2, $3, $4) 
			RETURNING id, userId, roomName, fileId, content, createdAt, updatedAt
		`

		err := r.db.QueryRow(query, data.UserId, roomName, data.Content, *data.FileId).Scan(
			&message.Id, &message.UserId, &message.RoomName, &message.FileId,
			&message.Content, &message.CreatedAt, &message.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}
	} else {
		query := `
			INSERT INTO messages (userId, roomName, content) 
			VALUES ($1, $2, $3) 
			RETURNING id, userId, roomName, fileId, content, createdAt, updatedAt
		`

		err := r.db.QueryRow(query, data.UserId, roomName, data.Content).Scan(
			&message.Id, &message.UserId, &message.RoomName, &message.FileId,
			&message.Content, &message.CreatedAt, &message.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}
	}

	return message, nil
}

func (r *chatRepository) JoinRoom(data *models.JoinRoomRequest) error {
	query := `INSERT INTO groupMembers (userId, roomName) VALUES ($1, $2)`
	_, err := r.db.Exec(query, data.UserId, data.RoomName)
	if err != nil {
		return err
	}
	return nil
}

func (r *chatRepository) JoinPrivateRoom(data *models.JoinRoomRequest) error {
	query := `INSERT INTO groupMembers (userId, roomName) VALUES ($1, $2)`
	_, err := r.db.Exec(query, data.UserId, data.RoomName)
	if err != nil {
		return err
	}
	return nil
}

func (r *chatRepository) GetAllJoinRoom(userId int) ([]*models.ChatRoom, error) {
	query, err := r.queries.Get("chat", "GetAllJoinRoom")
	if err != nil {
		return nil, err
	}

	rows, err := r.db.Query(query, userId)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var data []*models.ChatRoom
	for rows.Next() {
		room := &models.ChatRoom{}
		err := rows.Scan(&room.Id, &room.Name, &room.Description, &room.Private, &room.UserId, &room.Members)
		if err != nil {
			return nil, err
		}

		data = append(data, room)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return data, nil
}

func (r *chatRepository) LeaveRoom(userId int, roomName string) error {
	query := `DELETE FROM groupMembers WHERE userId = $1 AND roomName = $2`
	_, err := r.db.Exec(query, userId, roomName)
	if err != nil {
		return err
	}

	return nil
}

func (r *chatRepository) GetFile(fileId *int) (*models.UploadedFile, error) {
	data := &models.UploadedFile{}
	query := `SELECT id, publicid, secureurl, format, resourcetype, size, width, height, originalfilename FROM WHERE id = $1`
	err := r.db.QueryRow(query, fileId).Scan(&data.Id, &data.PublicId, &data.SecureUrl, &data.Format, &data.ResourceType, &data.Size, &data.Width, &data.Height, &data.OriginalFilename)
	if err != nil {
		return data, err
	}

	return data, nil
}
