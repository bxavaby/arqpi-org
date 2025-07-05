package api

import (
	"crypto/sha256"
	"encoding/json"
	"fmt"
	"log"
	"math/rand"
	"net/http"
	"os"
	"strconv"
	"sync"
	"time"

	"github.com/bxavaby/arqpi-org/internal/models"
	"github.com/bxavaby/arqpi-org/internal/search"
	"github.com/go-chi/render"
)

type API struct {
	Fragments    []models.Fragment
	Metadata     models.Metadata
	SearchIndex  *search.SearchIndex
	StartTime    time.Time
	RequestCount int64
	rng          *rand.Rand
	clients      map[string]*client
	clientsMu    sync.Mutex
}

type client struct {
	count     int
	lastReset time.Time
	isDonor   bool
}

func NewAPI(fragments []models.Fragment, metadata models.Metadata, rng *rand.Rand) *API {
	return &API{
		Fragments:    fragments,
		Metadata:     metadata,
		SearchIndex:  search.NewSearchIndex(fragments),
		StartTime:    time.Now(),
		RequestCount: 0,
		rng:          rng,
		clients:      make(map[string]*client),
	}
}

func (a *API) GetFragment(w http.ResponseWriter, r *http.Request) {
	idStr := r.URL.Query().Get("id")
	if idStr == "" {
		http.Error(w, "Missing id parameter", http.StatusBadRequest)
		return
	}

	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "Invalid id parameter", http.StatusBadRequest)
		return
	}

	for _, fragment := range a.Fragments {
		if fragment.ID == id {
			render.JSON(w, r, fragment.ToResponse())
			return
		}
	}

	http.Error(w, "Fragment not found", http.StatusNotFound)
}

func (a *API) GetRandomFragment(w http.ResponseWriter, r *http.Request) {
	if len(a.Fragments) == 0 {
		http.Error(w, "No fragments available", http.StatusInternalServerError)
		return
	}

	randIndex := a.rng.Intn(len(a.Fragments))
	render.JSON(w, r, a.Fragments[randIndex].ToResponse())
}

func (a *API) SearchFragments(w http.ResponseWriter, r *http.Request) {
	query := r.URL.Query().Get("q")
	if query == "" {
		http.Error(w, "Missing search query", http.StatusBadRequest)
		return
	}

	limitStr := r.URL.Query().Get("limit")
	limit := 10 // dflt
	if limitStr != "" {
		parsedLimit, err := strconv.Atoi(limitStr)
		if err == nil && parsedLimit > 0 {
			limit = parsedLimit
		}
	}

	results := a.SearchIndex.Search(query, limit)

	responses := make([]models.FragmentResponse, 0, len(results))
	for _, fragment := range results {
		responses = append(responses, fragment.ToResponse())
	}

	render.JSON(w, r, responses)
}
func (a *API) GetInfo(w http.ResponseWriter, r *http.Request) {
	render.JSON(w, r, a.Metadata)
}

func (a *API) GetQuote(w http.ResponseWriter, r *http.Request) {
	var potentialQuotes []models.Fragment
	for _, fragment := range a.Fragments {
		if fragment.Length > 50 && fragment.Length < 300 {
			potentialQuotes = append(potentialQuotes, fragment)
		}
	}

	if len(potentialQuotes) == 0 {
		for i := 0; i < 20 && i < len(a.Fragments); i++ {
			minIndex := i
			for j := i + 1; j < len(a.Fragments); j++ {
				if a.Fragments[j].Length < a.Fragments[minIndex].Length {
					minIndex = j
				}
			}
			if i != minIndex {
				a.Fragments[i], a.Fragments[minIndex] = a.Fragments[minIndex], a.Fragments[i]
			}
			potentialQuotes = append(potentialQuotes, a.Fragments[i])
		}
	}

	if len(potentialQuotes) == 0 {
		http.Error(w, "No suitable quotes found", http.StatusInternalServerError)
		return
	}

	randIndex := a.rng.Intn(len(potentialQuotes))
	quote := potentialQuotes[randIndex]

	render.JSON(w, r, map[string]any{
		"id":    quote.ID,
		"text":  quote.Text,
		"title": quote.Title,
		"url":   quote.URL,
	})
}

func (a *API) GetStatus(w http.ResponseWriter, r *http.Request) {
	uptime := time.Since(a.StartTime)

	clientID := getClientID(r)

	a.clientsMu.Lock()
	userRequests := 0
	remainingRequests := 0

	rateLimit := 60
	if envLimit := os.Getenv("API_RATE_LIMIT"); envLimit != "" {
		if val, err := strconv.Atoi(envLimit); err == nil && val > 0 {
			rateLimit = val
		}
	}

	if client, exists := a.clients[clientID]; exists {
		userRequests = client.count
		remainingRequests = max(rateLimit-client.count, 0)
	} else {
		remainingRequests = rateLimit
	}
	a.clientsMu.Unlock()

	apiKey := r.URL.Query().Get("key")
	isDonor := isDonorKey(apiKey)

	render.JSON(w, r, map[string]any{
		"status":          "operational",
		"uptime":          uptime.String(),
		"version":         "1.0.0",
		"fragment_count":  len(a.Fragments),
		"total_requests":  a.RequestCount,
		"your_requests":   userRequests,
		"remaining_limit": remainingRequests,
		"is_donor":        isDonor,
	})
}

func (a *API) HandleKofiWebhook(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	kofiToken := os.Getenv("KOFI_VERIFICATION_TOKEN")
	if kofiToken == "" {
		log.Println("Error: KOFI_VERIFICATION_TOKEN not set")
		http.Error(w, "Configuration error", http.StatusInternalServerError)
		return
	}

	token := r.Header.Get("Kofi-Verification-Token")
	if token != kofiToken {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	var kofiData struct {
		Data struct {
			Email             string `json:"email"`
			Name              string `json:"name"`
			Amount            string `json:"amount"`
			KofiTransactionID string `json:"kofi_transaction_id"`
		} `json:"data"`
		MessageType string `json:"message_type"`
	}

	if err := json.NewDecoder(r.Body).Decode(&kofiData); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	apiKey := generateAPIKey(kofiData.Data.Email)

	log.Printf("New donation received from %s (%s), Amount: %s",
		kofiData.Data.Name, kofiData.Data.Email, kofiData.Data.Amount)
	log.Printf("API Key generated: %s", apiKey)
	log.Printf("Please add this key to your DONOR_API_KEYS environment variable in Render")

	render.JSON(w, r, map[string]any{
		"status":  "success",
		"message": "Thank you for your support!",
		"key":     apiKey,
	})
}

func generateAPIKey(seed string) string {
	salt := os.Getenv("API_KEY_SALT")
	if salt == "" {
		salt = "]qv7U4Y(Xww2<r>MXy" // fallback
	}

	h := sha256.New()
	h.Write([]byte(seed + salt + time.Now().String()))
	return fmt.Sprintf("%x", h.Sum(nil))[:32]
}
