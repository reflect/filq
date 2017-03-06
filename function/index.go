package function

import (
	"reflect"

	"github.com/pkg/errors"
	"github.com/reflect/filq/context"
)

func init() {
	fn, _ := NewFunction(Index)
	register("index", fn)
}

func Index(ctx *context.Context, in context.Valuer, search []context.Valuer) ([]context.Valuer, error) {
	out := make([]context.Valuer, len(search))

	v, err := in.Value(ctx)
	if err != nil {
		return nil, err
	}

	vs, ok := v.(context.Indexable)
	if !ok {
		return nil, errors.WithStack(&context.UnexpectedTypeError{
			Wanted: []reflect.Type{
				reflect.TypeOf((*context.Indexable)(nil)).Elem(),
			},
			Got: reflect.TypeOf(v),
		})
	}

	for i, candidate := range search {
		r, err := vs.Index(ctx, candidate)
		if err != nil {
			return nil, err
		}

		out[i] = r
	}

	return out, nil
}
