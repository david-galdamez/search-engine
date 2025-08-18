package handlers

import (
	"encoding/json"
	"log"
	"net/http"
	"strings"
	"sync"
	"unicode"

	"github.com/david-galdamez/search-engine/database"
	"github.com/david-galdamez/search-engine/utils"
)

func AddTextToDB(docId, text string) {

	db, err := database.GetDB()
	if err != nil {
		log.Fatalf("Error opening database: %v\n", err)
	}
	defer db.Close()

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

type indexRequest struct {
	Id   string `json:"id"`
	Text string `json:"text,omitempty"`
	Url  string `json:"url,omitempty"`
}

func Index(w http.ResponseWriter, r *http.Request) {

	var mux sync.RWMutex
	var wg sync.WaitGroup
	request := []indexRequest{}

	err := json.NewDecoder(r.Body).Decode(&request)
	if err != nil {
		utils.RespondWithError(w, 500, "Error decoding json")
		return
	}
	defer r.Body.Close()

	for _, value := range request {
		wg.Add(1)
		go func(val indexRequest) {
			defer wg.Done()
			if val.Text != "" {
				mux.Lock()
				AddTextToDB(val.Id, val.Text)
				mux.Unlock()
			}
		}(value)
	}

	wg.Wait()

	utils.RespondWithJson(w, http.StatusOK, map[string]string{"status": "ok"})
}
