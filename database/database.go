package database

import (
	"fmt"
	"github.com/drewCoSoftware/UnicornKitchen/settings"
	"github.com/go-pg/pg/v10"
	"github.com/go-pg/pg/v10/orm"
)

// type DAL
// {

// }

func CreateDatabase() {
	db := pg.Connect(settings.GetDatabaseOptions())
	defer db.Close()

	// Let's make the Database first...
	exists, err := dbExists(db, settings.DB_NAME)
	if err != nil {
		fmt.Println("There was an error checking for the database!")
		fmt.Println(err.Error())
		return
	}

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

}

func GetConnection() *pg.DB {
	res := pg.Connect(settings.GetDatabaseOptions())
	return res
}

func AddDefaultData() {
	fmt.Println("Adding some default data....")

	db := GetConnection()
	defer db.Close()

	// NOTE: A string list or input file would be a good idea...
	// Maybe just a json file?
	i1 := &Ingredient{
		Name: "Potato",
	}
	insertData(db, i1)

}

// func (db *pg.DB)  SaveRecipe(r *Recipe) {
// }

func SaveRecipe(db *pg.DB, r *Recipe) {

	fmt.Println("Saving Recipe: " + r.Name)

	// Iterate over the ingredients, saving each if it doesn't exist...
	for item := range r.Ingredients {
		ingredient := r.Ingredients[item].Ingredient
		//		fmt.Println(ingredient.Name)

		var i Ingredient
		var match = db.Model(&i).First()

		fmt.Println(ingredient.Name)
		fmt.Println(match)
	}
}

func insertData(db *pg.DB, data interface{}) {
	_, err := db.Model(data).Insert()
	if err != nil {
		panic(err)
	}
}

func CreateTables(removeExisting bool) {
	db := GetConnection()
	defer db.Close()

	// We could do some stuff here.....
	err := createSchema(db, removeExisting)
	if err != nil {
		fmt.Println("There was an error creating the schema!")
		fmt.Println(err.Error())
	} else {
		fmt.Println("Schema creation is OK!")
	}
}

func dbExists(db *pg.DB, dbName string) (bool, error) {
	// Check for the database....
	qr, err := db.Exec("SELECT datname FROM pg_catalog.pg_database WHERE datname = '" + dbName + "'")
	if err != nil {
		return false, err
		// fmt.Println("There was an error checking for the database!")
		// fmt.Println(err.Error())
	}

	return qr.RowsReturned() > 0, nil
}

// Create the UnicorKitchen schema if it doesn't currently exist.
func createSchema(db *pg.DB, removeExistingTables bool) error {

	//	orm.RegisterTable((*RecipeIngredient)(nil))

	models := []interface{}{
		(*Ingredient)(nil),
		(*Recipe)(nil),
		(*RecipeIngredient)(nil),
	}

	for _, model := range models {
		// We could individually check for tables here and create one by one if we wanted to...
		if removeExistingTables {
			db.Model(model).DropTable(&orm.DropTableOptions{
				IfExists: true,
				Cascade:  true,
			})
		}

		err := db.Model(model).CreateTable(&orm.CreateTableOptions{
			Temp:          false,
			IfNotExists:   true,
			FKConstraints: true,
		})
		if err != nil {
			return err
		}
	}

	return nil
}
