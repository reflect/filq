package types

import (
	"bytes"
	"fmt"
	"reflect"
	"unicode"

	"unicode/utf8"

	"github.com/pkg/errors"
	"github.com/reflect/filq/context"
)

type Bytes []byte

func (b Bytes) Format(f fmt.State, c rune) {
	if c != 'v' || !f.Flag('+') {
		FormatDefault(f, c, []byte(b))
	}

	out := []byte(b)
	for len(out) > 0 {
		r, size := utf8.DecodeRune(out)
		if r == utf8.RuneError {
			fmt.Fprintf(f, "\\x%x", out[0])
			out = out[1:]
		} else {
			if unicode.IsGraphic(r) {
				fmt.Fprintf(f, "%c", r)
			} else {
				fmt.Fprintf(f, "\\u%04x", r)
			}

			out = out[size:]
		}
	}
}

func (b Bytes) Index(ctx *context.Context, key context.Valuer) (context.Valuer, error) {
	k, err := key.Value(ctx)
	if err != nil {
		return nil, err
	}

	var r int

	switch kt := k.(type) {
	case Str:
		r = bytes.Index([]byte(b), []byte(kt))
	case Bytes:
		r = bytes.Index([]byte(b), kt)
	case byte:
		r = bytes.IndexByte([]byte(b), kt)
	case rune:
		r = bytes.IndexRune([]byte(b), kt)
	default:
		return nil, errors.WithStack(&context.UnexpectedTypeError{
			Wanted: []reflect.Type{
				reflect.TypeOf(Str("")),
				reflect.TypeOf(Bytes([]byte{})),
				reflect.TypeOf(byte(0)),
				reflect.TypeOf(rune(0)),
			},
			Got: reflect.TypeOf(k),
		})
	}

	if r < 0 {
		return context.NewConstValuer(nil), nil
	}

	return context.NewConstValuer(int64(r)), nil
}

type BytesConverter struct{}

func (bco *BytesConverter) Convert(in interface{}) interface{} {
	return Bytes(in.([]byte))
}
