package filter

import (
	"github.com/reflect/filq/context"
)

type op2EqualFilter struct {
	l, r    context.Valuer
	inverse bool
}

func (f *op2EqualFilter) Value(ctx *context.Context) (interface{}, error) {
	lv, err := f.l.Value(ctx)
	if err != nil {
		return nil, err
	}

	rv, err := f.r.Value(ctx)
	if err != nil {
		return nil, err
	}

	if f.inverse {
		return lv != rv, nil
	}

	return lv == rv, nil
}
