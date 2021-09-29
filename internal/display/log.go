package display

import (
	"strangelet/pkg/app"

	"github.com/Kyrasuum/cview"
	"github.com/gdamore/tcell/v2"
)

var (
	logbgColor tcell.Color

	logWinW int = 30

	CurLogWin *logWin
)

type logWin struct {
	*cview.TextView

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

	//setup log
	log.TextView = cview.NewTextView()
	log.TextView.Box.SetBackgroundColor(logbgColor)

	subFlex.AddItem(log, logWinW, 1, false)
	log.parentFlex = subFlex
	CurLogWin = log
	// Default to closed
	log.ToggleDisplay()

	return log
}

func (log *logWin) IsVisible() bool {
	return log.TextView.Box.GetVisible()
}

func (log *logWin) HandleInput(tevent *tcell.EventKey) *tcell.EventKey {
	if tevent.Key() == tcell.KeyCtrlL {
		log.ToggleDisplay()
		return nil
	}
	return tevent
}

func (log *logWin) ToggleDisplay() {
	if log.TextView.Box.GetVisible() {
		log.parentFlex.ResizeItem(log, -1, 0)
		log.TextView.Box.SetVisible(false)
		app.CurApp.SetFocus(log)
	} else {
		log.parentFlex.ResizeItem(log, logWinW, 1)
		log.TextView.Box.SetVisible(true)
		app.CurApp.SetFocus(nil)
	}
}
