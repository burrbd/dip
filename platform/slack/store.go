package slack

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"
	"sync"
)

// store is an append-only JSONL file store. Each key maps to a separate file
// at {dir}/{key}.jsonl. Methods are safe for concurrent use.
type store struct {
	dir        string
	mu         sync.Mutex
	fprintlnFn func(w io.Writer, line string) error
	openFn     func(name string) (*os.File, error)
}

func defaultFprintln(w io.Writer, line string) error {
	_, err := fmt.Fprintln(w, line)
	return err
}

// newStore returns a store backed by dir, creating the directory if necessary.
func newStore(dir string) (*store, error) {
	if err := os.MkdirAll(dir, 0o755); err != nil {
		return nil, fmt.Errorf("slack: create store dir: %w", err)
	}
	return &store{dir: dir, fprintlnFn: defaultFprintln, openFn: os.Open}, nil
}

func (s *store) filePath(key string) string {
	return filepath.Join(s.dir, key+".jsonl")
}

// append writes line as a new entry in the file for key.
func (s *store) append(key, line string) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	f, err := os.OpenFile(s.filePath(key), os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0o644)
	if err != nil {
		return fmt.Errorf("slack: open store file: %w", err)
	}
	defer f.Close()
	if err := s.fprintlnFn(f, line); err != nil {
		return fmt.Errorf("slack: write store file: %w", err)
	}
	return nil
}

// readAll returns all non-empty lines for key, in order. If the file does not
// exist, nil is returned with no error.
func (s *store) readAll(key string) ([]string, error) {
	s.mu.Lock()
	defer s.mu.Unlock()
	f, err := s.openFn(s.filePath(key))
	if os.IsNotExist(err) {
		return nil, nil
	}
	if err != nil {
		return nil, fmt.Errorf("slack: open store file: %w", err)
	}
	defer f.Close()
	var lines []string
	sc := bufio.NewScanner(f)
	for sc.Scan() {
		if line := strings.TrimSpace(sc.Text()); line != "" {
			lines = append(lines, line)
		}
	}
	return lines, sc.Err()
}
