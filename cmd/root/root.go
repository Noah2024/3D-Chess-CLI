package root

import (
	"3DChessCLI/cmd/ls"
	"3DChessCLI/util/logger"

	"github.com/spf13/cobra"
)

func RootCommand() *cobra.Command {
	rootCmd := &cobra.Command{
		Use:   "3DChessCLI",
		Short: "Root command for 3d chess",
		Long:  "3DChessCLI is exactly as the name implies a CLI application for running 3D games of chess",
		Run: func(cmd *cobra.Command, args []string) {
			logger.Debug("RootCmd Registered")
		},
	}
	rootCmd.AddCommand(ls.LsCommand())
	return rootCmd
}
