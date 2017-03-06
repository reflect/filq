// Portions of this file are derived from jparsec, a parser combinator library for Java.
//
// Copyright (C) jparsec.org
//
// https://github.com/jparsec/jparsec/blob/d84f3065cf3cea1ec500ad81ade858741eee62c2/jparsec/src/main/java/org/jparsec/OperatorTable.java
// https://github.com/jparsec/jparsec/blob/d84f3065cf3cea1ec500ad81ade858741eee62c2/jparsec/src/main/java/org/jparsec/Parser.java

package parser

import (
	"sort"

	"github.com/jmikkola/parsego/parser"
)

type OperatorAssociativity int

const (
	OperatorPrefix OperatorAssociativity = iota
	OperatorPostfix
	OperatorLeftInfix
	OperatorNonAssociativeInfix
	OperatorRightInfix
)

type Operator struct {
	p             parser.Parser
	precedence    int
	associativity OperatorAssociativity
}

type operatorSort []Operator

func (s operatorSort) Len() int {
	return len(s)
}

func (s operatorSort) Less(i, j int) bool {
	if s[i].precedence > s[j].precedence {
		return true
	}

	if s[i].precedence < s[j].precedence {
		return false
	}

	return s[i].associativity < s[j].associativity
}

func (s operatorSort) Swap(i, j int) {
	s[i], s[j] = s[j], s[i]
}

type OperatorTable struct {
	ops []Operator
}

func (ot *OperatorTable) LeftInfix(in parser.Parser, precedence int) *OperatorTable {
	ot.ops = append(
		ot.ops,
		Operator{p: in, precedence: precedence, associativity: OperatorLeftInfix},
	)

	return ot
}

func (ot *OperatorTable) RightInfix(in parser.Parser, precedence int) *OperatorTable {
	ot.ops = append(
		ot.ops,
		Operator{p: in, precedence: precedence, associativity: OperatorRightInfix},
	)

	return ot
}

func (ot *OperatorTable) NonAssociativeInfix(in parser.Parser, precedence int) *OperatorTable {
	ot.ops = append(
		ot.ops,
		Operator{p: in, precedence: precedence, associativity: OperatorNonAssociativeInfix},
	)

	return ot
}

func (ot *OperatorTable) Prefix(in parser.Parser, precedence int) *OperatorTable {
	ot.ops = append(
		ot.ops,
		Operator{p: in, precedence: precedence, associativity: OperatorPrefix},
	)

	return ot
}

func (ot *OperatorTable) Postfix(in parser.Parser, precedence int) *OperatorTable {
	ot.ops = append(
		ot.ops,
		Operator{p: in, precedence: precedence, associativity: OperatorPostfix},
	)

	return ot
}

func (ot *OperatorTable) Parser(operand parser.Parser) parser.Parser {
	if len(ot.ops) == 0 {
		return operand
	}

	sort.Sort(operatorSort(ot.ops))

	r := operand

	begin := 0
	precedence := ot.ops[0].precedence
	associativity := ot.ops[0].associativity

	end := 0

	for i := 1; i < len(ot.ops); i++ {
		end = i
		if ot.ops[i].precedence == precedence && ot.ops[i].associativity == associativity {
			continue
		}

		next := ot.slice(begin, end)
		r = ot.build(next, associativity, r)

		begin = i
		precedence = ot.ops[i].precedence
		associativity = ot.ops[i].associativity
	}

	if end != len(ot.ops) {
		end = len(ot.ops)

		associativity = ot.ops[begin].associativity
		next := ot.slice(begin, end)
		r = ot.build(next, associativity, r)
	}

	return r
}

func (ot *OperatorTable) slice(begin, end int) parser.Parser {
	parsers := make([]parser.Parser, end-begin)
	for i := 0; i < len(parsers); i++ {
		parsers[i] = ot.ops[i+begin].p
	}

	return parser.Or(parsers...)
}

func (ot *OperatorTable) build(op parser.Parser, associativity OperatorAssociativity, operand parser.Parser) parser.Parser {
	fn := map[OperatorAssociativity]func(op, operand parser.Parser) parser.Parser{
		OperatorPrefix:              Prefix,
		OperatorPostfix:             Postfix,
		OperatorLeftInfix:           LeftInfix,
		OperatorRightInfix:          RightInfix,
		OperatorNonAssociativeInfix: NonAssociativeInfix,
	}[associativity]

	return fn(op, operand)
}

func NewOperatorTable() *OperatorTable {
	return &OperatorTable{}
}

type UnaryOperation struct {
	Operator interface{}
	Operand  interface{}
}

type BinaryOperation struct {
	Operator    interface{}
	Left, Right interface{}
}

func Prefix(op, operand parser.Parser) parser.Parser {
	return parser.Map([]parser.Named{
		{Name: "operator", Parser: parser.ListOf(op)},
		{Name: "operand", Parser: operand},
	}, func(m map[string]interface{}) interface{} {
		out := m["operand"]

		ops := m["operator"].([]interface{})
		for i := len(ops) - 1; i >= 0; i-- {
			out = &UnaryOperation{
				Operator: ops[i],
				Operand:  out,
			}
		}

		return out
	})
}

func Postfix(op, operand parser.Parser) parser.Parser {
	return parser.Map([]parser.Named{
		{Name: "operand", Parser: operand},
		{Name: "operator", Parser: parser.ListOf(op)},
	}, func(m map[string]interface{}) interface{} {
		out := m["operand"]

		for _, op := range m["operator"].([]interface{}) {
			out = &UnaryOperation{
				Operator: op,
				Operand:  out,
			}
		}

		return out
	})
}

func LeftInfix(op, operand parser.Parser) parser.Parser {
	return Next(operand, func(first interface{}) parser.Parser {
		cont := parser.Map([]parser.Named{
			{Name: "operator", Parser: op},
			{Name: "operand", Parser: operand},
		}, func(m map[string]interface{}) interface{} {
			// NB: Left not yet filled in.
			return &BinaryOperation{
				Operator: m["operator"],
				Right:    m["operand"],
			}
		})

		return parser.ParseWith(
			parser.ListOf(cont),
			func(in interface{}) interface{} {
				out := first

				for _, next := range in.([]interface{}) {
					op := next.(*BinaryOperation)
					op.Left = out

					out = op
				}

				return out
			},
		)
	})
}

type rhs struct {
	Operator, Operand interface{}
}

func RightInfix(op, operand parser.Parser) parser.Parser {
	next := parser.Map([]parser.Named{
		{Name: "operator", Parser: op},
		{Name: "operand", Parser: operand},
	}, func(m map[string]interface{}) interface{} {
		return rhs{Operator: m["operator"], Operand: m["operand"]}
	})

	return parser.Map([]parser.Named{
		{Name: "operand", Parser: operand},
		{Name: "operator", Parser: parser.ListOf(next)},
	}, func(m map[string]interface{}) interface{} {
		nexts := m["operator"].([]interface{})
		if len(nexts) == 0 {
			return m["operand"]
		}

		right := nexts[len(nexts)-1].(rhs)

		operand := right.Operand
		for i := len(nexts) - 1; i >= 1; i-- {
			left := nexts[i-1].(rhs)
			operand = &BinaryOperation{
				Operator: right.Operator,
				Left:     left.Operand,
				Right:    operand,
			}

			right = left
		}

		return &BinaryOperation{
			Operator: right.Operator,
			Left:     m["operand"],
			Right:    operand,
		}
	})
}

func NonAssociativeInfix(op, operand parser.Parser) parser.Parser {
	return Next(operand, func(a interface{}) parser.Parser {
		shift := parser.Map([]parser.Named{
			{Name: "operator", Parser: op},
			{Name: "operand", Parser: operand},
		}, func(m map[string]interface{}) interface{} {
			return &BinaryOperation{
				Operator: m["operator"],
				Left:     a,
				Right:    m["operand"],
			}
		})
		return parser.Or(shift, Always(a))
	})
}
