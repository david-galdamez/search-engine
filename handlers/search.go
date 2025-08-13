package handlers

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

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

	index := make(IndexedDocs)

	done := make(chan error, 1)

	go func() {
		data, err := os.ReadFile("sample.json")
		if err != nil {
			done <- fmt.Errorf("Error reading json: %v", err)
			return
		}

		err = json.Unmarshal(data, &index)
		if err != nil {
			done <- fmt.Errorf("Error decoding json: %v", err)
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

	utils.RespondWithJson(w, http.StatusOK, index[search])
}
