package errorx

import (
	"github.com/pkg/errors"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

var (
	ErrProtoPackageNotLoaded           = errors.New("package wasn't loaded before it's usage")
	ErrProtoMapKeyTypeInvalid          = errors.New("key of a proto can be every scalar except bytes and floating point types")
	ErrProtoMapValueTypeInvalid        = errors.New("unable to understand value type in the map")
	ErrEmptyMessageNotSupportedInField = errors.New("empty field message not supported")
	ErrInvalidGraphqlSchemaTemplate    = errors.New("unable to parse template for generating graphql schema")
	ErrUnableToCreateGraphqlSchema     = errors.New("unable to create graphql schema")
	ErrUnableToUseProtoGeneratedFile   = errors.New("unable to use proto to write graphql schema")
)

type Context struct {
	ctx zerolog.Context
	err error
}

func New() *Context {
	return &Context{
		ctx: log.With(),
	}
}

func (c *Context) Str(key string, value string) *Context {
	c.ctx = c.ctx.Str(key, value)
	return c
}

func (c *Context) Err(err error) *Context {
	c.ctx = c.ctx.Stack().Err(err)
	c.err = err
	return c
}

func (c *Context) Int(s string, i int) *Context {
	c.ctx = c.ctx.Int(s, i)
	return c
}

func (c *Context) Panic(msg string) {
	l := c.ctx.Logger()
	l.Panic().Msg(msg)
}

func (c *Context) Error() string {
	return c.err.Error()
}
