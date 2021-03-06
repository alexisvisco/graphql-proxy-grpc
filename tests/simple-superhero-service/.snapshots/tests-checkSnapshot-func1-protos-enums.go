// Code generated by graphql-proxy-grpc-protoc. DO NOT EDIT.

package superheroes

import (
	fmt "fmt"
	io "io"
	strconv "strconv"
)

// UnmarshalGQL for CallSuperHeroRequest_City
func (x *CallSuperHeroRequest_City) UnmarshalGQL(v interface{}) error {
	str, ok := v.(string)
	if !ok {
		return fmt.Errorf("enums must be strings")
	}

	val, ok := CallSuperHeroRequest_City_value[str]
	if !ok {
		return fmt.Errorf("%s is not a valid CallSuperHeroRequest_City", str)
	}

	*x = CallSuperHeroRequest_City(val)
	return nil
}

// MarshalGQL for CallSuperHeroRequest_City
func (x CallSuperHeroRequest_City) MarshalGQL(w io.Writer) {
	fmt.Fprint(w, strconv.Quote(x.String()))
}

// UnmarshalGQL for Heroes
func (x *Heroes) UnmarshalGQL(v interface{}) error {
	str, ok := v.(string)
	if !ok {
		return fmt.Errorf("enums must be strings")
	}

	val, ok := Heroes_value[str]
	if !ok {
		return fmt.Errorf("%s is not a valid Heroes", str)
	}

	*x = Heroes(val)
	return nil
}

// MarshalGQL for Heroes
func (x Heroes) MarshalGQL(w io.Writer) {
	fmt.Fprint(w, strconv.Quote(x.String()))
}

