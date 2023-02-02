package view

import (
	tea "github.com/charmbracelet/bubbletea"
)

type Elem interface {
	tea.Model
	ViewWH(int, int) string
	UpdateI(tea.Msg) (interface{}, tea.Cmd)
	SetActive(bool) (interface{}, tea.Cmd)
}

var ()

const ()
