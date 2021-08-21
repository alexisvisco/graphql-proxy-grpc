package gengql

import (
	"encoding/json"
	"fmt"
	"github.com/99designs/gqlgen/api"
	"github.com/99designs/gqlgen/codegen/config"
	"github.com/alexisvisco/graphql-proxy-grpc/internal/gengql/pluginresolvers"
	"github.com/alexisvisco/graphql-proxy-grpc/internal/virtual"
	"github.com/vektah/gqlparser/v2/ast"
	"io/ioutil"
	"path"
)

const (
	graphqlSchemaFileName   = "schema.graphql"
	protoDescriptorFileName = "descriptor.json"
)

func GenerateGoGqlFromDir(dir string) {
	protoDescriptorAsJson, err := ioutil.ReadFile(path.Join(dir, protoDescriptorFileName))
	if err != nil {
		fmt.Printf("graphql-proxy-grpc: error while reading proto descriptor file: %v\n", err)
		return
	}

	protoDescriptor := &virtual.Descriptor{}
	if err := json.Unmarshal(protoDescriptorAsJson, protoDescriptor); err != nil {
		fmt.Printf("graphql-proxy-grpc: error while unmarshaling proto file: %v\n", err)
		return
	}

	graphqlSchema, err := ioutil.ReadFile(path.Join(dir, graphqlSchemaFileName))
	if err != nil {
		fmt.Printf("graphql-proxy-grpc: error while reading graphql schema file: %v\n", err)
		return
	}

	if err := GenerateGoGql(*protoDescriptor, string(graphqlSchema)); err != nil {
		fmt.Printf("graphql-proxy-grpc: unable to generate graphql resolver via gqlgen: %v\n", err)
		return
	}
}

func GenerateGoGql(descriptor virtual.Descriptor, schema string) error {
	cfg := getGQLGenConfig(schema)

	fmt.Println(cfg, descriptor, schema)
	if err := api.Generate(cfg, api.AddPlugin(pluginresolvers.New(descriptor))); err != nil {
		return err
	}

	return nil
}

func getGQLGenConfig(schema string) *config.Config {
	cfg := &config.Config{
		Exec: config.PackageConfig{
			Filename: path.Join(dir, "exec_gen.go"),
			Package:  "graphql",
		},
		Model: config.PackageConfig{
			Filename: path.Join(dir, "models_gen.go"),
			Package:  "graphql",
		},
		Resolver: config.ResolverConfig{
			Filename: dir + "/resolvers.go",
			Package:  "graphql",
			Layout:   config.LayoutSingleFile,
		},
		Models: config.TypeMap{
			"Int64":   config.TypeMapEntry{Model: config.StringList{"github.com/99designs/gqlgen/graphql.Int64"}},
			"Uint32":  config.TypeMapEntry{Model: config.StringList{"github.com/alexisvisco/graphql-proxy-grpc/types.Uint32"}},
			"Uint64":  config.TypeMapEntry{Model: config.StringList{"github.com/alexisvisco/graphql-proxy-grpc/types.Uint64"}},
			"Float32": config.TypeMapEntry{Model: config.StringList{"github.com/alexisvisco/graphql-proxy-grpc/types.Float32"}},
			"Bytes":   config.TypeMapEntry{Model: config.StringList{"github.com/alexisvisco/graphql-proxy-grpc/types.Bytes"}},
		},
		Directives: map[string]config.DirectiveConfig{},

		Sources: []*ast.Source{
			{Name: "schema.graphql", Input: schema},
		},

		// Validation must be skipped because protoc generate go code after gqlgen
		SkipValidation: true,
	}

	return cfg
}
