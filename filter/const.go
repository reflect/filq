package filter

import (
	"github.com/reflect/filq/context"
)

type Const struct {
	Valuer context.Valuer
}

func (c *Const) Apply(ctx *context.Context, in context.Valuer) ([]context.Valuer, error) {
	return []context.Valuer{c.Valuer}, nil
}
