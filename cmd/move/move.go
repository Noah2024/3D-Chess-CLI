// Copyright © 2026 Noah Yurasko distributed under GNU GENERAL PUBLIC LICENSE V3

// move Provides CLI interface for moving pieces in the game
package move

import (
	"3DC/internal/move"

	"github.com/spf13/cobra"
)

// MoveCommand returns a cobra command that moves a piece from one location to another in the game.
// It takes two arguments: the current location of the piece and the target location to move it to.
//
//	Expected Arugment format
//	"a1A"
//
// 'a' refers to the column,
// '1' refers to the row,
// 'A' refers to the level of the piece in the 3D board.
func MoveCommand() *cobra.Command {
	return &cobra.Command{
		Use:   "move",
		Short: "Move a pice to a given location expectes two arguments like 'a1A'",
		Long:  "Takes two arguments, location of piece to be moved, location of where to move it. Expected format is 'a1A' where 'a' is the column, '1' is the row, and 'A' is the level of the piece in the 3D board.",
		Args:  cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			move.MoveCommand(args[0], args[1])
			return nil
		},
	}
}
