package database

import (
	"errors"
	"sort"
)

// Creates a new chirp and saves it to disk
func (db *DB) CreateChirp(body string, ch chan<- int) error {
	nextID := db.getNextID()
	dat, err := db.loadDB()
	if err != nil {
		return err
	}

	// Build a map of [int]Chirp and add the newest Chirp to it
	chirp := Chirp{
		Body: body,
		Id:   nextID,
	}
	chirpMap, userMap := generateDataMap(&dat)
	chirpMap[nextID] = chirp
	chirpStructure := &DBStructure{
		Chirps: chirpMap,
		Users:  userMap,
	}

	ch <- nextID
	db.writeDB(*chirpStructure)
	close(ch)
	return nil
}

// Returns all chirps in the database in ascending order based on ID
func (db *DB) GetChirps() ([]Chirp, error) {
	chirpSlice, err := db.getChirpsSlice()
	if err != nil {
		return nil, err
	}

	// And sort the slice to make it pretty
	sort.Slice(chirpSlice, func(i, j int) bool {
		return chirpSlice[i].Id < chirpSlice[j].Id
	})
	return chirpSlice, nil
}

func (db *DB) GetSingleChirp(chirpID int) (Chirp, error) {
	chirpSlice, err := db.getChirpsSlice()
	if err != nil {
		return Chirp{}, err
	}

	for i, chirp := range chirpSlice {
		if chirpID == chirp.Id {
			return chirpSlice[i], nil
		}
	}
	return Chirp{}, errors.New("Chirp not found")
}

// Returns all chirps as a Slice for easier manipulation
func (db *DB) getChirpsSlice() ([]Chirp, error) {
	data, err := db.loadDB()
	if err != nil {
		return nil, err
	}

	var chirpSlice []Chirp
	for _, chirp := range data.Chirps {
		chirpSlice = append(chirpSlice, chirp)
	}
	return chirpSlice, nil
}
