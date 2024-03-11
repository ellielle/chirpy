package main

import (
	"encoding/json"
	"errors"
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

// TODO : made ExpiresIn omitempty so it doesn't return in responses, finish jwt functions
// and uncomment PUT method for /api/users
type User struct {
	Id        int    `json:"id"`
	Email     string `json:"email"`
	ExpiresIn int    `json:"expires_in_seconds,omitempty"`
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
	err = validateUser(params.Email, params.Password)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, err.Error())
		return
	}

	_ = createJWT("bleep", "blep")
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
	err = validateUser(params.Email, params.Password)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, err.Error())
		return
	}
	// Attempt to log user in, a failture will result in a 401 not authorized error
	user, err := cfg.DB.LoginUser(params.Email, params.Password)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, err.Error())
		return
	}
	respondWithJSON(w, http.StatusOK, User{Id: user.Id, Email: user.Email})
}

func validateUser(email, password string) error {
	// No password validation other than existence
	if password == "" {
		return ErrInvalidPassword
	}
	// Most minimum of requirements for an email
	if !strings.Contains(email, "@") {
		return ErrInvalidEmail
	}
	return nil
}

func createJWT(userInfo string, jwtSecret string) jwt.Token {
	// TODO: working on JWT signing
	// WARN: don't forget to use Harpoon
	// FIXME: https://pkg.go.dev/time#Unix figure out how to timestamp and convert seconds to / from
	// and convert User.ExpiresIn (seconds) to a time format jwt accepts
	claims := &jwt.RegisteredClaims{
		ExpiresAt: jwt.NewNumericDate(time.Unix(time.Now().Unix(), 0)),
		Issuer:    "chirpy",
	}
	log.Println("")
	log.Printf("time: %v", time.Now())
	log.Printf("time unix: %v", time.Now().Unix())
	log.Printf("claims: %v", claims)
	log.Println("")
	return jwt.Token{}
}
