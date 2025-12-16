package tui

import (
	"fmt"
	"strings"

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
)

type Model struct {
	Entries []loader.LogEntry
	Columns []Column
	Cursor  int
	Offset  int // Scroll offset
	Width   int
	Height  int
	Ready   bool
}

func NewModel(entries []loader.LogEntry) Model {
	return Model{
		Entries: entries,
		Columns: detectColumns(entries),
		Cursor:  0,
	}
}

func (m Model) Init() tea.Cmd {
	return nil
}

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "q", "ctrl+c":
			return m, tea.Quit
		case "k", "up":
			if m.Cursor > 0 {
				m.Cursor--
				if m.Cursor < m.Offset {
					m.Offset--
				}
			}
		case "j", "down":
			if m.Cursor < len(m.Entries)-1 {
				m.Cursor++
				// TODO: Better scrolling logic based on height
				if m.Cursor >= m.Offset+m.Height-2 { // -2 for header
					m.Offset++
				}
			}
		}
	case tea.WindowSizeMsg:
		m.Width = msg.Width
		m.Height = msg.Height
		m.Ready = true
	}
	return m, nil
}

func (m Model) View() string {
	if !m.Ready {
		return "Initializing..."
	}

	s := strings.Builder{}

	// Header
	header := ""
	for _, col := range m.Columns {
		// Truncate or pad header
		h := fmt.Sprintf("%-*s", col.Width, col.Key)
		if len(h) > col.Width {
			h = h[:col.Width]
		}
		header += h + " "
	}
	s.WriteString(headerStyle.Render(header) + "\n")

	// Rows
	// available height for rows = Height - header (1) - footer (optional)
	cnt := 0
	for i := m.Offset; i < len(m.Entries); i++ {
		if cnt >= m.Height-1 {
			break
		}

		entry := m.Entries[i]
		rowStr := ""
		for _, col := range m.Columns {
			val := fmt.Sprintf("%v", entry[col.Key])
			// Truncate or pad
			if len(val) > col.Width {
				val = val[:col.Width]
			} else {
				val = fmt.Sprintf("%-*s", col.Width, val)
			}
			rowStr += val + " "
		}

		if i == m.Cursor {
			s.WriteString(selectedStyle.Render(rowStr) + "\n")
		} else {
			s.WriteString(cellStyle.Render(rowStr) + "\n")
		}
		cnt++
	}

	return s.String()
}
