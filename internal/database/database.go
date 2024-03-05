package database

import (
	"encoding/json"
	"errors"
	"log"
	"os"
	"sort"
	"sync"
)

type Chirp struct {
	Body string `json:"body"`
	Id   int    `json:"id"`
}

type DB struct {
	path string
	mu   *sync.RWMutex
}

type DBStructure struct {
	Chirps map[int]Chirp `json:"chirps"`
}

// Returns a new database connection
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

		// Initiaize 'database' JSON file
		birdSeed := &DBStructure{Chirps: map[int]Chirp{}}
		data, _ := json.Marshal(birdSeed)
		os.WriteFile(db.path, data, 0600)
	} else if err != nil {
		// Something else went wrong
		log.Fatal(err.Error())
	}

	return db, nil
}

// Returns an error if the database does not exist yet
func (db *DB) ensureDB() error {
	_, err := os.ReadFile(db.path)
	if err != nil {
		return err
	}
	return nil
}

// Reads the database file into memory as a DBStructure struct
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

// Creates a new chirp and saves it to disk
func (db *DB) CreateChirp(body string, ch chan int) error {
	nextID := db.getNextID()
	dat, err := db.loadDB()
	if err != nil {
		return err
	}

	chirp := Chirp{
		Body: body,
		Id:   nextID,
	}

	chirpMap := map[int]Chirp{}
	for i, c := range dat.Chirps {
		chirpMap[i] = Chirp{
			Body: c.Body,
			Id:   c.Id,
		}
	}
	chirpMap[nextID] = chirp

	chirpStructure := &DBStructure{
		Chirps: chirpMap,
	}

	ch <- nextID
	db.writeDB(*chirpStructure)
	return nil
}

// Returns all chirps in the database in ascending order based on ID
func (db *DB) GetChirps() ([]Chirp, error) {
	chirpSlice, err := db.getChirpsSlice()
	if err != nil {
		return nil, err
	}

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

// Writes the database file to disk
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

// Returns the next available ID in the database
func (db *DB) getNextID() int {
	dbSlice, err := db.getChirpsSlice()
	if err != nil {
		log.Fatal("getNextID failed")
	}

	return len(dbSlice) + 1
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
