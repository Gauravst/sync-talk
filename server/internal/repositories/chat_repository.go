package repositories

import (
	"database/sql"

	"github.com/gauravst/real-time-chat/internal/models"
)

type ChatRepository interface {
	GetAllChatRoom() ([]*models.ChatRoom, error)
	GetChatRoomByName(name string) (*models.ChatRoom, error)
	UpdateChatRoom(data *models.ChatRoomRequest) error
	DeleteChatRoom(name string) error
	CreateNewChatRoom(data *models.ChatRoomRequest) error
	CheckChatRoomMember(userId int, roomName string) (bool, error)
	GetOldMessages(roomName string, limit int) ([]*models.MessageRequest, error)
	CreateNewMessage(data *models.MessageRequest, roomName string) error
	JoinRoom(data *models.JoinRoomRequest) error
	GetAllJoinRoom(userId int) ([]*models.ChatRoom, error)
}

// userRepository implements the AuthRepository interface
type chatRepository struct {
	db *sql.DB
}

func NewChatRepository(db *sql.DB) ChatRepository {
	return &chatRepository{
		db: db,
	}
}

func (r *chatRepository) GetAllChatRoom() ([]*models.ChatRoom, error) {
	query := `SELECT id, name , userId, profilePic FROM chatRoom`

	rows, err := r.db.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var data []*models.ChatRoom
	for rows.Next() {
		room := &models.ChatRoom{}
		err := rows.Scan(&room.Id, &room.Name, &room.UserId, &room.ProfilePic)
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

func (r *chatRepository) GetChatRoomByName(name string) (*models.ChatRoom, error) {
	var data *models.ChatRoom
	query := `SELECT id, name, userId, profilePic FROM chatRoom WHERE name = $1`
	err := r.db.QueryRow(query, name).Scan(&data.Id, &data.Name, &data.UserId, &data.ProfilePic)
	if err != nil {
		return data, err
	}
	return data, nil
}

func (r *chatRepository) UpdateChatRoom(data *models.ChatRoomRequest) error {
	query := `UPDATE chatRoom SET name = $1, userId = $2, profilePic = $3 WHERE name = $4 RETURNING id, name, userId, profilePic`
	row := r.db.QueryRow(query, data.Name, data.UserId, data.ProfilePic, data.Name)
	err := row.Scan(&data.Id, &data.Name, &data.UserId, &data.ProfilePic)
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
	query := `INSERT INTO chatRoom (name, userId) VALUES ($1, $2)`
	_, err := r.db.Exec(query, data.Name, data.UserId)
	if err != nil {
		return err
	}
	return nil
}

func (r *chatRepository) CheckChatRoomMember(userId int, roomName string) (bool, error) {
	var exists bool
	query := `SELECT EXISTS(SELECT 1 FROM groupMembers WHERE userId = $1 AND roomName = $2)`
	err := r.db.QueryRow(query, userId, roomName).Scan(&exists)
	return exists, err
}

func (r *chatRepository) GetOldMessages(roomName string, limit int) ([]*models.MessageRequest, error) {
	var messages []*models.MessageRequest
	query := `SELECT userId, content, createdAt FROM messages WHERE roomName = $1 ORDER BY createdAt DESC LIMIT $2`

	rows, err := r.db.Query(query, roomName, limit)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var msg *models.MessageRequest
		err := rows.Scan(&msg.UserId, &msg.Content, &msg.CreatedAt)
		if err != nil {
			return nil, err
		}
		messages = append(messages, msg)
	}
	return messages, nil
}

func (r *chatRepository) CreateNewMessage(data *models.MessageRequest, roomName string) error {
	query := `INSERT INTO messages (userId, roomName, content ) VALUES ($1, $2, $3)`
	_, err := r.db.Exec(query, data.UserId, roomName, data.Content)
	if err != nil {
		return err
	}
	return nil
}

func (r *chatRepository) JoinRoom(data *models.JoinRoomRequest) error {
	query := `INSERT INTO groupMembers (userId, roomName) VALUES ($1, $2)`
	_, err := r.db.Exec(query, data.UserId, data.RoomName)
	if err != nil {
		return err
	}
	return nil
}

func (r *chatRepository) GetAllJoinRoom(userId int) ([]*models.ChatRoom, error) {
	query := `SELECT groupMembers.id, groupMembers.roomName FROM groupMembers JOIN chatRoom ON groupMembers.roomName = chatRoom.name WHERE groupMembers.userId = $1;
`

	rows, err := r.db.Query(query, userId)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var data []*models.ChatRoom
	for rows.Next() {
		room := &models.ChatRoom{}
		err := rows.Scan(&room.Id, &room.Name, &room.UserId, &room.ProfilePic)
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
