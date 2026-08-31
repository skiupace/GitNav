package ui

import (
	"github.com/gdamore/tcell/v2"
	"github.com/skiupace/gitnav/commands"
	"github.com/skiupace/gitnav/keymap"
)

var TreeKeyMap = keymap.Map{
	{
		Key:         keymap.Key{Char: 'j'},
		Cmd:         commands.MoveDown,
		Description: "Move down",
	},
	{
		Key:         keymap.Key{Code: tcell.KeyDown},
		Cmd:         commands.MoveDown,
		Description: "Move down",
	},
	{
		Key:         keymap.Key{Char: 'k'},
		Cmd:         commands.MoveUp,
		Description: "Move up",
	},
	{
		Key:         keymap.Key{Code: tcell.KeyUp},
		Cmd:         commands.MoveUp,
		Description: "Move up",
	},
	{
		Key:         keymap.Key{Char: 'h'},
		Cmd:         commands.MoveLeft,
		Description: "Collapse / Go to parent",
	},
	{
		Key:         keymap.Key{Code: tcell.KeyLeft},
		Cmd:         commands.MoveLeft,
		Description: "Collapse / Go to parent",
	},
	{
		Key:         keymap.Key{Char: 'l'},
		Cmd:         commands.MoveRight,
		Description: "Expand / Go to child",
	},
	{
		Key:         keymap.Key{Code: tcell.KeyRight},
		Cmd:         commands.MoveRight,
		Description: "Expand / Go to child",
	},
	{
		Key:         keymap.Key{Code: tcell.KeyEnter},
		Cmd:         commands.Select,
		Description: "Open file in editor / toggle folder",
	},
}

// GlobalKeyMap binds pane navigation and app-level actions. Vim-style
// uppercase HJKL and Ctrl+arrows move focus between panes; lowercase keys
// stay free for in-pane navigation.
var GlobalKeyMap = keymap.Map{
	{
		Key:         keymap.Key{Char: 'H'},
		Cmd:         commands.FocusSidebar,
		Description: "Focus tree pane",
	},
	{
		Key:         keymap.Key{Char: 'L'},
		Cmd:         commands.FocusPreview,
		Description: "Focus preview pane",
	},
	{
		Key:         keymap.Key{Char: 'J'},
		Cmd:         commands.Search,
		Description: "Focus search (insert mode)",
	},
	{
		Key:         keymap.Key{Char: 'K'},
		Cmd:         commands.FocusTree,
		Description: "Focus tree (up)",
	},
	{
		Key:         keymap.Key{Char: '/'},
		Cmd:         commands.Search,
		Description: "Search files (insert mode)",
	},
	{
		Key:         keymap.Key{Char: 'i'},
		Cmd:         commands.Search,
		Description: "Search files (insert mode)",
	},
	{
		Key:         keymap.Key{Char: 'q'},
		Cmd:         commands.Quit,
		Description: "Quit",
	},
	{
		Key:         keymap.Key{Char: '?'},
		Cmd:         commands.Help,
		Description: "Toggle keybinding cheat sheet",
	},
}

var PreviewKeyMap = keymap.Map{
	{
		Key:         keymap.Key{Char: 'j'},
		Cmd:         commands.ScrollDown,
		Description: "Scroll down",
	},
	{
		Key:         keymap.Key{Code: tcell.KeyDown},
		Cmd:         commands.ScrollDown,
		Description: "Scroll down",
	},
	{
		Key:         keymap.Key{Char: 'k'},
		Cmd:         commands.ScrollUp,
		Description: "Scroll up",
	},
	{
		Key:         keymap.Key{Code: tcell.KeyUp},
		Cmd:         commands.ScrollUp,
		Description: "Scroll up",
	},
	{
		Key:         keymap.Key{Char: 'g'},
		Cmd:         commands.ScrollTop,
		Description: "Scroll to top",
	},
	{
		Key:         keymap.Key{Char: 'G'},
		Cmd:         commands.ScrollBottom,
		Description: "Scroll to bottom",
	},
}
