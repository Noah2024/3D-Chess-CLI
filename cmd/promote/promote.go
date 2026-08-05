// Copyright © 2026 Noah Yurasko distributed under GNU GENERAL PUBLIC LICENSE V3

// move Provides CLI interface for promoting pawns in the game
package promote

import (
	"3DC/internal/promote"

	"github.com/spf13/cobra"
)

func PromoteCommand() *cobra.Command {
	return &cobra.Command{
		Use:   "promote",
		Short: "Promote a pawn takes two arguments, team of promotion & the target piece to promote to",
		Long:  "Takes two arguments, location of piece to be moved, location of where to move it. Expected format is 'a1A' where 'a' is the column, '1' is the row, and 'A' is the level of the piece in the 3D board.",
		Args:  cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			promote.AttemptPawnPromotion(args[0], args[1])
			return nil
		},
	}
}
