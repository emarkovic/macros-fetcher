package models

// IngredientDetails represent the amount and unit related to an ingredient in the recipe
type IngredientDetails struct {
	Amount string
	Unit   string
}

// Ingredients represent to names of ingredients used in the recipe as well as the amount and unit
type Ingredients map[string]IngredientDetails

// IngredientDbId represents the USDA database lookup ID for a given ingredient
type IngredientDbId int

// IngredientToDbIdMap represents names of ingredients to their ingredient DB ID
type IngredientToDbIdMap map[string]IngredientDbId
