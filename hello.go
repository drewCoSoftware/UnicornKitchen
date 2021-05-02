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

	database.CreateDatabase()

	//	db := database.
	// fmt.Println("Checking for the database....")
	// hasDB, err := database.DbExists(settings.DB_NAME)
	// if err != nil {
	// 	panic(err)
	// }

	// // Gee, a ternary sure would be nice here......
	// tf := "false"
	// if hasDB {
	// 	tf = "true"
	// }

	// fmt.Println(fmt.Sprintf("The database %s exists? (%s)", settings.DB_NAME, tf))
	// As a sanity check, I am going to print the db settings, which should use
	// the environment variables from the server....
	// s := settings.GetDatabaseOptions()

	// fmt.Println("Address = " + s.Addr)
	// fmt.Println("db name = " + s.Database)
	// fmt.Println("db user = " + s.User)

	// dbCfg := os.Getenv("DB_CONFIG_PATH")
	// if dbCfg == "" {
	// 	dbCfg = "local-db-cfg.json"
	// }

	// fmt.Println("The database cfg is: '" + dbCfg + "'")
	// hasFish := ingredients.HasIngredient("fish")
	// fmt.Println(fmt.Sprintf("has fish?: %t", hasFish))

	// We want to create our database resources if they don't currently exist.
	// database.CreateDatabase()
	// database.CreateTables(true)
	// database.AddDefaultData()

	//	Let's create a new recipe, and add the data, all in one go.
	r1 := &database.Recipe{
		Name: "Electric Potato",
	}
	r1.Ingredients = []database.RecipeIngredient{
		{
			Recipe:           r1,
			Ingredient:       &database.Ingredient{Name: "Potato"},
			IngredientAmount: "1",
		},
	}

	// db := database.GetConnection()
	// database.SaveRecipe(db, r1)

}
