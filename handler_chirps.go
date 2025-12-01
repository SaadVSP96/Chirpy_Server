package main

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/google/uuid"

	"github.com/SaadVSP96/Chirpy_Server.git/internal/database"
)

type CreateChirpRequest struct {
	Body   string `json:"body"`
	UserId string `json:"user_id"`
}

type CreateChirpResponse struct {
	ID        string `json:"id"`
	CreatedAt string `json:"created_at"`
	UpdatedAt string `json:"updated_at"`
	Body      string `json:"body"`
	UserID    string `json:"user_id"`
}

func (cfg *apiConfig) handlerChirps(w http.ResponseWriter, r *http.Request) {
	var params CreateChirpRequest

	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(&params)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Couldn't decode parameters", err)
		return
	}

	const maxChirpLength = 140
	if len(params.Body) > maxChirpLength {
		respondWithError(w, http.StatusBadRequest, "Chirp is too long", nil)
		return
	}

	if params.UserId == "" {
		respondWithError(w, http.StatusBadRequest, "user id is missing", nil)
		return
	}

	// Now add the profanity checker
	params.Body = profanityCleaner(params.Body)
	userUUID, err := uuid.Parse(params.UserId)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid user ID format", err)
		return
	}
	chirp, err := cfg.dbQueries.CreateChirp(r.Context(), database.CreateChirpParams{
		Body:   params.Body,
		UserID: userUUID,
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
