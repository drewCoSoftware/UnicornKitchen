package database

type Ingredient struct {
	Id   int64
	Name string `pg:",unique,notnull"`
}

type Recipe struct {
	Id         int64
	Name       string               `pg:",unique,notnull"`
	Ingredient []RecipeToIngredient `pg:"many2many:recipe_to_ingredients"`
}

type RecipeToIngredient struct {
	Id               int64
	RecipeId         int64
	IngredientId     int64
	IngredientAmount string
}
