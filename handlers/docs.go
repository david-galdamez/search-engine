package handlers

import (
	"context"
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/david-galdamez/search-engine/database"
	"github.com/david-galdamez/search-engine/models"
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
	document := &models.Document{}

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
		log.Println("Response timeout")
		utils.RespondWithError(w, http.StatusRequestTimeout, "Response timeout")
		return
	case err := <-done:
		if err != nil {
			if strings.Contains(err.Error(), "not found") {
				utils.RespondWithError(w, http.StatusNotFound, err.Error())
				return
			}
			utils.RespondWithError(w, http.StatusInternalServerError, err.Error())
			return
		}
	}

	utils.RespondWithJson(w, http.StatusOK, *document)
}
