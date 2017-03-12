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
	} else if selectable, ok := sel.(context.Sel); ok {
		return selectable.Select(ctx, tree[1:])
	}

	return nil, &context.UnexpectedTypeError{
		Wanted: []reflect.Type{
			reflect.TypeOf((*context.Sel)(nil)).Elem(),
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

func (a Array) Equal(ctx *context.Context, other context.Valuer) (bool, error) {
	ov, err := other.Value(ctx)
	if err != nil {
		return false, err
	}

	oa, ok := ov.(Array)
	if !ok {
		return false, nil
	}

	if len(oa) != len(a) {
		return false, nil
	}

	for i, test := range a {
		if ve, ok := test.(context.Eq); ok {
			eq, err := ve.Equal(ctx, context.NewConstValuer(oa[i]))
			if err != nil {
				return false, err
			}

			if !eq {
				return false, nil
			}
		} else if test != oa[i] {
			return false, nil
		}
	}

	return true, nil
}

func (a Array) Index(ctx *context.Context, key context.Valuer) (context.Valuer, error) {
	v, err := key.Value(ctx)
	if err != nil {
		return nil, err
	}

	for i, candidate := range a {
		if ve, ok := v.(context.Eq); ok {
			eq, err := ve.Equal(ctx, context.NewConstValuer(candidate))
			if err != nil {
				return nil, err
			}

			if eq {
				return context.NewConstValuer(i), nil
			}
		} else if v == ctx.Convert(candidate) {
			return context.NewConstValuer(i), nil
		}
	}

	return context.NewConstValuer(nil), nil
}

func (a Array) Select(ctx *context.Context, tree []context.Valuer) (context.Valuer, error) {
	v, err := tree[0].Value(ctx)
	if err != nil {
		return nil, err
	}

	switch vt := v.(type) {
	case Int:
		return a.selectIndex(ctx, int(vt), tree)
	case Float:
		return a.selectIndex(ctx, int(vt), tree)
	case Slice:
		return a.selectSlice(ctx, vt, tree)
	default:
		return nil, errors.WithStack(&context.UnexpectedTypeError{
			Wanted: []reflect.Type{
				reflect.TypeOf(Int(0)),
				reflect.TypeOf(Float(0)),
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
