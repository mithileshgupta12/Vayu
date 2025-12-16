package tui

import (
	"bytes"
	"encoding/json"
	"fmt"
	"strings"

	"github.com/alecthomas/chroma/v2/quick"
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
	Entries []loader.LogEntry
	Columns []Column
	Cursor  int
	Offset  int // Scroll offset
	Width   int
	Height  int
	Ready   bool

	detailView viewport.Model
}

func NewModel(entries []loader.LogEntry) Model {
	m := Model{
		Entries: entries,
		Columns: detectColumns(entries),
		Cursor:  0,
	}
	m.updateDetail()
	return m
}

func (m *Model) updateDetail() {
	if len(m.Entries) == 0 {
		return
	}
	entry := m.Entries[m.Cursor]
	b, _ := json.MarshalIndent(entry, "", "  ")

	var buf bytes.Buffer
	// Highlight JSON using Chroma
	// "terminal16m" = True Color, "monokai" = style
	err := quick.Highlight(&buf, string(b), "json", "terminal16m", "monokai")
	if err != nil {
		m.detailView.SetContent(string(b))
	} else {
		m.detailView.SetContent(buf.String())
	}
}

func (m Model) Init() tea.Cmd {
	return nil
}

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.MouseMsg:
		var cmd tea.Cmd
		m.detailView, cmd = m.detailView.Update(msg)
		return m, cmd

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
				m.updateDetail()
			}
		case "j", "down":
			if m.Cursor < len(m.Entries)-1 {
				m.Cursor++
				// TODO: Better scrolling logic based on height
				if m.Cursor >= m.Offset+m.Height-2 { // -2 for header
					m.Offset++
				}
				m.updateDetail()
			}
		}
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

	s := strings.Builder{}

	// Header
	header := ""
	// Calculate total width of columns to ensure we don't exceed list width
	listWidth := int(float64(m.Width) * 0.6)

	for _, col := range m.Columns {
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

		// Pad row to full list width
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

	// Fill remaining height with empty lines to maintain layout
	for cnt < m.Height-1 {
		s.WriteString(strings.Repeat(" ", listWidth) + "\n")
		cnt++
	}

	listView := s.String()
	detailView := m.detailView.View()

	return lipgloss.JoinHorizontal(lipgloss.Top, listView, detailView)
}
