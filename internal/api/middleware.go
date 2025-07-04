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
	log.Printf("RateLimiter initialization: Donor keys env variable length: %d", len(donorKeysStr))

	donorKeys := make(map[string]bool)
	if donorKeysStr != "" {
		for _, key := range strings.Split(donorKeysStr, ",") {
			trimmedKey := strings.TrimSpace(key)
			donorKeys[trimmedKey] = true
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
			log.Printf("Request with key: %s", maskKey(apiKey))

			if apiKey != "" {
				// direct check against map
				isDonor := donorKeys[apiKey]
				log.Printf("Key donor status: %v", isDonor)

				if isDonor {
					log.Printf("Bypassing rate limit for donor key")
					next.ServeHTTP(w, r)
					return
				}
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

func isDonorKey(apiKey string) bool {
	if apiKey == "" {
		log.Printf("isDonorKey: Empty key provided, returning false")
		return false
	}

	donorKeysEnv := os.Getenv("DONOR_API_KEYS")
	if donorKeysEnv == "" {
		log.Printf("isDonorKey: No donor keys in environment, returning false")
		return false
	}

	log.Printf("isDonorKey: Checking against env variable (length: %d)", len(donorKeysEnv))

	donorKeySlice := strings.Split(donorKeysEnv, ",")
	log.Printf("isDonorKey: Split into %d keys", len(donorKeySlice))

	for i, donorKey := range donorKeySlice {
		cleanKey := strings.TrimSpace(donorKey)
		log.Printf("isDonorKey: Key %d (masked): %s", i, maskKey(cleanKey))
		if cleanKey == apiKey {
			log.Printf("isDonorKey: Match found!")
			return true
		}
	}

	log.Printf("isDonorKey: No match found")
	return false
}

func maskKey(key string) string {
	if len(key) <= 8 {
		return "***"
	}
	return key[:4] + "..." + key[len(key)-4:]
}

func CacheControl(duration string) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Cache-Control", "public, max-age="+duration)
			next.ServeHTTP(w, r)
		})
	}
}
