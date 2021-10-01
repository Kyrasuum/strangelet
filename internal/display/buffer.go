package display

import (
	buff "strangelet/internal/buffer"
	"strangelet/internal/cursor"

	"github.com/Kyrasuum/cview"
	"github.com/gdamore/tcell/v2"
)

var ()

type buffer struct {
	*cview.TextView
	*buff.Buffer

	cursors     []*cursor.Cursor
	curCursor   int
	StartCursor cursor.Loc
}

func NewBuffer(subFlex *cview.Flex, b *buff.Buffer) (buff *buffer) {
	buff = &buffer{}
	buff.Buffer = b

	buff.TextView = cview.NewTextView()
	buff.TextView.SetTextAlign(cview.AlignLeft)
	buff.TextView.SetText("")
	buff.TextView.Box.SetBackgroundColor(tcell.NewRGBColor(20, 20, 20))
	buff.TextView.SetTextColor(tcell.NewRGBColor(230, 230, 230))

	subFlex.AddItem(buff.TextView, 0, 1, false)

	return buff
}
