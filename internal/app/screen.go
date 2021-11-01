package app

import ()

var ()

func (app application) Redraw(f func()) {
	cviewApp.QueueUpdateDraw(f)
}
