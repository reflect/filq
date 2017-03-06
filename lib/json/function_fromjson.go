package json

import (
	"bytes"
	"encoding/json"
	"io"
	"reflect"
	"strings"

	"github.com/pkg/errors"
	"github.com/reflect/filq/context"
	"github.com/reflect/filq/types"
)

func from(v interface{}) (context.Valuer, error) {
	d, ok := v.(*json.Decoder)
	if !ok {
		if r, ok := v.(io.Reader); ok {
			d = json.NewDecoder(r)
		} else {
			switch vt := v.(type) {
			case types.Str:
				d = json.NewDecoder(strings.NewReader(string(vt)))
			case types.Bytes:
				d = json.NewDecoder(bytes.NewReader(vt))
			default:
				return nil, errors.WithStack(&context.UnexpectedTypeError{
					Wanted: []reflect.Type{
						reflect.TypeOf(&json.Decoder{}),
						reflect.TypeOf((*io.Reader)(nil)).Elem(),
						reflect.TypeOf(types.Str("")),
						reflect.TypeOf(types.Bytes([]byte{})),
					},
					Got: reflect.TypeOf(v),
				})
			}
		}
	}

	lazy := func(ctx *context.Context) (interface{}, error) {
		var out interface{}
		if err := d.Decode(&out); err != nil {
			return nil, errors.Wrap(err, "parsing JSON")
		}

		return ctx.Convert(out), nil
	}

	return context.NewLazyValuer(lazy), nil
}

func FromJSON(ctx *context.Context, in context.Valuer) ([]context.Valuer, error) {
	v, err := in.Value(ctx)
	if err != nil {
		return nil, err
	}

	vr, err := from(v)
	if err != nil {
		return nil, err
	}

	return []context.Valuer{vr}, nil
}
