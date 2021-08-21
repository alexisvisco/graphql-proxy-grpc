package templatedata

import (
	"github.com/alexisvisco/graphql-proxy-grpc/internal/virtual"
	"strings"
)

func messageIntoInputType(msg virtual.Message) virtual.Message {
	msg.Name.GqlName = nameToInput(msg.Name.GqlName)
	return msg
}

func nameToInput(str string) string {
	return strings.ReplaceAll(str, "Request", "") + "Input"
}
