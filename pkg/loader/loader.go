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

// LoadLogs reads NDJSON from the given reader and returns all raw lines,
// and a sample of parsed entries for column detection.
func LoadLogs(r io.Reader) ([]string, []LogEntry, error) {
	var lines []string
	var sample []LogEntry
	scanner := bufio.NewScanner(r)

	// Increase buffer size
	buf := make([]byte, 0, 64*1024)
	scanner.Buffer(buf, 10*1024*1024)

	count := 0
	for scanner.Scan() {
		line := scanner.Bytes()
		if len(line) == 0 {
			continue
		}

		// Store raw line
		strLine := string(line)
		lines = append(lines, strLine)

		// Parse sample (first 100 lines)
		if count < 100 {
			var entry LogEntry
			if err := json.Unmarshal(line, &entry); err == nil {
				sample = append(sample, entry)
			}
		}
		count++
	}

	if err := scanner.Err(); err != nil {
		return nil, nil, err
	}

	return lines, sample, nil
}

// LoadFile is a utility to load logs from a file path.
func LoadFile(path string) ([]string, []LogEntry, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, nil, err
	}
	defer f.Close()
	return LoadLogs(f)
}
