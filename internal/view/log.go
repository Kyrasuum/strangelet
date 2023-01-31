package view

import (
	"log"

	config "strangelet/internal/config"
	events "strangelet/internal/events"
	pub "strangelet/pkg/app"

	tea "github.com/charmbracelet/bubbletea"
	lipgloss "github.com/charmbracelet/lipgloss"
)

type errorMsg error

type logWindow struct {
	visible bool

	height int

	log string
}

var (
	logstyle = lipgloss.NewStyle()
)

const ()

func NewLog(app pub.App) logWindow {
	return logWindow{
		visible: false,
	}
}

func (l logWindow) Init() tea.Cmd {
	return nil
}

func (l logWindow) Update(msg tea.Msg) (tea.Model, tea.Cmd) { return l.UpdateTyped(msg) }
func (l logWindow) UpdateTyped(msg tea.Msg) (logWindow, tea.Cmd) {
	switch msg := msg.(type) {
	case events.LogMessage:
		l.log = string(msg)
	case errorMsg:
		log.Println(string(msg.Error()))
	case tea.MouseMsg:
		// tea.MouseEvent(msg)
	}

	// Return the updated model to the Bubble Tea runtime for processing.
	// Note that we're not returning a command.
	return l, nil
}

func (l logWindow) View() string {
	return logstyle.Height(l.height).Width(int(config.GlobalSettings["logwidth"].(float64))).Render("Log:\n" + l.log)
}

func (l logWindow) SetHeight(h int) {
	l.height = h
}

func (l logWindow) ToggleVisible() logWindow {
	l.visible = !l.visible
	return l
}
