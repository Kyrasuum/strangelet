package display

import (
	"fmt"

	"github.com/Kyrasuum/cview"
)

var ()

type Panel struct {
	*cview.Flex
	subFlex *cview.Flex

	buffer Buffer
	gutter Gutter
	status StatusBar
}

func (panel *Panel) InitPanel(tabs *cview.TabbedPanels, index int) {
	panel.subFlex = cview.NewFlex()
	panel.gutter.InitGutter(panel.subFlex)
	panel.buffer.InitBuffer(panel.subFlex)

	panel.Flex = cview.NewFlex()
	panel.Flex.SetDirection(cview.FlexRow)
	panel.Flex.AddItem(panel.subFlex, 0, 1, false)
	panel.status.InitStatusBar(panel.Flex)

	tabs.AddTab(fmt.Sprintf("panel-%d", index), fmt.Sprintf("Panel #%d", index), panel)
}
