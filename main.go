// Command go_todo is a terminal task manager.
package main

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strconv"
)

func main() {
	os.Exit(run(os.Args[1:], os.Stdout, os.Stderr))
}

// run returns the exit code instead of calling os.Exit, for tests.
func run(args []string, stdout, stderr io.Writer) int {
	// TODO(stage 3): --file -> TODO_FILE -> ~/.go_todo.json.
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

func defaultTaskFile() (string, error) {
	dir, err := os.UserHomeDir()
	if err != nil {
		return "", err
	}
	return filepath.Join(dir, ".go_todo.json"), nil
}

func parseID(s string) (int, error) {
	id, err := strconv.Atoi(s)
	if err != nil {
		return 0, fmt.Errorf("invalid task id %q", s)
	}
	return id, nil
}
