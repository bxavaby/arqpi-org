package api

import (
	"encoding/json"
	"log"
	"net/http"
	"os"
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

	log.Printf("Rate limiter initialized: %d requests per %d seconds", limit, windowSecs)
	log.Printf("Loaded %d donor keys", len(donorKeys))

	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			log.Printf("Request from: %s, Path: %s", getClientID(r), r.URL.Path)

			if r.Method == "OPTIONS" {
				next.ServeHTTP(w, r)
				return
			}

			apiKey := r.URL.Query().Get("key")
			if apiKey != "" && isDonorKey(apiKey) {
				log.Printf("Donor key used: %s", apiKey[:8]+"...")
				next.ServeHTTP(w, r)
				return
			}

			clientID := getClientID(r)
			log.Printf("Client ID: %s", clientID)

			a.clientsMu.Lock()
			c, exists := a.clients[clientID]

			if !exists {
				log.Printf("New client: %s", clientID)
				a.clients[clientID] = &client{
					count:     1,
					lastReset: time.Now(),
				}
				a.clientsMu.Unlock()
				next.ServeHTTP(w, r)
				return
			}

			if time.Since(c.lastReset).Seconds() > float64(windowSecs) {
				log.Printf("Resetting count for client: %s (previous count: %d)", clientID, c.count)
				c.count = 1
				c.lastReset = time.Now()
				a.clientsMu.Unlock()
				next.ServeHTTP(w, r)
				return
			}

			if c.count >= limit {
				log.Printf("Rate limit exceeded for client: %s (count: %d)", clientID, c.count)
				a.clientsMu.Unlock()

				w.Header().Set("Retry-After", strconv.Itoa(windowSecs))
				w.Header().Set("Content-Type", "application/json")
				w.WriteHeader(http.StatusTooManyRequests)

				donationMessage := map[string]any{
					"error":               "Rate limit exceeded",
					"message":             "You've reached the free usage limit. Consider supporting this project on Ko-fi to get unlimited access.",
					"donate_url":          "https://ko-fi.com/bxav",
					"retry_after_seconds": windowSecs,
					"current_count":       c.count,
					"limit":               limit,
				}

				json.NewEncoder(w).Encode(donationMessage)
				return
			}

			c.count++
			log.Printf("Incremented count for client: %s (new count: %d)", clientID, c.count)
			a.clientsMu.Unlock()

			next.ServeHTTP(w, r)
		})
	}
}

func getClientID(r *http.Request) string {
	ip := r.Header.Get("X-Forwarded-For")
	if ip == "" {
		ip = r.RemoteAddr
	}

	if commaIdx := strings.Index(ip, ","); commaIdx > 0 {
		ip = strings.TrimSpace(ip[:commaIdx])
	}

	ua := r.Header.Get("User-Agent")
	if len(ua) > 50 {
		ua = ua[:50]
	}

	return ip + "|" + ua
}

func isDonorKey(key string) bool {
	if key == "" {
		return false
	}

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
