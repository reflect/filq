package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"os"

	"github.com/reflect/filq"
)

func main() {
	flag.Parse()

	b, err := ioutil.ReadAll(os.Stdin)
	outs, err := filq.Run(filq.NewContext(), flag.Arg(0), b)
	if err != nil {
		fmt.Fprintf(os.Stderr, "%+v\n", err)
		os.Exit(1)
	}

	for _, out := range outs {
		fmt.Printf("%+v\n", out)
	}
}
