package database

import (
	"database/sql"
	"fmt"
	"strings"

	"github.com/drewCoSoftware/UnicornKitchen/settings"
	_ "github.com/lib/pq"
)

// This will create the database and setup the schema.
func CreateDatabase() {

	// Let's make the Database first...
	exists := dbExists(settings.DB_NAME)

	db := getPostGresDB()
	defer db.Close()

	if !exists {
		fmt.Println("The database doesn't exist, so we will create it!")

		_, err := db.Exec("CREATE DATABASE " + settings.DB_NAME)
		if err != nil {
			fmt.Println("There was an error creating the Database!")
			fmt.Println(err.Error())
		} else {
			fmt.Println("Database creation OK!")
		}
	} else {
		fmt.Println("The database already exists.")
	}

	fmt.Println("Creating the schema....")
	createSchema(true)
}

func Connect() *sql.DB {
	ops := settings.GetDatabaseOptions()

	connectionString := getConnectionString(ops)
	res := OpenConnection(connectionString)
	return res
}

func OpenConnection(connectionString string) *sql.DB {
	res, err := sql.Open("postgres", connectionString)
	if err != nil {
		panic(err)
	}
	return res
}

func getConnectionString(ops *settings.DBOptions) string {
	fmtString := "port=%s host=%s user=%s password=%s dbname=%s sslmode=disable"
	useHost := ops.Address
	usePort := "5432"

	// Parse out the port / host from the address.
	parts := strings.Split(ops.Address, ":")
	pLen := len(parts)
	if pLen > 1 {
		useHost = parts[0]
		usePort = parts[1]
	}
	res := fmt.Sprintf(fmtString, usePort, useHost, ops.User, ops.Password, ops.Database)

	return res
}

func dbExists(dbName string) bool {
	// We want to check the postgres catalog this time around.
	db := getPostGresDB()

	defer db.Close()

	// Check for the database....
	qr, err := db.Query("SELECT datname FROM pg_catalog.pg_database WHERE datname = '" + dbName + "'")
	if err != nil {
		panic(err)
	}
	defer qr.Close()

	hasOne := qr.Next()
	return hasOne
}

func getPostGresDB() *sql.DB {
	ops := settings.GetDatabaseOptions()
	ops.Database = "postgres"

	cs := getConnectionString(ops)
	db := OpenConnection(cs)
	return db
}

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
		dropTable(db, "recipe_ingredients")
		dropTable(db, "ingredients")
		dropTable(db, "recipes")
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
									Description VARCHAR NULL )`

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
}

func exec(db *sql.DB, query string, confirm bool) {
	qr, err := db.Exec(query)
	if err != nil {
		panic(err)
	}
	if confirm {
		checkExec(qr)
	}
}

func checkExec(qr sql.Result) {
	count, err := qr.RowsAffected()
	if err != nil {
		panic(err)
	}
	if count < 1 {
		panic("The query did not succeed!")
	}
}
