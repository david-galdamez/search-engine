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

	errors := make(chan error, len(request))

	for _, val := range request {
		wg.Add(1)
		go func(val models.IndexRequest) {
			defer wg.Done()
			if val.Text != "" {
				if err := InsertOrReplaceText(db, val); err != nil {
					errors <- err
				}
			}

			if val.Url != "" {
				if err := InsertOrReplacePage(db, val); err != nil {
					errors <- err
				}
			}
		}(val)
	}

	go func() {
		wg.Wait()
		close(errors)
	}()

	for err := range errors {
		if err != nil {
			utils.RespondWithError(w, http.StatusInternalServerError, err.Error())
			return
		}
	}

	utils.RespondWithJson(w, http.StatusOK, map[string]string{"status": "ok"})
}

func InsertOrReplaceText(db *bolt.DB, val models.IndexRequest) error {
	document, err := services.GetDoc([]byte(val.Id), db)
	if err != nil && !(document == nil && strings.Contains(err.Error(), "not found")) {
		return err
	}

	newDoc := &models.Document{
		Id:      val.Id,
		Title:   val.Title,
		Length:  len(val.Text),
		Content: val.Text,
	}

	if document != nil {
		err := services.DeleteDocTerms(db, document)
		if err != nil {
			return err
		}
	}

	err = services.AddDoc(db, newDoc)
	if err != nil {
		return err
	}

	err = services.AddTextToDB(newDoc.Id, newDoc.Title, newDoc.Content, db)
	if err != nil {
		return err
	}

	if document == nil {
		err = services.IncrementDocCounter(db)
		if err != nil {
			return err
		}
	}

	return nil
}

func InsertOrReplacePage(db *bolt.DB, val models.IndexRequest) error {

	document, err := services.GetDoc([]byte(val.Id), db)
	if err != nil && !(document == nil && strings.Contains(err.Error(), "not found")) {
		return err
	}

	scrapedText, err := services.ScrapePage(val.Url)
	if err != nil {
		return err
	}

	newDoc := &models.Document{
		Id:      val.Id,
		Title:   val.Title,
		Content: scrapedText,
		Length:  len(scrapedText),
		Url:     &val.Url,
	}

	if document != nil {
		err := services.DeleteDocTerms(db, document)
		if err != nil {
			return err
		}
	}

	err = services.AddDoc(db, newDoc)
	if err != nil {
		return err
	}

	err = services.AddTextToDB(newDoc.Id, newDoc.Title, newDoc.Content, db)
	if err != nil {
		return err
	}

	if document == nil {
		err = services.IncrementDocCounter(db)
		if err != nil {
			return err
		}
	}

	return nil
}
