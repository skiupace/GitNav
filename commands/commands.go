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
	Select

	// Preview scrolling
	ScrollUp
	ScrollDown
	ScrollTop
	ScrollBottom

	// Focus switching
	FocusPreview
	FocusTree
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
		return "MoveLeft"
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
	case Select:
		return "Select"

	// Preview scrolling
	case ScrollUp:
		return "ScrollUp"
	case ScrollDown:
		return "ScrollDown"
	case ScrollTop:
		return "ScrollTop"
	case ScrollBottom:
		return "ScrollBottom"

	// Focus switching
	case FocusPreview:
		return "FocusPreview"
	case FocusTree:
		return "FocusTree"
	}

	return "Unknown"
}
