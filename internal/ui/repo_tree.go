package ui

import (
	"strings"

	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
	"github.com/skiupace/gitnav/commands"
	"github.com/skiupace/gitnav/internal/git"
)

func RepoTree(rootNode *git.Node) *tview.TreeView {
	root := tview.NewTreeNode("").
		SetColor(tcell.ColorBlue).
		SetSelectable(false)

	tree := tview.NewTreeView().
		SetRoot(root).
		SetCurrentNode(root)

	tree.SetSelectedFunc(func(node *tview.TreeNode) {
		toggleExpansion(node)
	})

	tree.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		cmd := TreeKeyMap.Resolve(event)

		switch cmd {
		case commands.MoveDown:
			return tcell.NewEventKey(tcell.KeyDown, 0, tcell.ModNone)
		case commands.MoveUp:
			return tcell.NewEventKey(tcell.KeyUp, 0, tcell.ModNone)
		case commands.MoveLeft:
			return tcell.NewEventKey(tcell.KeyLeft, 0, tcell.ModNone)
		case commands.MoveRight:
			return tcell.NewEventKey(tcell.KeyRight, 0, tcell.ModNone)
		case commands.Select:
			toggleExpansion(tree.GetCurrentNode())
			return nil
		}

		return event
	})

	addChildren(root, rootNode)

	title := " " + rootNode.Name + " "
	tree.Box.SetBorder(true).
		SetBorderColor(tcell.ColorBlue).
		SetTitleAlign(tview.AlignLeft).
		SetTitle(title)

	return tree
}

func addChildren(tnode *tview.TreeNode, gnode *git.Node) {
	for _, child := range gnode.Children {
		icon, color := GetIcon(child.Name, child.IsDir)

		// If dir, use GetFolderIcon with default collapsed state (expanded=false)
		// collapsed -> Open icon (as per inverted logic)
		if child.IsDir {
			icon = GetFolderIcon(false)
		}

		node := tview.NewTreeNode(icon + " " + child.Name).
			SetReference(child.Path).
			SetColor(color)

		if child.IsDir {
			// node.SetColor(tcell.ColorYellow) // handled by GetIcon
			node.SetExpanded(false)
		}

		tnode.AddChild(node)

		// Recursively add child folders/files
		if child.IsDir {
			addChildren(node, child)
		}
	}
}

func toggleExpansion(node *tview.TreeNode) {
	if node == nil {
		return
	}
	// If it has children (i.e., is a directory), toggle expansion
	if len(node.GetChildren()) > 0 {
		expanded := !node.IsExpanded()
		node.SetExpanded(expanded)

		// Update icon based on new state
		// We need to preserve the name while changing the icon
		text := node.GetText()
		parts := strings.SplitN(text, " ", 2)
		if len(parts) == 2 {
			name := parts[1]
			newIcon := GetFolderIcon(expanded)
			node.SetText(newIcon + " " + name)
		}
	}
}
