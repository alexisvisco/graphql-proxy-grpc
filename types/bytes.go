package types

import (
	"encoding/json"
	"fmt"
	"github.com/99designs/gqlgen/graphql"
	"io"
)

func MarshalBytes(i []byte) graphql.Marshaler {
	return graphql.WriterFunc(func(w io.Writer) {
		marshal, _ := json.Marshal(i)
		_, _ = io.WriteString(w, string(marshal))
	})
}

func UnmarshalBytes(v interface{}) ([]byte, error) {
	switch v := v.(type) {
	case string:
		return []byte(v), nil
	case []byte:
		return v, nil
	default:
		return nil, fmt.Errorf("UnmarshalBytes: %T is not a parseable", v)
	}
}
