package withoutauth

import (
	crand "crypto/rand" // Alias crypto/rand to crand
	"math/big"
	"math/rand" // Keep math/rand as rand
	"time"
)

// GenerateUsername creates a random username with a prefix
func GenerateUsername(prefix string, length int) (string, error) {
	const charset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	r := rand.New(rand.NewSource(time.Now().UnixNano())) // Use math/rand

	username := make([]byte, length)
	for i := range username {
		username[i] = charset[r.Intn(len(charset))]
	}

	return prefix + string(username), nil
}

// GeneratePassword creates a secure random password
func GeneratePassword(length int) (string, error) {
	const charset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789!@#$%^&*()-_=+"
	password := make([]byte, length)

	for i := range password {
		randomIndex, err := crand.Int(crand.Reader, big.NewInt(int64(len(charset)))) // Use crand (crypto/rand)
		if err != nil {
			return "", err
		}
		password[i] = charset[randomIndex.Int64()]
	}

	return string(password), nil
}
