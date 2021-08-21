package types

import (
	"encoding/json"
	"fmt"
	"github.com/99designs/gqlgen/graphql"
	"github.com/pkg/errors"
	"io"
	"strconv"
)

func MarshalFloat32(i float32) graphql.Marshaler {
	return graphql.WriterFunc(func(w io.Writer) {
		_, _ = io.WriteString(w, fmt.Sprintf(`%f`, i))
	})
}

func UnmarshalFloat32(v interface{}) (float32, error) {
	switch v := v.(type) {
	case string:
		i, err := strconv.ParseFloat(v, 32)
		if err != nil {
			return 0, errors.Wrap(err, "UnmarshalFloat32: unable to unmarshal from string")
		} else {
			return float32(i), nil
		}
	case int32:
		return float32(v), nil
	case int64:
		return float32(v), nil
	case uint32:
		return float32(v), nil
	case float32:
		return v, nil
	case json.Number:
		i, err := strconv.ParseFloat(string(v), 32)
		if err != nil {
			return 0, errors.Wrap(err, "UnmarshalFloat32: unable to unmarshal from json.Number")
		} else {
			return float32(i), nil
		}
	default:
		return 0, fmt.Errorf("UnmarshalFloat32: %T is not a parseable", v)
	}
}
