// Command go_todo is a terminal task manager.
package main

import (
	"fmt"
	"io"
	"os"
)

func main() {
	os.Exit(run(os.Args[1:], os.Stdout, os.Stderr))
}

// run returns the exit code instead of calling os.Exit, for tests.
func run(args []string, stdout, stderr io.Writer) int {
	// TODO(stage 2): cobra root command.
	// TODO(stage 3): --file -> TODO_FILE -> ~/.go_todo.json.
	_, _, _ = args, stdout, stderr
	fmt.Fprintln(stderr, "not implemented yet")
	return 1
}
