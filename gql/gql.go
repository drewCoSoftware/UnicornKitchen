package gql

import (
	"errors"
	"fmt"
	"reflect"
	"strings"

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

var resolvedTypes map[reflect.Type]*graphql.Object

var ii2 = CreateGqlObjectFromType("ingredient", gqlIngredient{})

var rt2 = graphql.NewObject(graphql.ObjectConfig{
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
		"instructions": &graphql.Field{
			Type: graphql.NewList(graphql.String),
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				if recipe, ok := p.Source.(*gqlRecipe); ok {
					res := database.GetInstructions(recipe.Id)
					return res, nil
				} else {
					return nil, nil
				}
			},
		},
		"ingredients": &graphql.Field{
			Type: graphql.NewList(ii2),
			// NOTE: If we wanted to come up with a special way to resolve the recipe ingredients:
			// This could do stuff like cache popular ingredients and their description, etc.
			// Essentially we could fine tune our system based on its performance characterisitics with
			// these resolvers.
			//
			// Resolve: func(p graphql.ResolveParams) (interface{}, error) {
			// 	if recipe, ok := p.Source.(*gqlRecipe); ok {
			// 		fmt.Println("I am resolving some ingredients! for recipe: ", recipe.Name)
			// 	} else {
			// 		fmt.Println("I can't resolve the ingredients! (src =", p.Source, ")")
			// 	}
			// 	return nil, nil
			// },
		},
	},
})

// Scratchpad function for testing out reflection....
func GqlReflect() {
	var thing []gqlIngredient
	i := gqlIngredient{}
	thing = append(thing, i)

	if isArrayOrSlice(thing) {
		arType := reflect.TypeOf(thing).Elem()
		fmt.Println("The array type is: ", arType.Name())
	}
}

func isArrayOrSlice(val interface{}) bool {
	v := reflect.ValueOf(val)
	k := v.Kind()
	return k == reflect.Slice || k == reflect.Array
}

// Why do it by hand when we have reflection?
func CreateGqlObjectFromType(name string, data interface{}) *graphql.Object {
	dataType := reflect.TypeOf(data)

	fieldCount := dataType.NumField()
	var fields graphql.Fields = make(map[string]*graphql.Field)

	// NOTE: This doesn't work because of how goes type system works.
	// Go typedefs ARE NOT simple aliases.
	// fields := make(map[string]*graphql.Field)

	for i := 0; i < fieldCount; i++ {
		field := dataType.Field(i)

		// Check any tags we have applied....
		tag := field.Tag.Get("gql")
		if tag == "ignore" {
			continue
		}

		fType, err := resolveGraphqlType(field.Type)
		if err != nil {
			panic(err)
		} else {
			// NOTE: We would want to use the json name here.....
			fields[strings.ToLower(field.Name)] = &graphql.Field{
				Type: fType,
			}
		}

	}

	res := graphql.NewObject(graphql.ObjectConfig{
		Name:   name,
		Fields: fields,
	})
	return res
}

func resolveGraphqlType(t reflect.Type) (graphql.Output, error) {

	switch t.Kind() {
	case reflect.Array, reflect.Slice:

		return nil, errors.New("Arrays + slices are not supported at this time!")
		// aType := t.Elem()
		// if gqlType, err := resolveGraphqlType(aType); err != nil {
		// 	panic(err)
		// }

		// return graphql.NewList()

	default:

		switch t.Name() {
		case "int64":
			return graphql.Int, nil
		case "string":
			return graphql.String, nil
		default:
			return nil, errors.New("Unknown graphql type name: " + t.Name())
		}

	}

}

var RecipeQuery = graphql.NewObject(graphql.ObjectConfig{
	Name: "Query",
	Fields: graphql.Fields{
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
