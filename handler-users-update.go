package main

import (
	"encoding/json"
	"net/http"
	"strings"

	auth "github.com/ellielle/chirpy/internal/auth"
)

// Updates user email or password. User verifies themselves with a JWT token,
// and sends an email and / or password to attempt to update along with it.
func (cfg apiConfig) handlerUsersUpdate(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()
	type parameters struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}

	// Create a new JSON decoder and check the validity of the JSON from the Request body
	decoder := json.NewDecoder(r.Body)
	params := parameters{}
	err := decoder.Decode(&params)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Bad Request")
		return
	}

	// Ensure User's email and password are valid
	if params.Email != "" {
		err = validateEmail(params.Email)
		if err != nil {
			respondWithError(w, http.StatusBadRequest, err.Error())
			return
		}
	}
	if params.Password != "" {
		err = validatePassword(params.Password)
		if err != nil {
			respondWithError(w, http.StatusBadRequest, err.Error())
			return
		}
	}

	// Check that an Authorization header exists
	tokenHeader, found := strings.CutPrefix(r.Header.Get("Authorization"), "Bearer ")
	if !found {
		respondWithError(w, http.StatusBadRequest, "Invalid headers")
		return
	}

	// Verify User by validating their JWT
	userID, err := auth.ValidateJWT(tokenHeader, cfg.jwtSecret)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "Unauthorized")
		return
	}

	updatedUser, err := cfg.DB.UpdateUser(userID, params.Email, params.Password)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, err.Error())
		return
	}

	respondWithJSON(w, http.StatusOK, User{Id: updatedUser.Id, Email: updatedUser.Email})
}
