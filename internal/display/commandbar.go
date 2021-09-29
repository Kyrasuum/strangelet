package display

import (
	"github.com/Kyrasuum/cview"
	"github.com/gdamore/tcell/v2"
)

var (
	commandBarBGColor tcell.Color
	commandBarFGColor tcell.Color
)

type commandBar struct {
	*cview.TextView
}

func NewCommandBar(flex *cview.Flex) (cbar *commandBar) {
	cbar = &commandBar{}

	if commandBarBGColor == 0 {
		commandBarBGColor = tcell.NewRGBColor(20, 20, 20)
	}
	if commandBarFGColor == 0 {
		commandBarFGColor = tcell.NewRGBColor(170, 170, 170)
	}

	cbar.TextView = cview.NewTextView()
	cbar.TextView.SetTextAlign(cview.AlignLeft)
	cbar.TextView.SetText("Commands Bar")
	cbar.TextView.Box.SetBackgroundColor(commandBarBGColor)
	cbar.TextView.SetTextColor(commandBarFGColor)

	flex.AddItem(cbar.TextView, 1, 1, false)

	return cbar
}
