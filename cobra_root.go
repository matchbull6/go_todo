package main

import "github.com/spf13/cobra"

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
	cmd.AddCommand(newListCmd())
	return cmd
}
