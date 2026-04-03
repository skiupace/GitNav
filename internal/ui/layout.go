package ui

import (
	"os"

	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"

	"github.com/skiupace/gitnav/internal/git"
)

func BaseLayout(repoPath string, app *tview.Application) tview.Primitive {
	// LEFT TREE
	repo, err := git.NewRepo(repoPath)
	if err != nil {
		panic("Not a git repository")
	}

	rootNode, err := repo.ScanRepo()
	if err != nil {
		panic(err)
	}

	tree := RepoTree(rootNode)

	// CENTER PREVIEW
	preview := NewPreviewPanel()

	// Track focusable panels for Tab cycling
	panels := []tview.Primitive{tree, preview.TextView}
	focusIndex := 0

	// Visual focus indicator helper
	updateFocusBorders := func(idx int) {
		if idx == 0 {
			tree.Box.SetBorderColor(tcell.ColorBlue)
			preview.TextView.SetBorderColor(tcell.ColorGray)
		} else {
			tree.Box.SetBorderColor(tcell.ColorGray)
			preview.TextView.SetBorderColor(tcell.ColorBlue)
		}
	}

	// Wire tree node changes to update the preview
	tree.SetChangedFunc(func(node *tview.TreeNode) {
		ref := node.GetReference()
		if ref == nil {
			return
		}
		filePath, ok := ref.(string)
		if !ok {
			return
		}

		// Only preview files, not directories
		info, err := os.Stat(filePath)
		if err != nil || info.IsDir() {
			preview.ClearPreview()
			return
		}

		preview.UpdatePreview(filePath)
	})

	// SEARCH BOX (bottom-left)
	search := tview.NewBox().
		SetBorder(true).
		SetTitleAlign(tview.AlignLeft).
		SetTitle(" Search ").
		SetBackgroundColor(tcell.ColorDefault)

	// STATS BOX (bottom-right)
	stats := tview.NewBox().
		SetBorder(true).
		SetTitle(" Repo Statistics ").
		SetTitleAlign(tview.AlignLeft).
		SetBackgroundColor(tcell.ColorDefault)

	// MAIN ROW (tree + preview)
	top := tview.NewFlex().
		AddItem(tree, 0, 1, true).
		AddItem(preview.TextView, 0, 3, false)

	bottom := tview.NewFlex().
		AddItem(search, 0, 1, false).
		AddItem(stats, 0, 3, false)

	// ROOT LAYOUT
	root := tview.NewFlex().SetDirection(tview.FlexRow)
	root.AddItem(top, 0, 3, true)
	root.AddItem(bottom, 3, 0, false)

	// Tab / Shift+Tab focus switching at the root level
	root.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		if event.Key() == tcell.KeyTab || event.Key() == tcell.KeyRight {
			focusIndex = (focusIndex + 1) % len(panels)
			updateFocusBorders(focusIndex)
			app.SetFocus(panels[focusIndex])
			return nil
		}
		if event.Key() == tcell.KeyBacktab || event.Key() == tcell.KeyLeft {
			focusIndex = (focusIndex - 1 + len(panels)) % len(panels)
			updateFocusBorders(focusIndex)
			app.SetFocus(panels[focusIndex])
			return nil
		}
		return event
	})

	// Set initial focus borders
	updateFocusBorders(0)

	return root
}
