
# Graphql proxy grpc

The graphql-proxy-grpc is a protoc plugin abd a CLI. It reads protobuf service definitions and generates a reverse-proxy 
server which translates a Graphql API into gRPC. The graphql server is generated using [gqlgen](https://gqlgen.com/).

There is two binaries, one protoc plugin that generate the graphql schema based on your proto specification and one 
binary that generate go resolvers based on the previous generated schema.


## Features

- [x] Enum
- [x] Messages
- [x] Map
- [x] Services (rpc method customization for query and mutation)
- [x] Multi packages
- [x] Grpc clients customizable
- [ ] Oneof (not planned)
- [ ] Documentation with comment (planned)
- [ ] Ovveride naming (planned)

## Usage/Examples

### Using buf

todo

### Using protoc

todo


## Related

Here are some related projects

[gqlgen](https://github.com/99designs/gqlgen)
[grpc gateway](https://github.com/grpc-ecosystem/grpc-gateway)

  