package services

import (
	"encoding/json"
	"fmt"
	"log"
	"math"
	"strconv"
	"strings"

	"github.com/boltdb/bolt"
	"github.com/david-galdamez/search-engine/models"
)

func SearchCalculate(db *bolt.DB, word []byte) (models.DocsScore, error) {

	docScore := make(models.DocsScore)
	var docsCounter int
	wordSearch := models.Terms{DF: 0, Docs: make(models.DocsFrequency)}

	err := db.View(func(tx *bolt.Tx) error {

		termB := tx.Bucket([]byte("terms"))

		termV := termB.Get(word)
		if termV == nil {
			return fmt.Errorf("word not registered")
		}

		err := json.Unmarshal(termV, &wordSearch)
		if err != nil {
			return err
		}

		termM := tx.Bucket([]byte("meta"))
		termMV := termM.Get([]byte("N"))

		counter, err := strconv.ParseInt(string(termMV), 10, 64)
		if err != nil {
			return err
		}

		docsCounter = int(counter)

		return nil
	})

	if err != nil {
		if strings.Contains(err.Error(), "not registered") {
			return nil, nil
		}

		return nil, err
	}

	for docId, counter := range wordSearch.Docs {
		doc, err := GetDoc([]byte(docId), db)
		if err != nil {
			return nil, err
		}

		tf := float64(counter) / float64(doc.Length)
		idf := math.Log((float64(docsCounter)+1)/(float64(wordSearch.DF)+1)) + 1.0
		log.Printf("docId: %v ,tf: %v, idf: %v\n", docId, tf, idf)
		docScore[docId] = (tf * idf)
	}

	return docScore, nil
}
