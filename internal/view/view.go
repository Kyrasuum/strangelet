package view

import (
	"fmt"

	config "strangelet/internal/config"
	pub "strangelet/pkg/app"

	tea "github.com/charmbracelet/bubbletea"
	lipgloss "github.com/charmbracelet/lipgloss"
)

type view struct {
	active int

	width  int
	height int

	s  split
	fb filebrowser
	l  log
	c  cmd
}

var (
	highlightColor   = lipgloss.AdaptiveColor{Light: "#874BFD", Dark: "#7D56F4"}
	inactiveWinStyle = lipgloss.NewStyle().BorderStyle(lipgloss.HiddenBorder())
	activeWinStyle   = inactiveTabStyle.Copy().BorderStyle(lipgloss.NormalBorder()).BorderForeground(highlightColor).Padding(0, 1)
)

const (
	splitView int = iota
	filesView
	logView
	cmdView
)

func NewView(app pub.App) view {
	v := view{
		active: 0,

		width:  0,
		height: 0,

		s:  NewSplit(app),
		fb: NewFileBrowser(app),
		l:  NewLog(app),
		c:  NewCmd(app),
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
	case tea.WindowSizeMsg:
		v.width = msg.Width
		v.height = msg.Height
	case tea.KeyMsg:
		switch msg.String() {
		// These keys should exit the program.
		case "ctrl+c", "esc", "q":
			return v, tea.Quit
		case "tab":
			switch v.active {
			case splitView:
				v.active = filesView
			case filesView:
				v.active = logView
			case logView:
				v.active = cmdView
			case cmdView:
				v.active = splitView
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
		v.s, cmd = v.s.UpdateTyped(msg)
		cmds = append(cmds, cmd)
	case filesView:
		v.fb, cmd = v.fb.UpdateTyped(msg)
		cmds = append(cmds, cmd)
	case logView:
		v.l, cmd = v.l.UpdateTyped(msg)
		cmds = append(cmds, cmd)
	case cmdView:
		v.c, cmd = v.c.UpdateTyped(msg)
		cmds = append(cmds, cmd)
	}

	// Return the updated model to the Bubble Tea runtime for processing.
	// Note that we're not returning a command.
	return v, tea.Batch(cmds...)
}

func (v view) View() string {
	//initial render
	c := inactiveWinStyle.Render(v.c.ViewW(v.width - inactiveWinStyle.GetHorizontalFrameSize()))
	fb := inactiveWinStyle.Render(v.fb.ViewH(v.height - lipgloss.Height(c) - inactiveWinStyle.GetVerticalFrameSize()))
	l := inactiveWinStyle.Render(v.l.ViewH(v.height - lipgloss.Height(c) - inactiveWinStyle.GetVerticalFrameSize()))
	s := inactiveWinStyle.Render(v.s.ViewWH(
		v.width-lipgloss.Width(fb)-lipgloss.Width(l)-inactiveWinStyle.GetHorizontalFrameSize(),
		v.height-lipgloss.Height(c)-inactiveWinStyle.GetVerticalFrameSize()))

	//switch to highlight
	switch v.active {
	case splitView:
		s = activeWinStyle.Render(fmt.Sprintf("%4s", v.s.ViewWH(
			v.width-lipgloss.Width(fb)-lipgloss.Width(l)-activeWinStyle.GetHorizontalFrameSize(),
			v.height-lipgloss.Height(c)-inactiveWinStyle.GetVerticalFrameSize())))
	case filesView:
		fb = activeWinStyle.Render(fmt.Sprintf("%4s", v.fb.ViewH(
			v.height-lipgloss.Height(c)-inactiveWinStyle.GetVerticalFrameSize())))
	case logView:
		l = activeWinStyle.Render(fmt.Sprintf("%4s", v.l.ViewH(
			v.height-lipgloss.Height(c)-inactiveWinStyle.GetVerticalFrameSize())))
	case cmdView:
		c = activeWinStyle.Render(fmt.Sprintf("%4s", v.c.ViewW(
			v.width-inactiveWinStyle.GetHorizontalFrameSize())))
	}

	display := lipgloss.JoinVertical(lipgloss.Left, lipgloss.JoinHorizontal(lipgloss.Top, fb, s, l), c)

	return display
}
