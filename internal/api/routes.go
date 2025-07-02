package api

import (
	"net/http"
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
