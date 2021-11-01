package app

import (
	"flag"
	"fmt"
	"os"

	"strangelet/internal/buffer"
	"strangelet/internal/config"
	"strangelet/internal/display"
	"strangelet/internal/event"
	"strangelet/internal/util"
	iapp "strangelet/pkg/app"

	"github.com/Kyrasuum/cview"

	_ "github.com/allan-simon/go-singleinstance"
	"github.com/go-errors/errors"
	lua "github.com/yuin/gopher-lua"
)

var (
	cviewApp *cview.Application

	focusStk *util.Stack
	focusMap map[cview.Primitive]struct{}
)

type application struct {
	iapp.App
}

func NewApp() (app application) {
	//handle calls post init
	if cviewApp != nil {
		return iapp.CurApp.(application)
	}

	//create singleton instance
	app = application{}
	cviewApp = cview.NewApplication()
	//run startup
	app.startApp()

	defer cviewApp.HandlePanic()

	//set cview flags
	cviewApp.EnableMouse(true)
	cviewApp.SetBeforeFocusFunc(app.focusHook)

	//initialize eventhandler
	event.InitEvents()

	//initialize focus handlers
	focusStk = &util.Stack{}
	focusMap = make(map[cview.Primitive]struct{})
	iapp.CurApp = app

	//start display
	frame := display.NewDisplay(cviewApp)

	//load buffers from flags
	args := flag.Args()
	b := app.LoadInput(args)

	if len(b) == 0 {
		// No buffers to open
		app.Stop()
	}
	//load tabs from buffers
	for _, bobj := range b {
		frame.AddTabToCurrentPanel(bobj)
	}
	//force a redraw
	app.Redraw(func() {})

	//postinit hook
	err := config.RunPluginFn("postinit")
	if err != nil {
		app.TermMessage(err)
	}
	//done
	return app
}

func (app application) startApp() {
	defer cviewApp.HandlePanic()

	//init log
	InitFlags()
	InitLog()

	//init config
	err := config.InitConfigDir(*flagConfigDir)
	if err != nil {
		app.TermMessage(err)
	}

	//init runtime
	config.InitRuntimeFiles()
	err = config.ReadSettings()
	if err != nil {
		app.TermMessage(err)
	}
	//init global
	err = config.InitGlobalSettings()
	if err != nil {
		app.TermMessage(err)
	}

	// flag options
	for k, v := range optionFlags {
		if *v != "" {
			nativeValue, err := config.GetNativeValue(k, config.DefaultAllSettings()[k], *v)
			if err != nil {
				app.TermMessage(err)
				continue
			}
			config.GlobalSettings[k] = nativeValue
		}
	}
	//process flags for plugins
	DoPluginFlags()

	//handler for errors
	defer func() {
		if err := recover(); err != nil {
			app.Stop()
			if e, ok := err.(*lua.ApiError); ok {
				fmt.Println("Lua API error:", e)
			} else {
				fmt.Println("Strangelet encountered an error:", errors.Wrap(err, 2).ErrorStack(), "\nIf you can reproduce this error, please report it")
			}
			// backup all open buffers
			for _, b := range buffer.OpenBuffers {
				b.Backup()
			}
			os.Exit(1)
		}
	}()

	//load all plugins
	err = config.LoadAllPlugins()
	if err != nil {
	}

	// action.InitBindings()
	// action.InitCommands()

	//load color scheme
	err = config.InitColorscheme()
	if err != nil {
		app.TermMessage(err)
	}

	//preinit hook
	err = config.RunPluginFn("preinit")
	if err != nil {
		app.TermMessage(err)
	}

	// action.InitGlobals()

	//init hook
	err = config.RunPluginFn("init")
	if err != nil {
		app.TermMessage(err)
	}

	// m := clipboard.SetMethod(config.GetGlobalOption("clipboard").(string))
	// go func() {
	// clipErr := clipboard.Initialize(m)
	//
	// if clipErr != nil {
	// log.Println(clipErr, " or change 'clipboard' option")
	// }
	// }()

	//setup autosave
	if a := config.GetGlobalOption("autosave").(float64); a > 0 {
		config.SetAutoTime(int(a))
		config.StartAutoSave()
	}

	go func() {
		for {
			if err := cviewApp.Run(); err != nil {
				panic(err)
			}
		}
	}()
}
