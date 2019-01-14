package handlers

import (
	"encoding/json"
	"io"
	"log"
	"net/http"
	"strings"

	"golang.org/x/net/html"
)

// IngredientDetails represent the amount and unit related to an ingredient in the recipe
type IngredientDetails struct {
	Amount string
	Unit   string
}

// Ingredients represent to names of ingredients used in the recipe as well as the amount and unit
type Ingredients map[string]IngredientDetails

func processWPRMIngredient(tokenizer *html.Tokenizer) (string, string, string) {
	/*
		contains 3 parts:
		<li class="wprm-recipe-ingredient">
			<span class="wprm-recipe-ingredient-amount">2</span>
			<span class="wprm-recipe-ingredient-unit">tablespoons</span>
			<span class="wprm-recipe-ingredient-name">butter</span>
		</li>
	*/
	// get out of the loop when its a closing li tag
	var ingredient, amount, unit string
	for {
		tokenType := tokenizer.Next()

		switch tokenType {
		case html.StartTagToken:
			token := tokenizer.Token()
			if token.Data == "span" {
				for _, attr := range token.Attr {
					if attr.Key == "class" {
						switch attr.Val {
						case "wprm-recipe-ingredient-amount":
							if tokenizer.Next() == html.TextToken {
								amount = tokenizer.Token().Data
							}
						case "wprm-recipe-ingredient-unit":
							if tokenizer.Next() == html.TextToken {
								unit = tokenizer.Token().Data
							}
						case "wprm-recipe-ingredient-name":
							if tokenizer.Next() == html.TextToken {
								ingredient = tokenizer.Token().Data
							}
						}
					}
				}
			}
			break
		case html.EndTagToken:
			token := tokenizer.Token()
			if token.Data == "li" {
				// return the ingredients when we see a closing li tag.
				return ingredient, amount, unit
			}
		}
	}
}

func processIngredient(tokenizer *html.Tokenizer) {
	/*
		contains:
		<li class="ingredient" itemprop="ingredients">6 cups chicken stock</li>

		ideas to process:
		split on strings
	*/

}

func processHTML(w http.ResponseWriter, resp *http.Response) (*Ingredients, error) {
	// create a new tokenizer over the response body
	tokenizer := html.NewTokenizer(resp.Body)
	ingredients := make(Ingredients)
	for {
		tokenType := tokenizer.Next()
		switch tokenType {
		case html.ErrorToken:
			err := tokenizer.Err()
			if err == io.EOF {
				// end of the file, break out of loop
				return &ingredients, nil
			}

			// if not the end of file, then there was an actual error tokenizing which
			// likely means the HTML was malfortatted.
			http.Error(w, "error tokenizing html", http.StatusInternalServerError)
			return nil, err
		case html.StartTagToken:
			token := tokenizer.Token()
			if token.Data == "li" {
				for _, attr := range token.Attr {
					if attr.Key == "class" {
						if strings.Contains(attr.Val, "wprm-recipe-ingredient") {
							ingredient, amount, unit := processWPRMIngredient(tokenizer)
							ingredients[ingredient] = IngredientDetails{
								Amount: amount,
								Unit:   unit,
							}
						} else if strings.Contains(attr.Val, "ingredient") {
							processIngredient(tokenizer)
						}
					}
				}
			}
		}

	}
}

func getIngredients(w http.ResponseWriter, url string) (*Ingredients, error) {
	// look if its in the wprm format..
	// or "ingredients"

	resp, err := http.Get(url)
	// report an error if there was one
	if err != nil {
		http.Error(w, "error fetching URL", http.StatusInternalServerError)
		return nil, err
	}

	// make sure the response body get closed at the end of the function
	defer resp.Body.Close()

	// check response status code
	if resp.StatusCode != http.StatusOK {
		http.Error(w, "response status code was not OK", http.StatusInternalServerError)
		return nil, err
	}

	// check response content type -> must be html
	ctype := resp.Header.Get("Content-Type")
	if !strings.HasPrefix(ctype, "text/html") {
		http.Error(w, "response content type was not test/htm", http.StatusInternalServerError)
	}

	return processHTML(w, resp)
}

// IngredientsHandler handles getting the list of ingredients from a recipe url
func IngredientsHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Add("Access-Control-Allow-Origin", "*")

	url := r.FormValue("url")

	if len(url) == 0 {
		http.Error(w, "no url provided", http.StatusBadRequest)
		return
	}

	ingredients, err := getIngredients(w, url)
	if err != nil {
		log.Fatalf("error getting ingredients: %v\n", err)
	}

	w.Header().Add("Content-Type", "application/json; charset=utf-8")
	encoder := json.NewEncoder(w)
	if err := encoder.Encode(ingredients); err != nil {
		http.Error(w, "error encoding json: "+err.Error(), http.StatusInternalServerError)
		return
	}
}
