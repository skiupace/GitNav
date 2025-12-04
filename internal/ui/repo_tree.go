package ui

import (
	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"

	"github.com/skiupace/gitnav/internal/git"
)

func RepoTree(rootNode *git.Node) *tview.TreeView {
	root := tview.NewTreeNode("")
	// SetSelectable(false)

	addChildren(root, rootNode)

	tree := tview.NewTreeView().
		SetRoot(root).
		SetCurrentNode(root)

	title := " " + rootNode.Name + " "
	tree.Box.SetBorder(true).
		SetBorderColor(tcell.ColorBlue).
		SetTitleAlign(tview.AlignLeft).
		SetTitle(title)

	return tree
}

func addChildren(tnode *tview.TreeNode, gnode *git.Node) {
	for _, child := range gnode.Children {
		node := tview.NewTreeNode(child.Name).
			SetReference(child.Path).
			SetColor(tcell.ColorWhite)

		if child.IsDir {
			node.SetColor(tcell.ColorYellow)
			// node.SetExpanded(false)
		}

		tnode.AddChild(node)

		// Recursively add child folders/files
		if child.IsDir {
			addChildren(node, child)
		}
	}
}
