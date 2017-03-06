package context

type Selectable interface {
	Select(ctx *Context, tree []Valuer) (Valuer, error)
}

type Indexable interface {
	Index(ctx *Context, key Valuer) (Valuer, error)
}

type Expandable interface {
	Expand(ctx *Context) ([]Valuer, error)
}
