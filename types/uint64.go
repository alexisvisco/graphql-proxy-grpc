package types

import (
	"encoding/json"
	"fmt"
	"github.com/99designs/gqlgen/graphql"
	"github.com/pkg/errors"
	"io"
	"strconv"
)

func MarshalUint64(i uint64) graphql.Marshaler {
	return graphql.WriterFunc(func(w io.Writer) {
		_, _ = io.WriteString(w, fmt.Sprintf(`"%d"`, i))
	})
}

func UnmarshalUint64(v interface{}) (uint64, error) {
	switch v := v.(type) {
	case string:
		i, err := strconv.ParseUint(v, 10, 64)
		if err != nil {
			return 0, errors.Wrap(err, "UnmarshalUint64: unable to unmarshal from string")
		} else {
			return i, nil
		}
	case int32:
		return uint64(v), nil
	case int64:
		return uint64(v), nil
	case uint32:
		return uint64(v), nil
	case uint64:
		return v, nil
	case json.Number:
		i, err := strconv.ParseUint(string(v), 10, 64)
		if err != nil {
			return 0, errors.Wrap(err, "UnmarshalUint64: unable to unmarshal from json.Number")
		} else {
			return i, nil
		}
	default:
		return 0, fmt.Errorf("UnmarshalUint64: %T is not a parseable", v)
	}
}
