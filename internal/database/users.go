package database

// Creates a new User and saves it to disk
func (db *DB) CreateUser(body string, ch chan<- int) error {
	// TODO: get next USER id, maybe get next ID needs to be generic?
	nextID := db.getNextUserID()
	dat, err := db.loadDB()
	if err != nil {
		return err
	}

	user := User{
		Email: body,
	}

	// Build a map of [int]User and add the new User to it
	chirpMap, userMap := generateDataMap(&dat)
	userMap[nextID] = user

	userStructure := &DBStructure{
		Chirps: chirpMap,
		Users:  userMap,
	}

	ch <- nextID
	db.writeDB(*userStructure)
	close(ch)
	return nil
}

// Returns all chirps as a Slice for easier manipulation
func (db *DB) getUsersSlice() ([]User, error) {
	data, err := db.loadDB()
	if err != nil {
		return nil, err
	}

	var userSlice []User
	for _, user := range data.Users {
		userSlice = append(userSlice, user)
	}

	return userSlice, nil
}
