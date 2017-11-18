package io

import (
	"io"
	"os"
	"reflect"

	"github.com/pkg/errors"
	"github.com/reflect/filq/context"
	"github.com/reflect/filq/types"
)

func writeFileWithMode(ctx *context.Context, in context.Valuer, paths []context.Valuer, flag int) ([]context.Valuer, error) {
	var d types.Bytes

	v, err := in.Value(ctx)
	if err != nil {
		return nil, err
	}

	switch vt := v.(type) {
	case types.Str:
		d = types.Bytes(vt)
	case types.Bytes:
		d = vt
	default:
		return nil, errors.WithStack(&context.UnexpectedTypeError{
			Wanted: []reflect.Type{
				reflect.TypeOf(types.Str("")),
				reflect.TypeOf(types.Bytes([]byte{})),
			},
			Got: reflect.TypeOf(v),
		})
	}

	for _, path := range paths {
		pv, err := path.Value(ctx)
		if err != nil {
			return nil, err
		}

		var p types.Str

		switch pvt := pv.(type) {
		case types.Str:
			p = pvt
		case types.Bytes:
			p = types.Str(pvt)
		default:
			return nil, errors.WithStack(&context.UnexpectedTypeError{
				Wanted: []reflect.Type{
					reflect.TypeOf(types.Str("")),
					reflect.TypeOf(types.Bytes([]byte{})),
				},
				Got: reflect.TypeOf(pv),
			})
		}

		f, err := os.OpenFile(string(p), os.O_WRONLY|os.O_CREATE|flag, 0644)
		if err != nil {
			return nil, err
		}
		defer f.Close()

		n, err := f.Write(d)
		if err != nil {
			return nil, err
		} else if n < len(d) {
			return nil, io.ErrShortWrite
		}
	}

	return []context.Valuer{in}, nil
}

func WriteFile(ctx *context.Context, in context.Valuer, paths []context.Valuer) ([]context.Valuer, error) {
	return writeFileWithMode(ctx, in, paths, os.O_TRUNC)
}
