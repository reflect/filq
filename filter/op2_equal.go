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

	var equal bool
	if eq, ok := lv.(context.Eq); ok {
		equal, err = eq.Equal(ctx, f.r)
		if err != nil {
			return nil, err
		}
	} else {
		rv, err := f.r.Value(ctx)
		if err != nil {
			return nil, err
		}

		equal = rv == lv
	}

	if f.inverse {
		return !equal, nil
	}

	return equal, nil
}
