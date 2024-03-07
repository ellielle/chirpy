package database

import (
	"errors"

	"golang.org/x/crypto/bcrypt"
)

type User struct {
	Id       int    `json:"id"`
	Email    string `json:"email"`
	Password string `json:"password"`
}

var ErrInvalidLogin = errors.New("invalid login")

// Creates a new User and saves it to disk
func (db *DB) CreateUser(email, password string) (User, error) {
	dbStructure, err := db.loadDB()
	if err != nil {
		return User{}, err
	}

	_, err = getUserByEmail(email, &dbStructure)
	if err == nil {
		return User{}, errors.New("username taken")
	}

	// hash password with bcrypt
	hash, err := bcrypt.GenerateFromPassword([]byte(password), 14)

	// Create a new User with the next incremental ID
	nextID := len(dbStructure.Users) + 1
	user := User{
		Id:       nextID,
		Email:    email,
		Password: string(hash),
	}
	dbStructure.Users[nextID] = user
	err = db.writeDB(dbStructure)
	if err != nil {
		return User{}, nil
	}

	return user, nil
}

func (db *DB) LoginUser(email, password string) (User, error) {
	dbStructure, err := db.loadDB()
	if err != nil {
		return User{}, err
	}

	foundUser, err := getUserByEmail(email, &dbStructure)
	if err != nil {
		return User{}, err
	}

	// Compare hashes and return a valid response if they match
	// Return invalid login error inf mismatched
	err = bcrypt.CompareHashAndPassword([]byte(foundUser.Password), []byte(password))
	if err != nil {
		return User{}, ErrInvalidLogin
	}

	return foundUser, nil
}

func getUserByEmail(email string, dbStructure *DBStructure) (User, error) {
	// Find User in database, and return it
	// Will return an error if the user does not exist
	foundUser := User{}
	for _, user := range dbStructure.Users {
		if user.Email == email {
			foundUser = user
			break
		}
	}
	if foundUser == (User{}) {
		return User{}, ErrInvalidLogin
	}

	return foundUser, nil
}
