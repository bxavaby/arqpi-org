package main

import (
	"encoding/json"
	"log"
	"log/slog"
	"math/rand"
	"net/http"
	"os"
	"time"

	"github.com/bxavaby/arqpi-org/internal/api"
	"github.com/bxavaby/arqpi-org/internal/models"
)

func main() {
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))

	fragmentsFile, err := os.ReadFile("data/all_fragments.json")
	if err != nil {
		log.Fatalf("Error reading fragments file: %v", err)
	}

	var fragments []models.Fragment
	if err := json.Unmarshal(fragmentsFile, &fragments); err != nil {
		log.Fatalf("Error parsing fragments: %v", err)
	}

	log.Printf("Loaded %d fragments", len(fragments))

	metadataFile, err := os.ReadFile("data/metadata.json")
	if err != nil {
		log.Fatalf("Error reading metadata file: %v", err)
	}

	var metadata models.Metadata
	if err := json.Unmarshal(metadataFile, &metadata); err != nil {
		log.Fatalf("Error parsing metadata: %v", err)
	}

	// make api
	apiInstance := api.NewAPI(fragments, metadata, rng)

	// make routes
	router := apiInstance.SetupRoutes()

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080" // dflt
	}

	log.Printf("Server starting on port %s", port)
	if err := http.ListenAndServe(":"+port, router); err != nil {
		log.Fatalf("Error starting server: %v", err)
	}

	logger := slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
		Level: slog.LevelInfo,
	}))
	slog.SetDefault(logger)
}
