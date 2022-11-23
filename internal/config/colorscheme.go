package config

import (
	"errors"
	"regexp"
	"strconv"
	"strings"

	"github.com/gdamore/tcell/v2"
)

// DefStyle is Micro's default style
var DefStyle tcell.Style = tcell.StyleDefault

// Colorscheme is the current colorscheme
var Colorscheme map[string]tcell.Style

// GetColor takes in a syntax group and returns the colorscheme's style for that group
func GetColor(color string) tcell.Style {
	st := DefStyle
	if color == "" {
		return st
	}
	groups := strings.Split(color, ".")
	if len(groups) > 1 {
		curGroup := ""
		for i, g := range groups {
			if i != 0 {
				curGroup += "."
			}
			curGroup += g
			if style, ok := Colorscheme[curGroup]; ok {
				st = style
			}
		}
	} else if style, ok := Colorscheme[color]; ok {
		st = style
	} else {
		st = StringToStyle(color)
	}

	return st
}

// ColorschemeExists checks if a given colorscheme exists
func ColorschemeExists(colorschemeName string) bool {
	return FindRuntimeFile(RTColorscheme, colorschemeName) != nil
}

// InitColorscheme picks and initializes the colorscheme when micro starts
func InitColorscheme() error {
	Colorscheme = make(map[string]tcell.Style)
	DefStyle = tcell.StyleDefault

	return LoadDefaultColorscheme()
}

// LoadDefaultColorscheme loads the default colorscheme from $(ConfigDir)/colorschemes
func LoadDefaultColorscheme() error {
	return LoadColorscheme(GlobalSettings["colorscheme"].(string))
}

// LoadColorscheme loads the given colorscheme from a directory
func LoadColorscheme(colorschemeName string) error {
	file := FindRuntimeFile(RTColorscheme, colorschemeName)
	if file == nil {
		return errors.New(colorschemeName + " is not a valid colorscheme")
	}
	if data, err := file.Data(); err != nil {
		return errors.New("Error loading colorscheme: " + err.Error())
	} else {
		Colorscheme, err = ParseColorscheme(string(data))
		if err != nil {
			return err
		}
	}
	return nil
}

// ParseColorscheme parses the text definition for a colorscheme and returns the corresponding object
// Colorschemes are made up of color-link statements linking a color group to a list of colors
// For example, color-link keyword (blue,red) makes all keywords have a blue foreground and
// red background
func ParseColorscheme(text string) (map[string]tcell.Style, error) {
	var err error
	parser := regexp.MustCompile(`color-link\s+(\S*)\s+"(.*)"`)

	lines := strings.Split(text, "\n")

	c := make(map[string]tcell.Style)

	for _, line := range lines {
		if strings.TrimSpace(line) == "" ||
			strings.TrimSpace(line)[0] == '#' {
			// Ignore this line
			continue
		}

		matches := parser.FindSubmatch([]byte(line))
		if len(matches) == 3 {
			link := string(matches[1])
			colors := string(matches[2])

			style := StringToStyle(colors)
			c[link] = style

			if link == "default" {
				DefStyle = style
			}
		} else {
			err = errors.New("Color-link statement is not valid: " + line)
		}
	}

	return c, err
}

// StringToStyle returns a style from a string
// The strings must be in the format "extra foregroundcolor,backgroundcolor"
// The 'extra' can be bold, reverse, italic or underline
func StringToStyle(str string) tcell.Style {
	var fg, bg string
	spaceSplit := strings.Split(str, " ")
	split := strings.Split(spaceSplit[len(spaceSplit)-1], ",")
	if len(split) > 1 {
		fg, bg = split[0], split[1]
	} else {
		fg = split[0]
	}
	fg = strings.TrimSpace(fg)
	bg = strings.TrimSpace(bg)

	var fgColor, bgColor tcell.Color
	var ok bool
	if fg == "" || fg == "default" {
		fgColor, _, _ = DefStyle.Decompose()
	} else {
		fgColor, ok = StringToColor(fg)
		if !ok {
			fgColor, _, _ = DefStyle.Decompose()
		}
	}
	if bg == "" || bg == "default" {
		_, bgColor, _ = DefStyle.Decompose()
	} else {
		bgColor, ok = StringToColor(bg)
		if !ok {
			_, bgColor, _ = DefStyle.Decompose()
		}
	}

	style := DefStyle.Foreground(fgColor).Background(bgColor)
	if strings.Contains(str, "bold") {
		style = style.Bold(true)
	}
	if strings.Contains(str, "italic") {
		style = style.Italic(true)
	}
	if strings.Contains(str, "reverse") {
		style = style.Reverse(true)
	}
	if strings.Contains(str, "underline") {
		style = style.Underline(true)
	}
	return style
}

// StringToColor returns a tcell color from a string representation of a color
// We accept either bright... or light... to mean the brighter version of a color
func StringToColor(str string) (tcell.Color, bool) {
	switch str {
	case "black":
		return tcell.ColorBlack, true
	case "red":
		return tcell.ColorMaroon, true
	case "green":
		return tcell.ColorGreen, true
	case "yellow":
		return tcell.ColorOlive, true
	case "blue":
		return tcell.ColorNavy, true
	case "magenta":
		return tcell.ColorPurple, true
	case "cyan":
		return tcell.ColorTeal, true
	case "white":
		return tcell.ColorSilver, true
	case "brightblack", "lightblack":
		return tcell.ColorGray, true
	case "brightred", "lightred":
		return tcell.ColorRed, true
	case "brightgreen", "lightgreen":
		return tcell.ColorLime, true
	case "brightyellow", "lightyellow":
		return tcell.ColorYellow, true
	case "brightblue", "lightblue":
		return tcell.ColorBlue, true
	case "brightmagenta", "lightmagenta":
		return tcell.ColorFuchsia, true
	case "brightcyan", "lightcyan":
		return tcell.ColorAqua, true
	case "brightwhite", "lightwhite":
		return tcell.ColorWhite, true
	case "default":
		return tcell.ColorDefault, true
	default:
		// Check if this is a 256 color
		if num, err := strconv.Atoi(str); err == nil {
			return GetColor256(num), true
		}
		// Check if this is a truecolor hex value
		if len(str) == 7 && str[0] == '#' {
			return tcell.GetColor(str), true
		}
		return tcell.ColorDefault, false
	}
}

// GetColor256 returns the tcell color for a number between 0 and 255
func GetColor256(color int) tcell.Color {
	if color == 0 {
		return tcell.ColorDefault
	}
	return tcell.PaletteColor(color)
}
