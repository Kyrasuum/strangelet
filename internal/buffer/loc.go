package buffer

import (
	"strangelet/internal/cursor"
	"strangelet/internal/util"
)

// The following functions require a buffer to know where newlines are

// Diff returns the distance between two locations
func DiffLA(a, b cursor.Loc, buf *LineArray) int {
	if a.Y == b.Y {
		if a.X > b.X {
			return a.X - b.X
		}
		return b.X - a.X
	}

	// Make sure a is guaranteed to be less than b
	if b.LessThan(a) {
		a, b = b, a
	}

	loc := 0
	for i := a.Y + 1; i < b.Y; i++ {
		// + 1 for the newline
		loc += util.CharacterCount(buf.LineBytes(i)) + 1
	}
	loc += util.CharacterCount(buf.LineBytes(a.Y)) - a.X + b.X + 1
	return loc
}

// This moves the location one character to the right
func right(l cursor.Loc, buf *LineArray) cursor.Loc {
	if l == buf.End() {
		return cursor.Loc{l.X + 1, l.Y}
	}
	var res cursor.Loc
	if l.X < util.CharacterCount(buf.LineBytes(l.Y)) {
		res = cursor.Loc{l.X + 1, l.Y}
	} else {
		res = cursor.Loc{0, l.Y + 1}
	}
	return res
}

// This moves the given location one character to the left
func left(l cursor.Loc, buf *LineArray) cursor.Loc {
	if l == buf.Start() {
		return cursor.Loc{l.X - 1, l.Y}
	}
	var res cursor.Loc
	if l.X > 0 {
		res = cursor.Loc{l.X - 1, l.Y}
	} else {
		res = cursor.Loc{util.CharacterCount(buf.LineBytes(l.Y - 1)), l.Y - 1}
	}
	return res
}

// MoveLA moves the cursor n characters to the left or right
// It moves the cursor left if n is negative
func MoveLA(l cursor.Loc, n int, buf *LineArray) cursor.Loc {
	if n > 0 {
		for i := 0; i < n; i++ {
			l = right(l, buf)
		}
		return l
	}
	for i := 0; i < util.Abs(n); i++ {
		l = left(l, buf)
	}
	return l
}

// Diff returns the difference between two locs
func Diff(l cursor.Loc, b cursor.Loc, buf *Buffer) int {
	return DiffLA(l, b, buf.LineArray)
}

// Move moves a loc n characters
func Move(l cursor.Loc, n int, buf *Buffer) cursor.Loc {
	return MoveLA(l, n, buf.LineArray)
}

// ByteOffset is just like ToCharPos except it counts bytes instead of runes
func ByteOffset(pos cursor.Loc, buf *Buffer) int {
	x, y := pos.X, pos.Y
	loc := 0
	for i := 0; i < y; i++ {
		// + 1 for the newline
		loc += len(buf.Line(i)) + 1
	}
	loc += len(buf.Line(y)[:x])
	return loc
}

// clamps a loc within a buffer
func clamp(pos cursor.Loc, la *LineArray) cursor.Loc {
	if pos.GreaterEqual(la.End()) {
		return la.End()
	} else if pos.LessThan(la.Start()) {
		return la.Start()
	}
	return pos
}
