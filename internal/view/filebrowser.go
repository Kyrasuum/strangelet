package view

import (
	config "strangelet/internal/config"
	events "strangelet/internal/events"
	pub "strangelet/pkg/app"

	tea "github.com/charmbracelet/bubbletea"
	lipgloss "github.com/charmbracelet/lipgloss"
)

type filebrowser struct {
	visible bool
	active  bool

	height int
}

var (
	fbstyle = lipgloss.NewStyle()
)

const ()

func NewFileBrowser(app pub.App) filebrowser {
	return filebrowser{
		visible: false,
	}
}

func (fb filebrowser) Init() tea.Cmd {
	return nil
}

func (fb filebrowser) Update(msg tea.Msg) (tea.Model, tea.Cmd)    { return fb.UpdateTyped(msg) }
func (fb filebrowser) UpdateI(msg tea.Msg) (interface{}, tea.Cmd) { return fb.UpdateTyped(msg) }
func (fb filebrowser) UpdateTyped(msg tea.Msg) (filebrowser, tea.Cmd) {
	var cmds []tea.Cmd

	switch msg := msg.(type) {
	case tea.KeyMsg:
		if action, ok := config.Bindings["Filebrowser"][msg.String()]; ok {
			if handler, ok := events.Actions[action]; ok {
				cmds = append(cmds, handler(msg))
			}
		}
	case tea.MouseMsg:
		// tea.MouseEvent(msg)
	}

	// Return the updated model to the Bubble Tea runtime for processing.
	// Note that we're not returning a command.
	return fb, tea.Batch(cmds...)
}

func (fb filebrowser) View() string {
	return fbstyle.
		Height(fb.height).
		Width(int(config.GlobalSettings["fbwidth"].(float64))).
		Render("File Browser:")
}
func (fb filebrowser) ViewWH(w int, h int) string { return fb.View() }

func (fb filebrowser) SetHeight(h int) {
	fb.height = h
}

func (fb filebrowser) ToggleVisible() filebrowser {
	fb.visible = !fb.visible
	return fb
}

func (fb filebrowser) SetActive(b bool) elem {
	fb.active = b
	return fb
}
