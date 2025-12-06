package auth

import (
	"errors"
	"net/http"
	"strings"

	"github.com/alexedwards/argon2id"
)

// HashPassword generates an Argon2id hash from a plaintext password.
func HashPassword(password string) (string, error) {
	hash, err := argon2id.CreateHash(password, argon2id.DefaultParams)
	if err != nil {
		return "", err
	}
	return hash, nil
}

// CheckPasswordHash compares a plaintext password with a stored hash.
func CheckPasswordHash(password, hash string) (bool, error) {
	match, err := argon2id.ComparePasswordAndHash(password, hash)
	if err != nil {
		return false, err
	}
	return match, nil
}

// GetBearerToken function which extracts bearer token from incoming request
func GetBearerToken(headers http.Header) (string, error) {
	// look for authorization Header
	authHeader := headers.Get("Authorization")
	if authHeader == "" {
		return "", errors.New("no authorization header included")
	}

	// Must Begin with Bearer
	if !strings.HasPrefix(authHeader, "Bearer ") {
		return "", errors.New("malformed authorization header")
	}

	// Extract token portion
	token := strings.TrimSpace(strings.TrimPrefix(authHeader, "Bearer "))
	if token == "" {
		return "", errors.New("empty bearer token")
	}
	return token, nil
}
