package view

import (
	"fmt"
	logger "log"
	"strings"

	"strangelet/internal/config"
	"strangelet/internal/util"

	"github.com/charmbracelet/bubbles/textarea"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/knipferrc/teacup/dirfs"
	"github.com/zyedidia/highlight"
)

var ()

const ()

type line struct {
	text  string
	match highlight.LineMatch
	state highlight.State
	err   error
}

type Content map[int]*line

type code struct {
	Filename string
	Syntax   *highlight.Highlighter
	Content
}

type syntaxMsg Content
type lineErrorMsg struct {
	err  error
	line int
}

func (c *code) Highlight(content string, fileName string) (Content, error) {
	//split content into lines
	text := strings.Split(content, "\n")

	// Load the syntax definition
	syntaxDef := config.DetectType(fileName, []byte(text[0]))

	// Make a new highlighter from the definition
	c.Syntax = highlight.NewHighlighter(syntaxDef)

	// Create initial line information
	lines := Content{}
	for lineN, li := range text {
		lines[lineN] = new(line)
		lines[lineN].text = li
	}

	// Syntax highlighting calls
	c.Syntax.HighlightStates(lines)
	c.Syntax.HighlightMatches(lines, 0, len(lines))

	return lines, nil
}

func (c *code) OpenFile(filename string) tea.Cmd {
	return func() tea.Msg {
		content, err := dirfs.ReadFileContent(filename)
		if err != nil {
			return errorMsg(err)
		}

		lines, err := c.Highlight(content, filename)
		if err != nil {
			return errorMsg(err)
		}

		c.Filename = filename
		return syntaxMsg(lines)
	}
}

func NewCode() code {
	viewPort := textarea.New()

	viewPort.Prompt = ""
	viewPort.Placeholder = ""
	viewPort.ShowLineNumbers = true

	return code{
		Filename: "",
		Content:  map[int]*line{},
	}
}

func (c code) Init() tea.Cmd {
	return nil
}

func (c code) Update(msg tea.Msg) (tea.Model, tea.Cmd)    { return c.UpdateTyped(msg) }
func (c code) UpdateI(msg tea.Msg) (interface{}, tea.Cmd) { return c.UpdateTyped(msg) }
func (c code) UpdateTyped(msg tea.Msg) (code, tea.Cmd) {
	var (
		// cmd  tea.Cmd
		cmds []tea.Cmd
	)

	switch msg := msg.(type) {
	case syntaxMsg:
		for name, grp := range highlight.Groups {
			logger.Printf("%+v: %+v\n", name, grp)
		}

		c.Content = Content(msg)
		return c, nil
	case lineErrorMsg:
		line := c.Content[msg.line]
		line.err = msg.err
		return c, nil
	}

	return c, tea.Batch(cmds...)
}

func (c code) View() string {
	display := []string{}
	var group highlight.Group = highlight.Group(len(highlight.Groups))

	for _, line := range c.Content {
		text := ""
		for j := 0; j < len(line.text); j++ {
			if newgrp, ok := line.match[j]; ok {
				group = newgrp
			}
			if grp, ok := config.ColorGroups[group]; ok {
				if style, ok := config.ColorScheme[grp]; ok {
					//print using style group
					text += style.Render(fmt.Sprintf("%s", []byte{line.text[j]}))
				} else {
					//look for parent defined style
					style := config.ColorScheme["default"]
					parents := append(strings.Split(grp, "."), strings.Split(grp, "-")...)
					for i, _ := range parents[:len(parents)-1] {
						parent := strings.Join(parents[:i], ".")
						if parstyle, ok := config.ColorScheme[parent]; ok {
							style = parstyle
						}
					}
					text += style.Render(fmt.Sprintf("%s", []byte{line.text[j]}))
				}
			} else {
				//default to default style
				text += config.ColorScheme["default"].Render(fmt.Sprintf("%s", []byte{line.text[j]}))
			}
		}
		display = append(display, text)
	}

	return "" + strings.Join(display, "\n")
}
func (c code) ViewWH(w, h int) string {
	display := []string{}
	var group highlight.Group = highlight.Group(len(highlight.Groups))

	for i := 0; i < h; i++ {
		if i >= len(c.Content) {
			break
		}
		line := c.Content[i]
		text := ""
		for j := 0; j < util.Min(w, len(line.text)); j++ {
			if newgrp, ok := line.match[j]; ok {
				group = newgrp
			}
			if grp, ok := config.ColorGroups[group]; ok {
				text += config.ColorScheme[grp].Inherit(config.ColorScheme["background"]).Render(fmt.Sprintf("%s", []byte{line.text[j]}))
			} else {
				text += config.ColorScheme["default"].Inherit(config.ColorScheme["background"]).Render(fmt.Sprintf("%s", []byte{line.text[j]}))
			}
		}
		display = append(display, text)
	}

	return "" + strings.Join(display, "\n")
}

func (c code) Line(n int) string {
	return c.Content[n].text
}
func (c Content) Line(n int) string {
	return c[n].text
}

func (c code) LinesNum() int {
	return len(c.Content)
}
func (c Content) LinesNum() int {
	return len(c)
}

func (c code) State(lineN int) highlight.State {
	return c.Content[lineN].state
}
func (c Content) State(lineN int) highlight.State {
	return c[lineN].state
}

func (c code) SetState(lineN int, s highlight.State) {
	c.Content[lineN].state = s
}
func (c Content) SetState(lineN int, s highlight.State) {
	c[lineN].state = s
}

func (c code) SetMatch(lineN int, m highlight.LineMatch) {
	c.Content[lineN].match = m
}
func (c Content) SetMatch(lineN int, m highlight.LineMatch) {
	c[lineN].match = m
}
