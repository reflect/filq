package filter

import (
	"reflect"

	"github.com/pkg/errors"
	"github.com/reflect/filq/context"
)

type Expand struct {
	Filter Filter
}

func (e *Expand) Apply(ctx *context.Context, in context.Valuer) ([]context.Valuer, error) {
	vrs, err := e.Filter.Apply(ctx, in)
	if err != nil {
		return nil, err
	}

	var out []context.Valuer
	for _, vr := range vrs {
		v, err := vr.Value(ctx)
		if err != nil {
			return nil, err
		}

		ex, ok := v.(context.Iter)
		if !ok {
			return nil, errors.WithStack(&context.UnexpectedTypeError{
				Wanted: []reflect.Type{
					reflect.TypeOf((*context.Iter)(nil)).Elem(),
				},
				Got: reflect.TypeOf(v),
			})
		}

		vs, err := ex.Expand(ctx)
		if err != nil {
			return nil, err
		}

		out = append(out, vs...)
	}

	return out, nil
}
