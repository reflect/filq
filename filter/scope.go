package filter

import (
	"github.com/reflect/filq/context"
)

type Scope struct {
	Filter Filter
}

func (s *Scope) Apply(ctx *context.Context, in context.Valuer) ([]context.Valuer, error) {
	cctx := context.OverlayContext(ctx)
	return s.Filter.Apply(cctx, in)
}
