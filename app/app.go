package app

import (
	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"

	"github.com/skiupace/gitnav/internal/ui"
)

var (
	App    *Application
	Styles *tview.Theme
)

type Application struct {
	app *tview.Application
}

func init() {
	App = &Application{
		app: tview.NewApplication(),
	}

	// Colors follow the terminal's own theme: backgrounds default, accents
	// use the 16-color ANSI palette (tcell named colors resolve to palette
	// indices, so they adapt to the user's color scheme).
	Styles = &tview.Theme{
		PrimitiveBackgroundColor:    tcell.ColorDefault,
		ContrastBackgroundColor:     tcell.ColorDefault,
		MoreContrastBackgroundColor: tcell.ColorDefault,
		BorderColor:                 tcell.ColorGray,
		TitleColor:                  tcell.ColorDefault,
		GraphicsColor:               tcell.ColorGray,
		PrimaryTextColor:            tcell.ColorDefault,
		SecondaryTextColor:          tcell.ColorYellow,
		TertiaryTextColor:           tcell.ColorGreen,
		InverseTextColor:            tcell.ColorDefault,
		ContrastSecondaryTextColor:  tcell.ColorGray,
	}

	tview.Styles = *Styles
}

func (a *Application) Run(repoPath string) error {
	layout := ui.BaseLayout(repoPath, a.app)

	a.app.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		// 'q' is handled by the layout (suppressed while the search
		// input field has focus); Ctrl+C always quits.
		if event.Key() == tcell.KeyCtrlC {
			a.app.Stop()
			return nil
		}
		return event
	})

	return a.app.SetRoot(layout, true).Run()
}