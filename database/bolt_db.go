package database

import "github.com/boltdb/bolt"

func GetDB() (*bolt.DB, error) {
	db, err := bolt.Open("search-engine.db", 0600, nil)
	if err != nil {
		return nil, err
	}

	return db, nil
}
