package database

import (
	"fmt"
)

type Ingredient struct {
	Id   int64
	Name string `pg:",unique,notnull"`
}

type Recipe struct {
	Id          int64
	Name        string `pg:",unique,notnull"`
	Ingredients []*IngredientEntry
}

type IngredientEntry struct {
	Id         int64
	Ingredient *Ingredient
	Amount     string // Quantities are encoded as strings.  i.e. '1' - '25mL' - '0.25#' and so on.  Not 100% ideal, but will do the trick for now.
}

// These are just type aliases.....
type (
	A1 = string
	A2 = int
)

func GetThing() (x int, y string) {
	//	decimal.De
	//	decimal.
	// res := &struct{
	// 	number = 100
	// }
	x = 100
	y = "abc"
	return x, y
}

// func GetStruct() (struct{ int, x string })

func check() {
	x, y := GetThing()
	fmt.Println(x)
	fmt.Println(y)
}

// func doit(x A1, y A2) A1 {
// 	z := 100 + y
// 	if z < 100 {
// 		return x
// 	}
// 	return x
// }
