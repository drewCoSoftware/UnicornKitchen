package gql

import (
	"errors"

	"github.com/drewCoSoftware/UnicornKitchen/database"
	"github.com/graphql-go/graphql"
)

type gqlRecipe struct {
	Id          int64  `json:"id"`
	Name        string `json:"name"`
	Description string `json:"description"`
}

var RecipeType = graphql.NewObject(graphql.ObjectConfig{
	Name: "recipe",
	Fields: graphql.Fields{
		"id": &graphql.Field{
			Type: graphql.Int,
		},
		"name": &graphql.Field{
			Type: graphql.String,
		},
		"description": &graphql.Field{
			Type: graphql.String,
		},
	},
})

var RecipeQuery = graphql.NewObject(graphql.ObjectConfig{
	Name: "Query",
	Fields: graphql.Fields{
		"recipe": &graphql.Field{
			Type: RecipeType,
			Args: graphql.FieldConfigArgument{
				"name": &graphql.ArgumentConfig{
					Type: graphql.String,
				},
			},
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				name, ok := p.Args["name"].(string)
				if ok {
					// This is where we will find the recipe name in the database...
					dbRecipe := database.GetRecipe(name)
					res := &gqlRecipe{
						Name:        dbRecipe.Name,
						Description: dbRecipe.Description,
					}

					return res, nil
				} else {
					return nil, errors.New("There is no argument for 'name'!")
				}

				// return &gqlRecipe{
				// 	Name:        "test-recipe",
				// 	Description: "test-descriptions",
				// }, nil
			},
		},
	},
})

// func CreateGQLObject(type t) {

// }

//var x map[int]string
// x := map[int]string {
// 	1, "abc",
// }
// fmt.PrintLn(x)
