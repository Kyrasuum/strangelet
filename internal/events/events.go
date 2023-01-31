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
type NextSplitMsg string
type PrevSplitMsg string

type SaveTabMsg string
type CloseTabMsg string
type NewTabMsg string
type NextTabMsg string
type PrevTabMsg string

type CursorUpMsg string
type CursorDownMsg string
type CursorLeftMsg string
type CursorRightMsg string

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
		"NextSplit":  func(_ tea.Msg) tea.Cmd { return func() tea.Msg { return NextSplitMsg("") } },
		"PrevSplit":  func(_ tea.Msg) tea.Cmd { return func() tea.Msg { return PrevSplitMsg("") } },

		"SaveTab":  func(_ tea.Msg) tea.Cmd { return func() tea.Msg { return SaveTabMsg("") } },
		"CloseTab": func(_ tea.Msg) tea.Cmd { return func() tea.Msg { return CloseTabMsg("") } },
		"NewTab":   func(_ tea.Msg) tea.Cmd { return func() tea.Msg { return NewTabMsg("") } },
		"NextTab":  func(_ tea.Msg) tea.Cmd { return func() tea.Msg { return NextTabMsg("") } },
		"PrevTab":  func(_ tea.Msg) tea.Cmd { return func() tea.Msg { return PrevTabMsg("") } },

		"CursorUp":    func(_ tea.Msg) tea.Cmd { return func() tea.Msg { return CursorUpMsg("") } },
		"CursorDown":  func(_ tea.Msg) tea.Cmd { return func() tea.Msg { return CursorDownMsg("") } },
		"CursorLeft":  func(_ tea.Msg) tea.Cmd { return func() tea.Msg { return CursorLeftMsg("") } },
		"CursorRight": func(_ tea.Msg) tea.Cmd { return func() tea.Msg { return CursorRightMsg("") } },

		"ToggleFileBrowser": func(_ tea.Msg) tea.Cmd { return func() tea.Msg { return ToggleFileBrowser("") } },
		"FocusFileBrowser":  func(_ tea.Msg) tea.Cmd { return func() tea.Msg { return FocusFileBrowser("") } },

		"ToggleLogWindow": func(_ tea.Msg) tea.Cmd { return func() tea.Msg { return ToggleLogWindow("") } },

		"FocusCommand": func(_ tea.Msg) tea.Cmd { return func() tea.Msg { return FocusCommand("") } },
	}

	return nil
}
