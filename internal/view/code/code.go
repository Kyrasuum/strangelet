package code

import (
	"fmt"
	"strings"

	buffer "strangelet/internal/buffer"
	cursor "strangelet/internal/cursor"
	events "strangelet/internal/events"

	textarea "github.com/charmbracelet/bubbles/textarea"
	tea "github.com/charmbracelet/bubbletea"
	lipgloss "github.com/charmbracelet/lipgloss"
)

var ()

const ()

type Code struct {
	Filename string
	Content  *buffer.Buffer

	Cursors []cursor.Cursor

	active bool
	dirty  bool
	frame  string
}

func (c *Code) OpenFile(filename string) tea.Cmd {
	return func() tea.Msg {
		var err events.ErrorMsg
		c.Content, err = buffer.OpenFile(filename)
		if err != nil {
			return err
		}

		pos := &cursor.Pos{X: 0, Y: 0}
		c.Cursors = []cursor.Cursor{cursor.Cursor{Begin: pos, End: pos}}
		c.Filename = filename
		c.dirty = true

		return *c
	}
}

func NewCode() *Code {
	viewPort := textarea.New()

	viewPort.Prompt = ""
	viewPort.Placeholder = ""
	viewPort.ShowLineNumbers = true

	c := Code{
		Filename: "",
		Content:  nil,
		Cursors:  []cursor.Cursor{},

		active: true,
		dirty:  true,
		frame:  "",
	}

	return &c
}

func (c *Code) Status() tea.Msg {
	ft := "ft:"
	if c.Content != nil && c.Content.Def != nil {
		ft += c.Content.Def.FileType
	}

	cs := []string{}
	for _, cursor := range c.Cursors {
		if cursor.Begin == cursor.End {
			cs = append(cs, fmt.Sprintf("(%d,%d)", cursor.Begin.X, cursor.Begin.Y))
		} else {
			cs = append(cs, fmt.Sprintf("(%d,%d->%d,%d)",
				cursor.Begin.X, cursor.Begin.Y, cursor.End.X, cursor.End.Y))
		}
	}
	sperc := ""
	if len(c.Cursors) > 0 {
		perc := float64(c.Cursors[len(c.Cursors)-1].End.Y) / float64(c.Content.LinesNum())
		sperc = fmt.Sprintf(" %d%%", int64(100*perc))
	}

	return events.StatusMsg([]string{ft, "errors: ", "git: ", "[" + strings.Join(cs, ",") + "]" + sperc})
}

func (c *Code) Init() tea.Cmd {
	return nil
}

func (c *Code) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var (
		// cmd  tea.Cmd
		cmds []tea.Cmd
	)

	switch msg.(type) {
	case Code:
		cmds = append(cmds, c.Status)
		c.Redraw(lipgloss.Size(c.frame))
	}

	return c, tea.Batch(cmds...)
}

func (c *Code) Redraw(w, h int) {
	c.frame = c.Content.Render(w, h, c.active)
}

func (c *Code) View() string {
	return c.ViewWH(lipgloss.Size(c.frame))
}

func (c *Code) ViewWH(w, h int) string {
	if c.dirty || w != lipgloss.Width(c.frame) || h != lipgloss.Height(c.frame) {
		c.Redraw(w, h)
		c.dirty = false
	}
	return c.frame
}

func (c *Code) SetActive(b bool) tea.Cmd {
	c.active = b
	c.dirty = true
	return c.Status
}

func (c *Code) SetDirty() {
	c.dirty = true
}
