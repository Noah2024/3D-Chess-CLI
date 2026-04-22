package ls

import (
	"fmt"

	"github.com/spf13/cobra"
)

func LsCommand() *cobra.Command {
	return &cobra.Command{
		Use:   "ls",
		Short: "View current slice of baord",
		Long:  "view current slice of board",
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Fprintf(cmd.OutOrStdout(), "Testing args %s", args[0])
		},
	}
}
