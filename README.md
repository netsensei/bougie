# Bougie

> *A tiny, sparking terminal browser for the smolweb*

[![License: GPL v3](https://img.shields.io/badge/License-GPLv3-blue.svg)](https://www.gnu.org/licenses/gpl-3.0)
[![Go Version](https://img.shields.io/badge/go-1.23+-00ADD8?logo=go)](https://go.dev/)

Bougie is a lightweight, terminal-based browser for exploring the [Gemini](https://geminiprotocol.net/) and [Gopher](https://en.wikipedia.org/wiki/Gopher_(protocol)) protocols. Built with Go and [Bubble Tea](https://github.com/charmbracelet/bubbletea), it brings the beauty and simplicity of the smolweb to your command line.

## What is the Smolweb?

The **smolweb** is a movement promoting practical, human-scale Internet networking:

- Lightweight protocols (Gemini, Gopher, Finger)
- Small servers with minimal system requirements  
- Intimate, community-focused spaces
- Content over complexity

Bougie is a love letter to this philosophy, a personal passion project celebrating the open source spirit and the joy of simple, text-based browsing.

## Features

### Current

- **Gemini Protocol Support** - Browse gemini:// capsules with full gemtext rendering
- **Gopher Protocol Support** - Navigate classic gopherspace with directory listings
- **Local File Support** - Open and view local gemtext files
- **View Source** - Toggle between rendered and raw document view
- **File Downloads** - Save binary files and documents to disk
- **Configurable Keybindings** - Customize keyboard shortcuts via TOML config
- **Cross-platform** - Works on macOS, Linux, and Windows

### Planned

- Finger protocol support
- Customizable color themes
- Bookmarks management
- Persistent browsing history
- Certificate management (TOFU for Gemini)
- Better error handling and status messages

## Installation

### From Source

Requires Go 1.23 or later:

```bash
git clone https://github.com/netsensei/bougie.git
cd bougie
go build -o bougie
```

### Running

```bash
# Start with default home page
./bougie

# Open a local file
./bougie --file path/to/file
```

## Usage

### Keyboard Shortcuts

| Key | Action |
|-----|--------|
| `Ctrl-n` | Enter navigation mode (address bar) |
| `Ctrl-v` | Return to view mode |
| `Ctrl-u` | View source |
| `Ctrl-h` | Go to home page |
| `Ctrl-r` | Reload current page |
| `b` | Navigate backward |
| `f` | Navigate forward |
| `j` / `↓` | Next link / scroll down |
| `k` / `↑` | Previous link / scroll up |
| `Enter` | Follow selected link / submit input |
| `Ctrl+d` / `Space` | Page down |
| `Ctrl+u` | Page up |
| `Tab` | Navigate forward through links |
| `Shift-Tab` | Navigate backward through links |
| `Tab` | Cycle through dialog components |
| `q` / `Ctrl+c` | Quit |

### Configuration

Bougie creates a configuration file at:

- **Linux/macOS**: `~/.config/bougie/config.toml`
- **Windows**: `%APPDATA%\bougie\config.toml`

Example configuration:

```toml
[general]
home = "gemini://geminiprotocol.net"
downloads_directory = "~/Downloads"

[keybindings]
quit = ["ctrl+c", "ctrl+q"]
nav = "ctrl+n"
view = "ctrl+v"
source = "ctrl+u"
home = "ctrl+h"
reload = "ctrl+r"
# ... customize any key binding
```

## Contributing

Contributions are welcome! This is a personal project, but pull requests, bug reports, and feature suggestions are appreciated.

Please ensure:

- Code follows Go conventions
- Changes are simple, focused and have minimal dependencies
- Updates are tested on at least one platform

## Acknowledgments

- [Bubble Tea](https://github.com/charmbracelet/bubbletea) - The delightful TUI framework
- [Lip Gloss](https://github.com/charmbracelet/lipgloss) - Style definitions for terminal layouts
- [Bubbles](https://github.com/charmbracelet/bubbles) - TUI components
- The Gemini and Gopher communities for keeping the smolweb alive

## Resources

- [Gemini Protocol](https://geminiprotocol.net/)
- [Gopher Protocol (RFC 1436)](https://datatracker.ietf.org/doc/html/rfc1436)
- [Gemini Software](https://geminiprotocol.net/software/)
- [Floodgap Gopher](gopher://gopher.floodgap.com)

## License

This project is licensed under the GNU General Public License v3.0. See the [LICENSE.txt](LICENSE.txt) file for details.