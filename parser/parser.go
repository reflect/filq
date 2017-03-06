package parser

import (
	"fmt"
	"strconv"
	"unicode"

	"github.com/jmikkola/parsego/parser"
	"github.com/reflect/filq/context"
	"github.com/reflect/filq/filter"
	"github.com/reflect/filq/types"
)

func nullParser() parser.Parser {
	return parser.TokenAs("null", nil)
}

func boolParser() parser.Parser {
	return parser.Or(
		parser.TokenAs("true", true),
		parser.TokenAs("false", false),
	)
}

func numberParser() parser.Parser {
	decimal := parser.Sequence(parser.Char('.'), parser.Digits())
	exp := parser.Sequence(parser.AnyChar('e', 'E'), parser.Maybe(parser.AnyChar('-', '+')), parser.Digits())

	return parser.Map([]parser.Named{
		{Name: "integer", Parser: parser.Sequence(parser.Maybe(parser.Char('-')), parser.Digits())},
		{Name: "decimal", Parser: parser.Maybe(decimal)},
		{Name: "exp", Parser: parser.Maybe(exp)},
	}, func(m map[string]interface{}) interface{} {
		if m["decimal"] == "" && m["exp"] == "" {
			r, _ := strconv.ParseInt(m["integer"].(string), 10, 64)
			return r
		}

		r, _ := strconv.ParseFloat(fmt.Sprintf("%s%s%s", m["integer"], m["decimal"], m["exp"]), 64)
		return r
	})
}

func sliceParser() parser.Parser {
	return parser.Map([]parser.Named{
		{Name: "left", Parser: numberParser()},
		{Parser: Sep(parser.Char(':'))},
		{Name: "right", Parser: numberParser()},
	}, func(m map[string]interface{}) interface{} {
		s := types.Slice{}

		switch lt := m["left"].(type) {
		case int64:
			s.Left = lt
		case float64:
			s.Left = int64(lt)
		}

		switch rt := m["right"].(type) {
		case int64:
			s.Right = rt
		case float64:
			s.Right = int64(rt)
		}

		return s
	})
}

func stringParser() parser.Parser {
	hex := parser.Or(parser.CharRange('a', 'f'), parser.CharRange('A', 'F'), parser.Digit())
	u := parser.Sequence(parser.Char('u'), hex, hex, hex, hex)
	escape := parser.Sequence(
		parser.Char('\\'),
		parser.Or(u, parser.AnyChar('"', '\\', '/', 'b', 'f', 'n', 'r', 't')),
	)
	chars := NotRune(
		unicode.IsControl,
		func(c rune) bool { return c == '"' || c == '\\' },
	)
	valid := parser.Or(chars, escape)
	str := parser.Many(valid)

	return parser.ParseWith(
		parser.Sequence(parser.Char('"'), str, parser.Char('"')),
		func(in interface{}) interface{} {
			s, _ := strconv.Unquote(in.(string))
			return types.Str(s)
		},
	)
}

func identParser() parser.Parser {
	start := parser.Or(parser.Char('_'), parser.LowerLetter(), parser.UpperLetter())
	cont := parser.Or(start, parser.Digit())

	return parser.Sequence(start, parser.Many(cont))
}

func variableParser() parser.Parser {
	return parser.Sequence(parser.Ignore(parser.Char('$')), identParser())
}

func constParser(in parser.Parser) parser.Parser {
	return parser.ParseWith(
		in,
		func(in interface{}) interface{} {
			return &filter.Const{
				Valuer: context.NewConstValuer(in),
			}
		},
	)
}

func subscriptParser() parser.Parser {
	return parser.Surround(
		parser.Sequence(parser.Char('['), parser.Whitespace()),
		parser.Or(constParser(sliceParser()), exprParser()),
		parser.Sequence(parser.Whitespace(), parser.Char(']')),
	)
}

func recallParser() parser.Parser {
	return parser.ParseWith(
		parser.Sequence(parser.Maybe(variableParser()), parser.Ignore(parser.Char('.'))),
		func(in interface{}) interface{} {
			v := in.(string)

			switch v {
			case "":
				return &context.PipeRecall{}
			default:
				return &context.VariableRecall{Name: v}
			}
		},
	)
}

func selectorMapper(m map[string]interface{}) interface{} {
	sel := &filter.Selector{}

	if recall, ok := m["recall"].(string); ok {
		sel.Recall = &context.VariableRecall{Name: recall}
	} else {
		sel.Recall = &context.PipeRecall{}
	}

	seq := m["selector"].([]interface{})

	tree := make([]filter.Filter, len(seq))
	for i, c := range seq {
		tree[i] = c.(filter.Filter)
	}

	sel.Tree = tree
	return sel
}

func selectorParser(recall bool) parser.Parser {
	ident := parser.ParseWith(
		identParser(),
		func(in interface{}) interface{} {
			return types.Str(in.(string))
		},
	)

	start := parser.Or(constParser(ident), constParser(stringParser()), subscriptParser())
	cont := parser.Or(
		Second(parser.Char('.'), constParser(ident)),
		Second(parser.Char('.'), constParser(stringParser())),
		subscriptParser(),
	)
	selector := Flatten(start, parser.ListOf(cont))

	all := parser.Map([]parser.Named{
		{Parser: parser.Char('.')},
		{Name: "selector", Parser: Default(selector, []interface{}{})},
	}, selectorMapper)

	if !recall {
		return all
	}

	return parser.Or(
		parser.Map([]parser.Named{
			{Name: "recall", Parser: variableParser()},
			{Name: "selector", Parser: Default(parser.ListOf(cont), []interface{}{})},
		}, selectorMapper),
		all,
	)
}

func expandParser() parser.Parser {
	var mapper func(in interface{}) interface{}
	mapper = func(in interface{}) interface{} {
		if op, ok := in.(*UnaryOperation); ok {
			return &filter.Expand{
				Filter: mapper(op.Operand).(filter.Filter),
			}
		}

		return in
	}

	expansion := parser.Sequence(parser.Char('['), Whitespace(), parser.Char(']'))
	start := parser.ParseWith(Postfix(expansion, selectorParser(true)), mapper)
	cont := parser.ParseWith(Postfix(expansion, selectorParser(false)), mapper)

	return parser.ParseWith(
		Flatten(start, parser.ListOf(cont)),
		func(in interface{}) interface{} {
			seq := in.([]interface{})
			if len(seq) == 1 {
				return seq[0]
			}

			pipe := &filter.Pipe{
				Filter: seq[len(seq)-1].(filter.Filter),
			}

			for i := len(seq) - 2; i >= 0; i-- {
				pipe = &filter.Pipe{
					Filter: seq[i].(filter.Filter),
					Next:   pipe,
				}
			}

			return pipe
		},
	)
}

func funcParser() parser.Parser {
	sep := Sep(parser.Char(';'))

	argumentParser := parser.Map([]parser.Named{
		{Parser: parser.Char('(')},
		{Parser: Whitespace()},
		{Name: "arguments", Parser: parser.ManySepBy(pipelineParser(), sep)},
		{Parser: Whitespace()},
		{Parser: parser.Char(')')},
	}, func(m map[string]interface{}) interface{} {
		seq := m["arguments"].([]interface{})

		arguments := make([]filter.Filter, len(seq))
		for i, c := range seq {
			arguments[i] = c.(filter.Filter)
		}

		return arguments
	})

	return parser.Map([]parser.Named{
		{Name: "function", Parser: identParser()},
		{Name: "arguments", Parser: parser.Maybe(argumentParser)},
	}, func(m map[string]interface{}) interface{} {
		arguments, _ := m["arguments"].([]filter.Filter)

		return &filter.Call{
			Function:  m["function"].(string),
			Arguments: arguments,
		}
	})
}

func exprParser() parser.Parser {
	pipeline := parser.Surround(
		Sep(parser.Char('(')),
		Scoped(pipelineParser()),
		Sep(parser.Char(')')),
	)

	var mapper func(in interface{}) interface{}
	mapper = func(in interface{}) interface{} {
		if op, ok := in.(*BinaryOperation); ok {
			return &filter.Op2{
				Operator: op.Operator.(string),
				Left:     mapper(op.Left).(filter.Filter),
				Right:    mapper(op.Right).(filter.Filter),
			}
		} else if op, ok := in.(*UnaryOperation); ok {
			return &filter.Op1{
				Operator: op.Operator.(string),
				Operand:  mapper(op.Operand).(filter.Filter),
			}
		}

		return in
	}

	return parser.Lazy(func() parser.Parser {
		base := funcParser()
		base = parser.Or(constParser(parser.Or(nullParser(), boolParser(), stringParser(), numberParser())), base)
		base = parser.Or(pipeline, expandParser(), base)
		base = NewOperatorTable().
			LeftInfix(Sep(parser.Token("*")), 90).
			LeftInfix(Sep(parser.Token("/")), 90).
			LeftInfix(Sep(parser.Token("+")), 90).
			LeftInfix(Sep(parser.Token("-")), 90).
			LeftInfix(Sep(parser.Token("<=")), 60).
			LeftInfix(Sep(parser.Token("<")), 60).
			LeftInfix(Sep(parser.Token(">=")), 60).
			LeftInfix(Sep(parser.Token(">")), 60).
			LeftInfix(Sep(parser.Token("==")), 60).
			LeftInfix(Sep(parser.Token("!=")), 60).
			Prefix(Sep(parser.Token("not")), 50).
			LeftInfix(Sep(parser.Token("and")), 40).
			LeftInfix(Sep(parser.Token("or")), 30).
			Parser(base)

		return parser.ParseWith(base, mapper)
	})
}

func assignmentParser() parser.Parser {
	return parser.ParseWith(
		variableParser(),
		func(in interface{}) interface{} {
			return &context.SimpleAssignment{Name: in.(string)}
		},
	)
}

func pipelineParser() parser.Parser {
	pipe := Sep(parser.Char('|'))
	assignment := N(3,
		parser.Whitespace1(),
		parser.Token("as"),
		parser.Whitespace1(),
		assignmentParser(),
	)

	mapper := func(m map[string]interface{}) interface{} {
		pipe := &filter.Pipe{}

		if assignment, ok := m["assignment"].(context.Assignment); ok {
			pipe.Assignment = assignment
		}

		if next, ok := m["next"].(*filter.Pipe); ok {
			pipe.Next = next
		}

		return pipe
	}

	return parser.Lazy(func() parser.Parser {
		next := Second(pipe, pipelineParser())

		conts := parser.Or(
			parser.Map([]parser.Named{
				{Name: "assignment", Parser: assignment},
				{Name: "next", Parser: next},
			}, mapper),
			parser.Map([]parser.Named{
				{Name: "next", Parser: parser.Maybe(next)},
			}, mapper),
		)

		return parser.Map([]parser.Named{
			{Name: "expression", Parser: exprParser()},
			{Name: "pipe", Parser: conts},
		}, func(m map[string]interface{}) interface{} {
			pipe := m["pipe"].(*filter.Pipe)
			pipe.Filter = m["expression"].(filter.Filter)

			return pipe
		})
	})
}

type Parser struct {
	backend parser.Parser
}

func (p *Parser) ParseString(in string) (filter.Filter, error) {
	f, err := parser.ParseString(p.backend, in)
	if err != nil {
		return nil, err
	}

	return f.(filter.Filter), nil
}

func NewParser() *Parser {
	initial := parser.Surround(
		Whitespace(),
		Scoped(pipelineParser()),
		Whitespace(),
	)

	consumed := parser.ParseWith(
		initial,
		func(in interface{}) interface{} {
			return &filter.Consume{Filter: in.(filter.Filter)}
		},
	)

	return &Parser{
		backend: First(consumed, parser.EOF()),
	}
}
