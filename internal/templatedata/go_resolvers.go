package templatedata

import (
	"github.com/99designs/gqlgen/codegen"
	"github.com/alexisvisco/graphql-proxy-grpc/internal/virtual"
)

type ClientNameData struct {
	Service string
	Package string
	Client  string
}

func NewClientNames(data *codegen.Data, descriptor virtual.Descriptor) (d []ClientNameData) {
	for path, pkg := range descriptor.Packages {
		for _, s := range pkg.Services {
			d = append(d, ClientNameData{
				Service: s.Name.Identifier.GoName,
				Package: data.Config.Packages.NameForPackage(path),
				Client:  s.Name.Identifier.GoName + "Client",
			})
		}
	}
	return
}
