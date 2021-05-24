package gql

import (
	"errors"
	"fmt"
	"strconv"

	"github.com/drewCoSoftware/UnicornKitchen/database"
	"github.com/graphql-go/graphql"
)

var RecipeQuery *graphql.Object
var IngredientQuery *graphql.Object

func resolvePageArgs(p graphql.ResolveParams) database.PageArgs {
	res := database.PageArgs{
		Before: asString(p.Args["before"]),
		After:  asString(p.Args["after"]),
	}
	res.First = asInt(p.Args["first"])
	res.Last = asInt(p.Args["last"])

	return res
}

func asInt(input interface{}) int {
	if input == nil {
		return 0
	}
	res := input.(int)
	return res
}

func asString(input interface{}) string {
	if input == nil {
		return ""
	}
	res := input.(string)
	return res
}

func tryParseInt64(input string, fallback int64) int64 {
	if val, err := strconv.ParseInt(input, 2, 64); err == nil {
		return val
	} else {
		fmt.Println("Could not parse an int64 from input string: ", input, ". The fallback value will be used!")
		return fallback
	}
}

func InitQueries() {

	RecipeQuery = graphql.NewObject(graphql.ObjectConfig{
		Name: "Query",
		Fields: graphql.Fields{

			// TODO: A way to get all of the ingredients, with pagination. (via connect)
			"ingredients": &graphql.Field{
				Type: ingredientConnection,
				Args: graphql.FieldConfigArgument{
					"first": &graphql.ArgumentConfig{
						Type: graphql.Int,
					},
					"after": &graphql.ArgumentConfig{
						Type: graphql.String,
					},
					"last": &graphql.ArgumentConfig{
						Type: graphql.Int,
					},
					"before": &graphql.ArgumentConfig{
						Type: graphql.String,
					},
				},
				Resolve: func(p graphql.ResolveParams) (interface{}, error) {

					// NOTE: Let's ignore the arguments for a sec...
					pageArgs := resolvePageArgs(p)

					match := database.GetIngredients(&pageArgs)
					size := len(match)

					res := gqlIngredientsConnection{
						Count: size,
						Edges: make([]gqlIngredientEdge, size),
					}
					for i, item := range match {
						edge := gqlIngredientEdge{
							Node:   Create(item),
							Cursor: fmt.Sprintf("%d", item.Id),
						}
						res.Edges[i] = edge
					}

					startCursor := res.Edges[0].Cursor
					endCursor := res.Edges[size-1].Cursor

					first, last := database.GetCursorBoundaries("ingredients", "ingredientid")
					res.PageInfo = PageInfo{
						StartCursor:     startCursor,
						EndCursor:       endCursor,
						HasPreviousPage: tryParseInt64(startCursor, 0) > first,
						HasNextPage:     tryParseInt64(endCursor, 0) < last,
					}

					return res, nil
				},
			},
			"ingredient": &graphql.Field{
				Type: ingredientDef,
				Args: graphql.FieldConfigArgument{
					"name": &graphql.ArgumentConfig{
						Type: graphql.String,
					},
				},
				Resolve: func(p graphql.ResolveParams) (interface{}, error) {
					name, ok := p.Args["name"].(string)
					if ok {
						dbIngredient := database.GetIngredient(name)
						res := &gqlIngredient{
							Name:        dbIngredient.Name,
							Description: dbIngredient.Description,
						}
						return res, nil
					} else {
						return nil, errors.New(fmt.Sprintf("There is no ingredient named: %s", name))
					}
				},
			},

			"recipe": &graphql.Field{
				Type: recipeDef,
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
							Id:          dbRecipe.Id,
							Name:        dbRecipe.Name,
							Description: dbRecipe.Description,
						}

						res.Ingredients = make([]gqlIngredient, len(dbRecipe.Ingredients))
						for i, ingredient := range dbRecipe.Ingredients {
							res.Ingredients[i] = gqlIngredient{
								Name:        ingredient.Ingredient.Name,
								Description: ingredient.Ingredient.Description,
								Amount:      ingredient.IngredientAmount,
							}
						}

						res.Instructions = make([]string, len(dbRecipe.Instructions))
						for i, instr := range dbRecipe.Instructions {
							res.Instructions[i] = instr.Content
						}

						return res, nil
					} else {
						return nil, errors.New("There is no argument for 'name'!")
					}
				},
			},
		},
	})

}
