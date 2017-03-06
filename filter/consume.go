package filter

import (
	"github.com/reflect/filq/context"
)

type Consume struct {
	Filter Filter
}

func (c *Consume) Apply(ctx *context.Context, in context.Valuer) ([]context.Valuer, error) {
	// This deliberately consumes the input in case external dependencies
	// require it to be in a processed state.
	v, err := in.Value(ctx)
	if err != nil {
		return nil, err
	}

	return c.Filter.Apply(ctx, context.NewConstValuer(v))
}
