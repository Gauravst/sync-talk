package repositories

import (
	"database/sql"

	"github.com/gauravst/real-time-chat/internal/models"
)

// UserRepository defines the interface for user-related database operations
type UserRepository interface {
	CreateUser(user *models.User) error
	GetUserByID(id int) (*models.User, error)
	UpdateUser(user *models.User) error
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
	query := `INSERT INTO users (name, email) VALUES ($1, $2) RETURNING id`
	err := r.db.QueryRow(query, user.Name, user.Email).Scan(&user.ID)
	if err != nil {
		return err
	}

	return nil
}

// GetUserByID retrieves a user by their ID from the database
func (r *userRepository) GetUserByID(id int) (*models.User, error) {
	user := &models.User{}
	query := `SELECT id, name, email FROM users WHERE id = $1`
	err := r.db.QueryRow(query, id).Scan(&user.ID, &user.Name, &user.Email)
	if err != nil {
		return nil, err
	}

	return user, nil
}

// UpdateUser updates an existing user in the database
func (r *userRepository) UpdateUser(user *models.User) error {
	query := `UPDATE users SET name = $1, email = $2 WHERE id = $3`
	_, err := r.db.Exec(query, user.Name, user.Email, user.ID)
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
