package services

import (
	"encoding/json"
	"fmt"
	"log"
	"strings"
	"unicode"

	"github.com/boltdb/bolt"
)

func AddTextToDB(docId, title, text string, db *bolt.DB) error {
	//trims punctuations and split by spaces
	cleanText := strings.Map(func(r rune) rune {
		if unicode.IsPunct(r) {
			return -1
		}
		return r
	}, text)

	wordsIterator := strings.FieldsSeq(strings.ToLower(cleanText))

	tx, err := db.Begin(true)
	if err != nil {
		return err
	}
	defer tx.Rollback()

	for word := range wordsIterator {
		if len(word) <= 2 {
			continue
		}

		termB := tx.Bucket([]byte("terms"))
		if termB == nil {
			return fmt.Errorf("bucket terms does not exist")
		}

		termV := termB.Get([]byte(word))
		if termV == nil {
			err := termB.Put([]byte(word), []byte("{}"))
			if err != nil {
				return err
			}
		}
		index := make(map[string]int)

		err := json.Unmarshal(termB.Get([]byte(word)), &index)
		if err != nil {
			return err
		}
		index[docId]++

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
