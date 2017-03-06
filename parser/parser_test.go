package parser

import (
	"testing"

	"github.com/jmikkola/parsego/parser"
	"github.com/reflect/filq/context"
	"github.com/reflect/filq/filter"
	"github.com/reflect/filq/types"
	"github.com/stretchr/testify/assert"
)

func TestStringParser(t *testing.T) {
	r, err := parser.ParseString(stringParser(), `"test"`)
	assert.NoError(t, err)
	assert.Equal(t, types.Str("test"), r)

	r, err = parser.ParseString(stringParser(), `"tes\u0074\n"`)
	assert.NoError(t, err)
	assert.Equal(t, types.Str("test\n"), r)

	r, err = parser.ParseString(stringParser(), `"not terminated`)
	assert.EqualError(t, err, "expected a character, got error Reached end of input at line 0, col 15")
	assert.Nil(t, r)

	r, err = parser.ParseString(stringParser(), "\"\n\"")
	assert.EqualError(t, err, "expected a character in the range '\"' to '\"', got error \n at line 1, col 0")
	assert.Nil(t, r)
}

func TestIdentParser(t *testing.T) {
	r, err := parser.ParseString(identParser(), "abcd")
	assert.NoError(t, err)
	assert.Equal(t, "abcd", r)

	r, err = parser.ParseString(identParser(), "_0AbCd")
	assert.NoError(t, err)
	assert.Equal(t, "_0AbCd", r)

	r, err = parser.ParseString(identParser(), "$test")
	assert.EqualError(t, err, "no parser matched at line 0, col 0")
	assert.Nil(t, r)

	r, err = parser.ParseString(identParser(), "01234")
	assert.EqualError(t, err, "no parser matched at line 0, col 0")
	assert.Nil(t, r)
}

func TestFuncParser(t *testing.T) {
	r, err := parser.ParseString(funcParser(), "fn")
	assert.NoError(t, err)
	assert.Equal(t, &filter.Call{
		Function: "fn",
	}, r)

	r, err = parser.ParseString(funcParser(), "fn()")
	assert.NoError(t, err)
	assert.Equal(t, &filter.Call{
		Function:  "fn",
		Arguments: []filter.Filter{},
	}, r)

	r, err = parser.ParseString(funcParser(), `fn("a"; 10)`)
	assert.NoError(t, err)
	assert.Equal(t, &filter.Call{
		Function: "fn",
		Arguments: []filter.Filter{
			&filter.Pipe{Filter: str("a")},
			&filter.Pipe{Filter: &filter.Const{context.NewConstValuer(int64(10))}},
		},
	}, r)
}

func TestSelectorParser(t *testing.T) {
	r, err := parser.ParseString(selectorParser(false), ".")
	assert.NoError(t, err)
	assert.Equal(t, &filter.Selector{
		Recall: &context.PipeRecall{},
		Tree:   []filter.Filter{},
	}, r)

	r, err = parser.ParseString(selectorParser(false), ".foo.bar")
	assert.NoError(t, err)
	assert.Equal(t, &filter.Selector{
		Recall: &context.PipeRecall{},
		Tree:   []filter.Filter{str("foo"), str("bar")},
	}, r)

	r, err = parser.ParseString(selectorParser(false), `."foo"."$bar$"`)
	assert.NoError(t, err)
	assert.Equal(t, &filter.Selector{
		Recall: &context.PipeRecall{},
		Tree:   []filter.Filter{str("foo"), str("$bar$")},
	}, r)

	r, err = parser.ParseString(selectorParser(true), "$bucket")
	assert.NoError(t, err)
	assert.Equal(t, &filter.Selector{
		Recall: &context.VariableRecall{Name: "bucket"},
		Tree:   []filter.Filter{},
	}, r)

	r, err = parser.ParseString(selectorParser(false), "$bucket")
	assert.EqualError(t, err, "expected a character in the range '.' to '.', got error $ at line 0, col 1")
	assert.Nil(t, r)

	r, err = parser.ParseString(selectorParser(true), "$bucket.foo.bar")
	assert.NoError(t, err)
	assert.Equal(t, &filter.Selector{
		Recall: &context.VariableRecall{Name: "bucket"},
		Tree:   []filter.Filter{str("foo"), str("bar")},
	}, r)
}

func TestExpandParser(t *testing.T) {
	r, err := parser.ParseString(expandParser(), ".")
	assert.NoError(t, err)
	assert.Equal(t, &filter.Selector{
		Recall: &context.PipeRecall{},
		Tree:   []filter.Filter{},
	}, r)

	r, err = parser.ParseString(expandParser(), ".[]")
	assert.NoError(t, err)
	assert.Equal(t, &filter.Expand{
		Filter: &filter.Selector{
			Recall: &context.PipeRecall{},
			Tree:   []filter.Filter{},
		},
	}, r)

	r, err = parser.ParseString(expandParser(), `$x[]."foo"[][].bar`)
	assert.NoError(t, err)
	assert.Equal(t, &filter.Pipe{
		Filter: &filter.Expand{
			Filter: &filter.Selector{
				Recall: &context.VariableRecall{Name: "x"},
				Tree:   []filter.Filter{},
			},
		},
		Next: &filter.Pipe{
			Filter: &filter.Expand{
				Filter: &filter.Expand{
					Filter: &filter.Selector{
						Recall: &context.PipeRecall{},
						Tree:   []filter.Filter{str("foo")},
					},
				},
			},
			Next: &filter.Pipe{
				Filter: &filter.Selector{
					Recall: &context.PipeRecall{},
					Tree:   []filter.Filter{str("bar")},
				},
			},
		},
	}, r)
}

func TestExprParser(t *testing.T) {
	r, err := parser.ParseString(exprParser(), ".foo.bar")
	assert.NoError(t, err)
	assert.Equal(t, &filter.Selector{
		Recall: &context.PipeRecall{},
		Tree:   []filter.Filter{str("foo"), str("bar")},
	}, r)

	r, err = parser.ParseString(exprParser(), `."foo"."$bar$"`)
	assert.NoError(t, err)
	assert.Equal(t, &filter.Selector{
		Recall: &context.PipeRecall{},
		Tree:   []filter.Filter{str("foo"), str("$bar$")},
	}, r)

	r, err = parser.ParseString(exprParser(), `($x + 1) * 40`)
	assert.NoError(t, err)
	assert.Equal(t, &filter.Op2{
		Operator: "*",
		Left: &filter.Scope{
			Filter: &filter.Pipe{
				Filter: &filter.Op2{
					Operator: "+",
					Left:     &filter.Selector{Recall: &context.VariableRecall{Name: "x"}, Tree: []filter.Filter{}},
					Right:    &filter.Const{context.NewConstValuer(int64(1))},
				},
			},
		},
		Right: &filter.Const{context.NewConstValuer(int64(40))},
	}, r)

	r, err = parser.ParseString(exprParser(), "(.a | .b)")
	assert.NoError(t, err)
	assert.Equal(t, &filter.Scope{
		Filter: &filter.Pipe{
			Filter: &filter.Selector{
				Recall: &context.PipeRecall{},
				Tree:   []filter.Filter{str("a")},
			},
			Next: &filter.Pipe{
				Filter: &filter.Selector{
					Recall: &context.PipeRecall{},
					Tree:   []filter.Filter{str("b")},
				},
			},
		},
	}, r)
}

func TestPipelineParser(t *testing.T) {
	r, err := parser.ParseString(pipelineParser(), ".foo.bar")
	assert.NoError(t, err)
	assert.Equal(t, &filter.Pipe{
		Filter: &filter.Selector{
			Recall: &context.PipeRecall{},
			Tree:   []filter.Filter{str("foo"), str("bar")},
		},
	}, r)

	r, err = parser.ParseString(pipelineParser(), ".a | .b")
	assert.NoError(t, err)
	assert.Equal(t, &filter.Pipe{
		Filter: &filter.Selector{
			Recall: &context.PipeRecall{},
			Tree:   []filter.Filter{str("a")},
		},
		Next: &filter.Pipe{
			Filter: &filter.Selector{
				Recall: &context.PipeRecall{},
				Tree:   []filter.Filter{str("b")},
			},
		},
	}, r)

	r, err = parser.ParseString(pipelineParser(), ". as $db | buckets | $db[.]")
	assert.NoError(t, err)
	assert.Equal(t, &filter.Pipe{
		Filter: &filter.Selector{
			Recall: &context.PipeRecall{},
			Tree:   []filter.Filter{},
		},
		Assignment: &context.SimpleAssignment{Name: "db"},
		Next: &filter.Pipe{
			Filter: &filter.Call{
				Function: "buckets",
			},
			Next: &filter.Pipe{
				Filter: &filter.Selector{
					Recall: &context.VariableRecall{Name: "db"},
					Tree: []filter.Filter{
						&filter.Selector{Recall: &context.PipeRecall{}, Tree: []filter.Filter{}},
					},
				},
			},
		},
	}, r)

	r, err = parser.ParseString(pipelineParser(), ".a as $a | .b")
	assert.NoError(t, err)
	assert.Equal(t, &filter.Pipe{
		Filter: &filter.Selector{
			Recall: &context.PipeRecall{},
			Tree:   []filter.Filter{str("a")},
		},
		Assignment: &context.SimpleAssignment{Name: "a"},
		Next: &filter.Pipe{
			Filter: &filter.Selector{
				Recall: &context.PipeRecall{},
				Tree:   []filter.Filter{str("b")},
			},
		},
	}, r)
}

func str(in string) *filter.Const {
	return &filter.Const{context.NewConstValuer(types.Str(in))}
}
