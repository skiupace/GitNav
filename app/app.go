package app

import (
	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"

	"github.com/skiupace/gitnav/internal/ui"
)

var (
	App    *Application
	Styles *Theme
)

type Application struct {
	app *tview.Application
}

type Theme struct {
	tview.Theme
	SidebarTitleBorderColor string
}

func init() {
	App = &Application{
		app: tview.NewApplication(),
	}

	Styles = &Theme{
		Theme: tview.Theme{
			PrimitiveBackgroundColor:    tcell.ColorDefault,
			ContrastBackgroundColor:     tcell.ColorBlue,
			MoreContrastBackgroundColor: tcell.ColorGreen,
			BorderColor:                 tcell.ColorWhite,
			TitleColor:                  tcell.ColorWhite,
			GraphicsColor:               tcell.ColorGray,
			PrimaryTextColor:            tcell.ColorDefault.TrueColor(),
			SecondaryTextColor:          tcell.ColorYellow,
			TertiaryTextColor:           tcell.ColorGreen,
			InverseTextColor:            tcell.ColorWhite,
			ContrastSecondaryTextColor:  tcell.ColorBlack,
		},
		SidebarTitleBorderColor: "#666A7E",
	}

	tview.Styles = Styles.Theme
}

func (a *Application) Run(repoPath string) error {
	layout := ui.BaseLayout(repoPath)

	a.app.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		if event.Rune() == 'q' || event.Key() == tcell.KeyCtrlC {
			a.app.Stop()
			return nil
		}
		return event
	})

	return a.app.SetRoot(layout, true).Run()
}
