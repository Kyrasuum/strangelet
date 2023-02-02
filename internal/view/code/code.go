package code

import (
	"fmt"
	"strings"

	buffer "strangelet/internal/buffer"
	cursor "strangelet/internal/cursor"
	events "strangelet/internal/events"

	textarea "github.com/charmbracelet/bubbles/textarea"
	tea "github.com/charmbracelet/bubbletea"
	dirfs "github.com/knipferrc/teacup/dirfs"
)

var ()

const ()

type Code struct {
	Filename string
	Content  buffer.Buffer

	Cursors []cursor.Cursor

	active bool
}

func (c *Code) OpenFile(filename string) tea.Cmd {
	return func() tea.Msg {
		content, err := dirfs.ReadFileContent(filename)
		if err != nil {
			return events.ErrorMsg(err)
		}

		c.Content, err = c.Content.Highlight(content, filename)
		if err != nil {
			return events.ErrorMsg(err)
		}

		pos := &cursor.Pos{X: 0, Y: 0}
		c.Cursors = []cursor.Cursor{cursor.Cursor{Begin: pos, End: pos}}
		c.Filename = filename

		return *c
	}
}

func NewCode() Code {
	viewPort := textarea.New()

	viewPort.Prompt = ""
	viewPort.Placeholder = ""
	viewPort.ShowLineNumbers = true

	c := Code{
		Filename: "",
		Content:  buffer.NewBuffer(),
		Cursors:  []cursor.Cursor{},

		active: true,
	}

	return c
}

func (c Code) Status() tea.Msg {
	ft := "ft:"
	if c.Content.Def != nil {
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

func (c Code) Init() tea.Cmd {
	return nil
}

func (c Code) Update(msg tea.Msg) (tea.Model, tea.Cmd)    { return c.UpdateTyped(msg) }
func (c Code) UpdateI(msg tea.Msg) (interface{}, tea.Cmd) { return c.UpdateTyped(msg) }
func (c Code) UpdateTyped(msg tea.Msg) (Code, tea.Cmd) {
	var (
		// cmd  tea.Cmd
		cmds []tea.Cmd
	)

	switch msg := msg.(type) {
	case Code:
		c = Code(msg)
		return c, c.Status
	}

	return c, tea.Batch(cmds...)
}

func (c Code) View() string {
	return c.ViewWH(0, 0)
}

func (c Code) ViewWH(w, h int) string {
	return c.Content.Render(w, h, c.active)
}

func (c Code) SetActive(b bool) (interface{}, tea.Cmd) {
	c.active = b
	return c, c.Status
}
