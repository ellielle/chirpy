package database

import (
	"errors"
	"strconv"
	"time"

	"github.com/golang-jwt/jwt/v5"

	auth "github.com/ellielle/chirpy/internal/auth"
)

var ErrUserNotFound = errors.New("User not found")
var ErrTokenRevoked = errors.New("Refresh token is revoked")
var ErrNoUserFound = errors.New("No user found")

// Takes a stringified version of a refresh token and adds it to the database as revoked, along with a timestamp
func (db *DB) RevokeToken(token string) error {
	dbStructure, err := db.loadDB()
	if err != nil {
		return err
	}

	err = tokenRevoker(token, &dbStructure)
	if err != nil {
		return err
	}

	err = db.writeDB(dbStructure)
	if err != nil {
		return err
	}
	return nil
}

// Revoke old refresh token and save it with a timestamp
// Takes a refresh token, finds the user by ID lookup, generates and returns a new access token
func (db *DB) RefreshToken(token *jwt.Token, stringToken, jwtSecret string) (string, error) {
	dbStructure, err := db.loadDB()
	if err != nil {
		return "", err
	}

	// Check for revoked status on the token before proceeding
	err = tokenRevokedStatus(stringToken, &dbStructure)
	if err != nil {
		return "", err
	}

	user, err := getUserBySubjectID(token, &dbStructure)
	if err != nil {
		return "", ErrUserNotFound
	}

	accessToken, err := auth.CreateJWT(auth.User{Id: user.Id}, jwtSecret, true)
	if err != nil {
		return "", err
	}

	// Revoke old refresh token and save it with a timestamp
	err = tokenRevoker(stringToken, &dbStructure)
	if err != nil {
		return "", err
	}
	err = db.writeDB(dbStructure)
	if err != nil {
		return "", err
	}

	return accessToken, nil
}

// Returns an error if a refresh token has been revoked
func tokenRevokedStatus(token string, dbStructure *DBStructure) error {
	for revToken := range dbStructure.RevokedTokens {
		if token == revToken {
			return ErrTokenRevoked
		}
	}
	return nil
}

// Revokes a refresh token and adds it to the database with a timestamp
// A nil return is successful
func tokenRevoker(token string, dbStructure *DBStructure) error {
	dbStructure.RevokedTokens[token] = time.Now()
	return nil
}

// Gets a User's ID from a validated JWT's Claims' Subject and returns the User
func getUserBySubjectID(token *jwt.Token, dbStructure *DBStructure) (User, error) {
	id, err := token.Claims.GetSubject()
	if err != nil {
		return User{}, err
	}
	userID, err := strconv.Atoi(id)
	if err != nil {
		return User{}, err
	}

	for _, user := range dbStructure.Users {
		if user.Id == userID {
			return user, nil
		}
	}

	return User{}, ErrNoUserFound
}
