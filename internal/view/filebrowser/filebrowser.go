package filebrowser

import (
	"fmt"
	"log"
	"os"
	// "path/filepath"

	config "strangelet/internal/config"
	events "strangelet/internal/events"
	util "strangelet/internal/util"
	pub "strangelet/pkg/app"

	tea "github.com/charmbracelet/bubbletea"
	lipgloss "github.com/charmbracelet/lipgloss"
)

type Filebrowser struct {
	visible bool
	active  bool

	root string
	sel  int
	scr  int
	max  int

	height int
}

var ()

const ()

func NewFileBrowser(app pub.App) Filebrowser {
	path, err := os.Getwd()
	if err != nil {
		log.Println(err)
	}

	fb := Filebrowser{
		visible: false,
		active:  false,
		root:    path,
		sel:     0,
		scr:     0,
		max:     0,
	}

	return fb
}

func (fb Filebrowser) Init() tea.Cmd {
	return nil
}

func (fb Filebrowser) UpdateI(msg tea.Msg) (interface{}, tea.Cmd) { return fb.UpdateTyped(msg) }
func (fb Filebrowser) UpdateTyped(msg tea.Msg) (Filebrowser, tea.Cmd) {
	var cmds []tea.Cmd

	switch msg := msg.(type) {
	case events.CursorDownMsg:
		if fb.sel < fb.max {
			fb.sel++
		}
	case events.CursorUpMsg:
		if fb.sel > 0 {
			fb.sel--
		}
	case tea.KeyMsg:
		if action, ok := config.Bindings["Filebrowser"][msg.String()]; ok {
			if handler, ok := events.Actions[action]; ok {
				cmds = append(cmds, handler(msg))
			}
		}
	case tea.MouseMsg:
		// tea.MouseEvent(msg)
	}

	// Return the updated model to the Bubble Tea runtime for processing.
	// Note that we're not returning a command.
	return fb, tea.Batch(cmds...)
}

func (fb Filebrowser) renderHeader() string {
	width := int(config.GlobalSettings["fbwidth"].(float64))
	header := ""
	var style lipgloss.Style

	if fb.active {
		if fb.sel == 0 {
			style = config.ColorScheme["filebrowser.sel"]
		} else {
			style = config.ColorScheme["filebrowser.dir"]
		}
	} else {
		style = config.ColorScheme["filebrowser.dir-inactive"]
	}
	header = config.ColorScheme["filebrowser.cwd"].Render(fb.root[util.Max(len(fb.root)-width+3, 0):] + string(os.PathSeparator))
	header = lipgloss.JoinHorizontal(lipgloss.Top, header,
		style.Width(width-lipgloss.Width(header)).Render(".."))
	return header
}

func (fb Filebrowser) renderFiles(content string, indent int, ch int) (Filebrowser, string) {
	width := int(config.GlobalSettings["fbwidth"].(float64))

	files, err := os.ReadDir(fb.root)
	if err != nil {
		log.Println(err.Error())
		content = err.Error()
	} else {
		fb.max += len(files)
		for _, file := range files[fb.scr:] {
			if ch >= fb.height {
				break
			}
			//set prefix and style
			var style lipgloss.Style
			prefix := util.Spaces(2 * indent)
			if fb.active {
				if file.IsDir() {
					style = config.ColorScheme["filebrowser.dir"]
					prefix = prefix[:len(prefix)-2] + "+ "
				} else {
					style = config.ColorScheme["filebrowser.file"]
				}
				if fb.sel == ch+fb.scr {
					style = config.ColorScheme["filebrowser.sel"]
				}
			} else {
				if file.IsDir() {
					style = config.ColorScheme["filebrowser.dir-inactive"]
					prefix = prefix[:len(prefix)-2] + "+ "
				} else {
					style = config.ColorScheme["filebrowser.file-inactive"]
				}
			}

			//create line text
			line := fmt.Sprintf("%s%+v", prefix,
				file.Name()[:util.Min(width-len(prefix), len(file.Name()))])

			//join into content
			content = lipgloss.JoinVertical(lipgloss.Left, content,
				style.Width(width).Render(line))
			ch++
		}
	}

	return fb, content
}

func (fb Filebrowser) View() (Filebrowser, string) {
	width := int(config.GlobalSettings["fbwidth"].(float64))
	fb.max = 0

	//render header
	header := fb.renderHeader()

	//get file content
	fb, content := fb.renderFiles(header, 1, 1)

	//render filebrowser
	return fb, config.ColorScheme["filebrowser"].
		Height(fb.height).
		Width(width).
		Render(content)
}

func (fb Filebrowser) SetHeight(h int) Filebrowser {
	fb.height = h
	return fb
}

func (fb Filebrowser) ToggleVisible() Filebrowser {
	fb.visible = !fb.visible
	return fb
}

func (fb Filebrowser) SetActive(b bool) (interface{}, tea.Cmd) {
	fb.active = b
	return fb, events.Actions["NOOP"]("")
}

func (fb Filebrowser) Visible() bool {
	return fb.visible
}
