package cursor

// The Cursor struct stores the location of the cursor in the buffer
// as well as the selection
type Cursor struct {
	Loc

	// Last cursor x position
	LastVisualX int

	// The current selection as a range of character numbers (inclusive)
	CurSelection [2]Loc
}

func NewCursor(l Loc) *Cursor {
	c := &Cursor{
		Loc: l,
	}
	return c
}

// Goto puts the cursor at the given cursor's location and gives
// the current cursor its selection too
func (c *Cursor) Goto(b Cursor) {
	c.X, c.Y, c.LastVisualX = b.X, b.Y, b.LastVisualX
	c.CurSelection = b.CurSelection
}

// SetSelectionStart sets the start of the selection
func (c *Cursor) SetSelectionStart(pos Loc) {
	c.CurSelection[0] = pos
}

// SetSelectionEnd sets the end of the selection
func (c *Cursor) SetSelectionEnd(pos Loc) {
	c.CurSelection[1] = pos
}

// HasSelection returns whether or not the user has selected anything
func (c *Cursor) HasSelection() bool {
	return c.CurSelection[0] != c.CurSelection[1]
}

// SelectTo selects from the current cursor location to the given
// location
func (c *Cursor) SelectTo(loc Loc) {
	if loc.GreaterThan(c.Loc) {
		c.SetSelectionStart(c.Loc)
		c.SetSelectionEnd(loc)
	} else {
		c.SetSelectionStart(loc)
		c.SetSelectionEnd(c.Loc)
	}
}
