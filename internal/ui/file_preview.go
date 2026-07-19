package ui

import (
	"fmt"
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

// PreviewPanel holds the preview text view and exposes update methods.
type PreviewPanel struct {
	TextView *tview.TextView
}

// NewPreviewPanel creates a scrollable, syntax-highlighted preview panel.
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

	// Input capture for vim scrolling and arrow keys
	tv.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		cmd := PreviewKeyMap.Resolve(event)
		switch cmd {
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
	})

	return pp
}

// UpdatePreview reads a file, applies syntax highlighting, and displays it.
func (pp *PreviewPanel) UpdatePreview(filePath string) {
	pp.TextView.SetTitle(fmt.Sprintf(" Preview: %s ", filepath.Base(filePath)))

	content, err := os.ReadFile(filePath)
	if err != nil {
		pp.TextView.SetText("[red]Error reading file[-]")
		return
	}

	// Check for binary content
	if isBinaryContent(content) {
		pp.TextView.SetText("[gray](binary file — cannot preview)[-]")
		return
	}

	text := string(content)
	highlighted := highlightContent(text, filePath)
	pp.TextView.SetText(highlighted)
	pp.TextView.ScrollToBeginning()
}

// ClearPreview resets the preview panel to its default state.
func (pp *PreviewPanel) ClearPreview() {
	pp.TextView.SetTitle(" Preview: <select-file> ")
	pp.TextView.SetText("")
}

// isBinaryContent checks if content contains null bytes or is not valid UTF-8.
func isBinaryContent(data []byte) bool {
	if len(data) == 0 {
		return false
	}
	// Check a sample of the content
	checkSize := len(data)
	if checkSize > 8192 {
		checkSize = 8192
	}
	sample := data[:checkSize]

	// Check for null bytes
	for _, b := range sample {
		if b == 0 {
			return true
		}
	}

	// Check for valid UTF-8
	return !utf8.Valid(sample)
}

// highlightContent uses chroma to syntax-highlight the file content
// and returns tview-compatible color-tagged text with line numbers.
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

		if entry.Colour.IsSet() {
			hex := fmt.Sprintf("#%02x%02x%02x", entry.Colour.Red(), entry.Colour.Green(), entry.Colour.Blue())
			sb.WriteString(fmt.Sprintf("[%s]%s[-]", hex, text))
		} else {
			sb.WriteString(text)
		}
	}

	return addLineNumbers(sb.String())
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
		// Gray line number + separator + content
		numStr := fmt.Sprintf("%*d", gutterWidth, lineNum)
		sb.WriteString(fmt.Sprintf("[#666666]%s │[-] %s", numStr, line))
		if i < len(lines)-1 {
			sb.WriteString("\n")
		}
	}

	return sb.String()
}
