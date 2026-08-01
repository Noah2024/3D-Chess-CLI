// Copyright © 2026 Noah Yurasko distributed under GNU GENERAL PUBLIC LICENSE V3

// moves Provides CLI interface for accessing the debug movement function
package moves

import (
	"3DC/config"
	"3DC/internal/game/load"
	"3DC/internal/move"
	"3DC/internal/move/genMoves"

	"3DC/util/bitutil"
	"fmt"

	"github.com/spf13/cobra"
)

// Re-executes algorithms contained in the genMoves.go file DOES NOT handle checking.
// But will handle all over movmenet validation rules
func DebugMoveGen() *cobra.Command {
	return &cobra.Command{
		Use:   "moves",
		Short: "Generates all valid moves for a given piece ",
		Long:  "Re-executes algorithms contained in the genMoves.go file DOES NOT handle checking",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {

			allLoadedPieces, _ := load.LoadGame(config.CurrentGame)

			uintLocFrom, _, _, _ := move.ParseLoc(args[0])
			visFrom, _ := move.PieceType(allLoadedPieces, uintLocFrom)

			bs, _ := load.GenerateBoardState(allLoadedPieces, visFrom)

			x, y, z := bitutil.UintToVec(uintLocFrom)
			moves := genMoves.MoveMap[visFrom](bs, uintLocFrom, x, y, z)

			fmt.Printf("Generating moves for %s @ %s\n", visFrom, args[0])
			fmt.Printf("%064b\n", moves)
			fmt.Println(moves)
			return nil
		},
	}
}
