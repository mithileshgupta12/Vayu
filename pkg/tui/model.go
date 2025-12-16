package tui

import (
	"bytes"
	"encoding/json"
	"fmt"
	"regexp"
	"strings"

	"github.com/alecthomas/chroma/v2/quick"
	"github.com/charmbracelet/bubbles/textinput"
	"github.com/charmbracelet/bubbles/viewport"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/mithileshgupta12/vayu/pkg/loader"
)

// Styles
var (
	headerStyle = lipgloss.NewStyle().
			Bold(true).
			Foreground(lipgloss.Color("#FAFAFA")).
			Background(lipgloss.Color("#7D56F4")).
			Padding(0, 1)

	cellStyle = lipgloss.NewStyle().
			Padding(0, 1)

	selectedStyle = lipgloss.NewStyle().
			Background(lipgloss.Color("#444"))

	detailStyle = lipgloss.NewStyle().
			Border(lipgloss.NormalBorder(), false, false, false, true). // Left border
			Padding(0, 1)
)

type Model struct {
	Entries         []string
	FilteredEntries []string
	Columns         []Column
	Cursor          int
	Offset          int // Scroll offset
	Width           int
	Height          int
	Ready           bool

	detailView viewport.Model
	textInput  textinput.Model
	filtering  bool

	colsMode  bool
	colCursor int
}

func NewModel(entries []string, sample []loader.LogEntry) Model {
	ti := textinput.New()
	ti.Placeholder = "Regex Filter..."
	ti.CharLimit = 156
	ti.Width = 20

	m := Model{
		Entries:         entries,
		FilteredEntries: entries, // Initially all entries are visible
		Columns:         detectColumns(sample),
		Cursor:          0,
		textInput:       ti,
	}
	m.updateDetail()
	return m
}

func (m *Model) updateDetail() {
	if len(m.FilteredEntries) == 0 {
		m.detailView.SetContent("")
		return
	}
	// Safety check for cursor
	if m.Cursor >= len(m.FilteredEntries) {
		m.Cursor = len(m.FilteredEntries) - 1
	}
	if m.Cursor < 0 {
		m.Cursor = 0
	}

	entryJSON := m.FilteredEntries[m.Cursor]
	// Try to pretty print it
	var entry map[string]interface{}
	var b []byte
	if err := json.Unmarshal([]byte(entryJSON), &entry); err == nil {
		b, _ = json.MarshalIndent(entry, "", "  ")
	} else {
		// fallback
		b = []byte(entryJSON)
	}

	var buf bytes.Buffer
	// Highlight JSON using Chroma
	err := quick.Highlight(&buf, string(b), "json", "terminal16m", "monokai")
	if err != nil {
		m.detailView.SetContent(string(b))
	} else {
		m.detailView.SetContent(buf.String())
	}
}

func (m Model) Init() tea.Cmd {
	return textinput.Blink
}

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd
	var cmds []tea.Cmd

	switch msg := msg.(type) {
	case tea.KeyMsg:
		// Column Mode Handling
		if m.colsMode {
			switch msg.String() {
			case "c", "esc", "q":
				m.colsMode = false
				return m, nil
			case "k", "up":
				if m.colCursor > 0 {
					m.colCursor--
				}
			case "j", "down":
				if m.colCursor < len(m.Columns)-1 {
					m.colCursor++
				}
			case " ", "enter":
				// Toggle visibility
				m.Columns[m.colCursor].Hidden = !m.Columns[m.colCursor].Hidden
			}
			return m, nil
		}

		if m.filtering {
			switch msg.String() {
			case "enter", "esc":
				m.filtering = false
				m.textInput.Blur()
				return m, nil
			}
			m.textInput, cmd = m.textInput.Update(msg)
			cmds = append(cmds, cmd)

			// Live filtering
			m.performFilter()
			return m, tea.Batch(cmds...)
		}

		switch msg.String() {
		case "q", "ctrl+c":
			return m, tea.Quit
		case "c":
			m.colsMode = !m.colsMode
			return m, nil
		case "s":
			// Sort Toggle (Reverse)
			m.reverseEntries()
			return m, nil
		case "/":
			m.filtering = true
			m.textInput.Focus()
			return m, textinput.Blink
		case "k", "up":
			if m.Cursor > 0 {
				m.Cursor--
				if m.Cursor < m.Offset {
					m.Offset--
				}
				m.updateDetail()
			}
		case "j", "down":
			if m.Cursor < len(m.FilteredEntries)-1 {
				m.Cursor++
				if m.Cursor >= m.Offset+m.Height-2 {
					m.Offset++
				}
				m.updateDetail()
			}
		}

	// ... existing WindowSizeMsg and Update logic ...
	case tea.WindowSizeMsg:
		m.Width = msg.Width
		m.Height = msg.Height
		m.Ready = true

		// Layout: 60% List, 40% Detail
		listWidth := int(float64(m.Width) * 0.6)
		detailWidth := m.Width - listWidth - 2 // -2 for border/padding

		// Update viewport size
		// Height - header (1)
		vpHeight := m.Height - 1
		m.detailView = viewport.New(detailWidth, vpHeight)
		m.detailView.Style = detailStyle
		m.updateDetail()
	}
	return m, nil
}

func (m Model) View() string {
	if !m.Ready {
		return "Initializing..."
	}

	if m.colsMode {
		return m.viewCols()
	}

	s := strings.Builder{}

	// Header
	header := ""
	// Calculate total width of columns to ensure we don't exceed list width
	listWidth := int(float64(m.Width) * 0.6)

	for _, col := range m.Columns {
		if col.Hidden {
			continue
		}
		// Truncate or pad header
		h := fmt.Sprintf("%-*s", col.Width, col.Key)
		if len(h) > col.Width {
			h = h[:col.Width]
		}
		header += h + " "
	}

	// Pad header to full list width
	if len(header) < listWidth {
		header += strings.Repeat(" ", listWidth-len(header))
	} else if len(header) > listWidth {
		header = header[:listWidth]
	}

	s.WriteString(headerStyle.Render(header) + "\n")

	// Rows
	// available height for rows = Height - header (1) - footer (optional)

	// Fill remaining height with empty lines to maintain layout
	// Reserve 1 line for footer/search bar if creating a footer
	contentHeight := m.Height - 1 // Header
	if m.filtering || m.textInput.Value() != "" {
		contentHeight-- // input bar
	}

	// Ensure we don't scan past available height
	cnt := 0
	for i := m.Offset; i < len(m.FilteredEntries); i++ {
		if cnt >= contentHeight {
			break
		}

		line := m.FilteredEntries[i]

		// Parse on demand for View
		var entry map[string]interface{}
		_ = json.Unmarshal([]byte(line), &entry) // If error, entry is nil map

		rowStr := ""
		for _, col := range m.Columns {
			if col.Hidden {
				continue
			}
			val := fmt.Sprintf("%v", entry[col.Key])
			if len(val) > col.Width {
				val = val[:col.Width]
			} else {
				val = fmt.Sprintf("%-*s", col.Width, val)
			}
			rowStr += val + " "
		}

		if len(rowStr) < listWidth {
			rowStr += strings.Repeat(" ", listWidth-len(rowStr))
		} else if len(rowStr) > listWidth {
			rowStr = rowStr[:listWidth]
		}

		if i == m.Cursor {
			s.WriteString(selectedStyle.Render(rowStr) + "\n")
		} else {
			s.WriteString(cellStyle.Render(rowStr) + "\n")
		}
		cnt++
	}

	for cnt < contentHeight {
		s.WriteString(strings.Repeat(" ", listWidth) + "\n")
		cnt++
	}

	listView := s.String()
	detailView := m.detailView.View()

	mainView := lipgloss.JoinHorizontal(lipgloss.Top, listView, detailView)

	// Footer (Filter Bar)
	if m.filtering || m.textInput.Value() != "" {
		return lipgloss.JoinVertical(lipgloss.Left, mainView, m.textInput.View())
	}

	return mainView
}

func (m *Model) performFilter() {
	// Simple regex filter
	pattern := m.textInput.Value()
	if pattern == "" {
		m.FilteredEntries = m.Entries
	} else {
		var filtered []string
		re, err := regexp.Compile("(?i)" + pattern) // Case insensitive
		if err != nil {
			return
		}

		for _, str := range m.Entries {
			// Fast filter on raw string
			if re.MatchString(str) {
				filtered = append(filtered, str)
			}
		}
		m.FilteredEntries = filtered
	}
	// Reset cursor and offset
	m.Cursor = 0
	m.Offset = 0
	m.updateDetail()
}

func (m *Model) reverseEntries() {
	for i, j := 0, len(m.FilteredEntries)-1; i < j; i, j = i+1, j-1 {
		m.FilteredEntries[i], m.FilteredEntries[j] = m.FilteredEntries[j], m.FilteredEntries[i]
	}
	m.Cursor = 0
	m.Offset = 0
	m.updateDetail()
}

func (m Model) viewCols() string {
	s := strings.Builder{}
	s.WriteString(headerStyle.Render("Column Management (Press 'Space' to toggle, 'c' to exit)") + "\n\n")

	for i, col := range m.Columns {
		cursor := "  "
		if i == m.colCursor {
			cursor = "> "
		}

		checked := "[x]"
		if col.Hidden {
			checked = "[ ]"
		}

		line := fmt.Sprintf("%s%s %s", cursor, checked, col.Key)
		if i == m.colCursor {
			s.WriteString(selectedStyle.Render(line) + "\n")
		} else {
			s.WriteString(cellStyle.Render(line) + "\n")
		}
	}
	return s.String()
}
