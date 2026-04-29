package delete

import (
	"3DC/internal/game/delete"

	"github.com/spf13/cobra"
)

func DeleteCommand() *cobra.Command {
	return &cobra.Command{
		Use:   "del",
		Short: "Delete a saved game",
		Long:  "Takes 1 argument; the name of the game or its number; use 'game list' to see a list of commands",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			//Again I should prolly do some more input validation
			delete.DeleteGame(args[0])
			return nil
		},
	}
}
