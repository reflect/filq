package function

import (
	"testing"

	"github.com/reflect/filq/context"
	"github.com/reflect/filq/types"
	"github.com/stretchr/testify/assert"
)

func TestIndex(t *testing.T) {
	ctx := context.OverlayContext(nil)

	conds := []struct {
		Str, Search string
		Expected    interface{}
	}{
		{"this longcat is long", "long", int64(5)},
		{"this longcat is long", "short", nil},
		{"zero is the loneliest number", "zero", int64(0)},
	}

	for _, cond := range conds {
		l, err := Index(
			ctx,
			context.NewConstValuer(types.Str(cond.Str)),
			[]context.Valuer{context.NewConstValuer(types.Str(cond.Search))},
		)
		assert.NoError(t, err)
		assert.Equal(t, []context.Valuer{context.NewConstValuer(cond.Expected)}, l)
	}
}
