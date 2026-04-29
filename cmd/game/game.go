package game

import (
	"3DC/cmd/game/list"
	"3DC/cmd/game/load"
	"3DC/cmd/game/new"
	"3DC/cmd/game/save"
	"3DC/util/logger"

	"github.com/spf13/cobra"
)

func GameCommand() *cobra.Command {
	gameCmd := &cobra.Command{
		Use:   "game",
		Short: "Create and manage games",
		Long:  "Create and manage games",
		RunE: func(cmd *cobra.Command, args []string) error {
			logger.Debug("Game root cmd")
			return nil
		},
	}
	gameCmd.AddCommand(new.NewCommand())
	gameCmd.AddCommand(list.ListCommand())
	gameCmd.AddCommand(save.SaveCommand())
	gameCmd.AddCommand(load.LoadCommand())
	return gameCmd
}
