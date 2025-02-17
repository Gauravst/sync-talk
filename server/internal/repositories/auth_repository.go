package repositories

import (
	"database/sql"
	"fmt"

	"github.com/gauravst/real-time-chat/internal/models"
)

// AuthRepository defines the interface for user-related database operations
type AuthRepository interface {
	LoginUser(data *models.LoginSession) error
	CreateNewUser(data *models.LoginRequest) error
	CheckUserByUsername(username string) (models.User, error)
}

// userRepository implements the AuthRepository interface
type authRepository struct {
	db *sql.DB
}

// NewAuthRepository creates a new instance of userRepository
func NewAuthRepository(db *sql.DB) AuthRepository {
	return &authRepository{
		db: db,
	}
}

func (r *authRepository) LoginUser(data *models.LoginSession) error {
	query := `INSERT INTO login_session (userId, token) VALUES ($1, $2)`
	_, err := r.db.Exec(query, data.UserId, data.Token)
	if err != nil {
		return err
	}
	return nil
}

func (r *authRepository) CreateNewUser(data *models.LoginRequest) error {
	query := `INSERT INTO users (username, password) VALUES ($1, $2)`
	_, err := r.db.Exec(query, data.Username, data.Password)
	if err != nil {
		return err
	}
	return err
}

func (r *authRepository) CheckUserByUsername(username string) (models.User, error) {
	var user models.User
	query := `SELECT id, username, password FROM users WHERE username = $1`
	err := r.db.QueryRow(query, username).Scan(&user.Id, &user.Username, &user.Password)
	if err != nil {
		if err == sql.ErrNoRows {
			return user, fmt.Errorf("user not found")
		}
		return user, err
	}
	return user, nil
}
