package auth_test

import (
	"testing"

	"github.com/SaadVSP96/Chirpy_Server.git/internal/auth"
)

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
