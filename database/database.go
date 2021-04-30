package database

import (
	"database/sql"
	"fmt"
	"strings"

	"github.com/drewCoSoftware/UnicornKitchen/settings"
	_ "github.com/lib/pq"
)

// func CreateDatabase() {
// 	db := Connect(settings.GetDatabaseOptions()) //pg.Connect(settings.GetDatabaseOptions())
// 	defer db.Close()

// 	// Let's make the Database first...
// 	exists, err := dbExists(db, settings.DB_NAME)
// 	if err != nil {
// 		fmt.Println("There was an error checking for the database!")
// 		fmt.Println(err.Error())
// 		return
// 	}

// 	if !exists {
// 		fmt.Println("The database doesn't exist, so we will create it!")

// 		_, err := db.Exec("CREATE DATABASE " + settings.DB_NAME)
// 		if err != nil {
// 			fmt.Println("There was an error creating the Database!")
// 			fmt.Println(err.Error())
// 		} else {
// 			fmt.Println("Database creation OK!")
// 		}
// 	} else {
// 		fmt.Println("The database already exists.")
// 	}

// }
func CreateDatabase() {

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

func DbExists(dbName string) (bool, error) {
	// We want to check the postgres catalog this time around.
	ops := settings.GetDatabaseOptions()
	ops.Database = "postgres"

	cs := getConnectionString(ops)
	db := OpenConnection(cs)
	defer db.Close()

	// Check for the database....
	qr, err := db.Query("SELECT datname FROM pg_catalog.pg_database WHERE datname = '" + dbName + "'")
	if err != nil {
		return false, err
		// fmt.Println("There was an error checking for the database!")
		// fmt.Println(err.Error())
	}
	defer qr.Close()

	hasOne := qr.Next()
	return hasOne, nil
}

// func AddDefaultData() {
// 	fmt.Println("Adding some default data....")

// 	db := GetConnection()
// 	defer db.Close()

// 	// NOTE: A string list or input file would be a good idea...
// 	// Maybe just a json file?
// 	i1 := &Ingredient{
// 		Name: "Potato",
// 	}
// 	insertData(db, i1)

// }

// func (db *pg.DB)  SaveRecipe(r *Recipe) {
// }

// func SaveRecipe(db *pg.DB, r *Recipe) {

// 	fmt.Println("Saving Recipe: " + r.Name)

// 	// Iterate over the ingredients, saving each if it doesn't exist...
// 	for item := range r.Ingredients {
// 		ingredient := r.Ingredients[item].Ingredient
// 		//		fmt.Println(ingredient.Name)

// 		var i Ingredient
// 		var match = db.Model(&i).First()

// 		fmt.Println(ingredient.Name)
// 		fmt.Println(match)
// 	}
// }

// func insertData(db *pg.DB, data interface{}) {
// 	_, err := db.Model(data).Insert()
// 	if err != nil {
// 		panic(err)
// 	}
// }

// func CreateTables(removeExisting bool) {
// 	db := GetConnection()
// 	defer db.Close()

// 	// We could do some stuff here.....
// 	err := createSchema(db, removeExisting)
// 	if err != nil {
// 		fmt.Println("There was an error creating the schema!")
// 		fmt.Println(err.Error())
// 	} else {
// 		fmt.Println("Schema creation is OK!")
// 	}
// }

// // Create the UnicorKitchen schema if it doesn't currently exist.
// func createSchema(db *pg.DB, removeExistingTables bool) error {

// 	//	orm.RegisterTable((*RecipeIngredient)(nil))

// 	models := []interface{}{
// 		(*Ingredient)(nil),
// 		(*Recipe)(nil),
// 		(*RecipeIngredient)(nil),
// 	}

// 	for _, model := range models {
// 		// We could individually check for tables here and create one by one if we wanted to...
// 		if removeExistingTables {
// 			db.Model(model).DropTable(&orm.DropTableOptions{
// 				IfExists: true,
// 				Cascade:  true,
// 			})
// 		}

// 		err := db.Model(model).CreateTable(&orm.CreateTableOptions{
// 			Temp:          false,
// 			IfNotExists:   true,
// 			FKConstraints: true,
// 		})
// 		if err != nil {
// 			return err
// 		}
// 	}

// 	return nil
// }
