package virtual

import (
	changecase "github.com/ku/go-change-case"
	"github.com/rs/zerolog/log"
	"google.golang.org/protobuf/compiler/protogen"
)

type Enum struct {
	Name Name

	Values []string
}

func (v Enum) String() string {
	return v.Name.String()
}

func (pkg *Package) getOrCreateEnum(enum *protogen.Enum) *Enum {
	ve := &Enum{
		Name: Name{
			Identifier: enum.GoIdent,
			GqlName:    changecase.Pascal(enum.GoIdent.GoName),
		},
		Values: nil,
	}

	if e := pkg.getEnum(ve.Name); e == nil {
		values := make([]string, len(enum.Values))
		for i, value := range enum.Values {
			values[i] = string(value.Desc.Name())
		}
		ve.Values = values

		pkg.addEnum(ve)
		log.Debug().Msgf("register enum %s in package %s", enum.Desc.Name(), enum.Desc.ParentFile().Package())
		return ve
	} else {
		return e
	}
}
