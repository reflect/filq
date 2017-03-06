package context

import (
	"reflect"

	"github.com/pkg/errors"
)

type Context struct {
	variables  map[string]Valuer
	functions  map[string]map[int]Function
	converters map[reflect.Type]Converter
	next       *Context
}

func (c *Context) DefineVariable(name string, v Valuer) {
	c.variables[name] = v
}

func (c *Context) Variable(name string) (Valuer, error) {
	if v, ok := c.variables[name]; ok {
		return v, nil
	}

	if c.next != nil {
		return c.next.Variable(name)
	}

	return nil, &VariableNotDefinedError{Name: name}
}

func (c *Context) DefineFunction(name string, f Function) {
	m, ok := c.functions[name]
	if !ok {
		m = make(map[int]Function)
		c.functions[name] = m
	}

	m[f.Arity()] = f
}

func (c *Context) Function(name string, arity int) (Function, error) {
	if fns, ok := c.functions[name]; ok {
		if fn, ok := fns[arity]; ok {
			return fn, nil
		}
	}

	if c.next != nil {
		return c.next.Function(name, arity)
	}

	return nil, errors.WithStack(&FunctionNotDefinedError{Name: name})
}

func (c *Context) Convert(in interface{}) interface{} {
	t := reflect.TypeOf(in)
	if co, ok := c.converters[t]; ok {
		return co.Convert(in)
	}

	if c.next != nil {
		return c.next.Convert(in)
	}

	return in
}

func (c *Context) DefineConverter(t reflect.Type, co Converter) {
	c.converters[t] = co
}

func OverlayContext(ctx *Context) *Context {
	return &Context{
		variables:  make(map[string]Valuer),
		functions:  make(map[string]map[int]Function),
		converters: make(map[reflect.Type]Converter),
		next:       ctx,
	}
}
