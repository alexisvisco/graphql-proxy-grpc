package virtual

import (
	"encoding/json"
	"github.com/pkg/errors"
	"google.golang.org/protobuf/compiler/protogen"
	"path"
)

type Descriptor struct {
	Packages map[string]*Package
}

func NewDescriptor() *Descriptor {
	return &Descriptor{Packages: make(map[string]*Package)}
}

func (d *Descriptor) addPackage(pkg *Package) {
	d.Packages[pkg.String()] = pkg
}

func (d *Descriptor) getPackage(path string) *Package {
	pkg, ok := d.Packages[path]
	if !ok {
		return nil
	}

	return pkg
}

func (d *Descriptor) InspectProto(gen *protogen.Plugin) error {
	for _, f := range gen.Files {
		isGooglePackage := false

		if f.Desc.Package() == "graphqlpb.v1" && len(f.Extensions) != 0 {
			continue
		}

		if f.Desc.Package() == "google.protobuf" {
			isGooglePackage = true
		}

		pkg := d.getPackage(string(f.GoImportPath))
		if pkg == nil {
			pkg = NewPackage(
				d,
				string(f.Desc.FullName()),
				string(f.GoPackageName),
				string(f.GoImportPath),
				path.Dir(f.GeneratedFilenamePrefix),
			)
			d.addPackage(pkg)
		}

		for _, enum := range f.Enums {
			if isGooglePackage {
				break
			}

			pkg.getOrCreateEnum(enum)
		}

		for _, message := range f.Messages {
			if isGooglePackage && !WhitelistedMessagesGoogle[string(message.Desc.Name())] {
				continue
			}
			if _, err := pkg.getOrCreateMessage(message, d); err != nil {
				return err
			}
		}

		for _, service := range f.Services {
			if isGooglePackage {
				break
			}

			if _, err := pkg.createService(service, d); err != nil {
				return err
			}
		}
	}

	return nil
}

// GenerateDescriptorFileForGQLGen save the descriptor as json file to
// be then used in the gql gen par.
func (d *Descriptor) GenerateDescriptorFileForGQLGen(pl *protogen.Plugin) error {
	indent, err := json.MarshalIndent(d, "", " ")
	if err != nil {
		return errors.Wrap(err, "unable to marshal into json descriptor")
	}

	g := pl.NewGeneratedFile("graphql/descriptor.json", "")
	_, err = g.Write(indent)
	if err != nil {
		return errors.Wrap(err, "unable to write into generated file")
	}

	return nil
}
