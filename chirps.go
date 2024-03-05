package main

import (
	"encoding/json"
	"log"
	"net/http"
	"strings"

	database "github.com/ellielle/chirpy/internal/database"
)

type Chirp struct {
	Id   int
	Body string
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
		Body string `json:"body"`
	}

	// Create a new JSON decoder and check the validity of the JSON from the Request body
	decoder := json.NewDecoder(r.Body)
	params := parameters{}
	err := decoder.Decode(&params)
	if err != nil {
		respondWithError(w, 400, "Bad Request")
		return
	}

	// Disallow any chirps longer than 140 characters
	if len(params.Body) > 140 {
		respondWithError(w, 400, "Chirp is too long")
		return
	}

	// Remove profanity because this is a Christian Minecraft server
	// Respond with error if database connection fails
	cleanedBody, hasProfanity := validateNoProfanity(params.Body)
	db, err := createDBConnection()
	if err != nil {
		respondWithError(w, 500, err.Error())
	}

	// Create a new chirp with the body and save it to database
	// Then respond with appropriate response
	db.CreateChirp(cleanedBody)

	// FIXME: should return with 'id' field also
	if hasProfanity {
		respondWithJSON(w, 200, returnBody{Body: cleanedBody})
		return
	}
	respondWithJSON(w, 200, returnBody{Body: params.Body})
}

// Checks for profanity usage by looping over theProfane slice and checking the words against a lower cased params.Body
func validateNoProfanity(bodyText string) (cleanedText string, hasProfanity bool) {
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

// Gets all chirps in database and returns them in ascending order
func getChirpsHandler(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()
	db, err := createDBConnection()
	if err != nil {
		log.Fatal(err)
	}

	chirps, err := db.GetChirps()
	if err != nil {
		log.Fatal(err)
	}

	respondWithJSON(w, 200, chirps)
}

// Creates a 'connection' to the database using a pointer to the JSON database in memory
func createDBConnection() (*database.DB, error) {
	db, err := database.CreateDB("./database.json")
	if err != nil {
		return nil, err
	}
	return db, nil
}
