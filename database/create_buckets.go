package database

import (
	"log"

	"github.com/boltdb/bolt"
)

func CreateBuckets() {
	db, err := GetDB()
	if err != nil {
		log.Fatalf("Error opening database: %v\n", err)
	}
	defer db.Close()

	err = db.Update(func(tx *bolt.Tx) error {

		_, err := tx.CreateBucketIfNotExists([]byte("terms"))
		if err != nil {
			return err
		}

		_, err = tx.CreateBucketIfNotExists([]byte("docs"))
		if err != nil {
			return err
		}

		b, err := tx.CreateBucketIfNotExists([]byte("meta"))
		if err != nil {
			return err
		}

		err = b.Put([]byte("N"), []byte("0"))
		if err != nil {
			return err
		}

		return nil
	})
}
