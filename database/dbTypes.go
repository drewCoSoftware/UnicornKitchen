package database

type Ingredient struct {
	Id   int64
	Name string `pg:",unique,notnull"`
}

type Recipe struct {
	Id          int64
	Name        string             `pg:",unique,notnull"`
	Ingredients []RecipeIngredient `pg:"rel:has-many"`
}

type RecipeIngredient struct {
	Id               int64
	Recipe           *Recipe     `pg:"rel:has-one"`
	Ingredient       *Ingredient `pg:"rel:has-one"`
	IngredientAmount string
}
