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
	dirty  bool
	frame  string
}

var (
	cmdstyle = lipgloss.NewStyle()
)

const ()

func NewCmd(app pub.App) *Cmd {
	c := Cmd{
		active: false,
		dirty:  true,
		frame:  "",
	}
	return &c
}

func (c *Cmd) Init() tea.Cmd {
	return nil
}

func (c *Cmd) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
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

func (c *Cmd) Redraw(w int) {
	c.frame = cmdstyle.Width(w).Render("Cmd")
}

func (c *Cmd) View() string { return c.ViewW(lipgloss.Width(c.frame)) }
func (c *Cmd) ViewW(w int) string {
	if c.dirty || w != lipgloss.Width(c.frame) {
		c.Redraw(w)
		c.dirty = false
	}
	return c.frame
}
func (c *Cmd) ViewWH(w int, h int) string { return c.ViewW(w) }

func (c *Cmd) SetActive(b bool) tea.Cmd {
	c.active = b
	c.dirty = true
	return events.Actions["NOOP"]("")
}

func (c *Cmd) SetDirty() {
	c.dirty = true
}
