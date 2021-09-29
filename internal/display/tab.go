package display

import (
	"fmt"

	"strangelet/pkg/app"

	"github.com/Kyrasuum/cview"
	"github.com/gdamore/tcell/v2"
)

var ()

type tab struct {
	*cview.Flex
	row *cview.Flex

	name  string
	label string

	buffer *buffer
	gutter *gutter
	status *statusBar
}

func NewTab(tabs *cview.TabbedPanels, index int) (t *tab) {
	t = &tab{}

	t.Flex = cview.NewFlex()
	t.Flex.SetDirection(cview.FlexRow)
	t.row = cview.NewFlex()

	t.gutter = NewGutter(t.row)
	t.buffer = NewBuffer(t.row)
	t.Flex.AddItem(t.row, 0, 1, false)
	t.status = NewStatusBar(t.Flex)

	t.name = fmt.Sprintf("tab-%d", index)
	t.label = fmt.Sprintf("Tab #%d", index)

	tabs.AddTab(t.name, t.label, t)
	tabs.SetCurrentTab(t.name)
	app.CurApp.SetFocus(t)

	return t
}

func (tab *tab) HandleInput(tevent *tcell.EventKey) (retEvent *tcell.EventKey) {
	return tevent
}

func (tab *tab) GetName() string {
	return tab.name
}
