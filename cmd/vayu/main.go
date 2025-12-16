package main

import (
	"fmt"
	"os"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/mithileshgupta12/vayu/pkg/loader"
	"github.com/mithileshgupta12/vayu/pkg/tui"
)

func main() {
	if len(os.Args) < 2 {
		fmt.Println("Usage: vayu <logfile>")
		os.Exit(1)
	}

	path := os.Args[1]
	fmt.Printf("Loading logs from %s...\n", path)

	// Load raw lines + sample for column detection
	lines, sample, err := loader.LoadFile(path)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error loading logs: %v\n", err)
		os.Exit(1)
	}

	fmt.Printf("Loaded %d log lines. Starting TUI...\n", len(lines))

	m := tui.NewModel(lines, sample)
	p := tea.NewProgram(m, tea.WithAltScreen())

	if _, err := p.Run(); err != nil {
		fmt.Fprintf(os.Stderr, "Error running program: %v\n", err)
		os.Exit(1)
	}
}
