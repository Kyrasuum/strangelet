package split

import (
	"math"

	config "strangelet/internal/config"
	events "strangelet/internal/events"
	util "strangelet/internal/util"
	pane "strangelet/internal/view/pane"
	pub "strangelet/pkg/app"
	view "strangelet/pkg/view"

	tea "github.com/charmbracelet/bubbletea"
	lipgloss "github.com/charmbracelet/lipgloss"
)

type Split struct {
	active    int
	direction int
	panes     []view.Elem

	size float64
}

var (
	splitStyle = lipgloss.NewStyle()
)

const (
	horizontal int = iota
	vertical
)

func NewSplit(app pub.App) Split {
	s := Split{
		direction: vertical,
		active:    0,
		panes:     []view.Elem{},
	}

	s.panes = append(s.panes, pane.NewPane())

	return s
}

func (s Split) Init() tea.Cmd {
	return nil
}

func (s Split) Update(msg tea.Msg) (tea.Model, tea.Cmd)    { return s.UpdateTyped(msg) }
func (s Split) UpdateI(msg tea.Msg) (interface{}, tea.Cmd) { return s.UpdateTyped(msg) }
func (s Split) UpdateTyped(msg tea.Msg) (Split, tea.Cmd) {
	var (
		cmd  tea.Cmd
		cmds []tea.Cmd
	)
	switch msg := msg.(type) {
	case events.NewSplitMsg:
		s.panes = append(s.panes, pane.NewPane())
	case events.CloseSplitMsg:
		if len(s.panes) > 1 {
			s.panes = append(s.panes[:s.active], s.panes[s.active+1:]...)
			s.active = (util.Max(0, s.active-1)) % len(s.panes)

			e, cmd := s.panes[s.active].SetActive(true)
			s.panes[s.active] = e.(view.Elem)
			cmds = append(cmds, cmd)
		} else {
			cmds = append(cmds, events.Actions["Quit"](msg))
		}
	case events.PrevSplitMsg:
		e, cmd := s.panes[s.active].SetActive(false)
		s.panes[s.active] = e.(view.Elem)
		cmds = append(cmds, cmd)

		s.active = (s.active - 1 + len(s.panes)) % len(s.panes)

		e, cmd = s.panes[s.active].SetActive(true)
		s.panes[s.active] = e.(view.Elem)
		cmds = append(cmds, cmd)
	case events.NextSplitMsg:
		e, cmd := s.panes[s.active].SetActive(false)
		s.panes[s.active] = e.(view.Elem)
		cmds = append(cmds, cmd)

		s.active = (s.active + 1) % len(s.panes)

		e, cmd = s.panes[s.active].SetActive(true)
		s.panes[s.active] = e.(view.Elem)
		cmds = append(cmds, cmd)
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
	s.panes[s.active] = e.(view.Elem)
	cmds = append(cmds, cmd)

	// Return the updated model to the Bubble Tea runtime for processing.
	// Note that we're not returning a command.
	return s, tea.Batch(cmds...)
}

func (s Split) View() string { return s.ViewWH(0, 0) }
func (s Split) ViewWH(w, h int) string {
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

func (s Split) SetActive(b bool) (interface{}, tea.Cmd) {
	var cmd tea.Cmd = nil
	if len(s.panes) > 0 {
		var e interface{}
		e, cmd = s.panes[s.active].SetActive(b)
		s.panes[s.active] = e.(view.Elem)
	}
	return s, cmd
}
