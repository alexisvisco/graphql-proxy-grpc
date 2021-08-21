package main

import (
	"fmt"
	"github.com/alexisvisco/graphql-proxy-grpc/internal/gengql"
	"os"
)

func main() {
	if len(os.Args) != 2 {
		fmt.Println("graphql-proxy-grpc: first argument is required: expect directory of generated files with protoc plugin")
		return
	}

	dir := os.Args[1]
	gengql.GenerateGoGqlFromDir(dir)
}
