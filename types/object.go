package types

import (
	"encoding/json"
	"fmt"
	"reflect"

	"github.com/pkg/errors"
	"github.com/reflect/filq/context"
)

type Object map[string]interface{}

func (o Object) Format(f fmt.State, c rune) {
	if c != 'v' || !f.Flag('+') {
		FormatDefault(f, c, map[string]interface{}(o))
	}

	b, err := json.MarshalIndent(o, "", "  ")
	if err != nil {
		FormatDefault(f, c, map[string]interface{}(o))
	}

	fmt.Fprintf(f, "%s", b)
}

func (o Object) Select(ctx *context.Context, tree []context.Valuer) (context.Valuer, error) {
	v, err := tree[0].Value(ctx)
	if err != nil {
		return nil, err
	}

	var key string

	switch vt := v.(type) {
	case Str:
		key = string(vt)
	case Bytes:
		key = string(vt)
	default:
		return nil, errors.WithStack(&context.UnexpectedTypeError{
			Wanted: []reflect.Type{
				reflect.TypeOf(Str("")),
				reflect.TypeOf(Bytes([]byte{})),
			},
			Got: reflect.TypeOf(v),
		})
	}

	sel, ok := o[key]
	if !ok {
		return context.NewConstValuer(nil), nil
	}

	sel = ctx.Convert(sel)

	if len(tree) == 1 {
		return context.NewConstValuer(sel), nil
	} else if selectable, ok := sel.(context.Sel); ok {
		return selectable.Select(ctx, tree[1:])
	}

	return nil, errors.WithStack(&context.UnexpectedTypeError{
		Wanted: []reflect.Type{
			reflect.TypeOf((*context.Sel)(nil)).Elem(),
		},
		Got: reflect.TypeOf(sel),
	})
}

func (o Object) Equal(ctx *context.Context, other context.Valuer) (bool, error) {
	ov, err := other.Value(ctx)
	if err != nil {
		return false, err
	}

	if ov == nil {
		return false, nil
	}

	oo, ok := ov.(Object)
	if !ok {
		return false, errors.WithStack(&context.UnexpectedTypeError{
			Wanted: []reflect.Type{
				reflect.TypeOf(Object(map[string]interface{}{})),
			},
			Got: reflect.TypeOf(ov),
		})
	}

	if len(oo) != len(o) {
		return false, nil
	}

	for key, value := range o {
		other, ok := oo[key]
		if !ok {
			return false, nil
		}

		if ve, ok := value.(context.Eq); ok {
			eq, err := ve.Equal(ctx, context.NewConstValuer(other))
			if err != nil {
				return false, err
			}

			if !eq {
				return false, nil
			}
		} else if value != other {
			return false, nil
		}
	}

	return true, nil
}

func (o Object) Expand(ctx *context.Context) ([]context.Valuer, error) {
	out := make([]context.Valuer, len(o))

	i := 0
	for key, value := range o {
		out[i] = NewEntryValuer(Str(key), ctx.Convert(value))
		i++
	}

	return out, nil
}

type ObjectConverter struct{}

func (oco *ObjectConverter) Convert(in interface{}) interface{} {
	return Object(in.(map[string]interface{}))
}
