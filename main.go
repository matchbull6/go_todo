// Command go_todo is a terminal task manager.
package main

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strconv"
	"time"

	"github.com/matchbull6/go_todo/internal/store"
	"github.com/matchbull6/go_todo/internal/task"
	"github.com/spf13/cobra"
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

func newRootCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:           "go_todo",
		Short:         "A terminal task manager",
		SilenceUsage:  true,
		SilenceErrors: true,
	}
	cmd.AddCommand(newAddCmd())
	cmd.AddCommand(newDoneCmd())
	cmd.AddCommand(newRmCmd())
	return cmd
}

func newAddCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "add <title>",
		Short: "Add a task",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			title := args[0]
			path, err := defaultPath()
			if err != nil {
				return err
			}
			l, err := store.Load(path)
			if err != nil {
				return err
			}
			t, err := l.Add(title, task.None, time.Now())
			if err != nil {
				return err
			}
			if err := store.Save(path, l); err != nil {
				return err
			}
			fmt.Fprintf(cmd.OutOrStdout(), "Added task %d: %s\n", t.ID, t.Title)
			return nil
		},
	}
}

func newDoneCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "done <id>",
		Short: "Mark a task as done",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			id, err := strconv.Atoi(args[0])
			if err != nil {
				return fmt.Errorf("invalid task id %q", args[0])
			}
			path, err := defaultPath()
			if err != nil {
				return err
			}
			l, err := store.Load(path)
			if err != nil {
				return err
			}
			t, err := l.Complete(id, time.Now())
			if err != nil {
				return err
			}
			if err := store.Save(path, l); err != nil {
				return err
			}
			fmt.Fprintf(cmd.OutOrStdout(), "Completed task %d: %s\n", t.ID, t.Title)
			return nil
		},
	}
}

func newRmCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "rm <id>",
		Short: "Delete a task",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			id, err := strconv.Atoi(args[0])
			if err != nil {
				return fmt.Errorf("invalid task id %q", args[0])
			}
			path, err := defaultPath()
			if err != nil {
				return err
			}
			l, err := store.Load(path)
			if err != nil {
				return err
			}
			t, err := l.Delete(id)
			if err != nil {
				return err
			}
			if err := store.Save(path, l); err != nil {
				return err
			}
			fmt.Fprintf(cmd.OutOrStdout(), "Deleted task %d: %s\n", t.ID, t.Title)
			return nil
		},
	}
}

func defaultPath() (string, error) {
	dir, err := os.UserHomeDir()
	if err != nil {
		return "", err
	}
	return filepath.Join(dir, ".go_todo.json"), nil
}
