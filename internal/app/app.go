package app

import (
	iapp "strangelet/pkg/app"

	"strangelet/internal/display"
	"strangelet/internal/event"
	"strangelet/internal/util"

	"github.com/Kyrasuum/cview"
)

var (
	cviewApp *cview.Application
	frame    iapp.Display

	focusStk *util.Stack
	focusMap map[cview.Primitive]struct{}
)

type application struct {
	iapp.App
}

func NewApp() (app application) {
	if cviewApp != nil {
		return iapp.CurApp.(application)
	}

	app = application{}
	cviewApp = cview.NewApplication()
	go app.startApp()

	defer cviewApp.HandlePanic()

	cviewApp.EnableMouse(true)
	cviewApp.SetBeforeFocusFunc(app.focusHook)

	event.InitEvents(cviewApp)

	focusStk = &util.Stack{}
	focusMap = make(map[cview.Primitive]struct{})
	iapp.CurApp = app

	frame = display.NewDisplay(cviewApp)

	return app
}

func (app application) startApp() {
	defer cviewApp.HandlePanic()

	if err := cviewApp.Run(); err != nil {
		panic(err)
	}
}

func (app application) focusHook(prim cview.Primitive) bool {
	defer cviewApp.HandlePanic()
	// if prim is nil then we are removing focus from current object
	if prim == nil {
		// attempt to set focus to next in stack
		prim = focusStk.Pop().(cview.Primitive)
		delete(focusMap, prim)
		if prim != nil {
			app.SetFocus(prim)
			return false
		}
		// nothing left in stack
		return true
	}
	// check if we are setting focus to something already in stack
	if _, ok := focusMap[prim]; ok {
		//remove from focusStk so that it can be added at top
		retrnStk := &util.Stack{}
		for i := 0; i < focusStk.Len(); i++ { //have to use len of stack so that we can pop it
			elem := focusStk.Pop()
			if prim == elem {
				break
			}
			retrnStk.Push(elem)
		}
		//can loop over elements of stack because we dont care about retrnStk after this
		for _, elem := range *retrnStk {
			focusStk.Push(elem)
		}
		//lastly remove from map
		delete(focusMap, prim)
	}
	focusStk.Push(prim)
	focusMap[prim] = struct{}{}
	return true
}

func (app application) SetFocus(prim cview.Primitive) {
	defer cviewApp.HandlePanic()
	cviewApp.SetFocus(prim)
}

func (app application) GetFocus() (prim cview.Primitive) {
	defer cviewApp.HandlePanic()
	return cviewApp.GetFocus()
}
