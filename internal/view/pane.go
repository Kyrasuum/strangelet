package view

import (
	pub "strangelet/pkg/app"

	tea "github.com/charmbracelet/bubbletea"
	lipgloss "github.com/charmbracelet/lipgloss"
)

type pane struct {
	active     int
	tabs       []string
	tabContent []elem
}

var (
	tabBorder = lipgloss.Border{
		Top:         "",
		Bottom:      "",
		Left:        "[",
		Right:       "]",
		TopLeft:     "",
		TopRight:    "",
		BottomLeft:  "",
		BottomRight: "",
	}
	activeTabBackground = lipgloss.Color("#B00000")
	tabBackground       = lipgloss.Color("#303030")
	tabForeground       = lipgloss.Color("#F0F0F0")
	tabsStyle           = lipgloss.NewStyle().Background(tabBackground).Foreground(tabForeground).BorderBackground(tabBackground).BorderForeground(tabForeground)
	tabStyle            = tabsStyle.Copy()
	inactiveTabStyle    = tabStyle.Copy().Border(lipgloss.HiddenBorder()).BorderTop(false).BorderBottom(false).BorderLeft(true).BorderRight(true)
	activeTabStyle      = inactiveTabStyle.Copy().Border(tabBorder).BorderTop(false).BorderBottom(false).BorderLeft(true).BorderRight(true).Background(activeTabBackground).BorderBackground(activeTabBackground)
	paneStyle           = lipgloss.NewStyle().Background(lipgloss.Color("#282828"))
	inactivePaneStyle   = lipgloss.NewStyle()
	activePaneStyle     = inactivePaneStyle.Copy()
)

const ()

func NewPane(app pub.App) pane {
	p := pane{
		active:     0,
		tabs:       []string{},
		tabContent: []elem{},
	}
	return p
}

func (p pane) Init() tea.Cmd {
	return nil
}

func (p pane) Update(msg tea.Msg) (tea.Model, tea.Cmd)    { return p.UpdateTyped(msg) }
func (p pane) UpdateI(msg tea.Msg) (interface{}, tea.Cmd) { return p.UpdateTyped(msg) }
func (p pane) UpdateTyped(msg tea.Msg) (pane, tea.Cmd) {
	var (
		cmds []tea.Cmd
	)
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+t":
			c := NewCode()
			cmd := c.OpenFile("internal/view/view.go")
			p.tabContent = append(p.tabContent, c)
			p.tabs = append(p.tabs, "internal/view/view.go")
			cmds = append(cmds, cmd)
		}
	case tea.MouseMsg:
		// tea.MouseEvent(msg)
	}
	if len(p.tabs) > 0 {
		e, cmd := p.tabContent[p.active].Update(msg)
		p.tabContent[p.active] = e.(elem)
		cmds = append(cmds, cmd)
	}

	// Return the updated model to the Bubble Tea runtime for processing.
	// Note that we're not returning a command.
	return p, tea.Batch(cmds...)
}

func (p pane) View() string { return p.ViewWH(0, 0) }
func (p pane) ViewWH(w, h int) string {
	ts := []string{}
	for i := 0; i < len(p.tabs); i++ {
		if i == p.active {
			ts = append(ts, activeTabStyle.Render(p.tabs[i]))
		} else {
			ts = append(ts, inactiveTabStyle.Render(p.tabs[i]))
		}
	}
	tabs := tabsStyle.Width(w).Height(1).Render(lipgloss.JoinHorizontal(lipgloss.Top, ts...))

	if len(p.tabs) > 0 {
		content := paneStyle.Width(w).Height(h - 1).Render(p.tabContent[p.active].ViewWH(w, h-1))
		return lipgloss.JoinVertical(lipgloss.Left, tabs, content)
	} else {
		return tabs
	}
}
