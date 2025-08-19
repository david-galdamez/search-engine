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

	if metaV == nil {
		counter = 1
	} else {
		counter, err := strconv.ParseUint(string(metaV), 10, 64)
		if err != nil {
			return err
		}
		counter++
	}

	err = metaB.Put([]byte("N"), []byte(strconv.FormatUint(counter, 10)))
	if err != nil {
		return err
	}

	if err := tx.Commit(); err != nil {
		return err
	}

	return nil
}
