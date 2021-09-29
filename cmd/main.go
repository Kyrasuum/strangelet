package main

import (
	"strangelet/internal/app"
	"strangelet/internal/sync"
)

var ()

func main() {
	app.NewApp()

	sync.Wait()
}
