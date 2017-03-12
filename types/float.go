package types

import (
	"reflect"

	"github.com/pkg/errors"
	"github.com/reflect/filq/context"
)

type Float float64

func (f Float) Equal(ctx *context.Context, other context.Valuer) (bool, error) {
	ov, err := other.Value(ctx)
	if err != nil {
		return false, err
	}

	if of, ok := ov.(Float); ok {
		return float64(f) == float64(of), nil
	}

	return false, nil
}

func (f Float) Compare(ctx *context.Context, other context.Valuer) (int, error) {
	ov, err := other.Value(ctx)
	if err != nil {
		return 0, err
	}

	of, ok := ov.(Float)
	if !ok {
		return 0, errors.WithStack(&context.UnexpectedTypeError{
			Wanted: []reflect.Type{reflect.TypeOf(Float(0))},
			Got:    reflect.TypeOf(ov),
		})
	}

	if f > of {
		return 1, nil
	} else if f < of {
		return -1, nil
	}

	return 0, nil
}

type FloatFloat64Converter struct{}

func (ico *FloatFloat64Converter) Convert(in interface{}) interface{} {
	return Float(in.(float64))
}

type FloatFloat32Converter struct{}

func (ico *FloatFloat32Converter) Convert(in interface{}) interface{} {
	return Float(in.(float32))
}
