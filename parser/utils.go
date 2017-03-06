package parser

import (
	"fmt"
	"unicode"

	"github.com/jmikkola/parsego/parser"
	"github.com/jmikkola/parsego/parser/result"
	"github.com/jmikkola/parsego/parser/scanner"
	"github.com/jmikkola/parsego/parser/textpos"
	"github.com/reflect/filq/filter"
)

type IsRuneFunc func(c rune) bool

type IsRuneParser struct {
	fns    []IsRuneFunc
	invert bool
}

func IsRune(fns ...IsRuneFunc) parser.Parser {
	return &IsRuneParser{fns: fns, invert: false}
}

func (p *IsRuneParser) Parse(sc scanner.Scanner) result.ParseResult {
	start := sc.GetPos()
	r, err := sc.Read()
	if err != nil {
		return result.Failed(textpos.Single(sc.GetPos()), fmt.Errorf("expected a character, got error %v", err))
	}

	if !p.invert {
		for _, fn := range p.fns {
			if fn(r) {
				return result.Success(textpos.Range(start, sc.GetPos()), string(r))
			}
		}

		return result.Failed(textpos.Single(sc.GetPos()), fmt.Errorf("expected a unicode character in range, got %c", r))
	} else {
		for _, fn := range p.fns {
			if fn(r) {
				return result.Failed(textpos.Single(sc.GetPos()), fmt.Errorf("expected a unicode character in range, got %c", r))
			}
		}

		return result.Success(textpos.Range(start, sc.GetPos()), string(r))
	}
}

type AlwaysParser struct {
	value interface{}
}

func Always(value interface{}) parser.Parser {
	return &AlwaysParser{value: value}
}

func (p *AlwaysParser) Parse(sc scanner.Scanner) result.ParseResult {
	return result.Success(textpos.Single(sc.GetPos()), p.value)
}

type NextParser struct {
	inner parser.Parser
	fn    func(interface{}) parser.Parser
}

func Next(in parser.Parser, next func(interface{}) parser.Parser) parser.Parser {
	return &NextParser{
		inner: in,
		fn:    next,
	}
}

func (p *NextParser) Parse(sc scanner.Scanner) result.ParseResult {
	result := p.inner.Parse(sc)
	if !result.Matched() {
		return result
	}

	next := p.fn(result.Result())
	return next.Parse(sc)
}

func NotRune(fns ...IsRuneFunc) parser.Parser {
	return &IsRuneParser{fns: fns, invert: true}
}

func Whitespace1() parser.Parser {
	return IsRune(unicode.IsSpace)
}

func Whitespace() parser.Parser {
	return parser.Maybe(Whitespace1())
}

func Sep(in parser.Parser) parser.Parser {
	return parser.Surround(Whitespace(), in, Whitespace())
}

func N(which int, parsers ...parser.Parser) parser.Parser {
	return parser.ParseWith(
		parser.Sequence(parsers...),
		func(in interface{}) interface{} {
			seq := in.([]interface{})
			return seq[which]
		},
	)
}

func First(parsers ...parser.Parser) parser.Parser {
	return N(0, parsers...)
}

func Second(parsers ...parser.Parser) parser.Parser {
	return N(1, parsers...)
}

func Default(in parser.Parser, otherwise interface{}) parser.Parser {
	return parser.Or(in, Always(otherwise))
}

func Scoped(in parser.Parser) parser.Parser {
	return parser.ParseWith(
		in,
		func(in interface{}) interface{} {
			return &filter.Scope{Filter: in.(filter.Filter)}
		},
	)
}

func Flatten(parsers ...parser.Parser) parser.Parser {
	var flatten func(in []interface{}) []interface{}
	flatten = func(in []interface{}) (out []interface{}) {
		for _, c := range in {
			if seq, ok := c.([]interface{}); ok {
				out = append(out, flatten(seq)...)
			} else {
				out = append(out, c)
			}
		}

		return
	}

	return parser.ParseWith(
		parser.Sequence(parsers...),
		func(in interface{}) interface{} {
			return flatten(in.([]interface{}))
		},
	)
}
