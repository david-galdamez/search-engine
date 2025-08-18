package services

import (
	"encoding/json"
	"log"
	"strings"
	"unicode"

	"github.com/boltdb/bolt"
)

func AddTextToDB(docId, text string, db *bolt.DB) {
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
		log.Fatalf("Error starting transactions: %v\n", err)
	}
	defer tx.Rollback()

	for word := range wordsIterator {
		if len(word) <= 2 {
			continue
		}

		termB := tx.Bucket([]byte("terms"))
		termV := termB.Get([]byte(word))
		if termV == nil {
			termB.Put([]byte(word), []byte("{}"))
		}
		index := make(map[string]int)

		err := json.Unmarshal(termB.Get([]byte(word)), &index)
		if err != nil {
			log.Fatalf("Error parsing json: %v\n", err)
		}
		index[docId]++

		data, err := json.Marshal(index)
		if err != nil {
			log.Fatalf("Error parsing to json: %v\n", err)
		}

		termB.Put([]byte(word), data)
	}

	if err := tx.Commit(); err != nil {
		log.Fatalf("Error commiting: %v\n", err)
	}
}
