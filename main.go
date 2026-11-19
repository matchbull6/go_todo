// Command go_todo is a terminal task manager.
package main

import (
	"os"

	"github.com/matchbull6/go_todo/internal/cli"
)

func main() {
	os.Exit(cli.Run(os.Args[1:], os.Stdout, os.Stderr))
}
