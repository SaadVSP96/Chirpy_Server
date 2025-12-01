package main

import (
	"fmt"
	"net/http"
)

func (cfg *apiConfig) handlerReset(w http.ResponseWriter, r *http.Request) {
	cfg.fileserverHits.Store(0)

	if cfg.platform != "dev" {
		respondWithError(w, http.StatusForbidden, "Forbidden", fmt.Errorf("reset only allowed in developer env"))
		return
	}

	err := cfg.dbQueries.DeletAllUsers(r.Context())
	if err != nil {
		respondWithError(w, http.StatusForbidden, "failed to reset users", fmt.Errorf("failed to reset users"))
		return
	}
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("All users Reset / Deleted and All Hits Set to 0"))

}
