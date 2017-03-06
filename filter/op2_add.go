package filter

import (
	"fmt"
	"reflect"

	"github.com/pkg/errors"
	"github.com/reflect/filq/context"
	"github.com/reflect/filq/types"
)

type op2AddFilter struct {
	l, r context.Valuer
}

func (f *op2AddFilter) Value(ctx *context.Context) (interface{}, error) {
	lv, err := f.l.Value(ctx)
	if err != nil {
		return nil, err
	}

	rv, err := f.r.Value(ctx)
	if err != nil {
		return nil, err
	}

	var out interface{}

	switch lt := lv.(type) {
	case int64:
		rt, ok := rv.(int64)
		if !ok {
			return nil, errors.WithStack(&context.UnexpectedTypeError{
				Wanted: []reflect.Type{reflect.TypeOf(int64(0))},
				Got:    reflect.TypeOf(rv),
			})
		}

		out = lt + rt
	case float64:
		rt, ok := rv.(float64)
		if !ok {
			return nil, errors.WithStack(&context.UnexpectedTypeError{
				Wanted: []reflect.Type{reflect.TypeOf(float64(0))},
				Got:    reflect.TypeOf(rv),
			})
		}

		out = lt + rt
	case types.Str:
		switch rt := rv.(type) {
		case types.Str:
			out = fmt.Sprintf("%s%s", lt, rt)
		case types.Bytes:
			out = fmt.Sprintf("%s%s", lt, rt)
		default:
			return nil, errors.WithStack(&context.UnexpectedTypeError{
				Wanted: []reflect.Type{
					reflect.TypeOf(types.Str("")),
					reflect.TypeOf(types.Bytes([]byte{})),
				},
				Got: reflect.TypeOf(rv),
			})
		}
	case types.Bytes:
		switch rt := rv.(type) {
		case types.Str:
			rtb := []byte(rt)
			b := make(types.Bytes, 0, len(lt)+len(rtb))
			b = append(b, lt...)
			b = append(b, rtb...)

			out = b
		case types.Bytes:
			b := make(types.Bytes, 0, len(lt)+len(rt))
			b = append(b, lt...)
			b = append(b, rt...)

			out = b
		default:
			return nil, errors.WithStack(&context.UnexpectedTypeError{
				Wanted: []reflect.Type{
					reflect.TypeOf(types.Str("")),
					reflect.TypeOf(types.Bytes([]byte{})),
				},
				Got: reflect.TypeOf(rv),
			})
		}
	default:
		return nil, errors.WithStack(&context.UnexpectedTypeError{
			Wanted: []reflect.Type{
				reflect.TypeOf(int64(0)),
				reflect.TypeOf(float64(0)),
				reflect.TypeOf(types.Str("")),
				reflect.TypeOf(types.Bytes([]byte{})),
			},
			Got: reflect.TypeOf(lv),
		})
	}

	return out, nil
}
