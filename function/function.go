package function

import (
	"reflect"

	"github.com/pkg/errors"
	"github.com/reflect/filq/context"
)

var (
	builtins = make(map[string]context.Function)
)

func register(name string, fn context.Function) {
	builtins[name] = fn
}

func DefineIn(ctx *context.Context) {
	for name, fn := range builtins {
		ctx.DefineFunction(name, fn)
	}
}

type Function struct {
	arity  int
	callee reflect.Value
}

func (f *Function) Arity() int {
	return f.arity
}

func (f *Function) Call(ctx *context.Context, in context.Valuer, arguments [][]context.Valuer) ([]context.Valuer, error) {
	ins := make([]reflect.Value, len(arguments)+2)
	ins[0] = reflect.ValueOf(ctx)
	ins[1] = reflect.ValueOf(in)

	for i, argument := range arguments {
		ins[i+2] = reflect.ValueOf(argument)
	}

	outs := f.callee.Call(ins)
	if ev := outs[1].Interface(); ev != nil {
		return nil, errors.Wrap(ev.(error), "calling function")
	}

	vv := outs[0].Interface()
	if vv == nil {
		return nil, nil
	}

	return vv.([]context.Valuer), nil
}

func NewFunction(fn interface{}) (context.Function, error) {
	callee := reflect.ValueOf(fn)
	t := callee.Type()

	if t.Kind() != reflect.Func {
		return nil, ErrInvalidFunction
	}

	// Not supported for now. Will reconsider in the future.
	if t.IsVariadic() {
		return nil, ErrInvalidFunction
	}

	// Make sure the first two arguments are *context.Context and
	// context.Valuer.
	if t.NumIn() < 2 {
		return nil, ErrInvalidFunction
	}

	if !reflect.TypeOf(&context.Context{}).AssignableTo(t.In(0)) {
		return nil, ErrInvalidFunction
	}

	if !reflect.TypeOf((*context.Valuer)(nil)).Elem().AssignableTo(t.In(1)) {
		return nil, ErrInvalidFunction
	}

	// The remaining arguments should accept []context.Valuer.
	for i := 2; i < t.NumIn(); i++ {
		if !reflect.TypeOf([]context.Valuer{}).AssignableTo(t.In(i)) {
			return nil, ErrInvalidFunction
		}
	}

	// The function should return ([]context.Valuer, error).
	if t.NumOut() != 2 {
		return nil, ErrInvalidFunction
	}

	if !t.Out(0).AssignableTo(reflect.TypeOf([]context.Valuer{})) {
		return nil, ErrInvalidFunction
	}

	if !t.Out(1).AssignableTo(reflect.TypeOf((*error)(nil)).Elem()) {
		return nil, ErrInvalidFunction
	}

	f := &Function{
		arity:  t.NumIn() - 2,
		callee: callee,
	}

	return f, nil
}
