package pluginresolvers

import (
	"bytes"
	_ "embed"
	"fmt"
	"github.com/99designs/gqlgen/codegen"
	"github.com/99designs/gqlgen/codegen/config"
	"github.com/99designs/gqlgen/codegen/templates"
	"github.com/99designs/gqlgen/plugin"
	"github.com/alexisvisco/graphql-proxy-grpc/internal/gengql/pluginresolvers/rewrite"
	"github.com/alexisvisco/graphql-proxy-grpc/internal/templatedata"
	"github.com/alexisvisco/graphql-proxy-grpc/internal/virtual"
	changecase "github.com/ku/go-change-case"
	"os"
	"strings"
	"text/template"
)

//go:embed resolvers.tmpl
var resolversTemplate string

//go:embed field_resolver.tmpl
var fieldResolverTemplate string

//go:embed service_resolver.tmpl
var serviceResolverTemplate string

func New(descriptor virtual.Descriptor) plugin.Plugin {
	return &Plugin{
		descriptor: descriptor,
	}
}

type Plugin struct {
	descriptor virtual.Descriptor
}

var _ plugin.CodeGenerator = &Plugin{}

func (m *Plugin) Name() string {
	return "resolvergen"
}

func (m *Plugin) GenerateCode(data *codegen.Data) error {
	if !data.Config.Resolver.IsDefined() {
		return nil
	}

	switch data.Config.Resolver.Layout {
	case config.LayoutSingleFile:
		return m.generateSingleFile(data)
	}

	return nil
}

func (m *Plugin) generateSingleFile(data *codegen.Data) error {
	file := file{}

	if _, err := os.Stat(data.Config.Resolver.Filename); err == nil {
		// file already exists and we dont support updating resolvers with layout = single so just return
		return nil
	}

	services := templatedata.NewClientNames(data, m.descriptor)
	var implSolutions = createImplementationSolutions(data.Config, m.descriptor, &file)

	for _, o := range data.Objects {
		if o.HasResolvers() {
			file.Objects = append(file.Objects, o)
		}
		for _, f := range o.Fields {
			if !f.IsResolver {
				continue
			}

			namespace, name := f.Object.Name, f.GoFieldName
			if namespace == "Query" || namespace == "Mutation" {
				for _, service := range services {
					if strings.HasPrefix(name, service.Service) {
						namespace = service.Service
						name = strings.TrimPrefix(name, service.Service)
						break
					}
				}
			}

			implementation := bytes.NewBufferString("panic(\"not implemented\")")
			namespaceAndName := fmt.Sprintf("%s.%s", namespace, name)
			implData, ok := implSolutions[namespaceAndName]
			if ok {
				implementation.Reset()
				tmpl := serviceResolverTemplate
				var data interface{} = implData.serviceResolverImplementation
				if implData.fieldResolverImplementation != nil {
					tmpl = fieldResolverTemplate
					data = implData.fieldResolverImplementation
				}

				parse, err := template.New(namespaceAndName).Parse(tmpl)
				if err != nil {
					return err
				}

				if err := parse.Execute(implementation, data); err != nil {
					return err
				}
			}

			resolver := resolver{o, f, implementation.String()}
			file.Resolvers = append(file.Resolvers, &resolver)
		}
	}

	resolverBuild := &resolverBuild{
		file:         &file,
		PackageName:  data.Config.Resolver.Package,
		ResolverType: data.Config.Resolver.Type,
		HasRoot:      true,
		Services:     services,
	}

	return templates.Render(templates.Options{
		Template:    resolversTemplate,
		PackageName: data.Config.Resolver.Package,
		FileNotice:  `// THIS CODE IS A STARTING POINT ONLY. IT WILL NOT BE UPDATED WITH SCHEMA CHANGES.`,
		Filename:    data.Config.Resolver.Filename,
		Data:        resolverBuild,
		Packages:    data.Config.Packages,
	})
}

func (m *Plugin) getAllServices(data *codegen.Data) (services []serviceWithPackageName) {
	for path, pkg := range m.descriptor.Packages {
		for _, service := range pkg.Services {
			services = append(services, serviceWithPackageName{service, data.Config.Packages.NameForPackage(path)})
		}
	}
	return services
}

type resolverBuild struct {
	*file
	HasRoot      bool
	PackageName  string
	ResolverType string
	Services     []templatedata.ClientNameData
}

type file struct {
	// These are separated because the type definition of the resolver object may live in a different file from the
	//resolver method implementations, for example when extending a type in a different graphql schema file
	Objects         []*codegen.Object
	Resolvers       []*resolver
	imports         []rewrite.Import
	RemainingSource string
}

type serviceWithPackageName struct {
	*virtual.Service
	Pkg string
}

func (f *file) Imports() string {
	for _, imp := range f.imports {
		if imp.Alias == "" {
			_, _ = templates.CurrentImports.Reserve(imp.ImportPath)
		} else {
			_, _ = templates.CurrentImports.Reserve(imp.ImportPath, imp.Alias)
		}
	}
	return ""
}

type resolver struct {
	Object         *codegen.Object
	Field          *codegen.Field
	Implementation string
}

type ServiceResolverImplementation struct {
	ServiceName       string
	ServiceMethodName string
	EmptyMessagePkg   string
	EmptyMessageName  string
	IsEmptyOutput     bool
	IsEmptyInput      bool
}

type FieldResolverImplementation struct {
	PkgName          string
	KeyValueTypeName string
	FieldName        string
}

type resolverImplementation struct {
	serviceResolverImplementation *ServiceResolverImplementation
	fieldResolverImplementation   *FieldResolverImplementation
}

type ImplementationSolutions map[string]resolverImplementation

func createImplementationSolutions(cfg *config.Config, descriptor virtual.Descriptor, file *file) ImplementationSolutions {
	is := ImplementationSolutions{}
	for path, pkg := range descriptor.Packages {
		packageName := cfg.Packages.NameForPackage(path)
		for _, service := range pkg.Services {
			for _, rpc := range service.RPCs {
				if rpc.Input.IsEmpty() {
					// we can have imported third party input message in case of empty input so we need to
					// add it's import otherwise there will not be present.
					file.imports = append(file.imports, rewrite.Import{
						Alias:      cfg.Packages.NameForPackage(string(rpc.Input.Name.Identifier.GoImportPath)),
						ImportPath: string(rpc.Input.Name.Identifier.GoImportPath),
					})
				}

				is[fmt.Sprintf("%s.%s", service.Name.Identifier.GoName, rpc.Name.Identifier.GoName)] = resolverImplementation{
					serviceResolverImplementation: &ServiceResolverImplementation{
						ServiceName:       service.Name.Identifier.GoName,
						ServiceMethodName: rpc.Name.Identifier.GoName,
						EmptyMessagePkg:   cfg.Packages.NameForPackage(string(rpc.Input.Name.Identifier.GoImportPath)),
						EmptyMessageName:  rpc.Input.Name.Identifier.GoName,
						IsEmptyOutput:     rpc.Output.IsEmpty(),
						IsEmptyInput:      rpc.Input.IsEmpty(),
					},
				}
			}
		}

		for _, message := range pkg.Messages {
			for _, field := range message.Fields {
				if field.Type.IsKv() {
					fieldName := changecase.Pascal(field.Name.GqlName)
					is[fmt.Sprintf("%s.%s", message.Name.Identifier.GoName, fieldName)] = resolverImplementation{
						fieldResolverImplementation: &FieldResolverImplementation{
							PkgName:          packageName,
							KeyValueTypeName: field.Type.Kv.Identifier.GoName,
							FieldName:        fieldName,
						},
					}
				}
			}
		}
	}

	return is
}
