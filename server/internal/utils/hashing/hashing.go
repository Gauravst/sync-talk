package hashing

import (
	"fmt"

	"golang.org/x/crypto/bcrypt"
)

func GenerateHashString(data string) (string, error) {
	hashedValue, err := bcrypt.GenerateFromPassword([]byte(data), bcrypt.DefaultCost)
	if err != nil {
		return "", nil
	}
	return string(hashedValue), nil
}

func CompareHashString(hashedValue string, normalValue string) error {
	err := bcrypt.CompareHashAndPassword([]byte(hashedValue), []byte(normalValue))
	if err != nil {
		return fmt.Errorf("invalid credentials")
	}
	return nil
}
