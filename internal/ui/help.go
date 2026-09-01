package ui

import (
	"fmt"
	"strings"

	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"

	"github.com/skiupace/gitnav/commands"
	"github.com/skiupace/gitnav/keymap"
)

func HelpPanel(onClose func()) *tview.TextView {
	tv := tview.NewTextView().
		SetDynamicColors(true).
		SetWrap(false)

	tv.SetBorder(true).
		SetTitleAlign(tview.AlignLeft).
		SetTitle(" Keybindings — j/k scroll, any other key closes ").
		SetBackgroundColor(tcell.ColorDefault)

	var b strings.Builder
	section := func(name, binds string) {
		b.WriteString("[::b]" + name + "[-]\n")
		b.WriteString(binds)
		b.WriteString("\n")
	}

	section("Global", formatBinds(GlobalKeyMap)+
		fmt.Sprintf("  %-26s %s\n", "Tab / Shift+Tab", "Cycle focus")+
		fmt.Sprintf("  %-26s %s\n", "Ctrl+C", "Quit"))
	section("Tree (normal mode)", formatBinds(TreeKeyMap))
	section("Preview", formatBinds(PreviewKeyMap))
	section("Search (insert mode)",
		fmt.Sprintf("  %-26s %s\n", "type", "Filter the tree")+
			fmt.Sprintf("  %-26s %s\n", "Enter", "Jump to first match")+
			fmt.Sprintf("  %-26s %s\n", "Esc", "Clear filter, back to NORMAL"))

	tv.SetText(b.String())

	tv.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		switch PreviewKeyMap.Resolve(event) {
		case commands.ScrollDown, commands.ScrollUp, commands.ScrollTop, commands.ScrollBottom:
			return scrollCapture(tv)(event)
		}
		onClose()
		return nil
	})

	return tv
}

func formatBinds(m keymap.Map) string {
	var order []string
	keysFor := map[string][]string{}
	for _, bind := range m {
		desc := bind.Description
		if _, seen := keysFor[desc]; !seen {
			order = append(order, desc)
		}
		keysFor[desc] = append(keysFor[desc], bind.Key.String())
	}

	var b strings.Builder
	for _, desc := range order {
		b.WriteString(fmt.Sprintf("  %-26s %s\n", strings.Join(keysFor[desc], " / "), desc))
	}
	return b.String()
}
