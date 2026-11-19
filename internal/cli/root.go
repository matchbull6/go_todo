package cli

import "github.com/spf13/cobra"

func newRootCmd() *cobra.Command {
	opts := &options{}
	cmd := &cobra.Command{
		Use:           "go_todo",
		Short:         "A terminal task manager",
		SilenceUsage:  true,
		SilenceErrors: true,
		PersistentPreRunE: func(cmd *cobra.Command, args []string) (err error) {
			opts.file, err = resolveTaskFile(opts.file)
			return err
		},
	}
	cmd.PersistentFlags().StringVarP(&opts.file, "file", "f", "",
		"task file (default $TODO_FILE or ~/.go_todo.json)")
	cmd.AddCommand(newAddCmd(opts), newDoneCmd(opts), newRmCmd(opts), newListCmd(opts))
	return cmd
}
