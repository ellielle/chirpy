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
	// wrap r in CORS headers
	corsMux := middlewareCors(r)
	fileseverHandler := apiCfg.middelwareMetricsInc(http.StripPrefix("/app", http.FileServer(http.Dir("."))))
	r.Handle("/app", fileseverHandler)
	r.Handle("/app/*", fileseverHandler)
	r.Route("/admin", func(r chi.Router) {
		r.Get("/metrics", apiCfg.metricsResponseHandler)
	})
	r.Route("/api", func(r chi.Router) {
		r.Get("/healthz", healthzResponseHandler)
		r.HandleFunc("/reset", apiCfg.metricsResetHandler)
	})

	server := &http.Server{
		Addr:    ":" + port,
		Handler: corsMux,
	}

	log.Printf("Serving on port: %s\n", port)
	log.Fatal(server.ListenAndServe())
}
