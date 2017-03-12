package function

import (
	"reflect"
	"strconv"

	"github.com/pkg/errors"
	"github.com/reflect/filq/context"
	"github.com/reflect/filq/types"
)

func init() {
	fn, _ := NewFunction(ToInt)
	register("toint", fn)
}

func ToInt(ctx *context.Context, in context.Valuer) ([]context.Valuer, error) {
	v, err := in.Value(ctx)
	if err != nil {
		return nil, err
	}

	var r int64
	if o, ok := v.(types.Int); ok {
		r = int64(o)
	} else if s, ok := v.(types.Str); ok {
		r, err = strconv.ParseInt(string(s), 10, 64)
		if err != nil {
			return nil, err
		}
	} else if f, ok := v.(types.Float); ok {
		r = int64(f)
	} else {
		return nil, errors.WithStack(&context.UnexpectedTypeError{
			Wanted: []reflect.Type{
				reflect.TypeOf(types.Int(0)),
				reflect.TypeOf(types.Float(0)),
				reflect.TypeOf(types.Str("")),
			},
			Got: reflect.TypeOf(v),
		})
	}

	return []context.Valuer{context.NewConstValuer(r)}, nil
}
