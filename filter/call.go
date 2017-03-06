package filter

import (
	"github.com/reflect/filq/context"
)

type Call struct {
	Function  string
	Arguments []Filter
}

func (c *Call) Apply(ctx *context.Context, in context.Valuer) ([]context.Valuer, error) {
	f, err := ctx.Function(c.Function, len(c.Arguments))
	if err != nil {
		return nil, err
	}

	arguments := make([][]context.Valuer, len(c.Arguments))
	for i, argument := range c.Arguments {
		arguments[i], err = argument.Apply(ctx, in)
		if err != nil {
			return nil, err
		}
	}

	return f.Call(ctx, in, arguments)
}
