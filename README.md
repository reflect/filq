# filq

filq is a library for implementing a stream processing pipeline for structured
data. filq is inspired by jq, and uses essentially the same language.

## Getting started

You can implement a working program using filq in just a few lines of code,
especially if you use the provided standard libraries.

```go
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
```

This is actually the
[`basicq` example](https://github.com/reflect/filq/tree/master/examples/basicq).
This handles input as a byte array, so to parse it as JSON, you'll need to
invoke the `fromjson` function as a filter first:

```
$ go build
$ echo '{"amazing":["just","like","jq"]}' | ./basicq 'fromjson | .amazing[]'
just
like
jq
```
