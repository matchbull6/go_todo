package store_test

import (
	"errors"
	"os"
	"path/filepath"
	"reflect"
	"strings"
	"testing"
	"time"

	"github.com/matchbull6/go_todo/internal/store"
	"github.com/matchbull6/go_todo/internal/task"
)

func TestLoad_Missing(t *testing.T) {
	l, err := store.Load(filepath.Join(t.TempDir(), "none.json"))
	if err != nil {
		t.Fatalf("Load: %v", err)
	}
	if l.LastID != 0 || len(l.Tasks) != 0 {
		t.Errorf("want empty list, got %+v", l)
	}
}

func TestLoad_BadJSON(t *testing.T) {
	path := filepath.Join(t.TempDir(), "bad.json")
	if err := os.WriteFile(path, []byte("{not json"), 0o600); err != nil {
		t.Fatal(err)
	}
	_, err := store.Load(path)
	if err == nil || !strings.Contains(err.Error(), path) {
		t.Errorf("err = %v, want error mentioning %s", err, path)
	}
}

func TestSaveLoad(t *testing.T) {
	path := filepath.Join(t.TempDir(), "tasks.json")
	now := time.Date(2026, 1, 2, 3, 4, 5, 0, time.UTC)

	var want task.List
	want.Add("a", task.High, now)
	want.Add("b", task.None, now)
	want.Complete(2, now)

	if err := store.Save(path, want); err != nil {
		t.Fatalf("Save: %v", err)
	}
	got, err := store.Load(path)
	if err != nil {
		t.Fatalf("Load: %v", err)
	}
	if !reflect.DeepEqual(got, want) {
		t.Errorf("got %+v\nwant %+v", got, want)
	}
}

func TestSave_NoTempLeft(t *testing.T) {
	dir := t.TempDir()
	if err := store.Save(filepath.Join(dir, "tasks.json"), task.List{}); err != nil {
		t.Fatalf("Save: %v", err)
	}
	entries, _ := os.ReadDir(dir)
	if len(entries) != 1 {
		t.Errorf("files in dir = %d, want 1 (temp file left?)", len(entries))
	}
}

func TestSave_ErrorShowsUserPath(t *testing.T) {
	path := filepath.Join(t.TempDir(), "missing-dir", "tasks.json")
	err := store.Save(path, task.List{})
	if err == nil {
		t.Fatal("want error for missing directory")
	}
	if !strings.Contains(err.Error(), path) || strings.Contains(err.Error(), ".tmp") {
		t.Errorf("err = %v, want user path, not temp file", err)
	}
}

func TestUpdate(t *testing.T) {
	path := filepath.Join(t.TempDir(), "tasks.json")
	add := func(l *task.List) error {
		_, err := l.Add("a", task.None, time.Now())
		return err
	}
	for range 2 {
		if err := store.Update(path, add); err != nil {
			t.Fatalf("Update: %v", err)
		}
	}
	l, _ := store.Load(path)
	if len(l.Tasks) != 2 || l.LastID != 2 {
		t.Errorf("got %+v, want 2 tasks", l)
	}
}

func TestUpdate_FnErrorDoesNotSave(t *testing.T) {
	path := filepath.Join(t.TempDir(), "tasks.json")
	errBoom := errors.New("boom")

	err := store.Update(path, func(l *task.List) error {
		l.Add("a", task.None, time.Now())
		return errBoom
	})
	if !errors.Is(err, errBoom) {
		t.Errorf("err = %v, want errBoom", err)
	}
	if _, err := os.Stat(path); !errors.Is(err, os.ErrNotExist) {
		t.Errorf("file was written: %v", err)
	}
}
