package virtual

import (
	"github.com/alexisvisco/graphql-proxy-grpc/internal/errorx"
	changecase "github.com/ku/go-change-case"
	"google.golang.org/protobuf/compiler/protogen"
	"google.golang.org/protobuf/reflect/protoreflect"
)

type Field struct {
	Name Name
	Type Type
}

func createField(field *protogen.Field, descriptor *Descriptor) (*Field, error) {
	vf := &Field{
		Name: Name{
			Identifier: field.GoIdent,
			GqlName:    changecase.Camel(field.GoName),
		},
		Type: Type{},
	}

	var (
		fieldType Type
		err       error
	)

	if field.Desc.Kind() == protoreflect.EnumKind {
		fieldType, err = createEnumType(field, descriptor)
	} else if goType, ok := nativeConversionFromProtoToGoType[field.Desc.Kind().String()]; ok {
		vf.Type = createTypeNative(field, goType)
		return vf, nil
	} else if field.Desc.IsMap() {
		fieldType, err = createKeyValueType(field, descriptor)
	} else if field.Desc.Kind() == protoreflect.MessageKind {
		fieldType, err = createMessageType(field, descriptor)
	}

	if err != nil {
		return nil, err
	}

	if fieldType.IsEmptyMessage() {
		return nil, errorx.New().Err(errorx.ErrEmptyMessageNotSupportedInField).
			Str("pkg", string(field.Desc.ParentFile().Package())).
			Str("msg", string(field.Parent.Desc.Name())).
			Str("field", string(field.Desc.Name()))
	}

	vf.Type = fieldType
	return vf, nil
}
