package handlers

import (
	"encoding/json"
	"log"
	"net/http"
	"sync"

	"github.com/david-galdamez/search-engine/database"
	"github.com/david-galdamez/search-engine/services"
	"github.com/david-galdamez/search-engine/utils"
)

type indexRequest struct {
	Id    string `json:"id"`
	Title string `json:"title"`
	Text  string `json:"text,omitempty"`
	Url   string `json:"url,omitempty"`
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

	db, err := database.GetDB()
	if err != nil {
		log.Fatalf("Error opening database: %v\n", err)
	}
	defer db.Close()

	for _, value := range request {
		wg.Add(1)
		go func(val indexRequest) {
			defer wg.Done()
			if val.Text != "" {
				mux.Lock()
				services.AddTextToDB(val.Id, val.Title, val.Text, db)
				doc := &services.Document{
					Id:     val.Id,
					Title:  val.Title,
					Length: len(val.Text),
					Text:   val.Text,
				}
				services.AddDoc(db, doc)
				services.IncrementDocCounter(db)
				mux.Unlock()
			}
		}(value)
	}

	wg.Wait()

	utils.RespondWithJson(w, http.StatusOK, map[string]string{"status": "ok"})
}
