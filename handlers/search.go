package handlers

import (
	"context"
	"log"
	"net/http"
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

	done := make(chan error, 1)
	searchedData := &models.Terms{}

	go func() {
		searchedData, err = services.SearchWordInDB([]byte(search), db)
		if err != nil {
			done <- err
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
