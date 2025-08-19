package services

import (
	"log"
	"strconv"

	"github.com/boltdb/bolt"
)

func IncrementDocCounter(db *bolt.DB) {
	tx, err := db.Begin(true)
	if err != nil {
		log.Fatalf("Error opening database: %v\n", err)
	}
	defer tx.Rollback()

	metaB := tx.Bucket([]byte("meta"))
	if metaB == nil {
		log.Fatalf("Bucket does not exist\n")
	}

	metaV := metaB.Get([]byte("N"))
	var counter uint64

	if metaV == nil {
		counter = 1
	} else {
		counter, err := strconv.ParseUint(string(metaV), 10, 64)
		if err != nil {
			log.Fatalf("Error parsing to uint: %v\n", err)
		}
		counter++
	}

	err = metaB.Put([]byte("N"), []byte(strconv.FormatUint(counter, 10)))
	if err != nil {
		log.Fatalf("Error inserting into database: %v\n", err)
	}

	if err := tx.Commit(); err != nil {
		log.Fatalf("Error committing: %v\n", err)
	}
}
