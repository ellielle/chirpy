package main

import (
	"encoding/json"
	"net/http"
	"strings"

	auth "github.com/ellielle/chirpy/internal/auth"
)

// Takes a user's email and password, and if valid, returns their email, id, JWT access token and JWT refresh token in a response
func (cfg apiConfig) handlerUsersLogin(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()
	type parameters struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}

	type response struct {
		User
		Token        string `json:"token"`
		RefreshToken string `json:"refresh_token"`
	}

	// Create a new JSON decoder and check the validity of the JSON from the Request body
	decoder := json.NewDecoder(r.Body)
	params := parameters{}
	err := decoder.Decode(&params)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, err.Error())
		return
	}

	// Ensure User's email and password are valid
	err = validateEmail(params.Email)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, err.Error())
		return
	}
	err = validatePassword(params.Password)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, err.Error())
		return
	}

	// Attempt to log user in, a failure will result in a 401 not authorized error
	user, err := cfg.DB.LoginUser(params.Email, params.Password)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, err.Error())
		return
	}

	// Create a JWT access token and refresh token to be sent back to the user in the response
	token, err := auth.CreateJWT(auth.User{Id: user.Id, Email: user.Email, Password: user.Password}, cfg.jwtSecret, true)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, err.Error())
	}

	refreshToken, err := auth.CreateJWT(auth.User{Id: user.Id, Email: user.Email, Password: user.Password}, cfg.jwtSecret, false)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, err.Error())
	}

	respondWithJSON(w, http.StatusOK, response{
		User: User{
			Id:    user.Id,
			Email: user.Email,
		},
		Token:        token,
		RefreshToken: refreshToken,
	})
}

// If refresh token is valid, responds with a new access token and newly created access token with a 1 hour expiration
// Note: The refresh token *should* be invalidated and a new one issued with the access token, but the assignment only wants an access token in the response
func (cfg apiConfig) handlerUsersRefresh(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()

	type response struct {
		AccessToken string `json:"token"`
	}

	// Grab Authorization Bearer token from headers
	headerToken, found := strings.CutPrefix(r.Header.Get("Authorization"), "Bearer ")
	if !found {
		respondWithError(w, http.StatusBadRequest, "Authorization header missing")
		return
	}

	if headerToken == "chirpy-access" {
		respondWithError(w, http.StatusBadRequest, "Access token used as refresh token")
		return
	}

	token, err := auth.ValidateJWT(headerToken, cfg.jwtSecret)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "Invalid token")
		return
	}

	// Check if refresh token has been revoked, and if not retrieve a new access token
	accessToken, err := cfg.DB.RefreshToken(token, headerToken, cfg.jwtSecret)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, err.Error())
		return
	}

	respondWithJSON(w, 200, response{AccessToken: accessToken})
}

// Revokes a refresh token and stores that token in the database as revoked, with a timestamp
func (cfg apiConfig) handlerTokensRevoke(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()

	// Grab Authorization Bearer token from headers
	headerToken, found := strings.CutPrefix(r.Header.Get("Authorization"), "Bearer ")
	if !found {
		respondWithError(w, http.StatusBadRequest, "Authorization header missing")
		return
	}

	if headerToken == "chirpy-access" {
		respondWithError(w, http.StatusBadRequest, "Access token used as refresh token")
		return
	}

	_, err := auth.ValidateJWT(headerToken, cfg.jwtSecret)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "Invalid token")
		return
	}

	err = cfg.DB.RevokeToken(headerToken)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	respondWithJSON(w, 200, "OK")
}
