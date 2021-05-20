package gql

import (
	"errors"
	"fmt"
	"reflect"

	"github.com/drewCoSoftware/UnicornKitchen/database"
	"github.com/graphql-go/graphql"
)

type gqlIngredient struct {
	Name        string `json:"name"`
	Description string `json:"description"`
	Amount      string `json:"amount"`
}

type gqlRecipe struct {
	Id           int64           `json:"id"`
	Name         string          `json:"name"`
	Description  string          `json:"description"`
	Ingredients  []gqlIngredient `json:"ingredients"`
	Instructions []string        `json:"instructions"`
}

var RecipeQuery *graphql.Object
var IngredientQuery *graphql.Object

func InitTypes() {
	ii2 := CreateGqlObjectFromType("ingredient", gqlIngredient{})
	rt2 := CreateGqlObjectFromType("recipe", gqlRecipe{})

	rt2.AddFieldConfig("ingredients", &graphql.Field{
		Type: graphql.NewList(ii2),
	})

	rt2.AddFieldConfig("instructions", &graphql.Field{
		Type: graphql.NewList(graphql.String),
		Resolve: func(p graphql.ResolveParams) (interface{}, error) {
			if recipe, ok := p.Source.(*gqlRecipe); ok {
				res := database.GetInstructions(recipe.Id)
				return res, nil
			} else {
				return nil, nil
			}
		},
	})

	RecipeQuery = graphql.NewObject(graphql.ObjectConfig{
		Name: "Query",
		Fields: graphql.Fields{
			// TODO: A way to get all of the ingredients, with pagination. (via connect)
			"ingredients": &graphql.Field{
				Type: graphql.NewList(ii2),
				Resolve: func(p graphql.ResolveParams) (interface{}, error) {
					match := database.GetIngredients()
					size := len(match)

					res := make([]*gqlIngredient, size)
					for i := 0; i < size; i++ {
						// NOTE: A DTOMapper might be nice here, but we probably
						// can't make an efficient one in GO (emit ASM or similar at runtime)
						res[i] = &gqlIngredient{
							Name:        match[i].Name,
							Description: match[i].Description,
						}
					}

					return res, nil
				},
			},
			"ingredient": &graphql.Field{
				Type: ii2,
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
				Type: rt2,
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

					// return &gqlRecipe{
					// 	Name:        "test-recipe",
					// 	Description: "test-descriptions",
					// }, nil
				},
			},
		},
	})

}

func isArrayOrSlice(val interface{}) bool {
	v := reflect.ValueOf(val)
	k := v.Kind()
	return k == reflect.Slice || k == reflect.Array
}
