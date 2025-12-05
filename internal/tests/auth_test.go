package auth_test

import (
	"net/http"
	"testing"
	"time"

	"github.com/SaadVSP96/Chirpy_Server.git/internal/auth"
	"github.com/google/uuid"
)

// auth.
func TestPasswordHashing(t *testing.T) {
	pw := "supersecret123"

	hash, err := auth.HashPassword(pw)
	if err != nil {
		t.Fatalf("hashing failed: %v", err)
	}

	ok, err := auth.CheckPasswordHash(pw, hash)
	if err != nil {
		t.Fatalf("comparison failed: %v", err)
	}

	if !ok {
		t.Fatalf("password should match hash but does not")
	}
}

func TestJWTLifecycle(t *testing.T) {
	secret := "supersecret"
	userID := uuid.New()

	token, err := auth.MakeJWT(userID, secret, time.Hour)
	if err != nil {
		t.Fatalf("failed to make jwt: %v", err)
	}

	gotID, err := auth.ValidateJWT(token, secret)
	if err != nil {
		t.Fatalf("failed to validate jwt: %v", err)
	}

	if gotID != userID {
		t.Fatalf("expected %v got %v", userID, gotID)
	}
}

func TestJWTExpired(t *testing.T) {
	secret := "abc123"
	userID := uuid.New()

	token, err := auth.MakeJWT(userID, secret, -time.Hour) // already expired
	if err != nil {
		t.Fatalf("make jwt failed: %v", err)
	}

	_, err = auth.ValidateJWT(token, secret)
	if err == nil {
		t.Fatalf("expected error for expired token, got none")
	}
}

func TestJWTWrongSecret(t *testing.T) {
	userID := uuid.New()

	token, err := auth.MakeJWT(userID, "correctSecret", time.Hour)
	if err != nil {
		t.Fatalf("make jwt failed: %v", err)
	}

	_, err = auth.ValidateJWT(token, "wrongSecret")
	if err == nil {
		t.Fatalf("expected signature error but got none")
	}
}

func TestGetBearerToken(t *testing.T) {
	headers := http.Header{}
	headers.Set("Authorization", "Bearer abc123")

	token, err := auth.GetBearerToken(headers)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if token != "abc123" {
		t.Fatalf("expected token 'abc123', got '%s'", token)
	}
}

func TestGetBearerToken_NoHeader(t *testing.T) {
	headers := http.Header{}
	_, err := auth.GetBearerToken(headers)
	if err == nil {
		t.Fatalf("expected error, got nil")
	}
}

func TestGetBearerToken_BadFormat(t *testing.T) {
	headers := http.Header{}
	headers.Set("Authorization", "Token xyz")

	_, err := auth.GetBearerToken(headers)
	if err == nil {
		t.Fatalf("expected error for malformed header")
	}
}
