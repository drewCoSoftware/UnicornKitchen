package database

// NOTE: If we want to add ORM type features to this, we would annotate the individual members.
// See previous versions of this code that used go-pg to do so.
type Ingredient struct {
	Id          int64
	Name        string
	Description string
}

type Recipe struct {
	Id           int64  `json:"id"`
	Name         string `json:"name"`
	Description  string `json:"description"`
	Ingredients  []*RecipeIngredient
	Instructions []*RecipeInstruction
	Yield        *RecipeYield
}

// A single instruction that is included in a recipe.
type RecipeInstruction struct {
	Id       int64
	RecipeId int64
	Order    int64
	Content  string
}

type RecipeIngredient struct {
	Id               int64
	Recipe           *Recipe
	Ingredient       *Ingredient
	IngredientAmount string // 1 cup, 2 dozen, 3 gallons, etc.
}

type RecipeYield struct {
	Amount     string // GO doesn't have a convenient decimal type, and we still need to represetn things like cups, pounds, mL, etc.
	Ingredient *Ingredient
}
