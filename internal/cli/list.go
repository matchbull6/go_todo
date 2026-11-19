package cli

import (
	"fmt"
	"text/tabwriter"

	"github.com/matchbull6/go_todo/internal/store"
	"github.com/matchbull6/go_todo/internal/task"
	"github.com/spf13/cobra"
)

const dateLayout = "02 Jan 15:04"

func newListCmd() *cobra.Command {
	var all bool
	cmd := &cobra.Command{
		Use:   "list",
		Short: "List tasks",
		Args:  cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			path, err := defaultTaskFile()
			if err != nil {
				return err
			}
			l, err := store.Load(path)
			if err != nil {
				return err
			}

			var shown []task.Task
			for _, t := range l.Tasks {
				if all || !t.Done {
					shown = append(shown, t)
				}
			}
			if len(shown) == 0 {
				fmt.Fprintln(cmd.OutOrStdout(), "No tasks")
				return nil
			}

			w := tabwriter.NewWriter(cmd.OutOrStdout(), 0, 0, 2, ' ', 0)
			if all {
				fmt.Fprint(w, "DONE\t")
			}
			fmt.Fprintln(w, "ID\tPRIORITY\tCREATED\tTITLE")
			for _, t := range shown {
				if all {
					done := ""
					if t.Done {
						done = "✓"
					}
					fmt.Fprintf(w, "%s\t", done)
				}
				fmt.Fprintf(w, "%d\t%s\t%s\t%s\n",
					t.ID, t.Priority, t.CreatedAt.Format(dateLayout), t.Title)
			}
			return w.Flush()
		},
	}
	cmd.Flags().BoolVar(&all, "all", false, "include done tasks")
	return cmd
}
