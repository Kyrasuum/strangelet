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
	dirty   bool
	frame   string

	height int

	log []string
}

var ()

const ()

func NewLog(app pub.App) *LogWindow {
	l := LogWindow{
		visible: false,
		dirty:   true,
		frame:   "",
	}
	return &l
}

func (l *LogWindow) Init() tea.Cmd {
	return nil
}

func (l *LogWindow) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmds []tea.Cmd

	switch msg := msg.(type) {
	case events.LogMessage:
		l.log = append(l.log, string(msg))
		log.Println(string(msg))
		cmds = append(cmds, events.Actions["NOOP"](""))
		l.Redraw()
	case events.ErrorMsg:
		l.log = append(l.log, string(msg.Error()))
		log.Println(string(msg.Error()))
		cmds = append(cmds, events.Actions["NOOP"](""))
		l.Redraw()
	case tea.MouseMsg:
		// tea.MouseEvent(msg)
	}

	// Return the updated model to the Bubble Tea runtime for processing.
	// Note that we're not returning a command.
	return l, tea.Batch(cmds...)
}

func (l *LogWindow) Redraw() {
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

	l.frame = config.ColorScheme["log"].
		Height(l.height).
		Width(width).
		Render(lipgloss.JoinVertical(lipgloss.Left, header, content))
}

func (l *LogWindow) View() string {
	if l.dirty {
		l.Redraw()
		l.dirty = false
	}

	return l.frame
}
func (l *LogWindow) ViewWH(w int, h int) string {
	l.SetHeight(h)
	return l.View()
}

func (l *LogWindow) SetHeight(h int) {
	l.height = h
	l.dirty = true
}

func (l *LogWindow) ToggleVisible() {
	l.visible = !l.visible
	l.dirty = true
}

func (l *LogWindow) Visible() bool {
	return l.visible
}

func (l *LogWindow) SetDirty() {
	l.dirty = true
}
