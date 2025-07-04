package api

import (
	"encoding/json"
	"net/http"
	"os"
	"slices"
	"strconv"
	"strings"
	"time"
)

func (a *API) RateLimiter(limit int, windowSecs int) func(http.Handler) http.Handler {
	donorKeysStr := os.Getenv("DONOR_API_KEYS")
	donorKeys := make(map[string]bool)
	if donorKeysStr != "" {
		for key := range strings.SplitSeq(donorKeysStr, ",") {
			donorKeys[strings.TrimSpace(key)] = true
		}
	}

	type client struct {
		count     int
		lastReset time.Time
	}
	clients := make(map[string]*client)

	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

			if origin := r.Header.Get("Origin"); origin != "" {
				allowedOrigins := []string{"https://arqpi.org", "https://www.arqpi.org", "https://bxavaby.github.io"}
				if slices.Contains(allowedOrigins, origin) {
					w.Header().Set("Access-Control-Allow-Origin", origin)
				}
			}

			// return immediately
			if r.Method == "OPTIONS" {
				w.Header().Set("Access-Control-Allow-Methods", "GET, POST, OPTIONS")
				w.Header().Set("Access-Control-Allow-Headers", "Accept, Content-Type, Kofi-Verification-Token")
				w.Header().Set("Access-Control-Max-Age", "300")
				w.WriteHeader(http.StatusOK)
				return
			}

			ip := r.Header.Get("X-Forwarded-For")
			if ip == "" {
				ip = r.RemoteAddr
			}

			apiKey := r.URL.Query().Get("key")
			if apiKey != "" && isDonorKey(apiKey) {
				next.ServeHTTP(w, r)
				return
			}

			c, exists := clients[ip]
			if !exists {
				clients[ip] = &client{
					count:     1,
					lastReset: time.Now(),
				}
				next.ServeHTTP(w, r)
				return
			}

			if time.Since(c.lastReset).Seconds() > float64(windowSecs) {
				c.count = 1
				c.lastReset = time.Now()
				next.ServeHTTP(w, r)
				return
			}

			if c.count >= limit {
				w.Header().Set("Retry-After", strconv.Itoa(windowSecs))
				w.Header().Set("Content-Type", "application/json")
				w.WriteHeader(http.StatusTooManyRequests)

				donationMessage := map[string]any{
					"error":               "Rate limit exceeded",
					"message":             "You've reached the free usage limit. Consider supporting this project on Ko-fi to get unlimited access.",
					"donate_url":          "https://ko-fi.com/bxav",
					"retry_after_seconds": windowSecs,
				}

				json.NewEncoder(w).Encode(donationMessage)
				return
			}

			c.count++
			next.ServeHTTP(w, r)
		})
	}
}

func isDonorKey(key string) bool {
	donorKeysEnv := os.Getenv("DONOR_API_KEYS")
	if donorKeysEnv == "" {
		return false
	}

	donorKeys := strings.SplitSeq(donorKeysEnv, ",")
	for k := range donorKeys {
		if strings.TrimSpace(k) == key {
			return true
		}
	}
	return false
}

func CacheControl(duration string) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Cache-Control", "public, max-age="+duration)
			next.ServeHTTP(w, r)
		})
	}
}
