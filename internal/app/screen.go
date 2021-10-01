package app

import ()

var ()

func (app application) Redraw() {
	cviewApp.QueueUpdateDraw(func() {})
}
