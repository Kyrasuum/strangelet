package config

import (
	"errors"
	"fmt"
	"regexp"
	"strings"

	lipgloss "github.com/charmbracelet/lipgloss"
	colorful "github.com/lucasb-eyer/go-colorful"
	highlight "github.com/zyedidia/highlight"
)

const ()

var (
	ColorScheme map[string]lipgloss.Style = make(map[string]lipgloss.Style)
	ColorGroups map[highlight.Group]string

	activetabBorder   lipgloss.Border
	inactivetabBorder lipgloss.Border
	outerHalfBorder   lipgloss.Border
	innerHalfBorder   lipgloss.Border

	paneBackground        lipgloss.TerminalColor
	tabForeground         lipgloss.TerminalColor
	tabBackground         lipgloss.TerminalColor
	activetabForeground   lipgloss.TerminalColor
	activetabBackground   lipgloss.TerminalColor
	inactivetabForeground lipgloss.TerminalColor
	inactivetabBackground lipgloss.TerminalColor

	TabsStyle         lipgloss.Style
	TabStyle          lipgloss.Style
	InactiveTabStyle  lipgloss.Style
	ActiveTabStyle    lipgloss.Style
	PaneStyle         lipgloss.Style
	InactivePaneStyle lipgloss.Style
	ActivePaneStyle   lipgloss.Style
)

func UpdateStyling() {
	activetabBorder = lipgloss.ThickBorder()
	inactivetabBorder = lipgloss.RoundedBorder()
	outerHalfBorder = lipgloss.Border{
		Top:         "▀",
		Bottom:      "▄",
		Left:        "▌",
		Right:       "▐",
		TopLeft:     "▛",
		TopRight:    "▜",
		BottomLeft:  "▙",
		BottomRight: "▟",
	}
	innerHalfBorder = lipgloss.Border{
		Top:         "▄",
		Bottom:      "▀",
		Left:        "▐",
		Right:       "▌",
		TopLeft:     "▗",
		TopRight:    "▖",
		BottomLeft:  "▝",
		BottomRight: "▘",
	}

	paneBackground = ColorScheme["background"].GetBackground()

	tabForeground = ColorScheme["tabbar"].GetForeground()
	tabBackground = ColorScheme["tabbar"].GetBackground()
	activetabForeground = ColorScheme["active-tab"].GetForeground()
	activetabBackground = ColorScheme["active-tab"].GetBackground()
	inactivetabForeground = ColorScheme["inactive-tab"].GetForeground()
	inactivetabBackground = ColorScheme["inactive-tab"].GetBackground()

	TabsStyle = lipgloss.NewStyle().Background(tabBackground).Foreground(tabForeground).BorderBackground(tabBackground).BorderForeground(tabForeground)
	TabStyle = TabsStyle.Copy()
	ActiveTabStyle = TabStyle.Copy().Border(activetabBorder).
		BorderTop(false).BorderBottom(false).BorderLeft(true).BorderRight(true).Bold(true).
		Background(activetabBackground).Foreground(activetabForeground).BorderBackground(activetabBackground).BorderForeground(activetabForeground)
	InactiveTabStyle = TabStyle.Copy().Border(inactivetabBorder).
		BorderTop(false).BorderBottom(false).BorderLeft(false).BorderRight(false).
		Background(inactivetabBackground).Foreground(inactivetabForeground).BorderBackground(inactivetabBackground).BorderForeground(inactivetabForeground)
	PaneStyle = lipgloss.NewStyle().Background(paneBackground)
	InactivePaneStyle = lipgloss.NewStyle()
	ActivePaneStyle = InactivePaneStyle.Copy()

	ColorScheme["log"] = ColorScheme["log"].Border(outerHalfBorder).
		BorderTop(false).BorderBottom(false).BorderLeft(true).BorderRight(false).
		BorderBackground(ColorScheme["log"].GetBackground()).BorderForeground(ColorScheme["log"].GetForeground())

	ColorScheme["filebrowser"] = ColorScheme["filebrowser"].Border(outerHalfBorder).
		BorderTop(false).BorderBottom(false).BorderLeft(false).BorderRight(true).
		BorderBackground(ColorScheme["filebrowser"].GetBackground()).BorderForeground(ColorScheme["filebrowser"].GetForeground())
	ColorScheme["filebrowser.cwd"] = ColorScheme["filebrowser.cwd"].Underline(true)
	ColorScheme["filebrowser.cwd-inactive"] = ColorScheme["filebrowser.cwd-inactive"].Underline(true)
}

// ColorschemeExists checks if a given colorscheme exists
func ColorschemeExists(colorschemeName string) bool {
	return FindRuntimeFile(RTColorscheme, colorschemeName) != nil
}

// InitColorscheme picks and initializes the colorscheme when micro starts
func InitColorscheme() error {
	initColorscheme()

	return LoadDefaultColorscheme()
}

func initColorscheme() {
	ColorScheme = make(map[string]lipgloss.Style)
	ColorGroups = make(map[highlight.Group]string)

	for name, grp := range highlight.Groups {
		ColorScheme[name] = lipgloss.NewStyle()
		ColorGroups[grp] = name
	}
}

// LoadDefaultColorscheme loads the default colorscheme from $(ConfigDir)/colorschemes
func LoadDefaultColorscheme() error {
	return loadColorscheme(GlobalSettings["colorscheme"].(string))
}

// LoadColorscheme loads the given colorscheme based of default
func LoadColorscheme(colorschemeName string) error {
	initColorscheme()

	err := LoadColorscheme(GlobalSettings["colorscheme"].(string))
	if err != nil {
		return err
	}
	return loadColorscheme(colorschemeName)
}

// loadColorscheme loads the given colorscheme from a directory
func loadColorscheme(colorschemeName string) error {
	file := FindRuntimeFile(RTColorscheme, colorschemeName)
	if file == nil {
		return errors.New(colorschemeName + " is not a valid colorscheme")
	}
	if data, err := file.Data(); err != nil {
		return errors.New("Error loading colorscheme: " + err.Error())
	} else {
		ColorScheme, err = ParseColorscheme(string(data))
		if err != nil {
			return err
		}
	}

	//load parent styles
	for key, style := range ColorScheme {
		if parstyle, ok := ColorScheme["background"]; ok {
			style = style.Inherit(parstyle)
		}
		if parstyle, ok := ColorScheme["default"]; ok {
			style = style.Inherit(parstyle)
		}
		parents := append(strings.Split(key, "."), strings.Split(key, "-")...)
		for i, _ := range parents[:len(parents)-1] {
			parent := strings.Join(parents[:i], ".")
			if parstyle, ok := ColorScheme[parent]; ok {
				style = style.Inherit(parstyle)
			}
		}
	}

	UpdateStyling()

	return nil
}

// ParseColorscheme parses the text definition for a colorscheme and returns the corresponding object
// Colorschemes are made up of color-link statements linking a color group to a list of colors
// For example, color-link keyword (blue,red) makes all keywords have a blue foreground and
// red background
func ParseColorscheme(text string) (map[string]lipgloss.Style, error) {
	var err error
	parser := regexp.MustCompile(`^"([\w\.\- ]*)" ?"?(\w*)?"? ?"?([\w#,]*)?"?$`)

	lines := strings.Split(text, "\n")

	c := make(map[string]lipgloss.Style)

	for key, style := range ColorScheme {
		c[key] = style
	}

	for _, line := range lines {
		if strings.TrimSpace(line) == "" ||
			strings.TrimSpace(line)[0] == '#' {
			// Ignore this line
			continue
		}

		matches := parser.FindSubmatch([]byte(line))
		if len(matches) == 4 {
			link := string(matches[1])
			mods := string(matches[2])
			colors := string(matches[3])

			style := lipgloss.NewStyle()
			style = StringToStyle(style, mods, colors)

			c[link] = style
		} else {
			err = errors.New("Color-link statement is not valid: " + line)
		}
	}

	//create inactive colors
	for key, style := range c {
		fg, err := colorful.Hex(fmt.Sprintf("%+v", style.GetForeground()))
		if err != nil {
			continue
		}
		bg, err := colorful.Hex(fmt.Sprintf("%+v", style.GetBackground()))
		if err != nil {
			c[key+"-inactive"] = style.Copy()
		}
		fg.R -= (fg.R - bg.R) / 2
		fg.G -= (fg.G - bg.G) / 2
		fg.B -= (fg.B - bg.B) / 2

		c[key+"-inactive"] = style.Copy().Foreground(lipgloss.Color(fg.Hex()))
	}

	return c, err
}

// StringToStyle returns a style from a string
// The strings must be in the format "extra foregroundcolor,backgroundcolor"
// The 'extra' can be bold, reverse, italic or underline
func StringToStyle(style lipgloss.Style, mstr string, cstr string) lipgloss.Style {
	var fg, bg string
	split := strings.Split(cstr, ",")
	if len(split) > 1 {
		fg, bg = split[0], split[1]
	} else {
		fg = split[0]
	}
	fg = strings.TrimSpace(fg)
	bg = strings.TrimSpace(bg)

	if fgColor, ok := StringToColor(fg); ok {
		style = style.Foreground(fgColor)
	}
	if bgColor, ok := StringToColor(bg); ok {
		style = style.Background(bgColor)
	}

	if strings.Contains(mstr, "bold") {
		style = style.Bold(true)
	}
	if strings.Contains(cstr, "italic") {
		style = style.Italic(true)
	}
	if strings.Contains(cstr, "reverse") {
		style = style.Reverse(true)
	}
	if strings.Contains(cstr, "underline") {
		style = style.Underline(true)
	}
	return style
}

// StringToColor returns a tcell color from a string representation of a color
// We accept either bright... or light... to mean the brighter version of a color
func StringToColor(str string) (lipgloss.Color, bool) {
	switch str {
	case "black":
		return lipgloss.Color("0"), true
	case "red":
		return lipgloss.Color("1"), true
	case "green":
		return lipgloss.Color("2"), true
	case "yellow":
		return lipgloss.Color("3"), true
	case "blue":
		return lipgloss.Color("4"), true
	case "magenta":
		return lipgloss.Color("5"), true
	case "cyan":
		return lipgloss.Color("6"), true
	case "white":
		return lipgloss.Color("7"), true
	case "brightblack", "lightblack":
		return lipgloss.Color("8"), true
	case "brightred", "lightred":
		return lipgloss.Color("9"), true
	case "brightgreen", "lightgreen":
		return lipgloss.Color("10"), true
	case "brightyellow", "lightyellow":
		return lipgloss.Color("11"), true
	case "brightblue", "lightblue":
		return lipgloss.Color("12"), true
	case "brightmagenta", "lightmagenta":
		return lipgloss.Color("13"), true
	case "brightcyan", "lightcyan":
		return lipgloss.Color("14"), true
	case "brightwhite", "lightwhite":
		return lipgloss.Color("15"), true
	default:
		// Check if this is a truecolor hex value
		if len(str) == 7 && str[0] == '#' {
			return lipgloss.Color(str), true
		}
		return lipgloss.Color(""), false
	}
}
