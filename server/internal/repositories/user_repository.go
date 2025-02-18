package repositories

import (
	"database/sql"

	"github.com/gauravst/real-time-chat/internal/models"
)

// UserRepository defines the interface for user-related database operations
type UserRepository interface {
	CreateUser(user *models.User) error
	GetAllUsers() ([]*models.User, error)
	GetUserByID(id int) (*models.User, error)
	UpdateUser(user *models.UserRequest) error
	DeleteUser(id int) error
}

// userRepository implements the UserRepository interface
type userRepository struct {
	db *sql.DB
}

// NewUserRepository creates a new instance of userRepository
func NewUserRepository(db *sql.DB) UserRepository {
	return &userRepository{
		db: db,
	}
}

// CreateUser inserts a new user into the database
func (r *userRepository) CreateUser(user *models.User) error {
	query := `INSERT INTO users (username, password) VALUES ($1, $2) RETURNING id`
	err := r.db.QueryRow(query, user.Username, user.Password).Scan(&user.Id)
	if err != nil {
		return err
	}

	return nil
}

// get all user
func (r *userRepository) GetAllUsers() ([]*models.User, error) {
	query := `SELECT id, username, role, password FROM users`

	rows, err := r.db.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var data []*models.User
	for rows.Next() {
		user := &models.User{}
		err := rows.Scan(&user.Id, &user.Username, &user.Role, &user.Password)
		if err != nil {
			return nil, err
		}

		data = append(data, user)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return data, nil
}

// GetUserByID retrieves a user by their ID from the database
func (r *userRepository) GetUserByID(id int) (*models.User, error) {
	user := &models.User{}
	query := `SELECT id, username, role, password FROM users WHERE id = $1`
	err := r.db.QueryRow(query, id).Scan(&user.Id, &user.Username, &user.Role, &user.Password)
	if err != nil {
		return nil, err
	}

	return user, nil
}

// UpdateUser updates an existing user in the database
func (r *userRepository) UpdateUser(user *models.UserRequest) error {
	query := `UPDATE users SET username = $1, role = $2, password = $3 WHERE id = $4`
	_, err := r.db.Exec(query, user.Username, user.Role, user.Password, user.Id)
	if err != nil {
		return err
	}

	return nil
}

// DeleteUser deletes a user by their ID from the database
func (r *userRepository) DeleteUser(id int) error {
	query := `DELETE FROM users WHERE id = $1`
	_, err := r.db.Exec(query, id)
	if err != nil {
		return err
	}

	return nil
}
