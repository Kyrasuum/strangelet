package view

import (
	config "strangelet/internal/config"
	events "strangelet/internal/events"
	pub "strangelet/pkg/app"

	tea "github.com/charmbracelet/bubbletea"
	lipgloss "github.com/charmbracelet/lipgloss"
)

type cmd struct{}

var (
	cmdstyle = lipgloss.NewStyle()
)

const ()

func NewCmd(app pub.App) cmd {
	return cmd{}
}

func (c cmd) Init() tea.Cmd {
	return nil
}

func (c cmd) Update(msg tea.Msg) (tea.Model, tea.Cmd) { return c.UpdateTyped(msg) }
func (c cmd) UpdateTyped(msg tea.Msg) (cmd, tea.Cmd) {
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

func (c cmd) View() string { return c.ViewW(0) }
func (c cmd) ViewW(w int) string {
	return cmdstyle.Width(w).Render("Cmd")
}
