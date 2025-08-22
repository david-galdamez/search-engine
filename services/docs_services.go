package services

import (
	"encoding/json"
	"fmt"

	"github.com/boltdb/bolt"
	"github.com/david-galdamez/search-engine/models"
)

func AddDoc(db *bolt.DB, doc *models.Document) error {
	tx, err := db.Begin(true)
	if err != nil {
		return err
	}
	defer tx.Rollback()

	if doc == nil {
		return fmt.Errorf("doc can not be empty")
	}

	docJson, err := json.Marshal(*doc)
	if err != nil {
		return err
	}

	docsB := tx.Bucket([]byte("docs"))
	if docsB == nil {
		return fmt.Errorf("bucket docs does not exist")
	}

	err = docsB.Put([]byte(doc.Id), docJson)
	if err != nil {
		return err
	}

	if err := tx.Commit(); err != nil {
		return err
	}

	return nil
}

func GetDoc(docId []byte, db *bolt.DB) (*models.Document, error) {

	doc := models.Document{}

	err := db.View(func(tx *bolt.Tx) error {
		docB := tx.Bucket([]byte("docs"))

		docV := docB.Get(docId)
		if docV == nil {
			return fmt.Errorf("document not found")
		}

		err := json.Unmarshal(docV, &doc)
		if err != nil {
			return err
		}

		return nil
	})
	if err != nil {
		return nil, err
	}

	return &doc, nil
}
