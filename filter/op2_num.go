package filter

import (
	"reflect"

	"github.com/pkg/errors"
	"github.com/reflect/filq/context"
	"github.com/reflect/filq/types"
)

type op2NumFunc struct {
	fnI func(a, b types.Int) interface{}
	fnF func(a, b types.Float) interface{}
}

func (f op2NumFunc) ApplyInt(a, b types.Int) interface{} {
	return f.fnI(a, b)
}

func (f op2NumFunc) ApplyFloat(a, b types.Float) interface{} {
	return f.fnF(a, b)
}

var (
	op2Mul = op2NumFunc{
		fnI: func(a, b types.Int) interface{} { return a * b },
		fnF: func(a, b types.Float) interface{} { return a * b },
	}
	op2Div = op2NumFunc{
		fnI: func(a, b types.Int) interface{} { return a / b },
		fnF: func(a, b types.Float) interface{} { return a / b },
	}
	op2Sub = op2NumFunc{
		fnI: func(a, b types.Int) interface{} { return a - b },
		fnF: func(a, b types.Float) interface{} { return a - b },
	}
)

type op2NumFilter struct {
	fn   op2NumFunc
	l, r context.Valuer
}

func (f *op2NumFilter) Value(ctx *context.Context) (interface{}, error) {
	lv, err := f.l.Value(ctx)
	if err != nil {
		return nil, err
	}

	rv, err := f.r.Value(ctx)
	if err != nil {
		return nil, err
	}

	var out interface{}

	if li, ok := lv.(types.Int); ok {
		ri, ok := rv.(types.Int)
		if !ok {
			return nil, errors.WithStack(&context.UnexpectedTypeError{
				Wanted: []reflect.Type{reflect.TypeOf(types.Int(0))},
				Got:    reflect.TypeOf(rv),
			})
		}

		out = f.fn.fnI(li, ri)
	} else if lf, ok := lv.(types.Float); ok {
		rf, ok := rv.(types.Float)
		if !ok {
			return nil, errors.WithStack(&context.UnexpectedTypeError{
				Wanted: []reflect.Type{reflect.TypeOf(types.Float(0))},
				Got:    reflect.TypeOf(rv),
			})
		}

		out = f.fn.fnF(lf, rf)
	} else {
		return nil, errors.WithStack(&context.UnexpectedTypeError{
			Wanted: []reflect.Type{reflect.TypeOf(types.Int(0)), reflect.TypeOf(types.Float(0))},
			Got:    reflect.TypeOf(lv),
		})
	}

	return out, nil
}
