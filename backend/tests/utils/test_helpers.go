package utils

import (
	"crypto/rand"
	"fmt"
	"math/big"
	"time"
)

// GenerateUniqueUsername creates a highly unique username for testing that fits database constraints
func GenerateUniqueUsername(prefix string) string {
	// Use shorter components to fit within 50 character limit
	randomNum, _ := rand.Int(rand.Reader, big.NewInt(99999))
	timestamp := time.Now().Unix() % 1000000 // Use last 6 digits of timestamp

	// Generate a short random string
	randomBytes := make([]byte, 3)
	rand.Read(randomBytes)

	// Create username that stays under 50 characters
	// Format: prefix_timestamp_random_hex (should be under 30 chars)
	username := fmt.Sprintf("%s_%d_%d_%x", prefix, timestamp, randomNum.Int64(), randomBytes)

	// Ensure it's under 50 characters
	if len(username) > 45 {
		// Truncate prefix if needed and use shorter format
		shortPrefix := prefix
		if len(prefix) > 8 {
			shortPrefix = prefix[:8]
		}
		username = fmt.Sprintf("%s_%d_%x", shortPrefix, timestamp, randomBytes)
	}

	return username
}
