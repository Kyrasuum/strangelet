package buffer

import (
	"strangelet/internal/cursor"
	"strangelet/internal/util"
)

// // SetCursors resets this buffer's cursors to a new list
// func (b *Buffer) SetCursors(c []*Cursor) {
// b.cursors = c
// b.EventHandler.cursors = b.cursors
// b.EventHandler.active = b.curCursor
// }
//
// // AddCursor adds a new cursor to the list
// func (b *Buffer) AddCursor(c *Cursor) {
// b.cursors = append(b.cursors, c)
// b.EventHandler.cursors = b.cursors
// b.EventHandler.active = b.curCursor
// b.UpdateCursors()
// }
//
// // SetCurCursor sets the current cursor
// func (b *Buffer) SetCurCursor(n int) {
// b.curCursor = n
// }
//
// // GetActiveCursor returns the main cursor in this buffer
// func (b *Buffer) GetActiveCursor() *Cursor {
// return b.cursors[b.curCursor]
// }
//
// // GetCursor returns the nth cursor
// func (b *Buffer) GetCursor(n int) *Cursor {
// return b.cursors[n]
// }
//
// // GetCursors returns the list of cursors in this buffer
// func (b *Buffer) GetCursors() []*Cursor {
// return b.cursors
// }
//
// // NumCursors returns the number of cursors
// func (b *Buffer) NumCursors() int {
// return len(b.cursors)
// }
//
// // MergeCursors merges any cursors that are at the same position
// // into one cursor
// func (b *Buffer) MergeCursors() {
// var cursors []*Cursor
// for i := 0; i < len(b.cursors); i++ {
// c1 := b.cursors[i]
// if c1 != nil {
// for j := 0; j < len(b.cursors); j++ {
// c2 := b.cursors[j]
// if c2 != nil && i != j && c1.Loc == c2.Loc {
// b.cursors[j] = nil
// }
// }
// cursors = append(cursors, c1)
// }
// }
//
// b.cursors = cursors
//
// for i := range b.cursors {
// b.cursors[i].Num = i
// }
//
// if b.curCursor >= len(b.cursors) {
// b.curCursor = len(b.cursors) - 1
// }
// b.EventHandler.cursors = b.cursors
// b.EventHandler.active = b.curCursor
// }
//
// // UpdateCursors updates all the cursors indicies
// func (b *Buffer) UpdateCursors() {
// b.EventHandler.cursors = b.cursors
// b.EventHandler.active = b.curCursor
// for i, c := range b.cursors {
// c.Num = i
// }
// }
//
// func (b *Buffer) RemoveCursor(i int) {
// copy(b.cursors[i:], b.cursors[i+1:])
// b.cursors[len(b.cursors)-1] = nil
// b.cursors = b.cursors[:len(b.cursors)-1]
// b.curCursor = util.Clamp(b.curCursor, 0, len(b.cursors)-1)
// b.UpdateCursors()
// }
//
// // ClearCursors removes all extra cursors
// func (b *Buffer) ClearCursors() {
// for i := 1; i < len(b.cursors); i++ {
// b.cursors[i] = nil
// }
// b.cursors = b.cursors[:1]
// b.UpdateCursors()
// b.curCursor = 0
// b.GetActiveCursor().ResetSelection()
// }
//
// // RelocateCursors relocates all cursors (makes sure they are in the buffer)
// func (b *Buffer) RelocateCursors() {
// for _, c := range b.cursors {
// c.Relocate()
// }
// }
//
// // CopySelection copies the user's selection to either "primary"
// // or "clipboard"
// func CopySelection(c *cursor.Cursor, target clipboard.Register) {
// if c.HasSelection() {
// if target != clipboard.PrimaryReg || c.buf.Settings["useprimary"].(bool) {
// clipboard.WriteMulti(string(c.GetSelection()), target, c.Num, c.buf.NumCursors())
// }
// }
// }

// InBounds returns whether the given location is a valid character position in the given buffer
func InBounds(pos cursor.Loc, buf *Buffer) bool {
	if pos.Y < 0 || pos.Y >= len(buf.lines) || pos.X < 0 || pos.X > util.CharacterCount(buf.LineBytes(pos.Y)) {
		return false
	}

	return true
}

// GetCharPosInLine gets the char position of a visual x y
// coordinate (this is necessary because tabs are 1 char but
// 4 visual spaces)
func GetCharPosInLine(c *cursor.Cursor, b []byte, visualPos int, buf *Buffer) int {
	tabsize := int(buf.Settings["tabsize"].(float64))
	return util.GetCharPosInLine(b, visualPos, tabsize)
}

// StartOfText moves the cursor to the first non-whitespace rune of
// the line it is on
func StartOfText(c *cursor.Cursor, buf *Buffer) {
	Start(c, buf)
	for util.IsWhitespace(RuneUnder(c, c.X, buf)) {
		if c.X == util.CharacterCount(buf.LineBytes(c.Y)) {
			break
		}
		Right(c, buf)
	}
}

// IsStartOfText returns whether the cursor is at the first
// non-whitespace rune of the line it is on
func IsStartOfText(c *cursor.Cursor, buf *Buffer) bool {
	x := 0
	for util.IsWhitespace(RuneUnder(c, x, buf)) {
		if x == util.CharacterCount(buf.LineBytes(c.Y)) {
			break
		}
		x++
	}
	return c.X == x
}

// End moves the cursor to the end of the line it is on
func End(c *cursor.Cursor, buf *Buffer) {
	c.X = util.CharacterCount(buf.LineBytes(c.Y))
	c.LastVisualX = GetVisualX(c, buf)
}

// Start moves the cursor to the start of the line it is on
func Start(c *cursor.Cursor, buf *Buffer) {
	c.X = 0
	c.LastVisualX = GetVisualX(c, buf)
}

func StoreVisualX(c *cursor.Cursor, buf *Buffer) {
	c.LastVisualX = GetVisualX(c, buf)
}

// DeleteSelection deletes the currently selected text
func DeleteSelection(c *cursor.Cursor, buf *Buffer) {
	if c.CurSelection[0].GreaterThan(c.CurSelection[1]) {
		buf.Remove(c.CurSelection[1], c.CurSelection[0])
		c.Loc = c.CurSelection[1]
	} else if !c.HasSelection() {
		return
	} else {
		buf.Remove(c.CurSelection[0], c.CurSelection[1])
		c.Loc = c.CurSelection[0]
	}
}

// GetSelection returns the cursor's selection
func GetSelection(c *cursor.Cursor, buf *Buffer) []byte {
	if InBounds(c.CurSelection[0], buf) && InBounds(c.CurSelection[1], buf) {
		if c.CurSelection[0].GreaterThan(c.CurSelection[1]) {
			return buf.Substr(c.CurSelection[1], c.CurSelection[0])
		}
		return buf.Substr(c.CurSelection[0], c.CurSelection[1])
	}
	return []byte{}
}

// SelectLine selects the current line
func SelectLine(c *cursor.Cursor, buf *Buffer) {
	Start(c, buf)
	c.SetSelectionStart(c.Loc)
	End(c, buf)
	if len(buf.lines)-1 > c.Y {
		c.SetSelectionEnd(Move(c.Loc, 1, buf))
	} else {
		c.SetSelectionEnd(c.Loc)
	}
}

// GetVisualX returns the x value of the cursor in visual spaces
func GetVisualX(c *cursor.Cursor, buf *Buffer) int {
	if c.X <= 0 {
		c.X = 0
		return 0
	}

	bytes := buf.LineBytes(c.Y)
	tabsize := int(buf.Settings["tabsize"].(float64))

	return util.StringWidth(bytes, c.X, tabsize)
}

// AddLineToSelection adds the current line to the selection
func AddLineToSelection(c *cursor.Cursor, buf *Buffer) {
	if c.Loc.LessThan(c.CurSelection[0]) {
		Start(c, buf)
		c.SetSelectionStart(c.Loc)
		c.SetSelectionEnd(c.CurSelection[1])
	}
	if c.Loc.GreaterThan(c.CurSelection[1]) {
		End(c, buf)
		c.SetSelectionEnd(Move(c.Loc, 1, buf))
		c.SetSelectionStart(c.CurSelection[0])
	}
}

// UpN moves the cursor up N lines (if possible)
func UpN(c *cursor.Cursor, amount int, buf *Buffer) {
	proposedY := c.Y - amount
	if proposedY < 0 {
		proposedY = 0
	} else if proposedY >= len(buf.lines) {
		proposedY = len(buf.lines) - 1
	}

	bytes := buf.LineBytes(proposedY)
	c.X = GetCharPosInLine(c, bytes, c.LastVisualX, buf)

	if c.X > util.CharacterCount(bytes) || (amount < 0 && proposedY == c.Y) {
		c.X = util.CharacterCount(bytes)
		StoreVisualX(c, buf)
	}

	if c.X < 0 || (amount > 0 && proposedY == c.Y) {
		c.X = 0
		StoreVisualX(c, buf)
	}

	c.Y = proposedY
}

// DownN moves the cursor down N lines (if possible)
func DownN(c *cursor.Cursor, amount int, buf *Buffer) {
	UpN(c, -amount, buf)
}

// Up moves the cursor up one line (if possible)
func Up(c *cursor.Cursor, buf *Buffer) {
	UpN(c, 1, buf)
}

// Down moves the cursor down one line (if possible)
func Down(c *cursor.Cursor, buf *Buffer) {
	DownN(c, 1, buf)
}

// Left moves the cursor left one cell (if possible) or to
// the previous line if it is at the beginning
func Left(c *cursor.Cursor, buf *Buffer) {
	if c.Loc == buf.Start() {
		return
	}
	if c.X > 0 {
		c.X--
	} else {
		Up(c, buf)
		End(c, buf)
	}
	StoreVisualX(c, buf)
}

// Right moves the cursor right one cell (if possible) or
// to the next line if it is at the end
func Right(c *cursor.Cursor, buf *Buffer) {
	if c.Loc == buf.End() {
		return
	}
	if c.X < util.CharacterCount(buf.LineBytes(c.Y)) {
		c.X++
	} else {
		Down(c, buf)
		Start(c, buf)
	}
	StoreVisualX(c, buf)
}

// GotoLoc puts the cursor at the given cursor's location and gives
// the current cursor its selection too
func GotoLoc(c *cursor.Cursor, l cursor.Loc, buf *Buffer) {
	c.X, c.Y = l.X, l.Y
	StoreVisualX(c, buf)
}

// Relocate makes sure that the cursor is inside the bounds
// of the buffer If it isn't, it moves it to be within the
// buffer's lines
func Relocate(c *cursor.Cursor, buf *SharedBuffer) {
	if c.Y < 0 {
		c.Y = 0
	} else if c.Y >= len(buf.lines) {
		c.Y = len(buf.lines) - 1
	}

	if c.X < 0 {
		c.X = 0
	} else if c.X > util.CharacterCount(buf.LineBytes(c.Y)) {
		c.X = util.CharacterCount(buf.LineBytes(c.Y))
	}
}

// SelectWord selects the word the cursor is currently on
func SelectWord(c *cursor.Cursor, buf *Buffer) {
	if len(buf.LineBytes(c.Y)) == 0 {
		return
	}

	if !util.IsWordChar(RuneUnder(c, c.X, buf)) {
		c.SetSelectionStart(c.Loc)
		c.SetSelectionEnd(Move(c.Loc, 1, buf))
		return
	}

	forward, backward := c.X, c.X

	for backward > 0 && util.IsWordChar(RuneUnder(c, backward-1, buf)) {
		backward--
	}

	c.SetSelectionStart(cursor.Loc{backward, c.Y})

	lineLen := util.CharacterCount(buf.LineBytes(c.Y)) - 1
	for forward < lineLen && util.IsWordChar(RuneUnder(c, forward+1, buf)) {
		forward++
	}

	c.SetSelectionEnd(Move(cursor.Loc{forward, c.Y}, 1, buf))
	c.Loc = c.CurSelection[1]
}

// AddWordToSelection adds the word the cursor is currently on
// to the selection
func AddWordToSelection(c *cursor.Cursor, buf *Buffer) {
	if c.Loc.GreaterThan(c.CurSelection[0]) && c.Loc.LessThan(c.CurSelection[1]) {
		return
	}

	if c.Loc.LessThan(c.CurSelection[0]) {
		backward := c.X

		for backward > 0 && util.IsWordChar(RuneUnder(c, backward-1, buf)) {
			backward--
		}

		c.SetSelectionStart(cursor.Loc{backward, c.Y})
		c.SetSelectionEnd(c.CurSelection[1])
	}

	if c.Loc.GreaterThan(c.CurSelection[1]) {
		forward := c.X

		lineLen := util.CharacterCount(buf.LineBytes(c.Y)) - 1
		for forward < lineLen && util.IsWordChar(RuneUnder(c, forward+1, buf)) {
			forward++
		}

		c.SetSelectionEnd(Move(cursor.Loc{forward, c.Y}, 1, buf))
		c.SetSelectionStart(c.CurSelection[0])
	}

	c.Loc = c.CurSelection[1]
}

// SelectTo selects from the current cursor location to the given
// location
func SelectTo(c *cursor.Cursor, loc cursor.Loc) {
	if loc.GreaterThan(c.CurSelection[0]) {
		c.SetSelectionStart(c.CurSelection[0])
		c.SetSelectionEnd(loc)
	} else {
		c.SetSelectionStart(loc)
		c.SetSelectionEnd(c.CurSelection[0])
	}
}

// ResetSelection resets the user's selection
func ResetSelection(c *cursor.Cursor, buf *Buffer) {
	c.CurSelection[0] = buf.Start()
	c.CurSelection[1] = buf.Start()
}

// Deselect closes the cursor's current selection
// Start indicates whether the cursor should be placed
// at the start or end of the selection
func Deselect(c *cursor.Cursor, buf *Buffer, start bool) {
	if c.HasSelection() {
		if start {
			c.Loc = c.CurSelection[0]
		} else {
			c.Loc = Move(c.CurSelection[1], -1, buf)
		}
		ResetSelection(c, buf)
		StoreVisualX(c, buf)
	}
}

// WordRight moves the cursor one word to the right
func WordRight(c *cursor.Cursor, buf *Buffer) {
	for util.IsWhitespace(RuneUnder(c, c.X, buf)) {
		if c.X == util.CharacterCount(buf.LineBytes(c.Y)) {
			Right(c, buf)
			return
		}
		Right(c, buf)
	}
	Right(c, buf)
	for util.IsWordChar(RuneUnder(c, c.X, buf)) {
		if c.X == util.CharacterCount(buf.LineBytes(c.Y)) {
			return
		}
		Right(c, buf)
	}
}

// WordLeft moves the cursor one word to the left
func WordLeft(c *cursor.Cursor, buf *Buffer) {
	Left(c, buf)
	for util.IsWhitespace(RuneUnder(c, c.X, buf)) {
		if c.X == 0 {
			return
		}
		Left(c, buf)
	}
	Left(c, buf)
	for util.IsWordChar(RuneUnder(c, c.X, buf)) {
		if c.X == 0 {
			return
		}
		Left(c, buf)
	}
	Right(c, buf)
}

// RuneUnder returns the rune under the given x position
func RuneUnder(c *cursor.Cursor, x int, buf *Buffer) rune {
	line := buf.LineBytes(c.Y)
	if len(line) == 0 || x >= util.CharacterCount(line) {
		return '\n'
	} else if x < 0 {
		x = 0
	}
	i := 0
	for len(line) > 0 {
		r, _, size := util.DecodeCharacter(line)
		line = line[size:]

		if i == x {
			return r
		}

		i++
	}
	return '\n'
}
