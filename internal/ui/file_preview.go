package ui

import (
	"fmt"
	"math"
	"os"
	"path/filepath"
	"strings"
	"unicode/utf8"

	"github.com/alecthomas/chroma/v2"
	"github.com/alecthomas/chroma/v2/lexers"
	"github.com/alecthomas/chroma/v2/styles"
	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
	"github.com/skiupace/gitnav/commands"
)

type PreviewPanel struct {
	TextView *tview.TextView
}

func NewPreviewPanel() *PreviewPanel {
	tv := tview.NewTextView().
		SetDynamicColors(true).
		SetRegions(false).
		SetWordWrap(false).
		SetScrollable(true)

	tv.SetBorder(true).
		SetTitle(" Preview: <select-file> ").
		SetTitleAlign(tview.AlignLeft).
		SetBackgroundColor(tcell.ColorDefault).
		SetBorderColor(tcell.ColorGray)

	pp := &PreviewPanel{TextView: tv}

	// Draw a visual scrollbar on the right inner edge
	tv.SetDrawFunc(func(screen tcell.Screen, x, y, width, height int) (int, int, int, int) {
		innerX := x + 1
		innerY := y + 1
		innerW := width - 2
		innerH := height - 2

		totalLines := tv.GetOriginalLineCount()
		if totalLines > innerH {
			scrollRow, _ := tv.GetScrollOffset()

			// Thumb height proportional to visible/total ratio
			thumbHeight := innerH * innerH / totalLines
			if thumbHeight < 1 {
				thumbHeight = 1
			}

			// Thumb position along the track
			maxOffset := totalLines - innerH
			thumbPos := 0
			if maxOffset > 0 {
				thumbPos = scrollRow * (innerH - thumbHeight) / maxOffset
			}

			trackX := innerX + innerW // rightmost inner column
			trackStyle := tcell.StyleDefault.Foreground(tcell.ColorGray)
			thumbStyle := tcell.StyleDefault.Foreground(tcell.ColorWhite)

			for i := 0; i < innerH; i++ {
				style := trackStyle
				ch := '│'
				if i >= thumbPos && i < thumbPos+thumbHeight {
					style = thumbStyle
					ch = '█'
				}
				screen.SetContent(trackX, innerY+i, ch, nil, style)
			}
		}

		return innerX, innerY, innerW, innerH
	})

	tv.SetInputCapture(scrollCapture(tv))

	return pp
}

// scrollCapture translates vim scroll keys into TextView scrolling. It is
// shared with the help panel.
func scrollCapture(tv *tview.TextView) func(event *tcell.EventKey) *tcell.EventKey {
	return func(event *tcell.EventKey) *tcell.EventKey {
		switch PreviewKeyMap.Resolve(event) {
		case commands.ScrollDown:
			row, col := tv.GetScrollOffset()
			tv.ScrollTo(row+1, col)
			return nil
		case commands.ScrollUp:
			row, col := tv.GetScrollOffset()
			if row > 0 {
				tv.ScrollTo(row-1, col)
			}
			return nil
		case commands.ScrollTop:
			tv.ScrollToBeginning()
			return nil
		case commands.ScrollBottom:
			tv.ScrollToEnd()
			return nil
		}
		return event
	}
}

func (pp *PreviewPanel) UpdatePreview(filePath string) {
	pp.TextView.SetTitle(fmt.Sprintf(" Preview: %s ", filepath.Base(filePath)))

	content, err := os.ReadFile(filePath)
	if err != nil {
		pp.TextView.SetText("[red]Error reading file[-]")
		return
	}

	if isBinaryContent(content) {
		pp.TextView.SetText("[gray](binary file — cannot preview)[-]")
		return
	}

	text := string(content)
	highlighted := highlightContent(text, filePath)
	pp.TextView.SetText(highlighted)
	pp.TextView.ScrollToBeginning()
}

func (pp *PreviewPanel) ClearPreview() {
	pp.TextView.SetTitle(" Preview: <select-file> ")
	pp.TextView.SetText("")
}

func isBinaryContent(data []byte) bool {
	if len(data) == 0 {
		return false
	}

	checkSize := len(data)
	if checkSize > 8192 {
		checkSize = 8192
	}
	sample := data[:checkSize]

	for _, b := range sample {
		if b == 0 {
			return true
		}
	}

	return !utf8.Valid(sample)
}

// highlightContent uses chroma to syntax-highlight the file content and
// returns tview-compatible color-tagged text with line numbers. Token colors
// are mapped to the terminal's 16-color ANSI palette, so highlighting
// follows the user's terminal theme instead of a fixed truecolor scheme.
func highlightContent(content, filePath string) string {
	lexer := lexers.Match(filePath)
	if lexer == nil {
		lexer = lexers.Analyse(content)
	}

	if lexer == nil {
		lexer = lexers.Fallback
	}
	lexer = chroma.Coalesce(lexer)

	style := styles.Get("monokai")

	iterator, err := lexer.Tokenise(nil, content)
	if err != nil {
		return addLineNumbers(tview.Escape(content))
	}

	var sb strings.Builder
	for _, token := range iterator.Tokens() {
		entry := style.Get(token.Type)
		text := tview.Escape(token.Value)

		tag := styleTag(entry)
		if tag != "" {
			sb.WriteString(tag)
		}
		sb.WriteString(text)
		if tag != "" {
			sb.WriteString("[-]")
		}
	}

	return addLineNumbers(sb.String())
}

// ansi16 mirrors chroma's own 16-color reference palette (xterm defaults,
// same as its terminal16 formatter), mapped to tcell color names. tcell
// resolves these names to palette indices, so they follow the terminal
// theme.
var ansi16 = []struct {
	hex  chroma.Colour
	name string
}{
	{chroma.MustParseColour("#000000"), "black"},
	{chroma.MustParseColour("#7f0000"), "maroon"},
	{chroma.MustParseColour("#007f00"), "green"},
	{chroma.MustParseColour("#7f7fe0"), "olive"},
	{chroma.MustParseColour("#00007f"), "navy"},
	{chroma.MustParseColour("#7f007f"), "purple"},
	{chroma.MustParseColour("#007f7f"), "teal"},
	{chroma.MustParseColour("#e5e5e5"), "silver"},
	{chroma.MustParseColour("#555555"), "gray"},
	{chroma.MustParseColour("#ff0000"), "red"},
	{chroma.MustParseColour("#00ff00"), "lime"},
	{chroma.MustParseColour("#ffff00"), "yellow"},
	{chroma.MustParseColour("#0000ff"), "blue"},
	{chroma.MustParseColour("#ff00ff"), "fuchsia"},
	{chroma.MustParseColour("#00ffff"), "aqua"},
	{chroma.MustParseColour("#ffffff"), "white"},
}

// nearestANSI returns the name of the palette color closest to c.
func nearestANSI(c chroma.Colour) string {
	name := "silver"
	best := math.MaxFloat64
	for _, e := range ansi16 {
		if d := e.hex.Distance(c); d < best {
			best, name = d, e.name
		}
	}
	return name
}

// styleTag builds a tview color tag for a chroma style entry:
// [foreground::flags] with flags b/i/u. Returns "" for default styling.
func styleTag(entry chroma.StyleEntry) string {
	flags := ""
	if entry.Bold == chroma.Yes {
		flags += "b"
	}
	if entry.Italic == chroma.Yes {
		flags += "i"
	}
	if entry.Underline == chroma.Yes {
		flags += "u"
	}

	fg := ""
	if entry.Colour.IsSet() {
		fg = nearestANSI(entry.Colour)
	}

	if fg == "" && flags == "" {
		return ""
	}
	return "[" + fg + "::" + flags + "]"
}

// addLineNumbers prepends line numbers to each line of the content.
// Handles tview color tags correctly when splitting lines.
func addLineNumbers(content string) string {
	lines := strings.Split(content, "\n")
	totalLines := len(lines)

	// Calculate gutter width based on total line count
	gutterWidth := len(fmt.Sprintf("%d", totalLines))
	if gutterWidth < 3 {
		gutterWidth = 3
	}

	var sb strings.Builder
	for i, line := range lines {
		lineNum := i + 1
		numStr := fmt.Sprintf("%*d", gutterWidth, lineNum)

		sb.WriteString(fmt.Sprintf("[gray]%s │[-] %s", numStr, line))

		if i < len(lines)-1 {
			sb.WriteString("\n")
		}
	}

	return sb.String()
}
