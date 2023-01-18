package view

import (
	logger "log"

	pub "strangelet/pkg/app"

	tea "github.com/charmbracelet/bubbletea"
	lipgloss "github.com/charmbracelet/lipgloss"
)

type errorMsg error

type log struct {
	visible bool
}

var (
	logstyle = lipgloss.NewStyle()
)

const ()

func NewLog(app pub.App) log {
	return log{}
}

func (l log) Init() tea.Cmd {
	return nil
}

func (l log) Update(msg tea.Msg) (tea.Model, tea.Cmd) { return l.UpdateTyped(msg) }
func (l log) UpdateTyped(msg tea.Msg) (log, tea.Cmd) {
	switch msg := msg.(type) {
	case errorMsg:
		logger.Println(string(msg.Error()))
	case tea.KeyMsg:
		switch msg.String() {
		}
	case tea.MouseMsg:
		// tea.MouseEvent(msg)
	}

	// Return the updated model to the Bubble Tea runtime for processing.
	// Note that we're not returning a command.
	return l, nil
}

func (l log) View() string { return l.ViewH(0) }
func (l log) ViewH(h int) string {
	return logstyle.Height(h).Render("Log")
}
