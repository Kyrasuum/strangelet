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
}

type Display interface{}
