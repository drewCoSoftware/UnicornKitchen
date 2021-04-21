package main

import (
	"fmt"
	"github.com/drewCoSoftware/UnicornKitchen/database"
	// "github.com/drewCoSoftware/UnicornKitchen/ingredients"
	//	"github.com/go-pg/pg/v10"
	//	"github.com/go-pg/pg/v10/orm"
)

//import "github.com/google/go-cmp/cmp"

func main() {
	fmt.Println("What's cookin' in the Unicorn Kitchen?")

	// hasFish := ingredients.HasIngredient("fish")
	// fmt.Println(fmt.Sprintf("has fish?: %t", hasFish))

	// We want to create our database resources if they don't currently exist.
	database.CreateDatabase()

}
