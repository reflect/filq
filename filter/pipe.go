package filter

import (
	"github.com/reflect/filq/context"
)

type Pipe struct {
	Filter     Filter
	Assignment context.Assignment
	Next       *Pipe
}

func (p *Pipe) Apply(ctx *context.Context, in context.Valuer) ([]context.Valuer, error) {
	m, err := p.Filter.Apply(ctx, in)
	if err != nil {
		return nil, err
	}

	var outs []context.Valuer
	if p.Assignment != nil {
		for _, mi := range m {
			cctx := context.OverlayContext(ctx)
			if err := p.Assignment.AssignIn(cctx, mi); err != nil {
				return nil, err
			}

			f, err := p.Next.Apply(cctx, in)
			if err != nil {
				return nil, err
			}

			outs = append(outs, f...)
		}
	} else if p.Next != nil {
		for _, mi := range m {
			f, err := p.Next.Apply(ctx, mi)
			if err != nil {
				return nil, err
			}

			outs = append(outs, f...)
		}
	} else {
		outs = m
	}

	return outs, nil
}
