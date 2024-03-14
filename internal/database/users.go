package database

import (
	"errors"
	"strconv"

	auth "github.com/ellielle/chirpy/internal/auth"
)

type User struct {
	Id       int    `json:"id"`
	Email    string `json:"email"`
	Password string `json:"-"`
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
	hash, err := auth.HashPassword(password)

	// Create a new User with the next incremental ID
	nextID := len(dbStructure.Users) + 1
	user := User{
		Id:       nextID,
		Email:    email,
		Password: hash,
	}
	dbStructure.Users[nextID] = user
	err = db.writeDB(dbStructure)
	if err != nil {
		return User{}, nil
	}

	return user, nil
}

// Logs user in by email and password by matching hashed
// passwords
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
	err = auth.CheckPasswordHash(foundUser.Password, password)
	if err != nil {
		return User{}, ErrInvalidLogin
	}

	return foundUser, nil
}

// Update user email/password using an authentication token
// Both email and password are optional parameters to the
// PUT endpoint api/users
func (db *DB) UpdateUser(id string, updates ...string) (User, error) {
	if len(updates) == 0 {
		return User{}, errors.New("No information to update")
	}

	dbStructure, err := db.loadDB()
	if err != nil {
		return User{}, err
	}

	// ID will be in string format after being parsed from token
	intID, err := strconv.Atoi(id)
	if err != nil {
		return User{}, nil
	}

	foundUser, err := getUserById(intID, &dbStructure)
	if err != nil {
		return User{}, err
	}

	newEmail := foundUser.Email
	newPassword := foundUser.Password

	if len(updates) > 0 && updates[0] != "" {
		newEmail = updates[0]
	}
	if len(updates) > 1 && updates[1] != "" {
		newPassword, err = auth.HashPassword(updates[1])
		if err != nil {
			return User{}, err
		}
	}

	user := User{
		Id:       foundUser.Id,
		Email:    newEmail,
		Password: newPassword,
	}
	dbStructure.Users[foundUser.Id] = user
	err = db.writeDB(dbStructure)
	if err != nil {
		return User{}, nil
	}

	return user, nil
}

// Find User by ID when user supplies an auth token
func getUserById(id int, dbStructure *DBStructure) (User, error) {
	foundUser := User{}
	for _, user := range dbStructure.Users {
		if user.Id == id {
			foundUser = user
			break
		}
	}
	if foundUser == (User{}) {
		return User{}, ErrInvalidLogin
	}

	return foundUser, nil
}

// Find User in database, and return it
// Will return an error if the user does not exist
func getUserByEmail(email string, dbStructure *DBStructure) (User, error) {
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
