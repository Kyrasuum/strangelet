package display

import (
	"fmt"

	buff "strangelet/internal/buffer"

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

func NewTab(tabs *cview.TabbedPanels, b *buff.Buffer) (t *tab) {
	t = &tab{}

	t.Flex = cview.NewFlex()
	t.Flex.SetDirection(cview.FlexRow)
	t.row = cview.NewFlex()

	t.gutter = NewGutter(t.row)
	t.buffer = NewBuffer(t.row, b)
	t.Flex.AddItem(t.row, 0, 1, false)
	t.status = NewStatusBar(t.Flex)

	t.name = fmt.Sprintf(b.GetName())
	t.label = fmt.Sprintf(b.GetName())

	tabs.AddTab(t.name, t.label, t)
	tabs.SetCurrentTab(t.name)

	return t
}

func (tab *tab) HandleInput(tevent *tcell.EventKey) (retEvent *tcell.EventKey) {
	return tevent
}

func (tab *tab) GetName() string {
	return tab.name
}
