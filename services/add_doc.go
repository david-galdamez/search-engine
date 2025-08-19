package services

import (
	"encoding/json"
	"log"

	"github.com/boltdb/bolt"
)

type Document struct {
	Id     string  `json:"id"`
	Title  string  `json:"title"`
	Length int     `json:"length"`
	Text   string  `json:"content"`
	Url    *string `json:"url,omitempty"`
}

func AddDoc(db *bolt.DB, doc *Document) {
	tx, err := db.Begin(true)
	if err != nil {
		log.Fatalf("Error starting transaction: %v\n", err)
	}
	defer tx.Rollback()

	if doc == nil {
		log.Fatalf("Document can not be empty")
	}

	docJson, err := json.Marshal(*doc)
	if err != nil {
		log.Fatalf("Error parsing to json: %v\n", err)
	}

	docsB := tx.Bucket([]byte("docs"))
	if docsB == nil {
		log.Fatalf("Bucket does not exist\n")
	}

	err = docsB.Put([]byte(doc.Id), docJson)
	if err != nil {
		log.Fatalf("Error inserting into db: %v\n", err)
	}

	if err := tx.Commit(); err != nil {
		log.Fatalf("Error committing into db: %v\n", err)
	}
}
