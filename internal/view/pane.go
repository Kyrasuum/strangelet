package view

import (
	config "strangelet/internal/config"
	events "strangelet/internal/events"
	util "strangelet/internal/util"

	tea "github.com/charmbracelet/bubbletea"
	lipgloss "github.com/charmbracelet/lipgloss"
)

type pane struct {
	active     int
	tabs       []string
	tabContent []elem
	tabx       int
}

var (
	activetabBorder   lipgloss.Border
	inactivetabBorder lipgloss.Border

	paneBackground        lipgloss.TerminalColor
	tabForeground         lipgloss.TerminalColor
	tabBackground         lipgloss.TerminalColor
	activetabForeground   lipgloss.TerminalColor
	activetabBackground   lipgloss.TerminalColor
	inactivetabForeground lipgloss.TerminalColor
	inactivetabBackground lipgloss.TerminalColor

	tabsStyle         lipgloss.Style
	tabStyle          lipgloss.Style
	inactiveTabStyle  lipgloss.Style
	activeTabStyle    lipgloss.Style
	paneStyle         lipgloss.Style
	inactivePaneStyle lipgloss.Style
	activePaneStyle   lipgloss.Style
)

const ()

func UpdateStyling() {
	activetabBorder = lipgloss.Border{
		Top:         "",
		Bottom:      "",
		Left:        "┃",
		Right:       "┃",
		TopLeft:     "",
		TopRight:    "",
		BottomLeft:  "",
		BottomRight: "",
	}
	inactivetabBorder = lipgloss.Border{
		Top:         "",
		Bottom:      "",
		Left:        "│",
		Right:       "│",
		TopLeft:     "",
		TopRight:    "",
		BottomLeft:  "",
		BottomRight: "",
	}

	paneBackground = config.ColorScheme["background"].GetBackground()

	tabForeground = config.ColorScheme["tabbar"].GetForeground()
	tabBackground = config.ColorScheme["tabbar"].GetBackground()
	activetabForeground = config.ColorScheme["active-tab"].GetForeground()
	activetabBackground = config.ColorScheme["active-tab"].GetBackground()
	inactivetabForeground = config.ColorScheme["inactive-tab"].GetForeground()
	inactivetabBackground = config.ColorScheme["inactive-tab"].GetBackground()

	tabsStyle = lipgloss.NewStyle().Background(tabBackground).Foreground(tabForeground).BorderBackground(tabBackground).BorderForeground(tabForeground)
	tabStyle = tabsStyle.Copy()
	activeTabStyle = tabStyle.Copy().Border(activetabBorder).BorderTop(false).BorderBottom(false).BorderLeft(true).BorderRight(true).Bold(true).
		Background(activetabBackground).Foreground(activetabForeground).BorderBackground(activetabBackground).BorderForeground(activetabForeground)
	inactiveTabStyle = tabStyle.Copy().Border(inactivetabBorder).BorderTop(false).BorderBottom(false).BorderLeft(false).BorderRight(false).
		Background(inactivetabBackground).Foreground(inactivetabForeground).BorderBackground(inactivetabBackground).BorderForeground(inactivetabForeground)
	paneStyle = lipgloss.NewStyle().Background(paneBackground)
	inactivePaneStyle = lipgloss.NewStyle()
	activePaneStyle = inactivePaneStyle.Copy()
}

func NewPane() pane {
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

	UpdateStyling()

	switch msg.(type) {
	case events.CloseTabMsg:
		if len(p.tabs) > 1 {
			p.tabs = append(p.tabs[:p.active], p.tabs[p.active+1:]...)
			p.tabContent = append(p.tabContent[:p.active], p.tabContent[p.active+1:]...)
			if len(p.tabs) > 0 {
				p.active = (util.Max(0, p.active-1)) % len(p.tabs)
			}
		} else {
			cmds = append(cmds, events.Actions["CloseSplit"](msg))
		}
	case events.PrevTabMsg:
		p.tabContent[p.active].SetActive(false)
		p.active = (p.active - 1 + len(p.tabs)) % len(p.tabs)
		p.tabContent[p.active].SetActive(true)
		cmds = append(cmds, func() tea.Msg {
			return ""
		})
	case events.NextTabMsg:
		p.tabContent[p.active].SetActive(false)
		p.active = (p.active + 1) % len(p.tabs)
		p.tabContent[p.active].SetActive(true)
		cmds = append(cmds, func() tea.Msg {
			return ""
		})
	case events.NewTabMsg:
		c := NewCode()
		cmd := c.OpenFile("internal/view/view.go")
		p.tabContent = append(p.tabContent, c)
		p.tabs = append(p.tabs, "internal/view/view.go")
		cmds = append(cmds, cmd)
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
	tw := 0
	for i := 0; i < len(p.tabs); i++ {
		cts := ""
		if i == p.active {
			cts = activeTabStyle.Render(p.tabs[i])
		} else {
			if i-1 != p.active {
				cts = inactiveTabStyle.BorderLeft(true).Render(p.tabs[i])
			} else {
				cts = inactiveTabStyle.BorderLeft(false).Render(p.tabs[i])
			}
		}
		ctw := lipgloss.Width(cts)
		if tw+ctw >= w {
			cts = inactiveTabStyle.Render(p.tabs[i][:w-tw-2])
			ctw = lipgloss.Width(cts)
		}
		tw += ctw
		ts = append(ts, cts)
		if tw >= w {
			break
		}
	}
	tabs := tabsStyle.Width(w).Height(1).Render(lipgloss.JoinHorizontal(lipgloss.Top, ts...))

	if len(p.tabs) > 0 {
		content := paneStyle.Width(w).Height(h - 1).Render(p.tabContent[p.active].ViewWH(w, h-1))
		return lipgloss.JoinVertical(lipgloss.Left, tabs, content)
	} else {
		content := paneStyle.Width(w).Height(h - 1).Render("")
		return lipgloss.JoinVertical(lipgloss.Left, tabs, content)
	}
}

func (p pane) SetActive(b bool) elem {
	if len(p.tabs) > 0 {
		p.tabContent[p.active] = p.tabContent[p.active].SetActive(b)
	}
	return p
}
