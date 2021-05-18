package main

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	//	"reflect"

	"io"
	"net/http"

	"github.com/drewCoSoftware/UnicornKitchen/database"
	"github.com/drewCoSoftware/UnicornKitchen/gql"
	"github.com/graphql-go/graphql"
	//	"github.com/drewCoSoftware/UnicornKitchen/gql"
)

var callCount int64 = 0

func main() {
	fmt.Println("What's cookin' in the Unicorn Kitchen?")

	schema, err := graphql.NewSchema(graphql.SchemaConfig{Query: gql.RecipeQuery})
	if err != nil {
		panic(err)
	}

	query := `{ recipe(name:"Electric Potato") { name, description } }`

	params := graphql.Params{Schema: schema, RequestString: query}
	r := graphql.Do(params)
	if len(r.Errors) > 0 {
		log.Fatalf("failed to execute graphql operation, errors: %+v", r.Errors)
	}

	rJSON, _ := json.Marshal(r)
	fmt.Printf("%s \n", rJSON)

	// Some reflection examples that we can hopefully use to auto-create gql defs....
	// x := database.RecipeIngredient{}
	// y := reflect.TypeOf(x)
	// fmt.Println(y)
	// fmt.Println(y.Kind())

	// fieldCount := y.NumField()
	// for i := 0; i < fieldCount; i++ {
	// 	fmt.Printf("Field: %s is a: %s\n", y.Field(i).Name, y.Field(i).Type)
	// }

	// // Set routing rules
	// http.HandleFunc("/", httpRoot)

	// //Use the default DefaultServeMux.
	// err := http.ListenAndServe(":8080", nil)
	// if err != nil {
	// 	log.Fatal(err)
	// }

	// x := gql.RecipeType
	// fmt.Printf("The name is: " + x.Name())
	// database.CreateDatabase()
	// addDefaultData()

	// After that we can look at setting up an HTTP endpoint to do some GraphQL queries
	// against our data....

}

func httpRoot(w http.ResponseWriter, r *http.Request) {

	io.WriteString(w, fmt.Sprintf("version %d", callCount))
	callCount++
}

func addDefaultData() {

	type defaultData struct {
		Ingredients []database.Ingredient
		Recipes     []database.Recipe
	}

	// Let's pull a recipe from file and add it that way...
	file, err := os.Open("data/default-data.json")
	if err != nil {
		panic(err)
	}
	defer file.Close()

	dec := json.NewDecoder(file)
	var dd defaultData
	if decodeErr := dec.Decode(&dd); decodeErr != nil {
		panic(decodeErr)
	}

	// Splat the default data into the database.
	for _, i := range dd.Ingredients {
		database.AddIngredient(&i)
	}

	for _, r := range dd.Recipes {
		database.AddRecipe(&r)
	}

}
