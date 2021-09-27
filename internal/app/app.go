package app

import (
	"strangelet/internal/event"
	"strangelet/internal/sync"
	"strangelet/internal/util"

	"github.com/Kyrasuum/cview"
)

var (
	app      *cview.Application
	focusStk *util.Stack
	focusMap map[cview.Primitive]struct{}
)

func InitApp(ap *cview.Application) {
	app = ap
	defer app.HandlePanic()

	app.EnableMouse(true)
	app.SetBeforeFocusFunc(focusHook)

	event.InitEvents(app)

	focusStk = &util.Stack{}
	focusMap = make(map[cview.Primitive]struct{})
}

func StartApp() {
	if err := app.Run(); err != nil {
		panic(err)
	}
	sync.Wait()
}

func focusHook(prim cview.Primitive) bool {
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

func SetFocus(prim cview.Primitive) {
	app.SetFocus(prim)
}

func GetFocus() (prim cview.Primitive) {
	return app.GetFocus()
}
