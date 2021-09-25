package display

import (
	"github.com/Kyrasuum/cview"
	"github.com/gdamore/tcell/v2"
)

var (
	displayFileBrowser = false
)

type Display struct {
	*cview.Flex

	subFlex *cview.Flex
	fb      Filebrowser
	cbar    CommandBar
	pane    Pane
}

func (dplay *Display) InitDisplay(app *cview.Application) {
	dplay.subFlex = cview.NewFlex()
	dplay.fb.InitFilebrowser(dplay.subFlex)
	dplay.pane.InitPane(dplay.subFlex)

	dplay.Flex = cview.NewFlex()
	dplay.Flex.SetDirection(cview.FlexRow)
	dplay.Flex.AddItem(dplay.subFlex, 0, 1, false)
	dplay.cbar.InitCommandBar(dplay.Flex)

	app.SetRoot(dplay, true)
}

func (dplay *Display) HandleInput(event *tcell.EventKey) (retEvent *tcell.EventKey) {
	if event.Key() == tcell.KeyCtrlD {
		dplay.fb.ToggleDisplay(dplay.subFlex)
	}
	if dplay.fb.IsVisible() {
		retEvent = dplay.fb.HandleInput(event)
		if retEvent != event {
			return retEvent
		}
	}
	return retEvent
}
