package time

import (
	"reflect"

	"github.com/pkg/errors"
	"github.com/reflect/filq/context"
	"github.com/reflect/filq/types"
	"github.com/reflect/xparse/xtime"
)

type parseTime struct {
	in string
}

func (pt *parseTime) Value(ctx *context.Context) (interface{}, error) {
	t, err := xtime.Parse(pt.in)
	if err != nil {
		return nil, err
	}

	return ctx.Convert(t), nil
}

func ParseTime(ctx *context.Context, in context.Valuer) ([]context.Valuer, error) {
	v, err := in.Value(ctx)
	if err != nil {
		return nil, err
	}

	var t string
	if s, ok := v.(types.Str); ok {
		t = string(s)
	} else if b, ok := v.(types.Bytes); ok {
		t = string(b)
	} else {
		return nil, errors.WithStack(&context.UnexpectedTypeError{
			Wanted: []reflect.Type{
				reflect.TypeOf(types.Str("")),
				reflect.TypeOf(types.Bytes([]byte{})),
			},
			Got: reflect.TypeOf(v),
		})
	}

	return []context.Valuer{&parseTime{t}}, nil
}
