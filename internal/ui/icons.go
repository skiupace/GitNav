package ui

import (
	"path/filepath"
	"strings"

	"github.com/gdamore/tcell/v2"
)

// GetIcon returns a Nerd Font icon and color for a given filename and type
func GetIcon(filename string, isDir bool) (string, tcell.Color) {
	if isDir {
		return GetFolderIcon(false), tcell.ColorBlue
	}

	ext := strings.ToLower(filepath.Ext(filename))
	fname := strings.ToLower(filename)

	// Specific filenames
	switch fname {
	case "go.mod", "go.sum":
		return "", tcell.ColorAqua // Cyan substitute
	case "makefile":
		return "", tcell.ColorOlive // Brown substitute
	case "dockerfile":
		return "", tcell.ColorBlue
	case "license", "license.md", "readme.md":
		return "", tcell.ColorYellow
	case ".gitignore", ".gitattributes":
		return "", tcell.ColorRed // OrangeRed substitute
	}

	// Extensions
	switch ext {
	// Languages
	case ".go":
		if strings.HasSuffix(fname, "_test.go") {
			return "", tcell.ColorGreen // Test files
		}
		return "", tcell.ColorAqua
	case ".js":
		return "", tcell.ColorYellow
	case ".ts", ".tsx":
		return "", tcell.ColorBlue
	case ".py":
		return "", tcell.ColorYellow
	case ".rs":
		return "", tcell.ColorRed
	case ".c", ".h":
		return "", tcell.ColorBlue
	case ".cpp", ".hpp":
		return "", tcell.ColorBlue
	case ".java":
		return "", tcell.ColorRed
	case ".php":
		return "", tcell.ColorLightCyan
	case ".html":
		return "", tcell.ColorOrange
	case ".css":
		return "", tcell.ColorBlue
	case ".scss", ".sass":
		return "", tcell.ColorFuchsia
	case ".json":
		return "", tcell.ColorYellow
	case ".yaml", ".yml":
		return "", tcell.ColorFuchsia
	case ".md":
		return "", tcell.ColorWhite
	case ".txt":
		return "", tcell.ColorWhite
	case ".sh", ".bash", ".zsh":
		return "", tcell.ColorGreen
	case ".sql":
		return "", tcell.ColorYellow

	// Images
	case ".png", ".jpg", ".jpeg", ".gif", ".svg", ".ico", ".webp":
		return "", tcell.ColorPurple

	// Archives
	case ".zip", ".tar", ".gz", ".7z", ".rar":
		return "", tcell.ColorRed
	}

	// Default file
	return "", tcell.ColorWhite
}

// GetFolderIcon returns the icon for a folder based on its expansion state.
// expanded: true -> Open Folder Icon
// expanded: false -> Closed Folder Icon
func GetFolderIcon(expanded bool) string {
	if expanded {
		// Expanded -> Open Folder
		return "" // nf-custom-folder_open
	}
	// Collapsed -> Closed Folder
	return "" // nf-custom-folder
}
