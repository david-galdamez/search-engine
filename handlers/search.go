package handlers

import (
	"context"
	"log"
	"net/http"
	"strings"
	"sync"
	"time"

	"github.com/david-galdamez/search-engine/database"
	"github.com/david-galdamez/search-engine/models"
	"github.com/david-galdamez/search-engine/services"
	"github.com/david-galdamez/search-engine/utils"
)

func Search(w http.ResponseWriter, r *http.Request) {

	ctx, cancel := context.WithTimeout(r.Context(), 10*time.Second)
	defer cancel()

	search := r.URL.Query().Get("q")
	if search == "" {
		utils.RespondWithError(w, http.StatusBadRequest, "Query parameter missing")
		return
	}
	defer r.Body.Close()

	db, err := database.GetDB()
	if err != nil {
		log.Fatalf("Error opening database: %v\n", err)
	}
	defer db.Close()

	var wg sync.WaitGroup
	var mux sync.Mutex
	searchedData := make(models.DocsScore)
	words := utils.Tokenizer(search)
	done := make(chan error, len(words))

	for _, word := range words {
		wg.Add(1)
		go func(word string) {
			defer wg.Done()
			tdfIdf, err := services.SearchCalculate(db, []byte(word))
			if err != nil {
				done <- err
				return
			}
			mux.Lock()
			for docId, score := range tdfIdf {
				searchedData[docId] += score
			}
			mux.Unlock()
		}(word)
	}

	go func() {
		wg.Wait()
		close(done)
	}()

	for {
		select {
		case <-ctx.Done():
			log.Println("Response timeout")
			utils.RespondWithError(w, http.StatusRequestTimeout, "Response timeout")
			return
		case err, ok := <-done:
			if !ok {
				goto BUILD_RESPONSE
			}
			if err != nil {
				if strings.Contains(err.Error(), "not found") {
					continue
				}
				utils.RespondWithError(w, http.StatusInternalServerError, err.Error())
				return
			}
		}
	}

BUILD_RESPONSE:
	searchedResponse := models.SearchResponse{Query: search, TotalResults: len(searchedData)}
	results := []models.SearchResults{}
	for docId, value := range searchedData {
		doc, err := services.GetDoc([]byte(docId), db)
		if err != nil {
			if strings.Contains(err.Error(), "not found") {
				continue
			}
			utils.RespondWithError(w, http.StatusInternalServerError, err.Error())
			return
		}

		newResult := models.SearchResults{DocId: doc.Id, Title: doc.Title, Score: value}
		utils.PushAndSort(&results, newResult)
	}

	searchedResponse.Results = results

	utils.RespondWithJson(w, http.StatusOK, searchedResponse)
}
