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
	Redraw()
	Stop()
	Pause(f func())
}

type Display interface{}
