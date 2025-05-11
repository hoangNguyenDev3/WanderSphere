package auth

import (
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt"
)

const (
	DefaultTokenLifespanHours = 24                        // Default token lifespan is 24 hours
	DefaultAPISecret          = "wandersphere_secret_key" // Default API secret
)

// NOTE: The JWT implementation is kept for backward compatibility.
// The preferred authentication approach is session-based authentication using Redis.
// New code should use the session-based authentication mechanism.

// GenerateToken creates a new JWT token for the given user ID
// Deprecated: Use session-based authentication with Redis instead.
func GenerateToken(user_id uint, secretKey string, tokenLifespanHours int) (string, error) {
	// Use provided values or defaults
	if secretKey == "" {
		secretKey = DefaultAPISecret
	}

	if tokenLifespanHours <= 0 {
		tokenLifespanHours = DefaultTokenLifespanHours
	}

	claims := jwt.MapClaims{}
	claims["authorized"] = true
	claims["user_id"] = user_id
	claims["exp"] = time.Now().Add(time.Hour * time.Duration(tokenLifespanHours)).Unix()
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	return token.SignedString([]byte(secretKey))
}

// TokenValid validates the JWT token from the request
// Deprecated: Use session-based authentication with Redis instead.
func TokenValid(c *gin.Context, secretKey string) error {
	if secretKey == "" {
		secretKey = DefaultAPISecret
	}

	tokenString := ExtractToken(c)
	_, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(secretKey), nil
	})
	if err != nil {
		return err
	}
	return nil
}

// ExtractToken extracts the JWT token from the request
// Deprecated: Use session-based authentication with Redis instead.
func ExtractToken(c *gin.Context) string {
	token := c.Query("token")
	if token != "" {
		return token
	}
	bearerToken := c.Request.Header.Get("Authorization")
	if len(strings.Split(bearerToken, " ")) == 2 {
		return strings.Split(bearerToken, " ")[1]
	}
	return ""
}

// ExtractTokenID extracts the user ID from the JWT token
// Deprecated: Use session-based authentication with Redis instead.
func ExtractTokenID(c *gin.Context, secretKey string) (uint, error) {
	if secretKey == "" {
		secretKey = DefaultAPISecret
	}

	tokenString := ExtractToken(c)
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(secretKey), nil
	})
	if err != nil {
		return 0, err
	}
	claims, ok := token.Claims.(jwt.MapClaims)
	if ok && token.Valid {
		uid, err := strconv.ParseUint(fmt.Sprintf("%.0f", claims["user_id"]), 10, 32)
		if err != nil {
			return 0, err
		}
		return uint(uid), nil
	}
	return 0, nil
}
