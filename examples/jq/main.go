package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"os"

	"github.com/pkg/errors"
	"github.com/reflect/filq"
	fj "github.com/reflect/filq/lib/json"
)

func run() error {
	ctx := filq.NewContext()
	filter, err := filq.NewParser().ParseString(flag.Arg(0))
	if err != nil {
		return err
	}

	decoder := json.NewDecoder(os.Stdin)

	for {
		vr, err := fj.NewValuer(decoder)
		if err != nil {
			return err
		}

		outs, err := filter.Apply(ctx, vr)
		if errors.Cause(err) == io.EOF {
			break
		} else if err != nil {
			return err
		}

		for _, out := range outs {
			v, err := out.Value(ctx)
			if err != nil {
				return err
			}

			fmt.Printf("%+v\n", v)
		}
	}

	return nil
}

func main() {
	flag.Parse()

	if err := run(); err != nil {
		fmt.Fprintf(os.Stderr, "%+v\n", err)
		os.Exit(1)
	}
}
