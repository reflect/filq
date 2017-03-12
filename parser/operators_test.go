package parser

import (
	"testing"

	"github.com/reflect/parsego/parser"
	"github.com/stretchr/testify/assert"
)

func TestPrefix(t *testing.T) {
	p := Prefix(parser.AnyChar('!', '~'), parser.Digits())
	r, err := parser.ParseString(p, "!~32")
	assert.NoError(t, err)
	assert.Equal(t, &UnaryOperation{
		Operator: "!",
		Operand: &UnaryOperation{
			Operator: "~",
			Operand:  "32",
		},
	}, r)
}

func TestPostfix(t *testing.T) {
	p := Postfix(parser.AnyChar('?', '!'), parser.Digits())
	r, err := parser.ParseString(p, "32?!")
	assert.NoError(t, err)
	assert.Equal(t, &UnaryOperation{
		Operator: "!",
		Operand: &UnaryOperation{
			Operator: "?",
			Operand:  "32",
		},
	}, r)
}

func TestLeftInfix(t *testing.T) {
	p := LeftInfix(parser.Token("+"), parser.Digits())
	r, err := parser.ParseString(p, "20+30+40")
	assert.NoError(t, err)
	assert.Equal(t, &BinaryOperation{
		Operator: "+",
		Left: &BinaryOperation{
			Operator: "+",
			Left:     "20",
			Right:    "30",
		},
		Right: "40",
	}, r)
}

func TestRightInfix(t *testing.T) {
	p := RightInfix(parser.Token("+"), parser.Digits())
	r, err := parser.ParseString(p, "20+30+40")
	assert.NoError(t, err)
	assert.Equal(t, &BinaryOperation{
		Operator: "+",
		Left:     "20",
		Right: &BinaryOperation{
			Operator: "+",
			Left:     "30",
			Right:    "40",
		},
	}, r)
}

func TestOperatorTable(t *testing.T) {
	p := NewOperatorTable().
		RightInfix(parser.Token("**"), 100).
		LeftInfix(parser.Token("*"), 90).
		LeftInfix(parser.Token("+"), 80).
		LeftInfix(parser.Token("-"), 80).
		Prefix(parser.Token("X"), 60).
		Parser(parser.Digits())
	r, err := parser.ParseString(p, "X10+20-3**4*2+30")
	assert.NoError(t, err)
	assert.Equal(t, &UnaryOperation{
		Operator: "X",
		Operand: &BinaryOperation{
			Operator: "+",
			Left: &BinaryOperation{
				Operator: "-",
				Left: &BinaryOperation{
					Operator: "+",
					Left:     "10",
					Right:    "20",
				},
				Right: &BinaryOperation{
					Operator: "*",
					Left: &BinaryOperation{
						Operator: "**",
						Left:     "3",
						Right:    "4",
					},
					Right: "2",
				},
			},
			Right: "30",
		},
	}, r)

	p = NewOperatorTable().
		LeftInfix(parser.Sep(parser.Token("<=")), 70).
		LeftInfix(parser.Sep(parser.Token("<")), 70).
		LeftInfix(parser.Sep(parser.Token(">=")), 70).
		LeftInfix(parser.Sep(parser.Token(">")), 70).
		Parser(parser.Digits())
	r, err = parser.ParseString(parser.First(p, parser.EOF()), "1 <= 2")
	assert.NoError(t, err)
	assert.Equal(t, &BinaryOperation{
		Operator: "<=",
		Left:     "1",
		Right:    "2",
	}, r)
}
