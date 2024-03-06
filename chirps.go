package main

import (
	"encoding/json"
	"net/http"
	"strconv"
	"strings"

	"github.com/go-chi/chi/v5"

	database "github.com/ellielle/chirpy/internal/database"
)

type Chirp struct {
	Body string
	Id   int
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
		Id   int    `json:"id"`
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

	// Read database.json into memory and give access to the db pointer
	db, err := createDBConnection()
	if err != nil {
		respondWithError(w, 500, err.Error())
		return
	}

	// Create a new chirp with the body and save it to database in a new goroutine
	// Then respond with appropriate response once the new ID number is received on the channel 'ch'
	ch := make(chan int)
	if hasProfanity {
		go db.CreateChirp(cleanedBody, ch)
		newID, ok := <-ch
		if !ok {
			respondWithError(w, 500, "Internal Server Error")
			return
		}
		respondWithJSON(w, 201, returnBody{Body: cleanedBody, Id: newID})
		return
	}
	go db.CreateChirp(params.Body, ch)
	newID := <-ch
	respondWithJSON(w, 201, returnBody{Body: params.Body, Id: newID})
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
		respondWithError(w, 500, err.Error())
		return
	}

	chirps, err := db.GetChirps()
	if err != nil {
		respondWithError(w, 500, err.Error())
		return
	}
	respondWithJSON(w, 200, chirps)
}

// Gets a single chirp by ID and returns it
func getSingleChirpHandler(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()
	chirp := chi.URLParam(r, "chirpID")
	chirpID, err := strconv.Atoi(chirp)
	if err != nil {
		respondWithError(w, 400, "Bad Request")
		return
	}

	db, err := createDBConnection()
	foundChirp, err := db.GetSingleChirp(chirpID)
	if err != nil {
		respondWithError(w, 404, "Not Found")
		return
	}
	respondWithJSON(w, 200, foundChirp)
}

// Creates a 'connection' to the database using a pointer to the JSON database in memory
func createDBConnection() (*database.DB, error) {
	db, err := database.CreateDB("./database.json")
	if err != nil {
		return nil, err
	}
	return db, nil
}
