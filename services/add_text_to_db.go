package services

import (
	"encoding/json"
	"fmt"
	"log"

	"github.com/boltdb/bolt"
	"github.com/david-galdamez/search-engine/models"
	"github.com/david-galdamez/search-engine/utils"
)

func AddTextToDB(docId, title, text string, db *bolt.DB) error {

	wordsIterator := utils.Tokenizer(text)

	tx, err := db.Begin(true)
	if err != nil {
		return err
	}
	defer tx.Rollback()

	for _, word := range wordsIterator {
		if len(word) <= 2 {
			continue
		}

		termB := tx.Bucket([]byte("terms"))
		if termB == nil {
			return fmt.Errorf("bucket terms does not exist")
		}

		termV := termB.Get([]byte(word))
		var index models.Terms
		if termV == nil {
			index = models.Terms{DF: 0, Docs: make(models.DocsFrequency)}
		} else {
			err := json.Unmarshal(termB.Get([]byte(word)), &index)
			if err != nil {
				return err
			}

			if index.Docs == nil {
				index.Docs = make(models.DocsFrequency)
			}
		}

		index.Docs[docId]++
		index.DF = len(index.Docs)

		data, err := json.Marshal(index)
		if err != nil {
			return err
		}

		err = termB.Put([]byte(word), data)
		if err != nil {
			return err
		}
	}

	if err := tx.Commit(); err != nil {
		return err
	}
	log.Printf("Document with id: %v successfully added\n", docId)

	return nil
}
