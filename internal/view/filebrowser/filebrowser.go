package filebrowser

import (
	"fmt"
	"log"
	"os"
	"path/filepath"

	config "strangelet/internal/config"
	events "strangelet/internal/events"
	util "strangelet/internal/util"
	pub "strangelet/pkg/app"

	tea "github.com/charmbracelet/bubbletea"
	lipgloss "github.com/charmbracelet/lipgloss"
)

type CurrEntry struct {
	path  string
	entry os.DirEntry
	pari  int
}

type Filebrowser struct {
	visible bool
	active  bool
	dirty   bool
	frame   string

	root string
	sel  int
	scr  int
	max  int

	open map[string]bool
	curr *CurrEntry

	height int
}

var ()

const ()

func NewFileBrowser(app pub.App) *Filebrowser {
	path, err := os.Getwd()
	if err != nil {
		log.Println(err)
	}

	fb := Filebrowser{
		visible: false,
		active:  false,
		dirty:   true,
		frame:   "",
		root:    path,
		sel:     0,
		scr:     0,
		max:     0,
		open:    map[string]bool{},
		curr:    nil,
		height:  0,
	}

	return &fb
}

func (fb *Filebrowser) Init() tea.Cmd {
	return nil
}

func (fb *Filebrowser) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmds []tea.Cmd

	switch msg := msg.(type) {
	case events.CursorDownMsg:
		if fb.sel < fb.max {
			fb.sel++
		}
		cmds = append(cmds, events.Actions["NOOP"](""))
		fb.Redraw()
	case events.CursorUpMsg:
		if fb.sel > 0 {
			fb.sel--
		}
		cmds = append(cmds, events.Actions["NOOP"](""))
		fb.Redraw()
	case events.FbOpenFolderMsg:
		if fb.curr.entry.IsDir() {
			fb.open[fb.curr.path] = true
		}
		fb.Redraw()
	case events.FbCloseFolderMsg:
		if _, ok := fb.open[fb.curr.path]; ok {
			delete(fb.open, fb.curr.path)
		} else {
			fb.sel = fb.curr.pari
		}
		fb.Redraw()
	case events.FbEnterEntryMsg:
		if fb.curr.entry.IsDir() {
			fb.root = fb.curr.path
		} else {
			cmds = append(cmds, events.Actions["NewTab"](fb.curr.path), events.Actions["FocusFileBrowser"](""))
		}
		fb.Redraw()
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

func (fb *Filebrowser) renderHeader() string {
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

func (fb *Filebrowser) renderFiles(content string, path string, indent int, ch int) string {
	width := int(config.GlobalSettings["fbwidth"].(float64))

	files, err := os.ReadDir(path)
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

					abs, err := filepath.Abs(filepath.Join(path, file.Name()))
					if err == nil {
						if _, ok := fb.open[abs]; ok {
							prefix = prefix[:len(prefix)-2] + "- "
						}
					}
				} else {
					style = config.ColorScheme["filebrowser.file"]
				}
				if fb.sel == ch+fb.scr {
					style = config.ColorScheme["filebrowser.sel"]
					abs, err := filepath.Abs(filepath.Join(path, file.Name()))
					if err == nil {
						fb.curr = &CurrEntry{entry: file, path: abs, pari: 0}
					}
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

			abs, err := filepath.Abs(filepath.Join(path, file.Name()))
			if err == nil {
				if _, ok := fb.open[abs]; ok {
					content = fb.renderFiles(content, filepath.Join(path, file.Name()), indent+1, ch)
					sub := lipgloss.Height(content)
					if ch <= fb.sel-fb.scr && sub >= fb.sel-fb.scr && fb.curr.pari == 0 {
						fb.curr.pari = ch + fb.scr - 1
					}
					ch = sub
				}
			}
		}
	}

	return content
}

func (fb *Filebrowser) Redraw() {
	width := int(config.GlobalSettings["fbwidth"].(float64))
	fb.max = 0

	//render header
	header := fb.renderHeader()

	//get file content
	content := fb.renderFiles(header, fb.root, 1, 1)

	//render filebrowser
	fb.frame = config.ColorScheme["filebrowser"].
		Height(fb.height).
		Width(width).
		Render(content)
}

func (fb *Filebrowser) View() string {
	if fb.dirty {
		fb.Redraw()
		fb.dirty = false
	}

	return fb.frame
}

func (fb *Filebrowser) SetHeight(h int) {
	fb.height = h
	fb.dirty = true
}

func (fb *Filebrowser) ToggleVisible() {
	fb.visible = !fb.visible
	fb.active = fb.visible
	fb.dirty = true
}

func (fb *Filebrowser) SetActive(b bool) tea.Cmd {
	fb.active = b
	fb.dirty = true
	return events.Actions["NOOP"]("")
}

func (fb *Filebrowser) Visible() bool {
	return fb.visible
}

func (fb *Filebrowser) SetDirty() {
	fb.dirty = true
}
