package main

import (
	"encoding/json"
	"fmt"
	"os"

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

	// Let's pull a recipe from file and add it that way...
	file, err := os.Open("data/default-data.json")
	if err != nil {
		panic(err)
	}
	//	file.Stat()
	defer file.Close()

	dec := json.NewDecoder(file)
	var r2 database.Recipe
	if decodeErr := dec.Decode(&r2); decodeErr != nil {
		panic(decodeErr)
	}

	database.AddRecipe(&r2)

	// Create a json file that will contain all of the 'default' data that we want to
	// appear.  Use this to create + populate the data, like all of the code above is
	// doing.

	// After that we can look at setting up an HTTP endpoint to do some GraphQL queries
	// against our data....

}
