package display

import (
	"fmt"

	"strangelet/internal/app"

	"github.com/Kyrasuum/cview"
	"github.com/gdamore/tcell/v2"
)

var ()

type Tab struct {
	*cview.Flex
	row *cview.Flex

	name  string
	label string

	buffer Buffer
	gutter Gutter
	status StatusBar
}

func (tab *Tab) InitTab(tabs *cview.TabbedPanels, index int) {
	tab.Flex = cview.NewFlex()
	tab.Flex.SetDirection(cview.FlexRow)
	tab.row = cview.NewFlex()

	tab.gutter.InitGutter(tab.row)
	tab.buffer.InitBuffer(tab.row)
	tab.Flex.AddItem(tab.row, 0, 1, false)
	tab.status.InitStatusBar(tab.Flex)

	tab.name = fmt.Sprintf("tab-%d", index)
	tab.label = fmt.Sprintf("Tab #%d", index)

	tabs.AddTab(tab.name, tab.label, tab)
	tabs.SetCurrentTab(tab.name)
	app.SetFocus(tab)
}

func (tab *Tab) HandleInput(tevent *tcell.EventKey) (retEvent *tcell.EventKey) {
	return tevent
}

func (tab *Tab) GetName() string {
	return tab.name
}
