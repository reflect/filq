package filter

import (
	"reflect"

	"github.com/pkg/errors"
	"github.com/reflect/filq/context"
)

type Selector struct {
	Recall context.Recall
	Tree   []Filter
}

func (s *Selector) tree(ctx *context.Context, in context.Valuer) ([]context.Valuer, error) {
	tree := make([]context.Valuer, len(s.Tree))
	for i, f := range s.Tree {
		vs, err := f.Apply(ctx, in)
		if err != nil {
			return nil, err
		}

		if len(vs) != 1 {
			return nil, errors.WithStack(&CannotSubscriptError{Dimensions: len(vs)})
		}

		tree[i] = vs[0]
	}

	return tree, nil
}

func (s *Selector) Apply(ctx *context.Context, in context.Valuer) ([]context.Valuer, error) {
	r, err := s.Recall.Resolve(ctx, in)
	if err != nil {
		return nil, err
	}

	var out context.Valuer

	if len(s.Tree) == 0 {
		out = r
	} else {
		v, err := r.Value(ctx)
		if err != nil {
			return nil, err
		}

		tree, err := s.tree(ctx, in)
		if err != nil {
			return nil, err
		}

		selector, ok := v.(context.Selectable)
		if !ok {
			return nil, errors.WithStack(&context.UnexpectedTypeError{
				Wanted: []reflect.Type{
					reflect.TypeOf((*context.Selectable)(nil)).Elem(),
				},
				Got: reflect.TypeOf(v),
			})
		}

		out, err = selector.Select(ctx, tree)
		if err != nil {
			return nil, err
		}
	}

	return []context.Valuer{out}, nil
}
