package view

import (
	pub "strangelet/pkg/app"

	tea "github.com/charmbracelet/bubbletea"
)

var ()

type view struct {
}

func NewView(app pub.App) view {
	return view{}
}

func (v view) Init() tea.Cmd {
	return tea.Batch(tea.EnterAltScreen)
}

func (v view) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		// These keys should exit the program.
		case "ctrl+c", "esc", "q":
			return v, tea.Quit
		}
	case tea.MouseMsg:
		// tea.MouseEvent(msg)
	}

	// Return the updated model to the Bubble Tea runtime for processing.
	// Note that we're not returning a command.
	return v, nil
}

func (v view) View() string {
	return "Hello world"
}
