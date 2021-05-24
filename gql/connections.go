package gql

// This is an implementation of the Connection spec from the Relay project (https://relay.dev/graphql/connections.htm) specific
// to our application.

type gqlIngredientsConnection struct {
	Count    int                 `json:"count"`
	Edges    []gqlIngredientEdge `json:"edges"`
	PageInfo PageInfo            `json:"pageInfo"`
}

type gqlIngredientEdge struct {
	Node   gqlIngredient `json:"node"`
	Cursor string        `json:"cursor"`
}

type PageInfo struct {
	HasPreviousPage bool   `json:"hasPreviousPage"`
	HasNextPage     bool   `json:"hasNextPage"`
	StartCursor     string `json:"startCursor"`
	EndCursor       string `json:"endCursor"`
}
