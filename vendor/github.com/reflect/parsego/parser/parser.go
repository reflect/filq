package parser

import (
	"bytes"
	"fmt"

	"github.com/reflect/parsego/parser/result"
	"github.com/reflect/parsego/parser/scanner"
	"github.com/reflect/parsego/parser/textpos"
)

// Parser defines the interface implemented by all combinable parsers.
type Parser interface {
	Parse(sc scanner.Scanner) result.ParseResult
}

func fail(at textpos.TextPos, format string, a ...interface{}) result.ParseResult {
	return result.Failed(textpos.Single(at), fmt.Errorf(format, a...))
}

// EOFParser expects just EOF.
type EOFParser struct{}

// EOF returns a parser that expects just EOF.
func EOF() Parser {
	return &EOFParser{}
}

// Parse parses the input.
func (p *EOFParser) Parse(sc scanner.Scanner) result.ParseResult {
	r, err := sc.Read()
	if err == nil {
		return fail(sc.GetPos(), "expected EOF, got %c", r)
	}
	return result.Success(textpos.Single(sc.GetPos()), "")
}

// CharRangeParser parses any single character in a range, inclusive.
type CharRangeParser struct {
	min rune // both inclusive
	max rune
}

// Char returns a parser that parses a single occurrence of that rune.
func Char(c rune) Parser {
	return &CharRangeParser{c, c}
}

// CharRange returns a parser that parses a single occurrence of any
// rune in the given range, inclusive.
func CharRange(min, max rune) Parser {
	return &CharRangeParser{min, max}
}

// Parse parses the input.
func (p *CharRangeParser) Parse(sc scanner.Scanner) result.ParseResult {
	start := sc.GetPos()
	r, err := sc.Read()
	if err != nil {
		return fail(sc.GetPos(), "expected a character, got error %v", err)
	}
	if r < p.min || r > p.max {
		return fail(sc.GetPos(),
			"expected a character in the range '%c' to '%c', got error %c",
			p.min, p.max, r)
	}
	return result.Success(
		textpos.Range(start, sc.GetPos()),
		string(r))
}

// TokenParser works like a series of CharRangeParsers, but is more
// efficient.
type TokenParser struct {
	token string
}

// Token returns a parser that parses the exact string given.
func Token(token string) Parser {
	return &TokenParser{token}
}

// Parse parses the input.
func (p *TokenParser) Parse(sc scanner.Scanner) result.ParseResult {
	start := sc.GetPos()
	seen := []rune{}
	for _, c := range p.token {
		r, err := sc.Read()
		seen = append(seen, r)
		if err != nil {
			return fail(sc.GetPos(), "expected '%s', got error %v", p.token, err)
		}
		if r != c {
			return fail(sc.GetPos(), "expected '%s', got '%s'", p.token, string(seen))
		}
	}
	return result.Success(
		textpos.Range(start, sc.GetPos()),
		string(seen))
}

// CharSetParser parses any single character in the set.
type CharSetParser struct {
	allowed map[rune]struct{}
	invert  bool
}

// AnyChar returns a parser that parses a single occurrence of any rune
// given.
func AnyChar(rs ...rune) Parser {
	return &CharSetParser{allowed: rarray2rmap(rs), invert: false}
}

// NoneOf returns a parser that parses a single occurrence of any rune
// other than the ones given.
func NoneOf(rs ...rune) Parser {
	return &CharSetParser{allowed: rarray2rmap(rs), invert: true}
}

// AnyCharIn returns a parser that parses a single occurrence of any
// rune in the given string.
func AnyCharIn(s string) Parser {
	return &CharSetParser{allowed: s2runemap(s), invert: false}
}

// AnyCharNotIn returns a parser that parses a single occurrence of any
// rune other than the ones in the given string.
func AnyCharNotIn(s string) Parser {
	return &CharSetParser{allowed: s2runemap(s), invert: true}
}

func s2runemap(s string) map[rune]struct{} {
	m := make(map[rune]struct{}, len(s))
	for _, r := range s {
		m[r] = struct{}{}
	}
	return m
}

func rarray2rmap(rs []rune) map[rune]struct{} {
	m := make(map[rune]struct{}, len(rs))
	for _, r := range rs {
		m[r] = struct{}{}
	}
	return m
}

// Parse parses the input.
func (p *CharSetParser) Parse(sc scanner.Scanner) result.ParseResult {
	start := sc.GetPos()
	r, err := sc.Read()
	if err != nil {
		return fail(sc.GetPos(), "expected a character, got error %v", err)
	}
	if _, ok := p.allowed[r]; ok == p.invert {
		return fail(sc.GetPos(), "expected a character in the set, got error %c", r)
	}
	return result.Success(textpos.Range(start, sc.GetPos()), string(r))
}

// SeqParser combines multiple parsers in sequence.
type SeqParser struct {
	parsers []Parser
	combine bool
}

// AllOf returns a parser that runs each parser in series.
func AllOf(parsers ...Parser) Parser {
	return &SeqParser{parsers: parsers, combine: false}
}

// Sequence returns a parser that runs each given parser in series and
// combines the result.
func Sequence(parsers ...Parser) Parser {
	return &SeqParser{parsers: parsers, combine: true}
}

// Parse parses the input.
func (p *SeqParser) Parse(sc scanner.Scanner) result.ParseResult {
	start := sc.GetPos()
	var end textpos.TextPos
	results := []interface{}{}

	for _, inner := range p.parsers {
		innerResult := inner.Parse(sc)
		// Return errors right away
		if !innerResult.Matched() {
			return innerResult
		}

		end = innerResult.TextRange().End()
		results = append(results, innerResult.Result())
	}

	var output interface{}
	if p.combine {
		output = cleanupResult(results)
	} else {
		output = results
	}

	return result.Success(textpos.Range(start, end), output)
}

func cleanupResult(results []interface{}) interface{} {
	var buffer bytes.Buffer
	allStr := true
	for _, result := range results {
		if result == "" {
			continue
		}
		if s, ok := result.(string); ok {
			buffer.WriteString(s)
		} else {
			allStr = false
			break
		}
	}
	if allStr {
		return buffer.String()
	}
	return results
}

// Wrapper modifies the result of a parser with a function.
type Wrapper struct {
	inner Parser
	fn    func(interface{}) interface{}
}

// ParseWith returns a parser that will apply the given function to
// the result of parsing, if the parser was successful.
func ParseWith(p Parser, fn func(interface{}) interface{}) Parser {
	return &Wrapper{inner: p, fn: fn}
}

// Parse parses the input.
func (p *Wrapper) Parse(sc scanner.Scanner) result.ParseResult {
	innerResult := p.inner.Parse(sc)
	if innerResult.Matched() {
		return result.Success(
			innerResult.TextRange(),
			p.fn(innerResult.Result()))
	}
	return innerResult
}

// MaybeParser tries to run the inner parser, but allows the inner
// parser to fail.
type MaybeParser struct {
	inner Parser
}

// Maybe returns a parser that parses 0 or 1 occurrences of the given
// parser.
func Maybe(inner Parser) Parser {
	return &MaybeParser{inner}
}

// Parse parses the input.
func (p *MaybeParser) Parse(sc scanner.Scanner) result.ParseResult {
	sc.StartSnapshot()

	innerResult := p.inner.Parse(sc)
	if innerResult.Matched() {
		sc.PopSnapshot()
		return innerResult
	}

	sc.RewindSnapshot()
	return result.Success(textpos.Single(sc.GetPos()), "")
}

// ManyParser Matches 0+ occurrences
type ManyParser struct {
	inner   Parser
	combine bool
}

// ListOf returns a parser that matches the given parser zero or more
// times, and returns a list of the results.
func ListOf(inner Parser) Parser {
	return &ManyParser{inner: inner, combine: false}
}

// Many returns a parser that matches the given parser zero or more
// times, and combines the results.
func Many(inner Parser) Parser {
	return &ManyParser{inner: inner, combine: true}
}

// Parse parses the input.
func (p *ManyParser) Parse(sc scanner.Scanner) result.ParseResult {
	start := sc.GetPos()
	results := []interface{}{}

	for true {
		sc.StartSnapshot()
		innerResult := p.inner.Parse(sc)

		if innerResult.Matched() {
			sc.PopSnapshot()
			results = append(results, innerResult.Result())
		} else {
			sc.RewindSnapshot()
			break
		}
	}

	var output interface{}
	if p.combine {
		output = cleanupResult(results)
	} else {
		output = results
	}

	return result.Success(textpos.Range(start, sc.GetPos()), output)
}

// OrParser parses at most one of the inner parses.
type OrParser struct {
	parsers []Parser
}

// Or returns a parser that accepts the union of the languages
// accepted by the given parsers.
func Or(parsers ...Parser) Parser {
	return &OrParser{parsers}
}

// Parse parses the input.
func (p *OrParser) Parse(sc scanner.Scanner) result.ParseResult {
	for _, inner := range p.parsers {
		sc.StartSnapshot()
		innerResult := inner.Parse(sc)

		if innerResult.Matched() {
			sc.PopSnapshot()
			return innerResult
		}
		sc.RewindSnapshot()
	}

	return fail(sc.GetPos(), "no parser matched")
}

// Named is used for arguments to Map
type Named struct {
	Name   string
	Parser Parser
}

// MapParser parses to a map of named components.
type MapParser struct {
	parsers []Named
	fn      func(map[string]interface{}) interface{}
}

// Map builds a parser that parses the named components in series,
// populating a map between the given names and the results of the
// given parsers. The output of parsers named "" is ignored.
func Map(parsers []Named, fn func(map[string]interface{}) interface{}) Parser {
	return &MapParser{
		parsers: parsers,
		fn:      fn,
	}
}

// Parse parses the input.
func (p *MapParser) Parse(sc scanner.Scanner) result.ParseResult {
	parsed := map[string]interface{}{}
	start := sc.GetPos()

	for _, named := range p.parsers {
		innerResult := named.Parser.Parse(sc)
		if !innerResult.Matched() {
			return innerResult
		}

		if named.Name != "" {
			parsed[named.Name] = innerResult.Result()
		}
	}

	return result.Success(textpos.Range(start, sc.GetPos()), p.fn(parsed))
}

// LazyFn contains a function that lazily constructs the real
// parser. Useful for constructing recursive parsers.
type LazyFn struct {
	fn func() Parser
}

// Lazy builds a lazily defined parser by calling the given function
// only when the parser is actually used.
func Lazy(fn func() Parser) Parser {
	return &LazyFn{fn}
}

// Parse parses the input.
func (p *LazyFn) Parse(sc scanner.Scanner) result.ParseResult {
	actual := p.fn()
	return actual.Parse(sc)
}

// IgnoreParser runs the inner parser, but then replaces the result
// with "".
type IgnoreParser struct {
	inner Parser
}

// Ignore ignores the result of the given parser.
func Ignore(inner Parser) Parser {
	return &IgnoreParser{inner}
}

// Parse parses the input.
func (p *IgnoreParser) Parse(sc scanner.Scanner) result.ParseResult {
	r := p.inner.Parse(sc)
	if r.Matched() {
		return result.Success(r.TextRange(), "")
	}
	return r
}

// IsRuneFunc is a function that returns true if the given rune satisfies its
// condition.
type IsRuneFunc func(c rune) bool

// IsRuneParser runs a single rune of input through the given satisfaction
// functions, and returns success if any functions succeed or fail, depending
// on whether it is inverted.
type IsRuneParser struct {
	fns    []IsRuneFunc
	invert bool
}

// IsRune tests whether a single character matches the given rune satisfaction
// functions, and returns success if any function succeeds.
func IsRune(fns ...IsRuneFunc) Parser {
	return &IsRuneParser{fns: fns, invert: false}
}

// Parse parses the input.
func (p *IsRuneParser) Parse(sc scanner.Scanner) result.ParseResult {
	start := sc.GetPos()
	r, err := sc.Read()
	if err != nil {
		return fail(sc.GetPos(), "expected a character, got error %v", err)
	}

	if p.invert {
		for _, fn := range p.fns {
			if fn(r) {
				return fail(sc.GetPos(), "expected a unicode character in range, got %c", r)
			}
		}

		return result.Success(textpos.Range(start, sc.GetPos()), string(r))
	}

	for _, fn := range p.fns {
		if fn(r) {
			return result.Success(textpos.Range(start, sc.GetPos()), string(r))
		}
	}

	return fail(sc.GetPos(), "expected a unicode character in range, got %c", r)
}

// NotRune tests whether a single character matches the given rune satisfaction
// functions, and returns success if any function fails.
func NotRune(fns ...IsRuneFunc) Parser {
	return &IsRuneParser{fns: fns, invert: true}
}

// AlwaysParser always returns the value associated with it.
type AlwaysParser struct {
	value interface{}
}

// Always always returns the given value.
func Always(value interface{}) Parser {
	return &AlwaysParser{value: value}
}

// Parse parses the input.
func (p *AlwaysParser) Parse(sc scanner.Scanner) result.ParseResult {
	return result.Success(textpos.Single(sc.GetPos()), p.value)
}

// NextParser calls an inner parser, and if it succeeds, passes the result to
// the given function, which determines the next parser to run.
type NextParser struct {
	inner Parser
	fn    func(interface{}) Parser
}

// Next runs an inner parser, and if it succeeds, passes the result to the
// given function, which determines the next parser to run.
func Next(in Parser, next func(interface{}) Parser) Parser {
	return &NextParser{
		inner: in,
		fn:    next,
	}
}

// Parse parses the input.
func (p *NextParser) Parse(sc scanner.Scanner) result.ParseResult {
	result := p.inner.Parse(sc)
	if !result.Matched() {
		return result
	}

	next := p.fn(result.Result())
	return next.Parse(sc)
}

// LongestParser runs all the parsers associated with it and returns the one
// with the longest match.
type LongestParser struct {
	parsers []Parser
}

// Longest runs each of the given parsers in order and returns the result from
// the one with the longest match.
func Longest(parsers ...Parser) Parser {
	return &LongestParser{parsers}
}

// Parse parses the input.
func (p *LongestParser) Parse(sc scanner.Scanner) result.ParseResult {
	if len(p.parsers) == 0 {
		return fail(sc.GetPos(), "no parser matched")
	} else if len(p.parsers) == 1 {
		return p.parsers[0].Parse(sc)
	}

	res, _ := p.one(sc, p.parsers, -1)

	if res == nil {
		return fail(sc.GetPos(), "no parser matched")
	}

	// The top snapshot is the winning end position. We'll rewind to it so the
	// next parser is in the right place.
	// S: []
	sc.RewindSnapshot()

	return res
}

func (p *LongestParser) one(sc scanner.Scanner, remaining []Parser, max int) (result.ParseResult, int) {
	if len(remaining) == 0 {
		// S: []
		return nil, -1
	}

	inner := remaining[0]
	rest := remaining[1:]

	// Create a snapshot for this run.
	// S: [Initial]
	sc.StartSnapshot()

	// Run this parser.
	res := inner.Parse(sc)
	if !res.Matched() {
		// If it didn't match we just try the next one.
		// S: []
		sc.RewindSnapshot()
		return p.one(sc, rest, max)
	}

	nlen := res.TextRange().Length()
	if nlen > max {
		max = nlen
	}

	// Create a new snapshot at the result location.
	// S: [Initial, OurMatch]
	sc.StartSnapshot()

	// Swap back to the start position to run the next parser.
	// S: [OurMatch, Initial]
	sc.SwapSnapshot()

	// S: [OurMatch]
	sc.RewindSnapshot()

	// If next is nil, S: [OurMatch]
	// Otherwise, S: [OurMatch, TheirMatch]
	next, nlen := p.one(sc, rest, max)
	if nlen > max {
		// They win. Remove our snapshot.
		// S: [TheirMatch, OurMatch]
		sc.SwapSnapshot()
		res, max = next, nlen
	}

	if next != nil {
		// If we got a result, we remove whatever the wrong one is from the
		// top of the stack.
		// S: [LongestMatch]
		sc.PopSnapshot()
	}

	return res, max
}
