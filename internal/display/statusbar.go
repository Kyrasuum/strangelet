package display

import (
	"github.com/Kyrasuum/cview"
	"github.com/gdamore/tcell/v2"
)

var (
	statusBarBGColor tcell.Color
	statusBarFGColor tcell.Color
)

// Statusbar display struct
type statusBar struct {
	*cview.TextView
}

func NewStatusBar(flex *cview.Flex) (stat *statusBar) {
	stat = &statusBar{}

	if statusBarBGColor == 0 {
		statusBarBGColor = tcell.NewRGBColor(160, 160, 160)
	}
	if statusBarFGColor == 0 {
		statusBarFGColor = tcell.NewRGBColor(20, 20, 20)
	}

	stat.TextView = cview.NewTextView()
	stat.TextView.SetTextAlign(cview.AlignLeft)
	stat.TextView.SetText("Status Bar")
	stat.TextView.Box.SetBackgroundColor(statusBarBGColor)
	stat.TextView.SetTextColor(statusBarFGColor)

	flex.AddItem(stat, 1, 1, false)

	return stat
}
