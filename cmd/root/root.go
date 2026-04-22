package root

import (
	"fmt"

	"github.com/spf13/cobra"

	"3DChessCLI/cmd/ls"
)

func RootCommand() *cobra.Command {
	rootCmd := &cobra.Command{
		Use:   "3DChessCLI",
		Short: "Root command for 3d chess",
		Long:  "3DChessCLI is exactly as the name implies a CLI application for running 3D games of chess",
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Fprintf(cmd.OutOrStdout(), "bruh")
		},
	}
	rootCmd.AddCommand(ls.LsCommand())
	return rootCmd
}
