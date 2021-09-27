package main

import (
	"strangelet/internal/app"
	"strangelet/internal/display"

	"github.com/Kyrasuum/cview"
)

var (
	frame display.Display
)

func main() {
	ap := cview.NewApplication()
	app.InitApp(ap)
	frame.InitDisplay(ap)

	app.StartApp()
}
