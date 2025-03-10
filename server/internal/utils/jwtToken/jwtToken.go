package jwtToken

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strings"
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

func SetAccessToken(w http.ResponseWriter, r *http.Request, token string, secure bool) {
	isLocal := isLocalRequest(r)

	http.SetCookie(w, &http.Cookie{
		Name:     "accessToken",
		Value:    token,
		Path:     "/",
		HttpOnly: true,
		Secure:   getSecureSetting(isLocal, secure),
		SameSite: getSameSiteMode(isLocal),
	})
}

// RemoveAccessToken removes the accessToken cookie
func RemoveAccessToken(w http.ResponseWriter, r *http.Request, secure bool) {
	isLocal := isLocalRequest(r)

	http.SetCookie(w, &http.Cookie{
		Name:     "accessToken",
		Value:    "",
		Path:     "/",
		HttpOnly: true,
		Secure:   getSecureSetting(isLocal, secure),
		SameSite: getSameSiteMode(isLocal),
		MaxAge:   -1, // Expire immediately
		Expires:  time.Unix(0, 0),
	})
}

// isLocalRequest checks if the request is from localhost
func isLocalRequest(r *http.Request) bool {
	host := r.Host
	forwarded := r.Header.Get("X-Forwarded-For")

	// Check if host is localhost or 127.0.0.1
	if strings.Contains(host, "localhost") || strings.HasPrefix(host, "127.0.0.1") {
		return true
	}

	// Handle proxies (e.g., Docker, Nginx, ngrok)
	if forwarded == "127.0.0.1" {
		return true
	}

	return false
}

// getSameSiteMode returns the correct SameSite mode based on the environment
func getSameSiteMode(isLocal bool) http.SameSite {
	if isLocal {
		return http.SameSiteLaxMode // Allow local dev
	}
	return http.SameSiteNoneMode // Required for cross-origin cookies in production
}

// getSecureSetting ensures Secure=true in production but respects param in localhost
func getSecureSetting(isLocal bool, secure bool) bool {
	if isLocal {
		return secure // Use whatever secure setting is passed in local
	}
	return true // Force Secure=true in production
}
