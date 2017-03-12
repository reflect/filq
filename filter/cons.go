package filter

import (
	"github.com/reflect/filq/context"
	"github.com/reflect/filq/types"
)

type Cons struct {
	Filters []Filter
}

func (c *Cons) Apply(ctx *context.Context, in context.Valuer) ([]context.Valuer, error) {
	var vrs []context.Valuer
	for _, filter := range c.Filters {
		outs, err := filter.Apply(ctx, in)
		if err != nil {
			return nil, err
		}

		vrs = append(vrs, outs...)
	}

	lazy := func(ctx *context.Context) (interface{}, error) {
		elements := make([]interface{}, len(vrs))
		for i, vr := range vrs {
			v, err := vr.Value(ctx)
			if err != nil {
				return nil, err
			}

			elements[i] = v
		}

		return types.Array(elements), nil
	}

	return []context.Valuer{context.NewLazyValuer(lazy)}, nil
}
