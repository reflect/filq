package filter

import (
	"fmt"

	"github.com/reflect/filq/context"
)

type Op2 struct {
	Operator    string
	Left, Right Filter
}

func (o *Op2) Apply(ctx *context.Context, in context.Valuer) ([]context.Valuer, error) {
	lvs, err := o.Left.Apply(ctx, in)
	if err != nil {
		return nil, err
	}

	rvs, err := o.Right.Apply(ctx, in)
	if err != nil {
		return nil, err
	}

	out := make([]context.Valuer, len(lvs)*len(rvs))
	for i, lv := range lvs {
		for j, rv := range rvs {
			k := len(lvs)*i + j
			switch o.Operator {
			case "*":
				out[k] = &op2NumFilter{fn: op2Mul, l: lv, r: rv}
			case "/":
				out[k] = &op2NumFilter{fn: op2Div, l: lv, r: rv}
			case "+":
				out[k] = &op2AddFilter{l: lv, r: rv}
			case "-":
				out[k] = &op2NumFilter{fn: op2Sub, l: lv, r: rv}
			case "<":
				out[k] = &op2CmpFilter{fn: op2Lt, l: lv, r: rv}
			case "<=":
				out[k] = &op2CmpFilter{fn: op2Lte, l: lv, r: rv}
			case ">":
				out[k] = &op2CmpFilter{fn: op2Gt, l: lv, r: rv}
			case ">=":
				out[k] = &op2CmpFilter{fn: op2Gte, l: lv, r: rv}
			case "==":
				out[k] = &op2EqualFilter{l: lv, r: rv}
			case "!=":
				out[k] = &op2EqualFilter{l: lv, r: rv, inverse: true}
			case "and":
				out[k] = &op2BoolFilter{fn: op2And, l: lv, r: rv}
			case "or":
				out[k] = &op2BoolFilter{fn: op2Or, l: lv, r: rv}
			default:
				panic(fmt.Errorf("binary operator %q not implemented", o.Operator))
			}
		}
	}

	return out, nil
}
