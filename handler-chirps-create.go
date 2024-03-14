package main

import (
	"encoding/json"
	"errors"
	"net/http"
	"strings"
)

func (cfg apiConfig) handlerChirpsCreate(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()

	type parameters struct {
		Body string `json:"body"`
	}

	// Create a new JSON decoder and check the validity of the JSON from the Request body
	decoder := json.NewDecoder(r.Body)
	params := parameters{}
	err := decoder.Decode(&params)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Malformed request body")
		return
	}

	// Validate that the Chirp meets all requirements
	cleanedChirp, err := validateChirp(params.Body)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, err.Error())
		return
	}

	// Create a new chirp with the body and save it to database
	chirp, err := cfg.DB.CreateChirp(cleanedChirp)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}
	respondWithJSON(w, http.StatusOK, chirp)
}

// Checks for profanity usage by looping over theProfane slice and checking the words against a lower cased params.Body
func validateChirp(bodyText string) (string, error) {
	// Disallow any chirps longer than 140 characters
	if len(bodyText) > 140 {
		return "", errors.New("Chirp is too long")
	}
	cleanedBody := getCleanedBody(bodyText)
	return cleanedBody, nil
}

func getCleanedBody(bodyText string) string {
	// Remove profanity because this is a Christian Minecraft server
	// The profanity list is created as a map with an empty struct to be easily matched against
	theProfane := map[string]struct{}{"kerfuffle": {}, "sharbert": {}, "fornax": {}}
	splitBody := strings.Split(bodyText, " ")

	for i, word := range splitBody {
		lowerWord := strings.ToLower(word)
		if _, ok := theProfane[lowerWord]; ok {
			splitBody[i] = "****"
		}
	}
	cleanedBody := strings.Join(splitBody, " ")
	return cleanedBody
}