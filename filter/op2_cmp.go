package filter

import (
	"reflect"

	"github.com/pkg/errors"
	"github.com/reflect/filq/context"
)

type op2CmpFunc func(cmp int) bool

func op2Gt(cmp int) bool  { return cmp > 0 }
func op2Gte(cmp int) bool { return cmp >= 0 }
func op2Lt(cmp int) bool  { return cmp < 0 }
func op2Lte(cmp int) bool { return cmp <= 0 }

type op2CmpFilter struct {
	fn   op2CmpFunc
	l, r context.Valuer
}

func (f *op2CmpFilter) Value(ctx *context.Context) (interface{}, error) {
	lv, err := f.l.Value(ctx)
	if err != nil {
		return nil, err
	}

	lc, ok := lv.(context.Cmp)
	if !ok {
		return nil, errors.WithStack(&context.UnexpectedTypeError{
			Wanted: []reflect.Type{
				reflect.TypeOf((*context.Cmp)(nil)).Elem(),
			},
			Got: reflect.TypeOf(lv),
		})
	}

	cmp, err := lc.Compare(ctx, f.r)
	if err != nil {
		return nil, err
	}

	return f.fn(cmp), nil
}
