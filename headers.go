package main

import (
	"encoding/json"
	"log"
	"net/http"
)

// Middleware handler to add basic (and open, which is necessary for the course to access the server) CORS headers
func middlewareCors(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE")
		w.Header().Set("Access-Control-Allow-Headers", "*")

		if r.Method == "OPTIONS" {
			w.WriteHeader(http.StatusOK)
			return
		}
		next.ServeHTTP(w, r)
	})
}

// Wraps respondWithJSON to respond with an Error in JSON format
func respondWithError(w http.ResponseWriter, code int, message string) error {
	type returnError struct {
		Error string `json:"error"`
	}
	return respondWithJSON(w, code, returnError{Error: message})
}

// Sends a JSON response with the request's ResponseWriter, reponse code, and payload. Validates JSON and responds with an error if invalid
func respondWithJSON(w http.ResponseWriter, code int, payload interface{}) error {
	response, err := json.Marshal(payload)
	if err != nil {
		log.Printf("Error marshalling JSON: %s", err)
		w.WriteHeader(500)
		return nil
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	w.Write(response)
	return nil
}
