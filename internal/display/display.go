package display

import (
	"strangelet/internal/event"

	"github.com/Kyrasuum/cview"
	"github.com/gdamore/tcell/v2"
)

var (
	displayFileBrowser = false
	CurDisplay         *Display
)

type Display struct {
	*cview.Flex

	cbar    CommandBar
	subFlex *cview.Flex

	fb Filebrowser

	rows   *cview.Flex
	cols   []*cview.Flex
	panels []Panel
}

func (dplay *Display) InitDisplay(app *cview.Application) {
	//enforce only one
	if CurDisplay != nil {
		return
	}
	//handle input
	app.SetInputCapture(dplay.HandleInput)

	//initialize overall flex for screen space (command bar, file browser, and panel display)
	dplay.Flex = cview.NewFlex()
	dplay.Flex.SetDirection(cview.FlexRow)
	dplay.subFlex = cview.NewFlex()
	dplay.fb.InitFilebrowser(dplay.subFlex)
	dplay.Flex.AddItem(dplay.subFlex, 0, 1, false)

	//initialize rows space
	dplay.rows = cview.NewFlex()
	dplay.rows.SetDirection(cview.FlexRow)

	//add initial row
	cols := dplay.AddPanelRow()
	//add dummy panels
	dplay.AddPanelToRow(cols, 0)
	dplay.AddPanelToRow(cols, 1)

	//add initial row
	cols = dplay.AddPanelRow()
	//add dummy panels
	dplay.AddPanelToRow(cols, 0)
	dplay.AddPanelToRow(cols, 1)

	//put it all together
	dplay.subFlex.AddItem(dplay.rows, 0, 1, false)
	dplay.cbar.InitCommandBar(dplay.Flex)

	//set root display area
	app.SetRoot(dplay, true)
	CurDisplay = dplay
}

func (dplay *Display) AddPanelRow() (cols *cview.Flex) {
	cols = cview.NewFlex()

	dplay.cols = append(dplay.cols, cols)
	dplay.rows.AddItem(cols, 0, 1, false)

	return cols
}

func (dplay *Display) AddPanelToRow(row *cview.Flex, panelIndex int) {
	var panel Panel
	panel.InitPanel(row, panelIndex)
	dplay.panels = append(dplay.panels, panel)
}

func (dplay *Display) HandleInput(tevent *tcell.EventKey) (retEvent *tcell.EventKey) {
	if tevent.Key() == tcell.KeyCtrlQ {
		event.Quit()
		return nil
	}
	if tevent.Key() == tcell.KeyCtrlC {
		return nil
	}
	if tevent.Key() == tcell.KeyCtrlW {
		return nil
	}
	retEvent = dplay.fb.HandleInput(tevent)
	if retEvent != tevent {
		return retEvent
	}
	for _, panel := range dplay.panels {
		retEvent = panel.HandleInput(tevent)
		if retEvent != tevent {
			return retEvent
		}
	}
	return tevent
}
