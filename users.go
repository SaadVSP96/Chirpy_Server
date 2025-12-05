package main

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/SaadVSP96/Chirpy_Server.git/internal/auth"
	"github.com/SaadVSP96/Chirpy_Server.git/internal/database"
)

type CreateUserRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type CreateUserResponse struct {
	ID        string `json:"id"`
	CreatedAt string `json:"created_at"`
	UpdatedAt string `json:"updated_at"`
	Email     string `json:"email"`
}

type LoginRequest struct {
	Email            string `json:"email"`
	Password         string `json:"password"`
	ExpiresInSeconds *int64 `json:"expires_in_seconds"` // pointer so it's optional
}

type LoginResponse struct {
	ID           string `json:"id"`
	CreatedAt    string `json:"created_at"`
	UpdatedAt    string `json:"updated_at"`
	Email        string `json:"email"`
	Token        string `json:"token"`
	RefreshToken string `json:"refresh_token"`
}

func (cfg *apiConfig) handlerCreateUser(w http.ResponseWriter, r *http.Request) {
	var req CreateUserRequest

	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(&req)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid JSON", err)
		return
	}

	if req.Email == "" || req.Password == "" {
		respondWithError(w, http.StatusBadRequest, "Email and password are required", nil)
		return
	}
	// hash the password
	hashedPassword, err := auth.HashPassword(req.Password)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Failed to hash password", err)
		return
	}
	// run the SQLC generated query function
	user, err := cfg.dbQueries.CreateUser(
		r.Context(),
		database.CreateUserParams{
			Email:          req.Email,
			HashedPassword: hashedPassword,
		},
	)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Could not create user", err)
		return
	}

	// build a response to send back in the response to post request
	respondWithJSON(w, http.StatusCreated, CreateUserResponse{
		ID:        user.ID.String(),
		CreatedAt: user.CreatedAt.Format(time.RFC3339),
		UpdatedAt: user.UpdatedAt.Format(time.RFC3339),
		Email:     user.Email,
	})
}

func (cfg *apiConfig) handlerLogin(w http.ResponseWriter, r *http.Request) {
	var req LoginRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid JSON", err)
		return
	}

	if req.Email == "" || req.Password == "" {
		respondWithError(w, http.StatusBadRequest, "Email and password are required", nil)
		return
	}
	// Look up the user (using your new SQLC query)
	user, err := cfg.dbQueries.GetUserByEmail(r.Context(), req.Email)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "Incorrect email or password", nil)
		return
	}

	// Compare password
	ok, _ := auth.CheckPasswordHash(req.Password, user.HashedPassword)
	if !ok {
		respondWithError(w, http.StatusUnauthorized, "Incorrect email or password", nil)
		return
	}

	// Create access token (JWT)
	accessToken, err := auth.MakeJWT(user.ID, cfg.jwtSecret, time.Hour)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Failed to create JWT", err)
		return
	}

	// Create refresh token (random string)
	refreshToken, err := auth.MakeRefreshToken()
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Failed to create refresh token", err)
		return
	}

	// Save refresh token to DB
	expiresAt := time.Now().UTC().Add(60 * 24 * time.Hour) // 60 days
	_, err = cfg.dbQueries.CreateRefreshToken(r.Context(), database.CreateRefreshTokenParams{
		Token:     refreshToken,
		UserID:    user.ID,
		ExpiresAt: expiresAt,
	})
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Could not save refresh token", err)
		return
	}

	// Build response
	resp := LoginResponse{
		ID:           user.ID.String(),
		CreatedAt:    user.CreatedAt.Format(time.RFC3339),
		UpdatedAt:    user.UpdatedAt.Format(time.RFC3339),
		Email:        user.Email,
		Token:        accessToken,
		RefreshToken: refreshToken,
	}

	respondWithJSON(w, http.StatusOK, resp)
}

func (cfg *apiConfig) handlerRefresh(w http.ResponseWriter, r *http.Request) {
	// Try to get refresh token from Authorization header first (as Bearer token)
	refreshTokenString, err := auth.GetBearerToken(r.Header)
	if err != nil {
		// If not in header, try to get from JSON body
		type RefreshRequest struct {
			RefreshToken string `json:"refresh_token"`
		}

		var req RefreshRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			respondWithError(w, http.StatusBadRequest, "Invalid JSON", err)
			return
		}

		if req.RefreshToken == "" {
			respondWithError(w, http.StatusUnauthorized, "Missing refresh token", nil)
			return
		}
		refreshTokenString = req.RefreshToken
	}

	if refreshTokenString == "" {
		respondWithError(w, http.StatusUnauthorized, "Missing refresh token", nil)
		return
	}

	// Get refresh token from database
	refreshToken, err := cfg.dbQueries.GetRefreshToken(r.Context(), refreshTokenString)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "Invalid token", nil)
		return
	}

	if time.Now().UTC().After(refreshToken.ExpiresAt) {
		respondWithError(w, http.StatusUnauthorized, "Expired token", nil)
		return
	}
	if refreshToken.RevokedAt.Valid && time.Now().UTC().After(refreshToken.RevokedAt.Time) {
		respondWithError(w, http.StatusUnauthorized, "Revoked token", nil)
		return
	}

	user, err := cfg.dbQueries.GetUserFromRefreshToken(r.Context(), refreshToken.Token)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "Invalid token", nil)
		return
	}

	// Create new access token
	newAccessToken, err := auth.MakeJWT(user.ID, cfg.jwtSecret, time.Hour)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Failed to create token", err)
		return
	}

	respondWithJSON(w, http.StatusOK, map[string]string{
		"token": newAccessToken,
	})

}

func (cfg *apiConfig) handlerRevoke(w http.ResponseWriter, r *http.Request) {
	tokenString, err := auth.GetBearerToken(r.Header)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "Missing token", nil)
		return
	}

	err = cfg.dbQueries.RevokeRefreshToken(r.Context(), tokenString)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Failed to revoke token", err)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
