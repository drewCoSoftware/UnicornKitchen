package main

import (
	"fmt"
	"github.com/drewCoSoftware/UnicornKitchen/database"
)

func main() {
	fmt.Println("What's cookin' in the Unicorn Kitchen?")

	database.CreateDatabase()

	potato := &database.Ingredient{
		Name:        "Potato",
		Description: "A starchy tuber!",
	}
	database.AddIngredient(potato)

	//	Let's create a new recipe, and add the data, all in one go.
	r1 := &database.Recipe{
		Name:        "Electric Potato",
		Description: "Create some volts from a potato, lemon and other common ingredients.",
		Yield: &database.RecipeYield{
			Amount: "0.5",
			Ingredient: &database.Ingredient{
				Name:        "Volt",
				Description: "A unit of potential electrical energy.",
			},
		},
	}
	r1.Ingredients = []*database.RecipeIngredient{
		{
			Ingredient:       potato,
			IngredientAmount: "1",
		},
		{
			Ingredient:       &database.Ingredient{Name: "Penny"},
			IngredientAmount: "1",
		},
		{
			Ingredient:       &database.Ingredient{Name: "Galvanized Iron Nail"},
			IngredientAmount: "1",
		},
		{
			Ingredient:       &database.Ingredient{Name: "18ga Copper Wire"},
			IngredientAmount: "6 Inches",
		},
		{
			Ingredient:       &database.Ingredient{Name: "Alligator Clip"},
			IngredientAmount: "2",
		},
	}

	database.AddRecipe(r1)

	// Complete the recipe adding code...
	// Find a way to define the recipes in a JSON file to make adding the default data less verbose.

}
