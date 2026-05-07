package load

import (
	"3DC/config"
	"3DC/internal/game/load"
	"3DC/internal/game/save"
	"3DC/util/dialog"
	"3DC/util/must"
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

func LoadCommand() *cobra.Command {
	return &cobra.Command{
		Use:   "load",
		Short: "loads the provided game into the data/CurrentGame folder",
		Long:  "Takes single name arugment (use 'game list' to see a list of games to load)",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			//Ik ik I really need to do better input validation here
			//But thats a later me problem

			//Checks if file exists
			if _, err := os.Stat(config.DataDir); err == nil {
				if !dialog.Confirm("Are you sure you want to overwrite your current game?") {
					return nil
				}
				return err
			}
			game := must.Must(load.LoadGame(config.DataDir))
			save.SaveGame(game, args[0])
			fmt.Println("Saved game loaded succesfully")
			return nil
		},
	}
}
