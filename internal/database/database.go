package database

import (
	"encoding/json"
	"errors"
	"log"
	"os"
	"sync"
)

type Chirp struct {
	Id   int    `json:"id"`
	Body string `json:"body"`
}

type DB struct {
	path string
	mu   *sync.RWMutex
}

type DBStructure struct {
	Chirps map[int]Chirp `json:"chirps"`
}

// CreateDB returns a new database connection
// and creates the database file if it doesn't exist
func CreateDB(path string) (*DB, error) {
	db := &DB{
		path: path,
		mu:   &sync.RWMutex{},
	}
	db.mu.Lock()
	defer db.mu.Unlock()

	// call ensureDB to check for database, and create one if it doesn't exist
	err := db.ensureDB()
	if err != nil && errors.Is(err, os.ErrNotExist) {

		// initiaize 'database' JSON file
		birdSeed := &DBStructure{Chirps: map[int]Chirp{}}
		data, _ := json.Marshal(birdSeed)
		os.WriteFile(db.path, data, 0600)
	}

	// something else went wrong
	if err != nil {
		log.Fatal(err.Error())
	}

	return db, nil
}

// ensureDB returns an error if the database does not exist yet
func (db *DB) ensureDB() error {
	_, err := os.ReadFile(db.path)
	if err != nil {
		return err
	}
	return nil
}

// loadDB reads the database file into memory as a DBStructure struct
func (db *DB) loadDB() (DBStructure, error) {
	db.mu.Lock()
	defer db.mu.Unlock()

	dat, err := os.ReadFile(db.path)
	if err != nil {
		return DBStructure{}, err
	}

	// Unmarshal the json data into DBStructure
	var data DBStructure
	json.Unmarshal(dat, &data)

	return data, nil
}

// CreateChirp creates a new chirp and saves it to disk
func (db *DB) CreateChirp(body string) error {
	nextID := db.getNextID()
	dat, err := db.loadDB()
	if err != nil {
		return err
	}

	chirp := Chirp{
		Id:   nextID,
		Body: body,
	}

	chirpMap := map[int]Chirp{}
	for i, c := range dat.Chirps {
		chirpMap[i] = Chirp{
			Id:   c.Id,
			Body: c.Body,
		}
	}
	chirpMap[nextID] = chirp

	chirpStructure := &DBStructure{
		Chirps: chirpMap,
	}

	db.writeDB(*chirpStructure)
	return nil
}

// GetChirps returns all chirps in the database
func (db *DB) GetChirps() ([]Chirp, error) {
	// TODO: test that this works once database is populated
	chirpSlice, err := db.getChirpsSlice()
	if err != nil {
		return nil, err
	}

	return chirpSlice, nil
}

// writeDB writes the database file to disk
func (db *DB) writeDB(dbStructure DBStructure) error {
	db.mu.Lock()
	defer db.mu.Unlock()

	jsonData, err := json.Marshal(dbStructure)
	if err != nil {
		return err
	}

	os.WriteFile(db.path, jsonData, 0600)
	return nil
}

// getID returns the next available ID in the database
func (db *DB) getNextID() int {
	dbSlice, err := db.getChirpsSlice()
	if err != nil {
		log.Fatal("getNextID failed")
	}

	return len(dbSlice)
}

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
