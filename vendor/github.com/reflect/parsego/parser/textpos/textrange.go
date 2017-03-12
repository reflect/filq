/*
Package textpos contains immutable structures for working with
positions in a text document.
*/
package textpos

// TextPos is a single character position in some text. Both line and
// col start from 0.
//
// Immutable data structures are somewhat inconvenient to write in Go.
type TextPos struct {
	offset int
	line   int
	col    int
}

// Offset returns the character offset in the stream.
func (t TextPos) Offset() int {
	return t.offset
}

// Line returns the line number, starting at 0
func (t TextPos) Line() int {
	return t.line
}

// Col returns the column number within the line, starting at 0
func (t TextPos) Col() int {
	return t.col
}

// TextRange is an (inclusive) range between two TextPos.
type TextRange struct {
	start TextPos
	end   TextPos
}

// Range constructs a new TextRange
func Range(start, end TextPos) TextRange {
	return TextRange{start, end}
}

// Single returns a zero-character range.
func Single(pos TextPos) TextRange {
	return TextRange{pos, pos}
}

// Start returns the position of the first character in the range.
func (t TextRange) Start() TextPos {
	return t.start
}

// End returns the position of the last character in the range.
func (t TextRange) End() TextPos {
	return t.end
}

// Length returns the length of this range.
func (t TextRange) Length() int {
	return t.end.Offset() - t.start.Offset()
}

// StartingPos returns the 0 position.
func StartingPos() TextPos {
	return TextPos{
		line: 0,
		col:  0,
	}
}

// Pos is a shorthand for creating a TextPos.
func Pos(offset, line, col int) TextPos {
	return TextPos{offset, line, col}
}

// AdvanceCol return a new TextPos with the column advanced by one.
func (t TextPos) AdvanceCol() TextPos {
	return TextPos{
		offset: t.offset + 1,
		col:    t.col + 1,
		line:   t.line,
	}
}

// AdvanceLine returns a new TextPos with the line advanced by one.
func (t TextPos) AdvanceLine() TextPos {
	return TextPos{
		offset: t.offset + 1,
		col:    0,
		line:   t.line + 1,
	}
}

// Advance returns a new TextPos advanced by the given character.
func (t TextPos) Advance(c rune) TextPos {
	if c == '\n' {
		return t.AdvanceLine()
	}
	return t.AdvanceCol()
}
