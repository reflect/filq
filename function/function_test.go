package function

import (
	"testing"

	"github.com/reflect/filq/context"
	"github.com/stretchr/testify/assert"
)

func TestNewFunction(t *testing.T) {
	fn1 := func(ctx *context.Context, in context.Valuer) ([]context.Valuer, error) {
		return nil, nil
	}

	fn, err := NewFunction(fn1)
	assert.NoError(t, err)
	assert.NotNil(t, fn)
	assert.Equal(t, 0, fn.Arity())

	fn2 := func(ctx *context.Context, in context.Valuer, arg1 []context.Valuer) ([]context.Valuer, error) {
		return nil, nil
	}

	fn, err = NewFunction(fn2)
	assert.NoError(t, err)
	assert.NotNil(t, fn)
	assert.Equal(t, 1, fn.Arity())

	fn3 := func() ([]context.Valuer, error) {
		return nil, nil
	}

	fn, err = NewFunction(fn3)
	assert.Equal(t, ErrInvalidFunction, err)
	assert.Nil(t, fn)

	fn4 := func(ctx *context.Context, in context.Valuer) error {
		return nil
	}

	fn, err = NewFunction(fn4)
	assert.Equal(t, ErrInvalidFunction, err)
	assert.Nil(t, fn)

	fn5 := func(ctx *context.Context, in string) error {
		return nil
	}

	fn, err = NewFunction(fn5)
	assert.Equal(t, ErrInvalidFunction, err)
	assert.Nil(t, fn)

	fn6 := func(ctx *context.Context, in context.Valuer, arg1 string) error {
		return nil
	}

	fn, err = NewFunction(fn6)
	assert.Equal(t, ErrInvalidFunction, err)
	assert.Nil(t, fn)
}

func TestFunctionCall(t *testing.T) {
	fn1 := func(ctx *context.Context, in context.Valuer, arg1 []context.Valuer) ([]context.Valuer, error) {
		v1, _ := in.Value(ctx)
		v2, _ := arg1[0].Value(ctx)

		return []context.Valuer{context.NewConstValuer(v1.(int) + v2.(int))}, nil
	}

	fn, err := NewFunction(fn1)
	assert.NoError(t, err)
	assert.NotNil(t, fn)

	vs, err := fn.Call(context.OverlayContext(nil), context.NewConstValuer(1), [][]context.Valuer{{context.NewConstValuer(2)}})
	assert.NoError(t, err)
	assert.Equal(t, []context.Valuer{context.NewConstValuer(3)}, vs)
}
