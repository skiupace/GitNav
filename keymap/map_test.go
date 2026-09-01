package keymap

import (
	"testing"

	"github.com/gdamore/tcell/v2"
	"github.com/skiupace/gitnav/commands"
)

func TestMap_Resolve(t *testing.T) {
	m := Map{
		{Key: Key{Char: 'q'}, Cmd: commands.Quit},
		{Key: Key{Code: tcell.KeyEnter}, Cmd: commands.Select},
	}

	tests := []struct {
		name  string
		event *tcell.EventKey
		want  commands.Command
	}{
		{
			name:  "Quit command with 'q'",
			event: tcell.NewEventKey(tcell.KeyRune, 'q', tcell.ModNone),
			want:  commands.Quit,
		},
		{
			name:  "Select command with Enter",
			event: tcell.NewEventKey(tcell.KeyEnter, 0, tcell.ModNone),
			want:  commands.Select,
		},
		{
			name:  "Noop for unknown key",
			event: tcell.NewEventKey(tcell.KeyRune, 'x', tcell.ModNone),
			want:  commands.Noop,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := m.Resolve(tt.event); got != tt.want {
				t.Errorf("Map.Resolve() = %v, want %v", got, tt.want)
			}
		})
	}
}
