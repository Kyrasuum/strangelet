package view

import (
	pub "strangelet/pkg/app"

	tea "github.com/charmbracelet/bubbletea"
	lipgloss "github.com/charmbracelet/lipgloss"
)

type filebrowser struct {
	visible bool
}

var (
	fbstyle = lipgloss.NewStyle()
)

const ()

func NewFileBrowser(app pub.App) filebrowser {
	return filebrowser{}
}

func (fb filebrowser) Init() tea.Cmd {
	return nil
}

func (fb filebrowser) Update(msg tea.Msg) (tea.Model, tea.Cmd) { return fb.UpdateTyped(msg) }
func (fb filebrowser) UpdateTyped(msg tea.Msg) (filebrowser, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		}
	case tea.MouseMsg:
		// tea.MouseEvent(msg)
	}

	// Return the updated model to the Bubble Tea runtime for processing.
	// Note that we're not returning a command.
	return fb, nil
}

func (fb filebrowser) View() string { return fb.ViewH(0) }
func (fb filebrowser) ViewH(h int) string {
	return fbstyle.Height(h).Render("File Browser")
}
