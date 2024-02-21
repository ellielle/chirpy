package main

import (
	"log"
	"net/http"
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

// FIXME: unused, and getMetricsData is probably unnecessary
func (cfg *apiConfig) metricsResponseHandler(w http.ResponseWriter, _ *http.Request) {
	w.Header().Set("Content-Type", "text/plain; charset=utf-8")
	w.WriteHeader(200)
	//w.Write([]byte())
}
