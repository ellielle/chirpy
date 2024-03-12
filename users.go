package main

import (
	"encoding/json"
	"errors"
	"net/http"
	"strings"

	auth "github.com/ellielle/chirpy/internal/auth"
)

type User struct {
	Id    int    `json:"id"`
	Email string `json:"email"`
	Token string `json:"token,omitempty"`
}

var ErrInvalidPassword = errors.New("password missing or invalid")
var ErrInvalidEmail = errors.New("email is invalid")

func (cfg apiConfig) handlerUsersCreate(w http.ResponseWriter, r *http.Request) {
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

	// Create a new user with the body and save it to database in a new goroutine
	user, err := cfg.DB.CreateUser(params.Email, params.Password)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}
	respondWithJSON(w, http.StatusCreated, User{Id: user.Id, Email: user.Email})
}

func (cfg apiConfig) handlerUsersLogin(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()
	type parameters struct {
		Email     string `json:"email"`
		Password  string `json:"password"`
		ExpiresIn int    `json:"expires_in_seconds"`
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
		respondWithError(w, http.StatusUnauthorized, "Bad Request")
		return
	}

	// Create a JWT to be sent back to the user in the response
	token := auth.CreateJWT(auth.User{Id: user.Id, Email: user.Email, Password: user.Password}, cfg.jwtSecret, params.ExpiresIn)
	respondWithJSON(w, http.StatusOK, User{Id: user.Id, Email: user.Email, Token: token})
}

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

// Validate User's email. For now, it's a basic check
func validateEmail(email string) error {
	// Most minimum of requirements for an email
	if !strings.Contains(email, "@") {
		return ErrInvalidEmail
	}

	return nil
}

// Validate User's password. No real password rules other than not being empty
func validatePassword(password string) error {
	// No password validation other than existence
	if password == "" {
		return ErrInvalidPassword
	}
	return nil
}
