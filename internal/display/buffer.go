package display

import (
	"github.com/Kyrasuum/cview"
	"github.com/gdamore/tcell/v2"
)

var ()

type buffer struct {
	*cview.TextView
}

func NewBuffer(subFlex *cview.Flex) (buff *buffer) {
	buff = &buffer{}

	buff.TextView = cview.NewTextView()
	buff.TextView.SetTextAlign(cview.AlignLeft)
	buff.TextView.SetText("Buffer Content")
	buff.TextView.Box.SetBackgroundColor(tcell.NewRGBColor(20, 20, 20))
	buff.TextView.SetTextColor(tcell.NewRGBColor(230, 230, 230))

	subFlex.AddItem(buff.TextView, 0, 1, false)

	return buff
}
