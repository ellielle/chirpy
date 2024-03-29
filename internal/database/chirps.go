package database

import (
	"errors"
	"strconv"
)

type Chirp struct {
	Id       int    `json:"id"`
	Body     string `json:"body"`
	AuthorId int    `json:"author_id"`
}

var ErrChirpNotFound = errors.New("Chirp not found")
var ErrUnauthorized = errors.New("Unauthorized")

// Creates a new chirp and saves it to disk
func (db *DB) CreateChirp(body, id string) (Chirp, error) {
	dbStructure, err := db.loadDB()
	if err != nil {
		return Chirp{}, err
	}
	userID, err := strconv.Atoi(id)
	if err != nil {
		return Chirp{}, err
	}

	// Create a new Chirp with the next incremental ID
	nextID := len(dbStructure.Chirps) + 1
	chirp := Chirp{
		Id:       nextID,
		Body:     body,
		AuthorId: userID,
	}
	dbStructure.Chirps[nextID] = chirp
	err = db.writeDB(dbStructure)
	if err != nil {
		return Chirp{}, nil
	}

	return chirp, nil
}

// Returns all chirps in the database in ascending order based on ID
func (db *DB) GetChirps(searchByAuthorID string) ([]Chirp, error) {
	dbStructure, err := db.loadDB()
	if err != nil {
		return nil, err
	}

	// Create a []Chirp slice and append all current chirps to it
	chirpSlice := make([]Chirp, 0, len(dbStructure.Chirps))
	for _, chirp := range dbStructure.Chirps {
		if searchByAuthorID == "" {
			chirpSlice = append(chirpSlice, chirp)
			continue
		}

		authorID, err := strconv.Atoi(searchByAuthorID)
		if err != nil {
			return nil, err
		}

		// Only get chirps from author with optional ID passed in
		if chirp.AuthorId == authorID {
			chirpSlice = append(chirpSlice, chirp)
		}
	}
	return chirpSlice, nil
}

// Get a specific Chirp from the database
func (db *DB) GetChirp(chirpID int) (Chirp, error) {
	dbStructure, err := db.loadDB()
	if err != nil {
		return Chirp{}, err
	}

	chirp, ok := dbStructure.Chirps[chirpID]
	if !ok {
		return Chirp{}, ErrChirpNotFound
	}
	return chirp, nil
}

func (db *DB) DeleteChirp(chirpID, authorID int) error {
	dbStructure, err := db.loadDB()
	if err != nil {
		return err
	}

	chirp, ok := dbStructure.Chirps[chirpID]
	if !ok {
		return ErrChirpNotFound
	}

	// User must be the owner of the chirp to delete it
	if chirp.AuthorId != authorID {
		return ErrUnauthorized
	}

	delete(dbStructure.Chirps, chirpID)
	err = db.writeDB(dbStructure)
	if err != nil {
		return err
	}

	return nil
}
