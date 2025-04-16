package services

import (
	"time"

	"github.com/gauravst/real-time-chat/internal/config"
	"github.com/gauravst/real-time-chat/internal/models"
	"github.com/gauravst/real-time-chat/internal/repositories"
	"github.com/gauravst/real-time-chat/internal/utils/hashing"
	withoutauth "github.com/gauravst/real-time-chat/internal/utils/withoutAuth"
	"github.com/golang-jwt/jwt/v5"
)

type AuthService interface {
	LoginUser(data *models.LoginRequest, cfg config.Config) (string, *models.User, error)
	LoginWithoutAuth(cfg config.Config) (string, *models.User, error)
	// RefreshToken(userId int, token string) error
	GetRefreshToken(userId int) (string, error)
	LogoutUser(userId int) error
}

type authService struct {
	authRepo repositories.AuthRepository
}

func NewAuthService(authRepo repositories.AuthRepository) AuthService {
	return &authService{
		authRepo: authRepo,
	}
}

func (s *authService) LoginUser(data *models.LoginRequest, cfg config.Config) (string, *models.User, error) {
	// check user exsit
	userData, err := s.authRepo.CheckUserByUsername(data.Username)
	var userId int
	var returnUserData *models.User

	// some err so return error
	if err != nil && err.Error() != "user not found" {
		return "", nil, err
	}

	// user not found create new user
	if err != nil {
		hashedPassword, err := hashing.GenerateHashString(data.Password)
		if err != nil {
			return "", nil, err
		}

		data.Password = hashedPassword
		// create new user
		createdUserData, err := s.authRepo.CreateNewUser(data)
		if err != nil {
			return "", nil, err
		}

		// seting userid from newUser
		userId = createdUserData.Id
		returnUserData = createdUserData
	}

	returnUserData = &userData

	// if user exist than
	if userData.Password != "" {
		// check user password and data password is same or not
		err = hashing.CompareHashString(userData.Password, data.Password)
		if err != nil {
			return "", nil, err
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
		return "", returnUserData, err
	}

	// remove all login here first
	err = s.authRepo.RemoveOtherLogin(userId)
	if err != nil {
		return "", nil, err
	}

	var loginData models.LoginSession
	loginData.Token = refreshTokenString
	loginData.UserId = userId

	// user exist create login
	err = s.authRepo.LoginUser(&loginData)
	if err != nil {
		return "", nil, err
	}

	// new accessToken for new user and login user
	claims1 := jwt.MapClaims{
		"userId":   userId,
		"username": data.Username,
		"role":     role,
		"exp":      time.Now().Add(30 * time.Minute).Unix(),
	}
	accessToken := jwt.NewWithClaims(jwt.SigningMethodHS256, claims1)
	accessTokenString, err := accessToken.SignedString([]byte(cfg.JwtPrivateKey))
	if err != nil {
		return "", nil, err
	}

	return accessTokenString, returnUserData, nil
}

func (s *authService) LoginWithoutAuth(cfg config.Config) (string, *models.User, error) {
	// create new username
	var returnUserData *models.User
	username, err := withoutauth.GenerateUsername("user_", 6)
	if err != nil {
		return "", nil, err
	}

	// create new password and hash
	password, err := withoutauth.GeneratePassword(12)
	if err != nil {
		return "", nil, err
	}

	hashedPassword, err := hashing.GenerateHashString(password)
	if err != nil {
		return "", nil, err
	}

	// create user here
	data := &models.LoginRequest{
		Username: username,
		Password: hashedPassword,
	}
	createdUserData, err := s.authRepo.CreateNewUser(data)
	if err != nil {
		return "", nil, err
	}

	// create accessToken & refreshToken here
	claims2 := jwt.MapClaims{
		"userId":   createdUserData.Id,
		"username": createdUserData.Username,
		"role":     createdUserData.Role,
		"exp":      time.Now().Add(24 * 30 * time.Hour).Unix(),
	}
	refreshToken := jwt.NewWithClaims(jwt.SigningMethodHS256, claims2)
	refreshTokenString, err := refreshToken.SignedString([]byte(cfg.JwtPrivateKey))
	if err != nil {
		return "", nil, err
	}

	returnUserData = createdUserData
	claims1 := jwt.MapClaims{
		"userId":   createdUserData.Id,
		"username": createdUserData.Username,
		"role":     createdUserData.Role,
		"exp":      time.Now().Add(30 * time.Minute).Unix(),
	}
	accessToken := jwt.NewWithClaims(jwt.SigningMethodHS256, claims1)
	accessTokenString, err := accessToken.SignedString([]byte(cfg.JwtPrivateKey))
	if err != nil {
		return "", nil, err
	}

	// create a login here with user info
	var loginData models.LoginSession
	loginData.Token = refreshTokenString
	loginData.UserId = createdUserData.Id

	err = s.authRepo.LoginUser(&loginData)
	if err != nil {
		return "", nil, err
	}

	// send back accessToken
	return accessTokenString, returnUserData, nil
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

func (s *authService) LogoutUser(userId int) error {
	err := s.authRepo.LogoutUser(userId)
	if err != nil {
		return err
	}

	return nil
}
