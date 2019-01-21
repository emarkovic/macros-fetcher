package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/macros-fetcher/models"
)

func getIngredientDbId(ingredientName string, API_KEY string) models.IngredientDbId {
	// make request to get the db id

	// ingredient was not found
	return -1
}

// MacrosHandler handlers getting the macros for a given list of ingredients
func (ctx *Context) MacrosHandler(w http.ResponseWriter, r *http.Request) {
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

	for key := range *ingredients {
		dbId := getIngredientDbId(key, ctx.UsdaAPIKey)
		fmt.Println(dbId)
	}

}
