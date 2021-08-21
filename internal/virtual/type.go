package virtual

import (
	"bytes"
	"fmt"
	"github.com/alexisvisco/graphql-proxy-grpc/internal/errorx"
	changecase "github.com/ku/go-change-case"
	"google.golang.org/protobuf/compiler/protogen"
)

type Type struct {
	Enum    *Enum
	Message *Message
	Kv      *KeyValue
	Native  *Native

	IsList bool
}

func createKeyValueType(field *protogen.Field, descriptor *Descriptor) (Type, error) {
	t := Type{}

	// we are creating a custom go identifier because we will create this go structure
	// Since graphql does not provide map we are create a list of typed k/v pair
	name := fmt.Sprintf("KeyValue%s", changecase.Pascal(field.GoIdent.GoName))

	t.Kv = &KeyValue{
		Name: Name{

			Identifier: protogen.GoIdent{
				GoName:       name,
				GoImportPath: field.Parent.GoIdent.GoImportPath,
			},
			GqlName: name,
		},
	}

	err := errorx.New().
		Str("msg", string(field.Parent.Desc.Name())).
		Str("field", string(field.Desc.Name())).
		Str("pkg", string(field.Desc.ParentFile().Package()))

	protoKeyKindType := field.Desc.MapKey().Kind().String()
	if goType, ok := nativeConversionFromProtoToGoType[protoKeyKindType]; ok {
		t.Kv.Key = *NewNative(goType)
	} else {
		return Type{}, err.Err(errorx.ErrProtoMapKeyTypeInvalid).Str("map-key-type-found", protoKeyKindType)
	}

	if len(field.Message.Fields) < 2 {
		return Type{}, err.Err(errorx.ErrProtoMapValueTypeInvalid).
			Int("map-message-fields-len", len(field.Message.Fields))
	}

	f, e := createField(field.Message.Fields[1], descriptor)
	if e != nil {
		return Type{}, e
	}

	t.Kv.Value = f.Type

	return t, nil
}

func createEnumType(field *protogen.Field, descriptor *Descriptor) (Type, error) {
	fieldEnum := field.Enum

	pkg := descriptor.getPackage(string(field.GoIdent.GoImportPath))
	if pkg == nil {
		return Type{}, errorx.New().
			Err(errorx.ErrProtoPackageNotLoaded).
			Str("msg", string(field.Parent.Desc.Name())).
			Str("field", field.GoName).
			Str("enum-to-load", string(fieldEnum.Desc.Name())).
			Str("pkg-to-load", string(field.Desc.ParentFile().Package()))
	}

	return Type{
		Enum:   pkg.getOrCreateEnum(fieldEnum),
		IsList: field.Desc.IsList(),
	}, nil
}

func createMessageType(field *protogen.Field, descriptor *Descriptor) (Type, error) {
	fieldMessage := field.Message

	pkg := descriptor.getPackage(string(fieldMessage.GoIdent.GoImportPath))
	if pkg == nil {
		return Type{}, errorx.New().
			Err(errorx.ErrProtoPackageNotLoaded).
			Str("msg", string(field.Parent.Desc.Name())).
			Str("field", field.GoName).
			Str("message-to-load", string(fieldMessage.Desc.Name())).
			Str("pkg-to-load", string(field.Desc.ParentFile().Package()))
	}

	msg, err := pkg.getOrCreateMessage(fieldMessage, descriptor)
	if err != nil {
		return Type{}, err
	}
	return Type{
		Message: msg,
		IsList:  field.Desc.IsList(),
	}, nil
}

func createTypeNative(field *protogen.Field, goType string) Type {
	return Type{
		Native: NewNative(goType),
		IsList: field.Desc.IsList(),
	}
}

type KeyValue struct {
	Name

	Key   Native
	Value Type

	IsInput bool
}

type Native struct {
	GoName  string
	GqlName string
}

func NewNative(goType string) *Native {
	return &Native{
		GoName:  goType,
		GqlName: nativeConversionFromGoToGqlType[goType],
	}
}

func (t Type) IsKv() bool {
	return t.Kv != nil
}

func (t Type) IsEmptyMessage() bool {
	return t.Message != nil && t.Message.IsEmpty()
}

func (t Type) TemplateTypeGql() string {
	buffer := bytes.NewBufferString("")
	if t.IsList || t.Kv != nil {
		buffer.WriteByte('[')
	}
	if t.Enum != nil {
		buffer.WriteString(t.Enum.Name.GqlName)
	}
	if t.Message != nil {
		buffer.WriteString(t.Message.Name.GqlName)
	}
	if t.Native != nil {
		buffer.WriteString(t.Native.GqlName)
	}
	if t.Kv != nil {
		buffer.WriteString(t.Kv.GqlName)
	}

	if t.IsList || t.Kv != nil {
		buffer.WriteByte(']')
	}

	return buffer.String()
}

func (t Type) IdentifierForKvValue() interface{} {
	if t.Enum != nil {
		return t.Enum.Name.Identifier
	}
	if t.Message != nil {
		return t.Message.Name.Identifier
	}
	if t.Native != nil {
		return t.Native.GoName
	}

	return ""
}

func (t Type) PointerForKeyValue() string {
	if t.Message != nil {
		return "*"
	}

	return ""
}
