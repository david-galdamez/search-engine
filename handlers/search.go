package handlers

import (
	"encoding/json"
	"net/http"
	"os"

	"github.com/david-galdamez/search-engine/utils"
)

func Search(w http.ResponseWriter, r *http.Request) {
	search := r.URL.Query().Get("q")
	if search == "" {
		utils.RespondWithError(w, http.StatusBadRequest, "Query parameter missing")
		return
	}
	defer r.Body.Close()

	index := make(IndexedDocs)

	data, err := os.ReadFile("sample.json")
	if err != nil {
		utils.RespondWithError(w, http.StatusInternalServerError, "Error reading json")
		return
	}

	err = json.NewDecoder(os.Stdin).Decode(&data)
	if err != nil {
		utils.RespondWithError(w, http.StatusInternalServerError, "Error converting json")
		return
	}

	utils.RespondWithJson(w, http.StatusOK, index[search])
}
