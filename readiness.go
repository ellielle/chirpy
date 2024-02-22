package main

import "net/http"

// Responds to the /api/healthz endpoint with readiness indication
func healthzResponseHandler(w http.ResponseWriter, _ *http.Request) {
	w.Header().Set("Content-Type", "text/plain; charset=utf-8")
	w.WriteHeader(200)
	w.Write([]byte("OK"))
}
