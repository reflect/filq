package context

type Valuer interface {
	Value(ctx *Context) (interface{}, error)
}

type ConstValuer struct {
	value interface{}
}

func (c *ConstValuer) Value(ctx *Context) (interface{}, error) {
	return ctx.Convert(c.value), nil
}

func NewConstValuer(v interface{}) *ConstValuer {
	return &ConstValuer{v}
}

type LazyValuer struct {
	fn    func(ctx *Context) (interface{}, error)
	value interface{}
	err   error
}

func (l *LazyValuer) Value(ctx *Context) (interface{}, error) {
	if l.value == nil && l.err == nil {
		l.value, l.err = l.fn(ctx)
	}

	return ctx.Convert(l.value), l.err
}

func NewLazyValuer(fn func(ctx *Context) (interface{}, error)) *LazyValuer {
	return &LazyValuer{fn: fn}
}
