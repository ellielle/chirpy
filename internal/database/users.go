package database

type User struct {
	Id    int    `json:"id"`
	Email string `json:"email"`
}

// Creates a new User and saves it to disk
func (db *DB) CreateUser(body string) (User, error) {
	dbStructure, err := db.loadDB()
	if err != nil {
		return User{}, err
	}

	// Create a new User with the next incremental ID
	nextID := len(dbStructure.Users)
	user := User{
		Id:    nextID,
		Email: body,
	}
	dbStructure.Users[nextID] = user
	err = db.writeDB(dbStructure)
	if err != nil {
		return User{}, nil
	}

	return user, nil
}
