package main

import (
	"encoding/json"
	"log"
	"net/http"
	"strings"

	database "github.com/ellielle/chirpy/internal/database"
)

type Chirp struct {
	id   int
	body string
}

func validateChirpHandler(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()
	type parameters struct {
		Body string `json:"body"`
	}

	type returnValid struct {
		Valid bool `json:"valid"`
	}

	type returnBody struct {
		Body string `json:"cleaned_body"`
	}

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

	cleanedBody, hasProfanity := validateNoProfanity(params.Body)
	// TODO: give unique ID and save to database.json
	// respondWithError with error returned if issues

	if hasProfanity {
		respondWithJSON(w, 200, returnBody{Body: cleanedBody})
		return
	}
	respondWithJSON(w, 200, returnBody{Body: params.Body})
}

func validateNoProfanity(bodyText string) (cleanedText string, hasProfanity bool) {
	// Checks for profanity usage by looping over theProfane slice and checking the words against a lower cased params.Body
	theProfane := []string{"kerfuffle", "sharbert", "fornax"}

	splitStr := strings.Split(bodyText, " ")
	for _, profanity := range theProfane {
		if strings.Contains(strings.ToLower(bodyText), strings.ToLower(profanity)) {
			for j, str := range splitStr {
				if strings.EqualFold(str, profanity) {
					splitStr[j] = "****"
				}
			}
			cleanedText = strings.Join(splitStr, " ")
			hasProfanity = true
		}
	}
	return cleanedText, hasProfanity
}

func buildChirp(bodyText string) (newChirp Chirp) {
	// TODO: build chirp with structure of Chirp type. Needs a unique ID per chirp
	return newChirp
}

func getChirpsHandler(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()
	// FIXME: don't use createDB in this handler
	db, err := database.CreateDB("./database.json")
	if err != nil {
		log.Fatal(err)
	}

	chirps, err := db.GetChirps()
	if err != nil {
		log.Fatal(err)
	}

	log.Printf("getchirpshandler chirps: %v", chirps)

}
