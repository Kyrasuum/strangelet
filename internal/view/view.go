package view

import (
	"fmt"

	config "strangelet/internal/config"
	events "strangelet/internal/events"
	cmd "strangelet/internal/view/cmd"
	filebrowser "strangelet/internal/view/filebrowser"
	logWindow "strangelet/internal/view/logWindow"
	split "strangelet/internal/view/split"
	pub "strangelet/pkg/app"

	tea "github.com/charmbracelet/bubbletea"
	lipgloss "github.com/charmbracelet/lipgloss"
	statusbar "github.com/knipferrc/teacup/statusbar"
)

type view struct {
	active int

	width  int
	height int

	s  *split.Split
	fb *filebrowser.Filebrowser
	l  *logWindow.LogWindow
	c  *cmd.Cmd
	sb *statusbar.Bubble
}

var ()

func NewView(app pub.App) view {
	s := split.NewSplit(app)
	fb := filebrowser.NewFileBrowser(app)
	l := logWindow.NewLog(app)
	c := cmd.NewCmd(app)

	sbstyle := config.ColorScheme["statusline"]
	sbfg := fmt.Sprintf("%+v", sbstyle.GetForeground())
	sbbg := fmt.Sprintf("%+v", sbstyle.GetBackground())

	sbftstyle := config.ColorScheme["statusline.ft"]
	sbftfg := fmt.Sprintf("%+v", sbftstyle.GetForeground())
	sbftbg := fmt.Sprintf("%+v", sbftstyle.GetBackground())

	sbgitstyle := config.ColorScheme["statusline.git"]
	sbgitfg := fmt.Sprintf("%+v", sbgitstyle.GetForeground())
	sbgitbg := fmt.Sprintf("%+v", sbgitstyle.GetBackground())

	sbcursorstyle := config.ColorScheme["statusline.cursor"]
	sbcursorfg := fmt.Sprintf("%+v", sbcursorstyle.GetForeground())
	sbcursorbg := fmt.Sprintf("%+v", sbcursorstyle.GetBackground())

	sb := statusbar.New(
		statusbar.ColorConfig{
			Foreground: lipgloss.AdaptiveColor{Light: sbfg, Dark: sbfg},
			Background: lipgloss.AdaptiveColor{Light: sbbg, Dark: sbbg},
		},
		statusbar.ColorConfig{
			Foreground: lipgloss.AdaptiveColor{Light: sbftfg, Dark: sbftfg},
			Background: lipgloss.AdaptiveColor{Light: sbftbg, Dark: sbftbg},
		},
		statusbar.ColorConfig{
			Foreground: lipgloss.AdaptiveColor{Light: sbgitfg, Dark: sbgitfg},
			Background: lipgloss.AdaptiveColor{Light: sbgitbg, Dark: sbgitbg},
		},
		statusbar.ColorConfig{
			Foreground: lipgloss.AdaptiveColor{Light: sbcursorfg, Dark: sbcursorfg},
			Background: lipgloss.AdaptiveColor{Light: sbcursorbg, Dark: sbcursorbg},
		},
	)

	v := view{
		active: config.SplitView,

		width:  0,
		height: 0,

		s:  &s,
		fb: &fb,
		l:  &l,
		c:  &c,
		sb: &sb,
	}

	return v
}

func (v view) Init() tea.Cmd {
	return tea.Batch(tea.EnterAltScreen)
}

func (v view) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var (
		c    tea.Cmd
		cmds []tea.Cmd
	)

	switch msg := msg.(type) {
	// These keys should exit the program.
	case events.StatusMsg:
		msg = append(msg, []string{"", "", "", ""}...)
		v.sb.SetContent(msg[0], msg[1], msg[2], msg[3])
	case events.FocusCommand:
		if v.active != config.CmdView {
			e, c := v.s.SetActive(false)
			*v.s = e.(split.Split)
			cmds = append(cmds, c)

			v.active = config.CmdView

			e, c = v.c.SetActive(true)
			*v.c = e.(cmd.Cmd)
			cmds = append(cmds, c)
		} else {
			e, c := v.s.SetActive(true)
			*v.s = e.(split.Split)
			cmds = append(cmds, c)

			v.active = config.SplitView

			e, c = v.c.SetActive(false)
			*v.c = e.(cmd.Cmd)
			cmds = append(cmds, c)
		}
	case events.ToggleLogWindow:
		*v.l = v.l.ToggleVisible()
	case events.ToggleFileBrowser:
		*v.fb = v.fb.ToggleVisible()
		cmds = append(cmds, events.Actions["FocusFileBrowser"](msg))
	case events.FocusFileBrowser:
		if v.active != config.FilesView {
			e, c := v.s.SetActive(false)
			*v.s = e.(split.Split)
			cmds = append(cmds, c)

			v.active = config.FilesView

			e, c = v.fb.SetActive(true)
			*v.fb = e.(filebrowser.Filebrowser)
			cmds = append(cmds, c)
		} else {
			e, c := v.s.SetActive(true)
			*v.s = e.(split.Split)
			cmds = append(cmds, c)

			v.active = config.SplitView

			e, c = v.fb.SetActive(false)
			*v.fb = e.(filebrowser.Filebrowser)
			cmds = append(cmds, c)
		}
	case events.CloseApp:
		return v, tea.Quit
	case tea.WindowSizeMsg:
		v.sb.SetSize(msg.Width)
		v.width = msg.Width
		v.height = msg.Height

		c := v.c.ViewW(v.width)
		sb := v.sb.View()
		h := v.height - lipgloss.Height(c) - lipgloss.Height(sb)
		*v.fb = v.fb.SetHeight(h)
		*v.l = v.l.SetHeight(h)
	case tea.KeyMsg:
		if action, ok := config.Bindings["Global"][msg.String()]; ok {
			if handler, ok := events.Actions[action]; ok {
				cmds = append(cmds, handler(msg))
			}
		}
		str := msg.String()
		if str[0] == config.PasteBeginKey && str[len(str)-1] == config.PasteEndKey {
			cmds = append(cmds, events.Actions["Paste"](msg))
		}
	case tea.MouseMsg:
		if !config.GlobalSettings["mouse"].(bool) {
			break
		}
		switch msg.Type {
		case tea.MouseWheelUp:
		case tea.MouseWheelDown:
		}
	}

	switch v.active {
	case config.SplitView:
		*v.s, c = v.s.UpdateTyped(msg)
		cmds = append(cmds, c)
	case config.FilesView:
		*v.fb, c = v.fb.UpdateTyped(msg)
		cmds = append(cmds, c)
	case config.CmdView:
		*v.c, c = v.c.UpdateTyped(msg)
		cmds = append(cmds, c)
	}

	*v.l, c = v.l.UpdateTyped(msg)
	cmds = append(cmds, c)

	empty := true
	for _, cmd := range cmds {
		if cmd != nil {
			empty = false
		}
	}
	if empty {
		switch msg := msg.(type) {
		case tea.KeyMsg:
			cmds = append(cmds, func() tea.Msg {
				return events.LogMessage("Unknown Keybind for scope[Global, " + config.Scopes[v.active] + "]: " + msg.String())
			})
		}
	}

	// Return the updated model to the Bubble Tea runtime for processing.
	// Note that we're not returning a command.
	return v, tea.Batch(cmds...)
}

func (v view) View() string {
	//initial render
	c := v.c.ViewW(v.width)
	sb := v.sb.View()

	fb := ""
	l := ""

	if v.fb.Visible() {
		*v.fb, fb = v.fb.View()
	}
	if v.l.Visible() {
		l = v.l.View()
	}

	s := v.s.ViewWH(
		v.width-lipgloss.Width(fb)-lipgloss.Width(l),
		v.height-lipgloss.Height(c)-lipgloss.Height(sb))

	display := lipgloss.JoinVertical(lipgloss.Left, lipgloss.JoinHorizontal(lipgloss.Top, fb, s, l), sb, c)

	return display
}
