package services

import (
	"encoding/json"
	"log"

	"github.com/boltdb/bolt"
)

type SearchedData map[string]int

func SearchWordInDB(word []byte, db *bolt.DB) SearchedData {

	tx, err := db.Begin(true)
	if err != nil {
		log.Fatalf("Error starting transaction: %v\n", err)
	}
	defer tx.Rollback()

	termB := tx.Bucket([]byte("terms"))
	if termB == nil {
		log.Fatalf("Bucket does not exist\n")
	}

	termV := termB.Get(word)
	if err != nil {
		return nil
	}

	searchedData := make(SearchedData)

	err = json.Unmarshal(termV, &searchedData)
	if err != nil {
		log.Fatalf("Error parsing json: %v\n", err)
	}

	return searchedData
}
