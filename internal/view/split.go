package view

import (
	"fmt"
	"math"

	pub "strangelet/pkg/app"

	tea "github.com/charmbracelet/bubbletea"
	lipgloss "github.com/charmbracelet/lipgloss"
)

type split struct {
	active    int
	direction int
	panes     []child

	lastw int
	lasth int
}

type elem interface {
	tea.Model
	ViewWH(int, int) string
	UpdateI(tea.Msg) (interface{}, tea.Cmd)
}

type child struct {
	elem
	size float64
}

var (
	splitStyle         = lipgloss.NewStyle()
	inactiveSplitStyle = lipgloss.NewStyle()
	activeSplitStyle   = inactiveSplitStyle.Copy()
)

const (
	horizontal int = iota
	vertical
)

func NewSplit(app pub.App) split {
	s := split{
		direction: vertical,
		active:    0,
		panes:     []child{},
	}

	s.panes = append(s.panes, child{size: 1, elem: NewPane(app)})

	return s
}

func (s split) Init() tea.Cmd {
	return nil
}

func (s split) Update(msg tea.Msg) (tea.Model, tea.Cmd)    { return s.UpdateTyped(msg) }
func (s split) UpdateI(msg tea.Msg) (interface{}, tea.Cmd) { return s.UpdateTyped(msg) }
func (s split) UpdateTyped(msg tea.Msg) (split, tea.Cmd) {
	var (
		cmd  tea.Cmd
		cmds []tea.Cmd
	)
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+n":
		}
	case tea.MouseMsg:
		// tea.MouseEvent(msg)
		switch msg.Type {
		case tea.MouseWheelUp:
		case tea.MouseWheelDown:
		}
	}
	e, cmd := s.panes[s.active].elem.Update(msg)
	s.panes[s.active].elem = e.(elem)
	cmds = append(cmds, cmd)

	// Return the updated model to the Bubble Tea runtime for processing.
	// Note that we're not returning a command.
	return s, tea.Batch(cmds...)
}

func (s split) View() string { return s.ViewWH(0, 0) }
func (s split) ViewWH(w, h int) string {
	s.lastw = w
	s.lasth = h

	pd := []string{}
	for i := 0; i < len(s.panes); i++ {
		if i == s.active {
			switch s.direction {
			case horizontal:
				pd = append(pd, activeSplitStyle.Render(fmt.Sprintf("%4s", s.panes[i].elem.ViewWH(
					int(math.Round(s.panes[i].size*float64(w))),
					h))))
			case vertical:
				pd = append(pd, activeSplitStyle.Render(fmt.Sprintf("%4s", s.panes[i].elem.ViewWH(
					w,
					int(math.Round(s.panes[i].size*float64(h)))))))
			}
		} else {
			switch s.direction {
			case horizontal:
				pd = append(pd, inactiveSplitStyle.Render(fmt.Sprintf("%4s", s.panes[i].elem.ViewWH(
					int(math.Round(s.panes[i].size*float64(w))),
					h))))
			case vertical:
				pd = append(pd, inactiveSplitStyle.Render(fmt.Sprintf("%4s", s.panes[i].elem.ViewWH(
					w,
					int(math.Round(s.panes[i].size*float64(h)))))))
			}
		}
	}

	display := ""
	switch s.direction {
	case horizontal:
		display += lipgloss.JoinHorizontal(lipgloss.Top, pd...)
	case vertical:
		display += lipgloss.JoinVertical(lipgloss.Left, pd...)
	}

	return splitStyle.Width(w).Height(h).Render(display)
}
