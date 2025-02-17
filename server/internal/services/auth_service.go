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

	// some err so return error
	if err != nil && err.Error() != "user not found" {
		return "", err
	}

	// new accessToken for new user and login user
	claims := jwt.MapClaims{
		"username": data.Username,
		"exp":      time.Now().Add(15 * time.Minute).Unix(),
	}
	accessToken := jwt.NewWithClaims(jwt.SigningMethodES256, claims)
	accessTokenString, err := accessToken.SignedString(cfg.JwtPrivateKey)
	data.AccessToken = accessTokenString

	// RefreshToken
	claims = jwt.MapClaims{
		"username": data.Username,
		"exp":      time.Now().Add(24 * 30 * time.Hour).Unix(),
	}
	refreshToken := jwt.NewWithClaims(jwt.SigningMethodES256, claims)
	refreshTokenString, err := refreshToken.SignedString(cfg.JwtPrivateKey)

	// user not found create new user
	if err != nil {

		hashedPassword, err := hashing.GenerateHashString(data.Password)
		if err != nil {
			return "", err
		}

		data.Password = hashedPassword
		// create new user
		err = s.authRepo.CreateNewUser(data)
		if err != nil {
			return "", err
		}
	}

	// check user password and data password is same or not
	err = hashing.CompareHashString(userData.Password, data.Password)
	if err != nil {
		return "", err
	}

	var loginData models.LoginSession
	loginData.Token = refreshTokenString
	loginData.UserId = userData.Id

	// user exist create login
	err = s.authRepo.LoginUser(&loginData)
	if err != nil {
		return "", err
	}

	return data.AccessToken, nil
}
