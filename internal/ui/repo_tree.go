package ui

import (
	"strings"

	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"

	"github.com/skiupace/gitnav/commands"
	"github.com/skiupace/gitnav/internal/git"
)

// RepoTree builds the file tree view from the scanned repo data.
func RepoTree(rootNode *git.Node) *tview.TreeView {
	root := newRootNode()

	tree := tview.NewTreeView().
		SetRoot(root).
		SetCurrentNode(root)

	// Note: SetSelectedFunc is set by BaseLayout in layout.go
	// to handle both directory toggling and opening files in editor.

	tree.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		cmd := TreeKeyMap.Resolve(event)

		switch cmd {
		case commands.MoveDown:
			return tcell.NewEventKey(tcell.KeyDown, 0, tcell.ModNone)
		case commands.MoveUp:
			return tcell.NewEventKey(tcell.KeyUp, 0, tcell.ModNone)
		case commands.MoveLeft:
			syncFolderIcon(tree.GetCurrentNode(), false)
			return tcell.NewEventKey(tcell.KeyLeft, 0, tcell.ModNone)
		case commands.MoveRight:
			syncFolderIcon(tree.GetCurrentNode(), true)
			return tcell.NewEventKey(tcell.KeyRight, 0, tcell.ModNone)
		case commands.Select:
			toggleExpansion(tree.GetCurrentNode())
			return nil
		}

		return event
	})

	addChildren(root, rootNode)

	tree.Box.SetBorder(true).
		SetBorderColor(tcell.ColorBlue).
		SetTitleAlign(tview.AlignLeft).
		SetTitle(" " + rootNode.Name + " ")

	return tree
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
	return tview.NewTreeNode("").
		SetColor(tcell.ColorBlue).
		SetSelectable(false)
}

func newFileNode(gnode *git.Node) *tview.TreeNode {
	icon, color := GetIcon(gnode.Name, gnode.IsDir)

	// SetColor also sets the selected background to the icon color, which
	// reads badly; give the selected state a controlled, theme-friendly
	// style instead (node color on the palette's gray).
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
