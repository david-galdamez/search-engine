package services

import (
	"encoding/json"
	"fmt"

	"github.com/boltdb/bolt"
)

type SearchedData map[string]int

func SearchWordInDB(word []byte, db *bolt.DB) (SearchedData, error) {

	tx, err := db.Begin(true)
	if err != nil {
		return nil, err
	}
	defer tx.Rollback()

	termB := tx.Bucket([]byte("terms"))

	termV := termB.Get(word)
	if err != nil {
		return nil, fmt.Errorf("word not found")
	}

	searchedData := make(SearchedData)

	err = json.Unmarshal(termV, &searchedData)
	if err != nil {
		return nil, err
	}

	return searchedData, err
}
