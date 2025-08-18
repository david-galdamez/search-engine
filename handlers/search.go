package handlers

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/david-galdamez/search-engine/database"
	"github.com/david-galdamez/search-engine/utils"
)

type SearchedData map[string]int

func SearchWordInDB(word []byte) SearchedData {

	db, err := database.GetDB()
	if err != nil {
		log.Fatalf("Error opening database: %v\n", err)
	}
	defer db.Close()

	tx, err := db.Begin(true)
	if err != nil {
		log.Fatalf("Error starting transaction: %v\n", err)
	}
	defer tx.Rollback()

	termB := tx.Bucket([]byte("terms"))
	termV := termB.Get(word)
	if err != nil {
		return nil
	}

	searchedData := make(SearchedData)

	err = json.Unmarshal(termV, &searchedData)
	if err != nil {
		log.Fatalf("Error parsing json: %v\n", err)
	}

	return searchedData
}

func Search(w http.ResponseWriter, r *http.Request) {

	ctx, cancel := context.WithTimeout(r.Context(), 10*time.Second)
	defer cancel()

	search := r.URL.Query().Get("q")
	if search == "" {
		utils.RespondWithError(w, http.StatusBadRequest, "Query parameter missing")
		return
	}
	defer r.Body.Close()

	done := make(chan error, 1)
	searchedData := make(SearchedData)

	go func() {
		searchedData = SearchWordInDB([]byte(search))
		if searchedData == nil {
			done <- fmt.Errorf("Documents for word not found")
			return
		}

		done <- nil
	}()

	select {
	case <-ctx.Done():
		log.Print("Response timeout")
		utils.RespondWithError(w, http.StatusRequestTimeout, "Response timeout")
		return
	case err := <-done:
		if err != nil {
			utils.RespondWithError(w, http.StatusNotFound, err.Error())
			return
		}
	}

	utils.RespondWithJson(w, http.StatusOK, searchedData)
}
