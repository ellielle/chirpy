package main

import (
	"encoding/json"
	"net/http"
)

type User struct {
	Email string
	Id    int
}

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
	// No password validation other than existence
	if params.Password == "" {
		respondWithError(w, http.StatusUnauthorized, "No password given")
		return
	}

	// Create a new user with the body and save it to database in a new goroutine
	user, err := cfg.DB.CreateUser(params.Email)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}
	respondWithJSON(w, http.StatusCreated, user)
}

func (cfg apiConfig) handlerUsersLogin(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()

	respondWithError(w, 500, "Method not finished")
	return
}
