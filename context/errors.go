package context

import (
	"fmt"
	"reflect"
)

type VariableNotDefinedError struct {
	Name string
}

func (e *VariableNotDefinedError) Error() string {
	return fmt.Sprintf("variable %q not defined", e.Name)
}

type FunctionNotDefinedError struct {
	Name string
}

func (e *FunctionNotDefinedError) Error() string {
	return fmt.Sprintf("function %q not defined", e.Name)
}

type UnexpectedTypeError struct {
	Wanted []reflect.Type
	Got    reflect.Type
}

func (e *UnexpectedTypeError) Error() string {
	return fmt.Sprintf("unexpected type %s (wanted one of %s)", e.Got, e.Wanted)
}
