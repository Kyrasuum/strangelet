package view

import (
	"fmt"
	logger "log"
	"strings"

	"strangelet/internal/config"
	"strangelet/internal/util"

	"github.com/charmbracelet/bubbles/textarea"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/knipferrc/teacup/dirfs"
	"github.com/zyedidia/highlight"
)

var (
	codeStyle      = paneStyle.Copy().Foreground(lipgloss.Color("#F8F8F2"))
	statementStyle = codeStyle.Copy().Foreground(lipgloss.Color("#F92672"))
	preprocStyle   = codeStyle.Copy().Foreground(lipgloss.Color("#CB4B16"))
	specialStyle   = codeStyle.Copy().Foreground(lipgloss.Color("#A6E22E"))
	stringStyle    = codeStyle.Copy().Foreground(lipgloss.Color("#E6DB74"))
	charStyle      = codeStyle.Copy().Foreground(lipgloss.Color("#BDE6AD"))
	typeStyle      = codeStyle.Copy().Foreground(lipgloss.Color("#66D9EF"))
	numberStyle    = codeStyle.Copy().Foreground(lipgloss.Color("#AE81FF"))
	commentStyle   = codeStyle.Copy().Foreground(lipgloss.Color("#75715E"))

// symbol.bracket: symbol.bracket
// brightblue: brightblue
// comment.bright: comment.bright
// constant: constant
// keyword: keyword
// constant.bool.false: constant.bool.false
// preproc.shebang: preproc.shebang
// symbol.tag: symbol.tag
// magenta: magenta
// error: error
// yellow: yellow
// default: default
// type: type
// indent-char.whitespace: indent-char.whitespace
// brightcyan: brightcyan
// brightwhite: brightwhite
// comment: comment
// type.keyword: type.keyword
// constant.number: constant.number
// symbol.operator: symbol.operator
// type.extended: type.extended
// constant.macro: constant.macro
// bold default: bold default
// special: special
// identifier.var: identifier.var
// brightyellow: brightyellow
// operator: operator
// statement: statement
// identifier.macro: identifier.macro
// black: black
// brightmagenta: brightmagenta
// symbol.tag.extended: symbol.tag.extended
// ignore: ignore
// todo: todo
// preproc: preproc
// symbol.brackets: symbol.brackets
// constant.string.url: constant.string.url
// identifier.micro: identifier.micro
// constant.specialChar: constant.specialChar
// symbol: symbol
// constant.string: constant.string
// brightgreen: brightgreen
// red: red
// green: green
// cyan: cyan
// constant.comment: constant.comment
// identifier.class: identifier.class
// underlined: underlined
// blue: blue
// constant.bool.true: constant.bool.true
// brightred: brightred
// indent-char: indent-char
// brightblack: brightblack
// constant-string: constant-string
// identifier: identifier
// constant.bool: constant.bool
)

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

	for _, line := range c.Content {
		display = append(display, codeStyle.Render(fmt.Sprintf("%s", line.text)))
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
			switch group {
			case highlight.Groups["statement"]:
				text += statementStyle.Render(fmt.Sprintf("%s", []byte{line.text[j]}))
			case highlight.Groups["preproc"]:
				text += preprocStyle.Render(fmt.Sprintf("%s", []byte{line.text[j]}))
			case highlight.Groups["special"]:
				text += specialStyle.Render(fmt.Sprintf("%s", []byte{line.text[j]}))
			case highlight.Groups["constant.string"]:
				text += stringStyle.Render(fmt.Sprintf("%s", []byte{line.text[j]}))
			case highlight.Groups["constant.specialChar"]:
				text += charStyle.Render(fmt.Sprintf("%s", []byte{line.text[j]}))
			case highlight.Groups["type"]:
				text += typeStyle.Render(fmt.Sprintf("%s", []byte{line.text[j]}))
			case highlight.Groups["constant.number"]:
				text += numberStyle.Render(fmt.Sprintf("%s", []byte{line.text[j]}))
			case highlight.Groups["comment"]:
				text += commentStyle.Render(fmt.Sprintf("%s", []byte{line.text[j]}))
			default:
				text += codeStyle.Render(fmt.Sprintf("%s", []byte{line.text[j]}))
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
