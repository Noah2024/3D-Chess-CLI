// Copyright © 2026 Noah Yurasko distributed under GNU GENERAL PUBLIC LICENSE V3

// save Provides CLI interface for saving the current game
package save

import (
	"3DC/config"
	"3DC/internal/game/load"
	"3DC/internal/game/save"
	"3DC/util/logger"
	"path/filepath"

	"github.com/spf13/cobra"
)

// SaveCommand returns a cobra command that saves the current game to the data/games folder.
// It takes a single argument, which is the name of the game to save.
func SaveCommand() *cobra.Command {
	return &cobra.Command{
		Use:   "save",
		Short: "save current game to data/games folder",
		Long:  "Takes single name arugment",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			//Ik ik I really need to do better input validation here
			//But thats a later me problem
			game, meta, err := load.LoadGame(config.CurrentGame)
			if err != nil {
				logger.Error("Could not load game in cobra save command")
			}
			gameLoc := filepath.Join(config.DataDir, args[0])
			save.SaveGame(game, meta, gameLoc)
			return nil
		},
	}
}
