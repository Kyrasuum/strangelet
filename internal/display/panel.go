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

type Panel struct {
	*cview.TabbedPanels
	tabs map[string]interface{}
}

func (panel *Panel) InitPanel(paneflex *cview.Flex, index int) {
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
	panel.tabs = make(map[string]interface{})

	//set colors
	panel.TabbedPanels = cview.NewTabbedPanels()
	panel.TabbedPanels.SetTabBackgroundColor(panelTabBGColor)
	panel.TabbedPanels.SetTabBackgroundColorFocused(panelTabBGColorF)
	panel.TabbedPanels.SetTabTextColor(panelTabFGColor)
	panel.TabbedPanels.SetTabTextColorFocused(panelTabFGColorF)

	//default tabs
	for panelIndex := 0; panelIndex < 5; panelIndex++ {
		panel.AddTab(panelIndex)
	}

	//add to paneflex
	paneflex.AddItem(panel.TabbedPanels, 0, 1, false)
}

func (panel *Panel) AddTab(tabIndex int) {
	var tab Tab
	tab.InitTab(panel.TabbedPanels, tabIndex)
	panel.tabs[tab.GetName()] = tab
}

func (panel *Panel) GetCurrentTab() (tab *Tab) {
	index := panel.TabbedPanels.GetCurrentTab()
	elem := panel.tabs[index].(Tab)
	return &elem
}

func (panel *Panel) HandleInput(tevent *tcell.EventKey) (retEvent *tcell.EventKey) {
	retEvent = panel.GetCurrentTab().HandleInput(tevent)
	if retEvent != tevent {
		return retEvent
	}
	return tevent
}
