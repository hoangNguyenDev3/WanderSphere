package auth

import (
	"crypto/rand"
	"strings"

	"golang.org/x/crypto/bcrypt"
)

func HashPassword(password string, salt []byte) (string, error) {
	// using bcrypt to hash the password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}
	return string(hashedPassword), nil
}

func CheckPasswordHash(hashedPassword string, password string) error {
	hashedPasswordBytes := []byte(hashedPassword)

	passwordBytes := []byte(password)

	return bcrypt.CompareHashAndPassword(hashedPasswordBytes, passwordBytes)
}

func SanitizePassword(password string) string {
	password = strings.TrimSpace(password)
	return password
}

func GenerateRandomSalt() ([]byte, error) {
	salt := make([]byte, 16)
	_, err := rand.Read(salt)
	if err != nil {
		return nil, err
	}
	return salt, nil
}
