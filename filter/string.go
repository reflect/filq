package filter

import (
	"bytes"
	"fmt"

	"github.com/reflect/filq/context"
	"github.com/reflect/filq/types"
)

type String struct {
	Filters []Filter
}

func (s *String) Apply(ctx *context.Context, in context.Valuer) ([]context.Valuer, error) {
	// Add root.
	vrs := [][]context.Valuer{[]context.Valuer{}}

	// Cartesian product of each filter.
	for _, filter := range s.Filters {
		fvr, err := filter.Apply(ctx, in)
		if err != nil {
			return nil, err
		}

		switch len(fvr) {
		case 0:
		case 1:
			for i := range vrs {
				vrs[i] = append(vrs[i], fvr[0])
			}
		default:
			nvrs := make([][]context.Valuer, len(vrs)*len(fvr))
			for i, prev := range vrs {
				for j, vr := range fvr {
					key := i*len(fvr) + j
					nvrs[key] = append(nvrs[key], prev...)
					nvrs[key] = append(nvrs[key], vr)
				}
			}

			vrs = nvrs
		}
	}

	lazy := make([]context.Valuer, len(vrs))
	for i, vr := range vrs {
		fn := func(ctx *context.Context) (interface{}, error) {
			var b bytes.Buffer

			for _, vri := range vr {
				v, err := vri.Value(ctx)
				if err != nil {
					return nil, err
				}

				if v == nil {
					b.WriteString("null")
					continue
				}

				switch vt := v.(type) {
				case types.Str:
					b.WriteString(string(vt))
				case types.Bytes:
					b.Write([]byte(vt))
				default:
					b.WriteString(fmt.Sprintf("%v", vt))
				}
			}

			return types.Str(b.Bytes()), nil
		}

		lazy[i] = context.NewLazyValuer(fn)
	}

	return lazy, nil
}
