package main

import (
	"encoding/json"
	"fmt"
	"github.com/drewCoSoftware/UnicornKitchen/database"
	"github.com/drewCoSoftware/UnicornKitchen/gql"
	"github.com/graphql-go/graphql"
	"io"
	"log"
	"net/http"
	"os"
)

var callCount int64 = 0

func main() {
	fmt.Println("What's cookin' in the Unicorn Kitchen?")

	gql.InitTypes()

	// Query via GQL:
	TestQuery()

}

func TestQuery() {
	schema, err := graphql.NewSchema(graphql.SchemaConfig{Query: gql.RecipeQuery})
	if err != nil {
		panic(err)
	}

	// query := `{ recipe(name:"Electric Potato") { name, description, ingredients {name, amount}, instructions } }`

	// query = `{ ingredient(name:"Potato") { name, description } }`
	query := `{ ingredients { name, description } }`

	params := graphql.Params{Schema: schema, RequestString: query}
	r := graphql.Do(params)
	if len(r.Errors) > 0 {
		log.Fatalf("failed to execute graphql operation, errors: %+v", r.Errors)
	}

	rJSON, _ := json.Marshal(r)
	fmt.Printf("%s \n", rJSON)
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
