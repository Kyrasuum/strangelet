package app

import (
	"github.com/Kyrasuum/cview"
)

var (
	CurApp App
)

type App interface {
	SetFocus(prim cview.Primitive)
	GetFocus() (prim cview.Primitive)
	Redraw(f func())
	Stop()
	Pause(f func())
	TermMessage(msg ...interface{})
}

type Display interface{}
