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
	dirty  bool
	frame  string

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
		dirty:  true,
		frame:  "",

		width:  0,
		height: 0,

		s:  s,
		fb: fb,
		l:  l,
		c:  c,
		sb: &sb,
	}

	return v
}

func (v *view) Init() tea.Cmd {
	return nil
}

func (v *view) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var (
		c    tea.Cmd
		cmds []tea.Cmd
	)

	switch msg := msg.(type) {
	// These keys should exit the program.
	case events.StatusMsg:
		msg = append(msg, []string{"", "", "", ""}...)
		v.sb.SetContent(msg[0], msg[1], msg[2], msg[3])
		v.Redraw()
	case events.FocusCommand:
		if v.active != config.CmdView {
			c := v.s.SetActive(false)
			cmds = append(cmds, c)

			v.active = config.CmdView

			c = v.c.SetActive(true)
			cmds = append(cmds, c)
		} else {
			c := v.s.SetActive(true)
			cmds = append(cmds, c)

			v.active = config.SplitView

			c = v.c.SetActive(false)
			cmds = append(cmds, c)
		}
		v.Redraw()
	case events.ToggleLogWindow:
		v.l.ToggleVisible()
		v.Redraw()
	case events.ToggleFileBrowser:
		if v.fb.Visible() {
			c := v.s.SetActive(true)
			cmds = append(cmds, c)
			v.active = config.SplitView
		} else {
			c := v.s.SetActive(false)
			cmds = append(cmds, c)
			v.active = config.FilesView
		}
		v.fb.ToggleVisible()
		v.Redraw()
	case events.FocusFileBrowser:
		if v.fb.Visible() {
			if v.active != config.FilesView {
				c := v.s.SetActive(false)
				cmds = append(cmds, c)

				v.active = config.FilesView

				c = v.fb.SetActive(true)
				cmds = append(cmds, c)
			} else {
				c := v.s.SetActive(true)
				cmds = append(cmds, c)

				v.active = config.SplitView

				c = v.fb.SetActive(false)
				cmds = append(cmds, c)
			}
			v.Redraw()
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
		v.fb.SetHeight(h)
		v.l.SetHeight(h)
		v.dirty = true
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
		_, c = v.s.Update(msg)
		if c != nil {
			cmds = append(cmds, c)
		}
		v.Redraw()
	case config.FilesView:
		_, c = v.fb.Update(msg)
		if c != nil {
			cmds = append(cmds, c)
		}
		v.Redraw()
	case config.CmdView:
		_, c = v.c.Update(msg)
		if c != nil {
			cmds = append(cmds, c)
		}
		v.Redraw()
	}

	_, c = v.l.Update(msg)
	if c != nil {
		v.dirty = true
		cmds = append(cmds, c)
		v.Redraw()
	}

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

func (v *view) Redraw() {
	c := v.c.ViewW(v.width)
	sb := v.sb.View()

	fb := ""
	l := ""

	if v.fb.Visible() {
		fb = v.fb.View()
	}
	if v.l.Visible() {
		l = v.l.View()
	}

	s := v.s.ViewWH(
		v.width-lipgloss.Width(fb)-lipgloss.Width(l),
		v.height-lipgloss.Height(c)-lipgloss.Height(sb))

	v.frame = lipgloss.JoinVertical(lipgloss.Left, lipgloss.JoinHorizontal(lipgloss.Top, fb, s, l), sb, c)
}

func (v *view) View() string {
	if v.dirty {
		//initial render
		v.s.SetDirty()
		v.c.SetDirty()
		v.fb.SetDirty()
		v.l.SetDirty()

		v.Redraw()
		v.dirty = false
	}

	return v.frame
}

func (v *view) SetDirty() {
	v.dirty = true
}
