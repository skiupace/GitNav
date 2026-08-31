package ui

import (
	"fmt"
	"time"

	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"

	"github.com/skiupace/gitnav/internal/git"
)

// StatsPanel shows branch, HEAD and worktree status for the open repo.
type StatsPanel struct {
	View *tview.TextView
	repo *git.Repo
}

func NewStatsPanel(repo *git.Repo) *StatsPanel {
	v := tview.NewTextView().
		SetDynamicColors(true).
		SetWrap(false)

	v.SetBorder(true).
		SetTitleAlign(tview.AlignLeft).
		SetTitle(" Repo Statistics ").
		SetBackgroundColor(tcell.ColorDefault)

	sp := &StatsPanel{View: v, repo: repo}
	sp.Refresh()
	return sp
}

// Refresh re-reads repository metadata from disk.
func (sp *StatsPanel) Refresh() {
	in := sp.repo.Info()

	text := ""
	if in.SHA != "" {
		text = fmt.Sprintf(" [::b]%s[-] · %s · %s · %s, %s",
			in.Branch, in.SHA, trunc(in.Subject, 42), in.Author, relTime(in.When))
	} else {
		text = in.Branch + text
	}

	sp.View.SetText(text)
	sp.View.SetTitle(statusTitle(in.Dirty))
}

func statusTitle(dirty int) string {
	status := "[green]clean"
	if dirty > 0 {
		status = fmt.Sprintf("[yellow]%d changed file(s)", dirty)
	}
	return status
}

// relTime renders a coarse relative timestamp.
func relTime(t time.Time) string {
	d := time.Since(t)
	switch {
	case d < time.Minute:
		return "just now"
	case d < time.Hour:
		return fmt.Sprintf("%dm ago", int(d.Minutes()))
	case d < 24*time.Hour:
		return fmt.Sprintf("%dh ago", int(d.Hours()))
	case d < 30*24*time.Hour:
		return fmt.Sprintf("%dd ago", int(d.Hours()/24))
	case d < 365*24*time.Hour:
		return fmt.Sprintf("%dmo ago", int(d.Hours()/(24*30)))
	default:
		return fmt.Sprintf("%dy ago", int(d.Hours()/(24*365)))
	}
}

func trunc(s string, n int) string {
	r := []rune(s)
	if len(r) > n {
		return string(r[:n-1]) + "…"
	}
	return s
}
