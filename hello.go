package main

import (
	"fmt"

	//"github.com/drewCoSoftware/UnicornKitchen/settings"
	// "os"
	"github.com/drewCoSoftware/UnicornKitchen/database"
	//	"github.com/drewCoSoftware/UnicornKitchen/settings"
	// "github.com/drewCoSoftware/UnicornKitchen/ingredients"
	//	"github.com/go-pg/pg/v10"
	//	"github.com/go-pg/pg/v10/orm"
)

//import "github.com/google/go-cmp/cmp"

func main() {
	fmt.Println("What's cookin' in the Unicorn Kitchen?")

	// database.CreateDatabase()
	potato := &database.Ingredient{
		Name:        "potato",
		Description: "A starchy tuber!",
	}
	database.AddIngredient(potato)

	//	Let's create a new recipe, and add the data, all in one go.
	r1 := &database.Recipe{
		Name:        "Electric Potato",
		Description: "Create some volts from a potato, lemon and other common ingredients.",
	}
	r1.Ingredients = []*database.RecipeIngredient{
		{
			Ingredient:       &database.Ingredient{Name: "Potato"},
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

}
