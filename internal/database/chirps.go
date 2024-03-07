package database

import (
	"errors"
)

type Chirp struct {
	Id   int    `json:"id"`
	Body string `json:"body"`
}

// Creates a new chirp and saves it to disk
func (db *DB) CreateChirp(body string) (Chirp, error) {
	dbStructure, err := db.loadDB()
	if err != nil {
		return Chirp{}, err
	}

	// Create a new Chirp with the next incremental ID
	nextID := len(dbStructure.Chirps) + 1
	chirp := Chirp{
		Id:   nextID,
		Body: body,
	}
	dbStructure.Chirps[nextID] = chirp
	err = db.writeDB(dbStructure)
	if err != nil {
		return Chirp{}, nil
	}

	return chirp, nil
}

// Returns all chirps in the database in ascending order based on ID
func (db *DB) GetChirps() ([]Chirp, error) {
	dbStructure, err := db.loadDB()
	if err != nil {
		return nil, err
	}

	// Create a []Chirp slice and append all current chirps to it
	chirpSlice := make([]Chirp, 0, len(dbStructure.Chirps))
	for _, chirp := range dbStructure.Chirps {
		chirpSlice = append(chirpSlice, chirp)
	}
	return chirpSlice, nil
}

func (db *DB) GetChirp(chirpID int) (Chirp, error) {
	dbStructure, err := db.loadDB()
	if err != nil {
		return Chirp{}, err
	}

	chirp, ok := dbStructure.Chirps[chirpID]
	if !ok {
		return Chirp{}, errors.New("Chirp not found")
	}
	return chirp, nil
}