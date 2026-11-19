package main

import (
	"fmt"

	"github.com/matchbull6/go_todo/internal/store"
	"github.com/matchbull6/go_todo/internal/task"
	"github.com/spf13/cobra"
)

func newRmCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "rm <id>",
		Short: "Delete a task",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			id, err := parseID(args[0])
			if err != nil {
				return err
			}
			path, err := defaultTaskFile()
			if err != nil {
				return err
			}
			var t task.Task
			err = store.Update(path, func(l *task.List) (err error) {
				t, err = l.Delete(id)
				return err
			})
			if err != nil {
				return err
			}
			fmt.Fprintf(cmd.OutOrStdout(), "Deleted task %d: %s\n", t.ID, t.Title)
			return nil
		},
	}
}
