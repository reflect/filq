package types

import (
	"reflect"
	"time"

	"github.com/pkg/errors"
	"github.com/reflect/filq/context"
)

type Time struct {
	time.Time
}

type timeSelector struct {
	t Time
	v context.Valuer
}

func (s *timeSelector) Value(ctx *context.Context) (interface{}, error) {
	v, err := s.v.Value(ctx)
	if err != nil {
		return nil, err
	}

	var sub string
	if b, ok := v.(Bytes); ok {
		sub = string(b)
	} else if s, ok := v.(Str); ok {
		sub = string(s)
	} else {
		return nil, errors.WithStack(&context.UnexpectedTypeError{
			Wanted: []reflect.Type{
				reflect.TypeOf(Str("")),
				reflect.TypeOf(Bytes([]byte{})),
			},
			Got: reflect.TypeOf(v),
		})
	}

	switch sub {
	case "year":
		return ctx.Convert(s.t.Year()), nil
	case "month":
		return ctx.Convert(int(s.t.Month())), nil
	case "day":
		return ctx.Convert(s.t.Day()), nil
	case "hour":
		return ctx.Convert(s.t.Hour()), nil
	case "minute":
		return ctx.Convert(s.t.Minute()), nil
	case "second":
		return ctx.Convert(s.t.Second()), nil
	case "nanosecond":
		return ctx.Convert(s.t.Nanosecond()), nil
	default:
		return nil, nil
	}
}

func (t Time) Equal(ctx *context.Context, other context.Valuer) (bool, error) {
	ov, err := other.Value(ctx)
	if err != nil {
		return false, err
	}

	if ot, ok := ov.(Time); ok {
		return t.Time == ot.Time, nil
	}

	return false, nil
}

func (t Time) Select(ctx *context.Context, tree []context.Valuer) (context.Valuer, error) {
	if len(tree) != 1 {
		return context.NewConstValuer(nil), nil
	}

	return &timeSelector{t: t, v: tree[0]}, nil
}

type TimeConverter struct{}

func (tco *TimeConverter) Convert(in interface{}) interface{} {
	return Time{in.(time.Time)}
}
