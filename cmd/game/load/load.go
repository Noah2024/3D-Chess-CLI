// Copyright © 2026 Noah Yurasko distributed under GNU GENERAL PUBLIC LICENSE V3

// load Provides CLI interface for loading a saved game
package load

import (
	"3DC/config"
	"3DC/internal/game/load"
	"3DC/internal/game/save"
	"3DC/util/dialog"
	"fmt"
	"os"
	"path"

	"github.com/spf13/cobra"
)

// LoadCommand returns a cobra command that loads a saved game into the CurrentGame folder. It takes a single argument, which is the name of the game to load.
// If the CurrentGame folder already contains a game, it prompts the user for confirmation before overwriting it.
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
			if _, err := os.Stat(config.CurrentGame); err == nil {
				fmt.Println("HERE ", err)
				if !dialog.Confirm("Are you sure you want to overwrite your current game?") {
					return nil
				}
			} else {
				return err
			}

			game, meta, err := load.LoadGame(path.Join(config.DataDir, args[0]))
			if err != nil {
				fmt.Println("Bro")
			}
			save.SaveGame(game, meta, config.CurrentGame)
			return nil
		},
	}
}
