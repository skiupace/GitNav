# GitNav

A powerful terminal-based Git repository navigator built with Go and [tview](https://github.com/rivo/tview). GitNav provides an intuitive TUI (Text User Interface) for exploring Git repositories, searching files, viewing file contents, all from the comfort of your terminal.

![Go Version](https://img.shields.io/badge/Go-1.21+-00ADD8?style=flat&logo=go)
![License](https://img.shields.io/badge/License-MIT-blue.svg)
![TUI](https://img.shields.io/badge/TUI-tview-2E8B57?style=flat)

## Features

- **Repository Navigation**: Browse your Git repository structure with an interactive tree view
- **File Preview**: View file contents directly in the terminal with syntax highlighting support
- **Powerful Search**: Search across files by name pattern matching
- **Git Integration**: Display branch information, and repository status \[clean, dirty\]
- **Keyboard-Driven**: Efficient navigation with vim-style keybindings
- **Statistics Dashboard**: View basic repository statistics including commit information
- **Help System**: Built-in help panel showing all available keybindings (_coming soon_)

## Installation

### Prerequisites

- Go 1.21 or higher
- Git (for repository operations)

### Build from Source

```bash
# Clone the repository
git clone https://github.com/skiupace/gitnav.git -b gitnav-go
cd gitnav

# Build the binary
go build .

# Install globally (optional)
go install -o $GOPATH/bin/gitnav
```

### Quick Start

```bash
# Navigate to a Git repository
cd /path/to/your/repo

# Run GitNav
gitnav

# Or specify a repository path
gitnav /path/to/your/repo
```

## Usage

### Basic Navigation

Once GitNav starts, you'll see a split-screen interface with:

- **Left Panel**: Repository tree structure showing all files and directories
- **Right Panel**: File preview, search results, or help information
- **Bottom Bar**: Repository statistics and current branch information

### Keybindings

GitNav uses vim-style keybindings for efficient navigation:

#### Navigation

| Key       | Action                                  |
| --------- | --------------------------------------- |
| `j` / `↓` | Move down                               |
| `k` / `↑` | Move up                                 |
| `h` / `←` | Collapse directory / Go back (_soon_)   |
| `l` / `→` | Expand directory / Select file (_soon_) |
| `g`       | Go to top                               |
| `G`       | Go to bottom                            |

#### File Operations (_soon_)

| Key     | Action                                                     |
| ------- | ---------------------------------------------------------- |
| `Enter` | Open file in preview                                       |
| `o`     | Open file in external editor (uses `$EDITOR` or `$VISUAL`) |
| `y`     | Copy file path to clipboard                                |
| `f`     | Toggle file preview                                        |

#### Search

| Key   | Action                          |
| ----- | ------------------------------- |
| `/`   | Start search mode               |
| `n`   | Next search result (_soon_)     |
| `N`   | Previous search result (_soon_) |
| `Esc` | Clear search / Exit search mode |

#### View Controls

| Key | Action |
|-----|--------|
| `Tab` | Switch between panels |
| `r` | Refresh repository view |
| `H` | Toggle help panel |
| `q` / `Ctrl+c` | Quit application |

#### Git Operations (_soon_)

| Key | Action |
|-----|--------|
| `b` | Show branch information |
| `c` | Show commit history |
| `s` | Show Git status |

## Configuration (_soon_)

GitNav respects the following environment variables:

- `EDITOR` or `VISUAL`: Default editor for opening files (default: `vim`)
- `GITNAV_CONFIG`: Path to configuration file (future feature)

## Development

### Running Tests

```bash
# Run all tests
go test ./...

# Run tests with coverage
go test -cover ./...

# Run specific package tests
go test ./internal/git/...
go test ./internal/ui/...
```

### Building for Production

```bash
# Build optimized binary
go build -ldflags="-s -w" -o gitnav

# Build for multiple platforms
GOOS=linux GOARCH=amd64 go build -o gitnav-linux-amd64
GOOS=darwin GOARCH=amd64 go build -o gitnav-darwin-amd64
GOOS=windows GOARCH=amd64 go build -o gitnav-windows-amd64.exe
```

### Code Style

This project follows standard Go conventions:

- Use `go fmt` for code formatting
- Run `go vet` for static analysis
- Keep functions small and focused
- Write tests for all public APIs

## Dependencies

- [tview](https://github.com/rivo/tview) - Terminal UI toolkit
- [tcell](https://github.com/rivo/tcell) - Terminal cell library
- [golang.org/x/term](https://pkg.go.dev/golang.org/x/term) - Terminal utilities

## Roadmap

- [ ] Mouse support improvements
- [ ] Commit history to jump between older ones
- [ ] Configuration file support (YAML/TOML)
- [ ] Syntax highlighting for more languages
- [ ] Improved search for content with fuzzy matching
- [ ] Multi-repository workspace support
- [ ] File diff viewer

## Contributing

Contributions are welcome! Please feel free to submit a Pull Request.

1. Fork the repository
2. Create your feature branch (`git checkout -b feature/amazing-feature`)
3. Commit your changes (`git commit -m 'Add some amazing feature'`)
4. Push to the branch (`git push origin feature/amazing-feature`)
5. Open a Pull Request

## License

This project is licensed under the [MIT License](LICENSE) - see the LICENSE file for details.

## Acknowledgments

- [tview](https://github.com/rivo/tview) for the excellent TUI framework
- [tcell](https://github.com/rivo/tcell) for terminal handling
- The Go community for continuous support and inspiration

## Support

If you encounter any issues or have questions:

1. Check existing issues in the repository
2. Create a new issue with detailed information
3. Include your Go version, OS, and GitNav version
