package main

import (
	"encoding/json"
	"net/http"
	"strings"
)

func validateChirpHandler(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()
	type parameters struct {
		Body string `json:"body"`
	}

	type returnValid struct {
		Valid bool `json:"valid"`
	}

	type returnCleaned struct {
		CleanedBody string `json:"cleaned_body"`
	}

	theProfane := []string{"kerfuffle", "sharbert", "fornax"}

	// Create a new JSON decoder and check the validity of the JSON from the Request body
	decoder := json.NewDecoder(r.Body)
	params := parameters{}
	err := decoder.Decode(&params)
	if err != nil {
		respondWithError(w, 400, "Bad Request")
		return
	}

	if len(params.Body) > 140 {
		respondWithError(w, 400, "Chirp is too long")
		return
	}

	// Checks for profanity usage by looping over theProfane slice and checking the words against a lower cased copy of params.Body
	// FIXME: lower casing the string is causing the original to be borked
	// redo algorithm
	tempStr := strings.ToLower(params.Body)
	hasProfanity := false
	for _, profanity := range theProfane {
		if strings.Contains(tempStr, profanity) {
			splitStr := strings.Split(tempStr, " ")
			for j, str := range splitStr {
				if str == profanity {
					splitStr[j] = "****"
				}
			}
			tempStr = strings.Join(splitStr, " ")
			//strings.ReplaceAll(params.Body, profanity, "****")
			hasProfanity = true
		}
	}

	if hasProfanity {
		respondWithJSON(w, 200, returnCleaned{CleanedBody: tempStr})
		return
	}
	respondWithJSON(w, 200, returnValid{Valid: true})
}
