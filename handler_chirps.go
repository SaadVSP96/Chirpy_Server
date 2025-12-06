package main

import (
	"database/sql"
	"encoding/json"
	"net/http"
	"strings"
	"time"

	"github.com/google/uuid"

	"github.com/SaadVSP96/Chirpy_Server.git/internal/auth"
	"github.com/SaadVSP96/Chirpy_Server.git/internal/database"
)

type CreateChirpRequest struct {
	Body string `json:"body"`
}

type CreateChirpResponse struct {
	ID        string `json:"id"`
	CreatedAt string `json:"created_at"`
	UpdatedAt string `json:"updated_at"`
	Body      string `json:"body"`
	UserID    string `json:"user_id"`
}

func (cfg *apiConfig) handlerGetChirpByID(w http.ResponseWriter, r *http.Request) {
	// Expected URL: /api/chirps/{chirpID}
	pathParts := strings.Split(r.URL.Path, "/")
	if len(pathParts) < 4 {
		http.NotFound(w, r)
		return
	}
	chirpIDStr := pathParts[3]

	chirpUUID, err := uuid.Parse(chirpIDStr)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid chirp ID format", err)
		return
	}

	chirp, err := cfg.dbQueries.GetChirpByID(r.Context(), chirpUUID)
	if err != nil {
		if err == sql.ErrNoRows {
			http.NotFound(w, r)
			return
		}
		respondWithError(w, http.StatusInternalServerError, "Error fetching chirp", err)
		return
	}

	resp := CreateChirpResponse{
		ID:        chirp.ID.String(),
		CreatedAt: chirp.CreatedAt.Format(time.RFC3339),
		UpdatedAt: chirp.UpdatedAt.Format(time.RFC3339),
		Body:      chirp.Body,
		UserID:    chirp.UserID.String(),
	}
	respondWithJSON(w, http.StatusOK, resp)
}

func (cfg *apiConfig) handlerGetChirps(w http.ResponseWriter, r *http.Request) {
	// Check if author_id query parameter is provided
	authorIDStr := r.URL.Query().Get("author_id")

	var chirps []database.Chirp
	var err error

	if authorIDStr != "" {
		// Filter by specific author
		authorUUID, err := uuid.Parse(authorIDStr)
		if err != nil {
			respondWithError(w, http.StatusBadRequest, "Invalid author ID format", err)
			return
		}
		chirps, err = cfg.dbQueries.ListChirpsByUser(r.Context(), authorUUID)
	} else {
		// Get all chirps
		chirps, err = cfg.dbQueries.ListChirps(r.Context())
	}

	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Error while getting chirps from db", err)
		return
	}

	// convert []Chirp -> []CreateChirpResponse
	response := make([]CreateChirpResponse, 0, len(chirps))

	for _, c := range chirps {
		response = append(response, CreateChirpResponse{
			ID:        c.ID.String(),
			CreatedAt: c.CreatedAt.Format(time.RFC3339),
			UpdatedAt: c.UpdatedAt.Format(time.RFC3339),
			Body:      c.Body,
			UserID:    c.UserID.String(),
		})
	}
	respondWithJSON(w, http.StatusOK, response)
}

func (cfg *apiConfig) handlerChirps(w http.ResponseWriter, r *http.Request) {
	var params CreateChirpRequest

	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(&params)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Couldn't decode parameters", err)
		return
	}

	tokenString, err := auth.GetBearerToken(r.Header)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "Invalid or missing token", err)
		return
	}

	userID, err := auth.ValidateJWT(tokenString, cfg.jwtSecret)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "Invalid token", err)
		return
	}

	const maxChirpLength = 140
	if len(params.Body) > maxChirpLength {
		respondWithError(w, http.StatusBadRequest, "Chirp is too long", nil)
		return
	}

	if userID.String() == "" {
		respondWithError(w, http.StatusBadRequest, "user id is missing", nil)
		return
	}

	// Now add the profanity checker
	params.Body = profanityCleaner(params.Body)
	chirp, err := cfg.dbQueries.CreateChirp(r.Context(), database.CreateChirpParams{
		Body:   params.Body,
		UserID: userID,
	})

	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Couldn't create chirp", err)
		return
	}
	respondWithJSON(w, http.StatusCreated, CreateChirpResponse{
		ID:        chirp.ID.String(),
		CreatedAt: chirp.CreatedAt.Format(time.RFC3339),
		UpdatedAt: chirp.UpdatedAt.Format(time.RFC3339),
		Body:      chirp.Body,
		UserID:    chirp.UserID.String(),
	})
}

func (cfg *apiConfig) handlerDeleteChirp(w http.ResponseWriter, r *http.Request) {
	// Parse chirpID from URL path manually (same approach as handlerGetChirpByID)
	pathParts := strings.Split(r.URL.Path, "/")
	if len(pathParts) < 4 {
		respondWithError(w, http.StatusBadRequest, "Invalid chirp ID format", nil)
		return
	}
	chirpIDStr := pathParts[3]

	chirpID, err := uuid.Parse(chirpIDStr)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid Chirp ID", err)
		return
	}

	// Rest of your existing code remains the same...
	tokenStr, err := auth.GetBearerToken(r.Header)
	if err != nil {
		respondWithError(w, 401, "Missing or invalid auth token", err)
		return
	}

	userID, err := auth.ValidateJWT(tokenStr, cfg.jwtSecret)
	if err != nil {
		respondWithError(w, 401, "Invalid token", err)
		return
	}

	chirp, err := cfg.dbQueries.GetChirpByID(r.Context(), chirpID)
	if err != nil {
		respondWithError(w, 404, "Chirp not found", err)
		return
	}

	if chirp.UserID != userID {
		respondWithError(w, 403, "Forbidden: you cannot delete another user's chirp", nil)
		return
	}

	rows, err := cfg.dbQueries.DeleteChirp(r.Context(), database.DeleteChirpParams{ID: chirpID, UserID: userID})
	if err != nil {
		respondWithError(w, 500, "Failed to delete chirp", err)
		return
	}
	if rows == 0 {
		respondWithError(w, 404, "Chirp not found", nil)
		return
	}

	w.WriteHeader(204)
}
