package main

import "fmt"

//import "github.com/google/go-cmp/cmp"

import "github.com/drewCoSoftware/UnicornKitchen/ingredients"

func main() {
	fmt.Println("What's cookin' in the Unicorn Kitchen?")
	//	fmt.Println(cmp.Diff("Hello World", "Hello Go"))

	hasFish := ingredients.HasIngredient("fish")
	fmt.Println(fmt.Sprintf("has fish?: %t", hasFish))
	//	fmt.Println("has fish?" + hasFish)

}
