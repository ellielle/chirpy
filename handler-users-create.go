package main

import (
	"encoding/json"
	"errors"
	"net/http"
	"strings"
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
