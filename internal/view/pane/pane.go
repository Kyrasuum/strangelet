package pane

import (
	config "strangelet/internal/config"
	events "strangelet/internal/events"
	util "strangelet/internal/util"
	code "strangelet/internal/view/code"
	view "strangelet/pkg/view"

	tea "github.com/charmbracelet/bubbletea"
	lipgloss "github.com/charmbracelet/lipgloss"
)

type Pane struct {
	active int
	dirty  bool
	frame  string

	tabs       []string
	tabContent []view.Elem
	tabx       int
}

var ()

const ()

func NewPane() *Pane {
	p := Pane{
		active:     0,
		dirty:      true,
		frame:      "",
		tabs:       []string{},
		tabContent: []view.Elem{},
	}
	return &p
}

func (p Pane) Init() tea.Cmd {
	return nil
}

func (p *Pane) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var (
		cmds []tea.Cmd
	)
	switch msg := msg.(type) {
	case events.CloseTabMsg:
		if len(p.tabs) > 1 {
			p.tabs = append(p.tabs[:p.active], p.tabs[p.active+1:]...)
			p.tabContent = append(p.tabContent[:p.active], p.tabContent[p.active+1:]...)
			if len(p.tabs) > 0 {
				p.active = (util.Max(0, p.active-1)) % len(p.tabs)
			}
			cmd := p.tabContent[p.active].SetActive(true)
			cmds = append(cmds, cmd)
		} else {
			cmds = append(cmds, events.Actions["CloseSplit"](msg))
		}
		p.Redraw(lipgloss.Size(p.frame))
	case events.PrevTabMsg:
		if len(p.tabs) > 0 {
			cmd := p.tabContent[p.active].SetActive(false)
			cmds = append(cmds, cmd)

			p.active = (p.active - 1 + len(p.tabs)) % len(p.tabs)

			cmd = p.tabContent[p.active].SetActive(true)
			cmds = append(cmds, cmd)
			p.Redraw(lipgloss.Size(p.frame))
		}
	case events.NextTabMsg:
		if len(p.tabs) > 0 {
			cmd := p.tabContent[p.active].SetActive(false)
			cmds = append(cmds, cmd)

			p.active = (p.active + 1) % len(p.tabs)

			cmd = p.tabContent[p.active].SetActive(true)
			cmds = append(cmds, cmd)
			p.Redraw(lipgloss.Size(p.frame))
		}
	case events.NewTabMsg:
		c := code.NewCode()
		cmd := c.OpenFile(string(msg))
		p.tabContent = append(p.tabContent, c)
		p.tabs = append(p.tabs, string(msg))
		cmds = append(cmds, cmd)
		p.Redraw(lipgloss.Size(p.frame))
	}
	if len(p.tabs) > 0 {
		_, cmd := p.tabContent[p.active].Update(msg)
		if cmd != nil {
			p.Redraw(lipgloss.Size(p.frame))
			cmds = append(cmds, cmd)
		}
	}

	// Return the updated model to the Bubble Tea runtime for processing.
	// Note that we're not returning a command.
	return p, tea.Batch(cmds...)
}

func (p *Pane) Redraw(w, h int) {
	ts := []string{}
	tw := 0
	for i := 0; i < len(p.tabs); i++ {
		cts := ""
		if i == p.active {
			cts = config.ActiveTabStyle.Render(p.tabs[i])
		} else {
			if i-1 != p.active {
				cts = config.InactiveTabStyle.BorderLeft(true).Render(p.tabs[i])
			} else {
				cts = config.InactiveTabStyle.BorderLeft(false).Render(p.tabs[i])
			}
		}
		ctw := lipgloss.Width(cts)
		if tw+ctw >= w {
			cts = config.InactiveTabStyle.Render(p.tabs[i][:util.Max(0, w-tw-2)])
			ctw = lipgloss.Width(cts)
		}
		tw += ctw
		ts = append(ts, cts)
		if tw >= w {
			break
		}
	}
	tabs := config.TabsStyle.Width(w).Height(1).Render(lipgloss.JoinHorizontal(lipgloss.Top, ts...))

	if len(p.tabs) > 0 {
		content := config.PaneStyle.Width(w).Height(h - 1).Render(p.tabContent[p.active].ViewWH(w, h-1))
		p.frame = lipgloss.JoinVertical(lipgloss.Left, tabs, content)

	} else {
		content := config.PaneStyle.Width(w).Height(h - 1).Render("")
		p.frame = lipgloss.JoinVertical(lipgloss.Left, tabs, content)
	}
}

func (p *Pane) View() string { return p.ViewWH(0, 0) }
func (p *Pane) ViewWH(w, h int) string {
	if p.dirty || w != lipgloss.Width(p.frame) || h != lipgloss.Height(p.frame) {
		for i := 0; i < len(p.tabs); i++ {
			p.tabContent[i].SetDirty()
		}
		p.Redraw(w, h)
		p.dirty = false
	}
	return p.frame
}

func (p *Pane) SetActive(b bool) tea.Cmd {
	var cmd tea.Cmd = events.Actions["NOOP"]("")
	if len(p.tabs) > 0 {
		cmd = p.tabContent[p.active].SetActive(b)
		p.dirty = true
	}
	return cmd
}

func (p *Pane) SetDirty() {
	p.dirty = true
}
