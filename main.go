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

var gqlSchema graphql.Schema

func main() {
	fmt.Println("What's cookin' in the Unicorn Kitchen?")

	// Initialize our GraphQL stuff.
	gql.InitTypes()
	gql.InitQueries()

	initGqlSchema()

	// Query via GQL:
	// TestQuery()

	http.HandleFunc("/gql", handleGqlQuery)
	err := http.ListenAndServe(":8080", nil)
	if err != nil {
		log.Fatal(err)
	}
}

func initGqlSchema() {

	var err error
	gqlSchema, err = graphql.NewSchema(graphql.SchemaConfig{Query: gql.RecipeQuery})
	if err != nil {
		panic(err)
	}

}

// This acts a proxy for running our gql queries.  We now have an active GraphQL service!
func handleGqlQuery(w http.ResponseWriter, r *http.Request) {
	queryVals := r.URL.Query()
	gqlQuery := queryVals.Get("query")

	json, errs := gql.Query(gqlSchema, gqlQuery)
	if len(errs) > 0 {
		for _, err := range errs {
			io.Writer.Write(w, []byte(err.Message))
		}
	} else {
		io.Writer.Write(w, json)
	}

}

func TestQuery() {

	// get ingredients w/ paging...
	query := `{ ingredients(first:3, after:"-1") 
				{ count, edges { cursor, node { name, description } }, 
				  pageInfo { hasPreviousPage, hasNextPage, startCursor, endCursor } } }`

	rJSON, _ := gql.Query(gqlSchema, query)

	fmt.Printf("%s \n", rJSON)
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
