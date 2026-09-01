package commands

import "testing"

func TestCommand_String(t *testing.T) {
	tests := []struct {
		cmd  Command
		want string
	}{
		{Noop, "Noop"},
		{MoveUp, "MoveUp"},
		{Quit, "Quit"},
		{Command(99), "Unknown"},
	}

	for _, tt := range tests {
		t.Run(tt.want, func(t *testing.T) {
			if got := tt.cmd.String(); got != tt.want {
				t.Errorf("Command.String() = %v, want %v", got, tt.want)
			}
		})
	}
}
