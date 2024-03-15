package main

import (
	"net/http"
	"strconv"
	"strings"

	auth "github.com/ellielle/chirpy/internal/auth"
)

func (cfg apiConfig) handlerChirpsDelete(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()
	chirpID := r.PathValue("chirpID")

	// Grab Authorization Bearer token from headers and then validate it
	headerToken, found := strings.CutPrefix(r.Header.Get("Authorization"), "Bearer ")
	if !found {
		respondWithError(w, http.StatusUnauthorized, "Authorization header missing")
		return
	}
	token, err := auth.ValidateJWT(headerToken, cfg.jwtSecret)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, err.Error())
		return
	}

	userID, err := auth.GetUserIDWithToken(*token)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, err.Error())
		return
	}

	// Convert string IDs to ints to be passed to DeleteChirp
	userIDInt, err := strconv.Atoi(userID)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, err.Error())
		return
	}
	chirpIDInt, err := strconv.Atoi(chirpID)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, err.Error())
		return
	}

	err = cfg.DB.DeleteChirp(chirpIDInt, userIDInt)
	if err != nil {
		respondWithError(w, http.StatusForbidden, "Unauthorized")
		return
	}

	respondWithJSON(w, http.StatusOK, "OK")
}
