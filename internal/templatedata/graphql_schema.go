package templatedata

import (
	"bytes"
	"fmt"
	"github.com/alexisvisco/graphql-proxy-grpc/internal/virtual"
	"go.buf.build/library/go/alexisvisco/graphql-proxy-grpc/graphqlpb/v1"
	"google.golang.org/protobuf/compiler/protogen"
	"sort"
	"strings"
)

type GraphQLSchemaData struct {
	Mutations []virtual.Rpc
	Queries   []virtual.Rpc
	Enums     []virtual.Enum
	Messages  []virtual.Message
}

func NewGraphQLSchemaData(descriptor virtual.Descriptor) GraphQLSchemaData {
	d := GraphQLSchemaData{
		Queries:   []virtual.Rpc{},
		Mutations: []virtual.Rpc{},
		Enums:     []virtual.Enum{},
		Messages:  []virtual.Message{},
	}
	for _, pkg := range descriptor.Packages {
		for _, service := range pkg.Services {
			for _, rpc := range service.RPCs {
				if rpc.Type == graphqlpbv1.GraphqlType_GRAPHQL_TYPE_QUERY {
					d.Queries = append(d.Queries, *rpc)
				} else {
					d.Mutations = append(d.Mutations, *rpc)
				}
			}
		}

		for _, enum := range pkg.Enums {
			d.Enums = append(d.Enums, *enum)
		}

		for _, message := range pkg.Messages {
			for _, field := range message.Fields {
				if field.Type.Kv != nil {
					msg := virtual.Message{
						Name: field.Type.Kv.Name,
						Fields: []*virtual.Field{
							{
								Name: virtual.Name{
									Identifier: protogen.GoIdent{GoName: "Key", GoImportPath: protogen.GoImportPath(pkg.Path)},
									GqlName:    "Key",
								},
								Type: virtual.Type{Native: &field.Type.Kv.Key},
							},
							{
								Name: virtual.Name{
									Identifier: protogen.GoIdent{GoName: "Value", GoImportPath: protogen.GoImportPath(pkg.Path)},
									GqlName:    "Value",
								},
								Type: field.Type.Kv.Value,
							},
						},
						IsInput: false,
					}

					d.Messages = append(d.Messages, msg)
					if field.Type.Kv.IsInput {
						d.Messages = append(d.Messages, messageIntoInputType(msg))
					}
				}
			}

			if message.IsEmpty() {
				continue
			}

			if message.IsInput {
				d.Messages = append(d.Messages, messageIntoInputType(*message))
			}

			msgCopy := *message
			msgCopy.IsInput = false
			d.Messages = append(d.Messages, msgCopy)
		}

		sort.Slice(d.Messages, func(i, j int) bool {
			return strings.Compare(d.Messages[i].Name.GqlName, d.Messages[j].Name.GqlName) == -1
		})

		sort.Slice(d.Queries, func(i, j int) bool {
			return strings.Compare(d.Queries[i].Name.GqlName, d.Queries[j].Name.GqlName) == -1
		})

		sort.Slice(d.Mutations, func(i, j int) bool {
			return strings.Compare(d.Mutations[i].Name.GqlName, d.Mutations[j].Name.GqlName) == -1
		})

		sort.Slice(d.Enums, func(i, j int) bool {
			return strings.Compare(d.Enums[i].Name.GqlName, d.Enums[j].Name.GqlName) == -1
		})
	}

	return d
}

func (d GraphQLSchemaData) FullGraphqlFieldName(isInput bool, f virtual.Field) string {
	field := bytes.NewBufferString(f.Name.GqlName)
	field.WriteString(": ")

	if f.Type.IsList || f.Type.Kv != nil {
		field.WriteByte('[')
	}
	if f.Type.Enum != nil {
		field.WriteString(f.Type.Enum.Name.GqlName)
	}
	if f.Type.Message != nil {
		if f.Type.Message.IsInput && isInput {
			field.WriteString(nameToInput(f.Type.Message.Name.GqlName))
		} else {
			field.WriteString(f.Type.Message.Name.GqlName)
		}
	}
	if f.Type.Native != nil {
		field.WriteString(f.Type.Native.GqlName)
	}
	if f.Type.Kv != nil {
		if f.Type.Kv.IsInput && isInput {
			field.WriteString(nameToInput(f.Type.Kv.Name.GqlName))
		} else {
			field.WriteString(f.Type.Kv.GqlName)
		}
	}

	if f.Type.IsList || f.Type.Kv != nil {
		field.WriteByte(']')
	}

	if f.Type.Kv != nil {
		field.WriteString(" @goField(forceResolver: true)")
	}

	return field.String()
}

func (d GraphQLSchemaData) FullGraphqlMethodName(r virtual.Rpc) string {
	method := bytes.NewBufferString(r.Name.GqlName)
	if r.Input.IsEmpty() {
		method.WriteString(": ")
	} else {
		msg := messageIntoInputType(*r.Input)
		method.WriteString(fmt.Sprintf("(in: %s): ", msg.Name.GqlName))
	}

	if r.Output.IsEmpty() {
		method.WriteString("Void")
	} else {
		method.WriteString(r.Output.Name.GqlName)
	}

	return method.String()
}

func (d GraphQLSchemaData) HasMutations() bool {
	return len(d.Mutations) > 0
}

func (d GraphQLSchemaData) HasQueries() bool {
	return len(d.Queries) > 0
}
