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

func Docs(w http.ResponseWriter, r *http.Request) {

	ctx, cancel := context.WithTimeout(r.Context(), 10*time.Second)
	defer cancel()

	db, err := database.GetDB()
	if err != nil {
		log.Fatalf("Error opening database: %v\n", err)
	}
	defer db.Close()

	docId := r.PathValue("id")

	done := make(chan error, 1)
	document := &services.Document{}

	go func() {
		document, err = services.GetDoc([]byte(docId), db)
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
			utils.RespondWithError(w, http.StatusInternalServerError, err.Error())
			return
		}
	}

	utils.RespondWithJson(w, http.StatusOK, *document)
}
