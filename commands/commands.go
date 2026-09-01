package commands

type Command uint8

const (
	Noop Command = iota

	MoveUp
	MoveDown
	MoveLeft
	MoveRight

	Search
	Quit
	ExpandAll
	FocusSidebar
	UnfocusSidebar
	ToggleSidebar
	Select

	ScrollUp
	ScrollDown
	ScrollTop
	ScrollBottom

	OpenEditor

	FocusPreview
	FocusTree

	Help
)

func (c Command) String() string {
	switch c {
	case Noop:
		return "Noop"

	case MoveUp:
		return "MoveUp"
	case MoveDown:
		return "MoveDown"
	case MoveLeft:
		return "MoveLeft"
	case MoveRight:
		return "MoveRight"

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

	case ScrollUp:
		return "ScrollUp"
	case ScrollDown:
		return "ScrollDown"
	case ScrollTop:
		return "ScrollTop"
	case ScrollBottom:
		return "ScrollBottom"

	case OpenEditor:
		return "OpenEditor"

	case FocusPreview:
		return "FocusPreview"
	case FocusTree:
		return "FocusTree"

	case Help:
		return "Help"
	}

	return "Unknown"
}
