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

	for _, word := range wordsIterator {
		index := models.Terms{}

		termV := termB.Get([]byte(word))
		if termV == nil {
			continue
		}
		err := json.Unmarshal(termV, &index)
		if err != nil {
			return err
		}

		delete(index.Docs, doc.Id)
		index.DF--

		if index.Docs == nil {
			if err := termB.Delete([]byte(word)); err != nil {
				return err
			}
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
