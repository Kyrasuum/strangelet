package display

import (
	"strangelet/pkg/app"

	"strangelet/internal/event"

	"github.com/Kyrasuum/cview"
	"github.com/gdamore/tcell/v2"
)

var (
	displayFileBrowser = false
	curDisplay         *display
)

type display struct {
	app.Display
	*cview.Flex

	cbar    *commandBar
	subFlex *cview.Flex

	fb  *filebrowser
	log *logWin

	rows   *cview.Flex
	cols   []*cview.Flex
	panels []*panel
}

func NewDisplay(app *cview.Application) (dplay *display) {
	//enforce only one
	if curDisplay != nil {
		return curDisplay
	}
	dplay = &display{}

	//handle input
	app.SetInputCapture(dplay.HandleInput)

	//initialize overall flex for screen space (command bar, file browser, and panel display)
	dplay.Flex = cview.NewFlex()
	dplay.Flex.SetDirection(cview.FlexRow)
	dplay.subFlex = cview.NewFlex()
	dplay.fb = NewFilebrowser(dplay.subFlex)
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
	dplay.cbar = NewCommandBar(dplay.Flex)
	dplay.log = NewLogWin(dplay.subFlex)

	//set root display area
	app.SetRoot(dplay, true)
	curDisplay = dplay

	return dplay
}

func (dplay *display) AddPanelRow() (cols *cview.Flex) {
	cols = cview.NewFlex()

	dplay.cols = append(dplay.cols, cols)
	dplay.rows.AddItem(cols, 0, 1, false)

	return cols
}

func (dplay *display) AddPanelToRow(row *cview.Flex, panelIndex int) {
	pan := NewPanel(row, panelIndex)
	dplay.panels = append(dplay.panels, pan)
}

func (dplay *display) HandleInput(tevent *tcell.EventKey) (retEvent *tcell.EventKey) {
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
	retEvent = dplay.log.HandleInput(tevent)
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
