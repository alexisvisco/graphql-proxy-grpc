package main

import (
	_ "embed"
	"flag"
	"fmt"
	"github.com/alexisvisco/graphql-proxy-grpc/internal/errorx"
	"github.com/alexisvisco/graphql-proxy-grpc/internal/gengo"
	"github.com/alexisvisco/graphql-proxy-grpc/internal/gengql"
	"github.com/alexisvisco/graphql-proxy-grpc/internal/virtual"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"os"
	"path/filepath"
	"time"

	"google.golang.org/protobuf/compiler/protogen"
)

var (
	flags flag.FlagSet
)

// main binary is responsible for generating the graphql schema.
// It will also produce a temporary file named 'descriptor.json'.
// This file (descriptor.json) serve to serialize proto information.
func main() {
	if len(os.Args) == 2 && os.Args[1] == "--version" {
		_, _ = fmt.Fprintf(os.Stdout, "%v %v\n", filepath.Base(os.Args[0]), "1.0")
		os.Exit(0)
	}
	if len(os.Args) == 2 && os.Args[1] == "--help" {
		_, _ = fmt.Fprintf(os.Stdout, "todo \n")
		os.Exit(0)
	}

	proto := virtual.NewDescriptor()
	plugin := protogen.Options{ParamFunc: flags.Set}
	plugin.Run(runPlugin(proto))
}

func runPlugin(proto *virtual.Descriptor) func(gen *protogen.Plugin) error {
	return func(gen *protogen.Plugin) error {
		if err := proto.InspectProto(gen); err != nil {
			err.(*errorx.Context).Panic("unable to inspect protobuf schema")
			return err
		}

		_, err := gengql.GenerateSchema(*proto, gen)
		if err != nil {
			err.(*errorx.Context).Panic("unable to generate graphql schema")
			return err
		}

		gengo.Generate(proto, gen)

		if err := proto.GenerateDescriptorFileForGQLGen(gen); err != nil {
			return err
		}

		return nil
	}
}

func init() {
	log.Logger = log.Output(zerolog.ConsoleWriter{
		Out:        os.Stderr,
		TimeFormat: time.RFC3339,
	}).With().Timestamp().Str("generator", "grpc2graphql").Logger()
}
