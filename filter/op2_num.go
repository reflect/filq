package filter

import (
	"reflect"

	"github.com/pkg/errors"
	"github.com/reflect/filq/context"
)

type op2NumFunc struct {
	fnI func(a, b int64) interface{}
	fnF func(a, b float64) interface{}
}

func (f op2NumFunc) ApplyInt(a, b int64) interface{} {
	return f.fnI(a, b)
}

func (f op2NumFunc) ApplyFloat(a, b float64) interface{} {
	return f.fnF(a, b)
}

type op2BoolFunc func(a, b bool) interface{}

var (
	op2Mul = op2NumFunc{
		fnI: func(a, b int64) interface{} { return a * b },
		fnF: func(a, b float64) interface{} { return a * b },
	}
	op2Div = op2NumFunc{
		fnI: func(a, b int64) interface{} { return a / b },
		fnF: func(a, b float64) interface{} { return a / b },
	}
	op2Sub = op2NumFunc{
		fnI: func(a, b int64) interface{} { return a - b },
		fnF: func(a, b float64) interface{} { return a - b },
	}
	op2Lt = op2NumFunc{
		fnI: func(a, b int64) interface{} { return a < b },
		fnF: func(a, b float64) interface{} { return a < b },
	}
	op2Lte = op2NumFunc{
		fnI: func(a, b int64) interface{} { return a <= b },
		fnF: func(a, b float64) interface{} { return a <= b },
	}
	op2Gt = op2NumFunc{
		fnI: func(a, b int64) interface{} { return a > b },
		fnF: func(a, b float64) interface{} { return a > b },
	}
	op2Gte = op2NumFunc{
		fnI: func(a, b int64) interface{} { return a >= b },
		fnF: func(a, b float64) interface{} { return a >= b },
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

	if li, ok := lv.(int64); ok {
		ri, ok := rv.(int64)
		if !ok {
			return nil, errors.WithStack(&context.UnexpectedTypeError{
				Wanted: []reflect.Type{reflect.TypeOf(int64(0))},
				Got:    reflect.TypeOf(rv),
			})
		}

		out = f.fn.fnI(li, ri)
	} else if lf, ok := lv.(float64); ok {
		rf, ok := rv.(float64)
		if !ok {
			return nil, errors.WithStack(&context.UnexpectedTypeError{
				Wanted: []reflect.Type{reflect.TypeOf(float64(0))},
				Got:    reflect.TypeOf(rv),
			})
		}

		out = f.fn.fnF(lf, rf)
	} else {
		return nil, errors.WithStack(&context.UnexpectedTypeError{
			Wanted: []reflect.Type{reflect.TypeOf(int64(0)), reflect.TypeOf(float64(0))},
			Got:    reflect.TypeOf(lv),
		})
	}

	return out, nil
}
