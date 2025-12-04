package app

import (
	"github.com/rivo/tview"
)

type AppConfig struct {
	DefaultPageSize int
	DisableSidebar  bool
	SidebarOverlay  bool
	DefaultEditor   string
}

type Config struct {
	ConfigFile string
	AppConfig  AppConfig `toml:"application"`
}

type UI struct {
	Tree    *tview.TreeView
	Preview *tview.TextView
	Search  *tview.InputField
	Stats   *tview.TextView
}
