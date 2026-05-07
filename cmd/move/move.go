package move

import (
	"3DC/internal/move"

	"github.com/spf13/cobra"
)

func MoveCommand() *cobra.Command {
	return &cobra.Command{
		Use:   "move",
		Short: "Move a pice to a given location",
		Long:  "Takes two arguments, location of piece to be moved, location of where to move it",
		Args:  cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			move.Move(args[0], args[1])
			return nil
		},
	}
}
