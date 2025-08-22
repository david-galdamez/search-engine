package services

import (
	"fmt"
	"strconv"

	"github.com/boltdb/bolt"
)

func IncrementDocCounter(db *bolt.DB) error {
	tx, err := db.Begin(true)
	if err != nil {
		return err
	}
	defer tx.Rollback()

	metaB := tx.Bucket([]byte("meta"))
	if metaB == nil {
		return fmt.Errorf("bucket meta does not exist")
	}

	metaV := metaB.Get([]byte("N"))
	var counter uint64

	counter, err = strconv.ParseUint(string(metaV), 10, 64)
	if err != nil {
		return err
	}
	counter++

	err = metaB.Put([]byte("N"), []byte(strconv.FormatUint(counter, 10)))
	if err != nil {
		return err
	}

	if err := tx.Commit(); err != nil {
		return err
	}

	return nil
}

func GetDocumentCounter(db *bolt.DB) (*int, error) {
	var counter int

	err := db.View(func(tx *bolt.Tx) error {
		metaB := tx.Bucket([]byte("meta"))

		metaV := metaB.Get([]byte("N"))
		if metaV == nil {
			return fmt.Errorf("counter was not initialized")
		}

		parsedCounter, err := strconv.ParseInt(string(metaV), 10, 64)
		if err != nil {
			return err
		}

		counter = int(parsedCounter)

		return nil
	})
	if err != nil {
		return nil, err
	}

	return &counter, nil
}
