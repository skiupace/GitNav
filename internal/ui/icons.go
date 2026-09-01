package ui

import (
	"path/filepath"
	"strings"

	"github.com/gdamore/tcell/v2"
)

func GetIcon(filename string, isDir bool) (string, tcell.Color) {
	if isDir {
		return GetFolderIcon(false), tcell.ColorBlue
	}

	ext := strings.ToLower(filepath.Ext(filename))
	fname := strings.ToLower(filename)

	switch fname {
	case "go.mod", "go.sum":
		return "", tcell.ColorAqua
	case "makefile":
		return "", tcell.ColorOlive
	case "dockerfile":
		return "", tcell.ColorBlue
	case "license", "license.md", "readme.md":
		return "", tcell.ColorYellow
	case ".gitignore", ".gitattributes":
		return "", tcell.ColorRed
	}

	switch ext {
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

	case ".png", ".jpg", ".jpeg", ".gif", ".svg", ".ico", ".webp":
		return "", tcell.ColorPurple

	case ".zip", ".tar", ".gz", ".7z", ".rar":
		return "", tcell.ColorRed
	}

	return "", tcell.ColorWhite
}

func GetFolderIcon(expanded bool) string {
	if expanded {
		return ""
	}
	return ""
}
