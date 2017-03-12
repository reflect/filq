package parser

import (
	"github.com/reflect/filq/filter"
	"github.com/reflect/parsego/parser"
)

func Scoped(in parser.Parser) parser.Parser {
	return parser.ParseWith(
		in,
		func(in interface{}) interface{} {
			return &filter.Scope{Filter: in.(filter.Filter)}
		},
	)
}
