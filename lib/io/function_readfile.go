package io

import (
	"io/ioutil"
	"os"

	"github.com/pkg/errors"
	"github.com/reflect/filq/context"
	"github.com/reflect/filq/types"
)

type readFile struct {
	path string
}

func (rf *readFile) Value(ctx *context.Context) (interface{}, error) {
	f, err := os.Open(rf.path)
	if err != nil {
		return nil, errors.Wrapf(err, "opening file %s", rf.path)
	}
	defer f.Close()

	b, err := ioutil.ReadAll(f)
	if err != nil {
		return nil, errors.Wrapf(err, "reading file %s", rf.path)
	}

	return ctx.Convert(b), nil
}

func ReadFile(ctx *context.Context, in context.Valuer) ([]context.Valuer, error) {
	v, err := in.Value(ctx)
	if err != nil {
		return nil, err
	}

	var path string
	if s, ok := v.(types.Str); ok {
		path = string(s)
	} else if b, ok := v.(types.Bytes); ok {
		path = string(b)
	}

	return []context.Valuer{&readFile{path}}, nil
}
