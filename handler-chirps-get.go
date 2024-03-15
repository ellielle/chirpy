package main

import (
	"net/http"
	"sort"
	"strconv"
)

// Gets all chirps in database and returns them in ascending order
func (cfg apiConfig) handlerChirpsGetAll(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()

	// Check for optional author_id query parameter
	authorID := r.URL.Query().Get("author_id")
	// Check for optional sort query parameter
	// sort can either be "asc" or "desc"
	// ascending is the default if no parameter is provided
	sortBy := r.URL.Query().Get("sort")

	// If an authorID was passed in, only chirps from that author will be returned
	// If authorID is "", all chirps will be returned
	chirps, err := cfg.DB.GetChirps(authorID)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	// Sort chirps by id in ascending order before sending response
	if sortBy == "" || sortBy == "asc" {
		sort.Slice(chirps, func(i, j int) bool {
			return chirps[i].Id < chirps[j].Id
		})
	}
	if sortBy == "desc" {
		sort.Slice(chirps, func(i, j int) bool {
			return chirps[i].Id > chirps[j].Id
		})
	}

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
