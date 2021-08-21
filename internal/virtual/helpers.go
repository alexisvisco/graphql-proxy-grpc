package virtual

import (
	"fmt"
	"google.golang.org/protobuf/compiler/protogen"
)

var nativeConversionFromProtoToGoType = map[string]string{
	"double":   "float64",
	"float":    "float32",
	"int32":    "int32",
	"int64":    "int64",
	"uint32":   "uint32",
	"uint64":   "uint64",
	"sint32":   "int32",
	"sint64":   "int64",
	"fixed32":  "uint32",
	"fixed64":  "uint64",
	"sfixed32": "int32",
	"sfixed64": "int64",
	"bool":     "bool",
	"string":   "string",
	"bytes":    "[]byte",
}

var nativeConversionFromGoToGqlType = map[string]string{
	"float64": "Float",
	"float32": "Float32",
	"int32":   "Int",
	"int64":   "Int64",
	"uint32":  "Uint32",
	"uint64":  "Uint64",
	"bool":    "Boolean",
	"string":  "String",
	"[]byte":  "Bytes",
}

type Name struct {
	Identifier protogen.GoIdent
	GqlName    string
}

func (n Name) String() string {
	return fmt.Sprintf("go: %q ; gql: %q", n.Identifier.GoName, n.GqlName)
}

func (n Name) ImportPath() string {
	return string(n.Identifier.GoImportPath)
}
