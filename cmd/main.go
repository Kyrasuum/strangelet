package main

import (
	"os"
	"os/signal"
	"sync"
	"syscall"

	"strangelet/internal/display"

	"github.com/Kyrasuum/cview"
	"github.com/gdamore/tcell/v2"
)

var (
	app = cview.NewApplication()

	frame display.Display

	done = make(chan struct{})

	waitgrp sync.WaitGroup
	sigterm chan os.Signal
	sighup  chan os.Signal
)

func main() {
	defer app.HandlePanic()

	sigterm = make(chan os.Signal, 1)
	sighup = make(chan os.Signal, 1)
	signal.Notify(sigterm, syscall.SIGTERM, syscall.SIGINT, syscall.SIGQUIT)
	signal.Notify(sighup, syscall.SIGHUP)

	app.EnableMouse(true)

	//handle input
	app.SetInputCapture(HandleInput)

	frame.InitDisplay(app)

	if err := app.Run(); err != nil {
		panic(err)
	}
	waitgrp.Wait()
}

func HandleInput(event *tcell.EventKey) *tcell.EventKey {
	if event.Key() == tcell.KeyCtrlC {
		Quit()
	}
	return frame.HandleInput(event)
}

func Quit() {
	close(done)
	app.Stop()
}

func DoEvent() {
	defer waitgrp.Done()

	select {
	case <-done:
		return
	case <-sighup:
		Quit()
	case <-sigterm:
		Quit()
	}
}
