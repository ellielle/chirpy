package database

import (
	"encoding/json"
	"errors"
	"log"
	"os"
	"sync"
)

type Chirp struct {
	id   int
	body string
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
	err := db.ensureDB()
	if err != nil {
		log.Printf("Error in ensureDB: %s", err.Error())
	}
	return db, nil
}

// CreateChirp creates a new chirp and saves it to disk
//func (db *DB) CreateChirp(body string) ([]Chirp, error) {}

// GetChirps returns all chirps in the database
func (db *DB) GetChirps() ([]Chirp, error) {
	db.mu.Lock()
	defer db.mu.Unlock()
	log.Print("starting file read")
	fi, err := os.ReadFile(db.path)
	if err != nil {
		log.Print(err)
		return nil, errors.New(err.Error())
	}

	log.Print("file read...")
	chirps := Chirp{}
	chirp_err := json.Unmarshal(fi, &chirps)
	if chirp_err != nil {
		return nil, chirp_err
	}

	log.Printf("chirps %v: ", chirps)
	log.Print("passed chirps")
	return nil, nil
}

// ensureDB creates a new database file if it doesn't exist
func (db *DB) ensureDB() error {
	_, err := os.ReadFile(db.path)
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			// initiaize 'database' JSON file
			birdSeed := &DBStructure{Chirps: map[int]Chirp{}}
			data, _ := json.Marshal(birdSeed)
			os.WriteFile(db.path, data, 0600)
			return nil
		}
		return err
	}
	return nil
}

// loadDB reads the database file into memory
//func (db *DB) loadDB() (DBStructure, error) {}

// writeDB writes the database file to disk
//func (db *DB) writeDB(dbStructure DBStructure) error {}
