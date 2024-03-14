package main

import (
	"net/http"
	"sort"
	"strconv"
)

// Gets all chirps in database and returns them in ascending order
func (cfg apiConfig) handlerChirpsGetAll(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()

	chirps, err := cfg.DB.GetChirps()
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	// Sort chirps by id before sending response
	sort.Slice(chirps, func(i, j int) bool {
		return chirps[i].Id < chirps[j].Id
	})

	respondWithJSON(w, http.StatusOK, chirps)
}

// Gets a single chirp by ID and returns it
func (cfg apiConfig) handlerChirpsGet(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()

	chirp := r.PathValue("chirpID")
	chirpID, err := strconv.Atoi(chirp)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid chirp ID")
		return
	}

	foundChirp, err := cfg.DB.GetChirp(chirpID)
	if err != nil {
		respondWithError(w, http.StatusNotFound, "Not Found")
		return
	}
	respondWithJSON(w, http.StatusOK, foundChirp)
}
