package main

import (
	"fmt"
	"net/http"
)

// Handler that returns hit count metric from /admin/metrics
func (cfg *apiConfig) handlerMetricsResponse(w http.ResponseWriter, _ *http.Request) {
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	w.WriteHeader(200)
	w.Write([]byte(fmt.Sprintf("<html><body><h1>Welcome, Chirpy Admin</h1><p>Chirpy has been visited %d times!</p></body></html>", cfg.fileserverHits)))
}

// Handler for resetting hit count when entering /api/reset
func (cfg *apiConfig) handlerMetricsReset(w http.ResponseWriter, _ *http.Request) {
	cfg.fileserverHits = 0
	w.Header().Set("Content-Type", "text/plain; charset=utf-8")
	w.WriteHeader(200)
	w.Write([]byte("OK"))
}

// Middleware handler to add a hit each time /app is visited
func (cfg *apiConfig) middelwareMetricsInc(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		cfg.fileserverHits++
		next.ServeHTTP(w, r)
	})
}
