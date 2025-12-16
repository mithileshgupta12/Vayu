# Vayu

**Vayu** is a high-performance, TUI-based JSON log explorer written in Go. It is designed to make reading Newline Delimited JSON (NDJSON) logs effortless and fast, even for large files.

## Features

*   **⚡ Performance**: Lazy-loads log lines, parsing only what is visible. Handles gigabytes of logs with ease.
*   **📊 Auto-Columns**: Automatically detects relevant fields (Time, Level, Msg) and displays them in a table.
*   **👁️ Split View**: Navigate the list on the left, see full pretty-printed JSON details on the right.
*   **🌈 Syntax Highlighting**: Beautiful JSON rendering with Chroma.
*   **🔍 Regex Filtering**: Type `/` to filter logs instantly with regex.
*   **🛠️ Column Management**: Type `c` to toggle column visibility.
*   **jj Toggle Sort**: Type `s` to reverse sort order (Newest/Oldest).

## Installation

```bash
go install github.com/mithileshgupta12/vayu/cmd/vayu@latest
```

## Usage

```bash
vayu <path-to-logfile>
```

## Keybindings

| Key | Action |
| --- | --- |
| `j` / `↓` | Move down |
| `k` / `↑` | Move up |
| `/` | Start Filtering |
| `c` | Enter Column Management Mode |
| `s` | Toggle Sort Order (Reverse) |
| `Space` | Toggle Column (in Column Mode) |
| `q` | Quit |

## License
MIT
