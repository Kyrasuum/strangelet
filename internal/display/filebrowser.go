package display

import (
	"io/ioutil"
	"path/filepath"

	"github.com/Kyrasuum/cview"
	"github.com/gdamore/tcell/v2"
)

var (
	dirColor tcell.Color
	filColor tcell.Color
	bgColor  tcell.Color

	filebrowserW int = 20
)

type Filebrowser struct {
	*cview.TreeView
	root *cview.TreeNode

	rootDir string
}

func (fb *Filebrowser) InitFilebrowser(subFlex *cview.Flex) {
	if dirColor == 0 {
		dirColor = tcell.NewRGBColor(220, 100, 100)
	}
	if filColor == 0 {
		filColor = tcell.NewRGBColor(220, 220, 220)
	}
	if bgColor == 0 {
		bgColor = tcell.NewRGBColor(30, 30, 30)
	}

	fb.rootDir = "."
	fb.root = cview.NewTreeNode(fb.rootDir)
	fb.root.SetColor(dirColor)
	fb.TreeView = cview.NewTreeView()
	fb.TreeView.SetRoot(fb.root)
	fb.TreeView.SetCurrentNode(fb.root)
	fb.TreeView.SetGraphics(false)
	fb.TreeView.Box.SetBackgroundColor(bgColor)

	// Add the current directory to the root node.
	fb.AddDirEntry(fb.root, fb.rootDir)
	// If a directory was selected, open it.
	fb.TreeView.SetSelectedFunc(fb.OpenDirectory)

	subFlex.AddItem(fb, filebrowserW, 1, false)
}

func (fb *Filebrowser) IsVisible() bool {
	return fb.TreeView.Box.GetVisible()
}

// A helper function which adds the files and directories of the given path
// to the given target node.
func (fb *Filebrowser) AddDirEntry(target *cview.TreeNode, path string) (err error) {
	files, err := ioutil.ReadDir(path)
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
func (fb *Filebrowser) OpenDirectory(node *cview.TreeNode) {
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

func (fb *Filebrowser) HandleInput(event *tcell.EventKey) *tcell.EventKey {
	return event
}

func (fb *Filebrowser) ToggleDisplay(subFlex *cview.Flex) {
	if fb.TreeView.Box.GetVisible() {
		subFlex.ResizeItem(fb, -1, 0)
		fb.TreeView.Box.SetVisible(false)
	} else {
		subFlex.ResizeItem(fb, filebrowserW, 1)
		fb.TreeView.Box.SetVisible(true)
	}
}
