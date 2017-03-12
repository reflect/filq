package filq

import (
	"github.com/reflect/filq/context"
	"github.com/reflect/filq/function"
	"github.com/reflect/filq/lib/io"
	"github.com/reflect/filq/lib/json"
	"github.com/reflect/filq/lib/regex"
	"github.com/reflect/filq/parser"
	"github.com/reflect/filq/types"
)

func NewParser() *parser.Parser {
	return parser.NewParser()
}

func NewContext() *context.Context {
	def := context.OverlayContext(nil)
	function.DefineIn(def)

	// Standard types.
	types.DefineIn(def)

	// Standard library.
	io.DefineIn(def)
	json.DefineIn(def)
	regex.DefineIn(def)

	return def
}

func NewConstValuer(v interface{}) context.Valuer {
	return context.NewConstValuer(v)
}

func Run(ctx *context.Context, statements string, v interface{}) ([]interface{}, error) {
	filter, err := NewParser().ParseString(statements)
	if err != nil {
		return nil, err
	}

	in, ok := v.(context.Valuer)
	if !ok {
		in = NewConstValuer(v)
	}

	vrs, err := filter.Apply(ctx, in)
	if err != nil {
		return nil, err
	}

	outs := make([]interface{}, len(vrs))
	for i, vr := range vrs {
		out, err := vr.Value(ctx)
		if err != nil {
			return nil, err
		}

		outs[i] = out
	}

	return outs, nil
}
