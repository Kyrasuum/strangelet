package cmd

import (
	config "strangelet/internal/config"
	events "strangelet/internal/events"
	pub "strangelet/pkg/app"

	tea "github.com/charmbracelet/bubbletea"
	lipgloss "github.com/charmbracelet/lipgloss"
)

type Cmd struct {
	active bool
}

var (
	cmdstyle = lipgloss.NewStyle()
)

const ()

func NewCmd(app pub.App) Cmd {
	return Cmd{
		active: false,
	}
}

func (c Cmd) Init() tea.Cmd {
	return nil
}

func (c Cmd) Update(msg tea.Msg) (tea.Model, tea.Cmd)    { return c.UpdateTyped(msg) }
func (c Cmd) UpdateI(msg tea.Msg) (interface{}, tea.Cmd) { return c.UpdateTyped(msg) }
func (c Cmd) UpdateTyped(msg tea.Msg) (Cmd, tea.Cmd) {
	var cmds []tea.Cmd

	switch msg := msg.(type) {
	case tea.KeyMsg:
		if action, ok := config.Bindings["Split"][msg.String()]; ok {
			if handler, ok := events.Actions[action]; ok {
				cmds = append(cmds, handler(msg))
			}
		}
	case tea.MouseMsg:
		// tea.MouseEvent(msg)
	}

	// Return the updated model to the Bubble Tea runtime for processing.
	// Note that we're not returning a command.
	return c, tea.Batch(cmds...)
}

func (c Cmd) View() string { return c.ViewW(0) }
func (c Cmd) ViewW(w int) string {
	return cmdstyle.Width(w).Render("Cmd")
}
func (c Cmd) ViewWH(w int, h int) string { return c.ViewW(w) }

func (c Cmd) SetActive(b bool) (interface{}, tea.Cmd) {
	c.active = b
	return c, events.Actions["NOOP"]("")
}
