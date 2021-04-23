package database

import "fmt"

type Ingredient struct {
	Id   int64
	Name string `pg:",unique"`
}

type Recipe struct {
	Id   int64
	Name string `pg:",unique"`
}

// These are just type aliases.....
type (
	A1 = string
	A2 = int
)

func GetThing() (x int, y string) {
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
