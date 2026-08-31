package ui

import (
	"os"
	"strings"

	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"

	"github.com/skiupace/gitnav/internal/git"
)

// SearchPanel incrementally filters the repo tree as you type. tview nodes
// have no visibility control, so filtering rebuilds the tree from the
// pristine git.Node data; the previous expansion state is snapshotted by
// path and restored when the query is cleared. Enter jumps to the first
// matching file, Escape clears the filter.
type SearchPanel struct {
	Field     *tview.InputField
	FocusTree func() // set by BaseLayout to move focus back to the tree

	tree     *tview.TreeView
	rootData *git.Node
	preview  *PreviewPanel
	snapshot map[string]bool // path -> expanded, captured before first filter
}

func NewSearchPanel(tree *tview.TreeView, rootData *git.Node, preview *PreviewPanel) *SearchPanel {
	sp := &SearchPanel{
		tree:     tree,
		rootData: rootData,
		preview:  preview,
	}

	sp.Field = tview.NewInputField().
		SetLabel(" / ").
		SetPlaceholder(" filter files ").
		SetLabelColor(tcell.ColorGray).
		SetFieldBackgroundColor(tcell.ColorDefault).
		SetFieldTextColor(tview.Styles.SecondaryTextColor).
		SetPlaceholderStyle(tcell.StyleDefault.Foreground(tcell.ColorGray))

	sp.Field.SetBorder(true).
		SetTitleAlign(tview.AlignLeft).
		SetTitle(" Search ").
		SetBackgroundColor(tcell.ColorDefault)

	sp.Field.SetChangedFunc(func(text string) { sp.apply(text) })

	sp.Field.SetDoneFunc(func(key tcell.Key) {
		switch key {
		case tcell.KeyEnter:
			if sp.jump() && sp.FocusTree != nil {
				sp.FocusTree()
			}
		case tcell.KeyEscape:
			sp.Field.SetText("")
			sp.apply("")
			if sp.FocusTree != nil {
				sp.FocusTree()
			}
		}
	})

	return sp
}

// apply swaps in a tree containing only files matching query (plus their
// ancestor folders, expanded).
func (sp *SearchPanel) apply(query string) {
	if query == "" {
		sp.restore()
		return
	}

	// Remember the current expansion state before the first filter pass.
	if sp.snapshot == nil {
		sp.snapshot = map[string]bool{}
		collectExpanded(sp.tree.GetRoot(), sp.snapshot)
	}

	root := newRootNode()
	for _, child := range sp.rootData.Children {
		if sub := buildFiltered(child, strings.ToLower(query)); sub != nil {
			root.AddChild(sub)
		}
	}

	sp.tree.SetRoot(root)
	if n := firstFile(root); n != nil {
		sp.tree.SetCurrentNode(n)
	} else {
		sp.tree.SetCurrentNode(root)
	}
}

// restore rebuilds the full tree and puts back the pre-filter expansion
// state.
func (sp *SearchPanel) restore() {
	root := newRootNode()
	addChildren(root, sp.rootData)

	if sp.snapshot != nil {
		byPath := map[string]*tview.TreeNode{}
		collectNodes(root, byPath)
		for path, expanded := range sp.snapshot {
			if n, ok := byPath[path]; ok {
				n.SetExpanded(expanded)
			}
		}
		sp.snapshot = nil
	}

	sp.tree.SetRoot(root)
	sp.tree.SetCurrentNode(root)
}

// jump selects the first matching file and previews it. Reports success.
func (sp *SearchPanel) jump() bool {
	n := firstFile(sp.tree.GetRoot())
	if n == nil {
		return false
	}
	sp.tree.SetCurrentNode(n)

	if path, ok := n.GetReference().(string); ok {
		if info, err := os.Stat(path); err == nil && !info.IsDir() {
			sp.preview.UpdatePreview(path)
		}
	}
	return true
}

// buildFiltered returns a new subtree for gnode if its name or that of any
// descendant matches query, nil otherwise. Included folders are expanded so
// matches stay visible.
func buildFiltered(gnode *git.Node, query string) *tview.TreeNode {
	if !gnode.IsDir {
		if !strings.Contains(strings.ToLower(gnode.Name), query) {
			return nil
		}
		return newFileNode(gnode)
	}

	node := newFileNode(gnode)
	for _, child := range gnode.Children {
		if sub := buildFiltered(child, query); sub != nil {
			node.AddChild(sub)
		}
	}
	if len(node.GetChildren()) == 0 {
		return nil
	}
	node.SetExpanded(true)
	return node
}

// firstFile returns the first file node (depth-first) in the tree, or nil.
// Filtered folders always have children, so childless nodes are files.
func firstFile(node *tview.TreeNode) *tview.TreeNode {
	for _, child := range node.GetChildren() {
		if len(child.GetChildren()) == 0 {
			return child
		}
		if n := firstFile(child); n != nil {
			return n
		}
	}
	return nil
}

func collectExpanded(node *tview.TreeNode, into map[string]bool) {
	if path, ok := node.GetReference().(string); ok {
		into[path] = node.IsExpanded()
	}
	for _, child := range node.GetChildren() {
		collectExpanded(child, into)
	}
}

func collectNodes(node *tview.TreeNode, into map[string]*tview.TreeNode) {
	if path, ok := node.GetReference().(string); ok {
		into[path] = node
	}
	for _, child := range node.GetChildren() {
		collectNodes(child, into)
	}
}
