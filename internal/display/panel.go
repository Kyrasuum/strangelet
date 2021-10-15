package display

import (
	buff "strangelet/internal/buffer"
	"strangelet/pkg/app"

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

func NewPanel(paneflex *cview.Flex) (pan *panel) {
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

	pan.AddEmptyTab()

	//add to paneflex
	paneflex.AddItem(pan.TabbedPanels, 0, 1, false)
	app.CurApp.SetFocus(pan)

	return pan
}

func (panel *panel) Render(scr tcell.Screen) bool {
	if panel.GetCurrentTab().Render(scr) {
		return true
	}

	return false
}

func (panel *panel) AddTab(b *buff.Buffer) {
	if !panel.HasNonEmptyTab() {
		panel.RemoveTabByIndex(0)
	}
	t := NewTab(panel.TabbedPanels, b)
	panel.tabs[t.GetName()] = t
	b.MarkModified(0, 0)
}

func (panel *panel) RemoveTabByName(name string) {
	panel.TabbedPanels.RemoveTab(name)
	panel.tabs[name] = nil
}

func (panel *panel) RemoveTabByIndex(ind int) {
	name := panel.TabbedPanels.TabName(ind)
	panel.RemoveTabByName(name)
}

func (panel *panel) AddEmptyTab() {
	b := buff.NewBufferFromString("", "", buff.BTDefault)
	t := NewTab(panel.TabbedPanels, b)
	panel.tabs[t.GetName()] = t
}

func (panel *panel) GetCurrentTab() (t *tab) {
	index := panel.TabbedPanels.GetCurrentTab()
	elem := panel.tabs[index].(*tab)
	return elem
}

func (panel *panel) HasNonEmptyTab() bool {
	if len(panel.tabs) > 1 {
		return true
	}
	if _, ok := panel.tabs["No name"]; ok {
		return false
	}
	return true
}

func (panel *panel) HandleInput(tevent *tcell.EventKey) (retEvent *tcell.EventKey) {
	retEvent = panel.GetCurrentTab().HandleInput(tevent)
	if retEvent != tevent {
		return retEvent
	}
	return tevent
}
