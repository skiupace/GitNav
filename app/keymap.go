package app

import (
	"github.com/gdamore/tcell/v2"

	cmd "github.com/skiupace/gitnav/commands"
	"github.com/skiupace/gitnav/keymap"
)

// local alias added for clarity purpose
type (
	Bind = keymap.Bind
	Key  = keymap.Key
	Map  = keymap.Map
)

// KeymapSystem is the actual key mapping system.
// A map can have several groups. But it always has a "Global" one.
type KeymapSystem struct {
	Groups map[string]Map
	Global Map
}

func (c KeymapSystem) Group(name string) Map {
	// Lookup the group
	if group, ok := c.Groups[name]; ok {
		return group
	}

	// Did not find any maps. Return a empty one
	return Map{}
}

// Resolve translates a tcell.EventKey into a command based on the mappings in
// the global group
func (c KeymapSystem) Resolve(event *tcell.EventKey) cmd.Command {
	return c.Global.Resolve(event)
}

const (
	HomeGroup   = "home"
	TreeGroup   = "tree"
	SearchGroup = "search"
)

// Define a global KeymapSystem object with default keybinds
var Keymaps = KeymapSystem{
	Groups: map[string]Map{
		HomeGroup: {
			Bind{Key: Key{Char: 'H'}, Cmd: cmd.MoveLeft, Description: "Focus tree"},
			Bind{Key: Key{Char: 'q'}, Cmd: cmd.Quit, Description: "Quit"},
		},
		TreeGroup: {
			Bind{Key: Key{Char: 'j'}, Cmd: cmd.MoveDown, Description: "Go down"},
			Bind{Key: Key{Code: tcell.KeyDown}, Cmd: cmd.MoveDown, Description: "Go down"},
			Bind{Key: Key{Char: 'k'}, Cmd: cmd.MoveUp, Description: "Go up"},
			Bind{Key: Key{Code: tcell.KeyUp}, Cmd: cmd.MoveUp, Description: "Go up"},
			Bind{Key: Key{Char: '/'}, Cmd: cmd.Search, Description: "Search"},
			Bind{Key: Key{Char: 'e'}, Cmd: cmd.ExpandAll, Description: "Expand all"},
		},
	},
}
