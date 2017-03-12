package filter

import (
	"reflect"

	"github.com/pkg/errors"
	"github.com/reflect/filq/context"
)

type op2BoolFunc func(a, b bool) interface{}

func op2And(a, b bool) interface{} { return a && b }
func op2Or(a, b bool) interface{}  { return a || b }

type op2BoolFilter struct {
	fn   op2BoolFunc
	l, r context.Valuer
}

func (f *op2BoolFilter) Value(ctx *context.Context) (interface{}, error) {
	lv, err := f.l.Value(ctx)
	if err != nil {
		return nil, err
	}

	rv, err := f.r.Value(ctx)
	if err != nil {
		return nil, err
	}

	l, ok := lv.(bool)
	if !ok {
		return nil, errors.WithStack(&context.UnexpectedTypeError{
			Wanted: []reflect.Type{reflect.TypeOf(false)},
			Got:    reflect.TypeOf(lv),
		})
	}

	r, ok := rv.(bool)
	if !ok {
		return nil, &context.UnexpectedTypeError{
			Wanted: []reflect.Type{reflect.TypeOf(false)},
			Got:    reflect.TypeOf(rv),
		}
	}

	return f.fn(l, r), nil
}
