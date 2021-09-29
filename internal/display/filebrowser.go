package display

import (
	"os"
	"path/filepath"

	"strangelet/pkg/app"

	"github.com/Kyrasuum/cview"
	"github.com/gdamore/tcell/v2"
)

var (
	dirColor  tcell.Color
	filColor  tcell.Color
	fbbgColor tcell.Color

	filebrowserW int = 30

	curFilebrowser *filebrowser
)

type filebrowser struct {
	*cview.TreeView
	root *cview.TreeNode

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

	//setup tree
	fb.rootDir = "."
	fb.root = cview.NewTreeNode(fb.rootDir)
	fb.root.SetColor(dirColor)
	fb.TreeView = cview.NewTreeView()
	fb.TreeView.SetRoot(fb.root)
	fb.TreeView.SetCurrentNode(fb.root)
	fb.TreeView.SetGraphics(false)
	fb.TreeView.Box.SetBackgroundColor(fbbgColor)

	// Add the current directory to the root node.
	fb.AddDirEntry(fb.root, fb.rootDir)
	// If a directory was selected, open it.
	fb.TreeView.SetSelectedFunc(fb.OpenDirectory)

	subFlex.AddItem(fb, filebrowserW, 1, false)
	fb.parentFlex = subFlex
	curFilebrowser = fb
	// Default to closed
	fb.ToggleDisplay()

	return fb
}

func (fb *filebrowser) IsVisible() bool {
	return fb.TreeView.Box.GetVisible()
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
	if reference == nil {
		return // Selecting the root node does nothing.
	}
	children := node.GetChildren()
	if len(children) == 0 {
		// Load and show files in this directory.
		path := reference.(string)
		if err := fb.AddDirEntry(node, path); err == nil {
			//was a file
		}
		node.SetExpanded(true)
	} else {
		// Collapse if visible, expand if collapsed.
		node.ClearChildren()
		node.SetChildren(nil)
		node.SetExpanded(false)
	}
}

func (fb *filebrowser) HandleInput(tevent *tcell.EventKey) *tcell.EventKey {
	if tevent.Key() == tcell.KeyCtrlD {
		fb.ToggleDisplay()
		return nil
	}
	return tevent
}

func (fb *filebrowser) ToggleDisplay() {
	if fb.TreeView.Box.GetVisible() {
		fb.parentFlex.ResizeItem(fb, -1, 0)
		fb.TreeView.Box.SetVisible(false)
		app.CurApp.SetFocus(fb)
	} else {
		fb.parentFlex.ResizeItem(fb, filebrowserW, 1)
		fb.TreeView.Box.SetVisible(true)
		app.CurApp.SetFocus(nil)
	}
}
