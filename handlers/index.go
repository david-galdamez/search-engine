package handlers

import (
	"encoding/json"
	"net/http"
	"os"
	"strings"

	"github.com/david-galdamez/search-engine/utils"
)

type IndexedDocs map[string]map[string]int

func (idx IndexedDocs) Add(docId, text string) {
	wordsIterator := strings.FieldsSeq(strings.ToLower(text))
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

var index = make(IndexedDocs)

func Index(w http.ResponseWriter, r *http.Request) {

	request := indexRequest{}

	err := json.NewDecoder(r.Body).Decode(&request)
	if err != nil {
		utils.RespondWithError(w, 500, "Error decoding json")
		return
	}
	defer r.Body.Close()

	if request.Text != "" {
		index.Add(request.Id, request.Text)
	}

	err = index.SaveIntoJson()
	if err != nil {
		utils.RespondWithError(w, 500, "Error saving into json")
		return
	}
}
