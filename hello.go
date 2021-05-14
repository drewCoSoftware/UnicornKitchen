package main

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/drewCoSoftware/UnicornKitchen/database"
	//	"github.com/drewCoSoftware/UnicornKitchen/ingredients"
)

func main() {
	fmt.Println("What's cookin' in the Unicorn Kitchen?")

	// potato := &database.Ingredient{
	// 	Name:        "Potato",
	// 	Description: "A starchy tuber!",
	// }
	// database.AddIngredient(potato)

	type defaultData struct {
		Ingredients []database.Ingredient
		Recipes     []database.Recipe
	}

	// Let's pull a recipe from file and add it that way...
	file, err := os.Open("data/default-data.json")
	if err != nil {
		panic(err)
	}
	defer file.Close()

	dec := json.NewDecoder(file)
	var dd defaultData
	if decodeErr := dec.Decode(&dd); decodeErr != nil {
		panic(decodeErr)
	}

	// Splat the default data into the database.
	database.CreateDatabase()
	for _, i := range dd.Ingredients {
		database.AddIngredient(&i)
	}

	for _, r := range dd.Recipes {
		database.AddRecipe(&r)
	}

	// Our recipes need instructions....

	// After that we can look at setting up an HTTP endpoint to do some GraphQL queries
	// against our data....

}
