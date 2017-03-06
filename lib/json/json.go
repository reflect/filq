package json

import (
	"github.com/reflect/filq/context"
	"github.com/reflect/filq/function"
)

func DefineIn(ctx *context.Context) {
	fn, _ := function.NewFunction(FromJSON)
	ctx.DefineFunction("fromjson", fn)

	fn, _ = function.NewFunction(ToJSON)
	ctx.DefineFunction("tojson", fn)
}

func NewValuer(in interface{}) (context.Valuer, error) {
	return from(in)
}
