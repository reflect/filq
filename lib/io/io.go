package io

import (
	"github.com/reflect/filq/context"
	"github.com/reflect/filq/function"
)

func DefineIn(ctx *context.Context) {
	fn, _ := function.NewFunction(AppendFile)
	ctx.DefineFunction("appendfile", fn)

	fn, _ = function.NewFunction(ReadFile)
	ctx.DefineFunction("readfile", fn)

	fn, _ = function.NewFunction(WriteFile)
	ctx.DefineFunction("writefile", fn)
}

func NewValuer(path string) (context.Valuer, error) {
	return &readFile{path}, nil
}
