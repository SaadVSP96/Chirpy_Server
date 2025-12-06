package main

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"os"

	"github.com/SaadVSP96/Chirpy_Server.git/internal/database"
	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
)

func main() {
	// Load .env file
	if err := godotenv.Load(); err != nil {
		log.Println("Warning: .env file not found")
	}

	dbURL := os.Getenv("DB_URL")
	if dbURL == "" {
		log.Fatal("DB_URL environment variable is required")
	}

	db, err := sql.Open("postgres", dbURL)
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}
	defer db.Close()

	// Test connection
	if err := db.Ping(); err != nil {
		log.Fatalf("Failed to ping database: %v", err)
	}

	q := database.New(db)

	// Delete all data
	if err := q.DeleteAllUsers(context.Background()); err != nil {
		log.Printf("Error deleting users: %v", err)
	}
	if err := q.DeleteAllChirps(context.Background()); err != nil {
		log.Printf("Error deleting chirps: %v", err)
	}

	fmt.Println("Database cleaned successfully!")
}
