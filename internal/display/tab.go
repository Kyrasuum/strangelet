package display

import (
	"bytes"
	"fmt"
	"math"
	"regexp"
	"strconv"

	buff "strangelet/internal/buffer"
	"strangelet/internal/config"
	"strangelet/internal/cursor"
	"strangelet/internal/util"
	"strangelet/pkg/app"

	"github.com/Kyrasuum/cview"
	"github.com/gdamore/tcell/v2"
)

var (
	gutterBGColor    tcell.Color
	gutterFGColor    tcell.Color
	statusBarBGColor tcell.Color
	statusBarFGColor tcell.Color

	formatParser = regexp.MustCompile(`\$\(.+?\)`)
)

var statusInfo = map[string]func(*tab) string{
	//support the following as values in the status format string
	"filename": func(t *tab) string {
		return t.GetName()
	},
	"line": func(t *tab) string {
		return strconv.Itoa(t.GetActiveCursor().Y + 1)
	},
	"col": func(t *tab) string {
		return strconv.Itoa(t.GetActiveCursor().X + 1)
	},
	"modified": func(t *tab) string {
		if t.Buffer.Modified() {
			return "+ "
		}
		if t.Type.Readonly {
			return "[ro] "
		}
		return ""
	},
}

// Overall tab display struct
type tab struct {
	*buff.Buffer
	*cview.Flex
	row *cview.Flex

	name  string
	label string

	dbuffer *buffer
	gutter  *gutter
	status  *statusBar

	cursors     []*cursor.Cursor
	curCursor   int
	StartCursor cursor.Loc
}

func NewTab(tabs *cview.TabbedPanels, b *buff.Buffer) (t *tab) {
	t = &tab{}
	t.Buffer = b

	t.Flex = cview.NewFlex()
	t.Flex.SetDirection(cview.FlexRow)
	t.row = cview.NewFlex()

	t.gutter = NewGutter(t.row)
	t.dbuffer = NewBuffer(t.row)
	t.Flex.AddItem(t.row, 0, 1, false)
	t.status = NewStatusBar(t.Flex)

	t.name = fmt.Sprintf(b.GetName())
	t.label = fmt.Sprintf(b.GetName())

	tabs.AddTab(t.name, t.label, t)
	tabs.SetCurrentTab(t.name)

	t.gutter.SetMouseCapture(t.wrapGutterMouse)

	return t
}

func (tab *tab) wrapGutterMouse(action cview.MouseAction, event *tcell.EventMouse) (cview.MouseAction, *tcell.EventMouse) {
	//pass scroll wheel actions to main buffer display
	switch action {
	case cview.MouseScrollUp:
		offx, offy := tab.dbuffer.GetScrollOffset()
		tab.dbuffer.ScrollTo(offx-1, offy)
	case cview.MouseScrollDown:
		offx, offy := tab.dbuffer.GetScrollOffset()
		tab.dbuffer.ScrollTo(offx+1, offy)
	}

	return action, event
}

func (tab *tab) Render(scr tcell.Screen) bool {
	//limit some redrawing to only when needed
	if tab.Buffer.ModifiedThisFrame {
		//update displayed buffer text
		tab.updateBufferDisplay()

		//check diff gutter
		if tab.Buffer.Settings["diffgutter"].(bool) {
			tab.Buffer.UpdateDiff(func(synchronous bool) {
				if !synchronous {
					app.CurApp.Redraw()
				}
			})
		}
		tab.Buffer.ModifiedThisFrame = false
	}
	//make gutter scroll to display buffer
	tab.gutter.ScrollTo(tab.dbuffer.GetScrollOffset())

	//update status bar displayed text
	tab.updateStatusBarDisplay()

	//done
	return false
}

func (tab *tab) updateBufferDisplay() {
	num_lines := tab.Buffer.LinesNum()

	//resize gutter
	gutter_size := (int)(math.Log10(float64(num_lines))+1) + 1
	tab.row.ResizeItem(tab.gutter, gutter_size, 1)

	//update displayed buffer text
	tab.dbuffer.SetBytes(nil)
	tab.gutter.TextView.SetBytes(nil)
	for line := 0; line < num_lines; line++ {
		//print line contents
		tab.dbuffer.Write(tab.Buffer.LineBytes(line))
		tab.dbuffer.Write([]byte("\n"))
		//print line numbers
		tab.gutter.TextView.Write([]byte(fmt.Sprintf("%d ", line)))
		tab.gutter.TextView.Write([]byte("\n"))
	}

	// _, _, width, height := tab.dbuffer.Box.GetRect()
	// _, _, width, height := tab.gutter.Box.GetRect()
	// for line := 0; line < height; line++ {
	// for col := 0; col < width; col++ {
	// }
	// }
}

func (tab *tab) updateStatusBarDisplay() {
	//clear status bar
	tab.status.SetBytes(nil)

	//find what character to use for dividers
	divchars := config.GetGlobalOption("divchars").(string)
	if util.CharacterCountInString(divchars) != 2 {
		divchars = "|-"
	}

	//define formatter for format string on status bar
	formatter := func(match []byte) []byte {
		name := match[2 : len(match)-1]
		if bytes.HasPrefix(name, []byte("opt")) {
			option := name[4:]
			return []byte(fmt.Sprint(tab.FindOpt(string(option))))
		} else if bytes.HasPrefix(name, []byte("bind")) {
			binding := string(name[5:])
			for k, v := range config.Bindings["buffer"] {
				if v == binding {
					return []byte(k)
				}
			}
			return []byte("null")
		} else {
			if fn, ok := statusInfo[string(name)]; ok {
				return []byte(fn(tab))
			}
			return []byte{}
		}
	}

	//grab formatted text for status bar
	leftText := []byte(tab.Buffer.Settings["statusformatl"].(string))
	leftText = formatParser.ReplaceAllFunc(leftText, formatter)
	rightText := []byte(tab.Buffer.Settings["statusformatr"].(string))
	rightText = formatParser.ReplaceAllFunc(rightText, formatter)

	//get size of parts of status bar
	leftLen := util.StringWidth(leftText, util.CharacterCount(leftText), 1)
	rightLen := util.StringWidth(rightText, util.CharacterCount(rightText), 1)

	//get width of status bar
	_, _, width, _ := tab.dbuffer.Box.GetRect()

	//print to status bar
	tab.status.Write(leftText)
	for i := leftLen; i < width-rightLen; i++ {
		tab.status.Write([]byte(" "))
	}
	tab.status.Write(rightText)
}

func (tab *tab) HandleInput(tevent *tcell.EventKey) (retEvent *tcell.EventKey) {
	return tevent
}

func (tab *tab) GetName() string {
	return tab.name
}

// FindOpt finds a given option in the current buffer's settings
func (tab *tab) FindOpt(opt string) interface{} {
	if val, ok := tab.Buffer.Settings[opt]; ok {
		return val
	}
	return "null"
}

func (tab *tab) GetActiveCursor() *cursor.Cursor {
	return &cursor.Cursor{Loc: cursor.Loc{X: 0, Y: 0}}
}

// Buffer display struct
type buffer struct {
	*cview.TextView
}

func NewBuffer(subFlex *cview.Flex) (buff *buffer) {
	buff = &buffer{}

	buff.TextView = cview.NewTextView()
	buff.TextView.SetTextAlign(cview.AlignLeft)
	buff.TextView.SetText("")
	buff.TextView.Box.SetBackgroundColor(tcell.NewRGBColor(20, 20, 20))
	buff.TextView.SetTextColor(tcell.NewRGBColor(230, 230, 230))

	subFlex.AddItem(buff, 0, 1, false)

	return buff
}

// Gutter display struct
type gutter struct {
	*cview.TextView
}

func NewGutter(subFlex *cview.Flex) (gutt *gutter) {
	gutt = &gutter{}

	if gutterBGColor == 0 {
		gutterBGColor = tcell.NewRGBColor(40, 40, 40)
	}
	if gutterFGColor == 0 {
		gutterFGColor = tcell.NewRGBColor(170, 170, 170)
	}
	gutt.TextView = cview.NewTextView()
	gutt.TextView.SetTextAlign(cview.AlignRight)
	gutt.TextView.Box.SetBackgroundColor(gutterBGColor)
	gutt.TextView.SetTextColor(gutterFGColor)
	gutt.TextView.SetScrollBarVisibility(cview.ScrollBarNever)

	subFlex.AddItem(gutt, 2, 1, false)

	return gutt
}

// Statusbar display struct
type statusBar struct {
	*cview.TextView
}

func NewStatusBar(flex *cview.Flex) (stat *statusBar) {
	stat = &statusBar{}

	if statusBarBGColor == 0 {
		statusBarBGColor = tcell.NewRGBColor(160, 160, 160)
	}
	if statusBarFGColor == 0 {
		statusBarFGColor = tcell.NewRGBColor(20, 20, 20)
	}

	stat.TextView = cview.NewTextView()
	stat.TextView.SetTextAlign(cview.AlignLeft)
	stat.TextView.SetText("Status Bar")
	stat.TextView.Box.SetBackgroundColor(statusBarBGColor)
	stat.TextView.SetTextColor(statusBarFGColor)

	flex.AddItem(stat, 1, 1, false)

	return stat
}
