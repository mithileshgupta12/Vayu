# Vayu Roadmap - 30 Day Challenge

## 🚀 Core Functionality (Days 1-7)
- [ ] **Live Tailing**: Implement `-f` flag to follow logs in real-time (like `tail -f`).
- [ ] **Clipboard Support**: Press `y` to copy the selected log line (or specific field) to clipboard.
- [ ] **Log Level Coloring**: detailed coloring rules (Red for ERROR, Green for INFO, Yellow for WARN).
- [ ] **JSON Key Navigation**: Allow tabbing through keys in the Detail View.
- [ ] **Error Handling**: Graceful fallback if file is deleted/rotated.
- [ ] **Performance**: optimize for files > 1GB (mmap?).
- [ ] **Wrapping**: Toggle line wrapping in the Detail View.

## 🛠️ User Experience (Days 8-14)
- [ ] **Search History**: Restore previous search queries with Up/Down arrows in filter bar.
- [ ] **Theme Support**: Allow selecting different Chroma themes (Dracula, Solarized).
- [ ] **Help Modal**: Press `?` to show a popup with keybindings.
- [ ] **Line Numbers**: Toggle line numbers in the list view.
- [ ] **Jump to Line**: Press `:` and type line number to jump.
- [ ] **Status Bar**: Show file size and percentage read.
- [ ] **Custom Columns**: Allow user to specify which keys to show via flags (e.g., `-cols=time,msg,user_id`).

## ⚙️ Advanced Features (Days 15-30)
- [ ] **Export**: Press `e` to save the currently *filtered* logs to a new file.
- [ ] **Time Filter**: Filter by time range (e.g., `>10m ago`).
- [ ] **Saved Views**: Save column configurations as named presets.
- [ ] **Flattening**: Option to flatten nested JSON objects in the list view.
- [ ] **Docker Integration**: `vayu docker <container-id>` implementation.
- [ ] **SSH Mode**: Stream logs from a remote server.
- [ ] **Config File**: `~/.vayu/config.yaml` for persistent settings.
