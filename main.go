package main

import (
	"log"
	"net/http"

	"github.com/go-chi/chi/v5"
)

func main() {
	const port = "8080"
	apiCfg := &apiConfig{
		fileserverHits: 0,
	}
	r := chi.NewRouter()
	// Wrap r in CORS headers
	corsMux := middlewareCors(r)

	// Fileserver for handling static pages
	fileseverHandler := apiCfg.middelwareMetricsInc(http.StripPrefix("/app", http.FileServer(http.Dir("."))))
	r.Handle("/app", fileseverHandler)
	r.Handle("/app/*", fileseverHandler)

	// Admin route, which only contains the metrics endpoint for now
	r.Route("/admin", func(r chi.Router) {
		r.Get("/metrics", apiCfg.metricsResponseHandler)
	})

	// Subroutes under /api
	r.Route("/api", func(r chi.Router) {
		// Health check endpoing
		r.Get("/healthz", healthzResponseHandler)
		// Page hit count reset endpoint
		r.HandleFunc("/reset", apiCfg.metricsResetHandler)
		// POST endpoint to submit "Chirps". Chrips must be 140 chars or less, and should be in JSON
		r.Post("/chirps", validateChirpHandler)
		// GET endpoint for retrieving all chirps
		r.Get("/chirps", getChirpsHandler)
		// GET endpoint for retrieving a single chirp
		r.Get("/chirps/{chirpID}", getSingleChirpHandler)
	})

	server := &http.Server{
		Addr:    ":" + port,
		Handler: corsMux,
	}

	log.Printf("Serving on port: %s\n", port)
	log.Fatal(server.ListenAndServe())
}
