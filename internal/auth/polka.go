package auth

import (
	"errors"
	"net/http"
	"strings"
)

func GetAPIKey(headers http.Header) (string, error) {
	authHeader := headers.Get("Authorization")
	if authHeader == "" {
		return "", errors.New("missing authorization header")
	}

	// Expected: "ApiKey <key>"
	const prefix = "ApiKey "

	if !strings.HasPrefix(authHeader, prefix) {
		return "", errors.New("malformed authorization header")
	}

	// Remove prefix and trim whitespace
	apiKey := strings.TrimSpace(strings.TrimPrefix(authHeader, prefix))
	if apiKey == "" {
		return "", errors.New("api key missing")
	}

	return apiKey, nil
}
