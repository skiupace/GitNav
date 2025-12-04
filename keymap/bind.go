package keymap

import (
	"fmt"

	"github.com/skiupace/gitnav/commands"
)

// Struct that holds a key and a command
type Bind struct {
	Key         Key
	Cmd         commands.Command
	Description string
}

func (b Bind) String() string {
	return fmt.Sprintf("%s = %s", b.Key.String(), b.Cmd.String())
}
