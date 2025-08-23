package services

import (
	"encoding/json"
	"fmt"

	"github.com/boltdb/bolt"
	"github.com/david-galdamez/search-engine/models"
)

func SearchWordInDB(word []byte, db *bolt.DB) (*models.Terms, error) {

	tx, err := db.Begin(true)
	if err != nil {
		return nil, err
	}
	defer tx.Rollback()

	termB := tx.Bucket([]byte("terms"))

	termV := termB.Get(word)
	if termV == nil {
		return nil, fmt.Errorf("word not found")
	}

	searchedData := models.Terms{}

	err = json.Unmarshal(termV, &searchedData)
	if err != nil {
		return nil, err
	}

	return &searchedData, err
}
