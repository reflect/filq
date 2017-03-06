package context

type Function interface {
	Arity() int
	Call(ctx *Context, in Valuer, arguments [][]Valuer) ([]Valuer, error)
}
