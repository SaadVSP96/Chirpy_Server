package main

import (
	"database/sql"
	"log"
	"net/http"
	"os"
	"sync/atomic"

	"github.com/SaadVSP96/Chirpy_Server.git/internal/database"
	"github.com/joho/godotenv"
	_ "github.com/lib/pq" // Postgres driver
	// SQLC generated package
)

type apiConfig struct {
	fileserverHits atomic.Int32
	dbQueries      *database.Queries
	platform       string
}

func main() {
	const filepathRoot = "."
	const port = "8080"
	// Load .env file
	godotenv.Load()
	// get db url and connect to db
	dbURL := os.Getenv("DB_URL")
	platform := os.Getenv("PLATFORM")
	db, err := sql.Open("postgres", dbURL)
	if err != nil {
		log.Fatalf("Failed to connect to DB: %v", err)
	}
	dbQueries := database.New(db)
	apiCfg := apiConfig{
		fileserverHits: atomic.Int32{},
		dbQueries:      dbQueries,
		platform:       platform,
	}
	defer db.Close()
	// server and endpoints logic.
	mux := http.NewServeMux()
	fsHandler := apiCfg.middlewareMetricsInc(http.StripPrefix("/app", http.FileServer(http.Dir(filepathRoot))))
	mux.Handle("/app/", fsHandler)

	mux.HandleFunc("GET /api/healthz", handlerReadiness)
	mux.HandleFunc("POST /api/validate_chirp", handlerChirpsValidate)
	mux.HandleFunc("POST /api/users", apiCfg.handlerCreateUser)
	mux.HandleFunc("POST /api/chirps", apiCfg.handlerChirps)
	mux.HandleFunc("GET /api/chirps/", apiCfg.handlerGetChirpByID)
	mux.HandleFunc("GET /api/chirps", apiCfg.handlerGetChirps)

	mux.HandleFunc("POST /admin/reset", apiCfg.handlerReset)
	mux.HandleFunc("GET /admin/metrics", apiCfg.handlerMetrics)

	srv := &http.Server{
		Addr:    ":" + port,
		Handler: mux,
	}

	log.Printf("Serving files from %s on port: %s\n", filepathRoot, port)
	log.Fatal(srv.ListenAndServe())
}
