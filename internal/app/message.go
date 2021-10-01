package app

import (
	"bufio"
	"fmt"
	"os"
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

		reader := bufio.NewReader(os.Stdin)
		reader.ReadString('\n')
	})
}
