package app

import (
	"bufio"
	"flag"
	"fmt"
	"log"
	"os"
	"os/signal"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"
	"syscall"

	clipboard "strangelet/internal/clipboard"
	config "strangelet/internal/config"
	util "strangelet/internal/util"
	view "strangelet/internal/view"
	pub "strangelet/pkg/app"

	singleinstance "github.com/allan-simon/go-singleinstance"
	tea "github.com/charmbracelet/bubbletea"
)

var (
	// Event channel
	autosave chan bool

	// Command line flags
	flagVersion   = flag.Bool("version", false, "Show the version number and information")
	flagConfigDir = flag.String("config-dir", "", "Specify a custom location for the configuration directory")
	flagClean     = flag.Bool("clean", false, "Clean configuration directory")

	sigterm chan os.Signal
	sighup  chan os.Signal
)

type subapp struct {
	Shutdown chan int

	Log  *os.File
	Lock *os.File
	Pipe *os.File
}

func NewApp() (app pub.App) {
	//create private app space
	priv := &subapp{Shutdown: make(chan int), Lock: nil}
	app.Priv = priv

	//grab flags
	InitFlags()
	args := flag.Args()

	//setup home
	err := config.InitConfigDir("")
	if err != nil {
		fmt.Println(err)
		os.Exit(0)
	}

	//setup settings
	config.InitRuntimeFiles()
	err = config.ReadSettings()
	if err != nil {
		fmt.Println(err)
		os.Exit(0)
	}
	err = config.InitGlobalSettings()
	if err != nil {
		fmt.Println(err)
		os.Exit(0)
	}
	err = config.InitColorscheme()
	if err != nil {
		fmt.Println(err)
		os.Exit(0)
	}
	err = config.InitBindings()
	if err != nil {
		fmt.Println(err)
		os.Exit(0)
	}

	//get single instance lock
	lockFile, err := singleinstance.CreateLockFile(filepath.Join(config.ConfigDir, config.GlobalSettings["lockname"].(string)))
	if err != nil {
		//another instance is already running attempt to pass to other instance
		_, err := singleinstance.GetLockFilePid(filepath.Join(config.ConfigDir, config.GlobalSettings["lockname"].(string)))
		if err != nil {
			//error occured, unrecoverable
			fmt.Println("Cannot get PID:", err)
			app.StartApp = func() { CloseApp(app, 1) }
			return app
		}

		//pass arguements
		app.StartApp = func() { PassArgs(app, args) }
		return app
	}
	priv.Lock = lockFile

	//load arguments
	LoadArgs(args)

	//setup logging
	priv.Log, err = os.OpenFile(filepath.Join(config.ConfigDir, config.GlobalSettings["logname"].(string)), os.O_RDWR|os.O_CREATE, 0777)
	if err == nil {
		log.SetOutput(priv.Log)
	} else {
		fmt.Println(err)
		app.StartApp = func() { CloseApp(app, 1) }
		return app
	}

	//create named pipe for IPC
	os.Remove(filepath.Join(config.ConfigDir, config.GlobalSettings["pipename"].(string)))
	err = syscall.Mkfifo(filepath.Join(config.ConfigDir, config.GlobalSettings["pipename"].(string)), 0666)
	if err != nil {
		log.Println("Error making named pipe")
		log.Println(err)
		app.StartApp = func() { CloseApp(app, 1) }
		return app
	}
	//open named pipe for IPC
	priv.Pipe, err = os.OpenFile(filepath.Join(config.ConfigDir, config.GlobalSettings["pipename"].(string)), os.O_RDWR, os.ModeNamedPipe)

	//setup clipboard
	method := clipboard.SetMethod(config.GetGlobalOption("clipboard").(string))
	err = clipboard.Initialize(method)
	if err != nil {
		log.Println(err, " or change 'clipboard' option")
	}

	//setup signal handlers
	sigterm = make(chan os.Signal, 1)
	sighup = make(chan os.Signal, 1)
	signal.Notify(sigterm, syscall.SIGTERM, syscall.SIGINT, syscall.SIGQUIT)
	signal.Notify(sighup, syscall.SIGHUP)

	//continue startup
	app.StartApp = func() { StartApp(app) }

	return app
}

func InitFlags() {
	flag.Usage = func() {
		fmt.Println("Usage: micro [OPTIONS] [FILE]...")
		fmt.Println("-clean")
		fmt.Println("    \tCleans the configuration directory")
		fmt.Println("-config-dir dir")
		fmt.Println("    \tSpecify a custom location for the configuration directory")
		fmt.Println("-version")
		fmt.Println("    \tShow the version number and information")
		fmt.Println("[FILE]:LINE:COL")
	}

	flag.Parse()

	if *flagVersion {
		// If -version was passed
		fmt.Println("Version:", util.Version)
		fmt.Println("Commit hash:", util.CommitHash)
		fmt.Println("Compiled on", util.CompileDate)
		os.Exit(0)
	}
}

func LoadArgs(args []string) []struct {
	name string
	line int
	col  int
} {
	var err error
	files := make([]struct {
		name string
		line int
		col  int
	}, 0, len(args))
	flagr := regexp.MustCompile(`^\+(\d+)(?::(\d+))?$`)

	for _, a := range args {
		match := flagr.FindStringSubmatch(a)
		line := 0
		col := 0

		if len(match) == 3 && match[2] != "" {
			line, err = strconv.Atoi(match[1])
			if err != nil {
				log.Println(err)
				continue
			}
			col, err = strconv.Atoi(match[2])
			if err != nil {
				log.Println(err)
				continue
			}
		}

		if len(match) == 3 && match[2] == "" {
			line, err = strconv.Atoi(match[1])
			if err != nil {
				log.Println(err)
				continue
			}
		}

		files = append(files, struct {
			name string
			line int
			col  int
		}{name: a, line: line, col: col})
	}

	return files
}

func PassArgs(app pub.App, args []string) {
	//get connection to locking instance
	pipe, err := os.OpenFile(filepath.Join(config.ConfigDir, config.GlobalSettings["pipename"].(string)), os.O_WRONLY, os.ModeNamedPipe)
	if err != nil {
		fmt.Println("Failed to open named pipe")
		fmt.Println(err)
		CloseApp(app, 1)
	}

	pipe.WriteString(fmt.Sprintf("%s\n", strings.Join(args, " ")))
	pipe.Close()

	go func() {
		<-app.Priv.(*subapp).Shutdown
	}()

	CloseApp(app, 0)
}

func ReadArgs(app pub.App) {
	ch := make(chan string)
	go func() {
		for {
			select {
			case flag := <-app.Priv.(*subapp).Shutdown:
				app.Priv.(*subapp).Shutdown <- flag
				break
			case msg := <-ch:
				fmt.Print(msg)
			}
		}
	}()

	//read input from secondary instance
	reader := bufio.NewReader(app.Priv.(*subapp).Pipe)
	for {
		//check for message
		line, err := reader.ReadString('\n')
		if err != nil {
			log.Println("Error reading from named pipe")
			log.Println(err)
			CloseApp(app, 1)
			return
		}

		//check for shutdown
		select {
		case flag := <-app.Priv.(*subapp).Shutdown:
			app.Priv.(*subapp).Shutdown <- flag
			break
		case ch <- line:
			//handle message
		}
	}
}

func HandleEvents(app pub.App) {
	for {
		select {
		case <-sighup:
			CloseApp(app, 0)
			break
		case <-sigterm:
			CloseApp(app, 0)
			break
		case flag := <-app.Priv.(*subapp).Shutdown:
			app.Priv.(*subapp).Shutdown <- flag
			break
		default:
		}
	}
}

func CloseApp(app pub.App, flag int) {
	//send shutdown flag
	app.Priv.(*subapp).Shutdown <- flag

	//release single instance lock
	if app.Priv.(*subapp).Lock != nil {
		app.Priv.(*subapp).Lock.Close()
	}
	//close named pipe
	if app.Priv.(*subapp).Pipe != nil {
		app.Priv.(*subapp).Pipe.Close()
	}
	//close log
	if app.Priv.(*subapp).Log != nil {
		app.Priv.(*subapp).Log.Close()
	}

	os.Exit(flag)
}

func StartApp(app pub.App) {
	//create view
	v := view.NewView(app)
	app.View = v

	//setup server for secondary process creations
	go ReadArgs(app)

	//setup event handling routine
	go HandleEvents(app)

	//start UI
	p := tea.NewProgram(v, tea.WithAltScreen())
	if _, err := p.Run(); err != nil {
		log.Printf("Error starting UI: %v", err)
		CloseApp(app, 1)
	}
}
