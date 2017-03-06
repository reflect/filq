// Portions of this file are derived from pretty, a pretty-printing package for Go values.
//
// Copyright 2012 Keith Rarick
//
// https://github.com/kr/pretty/blob/cfb55aafdaf3ec08f0db22699ab822c50091b1c4/formatter.go

package types

import (
	"fmt"
	"unicode"
)

func FormatDefault(f fmt.State, c rune, v interface{}) {
	s := "%"
	for i := 0; i < unicode.MaxASCII; i++ {
		if f.Flag(i) {
			s += string(i)
		}
	}
	if w, ok := f.Width(); ok {
		s += fmt.Sprintf("%d", w)
	}
	if p, ok := f.Precision(); ok {
		s += fmt.Sprintf(".%d", p)
	}
	s += string(c)
	fmt.Fprintf(f, s, v)
}
