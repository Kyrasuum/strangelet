package app

import (
	"strangelet/internal/util"

	"github.com/Kyrasuum/cview"
)

var ()

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
