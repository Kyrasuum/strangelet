package view

import (
	config "strangelet/internal/config"
	events "strangelet/internal/events"
	pub "strangelet/pkg/app"

	tea "github.com/charmbracelet/bubbletea"
	lipgloss "github.com/charmbracelet/lipgloss"
	statusbar "github.com/knipferrc/teacup/statusbar"
)

type view struct {
	active int

	width  int
	height int

	s  *split
	fb *filebrowser
	l  *logWindow
	c  *cmd
	sb *statusbar.Bubble
}

var (
	highlightColor = lipgloss.AdaptiveColor{Light: "#874BFD", Dark: "#7D56F4"}
	scopes         = map[int]string{
		splitView: "Split",
		filesView: "File Browser",
		logView:   "Log Window",
		cmdView:   "Command Bar",
	}
)

const (
	splitView int = iota
	filesView
	logView
	cmdView
)

func NewView(app pub.App) view {
	s := NewSplit(app)
	fb := NewFileBrowser(app)
	l := NewLog(app)
	c := NewCmd(app)
	sb := statusbar.New(
		statusbar.ColorConfig{
			Foreground: lipgloss.AdaptiveColor{Dark: "#ffffff", Light: "#ffffff"},
			Background: lipgloss.AdaptiveColor{Light: "#F25D94", Dark: "#F25D94"},
		},
		statusbar.ColorConfig{
			Foreground: lipgloss.AdaptiveColor{Light: "#ffffff", Dark: "#ffffff"},
			Background: lipgloss.AdaptiveColor{Light: "#3c3836", Dark: "#3c3836"},
		},
		statusbar.ColorConfig{
			Foreground: lipgloss.AdaptiveColor{Light: "#ffffff", Dark: "#ffffff"},
			Background: lipgloss.AdaptiveColor{Light: "#A550DF", Dark: "#A550DF"},
		},
		statusbar.ColorConfig{
			Foreground: lipgloss.AdaptiveColor{Light: "#ffffff", Dark: "#ffffff"},
			Background: lipgloss.AdaptiveColor{Light: "#6124DF", Dark: "#6124DF"},
		},
	)

	v := view{
		active: splitView,

		width:  0,
		height: 0,

		s:  &s,
		fb: &fb,
		l:  &l,
		c:  &c,
		sb: &sb,
	}

	return v
}

func (v view) Init() tea.Cmd {
	return tea.Batch(tea.EnterAltScreen)
}

func (v view) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var (
		cmd  tea.Cmd
		cmds []tea.Cmd
	)

	switch msg := msg.(type) {
	// These keys should exit the program.
	case events.FocusCommand:
		if v.active != cmdView {
			v.active = cmdView
		} else {
			v.active = splitView
		}
	case events.ToggleLogWindow:
		*v.l = v.l.ToggleVisible()
	case events.ToggleFileBrowser:
		*v.fb = v.fb.ToggleVisible()
		if v.fb.visible {
			v.active = filesView
		} else {
			v.active = splitView
		}
	case events.FocusFileBrowser:
		if v.active != filesView {
			v.active = filesView
		} else {
			v.active = splitView
		}
	case events.CloseApp:
		return v, tea.Quit
	case tea.WindowSizeMsg:
		v.sb.SetSize(msg.Width)
		v.sb.SetContent("test.txt", "~/.config/nvim", "1/23", "SB")
		v.width = msg.Width
		v.height = msg.Height

		c := v.c.ViewW(v.width)
		sb := v.sb.View()
		h := v.height - lipgloss.Height(c) - lipgloss.Height(sb)
		v.fb.SetHeight(h)
		v.l.SetHeight(h)
	case tea.KeyMsg:
		if action, ok := config.Bindings["Global"][msg.String()]; ok {
			if handler, ok := events.Actions[action]; ok {
				cmds = append(cmds, handler(msg))
			}
		}
	case tea.MouseMsg:
		if !config.GlobalSettings["mouse"].(bool) {
			break
		}
		switch msg.Type {
		case tea.MouseWheelUp:
		case tea.MouseWheelDown:
		}
	}

	switch v.active {
	case splitView:
		*v.s, cmd = v.s.UpdateTyped(msg)
		cmds = append(cmds, cmd)
	case filesView:
		*v.fb, cmd = v.fb.UpdateTyped(msg)
		cmds = append(cmds, cmd)
	case cmdView:
		*v.c, cmd = v.c.UpdateTyped(msg)
		cmds = append(cmds, cmd)
	}

	*v.l, cmd = v.l.UpdateTyped(msg)
	cmds = append(cmds, cmd)

	empty := true
	for _, cmd := range cmds {
		if cmd != nil {
			empty = false
		}
	}
	if empty {
		switch msg := msg.(type) {
		case tea.KeyMsg:
			cmds = append(cmds, func() tea.Msg {
				return events.LogMessage("Unknown Keybind for scope[Global, " + scopes[v.active] + "]: " + msg.String())
			})
		}
	}

	// Return the updated model to the Bubble Tea runtime for processing.
	// Note that we're not returning a command.
	return v, tea.Batch(cmds...)
}

func (v view) View() string {
	//initial render
	c := v.c.ViewW(v.width)
	sb := v.sb.View()

	fb := ""
	l := ""

	if v.fb.visible {
		fb = v.fb.View()
	}
	if v.l.visible {
		l = v.l.View()
	}

	s := v.s.ViewWH(
		v.width-lipgloss.Width(fb)-lipgloss.Width(l),
		v.height-lipgloss.Height(c)-lipgloss.Height(sb))

	display := lipgloss.JoinVertical(lipgloss.Left, lipgloss.JoinHorizontal(lipgloss.Top, fb, s, l), sb, c)

	return display
}
