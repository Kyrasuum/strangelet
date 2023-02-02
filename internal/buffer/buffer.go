package buffer

import (
	"fmt"
	"strings"

	config "strangelet/internal/config"
	util "strangelet/internal/util"

	highlight "github.com/zyedidia/highlight"
)

var ()

const ()

type Line struct {
	Text  string
	Match highlight.LineMatch
	State highlight.State
	Err   error
}

type Content map[int]*Line

func NewContent() *Content {
	c := Content{}
	return &c
}

func (c Content) SetMatch(lineN int, m highlight.LineMatch) {
	c[lineN].Match = m
}
func (c Content) SetState(lineN int, s highlight.State) {
	c[lineN].State = s
}
func (c Content) State(lineN int) highlight.State {
	return c[lineN].State
}
func (c Content) LinesNum() int {
	return len(c)
}
func (c Content) Line(n int) string {
	return c[n].Text
}

type Buffer struct {
	Lines *Content

	Syntax *highlight.Highlighter
	Def    *highlight.Def
}

func NewBuffer() Buffer {
	b := Buffer{
		Lines: NewContent(),
	}
	return b
}

func (b Buffer) LinesNum() int {
	return b.Lines.LinesNum()
}

func (b Buffer) Line(n int) string {
	return b.Lines.Line(n)
}

func (b Buffer) State(lineN int) highlight.State {
	return b.Lines.State(lineN)
}

func (b Buffer) SetState(lineN int, s highlight.State) {
	b.Lines.SetState(lineN, s)
}

func (b Buffer) SetMatch(lineN int, m highlight.LineMatch) {
	b.Lines.SetMatch(lineN, m)
}

func (b Buffer) Highlight(content string, fileName string) (Buffer, error) {
	//split content into lines
	text := strings.Split(content, "\n")

	// Load the syntax definition
	b.Def = config.DetectType(fileName, []byte(text[0]))

	// Make a new highlighter from the definition
	b.Syntax = highlight.NewHighlighter(b.Def)

	// Create initial line information
	b.Lines = new(Content)
	*b.Lines = map[int]*Line{}
	for lineN, li := range text {
		(*b.Lines)[lineN] = new(Line)
		(*b.Lines)[lineN].Text = li
	}

	// Syntax highlighting calls
	b.Syntax.HighlightStates(b.Lines)
	b.Syntax.HighlightMatches(b.Lines, 0, b.LinesNum())

	return b, nil
}

func (b Buffer) renderChar(char []byte, group highlight.Group, active bool) string {
	if grp, ok := config.ColorGroups[group]; ok {
		if active {
			//use active style
			if style, ok := config.ColorScheme[grp]; ok {
				//print using style group
				return style.Inherit(config.ColorScheme["background"]).Render(fmt.Sprintf("%s", char))
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
				return style.Inherit(config.ColorScheme["background"]).Render(fmt.Sprintf("%s", char))
			}
		} else {
			//use inactive style
			if style, ok := config.ColorScheme[grp+"-inactive"]; ok {
				return style.Inherit(config.ColorScheme["background-inactive"]).Render(fmt.Sprintf("%s", char))
			} else {
				//look for parent defined style
				style := config.ColorScheme["default-inactive"]
				parents := append(strings.Split(grp, "."), strings.Split(grp, "-")...)
				for i, _ := range parents[:len(parents)-1] {
					parent := strings.Join(parents[:i], ".")
					if parstyle, ok := config.ColorScheme[parent]; ok {
						style = parstyle
					}
				}
				return style.Inherit(config.ColorScheme["background-inactive"]).Render(fmt.Sprintf("%s", char))
			}
		}
	} else {
		//style does not exist use defaults
		if active {
			return config.ColorScheme["default"].Inherit(config.ColorScheme["background"]).Render(fmt.Sprintf("%s", char))
		} else {
			return config.ColorScheme["default-inactive"].Inherit(config.ColorScheme["background-inactive"]).Render(fmt.Sprintf("%s", char))
		}
	}
}

func (b Buffer) Render(w int, h int, active bool) string {
	display := []string{}
	var group highlight.Group = highlight.Group(len(highlight.Groups))

	for i := 0; i < h; i++ {
		if i >= b.LinesNum() {
			break
		}
		line := (*b.Lines)[i]
		text := ""
		cw := 0
		for j := 0; j < len(line.Text); j++ {
			//get syntax group
			if newgrp, ok := line.Match[j]; ok {
				group = newgrp
			}
			//avoiding doing logic if should skip
			if cw <= w {
				char := []byte{line.Text[j]}
				//get character width
				cw += util.StringWidth(char, 1, int(config.GlobalSettings["tabsize"].(float64)))
				//avoid printing if we will go past renderable area
				if cw <= w {
					//print tab as spaces (show tab as tabsize)
					if char[0] == byte('\t') {
						char = []byte(util.Spaces(int(config.GlobalSettings["tabsize"].(float64))))
					}
					text += b.renderChar(char, group, active)
				}
			}
		}
		display = append(display, text)
	}

	return "" + strings.Join(display, "\n")
}
