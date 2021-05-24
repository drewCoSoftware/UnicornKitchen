package gql

import (
	"reflect"

	"github.com/drewCoSoftware/UnicornKitchen/database"
	"github.com/graphql-go/graphql"
)

type gqlIngredient struct {
	Name        string `json:"name"`
	Description string `json:"description"`
	Amount      string `json:"amount"`
}

func Create(input database.Ingredient) gqlIngredient {
	res := gqlIngredient{
		Name:        input.Name,
		Description: input.Description,
	}
	return res
}

type gqlRecipe struct {
	Id           int64           `json:"id"`
	Name         string          `json:"name"`
	Description  string          `json:"description"`
	Ingredients  []gqlIngredient `json:"ingredients"`
	Instructions []string        `json:"instructions"`
}

var ingredientConnection *graphql.Object
var ingredientDef *graphql.Object
var recipeDef *graphql.Object

func InitTypes() {
	ingredientConnection = CreateGqlDefFromInstance("ingredientsConnection", gqlIngredientsConnection{})
	ingredientDef = CreateGqlDefFromInstance("ingredient", gqlIngredient{})
	recipeDef = CreateGqlDefFromInstance("recipe", gqlRecipe{})

	recipeDef.AddFieldConfig("ingredients", &graphql.Field{
		Type: graphql.NewList(ingredientDef),
	})

	recipeDef.AddFieldConfig("instructions", &graphql.Field{
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

}

func isArrayOrSlice(val interface{}) bool {
	v := reflect.ValueOf(val)
	k := v.Kind()
	return k == reflect.Slice || k == reflect.Array
}
