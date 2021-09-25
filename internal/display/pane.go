package display

import (
	"github.com/Kyrasuum/cview"
	"github.com/gdamore/tcell/v2"
)

var (
	paneTabBGColor  tcell.Color
	paneTabFGColor  tcell.Color
	paneTabBGColorF tcell.Color
	paneTabFGColorF tcell.Color
)

type Pane struct {
	*cview.TabbedPanels
	panels []Panel
}

func (pane *Pane) InitPane(subFlex *cview.Flex) {
	if paneTabBGColor == 0 {
		paneTabBGColor = tcell.NewRGBColor(200, 0, 0)
	}
	if paneTabFGColor == 0 {
		paneTabFGColor = tcell.NewRGBColor(220, 220, 220)
	}
	if paneTabBGColorF == 0 {
		paneTabBGColorF = tcell.NewRGBColor(200, 200, 200)
	}
	if paneTabFGColorF == 0 {
		paneTabFGColorF = tcell.NewRGBColor(20, 20, 20)
	}

	pane.TabbedPanels = cview.NewTabbedPanels()
	pane.TabbedPanels.SetTabBackgroundColor(paneTabBGColor)
	pane.TabbedPanels.SetTabBackgroundColorFocused(paneTabBGColorF)
	pane.TabbedPanels.SetTabTextColor(paneTabFGColor)
	pane.TabbedPanels.SetTabTextColorFocused(paneTabFGColorF)

	for panelIndex := 0; panelIndex < 5; panelIndex++ {
		pane.AddPanel(panelIndex)
	}
	subFlex.AddItem(pane, 0, 1, false)
}

func (pane *Pane) AddPanel(panelIndex int) {
	var panel Panel
	panel.InitPanel(pane.TabbedPanels, panelIndex)
	pane.panels = append(pane.panels, panel)
}
