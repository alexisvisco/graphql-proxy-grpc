package types

import (
	"encoding/json"
	"fmt"
	"github.com/99designs/gqlgen/graphql"
	"github.com/pkg/errors"
	"io"
	"strconv"
)

func MarshalUint32(i uint32) graphql.Marshaler {
	return graphql.WriterFunc(func(w io.Writer) {
		_, _ = io.WriteString(w, fmt.Sprintf("%d", i))
	})
}

func UnmarshalUint32(v interface{}) (uint32, error) {
	switch v := v.(type) {
	case string:
		i, err := strconv.ParseUint(v, 10, 32)
		if err != nil {
			return 0, errors.Wrap(err, "UnmarshalUint32: unable to unmarshal from string")
		} else {
			return uint32(i), nil
		}
	case int32:
		return uint32(v), nil
	case int64:
		return uint32(v), nil
	case uint32:
		return v, nil
	case uint64:
		return uint32(v), nil
	case json.Number:
		i, err := strconv.ParseUint(string(v), 10, 32)
		if err != nil {
			return 0, errors.Wrap(err, "UnmarshalUint32: unable to unmarshal from json.Number")
		} else {
			return uint32(i), nil
		}
	default:
		return 0, fmt.Errorf("UnmarshalUint32: %T is not a parseable", v)
	}
}
