package main

import (
	"database/sql"
	"log"
	"net/http"
	"os"
	"sync/atomic"

	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
	"github.com/peridan9/learn-http-server/internal/database"
	"github.com/peridan9/learn-http-server/internal/handlers"
)

// type apiConfig struct {
// 	fileserverHits atomic.Int32
// 	db             *database.Queries
// 	platform       string
// }

func main() {
	const filepathRoot = "."
	const port = "8080"

	godotenv.Load()
	dbURL := os.Getenv("DB_URL")
	if dbURL == "" {
		log.Fatal("DB_URL is not set")
	}

	dbConn, err := sql.Open("postgres", dbURL)
	if err != nil {
		log.Fatalf("Error opening database connection: %v", err)
	}

	dbQueries := database.New(dbConn)
	platform := os.Getenv("PLATFORM")
	if platform == "" {
		log.Fatal("PLATFORM must be set")
	}

	secretKey := os.Getenv("SECRET_KEY")
	if secretKey == "" {
		log.Fatal("SECRET_KEY must be set")
	}

	polkaKey := os.Getenv("POLKA_KEY")
	if polkaKey == "" {
		log.Fatal("POLKA_KEY must be set")
	}

	cfg := handlers.APIConfig{
		FileserverHits: atomic.Int32{},
		DB:             dbQueries,
		Platform:       platform,
		SecretKey:      secretKey,
		PolkaKey:       polkaKey,
	}

	mux := cfg.SetupRoutes()

	srv := &http.Server{
		Addr:    ":" + port,
		Handler: mux,
	}

	log.Printf("Serving files from %s on port: %s\n", filepathRoot, port)
	log.Fatal(srv.ListenAndServe())
}
