package logWindow

import (
	"log"
	"strings"

	config "strangelet/internal/config"
	events "strangelet/internal/events"
	pub "strangelet/pkg/app"

	tea "github.com/charmbracelet/bubbletea"
	lipgloss "github.com/charmbracelet/lipgloss"
	wordwrap "github.com/muesli/reflow/wordwrap"
)

type LogWindow struct {
	visible bool

	height int

	log []string
}

var ()

const ()

func NewLog(app pub.App) LogWindow {
	return LogWindow{
		visible: false,
	}
}

func (l LogWindow) Init() tea.Cmd {
	return nil
}

func (l LogWindow) Update(msg tea.Msg) (tea.Model, tea.Cmd) { return l.UpdateTyped(msg) }
func (l LogWindow) UpdateTyped(msg tea.Msg) (LogWindow, tea.Cmd) {
	switch msg := msg.(type) {
	case events.LogMessage:
		l.log = append(l.log, string(msg))
		log.Println(string(msg))
	case events.ErrorMsg:
		l.log = append(l.log, string(msg.Error()))
		log.Println(string(msg.Error()))
	case tea.MouseMsg:
		// tea.MouseEvent(msg)
	}

	// Return the updated model to the Bubble Tea runtime for processing.
	// Note that we're not returning a command.
	return l, nil
}

func (l LogWindow) View() string {
	width := int(config.GlobalSettings["logwidth"].(float64))
	header := "Log:"

	ch := 1
	lines := []string{}
	for i := len(l.log) - 1; i >= 0; i-- {
		msg := wordwrap.String(l.log[i], width)
		h := lipgloss.Height(msg)
		ch += h
		if ch > l.height {
			break
		}
		lines = append([]string{msg}, lines...)
	}
	content := strings.Join(lines, "\n")

	res := config.ColorScheme["log"].
		Height(l.height).
		Width(width).
		Render(lipgloss.JoinVertical(lipgloss.Left, header, content))

	return res
}
func (l LogWindow) ViewWH(w int, h int) string {
	l.SetHeight(h)
	return l.View()
}

func (l LogWindow) SetHeight(h int) LogWindow {
	l.height = h
	return l
}

func (l LogWindow) ToggleVisible() LogWindow {
	l.visible = !l.visible
	return l
}

func (l LogWindow) Visible() bool {
	return l.visible
}
