package ui

import (
	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

func FilePreview() tview.Primitive {
	preview := tview.NewBox().
		SetBorder(true).
		SetTitle(" Preview: <select-file> ").
		SetTitleAlign(tview.AlignLeft).
		SetBackgroundColor(tcell.ColorDefault)

	flex := tview.NewFlex().
		SetDirection(tview.FlexRow).
		AddItem(preview, 0, 1, false)

	return flex
}
