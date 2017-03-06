package filter

import (
	"github.com/reflect/filq/context"
)

type Filter interface {
	Apply(ctx *context.Context, in context.Valuer) ([]context.Valuer, error)
}
