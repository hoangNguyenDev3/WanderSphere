package auth

import (
	"crypto/rand"
	"html"
	"strings"

	"golang.org/x/crypto/bcrypt"
)

func HashPassword(password string, salt []byte) (string, error) {
	// Convert password string to byte slice
	var passwordBytes = []byte(password)

	// Append salt to password
	passwordBytes = append(passwordBytes, salt...)

	// Get the bcrypt hashed password
	hashedPasswordBytes, err := bcrypt.GenerateFromPassword(passwordBytes, 4)
	if err != nil {
		return "", err
	}

	return string(hashedPasswordBytes), err
}

func CheckPasswordHash(hashedPassword, password string, salt []byte) error {
	// Convert password string to byte slice
	var passwordBytes = []byte(password)

	// Append salt to password
	passwordBytes = append(passwordBytes, salt...)

	return bcrypt.CompareHashAndPassword([]byte(hashedPassword), passwordBytes)
}

func Santinize(data string) string {
	data = html.EscapeString(strings.TrimSpace(data))
	return data
}

func GenerateRandomSalt() ([]byte, error) {
	salt := make([]byte, 4)

	_, err := rand.Read(salt[:])
	if err != nil {
		return []byte("error"), err
	}

	return salt, nil
}
