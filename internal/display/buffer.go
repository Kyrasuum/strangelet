package display

import (
	"math"

	"github.com/Kyrasuum/cview"
	"github.com/gdamore/tcell/v2"
)

var (
	bufferFGColor tcell.Color
	bufferBGColor tcell.Color
)

// Buffer display struct
type buffer struct {
	*cview.TextView
}

func NewBuffer(subFlex *cview.Flex) (buff *buffer) {
	buff = &buffer{}

	if bufferBGColor == 0 {
		bufferBGColor = tcell.NewRGBColor(30, 30, 30)
	}
	if bufferFGColor == 0 {
		bufferFGColor = tcell.NewRGBColor(230, 230, 230)
	}

	buff.TextView = cview.NewTextView()
	buff.TextView.SetTextAlign(cview.AlignLeft)
	buff.TextView.SetRegions(true)
	buff.TextView.SetText("")
	buff.TextView.Box.SetBackgroundColor(bufferBGColor)
	buff.TextView.SetTextColor(bufferFGColor)
	buff.TextView.SetHighlightBackgroundColor(gutterBGColor)
	buff.TextView.SetHighlightForegroundColor(bufferFGColor)

	subFlex.AddItem(buff, 0, 1, false)

	return buff
}

func (dbuffer *buffer) HandleMouse(event *tcell.EventMouse, action cview.MouseAction) (*tcell.EventMouse, cview.MouseAction) {
	switch action {
	case cview.MouseLeftClick:
	case cview.MouseLeftDoubleClick:
	case cview.MouseRightClick:
	case cview.MouseRightDoubleClick:
	case cview.MouseMiddleClick:
	case cview.MouseMiddleDoubleClick:
	case cview.MouseScrollLeft:
		row, column := dbuffer.GetScrollOffset()
		dbuffer.ScrollTo(row, int(math.Max(0, float64(column-1))))
		return nil, 0
	case cview.MouseScrollRight:
		row, column := dbuffer.GetScrollOffset()
		dbuffer.ScrollTo(row, int(math.Max(0, float64(column+1))))
		return nil, 0
	case cview.MouseScrollUp:
		row, column := dbuffer.GetScrollOffset()
		dbuffer.ScrollTo(int(math.Max(0, float64(row-1))), column)
		return nil, 0
	case cview.MouseScrollDown:
		row, column := dbuffer.GetScrollOffset()
		dbuffer.ScrollTo(int(math.Max(0, float64(row+1))), column)
		return nil, 0
	}
	return event, action
}

func (dbuffer *buffer) HandleInput(event *tcell.EventKey) *tcell.EventKey {
	return event
}
