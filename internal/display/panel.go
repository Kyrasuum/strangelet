package display

import (
	"github.com/Kyrasuum/cview"
	"github.com/gdamore/tcell/v2"
)

var (
	panelTabBGColor  tcell.Color
	panelTabFGColor  tcell.Color
	panelTabBGColorF tcell.Color
	panelTabFGColorF tcell.Color
)

type panel struct {
	*cview.TabbedPanels
	tabs map[string]interface{}
}

func NewPanel(paneflex *cview.Flex, index int) (pan *panel) {
	pan = &panel{}

	//init colors
	if panelTabBGColor == 0 {
		panelTabBGColor = tcell.NewRGBColor(200, 0, 0)
	}
	if panelTabFGColor == 0 {
		panelTabFGColor = tcell.NewRGBColor(220, 220, 220)
	}
	if panelTabBGColorF == 0 {
		panelTabBGColorF = tcell.NewRGBColor(200, 200, 200)
	}
	if panelTabFGColorF == 0 {
		panelTabFGColorF = tcell.NewRGBColor(20, 20, 20)
	}

	//setup map
	pan.tabs = make(map[string]interface{})

	//set colors
	pan.TabbedPanels = cview.NewTabbedPanels()
	pan.TabbedPanels.SetTabBackgroundColor(panelTabBGColor)
	pan.TabbedPanels.SetTabBackgroundColorFocused(panelTabBGColorF)
	pan.TabbedPanels.SetTabTextColor(panelTabFGColor)
	pan.TabbedPanels.SetTabTextColorFocused(panelTabFGColorF)
	pan.TabbedPanels.SetTabSwitcherHeight(1)

	//default tabs
	for panelIndex := 0; panelIndex < 20; panelIndex++ {
		pan.AddTab(panelIndex)
	}

	//add to paneflex
	paneflex.AddItem(pan.TabbedPanels, 0, 1, false)

	return pan
}

func (panel *panel) AddTab(tabIndex int) {
	t := NewTab(panel.TabbedPanels, tabIndex)
	panel.tabs[t.GetName()] = t
}

func (panel *panel) GetCurrentTab() (t *tab) {
	index := panel.TabbedPanels.GetCurrentTab()
	elem := panel.tabs[index].(tab)
	return &elem
}

func (panel *panel) HandleInput(tevent *tcell.EventKey) (retEvent *tcell.EventKey) {
	retEvent = panel.GetCurrentTab().HandleInput(tevent)
	if retEvent != tevent {
		return retEvent
	}
	return tevent
}
