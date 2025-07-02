package api

import (
	"encoding/json"
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

	type client struct {
		count     int
		lastReset time.Time
	}
	clients := make(map[string]*client)

	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			ip := r.RemoteAddr

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
