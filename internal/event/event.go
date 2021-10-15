package event

import (
	"os"
	"os/signal"
	"syscall"

	"strangelet/internal/sync"
	"strangelet/pkg/app"

	"github.com/rjeczalik/notify"
)

var (
	done    = make(chan struct{})
	sigterm chan os.Signal
	sighup  chan os.Signal

	run = true
)

func InitEvents() {
	sigterm = make(chan os.Signal, 1)
	sighup = make(chan os.Signal, 1)
	signal.Notify(sigterm, syscall.SIGTERM, syscall.SIGINT, syscall.SIGQUIT)
	signal.Notify(sighup, syscall.SIGHUP)

	go listenEvents()
	go netEvents()
}

func Quit() {
	close(done)
}

func netEvents() {
	// Make the channel buffered to ensure no event is dropped. Notify will drop
	// an event if the receiver is not able to keep up the sending pace.
	c := make(chan notify.EventInfo, 1)

	for run {
		// Set up a watchpoint listening on events within current working directory.
		// Dispatch each create and remove events separately to c.
		if err := notify.Watch(".", c, notify.Create, notify.Remove); err != nil {
			app.CurApp.TermMessage(err.Error())
		}
		defer notify.Stop(c)

		// Block until an event is received.
		ei := <-c
		app.CurApp.TermMessage(ei)
	}
}

func listenEvents() {
	defer func() {
		sync.Done()
		run = false
	}()
	sync.Add(1)
	for {
		select {
		case <-done:
			quit()
			return
		case <-sighup:
			quit()
			return
		case <-sigterm:
			quit()
			return
		}
	}
}

func quit() {
	app.CurApp.Pause(func() { app.CurApp.Stop() })
}
