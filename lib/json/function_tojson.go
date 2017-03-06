package json

import (
	"encoding/json"

	"github.com/pkg/errors"
	"github.com/reflect/filq/context"
	"github.com/reflect/filq/types"
)

type toJSON struct {
	v interface{}
}

func (tj *toJSON) Value(ctx *context.Context) (interface{}, error) {
	b, err := json.Marshal(tj.v)
	if err != nil {
		return nil, errors.Wrap(err, "marshaling JSON")
	}

	return types.Bytes(b), nil
}

func ToJSON(ctx *context.Context, in context.Valuer) ([]context.Valuer, error) {
	v, err := in.Value(ctx)
	if err != nil {
		return nil, err
	}

	return []context.Valuer{&toJSON{v}}, nil
}
