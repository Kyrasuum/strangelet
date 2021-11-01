package display

import (
	"strangelet/pkg/app"

	"github.com/Kyrasuum/cview"
	"github.com/gdamore/tcell/v2"
)

var (
	logbgColor       tcell.Color
	logheaderbgColor tcell.Color

	logWinW int = 30

	CurLogWin *logWin
)

type logWin struct {
	*cview.Flex

	header *cview.Table
	log    *cview.TextView

	parentFlex *cview.Flex
}

func NewLogWin(subFlex *cview.Flex) (log *logWin) {
	log = &logWin{}
	//enforce only one
	if CurLogWin != nil {
		return
	}

	//init colors
	if logbgColor == 0 {
		logbgColor = tcell.NewRGBColor(30, 30, 30)
	}
	if logheaderbgColor == 0 {
		logheaderbgColor = tcell.NewRGBColor(150, 0, 0)
	}

	//setup log
	log.log = cview.NewTextView()
	log.log.Box.SetBackgroundColor(logbgColor)

	//setup header
	log.header = cview.NewTable()
	cell := cview.NewTableCell("")
	cell.SetText("Log")
	cell.SetAttributes(tcell.AttrUnderline)
	log.header.SetCell(0, 0, cell)
	log.header.Box.SetBackgroundColor(logheaderbgColor)

	//setup flex
	log.Flex = cview.NewFlex()
	log.Flex.SetDirection(cview.FlexRow)
	log.Flex.AddItem(log.header, 1, 1, false)
	log.Flex.AddItem(log.log, 0, 1, false)

	subFlex.AddItem(log, logWinW, 1, false)
	log.parentFlex = subFlex
	CurLogWin = log
	// Default to closed
	log.ToggleDisplay()

	return log
}

func (log *logWin) Render(scr tcell.Screen) bool {
	return false
}

func (log *logWin) IsVisible() bool {
	return log.Box.GetVisible()
}

func (log *logWin) HandleInput(tevent *tcell.EventKey) *tcell.EventKey {
	if tevent.Key() == tcell.KeyCtrlL {
		log.ToggleDisplay()
		return nil
	}
	return tevent
}

func (log *logWin) HandleMouse(event *tcell.EventMouse, action cview.MouseAction) (*tcell.EventMouse, cview.MouseAction) {
	return event, action
}

func (log *logWin) ToggleDisplay() {
	if log.Box.GetVisible() {
		log.parentFlex.ResizeItem(log, -1, 0)
		log.Box.SetVisible(false)
		app.CurApp.SetFocus(log)
	} else {
		log.parentFlex.ResizeItem(log, logWinW, 1)
		log.Box.SetVisible(true)
		app.CurApp.SetFocus(nil)
	}
	app.CurApp.Redraw(func() {})
}
