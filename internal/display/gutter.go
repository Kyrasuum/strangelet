package display

import (
	"github.com/Kyrasuum/cview"
	"github.com/gdamore/tcell/v2"
)

var (
	gutterBGColor tcell.Color
	gutterFGColor tcell.Color
)

type Gutter struct {
	*cview.TextView
}

func (gutt *Gutter) InitGutter(subFlex *cview.Flex) {
	if gutterBGColor == 0 {
		gutterBGColor = tcell.NewRGBColor(40, 40, 40)
	}
	if gutterFGColor == 0 {
		gutterFGColor = tcell.NewRGBColor(170, 170, 170)
	}
	gutt.TextView = cview.NewTextView()
	gutt.TextView.SetTextAlign(cview.AlignRight)
	gutt.TextView.SetText("1 ")
	gutt.TextView.Box.SetBackgroundColor(gutterBGColor)
	gutt.TextView.SetTextColor(gutterFGColor)

	subFlex.AddItem(gutt.TextView, 2, 1, false)
}
