package database

import (
	"database/sql"
	"fmt"
)

func dropTable(db *sql.DB, name string) {
	query := "DROP TABLE IF EXISTS " + name
	_, err := db.Exec(query)
	if err != nil {
		panic(err)
	}
}

// // Create the UnicorKitchen schema if it doesn't currently exist.
func createSchema(removeExistingTables bool) {

	db := Connect()
	defer db.Close()

	if removeExistingTables {
		fmt.Println("The old tables will be removed.")
		dropTable(db, "recipe_instructions")
		dropTable(db, "recipe_ingredients")
		dropTable(db, "recipes")
		dropTable(db, "ingredients")
	}

	// NOTE: Reflection would be cool, but that makes it more difficult to make
	// our otherwise simple tables.
	fmt.Println("Creating table for 'ingredients'.")
	query := `CREATE TABLE ingredients ( IngredientId BIGSERIAL PRIMARY KEY,
										 Name VARCHAR NOT NULL UNIQUE, 
										 Description VARCHAR NULL )`

	exec(db, query, false)

	fmt.Println("Creating table for 'recipes'")
	query = `CREATE TABLE recipes ( RecipeId BIGSERIAL PRIMARY KEY,
									Name VARCHAR NOT NULL UNIQUE,
									Description VARCHAR NULL,
									YieldAmount VARCHAR NOT NULL,
									YieldIngredientId BIGINT NOT NULL REFERENCES ingredients on DELETE CASCADE )`

	exec(db, query, false)

	fmt.Println("Creating table 'recipe_ingredients'")
	query = `CREATE TABLE recipe_ingredients ( 
				Id BIGSERIAL PRIMARY KEY,
				RecipeId BIGINT NOT NULL REFERENCES recipes ON DELETE CASCADE,
				IngredientId BIGINT NOT NULL REFERENCES ingredients ON DELETE CASCADE,
				Amount VARCHAR NOT NULL
					--, FOREIGN KEY(RecipeId) REFERENCES recipes (RecipeId)
					--, FOREIGN KEY(IngredientId) REFERENCES ingredients (IngredientId)
				) -- NOTE: It appears that PostGres will auto-create the FK refs as a result of the DELETE CASCADEs mentioned above.
				  -- I have left the explicit FK syntax in place so that one might contemplate the differences.`

	exec(db, query, false)

	fmt.Println("Creating table 'recipe_instructions'")
	query = `CREATE TABLE recipe_instructions (
			 Id BIGSERIAL PRIMARY KEY,
			 RecipeId BIGINT NOT NULL REFERENCES recipes (RecipeId) ON DELETE CASCADE,
			 InstructionOrder BIGINT NOT NULL,
			 Content VARCHAR NOT NULL )`

	exec(db, query, false)

}
