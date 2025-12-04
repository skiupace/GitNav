package ui

import (
	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"

	"github.com/skiupace/gitnav/internal/git"
)

func BaseLayout(repoPath string) tview.Primitive {
	// LEFT TREE
	repo, err := git.NewRepo(repoPath)
	if err != nil {
		panic("Not a git repository")
	}

	rootNode, err := repo.ScanRepo()
	if err != nil {
		panic(err)
	}

	treeFlex := RepoTree(rootNode)

	// CENTER PREVIEW
	previewFlex := FilePreview()

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
		AddItem(treeFlex, 0, 1, false).
		AddItem(previewFlex, 0, 3, false)

	bottom := tview.NewFlex().
		AddItem(search, 0, 1, false).
		AddItem(stats, 0, 3, false)

	// ROOT LAYOUT
	root := tview.NewFlex().SetDirection(tview.FlexRow)
	root.AddItem(top, 0, 3, false)
	root.AddItem(bottom, 3, 0, false)

	return root
}
