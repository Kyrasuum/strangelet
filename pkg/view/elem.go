package view

import (
	tea "github.com/charmbracelet/bubbletea"
)

type Elem interface {
	tea.Model
	ViewWH(int, int) string
	SetActive(bool) tea.Cmd
	SetDirty()
}

var ()

const ()
