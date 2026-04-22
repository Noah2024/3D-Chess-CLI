package ls

import (
	"fmt"

	"github.com/spf13/cobra"
)

func lsCommand() *cobra.Command {
	lsCommand := &cobra.Command{
		Use:   "ls",
		Short: "View current slice of baord",
		Long:  "view current slice of board",
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Println("bruh")
		},
	}
	return lsCommand
}
