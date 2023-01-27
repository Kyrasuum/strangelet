package events

import (
	tea "github.com/charmbracelet/bubbletea"
)

var (
	Actions map[string]keyHandler
)

const ()

type keyHandler func(tea.Msg) tea.Cmd

type CloseApp string

type CloseSplitMsg string
type NewSplitMsg string

type SaveTabMsg string
type CloseTabMsg string
type NewTabMsg string

type ToggleFileBrowser string
type FocusFileBrowser string

type ToggleLogWindow string
type LogMessage string

type FocusCommand string

func InitActions() error {
	Actions = map[string]keyHandler{
		"Quit": func(_ tea.Msg) tea.Cmd { return func() tea.Msg { return CloseApp("") } },

		"CloseSplit": func(_ tea.Msg) tea.Cmd { return func() tea.Msg { return CloseSplitMsg("") } },
		"NewSplit":   func(_ tea.Msg) tea.Cmd { return func() tea.Msg { return NewSplitMsg("") } },

		"SaveTab":  func(_ tea.Msg) tea.Cmd { return func() tea.Msg { return SaveTabMsg("") } },
		"CloseTab": func(_ tea.Msg) tea.Cmd { return func() tea.Msg { return CloseTabMsg("") } },
		"NewTab":   func(_ tea.Msg) tea.Cmd { return func() tea.Msg { return NewTabMsg("") } },

		"ToggleFileBrowser": func(_ tea.Msg) tea.Cmd { return func() tea.Msg { return ToggleFileBrowser("") } },
		"FocusFileBrowser":  func(_ tea.Msg) tea.Cmd { return func() tea.Msg { return FocusFileBrowser("") } },

		"ToggleLogWindow": func(_ tea.Msg) tea.Cmd { return func() tea.Msg { return ToggleLogWindow("") } },

		"FocusCommand": func(_ tea.Msg) tea.Cmd { return func() tea.Msg { return FocusCommand("") } },
	}

	return nil
}
