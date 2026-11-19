package cli

import (
	"bytes"
	"path/filepath"
	"strings"
	"testing"
)

// run executes go_todo with a fresh task file per test.
func run(t *testing.T, file string, args ...string) (stdout, stderr string, code int) {
	t.Helper()
	var out, errOut bytes.Buffer
	code = Run(append([]string{"--file", file}, args...), &out, &errOut)
	return out.String(), errOut.String(), code
}

func TestCommands(t *testing.T) {
	file := filepath.Join(t.TempDir(), "tasks.json")

	steps := []struct {
		args     []string
		wantOut  []string
		dontWant []string
	}{
		{args: []string{"list"}, wantOut: []string{"No tasks"}},
		{args: []string{"add", "buy milk"}, wantOut: []string{"Added task 1: buy milk"}},
		{args: []string{"add", "fix bug", "-p", "high"}, wantOut: []string{"Added task 2"}},
		{args: []string{"list"}, wantOut: []string{"ID", "buy milk", "fix bug", "high"}},
		{args: []string{"done", "1"}, wantOut: []string{"Completed task 1"}},
		{args: []string{"list"}, wantOut: []string{"fix bug"}, dontWant: []string{"buy milk"}},
		{args: []string{"list", "--all"}, wantOut: []string{"DONE", "✓", "buy milk"}},
		{args: []string{"rm", "2"}, wantOut: []string{"Deleted task 2"}},
		{args: []string{"add", "next"}, wantOut: []string{"Added task 3"}},
	}
	for _, s := range steps {
		out, errOut, code := run(t, file, s.args...)
		name := strings.Join(s.args, " ")
		if code != 0 {
			t.Fatalf("%s: code %d, want 0; stderr: %s", name, code, errOut)
		}
		for _, w := range s.wantOut {
			if !strings.Contains(out, w) {
				t.Errorf("%s: output missing %q:\n%s", name, w, out)
			}
		}
		for _, w := range s.dontWant {
			if strings.Contains(out, w) {
				t.Errorf("%s: output has %q:\n%s", name, w, out)
			}
		}
	}
}

func TestErrors(t *testing.T) {
	tests := []struct {
		name    string
		args    []string
		wantErr string
	}{
		{"bad id", []string{"done", "abc"}, `invalid task id "abc"`},
		{"missing task", []string{"rm", "42"}, "not found"},
		{"empty title", []string{"add", "  "}, "empty title"},
		{"bad priority", []string{"add", "x", "-p", "urgent"}, "invalid priority"},
		{"no args", []string{"add"}, "accepts 1 arg"},
		{"unknown command", []string{"nope"}, "unknown command"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			file := filepath.Join(t.TempDir(), "tasks.json")
			_, errOut, code := run(t, file, tt.args...)
			if code != 1 {
				t.Errorf("code = %d, want 1", code)
			}
			if !strings.Contains(errOut, tt.wantErr) {
				t.Errorf("stderr = %q, want %q", errOut, tt.wantErr)
			}
		})
	}
}

func TestResolveTaskFile(t *testing.T) {
	home := t.TempDir()
	t.Setenv("HOME", home)
	t.Setenv("USERPROFILE", home)

	t.Setenv("TODO_FILE", "")
	if got, _ := resolveTaskFile(""); got != filepath.Join(home, ".go_todo.json") {
		t.Errorf("default = %q", got)
	}

	t.Setenv("TODO_FILE", "env.json")
	if got, _ := resolveTaskFile(""); got != "env.json" {
		t.Errorf("env = %q, want env.json", got)
	}
	if got, _ := resolveTaskFile("flag.json"); got != "flag.json" {
		t.Errorf("flag = %q, want flag.json", got)
	}
}
