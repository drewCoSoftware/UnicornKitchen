package gql

import (
	"errors"
	"fmt"
	"reflect"
	"strings"

	"github.com/graphql-go/graphql"
)

func CreateGqlDefFromType(name string, dataType reflect.Type) *graphql.Object {

	fieldCount := dataType.NumField()
	var fields graphql.Fields = make(map[string]*graphql.Field)

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

// Why do it by hand when we have reflection?
func CreateGqlDefFromInstance(name string, data interface{}) *graphql.Object {
	dataType := reflect.TypeOf(data)
	return CreateGqlDefFromType(name, dataType)
}

// var resolvedTypes map[reflect.Type]*graphql.Object
// NOTE: We could do some work on a nested type resolver if we really wanted to.....
func resolveGraphqlType(t reflect.Type) (graphql.Output, error) {

	switch t.Kind() {
	case reflect.Array, reflect.Slice:

		var gqlDef graphql.Output

		elemType := t.Elem()
		elemKind := elemType.Kind()
		if elemKind == reflect.Struct {
			// If we have a struct, then we need to create a new def from scratch.
			// In the real world we WOULD cache this to deal with circular dependencies + speed issues.
			// We also need a name for our type it appears.....

			// NOTE: This name thing is lkely to be the next big problem that we need to solve....
			name := elemType.Name()
			gqlDef = CreateGqlDefFromType(name, elemType)
		} else {
			// This is probably a normal type.....
			if gqlDef, err := resolveGraphqlType(elemType); err != nil {
				return nil, err
			} else {
				return gqlDef, nil
			}
		}

		res := graphql.NewList(gqlDef)
		return res, nil

	case reflect.Struct:
		name := t.Name()
		res := CreateGqlDefFromType(name, t)
		return res, nil

	default:

		switch t.Name() {
		case "int", "int64":
			return graphql.Int, nil
		case "string":
			return graphql.String, nil
		case "bool":
			return graphql.Boolean, nil
		default:
			return nil, errors.New("Unknown graphql type name: " + t.Name())
		}

	}

}
