package event

import (
	"os"
	"os/signal"
	"syscall"

	"strangelet/internal/sync"

	"github.com/Kyrasuum/cview"
)

var (
	done    = make(chan struct{})
	sigterm chan os.Signal
	sighup  chan os.Signal
	app     *cview.Application
)

func InitEvents(ap *cview.Application) {
	sigterm = make(chan os.Signal, 1)
	sighup = make(chan os.Signal, 1)
	signal.Notify(sigterm, syscall.SIGTERM, syscall.SIGINT, syscall.SIGQUIT)
	signal.Notify(sighup, syscall.SIGHUP)

	app = ap

	go listenEvents()
}

func Quit() {
	close(done)
}

func listenEvents() {
	defer sync.Done()

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
	app.Stop()
	os.Exit(0)
}
