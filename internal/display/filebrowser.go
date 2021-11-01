package display

import (
	"os"
	"path/filepath"

	ibuffer "strangelet/internal/buffer"
	"strangelet/pkg/app"

	"github.com/Kyrasuum/cview"
	"github.com/gdamore/tcell/v2"
)

var (
	dirColor        tcell.Color
	filColor        tcell.Color
	fbbgColor       tcell.Color
	fbheaderbgColor tcell.Color

	filebrowserW int = 30

	curFilebrowser *filebrowser
)

type filebrowser struct {
	*cview.Flex

	tree   *cview.TreeView
	header *cview.Table
	root   *cview.TreeNode

	parentFlex *cview.Flex

	rootDir string
}

func NewFilebrowser(subFlex *cview.Flex) (fb *filebrowser) {
	//enforce only one
	if curFilebrowser != nil {
		return curFilebrowser
	}
	fb = &filebrowser{}

	//init colors
	if dirColor == 0 {
		dirColor = tcell.NewRGBColor(220, 100, 100)
	}
	if filColor == 0 {
		filColor = tcell.NewRGBColor(220, 220, 220)
	}
	if fbbgColor == 0 {
		fbbgColor = tcell.NewRGBColor(30, 30, 30)
	}
	if fbheaderbgColor == 0 {
		fbheaderbgColor = tcell.NewRGBColor(150, 0, 0)
	}

	//setup tree
	fb.rootDir = ".."
	fb.root = cview.NewTreeNode(fb.rootDir)
	fb.root.SetColor(dirColor)
	fb.root.SetReference("../")
	fb.tree = cview.NewTreeView()
	fb.tree.SetRoot(fb.root)
	fb.tree.SetCurrentNode(fb.root)
	fb.tree.SetGraphics(false)
	fb.tree.SetScrollBarVisibility(cview.ScrollBarNever)
	fb.tree.Box.SetBackgroundColor(fbbgColor)

	// Add the current directory to the root node.
	fb.AddDirEntry(fb.root, ".")

	//setup header
	fb.header = cview.NewTable()
	cell := cview.NewTableCell("")
	cwd, err := os.Getwd()
	if err != nil {
		cwd = fb.rootDir
	}
	cell.SetText(cwd)
	cell.SetAlign(cview.AlignRight)
	fb.rootDir = cwd
	cell.SetAttributes(tcell.AttrUnderline)
	fb.header.SetCell(0, 0, cell)
	fb.header.Box.SetBackgroundColor(fbheaderbgColor)

	//setup flex
	fb.Flex = cview.NewFlex()
	fb.Flex.SetDirection(cview.FlexRow)
	fb.Flex.AddItem(fb.header, 1, 1, false)
	fb.Flex.AddItem(fb.tree, 0, 1, false)

	// If a directory was selected, open it.
	fb.tree.SetSelectedFunc(fb.OpenDirectory)

	subFlex.AddItem(fb, filebrowserW, 1, false)
	fb.parentFlex = subFlex
	curFilebrowser = fb
	// Default to closed
	fb.ToggleDisplay()

	return fb
}

func (fb *filebrowser) Render(scr tcell.Screen) bool {
	return false
}

func (fb *filebrowser) IsVisible() bool {
	return fb.Box.GetVisible()
}

// A helper function which adds the files and directories of the given path
// to the given target node.
func (fb *filebrowser) AddDirEntry(target *cview.TreeNode, path string) (err error) {
	files, err := os.ReadDir(path)
	if err != nil {
		return err
	}
	for _, file := range files {
		node := cview.NewTreeNode(file.Name())
		node.SetReference(filepath.Join(path, file.Name()))
		node.SetSelectable(true)
		if file.IsDir() {
			node.SetColor(dirColor)
		} else {
			node.SetColor(filColor)
		}
		target.AddChild(node)
	}
	return nil
}

// helper function to open a directory on the tree view
func (fb *filebrowser) OpenDirectory(node *cview.TreeNode) {
	reference := node.GetReference()
	if reference == nil || reference.(string) == "../" {
		return
	}
	children := node.GetChildren()
	if len(children) == 0 {
		// Load and show files in this directory.
		path := reference.(string)
		if err := fb.AddDirEntry(node, path); err != nil {
			//file
		} else {
			//directory
			node.SetExpanded(true)
		}
	} else {
		// Collapse if visible, expand if collapsed.
		node.ClearChildren()
		node.SetChildren(nil)
		node.SetExpanded(false)
	}
}

func (fb *filebrowser) ChangeRootDir(path string) {
	//remove old
	fb.root.ClearChildren()
	fb.root.SetChildren(nil)

	//go up a directory
	os.Chdir(path)
	rootDir := ".."
	root := cview.NewTreeNode(rootDir)
	root.SetColor(dirColor)
	root.SetReference("../")
	fb.tree.SetRoot(root)
	fb.tree.SetCurrentNode(root)
	fb.AddDirEntry(root, ".")

	//update header
	cell := fb.header.GetCell(0, 0)
	cwd, err := os.Getwd()
	if err != nil {
		cwd = fb.rootDir
	}
	cell.SetText(cwd)

	//check root nodes
	fb.rootDir = cwd
	fb.root = root
}

func (fb *filebrowser) HandleInput(tevent *tcell.EventKey) *tcell.EventKey {
	if tevent.Key() == tcell.KeyCtrlD {
		fb.ToggleDisplay()
		return nil
	}
	if fb.Box.GetVisible() {
		if tevent.Key() == tcell.KeyLeft {
			node := fb.tree.GetCurrentNode()
			if len(node.GetChildren()) != 0 {
				fb.OpenDirectory(node)
			} else {
				fb.tree.SetCurrentNode(node.GetParent())
			}
			return nil
		}
		if tevent.Key() == tcell.KeyRight {
			node := fb.tree.GetCurrentNode()
			fb.OpenDirectory(node)
			return nil
		}
		if tevent.Key() == tcell.KeyEnter {
			node := fb.tree.GetCurrentNode()
			path := node.GetReference().(string)
			_, err := os.ReadDir(path)
			if err != nil {
				//might be file
				buff, err := ibuffer.NewBufferFromFile(path, ibuffer.BTDefault)
				if err == nil {
					curDisplay.AddTabToCurrentPanel(buff)
				}
			} else {
				//directory
				fb.ChangeRootDir(node.GetReference().(string))
			}
			return nil
		}
		fb.tree.InputHandler()(tevent, app.CurApp.SetFocus)
	}
	return tevent
}

func (fb *filebrowser) ToggleDisplay() {
	if fb.Box.GetVisible() {
		fb.parentFlex.ResizeItem(fb, -1, 0)
		fb.Box.SetVisible(false)
		app.CurApp.SetFocus(fb)
	} else {
		fb.parentFlex.ResizeItem(fb, filebrowserW, 1)
		fb.Box.SetVisible(true)
		app.CurApp.SetFocus(nil)
	}
	app.CurApp.Redraw(func() {})
}

func (fb *filebrowser) HandleMouse(event *tcell.EventMouse, action cview.MouseAction) (*tcell.EventMouse, cview.MouseAction) {
	if fb.Box.GetVisible() {
		switch action {
		case cview.MouseLeftClick:
		case cview.MouseLeftDoubleClick:
		case cview.MouseRightClick:
		case cview.MouseRightDoubleClick:
		case cview.MouseMiddleClick:
		case cview.MouseMiddleDoubleClick:
		case cview.MouseScrollLeft:
			offx, offy := fb.tree.GetScrollOffset()
			fb.tree.SetScrollOffset(offx-1, offy)
			return nil, 0
		case cview.MouseScrollRight:
			offx, offy := fb.tree.GetScrollOffset()
			fb.tree.SetScrollOffset(offx+1, offy)
			return nil, 0
		case cview.MouseScrollUp:
			offx, offy := fb.tree.GetScrollOffset()
			fb.tree.SetScrollOffset(offx, offy-1)
			return nil, 0
		case cview.MouseScrollDown:
			offx, offy := fb.tree.GetScrollOffset()
			fb.tree.SetScrollOffset(offx, offy+1)
			return nil, 0
		}
	}
	return event, action
}
