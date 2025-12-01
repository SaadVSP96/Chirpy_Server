package main

import (
	"encoding/json"
	"net/http"
	"strings"
)

func handlerChirpsValidate(w http.ResponseWriter, r *http.Request) {
	type parameters struct {
		Body string `json:"body"`
	}
	type returnVals struct {
		CleanedBody string `json:"cleaned_body"`
	}

	decoder := json.NewDecoder(r.Body)
	params := parameters{}
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

	// Now add the profanity checker
	params.Body = profanityCleaner(params.Body)

	respondWithJSON(w, http.StatusOK, returnVals{
		CleanedBody: params.Body,
	})
}

func profanityCleaner(s string) string {
	words := strings.Fields(s)

	// make the new cleaned string slice to later join and send back
	resultSlice := []string{}

	for _, w := range words {
		lw := strings.ToLower(w) // case-insensitive match
		if lw == "kerfuffle" || lw == "sharbert" || lw == "fornax" {
			resultSlice = append(resultSlice, "****")
			continue
		}
		resultSlice = append(resultSlice, w)
	}
	return strings.Join(resultSlice, " ")
}
