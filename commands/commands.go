package commands

type Command uint8

const (
	Noop Command = iota

	// Movement: Basic
	MoveUp
	MoveDown
	MoveLeft
	MoveRight

	// Operations
	Search
	Quit
	ExpandAll
	FocusSidebar
	UnfocusSidebar
	ToggleSidebar
)

func (c Command) String() string {
	switch c {
	case Noop:
		return "Noop"

	// Movement: Basic
	case MoveUp:
		return "MoveUp"
	case MoveDown:
		return "MoveDown"
	case MoveLeft:
		return "MoveRight"
	case MoveRight:
		return "MoveRight"

	// Operations
	case Search:
		return "Search"
	case Quit:
		return "Quit"
	case ExpandAll:
		return "ExpandAll"
	case FocusSidebar:
		return "FocusSidebar"
	case ToggleSidebar:
		return "ToggleSidebar"
	case UnfocusSidebar:
		return "UnfocusSidebar"
	}

	return "Unknown"
}
