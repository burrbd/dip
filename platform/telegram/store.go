// Package telegram adapts incoming Telegram Bot API updates into bot.Command
// values, and sends bot responses back to Telegram chats.
package telegram

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"
	"sync"
)

// Store is an append-only JSONL file store. Each key maps to a separate file
// at {dir}/{key}.jsonl. Methods are safe for concurrent use.
type Store struct {
	dir        string
	mu         sync.Mutex
	fprintlnFn func(w io.Writer, line string) error // injectable; defaults to defaultFprintln
	openFn     func(name string) (*os.File, error)   // injectable; defaults to os.Open
}

func defaultFprintln(w io.Writer, line string) error {
	_, err := fmt.Fprintln(w, line)
	return err
}

// NewStore returns a Store backed by dir, creating the directory if necessary.
func NewStore(dir string) (*Store, error) {
	if err := os.MkdirAll(dir, 0o755); err != nil {
		return nil, fmt.Errorf("telegram: create store dir: %w", err)
	}
	return &Store{dir: dir, fprintlnFn: defaultFprintln, openFn: os.Open}, nil
}

func (s *Store) filePath(key string) string {
	return filepath.Join(s.dir, key+".jsonl")
}

// Append writes line as a new entry in the file for key.
func (s *Store) Append(key, line string) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	f, err := os.OpenFile(s.filePath(key), os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0o644)
	if err != nil {
		return fmt.Errorf("telegram: open store file: %w", err)
	}
	defer f.Close()
	if err := s.fprintlnFn(f, line); err != nil {
		return fmt.Errorf("telegram: write store file: %w", err)
	}
	return nil
}

// ReadAll returns all non-empty lines for key, in order. If the file does not
// exist, nil is returned with no error.
func (s *Store) ReadAll(key string) ([]string, error) {
	s.mu.Lock()
	defer s.mu.Unlock()
	f, err := s.openFn(s.filePath(key))
	if os.IsNotExist(err) {
		return nil, nil
	}
	if err != nil {
		return nil, fmt.Errorf("telegram: open store file: %w", err)
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
