package graphql

import (
	_ "embed"
)

//go:embed schema.graphql
var GraphQLSchema string

