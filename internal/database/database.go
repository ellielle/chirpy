package database

import (
	"encoding/json"
	"errors"
	"log"
	"os"
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

type User struct {
	Email string `json:"email"`
}

type DBStructure struct {
	Chirps map[int]Chirp `json:"chirps"`
	Users  map[int]User  `json:"users"`
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
		birdSeed := &DBStructure{Chirps: map[int]Chirp{}, Users: map[int]User{}}
		data, _ := json.Marshal(birdSeed)
		os.WriteFile(db.path, data, 0600)
	} else if err != nil {
		// Something else went wrong
		log.Fatal(err.Error())
	}

	return db, nil
}

func DebugWipeTestDatabase(path string) error {
	db, err := CreateDB(path)
	if err != nil {
		return err
	}
	dbErr := db.ensureDB()
	if dbErr != nil && errors.Is(dbErr, os.ErrNotExist) {
		return nil
	}
	os.Remove(path)

	return nil
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

// Temp function to check functionality
func (db *DB) getNextUserID() int {
	dbSlice, err := db.getUsersSlice()
	if err != nil {
		log.Fatal("getNextUserID failed")
	}

	return len(dbSlice) + 1
}

func generateDataMap(data *DBStructure) (map[int]Chirp, map[int]User) {
	chirpMap := map[int]Chirp{}
	for i, c := range data.Chirps {
		chirpMap[i] = Chirp{
			Body: c.Body,
			Id:   c.Id,
		}
	}

	userMap := map[int]User{}
	for i, u := range data.Users {
		userMap[i] = User{
			Email: u.Email,
		}
	}

	return chirpMap, userMap
}
