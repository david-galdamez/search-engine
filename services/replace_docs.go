package services

import (
	"encoding/json"

	"github.com/boltdb/bolt"
	"github.com/david-galdamez/search-engine/models"
	"github.com/david-galdamez/search-engine/utils"
)

func DeleteDocTerms(db *bolt.DB, doc *models.Document) error {

	wordsIterator := utils.Tokenizer(doc.Content)

	tx, err := db.Begin(true)
	if err != nil {
		return err
	}
	defer tx.Rollback()

	termB := tx.Bucket([]byte("terms"))

	for word := range wordsIterator {
		index := make(map[string]int)

		termV := termB.Get([]byte(word))
		err := json.Unmarshal(termV, &index)
		if err != nil {
			return err
		}

		delete(index, doc.Id)
		if index == nil {
			termB.Delete([]byte(word))
			continue
		}

		newIndex, err := json.Marshal(index)
		if err != nil {
			return err
		}

		err = termB.Put([]byte(word), newIndex)
		if err != nil {
			return err
		}
	}

	if err := tx.Commit(); err != nil {
		return err
	}

	return nil
}
