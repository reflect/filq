package io

import (
	"os"

	"github.com/reflect/filq/context"
)

func AppendFile(ctx *context.Context, in context.Valuer, paths []context.Valuer) ([]context.Valuer, error) {
	return writeFileWithMode(ctx, in, paths, os.O_APPEND)
}
