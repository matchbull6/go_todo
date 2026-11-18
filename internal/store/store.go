// Package store keeps a task list in a JSON file.
package store

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"

	"github.com/matchbull6/go_todo/internal/task"
)

// Load returns an empty list if the file does not exist.
func Load(path string) (task.List, error) {
	data, err := os.ReadFile(path)
	if errors.Is(err, fs.ErrNotExist) {
		return task.List{}, nil
	}
	if err != nil {
		return task.List{}, fmt.Errorf("store: %w", err)
	}

	var l task.List
	if err := json.Unmarshal(data, &l); err != nil {
		return task.List{}, fmt.Errorf("store: %s: %w", path, err)
	}
	return l, nil
}

// Save writes atomically (temp file + rename), mode 0600.
func Save(path string, l task.List) (err error) {
	defer func() {
		if err != nil {
			err = fmt.Errorf("store: %w", err)
		}
	}()

	data, err := json.MarshalIndent(l, "", "  ")
	if err != nil {
		return err
	}
	data = append(data, '\n')

	tmp, err := os.CreateTemp(filepath.Dir(path), ".go_todo-*.tmp")
	if err != nil {
		return err
	}
	defer func() {
		tmp.Close()
		os.Remove(tmp.Name())
	}()
	if _, err := tmp.Write(data); err != nil {
		return err
	}
	if err := tmp.Sync(); err != nil {
		return err
	}
	if err := tmp.Close(); err != nil {
		return err
	}
	return os.Rename(tmp.Name(), path)
}
