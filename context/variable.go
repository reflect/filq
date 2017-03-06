package context

type Assignment interface {
	AssignIn(ctx *Context, in Valuer) error
}

type SimpleAssignment struct {
	Name string
}

func (a *SimpleAssignment) AssignIn(ctx *Context, in Valuer) error {
	ctx.DefineVariable(a.Name, in)
	return nil
}

type Recall interface {
	Resolve(ctx *Context, in Valuer) (Valuer, error)
}

type PipeRecall struct{}

func (r *PipeRecall) Resolve(ctx *Context, in Valuer) (Valuer, error) {
	return in, nil
}

type VariableRecall struct {
	Name string
}

func (r *VariableRecall) Resolve(ctx *Context, in Valuer) (Valuer, error) {
	return ctx.Variable(r.Name)
}
