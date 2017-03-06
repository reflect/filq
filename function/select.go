package function

import (
	"reflect"

	"github.com/pkg/errors"
	"github.com/reflect/filq/context"
)

func init() {
	fn, _ := NewFunction(Select)
	register("select", fn)
}

func Select(ctx *context.Context, in context.Valuer, criteria []context.Valuer) ([]context.Valuer, error) {
	for _, criterion := range criteria {
		v, err := criterion.Value(ctx)
		if err != nil {
			return nil, err
		}

		t, ok := v.(bool)
		if !ok {
			return nil, errors.WithStack(&context.UnexpectedTypeError{
				Wanted: []reflect.Type{reflect.TypeOf(false)},
				Got:    reflect.TypeOf(v),
			})
		}

		if !t {
			return []context.Valuer{}, nil
		}
	}

	return []context.Valuer{in}, nil
}
