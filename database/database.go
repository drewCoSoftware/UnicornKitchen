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

func AddIngredient(ingredient *Ingredient) {
	// db := Connect()
	// defer db.Close()
	exec := CreateExecutor(Connect, nil)
	defer exec.Complete()

	addIngredientInternal(exec, ingredient)

}

func addIngredientInternal(exec *dbExecutor, ingredient *Ingredient) {
	i := getIngredientInternal(exec, ingredient.Name)

	if i == nil {
		fmt.Println("Adding the ingredient " + ingredient.Name)
		query := "INSERT INTO ingredients (Name, Description) VALUES ($1,$2)"

		_, err := exec.Exec(query, ingredient.Name, ingredient.Description)
		if err != nil {
			panic(err)
		}
	} else {
		fmt.Println("The ingredient '" + ingredient.Name + "' already exists!")
	}
}

// NOTE: We aren't doing anything about case insensitive search at this time.  Ideally we would
// create a case-insenstive unique index and then perform LOWER checks on name and $1.
// For the purposes of this toy applicaiton, we can live without the additional complexity.
func GetIngredient(name string) *Ingredient {
	exec := CreateExecutor(Connect, nil)
	defer exec.Complete()

	return getIngredientInternal(exec, name)
}

func getIngredientInternal(exec *dbExecutor, name string) *Ingredient {

	query := "SELECT Name, Description FROM ingredients WHERE Name = $1"
	rows, err := exec.Query(query, name)
	defer rows.Close()

	if err != nil {
		panic(err)
	}

	// We should have one or none....
	hasMatch := rows.Next()
	if !hasMatch {
		return nil
	}

	// Deserialize the ingredient.....
	var i Ingredient
	scanErr := rows.Scan(&i.Name, &i.Description)
	if scanErr != nil {
		panic(scanErr)
	}

	return &i
}

// Add a recipe to the database, saving new ingredients, etc. as we go.
func AddRecipe(recipe *Recipe) {

	fmt.Println("Adding the recipe: " + recipe.Name)

	// Set the appropriate refs + validate data (incomplete at this time.)
	validateRecioe(recipe)

	exec := CreateExecutor(Connect, &sql.TxOptions{ReadOnly: false})
	defer exec.Complete()

	// Add All Ingredients first.
	for x := range recipe.Ingredients {
		addIngredientInternal(exec, recipe.Ingredients[x].Ingredient)
	}

	// Add the yield ingredient as well...
	addIngredientInternal(exec, recipe.Yield.Ingredient)

	// If we don't flag the executor as being OK, then the transaction
	// won't get committed...
	exec.SetTransationFlag(true)

}

// NOTE: This validation is very basic....
func validateRecioe(recipe *Recipe) {
	if recipe.Yield == nil {
		panic("There is no yield ingredient listed!")
	}
	for x := range recipe.Ingredients {
		recipe.Ingredients[x].Recipe = recipe
	}
}

func query(db *sql.DB, query string) *sql.Rows {
	res, err := db.Query(query)
	if err != nil {
		panic(err)
	}
	return res
}

// func AddRecipe(*Recipe recipe){

// }

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
