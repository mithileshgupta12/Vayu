package loader

import (
	"bufio"
	"encoding/json"
	"io"
	"os"
)

// LogEntry represents a single parsed log line.
// We use map[string]interface{} to be flexible with any JSON structure.
type LogEntry map[string]interface{}

// LoadLogs reads NDJSON from the given reader and returns a slice of LogEntries.
func LoadLogs(r io.Reader) ([]LogEntry, error) {
	var entries []LogEntry
	scanner := bufio.NewScanner(r)

	// Increase buffer size for large log lines (default is 64KB)
	buf := make([]byte, 0, 64*1024)
	scanner.Buffer(buf, 10*1024*1024) // 10 MB max line size

	for scanner.Scan() {
		line := scanner.Bytes()
		if len(line) == 0 {
			continue
		}

		var entry LogEntry
		if err := json.Unmarshal(line, &entry); err != nil {
			// If it's not JSON, we might want to wrap it in a struct
			// or just skip it. For now, let's treat it as a raw message.
			entry = LogEntry{"raw": string(line)}
		}
		entries = append(entries, entry)
	}

	if err := scanner.Err(); err != nil {
		return nil, err
	}

	return entries, nil
}

// LoadFile is a utility to load logs from a file path.
func LoadFile(path string) ([]LogEntry, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer f.Close()
	return LoadLogs(f)
}
