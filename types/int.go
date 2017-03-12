package types

import (
	"reflect"

	"github.com/pkg/errors"
	"github.com/reflect/filq/context"
)

type Int int64

func (i Int) Equal(ctx *context.Context, other context.Valuer) (bool, error) {
	ov, err := other.Value(ctx)
	if err != nil {
		return false, err
	}

	if oi, ok := ov.(Int); ok {
		return int64(i) == int64(oi), nil
	}

	return false, nil
}

func (i Int) Compare(ctx *context.Context, other context.Valuer) (int, error) {
	ov, err := other.Value(ctx)
	if err != nil {
		return 0, err
	}

	oi, ok := ov.(Int)
	if !ok {
		return 0, errors.WithStack(&context.UnexpectedTypeError{
			Wanted: []reflect.Type{reflect.TypeOf(Int(0))},
			Got:    reflect.TypeOf(ov),
		})
	}

	if i > oi {
		return 1, nil
	} else if i < oi {
		return -1, nil
	}

	return 0, nil
}

type IntInt64Converter struct{}

func (ico *IntInt64Converter) Convert(in interface{}) interface{} {
	return Int(in.(int64))
}

type IntInt32Converter struct{}

func (ico *IntInt32Converter) Convert(in interface{}) interface{} {
	return Int(in.(int32))
}

type IntInt16Converter struct{}

func (ico *IntInt16Converter) Convert(in interface{}) interface{} {
	return Int(in.(int16))
}

type IntInt8Converter struct{}

func (ico *IntInt8Converter) Convert(in interface{}) interface{} {
	return Int(in.(int8))
}

type IntIntConverter struct{}

func (ico *IntIntConverter) Convert(in interface{}) interface{} {
	return Int(in.(int))
}
