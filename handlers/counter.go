package handlers

import (
	"context"
	"log"
	"net/http"
	"time"

	"github.com/david-galdamez/search-engine/database"
	"github.com/david-galdamez/search-engine/services"
	"github.com/david-galdamez/search-engine/utils"
)

type CounterResponse struct {
	Counter int `json:"docs_counter"`
}

func Counter(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(r.Context(), 10*time.Second)
	defer cancel()

	db, err := database.GetDB()
	if err != nil {
		log.Fatalf("Error opening database")
	}
	defer db.Close()

	done := make(chan error, 1)
	var counter *int

	go func() {
		counter, err = services.GetDocumentCounter(db)
		if err != nil {
			done <- err
			return
		}

		done <- nil
	}()

	select {
	case <-ctx.Done():
		log.Println("Request timeout")
		utils.RespondWithError(w, http.StatusRequestTimeout, "request timeout")
		return
	case err := <-done:
		if err != nil {
			utils.RespondWithError(w, http.StatusInternalServerError, err.Error())
			return
		}
	}

	utils.RespondWithJson(w, http.StatusOK, CounterResponse{Counter: *counter})
}
