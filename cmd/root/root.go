package root

import (
	"3DC/util/logger"

	"3DC/cmd/board"

	"3DC/cmd/game"

	"github.com/spf13/cobra"
)

func RootCommand() *cobra.Command {
	rootCmd := &cobra.Command{
		Use:   "3DC",
		Short: "Root command for 3d chess",
		Long:  "3DC is exactly as the name implies a CLI application for running 3D games of chess",
		RunE: func(cmd *cobra.Command, args []string) error {
			logger.Debug("Init root command")
			return nil
		},
	}
	rootCmd.AddCommand(board.Board())
	rootCmd.AddCommand(game.GameCommand())
	return rootCmd
}
