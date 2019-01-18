package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/macros-fetcher/models"
)

// MacrosHandler handlers getting the macros for a given list of ingredients
func MacrosHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		http.Error(w, "invalid request method", http.StatusBadRequest)
		return
	}

	decoder := json.NewDecoder(r.Body)
	ingredients := &models.Ingredients{}
	if err := decoder.Decode(ingredients); err != nil {
		http.Error(w, "invalid json", http.StatusBadRequest)
		return
	}

}
