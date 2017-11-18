package parser

import (
	"testing"

	"github.com/reflect/filq/context"
	"github.com/reflect/filq/filter"
	"github.com/reflect/parsego/parser"
	"github.com/stretchr/testify/assert"
)

func TestStringParser(t *testing.T) {
	r, err := parser.ParseString(stringParser(), `"test"`)
	assert.NoError(t, err)
	assert.Equal(t, &filter.String{
		Filters: []filter.Filter{constFilter("test")},
	}, r)

	r, err = parser.ParseString(stringParser(), `"tes\u0074\n"`)
	assert.NoError(t, err)
	assert.Equal(t, &filter.String{
		Filters: []filter.Filter{constFilter("tes"), constFilter("t"), constFilter("\n")},
	}, r)

	r, err = parser.ParseString(stringParser(), `"test \(.interpolation)!"`)
	assert.NoError(t, err)
	assert.Equal(t, &filter.String{
		Filters: []filter.Filter{
			constFilter("test "),
			&filter.Scope{
				&filter.Pipe{
					Filter: &filter.Selector{
						Recall: &context.PipeRecall{},
						Tree:   []filter.Filter{constFilter("interpolation")},
					},
				},
			},
			constFilter("!"),
		},
	}, r)

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
			&filter.Pipe{Filter: stringFilter("a")},
			&filter.Pipe{Filter: constFilter(int64(10))},
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
		Tree:   []filter.Filter{constFilter("foo"), constFilter("bar")},
	}, r)

	r, err = parser.ParseString(selectorParser(false), `."foo"."$bar$"`)
	assert.NoError(t, err)
	assert.Equal(t, &filter.Selector{
		Recall: &context.PipeRecall{},
		Tree:   []filter.Filter{stringFilter("foo"), stringFilter("$bar$")},
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
		Tree:   []filter.Filter{constFilter("foo"), constFilter("bar")},
	}, r)
}

func TestArrayParser(t *testing.T) {
	r, err := parser.ParseString(arrayParser(), "[]")
	assert.NoError(t, err)
	assert.Equal(t, &filter.Cons{Filters: []filter.Filter{}}, r)

	r, err = parser.ParseString(arrayParser(), "[1, 2, 3]")
	assert.NoError(t, err)
	assert.Equal(t, &filter.Cons{Filters: []filter.Filter{
		&filter.Pipe{Filter: constFilter(int64(1))},
		&filter.Pipe{Filter: constFilter(int64(2))},
		&filter.Pipe{Filter: constFilter(int64(3))},
	}}, r)

	r, err = parser.ParseString(arrayParser(), "[.[] | .test, []]")
	assert.NoError(t, err)
	assert.Equal(t, &filter.Cons{Filters: []filter.Filter{
		&filter.Pipe{
			Filter: &filter.Expand{
				Filter: &filter.Selector{Recall: &context.PipeRecall{}, Tree: []filter.Filter{}},
			},
			Next: &filter.Pipe{
				Filter: &filter.Selector{
					Recall: &context.PipeRecall{},
					Tree:   []filter.Filter{constFilter("test")},
				},
			},
		},
		&filter.Pipe{
			Filter: &filter.Cons{Filters: []filter.Filter{}},
		},
	}}, r)
}

func TestObjectParser(t *testing.T) {
	r, err := parser.ParseString(objectParser(), "{}")
	assert.NoError(t, err)
	assert.Equal(t, &filter.Object{Entries: []filter.ObjectEntry{}}, r)
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

	r, err = parser.ParseString(expandParser(), "[][]")
	assert.NoError(t, err)
	assert.Equal(t, &filter.Expand{
		Filter: &filter.Cons{
			Filters: []filter.Filter{},
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
						Tree:   []filter.Filter{stringFilter("foo")},
					},
				},
			},
			Next: &filter.Pipe{
				Filter: &filter.Selector{
					Recall: &context.PipeRecall{},
					Tree:   []filter.Filter{constFilter("bar")},
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
		Tree:   []filter.Filter{constFilter("foo"), constFilter("bar")},
	}, r)

	r, err = parser.ParseString(exprParser(), `."foo"."$bar$"`)
	assert.NoError(t, err)
	assert.Equal(t, &filter.Selector{
		Recall: &context.PipeRecall{},
		Tree:   []filter.Filter{stringFilter("foo"), stringFilter("$bar$")},
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
					Right:    constFilter(int64(1)),
				},
			},
		},
		Right: constFilter(int64(40)),
	}, r)

	r, err = parser.ParseString(exprParser(), "(.a | .b)")
	assert.NoError(t, err)
	assert.Equal(t, &filter.Scope{
		Filter: &filter.Pipe{
			Filter: &filter.Selector{
				Recall: &context.PipeRecall{},
				Tree:   []filter.Filter{constFilter("a")},
			},
			Next: &filter.Pipe{
				Filter: &filter.Selector{
					Recall: &context.PipeRecall{},
					Tree:   []filter.Filter{constFilter("b")},
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
			Tree:   []filter.Filter{constFilter("foo"), constFilter("bar")},
		},
	}, r)

	r, err = parser.ParseString(pipelineParser(), ".a | .b")
	assert.NoError(t, err)
	assert.Equal(t, &filter.Pipe{
		Filter: &filter.Selector{
			Recall: &context.PipeRecall{},
			Tree:   []filter.Filter{constFilter("a")},
		},
		Next: &filter.Pipe{
			Filter: &filter.Selector{
				Recall: &context.PipeRecall{},
				Tree:   []filter.Filter{constFilter("b")},
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
			Tree:   []filter.Filter{constFilter("a")},
		},
		Assignment: &context.SimpleAssignment{Name: "a"},
		Next: &filter.Pipe{
			Filter: &filter.Selector{
				Recall: &context.PipeRecall{},
				Tree:   []filter.Filter{constFilter("b")},
			},
		},
	}, r)
}

func stringFilter(in string) *filter.String {
	return &filter.String{
		Filters: []filter.Filter{constFilter(in)},
	}
}

func constFilter(in interface{}) *filter.Const {
	return &filter.Const{context.NewConstValuer(in)}
}
