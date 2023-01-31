package view

import (
	"math"

	config "strangelet/internal/config"
	events "strangelet/internal/events"
	util "strangelet/internal/util"
	pub "strangelet/pkg/app"

	tea "github.com/charmbracelet/bubbletea"
	lipgloss "github.com/charmbracelet/lipgloss"
)

type split struct {
	active    int
	direction int
	panes     []elem

	size float64
}

type elem interface {
	tea.Model
	ViewWH(int, int) string
	UpdateI(tea.Msg) (interface{}, tea.Cmd)
	SetActive(bool) elem
}

var (
	splitStyle = lipgloss.NewStyle()
)

const (
	horizontal int = iota
	vertical
)

func NewSplit(app pub.App) split {
	s := split{
		direction: vertical,
		active:    0,
		panes:     []elem{},
	}

	s.panes = append(s.panes, NewPane())

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
	case events.NewSplitMsg:
		s.panes = append(s.panes, NewPane())
	case events.CloseSplitMsg:
		if len(s.panes) > 1 {
			s.panes = append(s.panes[:s.active], s.panes[s.active+1:]...)
			s.active = (util.Max(0, s.active-1)) % len(s.panes)
		} else {
			cmds = append(cmds, events.Actions["Quit"](msg))
		}
	case events.PrevSplitMsg:
		s.panes[s.active].SetActive(false)
		s.active = (s.active - 1 + len(s.panes)) % len(s.panes)
		s.panes[s.active].SetActive(true)
		cmds = append(cmds, func() tea.Msg {
			return ""
		})
	case events.NextSplitMsg:
		s.panes[s.active].SetActive(false)
		s.active = (s.active + 1) % len(s.panes)
		s.panes[s.active].SetActive(true)
		cmds = append(cmds, func() tea.Msg {
			return ""
		})
	case tea.KeyMsg:
		if action, ok := config.Bindings["Split"][msg.String()]; ok {
			if handler, ok := events.Actions[action]; ok {
				cmds = append(cmds, handler(msg))
			}
		}
	case tea.MouseMsg:
		// tea.MouseEvent(msg)
		switch msg.Type {
		case tea.MouseWheelUp:
		case tea.MouseWheelDown:
		}
	}
	e, cmd := s.panes[s.active].Update(msg)
	s.panes[s.active] = e.(elem)
	cmds = append(cmds, cmd)

	// Return the updated model to the Bubble Tea runtime for processing.
	// Note that we're not returning a command.
	return s, tea.Batch(cmds...)
}

func (s split) View() string { return s.ViewWH(0, 0) }
func (s split) ViewWH(w, h int) string {
	pd := []string{}
	for i := 0; i < len(s.panes); i++ {
		switch s.direction {
		case horizontal:
			pd = append(pd, splitStyle.Render(s.panes[i].ViewWH(
				int(math.Round(float64(w-1*i%2)/float64(len(s.panes)))),
				h)))
		case vertical:
			pd = append(pd, splitStyle.Render(s.panes[i].ViewWH(
				w,
				int(math.Round(float64(h-1*i%2)/float64(len(s.panes)))))))
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

func (s split) SetActive(b bool) split {
	for i, p := range s.panes {
		s.panes[i] = p.SetActive(b)
	}
	return s
}
