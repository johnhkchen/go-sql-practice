package routes

import (
	"crypto/rand"
	"crypto/subtle"
	"encoding/hex"
)

// TokenLength defines the byte length of tokens before encoding
const TokenLength = 32

// GenerateToken generates a cryptographically secure random token
func GenerateToken() (string, error) {
	bytes := make([]byte, TokenLength)
	if _, err := rand.Read(bytes); err != nil {
		return "", err
	}
	return hex.EncodeToString(bytes), nil
}

// ValidateToken performs constant-time comparison of tokens
func ValidateToken(provided, stored string) bool {
	if len(provided) != len(stored) {
		return false
	}
	return subtle.ConstantTimeCompare([]byte(provided), []byte(stored)) == 1
}