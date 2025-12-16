package tui

import (
	"sort"

	"github.com/mithileshgupta12/vayu/pkg/loader"
)

// Column represents a data column in the table.
type Column struct {
	Key   string
	Width int
}

// detectColumns scans the first N entries to find common keys.
// It prioritizes specific keys like "time", "level", "msg".
func detectColumns(entries []loader.LogEntry) []Column {
	if len(entries) == 0 {
		return nil
	}

	// Scan first 50 entries
	limit := 50
	if len(entries) < limit {
		limit = len(entries)
	}

	keyFrequency := make(map[string]int)
	for i := 0; i < limit; i++ {
		for k := range entries[i] {
			keyFrequency[k]++
		}
	}

	// Filter keys that appear in at least 50% of the sample
	// (Or just take all keys for now, but sort them)
	var keys []string
	for k := range keyFrequency {
		keys = append(keys, k)
	}

	// Custom sort: time > level > msg > others
	sort.Slice(keys, func(i, j int) bool {
		score := func(k string) int {
			switch k {
			case "time", "@timestamp", "ts":
				return 0
			case "level", "severity", "levelname":
				return 1
			case "msg", "message":
				return 2
			case "error":
				return 3
			default:
				return 4
			}
		}
		si, sj := score(keys[i]), score(keys[j])
		if si != sj {
			return si < sj
		}
		return keys[i] < keys[j]
	})

	var cols []Column
	for _, k := range keys {
		cols = append(cols, Column{Key: k, Width: 20}) // Default width
	}
	return cols
}
