package save

import (
	"3DC/config"
	"3DC/internal/game/load"
	"3DC/internal/game/save"
	"3DC/util/must"
	"path/filepath"

	"github.com/spf13/cobra"
)

func SaveCommand() *cobra.Command {
	return &cobra.Command{
		Use:   "save",
		Short: "save current game to data/games folder",
		Long:  "Takes single name arugment",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			//Ik ik I really need to do better input validation here
			//But thats a later me problem
			game := must.Must(load.LoadGame(config.CurrentGame))
			gameLoc := filepath.Join(config.DataDir, args[0])
			save.SaveGame(game, gameLoc)
			return nil
		},
	}
}
