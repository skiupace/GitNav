package ui

import (
	"os"
	"os/exec"

	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"

	"github.com/skiupace/gitnav/commands"
	"github.com/skiupace/gitnav/internal/git"
)

func BaseLayout(repoPath string, a *tview.Application) tview.Primitive {
	repo, err := git.NewRepo(repoPath)
	if err != nil {
		panic(err)
	}

	rootNode, err := repo.ScanRepo()
	if err != nil {
		panic(err)
	}

	tree := RepoTree(rootNode)
	preview := NewPreviewPanel()
	stats := NewStatsPanel(repo)
	search := NewSearchPanel(tree, rootNode, preview)

	// Open file in the user's editor, suspending tview until it exits.
	openInEditor := func(filePath string) {
		a.Suspend(func() {
			cmd := exec.Command(pickEditor(), filePath)
			cmd.Stdin, cmd.Stdout, cmd.Stderr = os.Stdin, os.Stdout, os.Stderr
			_ = cmd.Run()
		})
		stats.Refresh()
	}

	// Wire tree node changes to update the preview
	tree.SetChangedFunc(func(node *tview.TreeNode) {
		filePath, ok := node.GetReference().(string)
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

	// Handle Enter on tree: open file in editor, toggle dirs
	tree.SetSelectedFunc(func(node *tview.TreeNode) {
		filePath, ok := node.GetReference().(string)
		if !ok {
			return
		}

		info, err := os.Stat(filePath)
		if err != nil {
			return
		}

		if info.IsDir() {
			toggleExpansion(node)
		} else {
			openInEditor(filePath)
		}
	})

	// MAIN ROW (tree + preview)
	top := tview.NewFlex().
		AddItem(tree, 0, 1, true).
		AddItem(preview.TextView, 0, 3, false)

	// BOTTOM ROW (search + stats). Fixed height: a 1:3 proportional split
	// inside a 3-row flex gave the search box zero rows.
	bottom := tview.NewFlex().
		AddItem(search.Field, 0, 1, false).
		AddItem(stats.View, 0, 3, false)

	// ROOT LAYOUT
	root := tview.NewFlex().SetDirection(tview.FlexRow)
	root.AddItem(top, 0, 1, true)
	root.AddItem(bottom, 3, 0, false)

	// Focusable panels (tree -> preview -> search); stats is passive.
	// Conceptually a 2x2 grid: tree, preview / search, (stats).
	panels := []tview.Primitive{tree, preview.TextView, search.Field}
	boxes := []*tview.Box{tree.Box, preview.TextView.Box, search.Field.Box}
	focusIndex := 0

	updateFocusBorders := func(idx int) {
		for i, b := range boxes {
			if i == idx {
				b.SetBorderColor(tcell.ColorBlue)
			} else {
				b.SetBorderColor(tcell.ColorGray)
			}
		}
	}

	// Vim-style mode indicator: NORMAL in the panes, INSERT in search.
	setMode := func(insert bool) {
		if insert {
			search.Field.SetTitle(" Search · INSERT ")
		} else {
			search.Field.SetTitle(" Search · NORMAL ")
		}
	}

	setFocus := func(idx int) {
		focusIndex = idx
		a.SetFocus(panels[idx])
		updateFocusBorders(idx)
		setMode(idx == 2)
	}

	// Search hands focus back to the tree (Enter/Esc).
	search.FocusTree = func() { setFocus(0) }

	// Vim pane navigation: uppercase HJKL and Ctrl+arrows move focus
	// between panes; lowercase keys stay free for in-pane navigation.
	moveFocus := func(cmd commands.Command) {
		switch cmd {
		case commands.FocusSidebar: // left
			if focusIndex == 1 {
				setFocus(0)
			}
		case commands.FocusPreview: // right
			if focusIndex != 1 {
				setFocus(1)
			}
		case commands.Search: // down / insert
			if focusIndex != 2 {
				setFocus(2)
			}
		case commands.FocusTree: // up
			if focusIndex == 2 {
				setFocus(0)
			}
		}
	}

	// Keybinding cheat sheet overlay.
	helpVisible := false
	toggleHelp := func() {
		helpVisible = !helpVisible
		if helpVisible {
			// a.SetRoot(HelpPanel(toggleHelp), true)
		} else {
			a.SetRoot(root, true)
			setFocus(focusIndex)
		}
	}

	// Global keys: Tab/Shift+Tab cycle panels; '/' 'i' 'H' 'J' 'K' 'L',
	// Ctrl+arrows and 'q'/'?' act outside the search field, where letters
	// are reserved for typing.
	root.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		switch event.Key() {
		case tcell.KeyTab:
			setFocus((focusIndex + 1) % len(panels))
			return nil
		case tcell.KeyBacktab:
			setFocus((focusIndex - 1 + len(panels)) % len(panels))
			return nil
		}

		if a.GetFocus() == search.Field {
			return event
		}

		switch cmd := GlobalKeyMap.Resolve(event); cmd {
		case commands.Quit:
			a.Stop()
			return nil
		case commands.Help:
			toggleHelp()
			return nil
		case commands.FocusSidebar, commands.FocusPreview, commands.FocusTree, commands.Search:
			moveFocus(cmd)
			return nil
		}

		return event
	})

	// Initial focus + mode title
	setFocus(0)

	return root
}

// pickEditor returns the editor to open files with: $VISUAL, $EDITOR, or vi.
func pickEditor() string {
	if e := os.Getenv("VISUAL"); e != "" {
		return e
	}
	if e := os.Getenv("EDITOR"); e != "" {
		return e
	}
	return "vi"
}
