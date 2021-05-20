package gql

import (
	"errors"
	"fmt"
	"reflect"
	"strings"

	"github.com/graphql-go/graphql"
)

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
			msg := fmt.Sprintf("There was an error while resolving the type for property: %s.  You will need to map it manually!\n", field.Name)
			fmt.Println(msg)
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

// var resolvedTypes map[reflect.Type]*graphql.Object
// NOTE: We could do some work on a nested type resolver if we really wanted to.....
func resolveGraphqlType(t reflect.Type) (graphql.Output, error) {

	switch t.Kind() {
	case reflect.Array, reflect.Slice:

		return nil, errors.New("Arrays + slices are not supported at this time!")

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
