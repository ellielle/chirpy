package main

import (
	"encoding/json"
	"net/http"
)

type User struct {
	Email string
}

func validateUserHandler(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()

	type parameters struct {
		Email string `json:"email"`
	}

	type returnValid struct {
		Valid bool `json:"valid"`
	}

	type returnBody struct {
		Email string `json:"email"`
		Id    int    `json:"id"`
	}

	// Create a new JSON decoder and check the validity of the JSON from the Request body
	decoder := json.NewDecoder(r.Body)
	params := parameters{}
	err := decoder.Decode(&params)
	if err != nil {
		respondWithError(w, 400, "Bad Request")
		return
	}

	// Read database.json into memory and give access to the db pointer
	db, err := createDBConnection()
	if err != nil {
		respondWithError(w, 500, err.Error())
		return
	}

	// Create a new user with the body and save it to database in a new goroutine
	// Then respond with appropriate response once the new ID number is received on the channel 'ch'
	ch := make(chan int)
	go db.CreateUser(params.Email, ch)
	newID, ok := <-ch
	if !ok {
		respondWithError(w, 500, "Internal Server Error")
		return
	}
	respondWithJSON(w, 201, returnBody{Email: params.Email, Id: newID})
}
