// Package cli implements the go_todo commands.
package cli

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strconv"
)

// Run returns the exit code instead of calling os.Exit, for tests.
func Run(args []string, stdout, stderr io.Writer) int {
	cmd := newRootCmd()
	cmd.SetArgs(args)
	cmd.SetOut(stdout)
	cmd.SetErr(stderr)
	if err := cmd.Execute(); err != nil {
		fmt.Fprintln(stderr, "Error:", err)
		return 1
	}
	return 0
}

// options are shared by all commands; file is set from --file
// and resolved to the final path before any command runs.
type options struct {
	file string
}

func resolveTaskFile(flag string) (string, error) {
	if flag != "" {
		return flag, nil
	}
	if env := os.Getenv("TODO_FILE"); env != "" {
		return env, nil
	}
	home, err := os.UserHomeDir()
	if err != nil {
		return "", err
	}
	return filepath.Join(home, ".go_todo.json"), nil
}

func parseID(s string) (int, error) {
	id, err := strconv.Atoi(s)
	if err != nil {
		return 0, fmt.Errorf("invalid task id %q", s)
	}
	return id, nil
}
