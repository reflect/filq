package types

import (
	"fmt"
	"reflect"

	"github.com/pkg/errors"
	"github.com/reflect/filq/context"
)

type entrySelector struct {
	e *Entry
	v context.Valuer
}

func (s *entrySelector) Value(ctx *context.Context) (interface{}, error) {
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
	case "key":
		return ctx.Convert(s.e.Key), nil
	case "value":
		return ctx.Convert(s.e.Value), nil
	default:
		return nil, nil
	}
}

type Entry struct {
	Key, Value interface{}
}

func (e *Entry) Format(f fmt.State, c rune) {
	fmt.Fprintf(f, "%+v = %+v", e.Key, e.Value)
}

func (e *Entry) Select(ctx *context.Context, tree []context.Valuer) (context.Valuer, error) {
	if len(tree) != 1 {
		return context.NewConstValuer(nil), nil
	}

	return &entrySelector{e: e, v: tree[0]}, nil
}

func NewEntryValuer(key, value interface{}) context.Valuer {
	return context.NewConstValuer(&Entry{Key: key, Value: value})
}
