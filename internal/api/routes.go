package api

import (
	"net/http"
	"os"
	"strconv"
	"sync/atomic"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/cors"
)

func (a *API) SetupRoutes() http.Handler {
	r := chi.NewRouter()

	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)
	r.Use(a.requestCounter)

	// 60 requests/h for non-donors
	rateLimit := 60
	rateWindow := 3600

	if envLimit := os.Getenv("API_RATE_LIMIT"); envLimit != "" {
		if val, err := strconv.Atoi(envLimit); err == nil && val > 0 {
			rateLimit = val
		}
	}

	if envWindow := os.Getenv("API_RATE_WINDOW"); envWindow != "" {
		if val, err := strconv.Atoi(envWindow); err == nil && val > 0 {
			rateWindow = val
		}
	}

	// Apply rate limiter middleware
	r.Use(a.RateLimiter(rateLimit, rateWindow))

	// CORS for frontend
	r.Use(cors.Handler(cors.Options{
		AllowedOrigins:   []string{"https://arqpi.org", "https://www.arqpi.org"},        // Consider restricting this in production
		AllowedMethods:   []string{"GET", "POST", "OPTIONS"},                            // Added POST for webhook
		AllowedHeaders:   []string{"Accept", "Content-Type", "Kofi-Verification-Token"}, // Added Ko-fi header
		ExposedHeaders:   []string{"Link"},
		AllowCredentials: false,
		MaxAge:           300,
	}))

	// API routes
	r.Get("/fragment", a.GetFragment)
	r.Get("/random", a.GetRandomFragment)
	r.Get("/search", a.SearchFragments)
	r.Get("/info", a.GetInfo)
	r.Get("/quote", a.GetQuote)
	r.Get("/status", a.GetStatus)

	// webhook endpoint
	r.Post("/webhook/kofi", a.HandleKofiWebhook)

	return r
}

func (a *API) requestCounter(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		atomic.AddInt64(&a.RequestCount, 1)
		next.ServeHTTP(w, r)
	})
}
