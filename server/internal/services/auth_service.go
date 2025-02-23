package services

import (
	"time"

	"github.com/gauravst/real-time-chat/internal/config"
	"github.com/gauravst/real-time-chat/internal/models"
	"github.com/gauravst/real-time-chat/internal/repositories"
	"github.com/gauravst/real-time-chat/internal/utils/hashing"
	"github.com/golang-jwt/jwt/v5"
)

type AuthService interface {
	LoginUser(data *models.LoginRequest, cfg config.Config) (string, error)
	// RefreshToken(userId int, token string) error
	GetRefreshToken(userId int) (string, error)
}

type authService struct {
	authRepo repositories.AuthRepository
}

func NewAuthService(authRepo repositories.AuthRepository) AuthService {
	return &authService{
		authRepo: authRepo,
	}
}

func (s *authService) LoginUser(data *models.LoginRequest, cfg config.Config) (string, error) {
	// check user exsit
	userData, err := s.authRepo.CheckUserByUsername(data.Username)
	var userId int

	// some err so return error
	if err != nil && err.Error() != "user not found" {
		return "", err
	}

	// user not found create new user
	if err != nil {
		hashedPassword, err := hashing.GenerateHashString(data.Password)
		if err != nil {
			return "", err
		}

		data.Password = hashedPassword
		// create new user
		createdUserData, err := s.authRepo.CreateNewUser(data)
		if err != nil {
			return "", err
		}

		// seting userid from newUser
		userId = createdUserData.Id
	}

	// if user exist than
	if userData.Password != "" {
		// check user password and data password is same or not
		err = hashing.CompareHashString(userData.Password, data.Password)
		if err != nil {
			return "", err
		}

		// seting userId here from userData
		userId = userData.Id
	}

	role := userData.Role
	if role == "" {
		role = "USER"
	}

	// RefreshToken
	claims2 := jwt.MapClaims{
		"userId":   userId,
		"username": data.Username,
		"role":     role,
		"exp":      time.Now().Add(24 * 30 * time.Hour).Unix(),
	}
	refreshToken := jwt.NewWithClaims(jwt.SigningMethodHS256, claims2)
	refreshTokenString, err := refreshToken.SignedString([]byte(cfg.JwtPrivateKey))
	if err != nil {
		return "", err
	}

	// remove all login here first
	err = s.authRepo.RemoveOtherLogin(userId)
	if err != nil {
		return "", err
	}

	var loginData models.LoginSession
	loginData.Token = refreshTokenString
	loginData.UserId = userId

	// user exist create login
	err = s.authRepo.LoginUser(&loginData)
	if err != nil {
		return "", err
	}

	// new accessToken for new user and login user
	claims1 := jwt.MapClaims{
		"userId":   userId,
		"username": data.Username,
		"role":     role,
		"exp":      time.Now().Add(2 * time.Minute).Unix(),
	}
	accessToken := jwt.NewWithClaims(jwt.SigningMethodHS256, claims1)
	accessTokenString, err := accessToken.SignedString([]byte(cfg.JwtPrivateKey))
	if err != nil {
		return "", err
	}

	return accessTokenString, nil
}

// func (s *authService) RefreshToken(userId int, token string) error {
// }

func (s *authService) GetRefreshToken(userId int) (string, error) {
	token, err := s.authRepo.GetRefreshToken(userId)
	if err != nil {
		return "", err
	}

	return token, nil
}
