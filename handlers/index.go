package handlers

import (
	"encoding/json"
	"log"
	"net/http"
	"strings"
	"sync"

	"github.com/boltdb/bolt"
	"github.com/david-galdamez/search-engine/database"
	"github.com/david-galdamez/search-engine/models"
	"github.com/david-galdamez/search-engine/services"
	"github.com/david-galdamez/search-engine/utils"
)

func Index(w http.ResponseWriter, r *http.Request) {

	var mux sync.RWMutex
	var wg sync.WaitGroup
	request := []models.IndexRequest{}

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
		go func(val models.IndexRequest) {
			defer wg.Done()
			if val.Text != "" {
				err := InsertOrReplaceText(db, val, &mux)
				if err != nil {
					utils.RespondWithError(w, http.StatusInternalServerError, err.Error())
					return
				}
			}
		}(value)
	}
	wg.Wait()

	utils.RespondWithJson(w, http.StatusOK, map[string]string{"status": "ok"})
}

func InsertOrReplaceText(db *bolt.DB, val models.IndexRequest, mux *sync.RWMutex) error {
	document, err := services.GetDoc([]byte(val.Id), db)
	if err != nil {
		if document == nil && strings.Contains(err.Error(), "not found") {
			mux.Lock()

			newDoc := &models.Document{
				Id:      val.Id,
				Title:   val.Title,
				Length:  len(val.Text),
				Content: val.Text,
			}

			err := services.AddDoc(db, newDoc)
			if err != nil {
				return err
			}

			err = services.AddTextToDB(val.Id, val.Title, val.Text, db)
			if err != nil {
				return err
			}

			err = services.IncrementDocCounter(db)
			if err != nil {
				return err
			}
			mux.Unlock()

			return nil
		}
		return err
	}

	if document != nil {
		mux.Lock()
		err := services.DeleteDocTerms(db, document)
		if err != nil {
			return err
		}

		newDoc := &models.Document{
			Id:      val.Id,
			Title:   val.Title,
			Length:  len(val.Text),
			Content: val.Text,
		}
		err = services.AddDoc(db, newDoc)
		if err != nil {
			return err
		}

		err = services.AddTextToDB(val.Id, val.Text, val.Title, db)
		if err != nil {
			return err
		}
		mux.Unlock()
	}

	return nil
}
