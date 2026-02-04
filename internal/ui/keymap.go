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
		Description: "Toggle selection",
	},
}
