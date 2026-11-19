package cli

import (
	"fmt"
	"time"

	"github.com/matchbull6/go_todo/internal/store"
	"github.com/matchbull6/go_todo/internal/task"
	"github.com/spf13/cobra"
)

func parsePriorityFlag(s string) (task.Priority, error) {
	if s == "" {
		return task.None, nil
	}
	return task.ParsePriority(s)
}

func newAddCmd(opts *options) *cobra.Command {
	var priority string
	cmd := &cobra.Command{
		Use:   "add <title>",
		Short: "Add a task",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			title := args[0]
			prio, err := parsePriorityFlag(priority)
			if err != nil {
				return err
			}
			var t task.Task
			err = store.Update(opts.file, func(l *task.List) (err error) {
				t, err = l.Add(title, prio, time.Now())
				return err
			})
			if err != nil {
				return err
			}
			fmt.Fprintf(cmd.OutOrStdout(), "Added task %d: %s\n", t.ID, t.Title)
			return nil
		},
	}
	cmd.Flags().StringVarP(&priority, "priority", "p", "", "low, medium or high")
	return cmd
}
