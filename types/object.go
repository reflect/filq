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
	} else if selectable, ok := sel.(context.Selectable); ok {
		return selectable.Select(ctx, tree[1:])
	}

	return nil, errors.WithStack(&context.UnexpectedTypeError{
		Wanted: []reflect.Type{
			reflect.TypeOf((*context.Selectable)(nil)).Elem(),
		},
		Got: reflect.TypeOf(sel),
	})
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
