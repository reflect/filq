package filter

import (
	"fmt"
	"reflect"

	"github.com/pkg/errors"
	"github.com/reflect/filq/context"
)

type op1BoolFunc func(a bool) interface{}

func op1Not(a bool) interface{} { return !a }

type op1BoolFilter struct {
	fn op1BoolFunc
	v  context.Valuer
}

func (f *op1BoolFilter) Value(ctx *context.Context) (interface{}, error) {
	v, err := f.v.Value(ctx)
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

	return f.fn(t), nil
}

type Op1 struct {
	Operator string
	Operand  Filter
}

func (o *Op1) Apply(ctx *context.Context, in context.Valuer) ([]context.Valuer, error) {
	vs, err := o.Operand.Apply(ctx, in)
	if err != nil {
		return nil, err
	}

	out := make([]context.Valuer, len(vs))
	for i, v := range vs {
		switch o.Operator {
		case "not":
			out[i] = &op1BoolFilter{fn: op1Not, v: v}
		default:
			panic(fmt.Errorf("unary operator %q not implemented", o.Operator))
		}
	}

	return out, nil
}
