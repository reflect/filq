package filter

import (
	"reflect"

	"github.com/pkg/errors"
	"github.com/reflect/filq/context"
	"github.com/reflect/filq/types"
)

type ObjectEntry struct {
	Key, Value Filter
}

type Object struct {
	Entries []ObjectEntry
}

func (o *Object) Apply(ctx *context.Context, in context.Valuer) ([]context.Valuer, error) {
	type kv struct{ Key, Value context.Valuer }
	var kvs [][]kv

	// Add the root object.
	kvs = append(kvs, []kv{})

	for _, entry := range o.Entries {
		keys, err := entry.Key.Apply(ctx, in)
		if err != nil {
			return nil, err
		}

		var next []kv

		for _, key := range keys {
			vf := entry.Value
			if vf == nil {
				vf = &Selector{
					Recall: &context.PipeRecall{},
					Tree:   []Filter{&Const{Valuer: key}},
				}
			}

			values, err := vf.Apply(ctx, in)
			if err != nil {
				return nil, err
			}

			for _, value := range values {
				next = append(next, kv{Key: key, Value: value})
			}
		}

		nkvs := make([][]kv, len(kvs)*len(next))
		for i, prev := range kvs {
			for j, kvr := range next {
				key := i*len(next) + j
				nkvs[key] = append(nkvs[key], prev...)
				nkvs[key] = append(nkvs[key], kvr)
			}
		}

		kvs = nkvs
	}

	vrs := make([]context.Valuer, len(kvs))
	for i := range kvs {
		// Bind in scope.
		mi := kvs[i]

		lazy := func(ctx *context.Context) (interface{}, error) {
			m := make(map[string]interface{})

			for _, entry := range mi {
				key, err := entry.Key.Value(ctx)
				if err != nil {
					return nil, err
				}

				value, err := entry.Value.Value(ctx)
				if err != nil {
					return nil, err
				}

				var ks string
				if s, ok := key.(types.Str); ok {
					ks = string(s)
				} else if b, ok := key.(types.Bytes); ok {
					ks = string(b)
				} else {
					return nil, errors.WithStack(&context.UnexpectedTypeError{
						Wanted: []reflect.Type{
							reflect.TypeOf(types.Str("")),
							reflect.TypeOf(types.Bytes([]byte{})),
						},
						Got: reflect.TypeOf(key),
					})
				}

				m[ks] = value
			}

			return ctx.Convert(m), nil
		}

		vrs[i] = context.NewLazyValuer(lazy)
	}

	return vrs, nil
}
