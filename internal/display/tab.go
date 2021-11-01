package display

import (
	"bytes"
	"fmt"
	"math"
	"regexp"
	"strconv"
	"strings"

	buff "strangelet/internal/buffer"
	"strangelet/internal/config"
	"strangelet/internal/cursor"
	"strangelet/internal/util"
	"strangelet/pkg/app"

	"github.com/Kyrasuum/cview"
	"github.com/gdamore/tcell/v2"
)

var (
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

	t.StartCursor = cursor.Loc{X: 0, Y: 0}
	t.cursors = append(t.cursors, &cursor.Cursor{Loc: t.StartCursor, CurSelection: [2]cursor.Loc{t.StartCursor, t.StartCursor}})
	t.curCursor = 0

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

	//wrap settings
	if !b.Settings["softwrap"].(bool) {
		t.dbuffer.TextView.SetWrap(false)
	} else {
		t.dbuffer.TextView.SetWrap(true)
	}

	//setup buffer callback
	b.OptionCallback = func(option string, nativeValue interface{}) {
		if option == "softwrap" {
			if nativeValue.(bool) {
				t.dbuffer.TextView.SetWrap(false)
			} else {
				t.dbuffer.TextView.SetWrap(true)
			}
		}
		t.Buffer.ModifiedThisFrame = false
		app.CurApp.Redraw(func() {})
	}

	//schedule redraw
	app.CurApp.Redraw(func() {
		b.MarkModified(0, 0)
	})

	return t
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
					app.CurApp.Redraw(func() {})
				}
			})
		}
		tab.Buffer.ModifiedThisFrame = false
	}
	//make gutter scroll to display buffer
	tab.gutter.ScrollTo(tab.dbuffer.GetScrollOffset())

	//update status bar displayed text
	tab.updateStatusBarDisplay()

	//update rendering of cursor
	tab.updateCursorDisplay(scr)

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
	_, _, width, _ := tab.dbuffer.GetInnerRect()
	for line := 0; line < num_lines; line++ {
		//print line contents
		line_str := tab.Buffer.Line(line)
		tab.dbuffer.Write([]byte(fmt.Sprintf(`["%d"]%s[""]`, line, line_str)))
		tab.dbuffer.Write([]byte("\n"))
		//print line numbers
		prefix_len := math.Max(0, float64(gutter_size)-math.Log10(float64(line+1))-2)
		tab.gutter.Write([]byte(fmt.Sprintf(`["%d"]%s%d [""]`, line, strings.Repeat(" ", int(prefix_len)), line)))

		if tab.Buffer.Settings["softwrap"].(bool) {
			line_len := len(line_str)
			for li := 0; li <= int(math.Max(1.0, float64(line_len))/math.Max(1.0, float64(width))); li++ {
				tab.gutter.Write([]byte("\n"))
			}
		} else {
			tab.gutter.Write([]byte("\n"))
		}
	}
}

func (tab *tab) updateCursorDisplay(scr tcell.Screen) {
	//grab data for all cursors
	regions := []string{}
	x, y, _, _ := tab.dbuffer.GetInnerRect()
	row, column := tab.dbuffer.GetScrollOffset()

	//main cursor logic
	// mainCurs := tab.cursors[tab.curCursor]

	//each cursor logic
	for _, curs := range tab.cursors {
		regions = append(regions, fmt.Sprintf("%d", curs.Y))
		offx := curs.X + x - column
		offy := curs.Y + y - row
		if offx > x || offy >= y {
			scr.ShowCursor(offx, offy)
		} else {
			scr.HideCursor()
		}
	}

	//highlight line
	tab.dbuffer.Highlight(regions...)
	tab.gutter.Highlight(regions...)

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

func (tab *tab) HandleInput(tevent *tcell.EventKey) *tcell.EventKey {
	return tab.dbuffer.HandleInput(tevent)
}

func (tab *tab) HandleMouse(event *tcell.EventMouse, action cview.MouseAction) (*tcell.EventMouse, cview.MouseAction) {
	return tab.dbuffer.HandleMouse(event, action)
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
	return tab.cursors[tab.curCursor]
}
