package types

import (
	"encoding/json"
	"fmt"
	"reflect"

	"github.com/pkg/errors"
	"github.com/reflect/filq/context"
)

type Array []interface{}

func (a Array) Format(f fmt.State, c rune) {
	if c != 'v' || !f.Flag('+') {
		FormatDefault(f, c, []interface{}(a))
	}

	b, err := json.MarshalIndent(a, "", "  ")
	if err != nil {
		FormatDefault(f, c, []interface{}(a))
	}

	fmt.Fprintf(f, "%s", b)
}

func (a Array) selectIndex(ctx *context.Context, idx int, tree []context.Valuer) (context.Valuer, error) {
	if idx < 0 || idx >= len(a) {
		return context.NewConstValuer(nil), nil
	}

	sel := ctx.Convert(a[idx])

	if len(tree) == 1 {
		return context.NewConstValuer(sel), nil
	} else if selectable, ok := sel.(context.Selectable); ok {
		return selectable.Select(ctx, tree[1:])
	}

	return nil, &context.UnexpectedTypeError{
		Wanted: []reflect.Type{
			reflect.TypeOf((*context.Selectable)(nil)).Elem(),
		},
		Got: reflect.TypeOf(sel),
	}
}

func (a Array) selectSlice(ctx *context.Context, slice Slice, tree []context.Valuer) (context.Valuer, error) {
	var out []interface{}

	min := int(slice.Left)
	if min < 0 {
		min = 0
	}

	max := int(slice.Right)
	if max >= len(a) {
		max = len(a) - 1
	}

	for i := min; i <= max; i++ {
		out = append(out, a[i])
	}

	if len(tree) == 1 {
		return context.NewConstValuer(Array(out)), nil
	}

	return Array(out).Select(ctx, tree[1:])
}

func (a Array) Select(ctx *context.Context, tree []context.Valuer) (context.Valuer, error) {
	v, err := tree[0].Value(ctx)
	if err != nil {
		return nil, err
	}

	switch vt := v.(type) {
	case int64:
		return a.selectIndex(ctx, int(vt), tree)
	case float64:
		return a.selectIndex(ctx, int(vt), tree)
	case Slice:
		return a.selectSlice(ctx, vt, tree)
	default:
		return nil, errors.WithStack(&context.UnexpectedTypeError{
			Wanted: []reflect.Type{
				reflect.TypeOf(int64(0)),
				reflect.TypeOf(float64(0)),
				reflect.TypeOf(Slice{}),
			},
			Got: reflect.TypeOf(v),
		})
	}
}

func (a Array) Expand(ctx *context.Context) ([]context.Valuer, error) {
	out := make([]context.Valuer, len(a))

	for i, value := range a {
		out[i] = context.NewConstValuer(ctx.Convert(value))
	}

	return out, nil
}

type ArrayConverter struct{}

func (aco *ArrayConverter) Convert(in interface{}) interface{} {
	return Array(in.([]interface{}))
}
