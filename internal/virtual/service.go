package virtual

import (
	"github.com/alexisvisco/graphql-proxy-grpc/internal/errorx"
	changecase "github.com/ku/go-change-case"
	"github.com/rs/zerolog/log"
	"go.buf.build/library/go/alexisvisco/graphql-proxy-grpc/graphqlpb/v1"
	"google.golang.org/protobuf/compiler/protogen"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/types/descriptorpb"
)

type Service struct {
	Name Name

	RPCs []*Rpc
}

func (s Service) String() string {
	return s.Name.String()
}

func (pkg *Package) createService(service *protogen.Service, descriptor *Descriptor) (*Service, error) {
	svc := &Service{
		Name: Name{
			Identifier: protogen.GoIdent{
				GoName:       service.GoName,
				GoImportPath: protogen.GoImportPath(pkg.Path),
			},
		},
		RPCs: make([]*Rpc, len(service.Methods)),
	}

	for i, rpc := range service.Methods {
		inputPkg := descriptor.getPackage(string(rpc.Input.GoIdent.GoImportPath))
		if inputPkg == nil {
			return nil, errorx.New().Err(errorx.ErrProtoPackageNotLoaded).
				Str("service", string(service.Desc.Name())).
				Str("service-pkg", string(service.Desc.ParentFile().Package())).
				Str("pkg-to-load", string(rpc.Input.Desc.ParentFile().Package())).
				Str("msg-to-load", string(rpc.Input.Desc.Name()))
		}

		msgInput, err := inputPkg.getOrCreateMessage(rpc.Input, descriptor)
		if msgInput == nil {
			return nil, err
		}

		msgInput.markAsInput()

		outputPkg := descriptor.getPackage(string(rpc.Output.GoIdent.GoImportPath))
		if err != nil {
			return nil, errorx.New().Err(errorx.ErrProtoPackageNotLoaded).
				Str("service", string(service.Desc.Name())).
				Str("service-pkg", string(service.Desc.ParentFile().Package())).
				Str("pkg-to-load", string(rpc.Output.Desc.ParentFile().Package())).
				Str("message-to-load", string(rpc.Output.Desc.Name()))
		}

		msgOutput, err := outputPkg.getOrCreateMessage(rpc.Output, descriptor)
		if err != nil {
			return nil, err
		}

		graphqlType := graphqlpbv1.GraphqlType_GRAPHQL_TYPE_QUERY
		rpcOption, ok := rpc.Desc.Options().(*descriptorpb.MethodOptions)
		if ok {
			t, ok := proto.GetExtension(rpcOption, graphqlpbv1.E_Type).(graphqlpbv1.GraphqlType)
			if ok {
				graphqlType = t
			}
		}
		svc.RPCs[i] = &Rpc{
			Name: Name{
				Identifier: protogen.GoIdent{
					GoName:       rpc.GoName,
					GoImportPath: svc.Name.Identifier.GoImportPath,
				},
				GqlName: changecase.Camel(service.GoName + rpc.GoName),
			},
			Input:  msgInput,
			Output: msgOutput,
			Type:   graphqlType,
		}
	}

	pkg.addService(svc)
	log.Debug().Msgf("register service %s in package %s", service.Desc.Name(), service.Desc.ParentFile().Package())

	return svc, nil
}

type Rpc struct {
	Name Name

	Input  *Message
	Output *Message

	Type graphqlpbv1.GraphqlType
}
