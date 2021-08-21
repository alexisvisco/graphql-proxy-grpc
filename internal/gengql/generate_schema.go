package gengql

import (
	"bytes"
	_ "embed"
	"github.com/alexisvisco/graphql-proxy-grpc/internal/errorx"
	"github.com/alexisvisco/graphql-proxy-grpc/internal/templatedata"
	"google.golang.org/protobuf/compiler/protogen"
	"reflect"
	"text/template"
)
import "github.com/alexisvisco/graphql-proxy-grpc/internal/virtual"

//go:embed schema.tmpl
var schemaTemplate string

func GenerateSchema(descriptor virtual.Descriptor, plugin *protogen.Plugin) (string, error) {
	schema, err := generateSchema(descriptor, plugin)
	if err != nil {
		return "", err
	}

	generateSchemaEmbed(plugin)

	return schema, nil
}

func generateSchema(descriptor virtual.Descriptor, plugin *protogen.Plugin) (string, error) {
	data := templatedata.NewGraphQLSchemaData(descriptor)

	var fns = template.FuncMap{
		"notlast": func(x int, a interface{}) bool {
			return x != reflect.ValueOf(a).Len()-1
		},
	}

	tmpl, err := template.New("schema.gql").Funcs(fns).Parse(schemaTemplate)
	if err != nil {
		return "", errorx.New().Str("origin-err", err.Error()).Err(errorx.ErrInvalidGraphqlSchemaTemplate)
	}

	buff := bytes.NewBufferString("")
	if err := tmpl.Execute(buff, data); err != nil {
		return "", errorx.New().Str("origin-err", err.Error()).Err(errorx.ErrUnableToCreateGraphqlSchema)
	}

	g := plugin.NewGeneratedFile("graphql/"+schemaFileName, "")
	_, err = g.Write(buff.Bytes())
	if err != nil {
		return "", errorx.New().Str("origin-err", err.Error()).Err(errorx.ErrUnableToUseProtoGeneratedFile)
	}

	return buff.String(), nil
}

func generateSchemaEmbed(plugin *protogen.Plugin) {
	g := plugin.NewGeneratedFile("graphql/"+schemaEmbedFile, "")
	g.P("package graphql")
	g.Import("embed")
	g.Write([]byte(`


//go:embed schema.graphql
var GraphQLSchema string
`))
}
