package repositories

import (
	"database/sql"
	"fmt"

	"github.com/gauravst/real-time-chat/internal/models"
)

// AuthRepository defines the interface for user-related database operations
type AuthRepository interface {
	RemoveOtherLogin(userId int) error
	LoginUser(data *models.LoginSession) error
	CreateNewUser(data *models.LoginRequest) (*models.User, error)
	CheckUserByUsername(username string) (models.User, error)
	// RefreshToken(userId int, token string) error
	GetRefreshToken(userId int) (string, error)
	LogoutUser(userId int) error
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

func (r *authRepository) RemoveOtherLogin(userId int) error {
	query := `DELETE FROM loginSession WHERE userId = $1`
	_, err := r.db.Exec(query, userId)
	if err != nil {
		return err
	}

	return nil
}

func (r *authRepository) LoginUser(data *models.LoginSession) error {
	query := `INSERT INTO loginSession (userId, token) VALUES ($1, $2)`
	_, err := r.db.Exec(query, data.UserId, data.Token)
	if err != nil {
		return err
	}
	return nil
}

func (r *authRepository) CreateNewUser(data *models.LoginRequest) (*models.User, error) {
	user := &models.User{}
	query := `INSERT INTO users (username, password) VALUES ($1, $2) RETURNING id, username, role, createdAt`
	err := r.db.QueryRow(query, data.Username, data.Password).Scan(&user.Id, &user.Username, &user.Role, &user.CreatedAt)
	if err != nil {
		return nil, err
	}

	return user, nil
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

// func (r *authRepository) RefreshToken(userId int, token string) error {
// 	query := `UPDATE`
// 	_, err := r.db.Exec(query, token, userId)
// 	if err != nil {
// 		return err
// 	}
// 	return nil
// }

func (r *authRepository) GetRefreshToken(userId int) (string, error) {
	var token string
	query := `SELECT token FROM loginSession WHERE userId = $1`
	err := r.db.QueryRow(query, userId).Scan(&token)
	if err != nil {
		return token, err
	}
	return token, nil
}

func (r *authRepository) LogoutUser(userId int) error {
	query := `DELETE FROM loginSession WHERE userId = $1`
	_, err := r.db.Exec(query, userId)
	if err != nil {
		return err
	}
	return nil
}
