package view

import (
	pub "strangelet/pkg/app"

	tea "github.com/charmbracelet/bubbletea"
	lipgloss "github.com/charmbracelet/lipgloss"
)

type pane struct {
	active     int
	Tabs       []string
	TabContent []tea.Model
}

var (
	paneStyle         = lipgloss.NewStyle()
	inactivePaneStyle = lipgloss.NewStyle().BorderStyle(lipgloss.HiddenBorder())
	activePaneStyle   = inactiveTabStyle.Copy().BorderStyle(lipgloss.NormalBorder()).BorderForeground(highlightColor).Padding(0, 1)
)

const ()

func NewPane(app pub.App) pane {
	return pane{}
}

func (p pane) Init() tea.Cmd {
	return nil
}

func (p pane) Update(msg tea.Msg) (tea.Model, tea.Cmd)    { return p.UpdateTyped(msg) }
func (p pane) UpdateI(msg tea.Msg) (interface{}, tea.Cmd) { return p.UpdateTyped(msg) }
func (p pane) UpdateTyped(msg tea.Msg) (pane, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+t":
		case "tab":
			p.active = (p.active + 1) % len(p.Tabs)
		}
	case tea.MouseMsg:
		// tea.MouseEvent(msg)
	}

	// Return the updated model to the Bubble Tea runtime for processing.
	// Note that we're not returning a command.
	return p, nil
}

func (p pane) View() string { return p.ViewWH(0, 0) }
func (p pane) ViewWH(w, h int) string {
	return paneStyle.Width(w).Height(h).Render("Pane")
}
