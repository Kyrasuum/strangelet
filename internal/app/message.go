package app

import (
	"fmt"
	"runtime"
)

var ()

func (app application) Stop() {
	cviewApp.Stop()
	runtime.Goexit()
}

func (app application) Pause(f func()) {
	cviewApp.Suspend(f)
}

func (app application) TermMessage(msg ...interface{}) {
	app.Pause(func() {
		fmt.Println(msg...)
		fmt.Println("\nPress enter to continue")

		fmt.Scanln()
	})
}
