package main

import (
	"log"
	"net/http"
	"strconv"
)

func main() {
	const port = "8080"
	apiCfg := &apiConfig{
		fileserverHits: 0,
	}
	mux := http.NewServeMux()
	// wrap mux in CORS headers
	corsMux := middlewareCors(mux)
	fileseverHandler := http.StripPrefix("/app", http.FileServer(http.Dir(".")))
	mux.Handle("/app/", apiCfg.middelwareMetricsInc(fileseverHandler))
	//mux.Handle("/app/assets", http.StripPrefix("/app", http.FileServer(http.Dir("./assets"))))
	mux.HandleFunc("/healthz", healthzResponseHandler)
	mux.HandleFunc("/metrics", apiCfg.metricsResponseHandler)
	mux.HandleFunc("/reset", apiCfg.metricsResetHandler)

	server := &http.Server{
		Addr:    ":" + port,
		Handler: corsMux,
	}

	log.Printf("Serving on port: %s\n", port)
	log.Fatal(server.ListenAndServe())
}

func healthzResponseHandler(w http.ResponseWriter, _ *http.Request) {
	w.Header().Set("Content-Type", "text/plain; charset=utf-8")
	w.WriteHeader(200)
	w.Write([]byte("OK"))
}

func (cfg *apiConfig) metricsResponseHandler(w http.ResponseWriter, _ *http.Request) {
	numRequests := strconv.Itoa(cfg.fileserverHits)
	hitsStr := "Hits: " + numRequests
	w.Header().Set("Content-Type", "text/plain; charset=utf-8")
	w.WriteHeader(200)
	w.Write([]byte(hitsStr))
}

func (cfg *apiConfig) metricsResetHandler(w http.ResponseWriter, _ *http.Request) {
	cfg.fileserverHits = 0
	w.Header().Set("Content-Type", "text/plain; charset=utf-8")
	w.WriteHeader(200)
	w.Write([]byte("OK"))
}
