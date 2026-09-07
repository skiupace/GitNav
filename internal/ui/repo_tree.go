package ui

import (
	"strings"

	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"

	"github.com/skiupace/gitnav/commands"
	"github.com/skiupace/gitnav/internal/git"
)

type TreePanel struct {
	*tview.TreeView
	OnCopyPath func()
}

// RepoTree builds the file tree view from the scanned repo data.
func RepoTree(rootNode *git.Node) *TreePanel {
	root := newRootNode()

	treeView := tview.NewTreeView().
		SetRoot(root).
		SetCurrentNode(root)

	tp := &TreePanel{TreeView: treeView}

	// Note: SetSelectedFunc is set by BaseLayout in layout.go
	// to handle both directory toggling and opening files in editor.

	treeView.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		cmd := TreeKeyMap.Resolve(event)

		switch cmd {
		case commands.MoveDown:
			return tcell.NewEventKey(tcell.KeyDown, 0, tcell.ModNone)
		case commands.MoveUp:
			return tcell.NewEventKey(tcell.KeyUp, 0, tcell.ModNone)
		case commands.MoveLeft:
			node := treeView.GetCurrentNode()
			if node.IsExpanded() {
				syncFolderIcon(node, false)
				node.SetExpanded(false)
			}
			return nil
		case commands.MoveRight:
			node := treeView.GetCurrentNode()
			if !node.IsExpanded() {
				syncFolderIcon(node, true)
				node.SetExpanded(true)
			}
			return nil
		case commands.ScrollTop:
			first, _ := visibleTreeNodes(root)
			if first != nil {
				treeView.SetCurrentNode(first)
			}
			return nil
		case commands.ScrollBottom:
			_, last := visibleTreeNodes(root)
			if last != nil {
				treeView.SetCurrentNode(last)
			}
			return nil
		case commands.CopyPath:
			if tp.OnCopyPath != nil {
				tp.OnCopyPath()
			}
			return nil
		case commands.Select:
			toggleExpansion(treeView.GetCurrentNode())
			return nil
		}

		return event
	})

	addChildren(root, rootNode)

	treeView.Box.SetBorder(true).
		SetBorderColor(tcell.ColorBlue).
		SetTitleAlign(tview.AlignLeft).
		SetTitle(" " + rootNode.Name + " ")

	return tp
}

func visibleTreeNodes(root *tview.TreeNode) (first, last *tview.TreeNode) {
	var traverse func(n *tview.TreeNode)
	traverse = func(n *tview.TreeNode) {
		for _, child := range n.GetChildren() {
			if first == nil {
				first = child
			}
			last = child
			if child.IsExpanded() {
				traverse(child)
			}
		}
	}
	traverse(root)
	return
}

func addChildren(tnode *tview.TreeNode, gnode *git.Node) {
	for _, child := range gnode.Children {
		node := newFileNode(child)
		tnode.AddChild(node)

		if child.IsDir {
			addChildren(node, child)
		}
	}
}

func newRootNode() *tview.TreeNode {
	return tview.NewTreeNode("\ue21c").
		SetColor(tcell.ColorBlue).
		SetSelectable(false)
}

func newFileNode(gnode *git.Node) *tview.TreeNode {
	icon, color := GetIcon(gnode.Name, gnode.IsDir)

	node := tview.NewTreeNode(icon + " " + gnode.Name).
		SetReference(gnode.Path).
		SetColor(color).
		SetSelectedTextStyle(tcell.StyleDefault.
			Foreground(color).
			Background(tcell.ColorGray))

	if gnode.IsDir {
		node.SetExpanded(false)
	}

	return node
}

func toggleExpansion(node *tview.TreeNode) {
	if node == nil || len(node.GetChildren()) == 0 {
		return
	}

	expanded := !node.IsExpanded()
	node.SetExpanded(expanded)
	node.SetText(GetFolderIcon(expanded) + " " + nodeName(node.GetText()))
}

// syncFolderIcon pre-updates the folder icon for h/l navigation, since tview
// expands/collapses nodes itself without notifying us. right=true means the
// node is about to expand, right=false means it is about to collapse.
func syncFolderIcon(node *tview.TreeNode, right bool) {
	if node == nil || len(node.GetChildren()) == 0 {
		return
	}
	if right != node.IsExpanded() {
		node.SetText(GetFolderIcon(right) + " " + nodeName(node.GetText()))
	}
}

// nodeName returns the text after the icon ("icon name" -> "name").
func nodeName(text string) string {
	if _, name, ok := strings.Cut(text, " "); ok {
		return name
	}
	return text
}
