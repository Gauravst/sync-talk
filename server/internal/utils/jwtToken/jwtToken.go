package jwtToken

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

func VerifyJwtAndGetData[T any](jwtToken string, key string) (*T, error) {
	token, err := jwt.Parse(jwtToken, func(token *jwt.Token) (interface{}, error) {
		// Ensure the token uses the correct signing method
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		// Convert key (string) to []byte for HS256
		return []byte(key), nil
	})

	// Extract claims regardless of token validity
	var claims jwt.MapClaims
	if token != nil {
		if data, ok := token.Claims.(jwt.MapClaims); ok {
			claims = data
		}
	}

	// Convert claims into JSON and then unmarshal into struct T
	jsonData, _ := json.Marshal(claims)
	var result T
	errUnmarshal := json.Unmarshal(jsonData, &result)

	// Check if there's an error, but still return claims if possible
	if err != nil {
		// If token expired, return claims with specific error
		if errors.Is(err, jwt.ErrSignatureInvalid) {
			return nil, fmt.Errorf("invalid token signature")
		}
		if claims != nil {
			return &result, fmt.Errorf("token has expired")
		}
		return nil, err
	}

	// Handle potential JSON unmarshal errors separately
	if errUnmarshal != nil {
		return nil, errUnmarshal
	}

	return &result, nil
}

func CreateNewToken(data interface{}, key string) (string, error) {
	claims, ok := data.(jwt.MapClaims)
	if !ok {
		return "", errors.New("invalid claims type")
	}

	token := jwt.NewWithClaims(jwt.SigningMethodES256, claims)
	tokenString, err := token.SignedString(key)

	if err != nil {
		return "", err
	}

	return tokenString, nil
}

func SetAccessToken(w http.ResponseWriter, token string, secure bool) {
	http.SetCookie(w, &http.Cookie{
		Name:     "accessToken",
		Value:    token,
		Path:     "/",
		HttpOnly: true,
		Secure:   secure,
		SameSite: http.SameSiteLaxMode,
	})
}

func RemoveAccessToken(w http.ResponseWriter, secure bool) {
	http.SetCookie(w, &http.Cookie{
		Name:     "accessToken",
		Value:    "",
		Path:     "/",
		HttpOnly: true,
		Secure:   secure, // Set based on your environment
		SameSite: http.SameSiteLaxMode,
		MaxAge:   -1, // Immediately expire the cookie
		Expires:  time.Unix(0, 0),
	})
}
