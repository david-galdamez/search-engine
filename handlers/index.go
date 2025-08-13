package handlers

import (
	"encoding/json"
	"net/http"
	"os"
	"strings"
	"sync"
	"unicode"

	"github.com/david-galdamez/search-engine/utils"
)

type IndexedDocs map[string]map[string]int

func (idx IndexedDocs) AddText(docId, text string) {
	//trims punctuations and split by spaces
	cleanText := strings.Map(func(r rune) rune {
		if unicode.IsPunct(r) {
			return -1
		}

		return r
	}, text)

	wordsIterator := strings.FieldsSeq(strings.ToLower(cleanText))
	for word := range wordsIterator {
		if len(word) <= 2 {
			continue
		}

		if idx[word] == nil {
			idx[word] = make(map[string]int)
		}
		idx[word][docId]++
	}
}

func (idx IndexedDocs) SaveIntoJson() error {
	file, err := os.Create("sample.json")
	if err != nil {
		return err
	}
	defer file.Close()

	data, err := json.Marshal(idx)
	if err != nil {
		return err
	}

	_, err = file.Write(data)
	if err != nil {
		return err
	}

	return nil
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

	index := make(IndexedDocs)

	if data, err := os.ReadFile("sample.json"); err == nil {
		_ = json.Unmarshal(data, &index)
	}

	for _, value := range request {
		wg.Add(1)
		go func(val indexRequest) {
			defer wg.Done()
			if val.Text != "" {
				mux.Lock()
				index.AddText(val.Id, val.Text)
				mux.Unlock()
			}
		}(value)
	}

	wg.Wait()

	err = index.SaveIntoJson()
	if err != nil {
		utils.RespondWithError(w, 500, "Error saving into json")
		return
	}

	utils.RespondWithJson(w, http.StatusOK, map[string]string{"status": "ok"})
}
