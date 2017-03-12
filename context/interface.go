package context

// Eq is the interface for equality. It is expected to be transitive.
type Eq interface {
	Equal(ctx *Context, other Valuer) (bool, error)
}

// Cmp is the interface for comparability.
type Cmp interface {
	Compare(ctx *Context, other Valuer) (int, error)
}

// Sel is the interface for types that allow their contents to be picked by a
// tree.
type Sel interface {
	Select(ctx *Context, tree []Valuer) (Valuer, error)
}

// Idx is the interface for types that allow searching for a value within them.
type Idx interface {
	Index(ctx *Context, key Valuer) (Valuer, error)
}

// Iter is the interface for types that can be expanded into a set of more
// values.
type Iter interface {
	Expand(ctx *Context) ([]Valuer, error)
}
