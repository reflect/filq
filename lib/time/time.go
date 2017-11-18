package time

import (
	"github.com/reflect/filq/context"
	"github.com/reflect/filq/function"
)

func DefineIn(ctx *context.Context) {
	fn, _ := function.NewFunction(ParseTime)
	ctx.DefineFunction("parsetime", fn)
}
