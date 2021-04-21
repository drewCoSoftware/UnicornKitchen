package database

import (
	//	"fmt"

	"fmt"

	"github.com/go-pg/pg/v10"
	"github.com/go-pg/pg/v10/orm"
	//	"github.com/go-pg/pg/v10/orm"
)

const DB_NAME string = "unicornkitchen"

func CreateDatabase() {
	db := pg.Connect(&pg.Options{
		// NOTE: This is dev data.  In real life, one would use 'secrets' from their container or server provider.
		User:     "postgres",
		Password: "abc123",
	})
	defer db.Close()

	// Let's make the Database first...
	exists, err := dbExists(db, DB_NAME)
	if err != nil {
		fmt.Println("There was an error checking for the database!")
		fmt.Println(err.Error())
		return
	}

	if !exists {
		fmt.Println("The database doesn't exist, so we will create it!")

		_, err := db.Exec("CREATE DATABASE " + DB_NAME)
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

func CreateTables() {

	db := pg.Connect(&pg.Options{
		// NOTE: This is dev data.  In real life, one would use 'secrets' from their container or server provider.
		User:     "postgres",
		Password: "abc123",
		Database: DB_NAME,
	})
	defer db.Close()

	// We could do some stuff here.....
	err := createSchema(db)
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
func createSchema(db *pg.DB) error {

	models := []interface{}{
		(*Ingredient)(nil),
		(*Recipe)(nil),
	}

	for _, model := range models {
		// We could individually check for tables here and create one by one if we wanted to...

		err := db.Model(model).CreateTable(&orm.CreateTableOptions{
			Temp:        false,
			IfNotExists: false,
		})
		if err != nil {
			return err
		}
	}

	return nil
}
